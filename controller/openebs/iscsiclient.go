package openebs

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/k8s"
	"mayadata.io/openebs-upgrade/types"
	"strings"
)

// getISCSISetupComponentsStatus checks if the components required to setup ISCSI client on nodes
// are running or not, if they are in error, it will throw error. If they are already running which means
// setup is already done, it will mark them for deletion.
func (p *Planner) getISCSISetupComponentsStatus() (bool, error) {
	var (
		isRunning        bool
		desiredDaemonset *unstructured.Unstructured
		desiredConfigmap *unstructured.Unstructured
	)
	for _, component := range p.observedOpenEBSComponents {
		if component.GetKind() == types.KindDaemonSet {
			if component.GetName() == types.OpenEBSNodeSetupDaemonsetNameKey {
				// get the .spec.status field and compare the currentNumberScheduled
				// and desiredNumberScheduled field's values.
				componentStatus, _, err := unstructured.NestedMap(component.Object, "status")
				if err != nil {
					return isRunning, err
				}
				if !(componentStatus == nil || len(componentStatus) == 0) {
					desiredReplicas := componentStatus["desiredNumberScheduled"]
					scheduledReplicas := componentStatus["currentNumberScheduled"]
					readyReplicas := componentStatus["numberReady"]
					// if desired replicas is equal to current replicas is equal to no of ready replicas
					// then we can determine that ISCSI setup has completed or it was already installed
					// and we will go ahead and clean up the daemonset.
					if (desiredReplicas == scheduledReplicas) && (desiredReplicas == readyReplicas) {
						desiredDaemonset = component
					} else {
						return isRunning, errors.Errorf("No of ready replicas: %d for daemonset[name: openebs-node-setup, namespace: %s] is not equal to no of desired replicas: %d",
							readyReplicas, p.ObservedOpenEBS.Namespace, desiredReplicas)
					}
				}
			}
		} else if component.GetKind() == types.KindConfigMap {
			if component.GetName() == types.OpenEBSNodeSetupConfigmapNameKey {
				desiredConfigmap = component
			}
		}
	}
	// add daemonset and configmap to explicit deletes if and only if both are present.
	if desiredDaemonset != nil {
		// update the isSetupDone field in OpenEBS CR since it is a one time process
		// and should not be done again once completed.
		p.ObservedOpenEBS.Spec.PreInstallation.ISCSIClient.IsSetupDone = true
		var openebs *unstructured.Unstructured
		openebsRaw, err := json.Marshal(p.ObservedOpenEBS)
		if err != nil {
			return isRunning,
				errors.Errorf("Error marshalling updated OpenEBS for updating ISCSI fields: %+v", err)
		}
		err = json.Unmarshal(openebsRaw, &openebs)
		if err != nil {
			return isRunning,
				errors.Errorf("Error unmarshalling updated OpenEBS for updating ISCSI fields: %+v", err)
		}
		p.ExplicitUpdates = append(p.ExplicitUpdates, openebs)
		p.ExplicitDeletes = append(p.ExplicitDeletes, desiredDaemonset)
		if desiredConfigmap != nil {
			p.ExplicitDeletes = append(p.ExplicitDeletes, desiredConfigmap)
		}
		isRunning = true
	}
	return isRunning, nil
}

// getDesiredISCSIManifests updates all the ISCSI component's manifest as per the provided
// or the default values.
func (p *Planner) getDesiredISCSIManifests(iscsiManifests map[string]*unstructured.Unstructured) (
	map[string]*unstructured.Unstructured, error) {
	var err error
	for key, value := range iscsiManifests {
		// set the common label i.e., openebs-upgrade.dao.mayadata.io/managed: true
		// here since this label should be present in all the components irrespective
		// of their k8s kind, however some specific labels could be set per component
		// such as openebs-upgrade.dao.mayadata.io/component-name: ndm for NDM components,
		// etc.
		componentLabels := value.GetLabels()
		if componentLabels == nil {
			componentLabels = make(map[string]string, 0)
		} else {
			if _, exist := componentLabels[types.OpenEBSVersionLabelKey]; exist {
				// update the version label as per the given OpenEBS version and imageTagSuffix
				// if given.
				componentLabels[types.OpenEBSVersionLabelKey] =
					p.ObservedOpenEBS.Spec.Version + p.ObservedOpenEBS.Spec.ImageTagSuffix
			}
		}
		componentLabels[types.OpenEBSUpgradeDAOManagedLabelKey] =
			types.OpenEBSUpgradeDAOManagedLabelValue
		value.SetLabels(componentLabels)

		kind := strings.Split(key, "_")[1]
		switch kind {
		case types.KindDaemonSet:
			value, err = p.getDesiredDaemonSet(value)
		case types.KindConfigMap:
			value, err = p.getDesiredConfigmap(value)
		default:
			// Doing nothing if an unknown kind
			continue
		}
		if err != nil {
			return iscsiManifests, errors.Errorf("Error updating ISCSI manifests: %+v", err)
		}
		// update manifest with the updated values
		iscsiManifests[key] = value
	}
	return iscsiManifests, nil
}

// getISCSIInstallationManifest forms the YAML for ISCSI client installation on all the
// desired nodes of a cluster.
func (p *Planner) getISCSIInstallationManifest() (map[string]*unstructured.Unstructured, error) {
	componentsYAMLMap := make(map[string]*unstructured.Unstructured)
	var yamlFile string
	// get the OS image running on the underlying node
	osImage, err := k8s.GetOSImage()
	if err != nil {
		return componentsYAMLMap, errors.Errorf("Error getting OS image of a node, error: %+v", err)
	}
	osImageInLowercase := strings.ToLower(osImage)
	switch true {
	case strings.Contains(osImageInLowercase, "ubuntu"):
		yamlFile = "/templates/iscsi-ubuntu-setup.yaml"
	case strings.Contains(osImageInLowercase, strings.ToLower("Red Hat Enterprise Linux")) ||
		strings.Contains(osImageInLowercase, "centos") ||
		strings.Contains(osImageInLowercase, "amazon linux"):
		yamlFile = "/templates/iscsi-amazonlinux-setup.yaml"
	default:
		glog.V(3).Infof("ISCSI installation is not yet supported for %s.", osImage)
		return componentsYAMLMap, nil
	}
	iscsiYaml, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return componentsYAMLMap, errors.New("Error reading ISCSI installation YAML file.")
	}

	// form the mapping from component's "name_kind" as key to YAML
	// string as value using ISCSI yaml.
	componentsYAML := strings.Split(string(iscsiYaml), "---")
	for _, componentYAML := range componentsYAML {
		if componentYAML == "" {
			continue
		}
		unstructuredYAML := unstructured.Unstructured{}
		if err = yaml.Unmarshal([]byte(componentYAML), &unstructuredYAML.Object); err != nil {
			return componentsYAMLMap, errors.Errorf("Error unmarshalling YAML string:%s, Error: %+v", componentYAML, err)
		}
		kind := unstructuredYAML.GetKind()
		name := unstructuredYAML.GetName()

		// Form the key using component's Name and kind separated
		// by underscore
		keyForStoringYaml := name + "_" + kind
		// Store the latest yaml of each component in a map where the key
		// is componentName_kind
		componentsYAMLMap[keyForStoringYaml] = &unstructuredYAML
	}
	// get the desired ISCSI manifest with all the labels and values provided.
	componentsYAMLMap, err = p.getDesiredISCSIManifests(componentsYAMLMap)
	if err != nil {
		return componentsYAMLMap, err
	}
	return componentsYAMLMap, nil
}
