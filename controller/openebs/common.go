/*
Copyright 2020 The MayaData Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    https://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package openebs

import (
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/types"
	"mayadata.io/openebs-upgrade/unstruct"
)

// setDefaultImagePullPolicyIfNotSet sets the default imagePullPolicy
// to "IfNotPresent" for all the components.
// TODO: See if this is required component wise and not at the global
// level.
func (p *Planner) setDefaultImagePullPolicyIfNotSet() error {
	if p.ObservedOpenEBS.Spec.ImagePullPolicy == "" {
		p.ObservedOpenEBS.Spec.ImagePullPolicy = "IfNotPresent"
	}
	return nil
}

// setDefaultStoragePathIfNotSet sets the default storage path for
// OpenEBS to "/var/openebs" if not already set.
func (p *Planner) setDefaultStoragePathIfNotSet() error {
	if p.ObservedOpenEBS.Spec.DefaultStoragePath == "" {
		p.ObservedOpenEBS.Spec.DefaultStoragePath = "/var/openebs"
	} else if strings.HasSuffix(p.ObservedOpenEBS.Spec.DefaultStoragePath, "/") {
		p.ObservedOpenEBS.Spec.DefaultStoragePath = strings.TrimRight(
			p.ObservedOpenEBS.Spec.DefaultStoragePath, "/")
	}
	return nil
}

// setDefaultImagePrefixIfNotSet sets the default registry prefix for
// all the container images if not already set.
// It also checks if the given image registry ends with a forward slash
// or not, if not then it adds one.
func (p *Planner) setDefaultImagePrefixIfNotSet() error {
	if p.ObservedOpenEBS.Spec.ImagePrefix == "" {
		// Default docker registry for OpenEBS enterprise installation will
		// be "mayadataio/" while for community edition will be "quay.io/openebs/".
		if strings.Contains(p.ObservedOpenEBS.Spec.Version, "ee") {
			p.ObservedOpenEBS.Spec.ImagePrefix = "mayadataio/"
		} else {
			p.ObservedOpenEBS.Spec.ImagePrefix = "quay.io/openebs/"
		}
	} else if !strings.HasSuffix(p.ObservedOpenEBS.Spec.ImagePrefix, "/") {
		p.ObservedOpenEBS.Spec.ImagePrefix = p.ObservedOpenEBS.Spec.ImagePrefix + "/"
	}
	return nil
}

// setImageTagSuffixIfPresent sets a custom image tag suffix that can be specified
// for pulling the release candidate images for containers such as 1.10.0-RC1, etc.
//
// The value for this field can be RC1, RC2, etc which will be appended
// to the given OpenEBS version.
// For example, if version is 1.10.0 and the value of imageTagSuffix is RC1,
// the images that will be used for configurable OpenEBS components will be
// 1.10.0-RC1.
func (p *Planner) setImageTagSuffixIfPresent() error {
	if p.ObservedOpenEBS.Spec.ImageTagSuffix != "" {
		// prepend the imageTagSuffix with a hyphen(-) which will be used to append
		// the suffix to the given OpenEBS version.
		// For example if version is 1.10.0 and imageTagSuffix is RC1 then the resultant
		// image will be 1.10.0-RC1.
		p.ObservedOpenEBS.Spec.ImageTagSuffix = "-" + p.ObservedOpenEBS.Spec.ImageTagSuffix
	}
	return nil
}

// setDefaultStorageConfigIfNotSet sets the defaultStorageConfig value
// to "true" if not already set.
func (p *Planner) setDefaultStorageConfigIfNotSet() error {
	if p.ObservedOpenEBS.Spec.CreateDefaultStorageConfig == nil {
		p.ObservedOpenEBS.Spec.CreateDefaultStorageConfig = new(bool)
		*p.ObservedOpenEBS.Spec.CreateDefaultStorageConfig = true
	}
	return nil
}

// BasicComponentDetails stores only the component's kind and name
type BasicComponentDetails struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

// getManifests returns a mapping of component's "name_kind" to YAML of
// the respective components based on a particular version.
// Note: This method makes use of the various operator YAMLs to form this
// mapping.
func (p *Planner) getManifests() error {
	componentsYAMLMap := make(map[string]*unstructured.Unstructured)
	var yamlFile string

	switch p.ObservedOpenEBS.Spec.Version {
	case types.OpenEBSVersion150:
		yamlFile = "/templates/openebs-operator-1.5.0.yaml"
	case types.OpenEBSVersion160:
		yamlFile = "/templates/openebs-operator-1.6.0.yaml"
	case types.OpenEBSVersion170:
		yamlFile = "/templates/openebs-operator-1.7.0.yaml"
	case types.OpenEBSVersion180:
		yamlFile = "/templates/openebs-operator-1.8.0.yaml"
	case types.OpenEBSVersion190:
		yamlFile = "/templates/openebs-operator-1.9.0.yaml"
	case types.OpenEBSVersion1100:
		yamlFile = "/templates/openebs-operator-1.10.0.yaml"
	case types.OpenEBSVersion1100EE:
		yamlFile = "/templates/openebs-operator-1.10.0-ee.yaml"
	case types.OpenEBSVersion1110EE:
		yamlFile = "/templates/openebs-operator-1.11.0-ee.yaml"
	default:
		return errors.Errorf(
			"Unsupported OpenEBS version provided, version: %+v", p.ObservedOpenEBS.Spec.Version)
	}
	openEBSOperatorYaml, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return errors.Errorf(
			"Error reading YAML file for version %s: %+v", p.ObservedOpenEBS.Spec.Version, err)
	}

	// form the mapping from component's "name_kind" as key to YAML
	// string as value using operator yaml.
	componentsYAML := strings.Split(string(openEBSOperatorYaml), "---")
	for _, componentYAML := range componentsYAML {
		if componentYAML == "" {
			continue
		}
		unstructuredYAML := unstructured.Unstructured{}
		if err = yaml.Unmarshal([]byte(componentYAML), &unstructuredYAML.Object); err != nil {
			return errors.Errorf("Error unmarshalling YAML string:%s, Error: %+v", componentYAML, err)
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
	p.ComponentManifests = componentsYAMLMap
	return nil
}

// removeDisabledManifests removes the manifests which are disabled so that
// these components does not get installed.
// TODO: Delete the components if the components are disabled after installation.
func (p *Planner) removeDisabledManifests() error {
	if *p.ObservedOpenEBS.Spec.APIServer.Enabled == false {
		delete(p.ComponentManifests, types.MayaAPIServerManifestKey)
		delete(p.ComponentManifests, types.MayaAPIServerServiceManifestKey)
	}
	if *p.ObservedOpenEBS.Spec.AdmissionServer.Enabled == false {
		delete(p.ComponentManifests, types.AdmissionServerManifestKey)
	}
	if *p.ObservedOpenEBS.Spec.Provisioner.Enabled == false {
		delete(p.ComponentManifests, types.ProvisionerManifestKey)
	}
	if *p.ObservedOpenEBS.Spec.SnapshotOperator.Enabled == false {
		delete(p.ComponentManifests, types.SnapshotOperatorManifestKey)
	}
	if *p.ObservedOpenEBS.Spec.NDMDaemon.Enabled == false {
		delete(p.ComponentManifests, types.NDMConfigManifestKey)
		delete(p.ComponentManifests, types.NDMManifestKey)
	}
	if *p.ObservedOpenEBS.Spec.NDMOperator.Enabled == false {
		delete(p.ComponentManifests, types.NDMOperatorManifestKey)
	}
	if *p.ObservedOpenEBS.Spec.LocalProvisioner.Enabled == false {
		delete(p.ComponentManifests, types.LocalProvisionerManifestKey)
	}

	if *p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.Enabled == false &&
		*p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.Enabled == false {
		delete(p.ComponentManifests, types.CSINodeInfoCRDManifestKey)
		delete(p.ComponentManifests, types.CSIVolumeCRDManifestKey)
		delete(p.ComponentManifests, types.VolumeSnapshotClassCRDManifestKey)
		delete(p.ComponentManifests, types.VolumeSnapshotContentCRDManifestKey)
		delete(p.ComponentManifests, types.VolumeSnapshotCRDManifestKey)
		delete(p.ComponentManifests, types.CStorCSISnapshottterBindingManifestKey)
		delete(p.ComponentManifests, types.CStorCSISnapshottterRoleManifestKey)
		delete(p.ComponentManifests, types.CStorCSIControllerSAManifestKey)
		delete(p.ComponentManifests, types.CStorCSIProvisionerRoleManifestKey)
		delete(p.ComponentManifests, types.CStorCSIProvisionerBindingManifestKey)
		delete(p.ComponentManifests, types.CStorCSIControllerManifestKey)
		delete(p.ComponentManifests, types.CStorCSIAttacherRoleManifestKey)
		delete(p.ComponentManifests, types.CStorCSIAttacherBindingManifestKey)
		delete(p.ComponentManifests, types.CStorCSIClusterRegistrarRoleManifestKey)
		delete(p.ComponentManifests, types.CStorCSIClusterRegistrarBindingManifestKey)
		delete(p.ComponentManifests, types.CStorCSIRegistrarRoleManifestKey)
		delete(p.ComponentManifests, types.CStorCSIRegistrarBindingManifestKey)
		delete(p.ComponentManifests, types.CStorCSINodeSAManifestKey)
		delete(p.ComponentManifests, types.CStorCSINodeManifestKey)
		delete(p.ComponentManifests, types.CStorCSIDriverManifestKey)
		delete(p.ComponentManifests, types.CStorVolumeAttachmentCRDManifestKey)
	}

	if *p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.Enabled == false {
		delete(p.ComponentManifests, types.CStorCSIControllerManifestKey)
	}

	if *p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.Enabled == false {
		delete(p.ComponentManifests, types.CStorCSINodeManifestKey)
	}

	if *p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.Enabled == false {
		delete(p.ComponentManifests, types.CSPCOperatorManifestKey)
		delete(p.ComponentManifests, types.CSPCCRDV1alpha1ManifestKey)
		delete(p.ComponentManifests, types.CSPCCRDV1ManifestKey)
		delete(p.ComponentManifests, types.CSPICRDV1ManifestKey)
	}
	if *p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Enabled == false {
		delete(p.ComponentManifests, types.CVCOperatorManifestKey)
	}
	if *p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.Enabled == false &&
		*p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Enabled == false {
		delete(p.ComponentManifests, types.CStorAdmissionServerManifestKey)
	}

	p.removeMayastorManifests()
	return nil
}

// getDesiredManifests updates all the component's manifest as per the provided
// or the default values.
func (p *Planner) getDesiredManifests() error {
	var err error

	for key, value := range p.ComponentManifests {
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
		case types.KindNamespace:
			value, err = p.getDesiredNamespace(value)
		case types.KindServiceAccount:
			value, err = p.getDesiredServiceAccount(value)
		case types.KindClusterRole:
			value, err = p.getDesiredClusterRole(value)
		case types.KindClusterRoleBinding:
			value, err = p.getDesiredClusterRoleBinding(value)
		case types.KindDeployment:
			value, err = p.getDesiredDeployment(value)
		case types.KindDaemonSet:
			value, err = p.getDesiredDaemonSet(value)
		case types.KindConfigMap:
			value, err = p.getDesiredConfigmap(value)
		case types.KindService:
			value, err = p.getDesiredService(value)
		case types.KindStatefulset:
			value, err = p.getDesiredStatefulSet(value)
		case types.KindCustomResourceDefinition:
			value, err = p.getDesiredCustomResourceDefinition(value)
		case types.KindCSIDriver:
			value, err = p.getDesiredCSIDriver(value)
		default:
			// Doing nothing if an unknown kind
			continue
		}
		if err != nil {
			return errors.Errorf("Error updating manifests: %+v", err)
		}
		// update manifest with the updated values
		p.ComponentManifests[key] = value
	}
	return nil
}

// getDesiredDeployment updates the deployment manifest as per the given configuration.
func (p *Planner) getDesiredDeployment(deploy *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	var (
		replicas *int32
		err      error
	)
	nodeSelector := make(map[string]string)
	tolerations := make([]interface{}, 0)
	affinity := make(map[string]interface{})
	resources := make(map[string]interface{})
	// update the namespace
	deploy.SetNamespace(p.ObservedOpenEBS.Namespace)

	switch deploy.GetName() {
	case types.MayaAPIServerNameKey:
		replicas = p.ObservedOpenEBS.Spec.APIServer.Replicas
		resources = p.ObservedOpenEBS.Spec.APIServer.Resources
		nodeSelector = p.ObservedOpenEBS.Spec.APIServer.NodeSelector
		tolerations = p.ObservedOpenEBS.Spec.APIServer.Tolerations
		affinity = p.ObservedOpenEBS.Spec.APIServer.Affinity
		err = p.updateMayaAPIServer(deploy)

	case types.ProvisionerNameKey:
		replicas = p.ObservedOpenEBS.Spec.Provisioner.Replicas
		resources = p.ObservedOpenEBS.Spec.Provisioner.Resources
		nodeSelector = p.ObservedOpenEBS.Spec.Provisioner.NodeSelector
		tolerations = p.ObservedOpenEBS.Spec.Provisioner.Tolerations
		affinity = p.ObservedOpenEBS.Spec.Provisioner.Affinity
		err = p.updateOpenEBSProvisioner(deploy)

	case types.SnapshotOperatorNameKey:
		replicas = p.ObservedOpenEBS.Spec.SnapshotOperator.Replicas
		resources = p.ObservedOpenEBS.Spec.SnapshotOperator.Resources
		nodeSelector = p.ObservedOpenEBS.Spec.SnapshotOperator.NodeSelector
		tolerations = p.ObservedOpenEBS.Spec.SnapshotOperator.Tolerations
		affinity = p.ObservedOpenEBS.Spec.SnapshotOperator.Affinity
		err = p.updateSnapshotOperator(deploy)

	case types.NDMOperatorNameKey:
		replicas = p.ObservedOpenEBS.Spec.NDMOperator.Replicas
		resources = p.ObservedOpenEBS.Spec.NDMOperator.Resources
		nodeSelector = p.ObservedOpenEBS.Spec.NDMOperator.NodeSelector
		tolerations = p.ObservedOpenEBS.Spec.NDMOperator.Tolerations
		affinity = p.ObservedOpenEBS.Spec.NDMOperator.Affinity
		err = p.updateNDMOperator(deploy)

	case types.LocalProvisionerNameKey:
		replicas = p.ObservedOpenEBS.Spec.LocalProvisioner.Replicas
		resources = p.ObservedOpenEBS.Spec.LocalProvisioner.Resources
		nodeSelector = p.ObservedOpenEBS.Spec.LocalProvisioner.NodeSelector
		tolerations = p.ObservedOpenEBS.Spec.LocalProvisioner.Tolerations
		affinity = p.ObservedOpenEBS.Spec.LocalProvisioner.Affinity
		err = p.updateLocalProvisioner(deploy)

	case types.AdmissionServerNameKey:
		replicas = p.ObservedOpenEBS.Spec.AdmissionServer.Replicas
		resources = p.ObservedOpenEBS.Spec.AdmissionServer.Resources
		nodeSelector = p.ObservedOpenEBS.Spec.AdmissionServer.NodeSelector
		tolerations = p.ObservedOpenEBS.Spec.AdmissionServer.Tolerations
		affinity = p.ObservedOpenEBS.Spec.AdmissionServer.Affinity
		err = p.updateAdmissionServer(deploy)

	case types.CSPCOperatorNameKey:
		replicas = p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.Replicas
		resources = p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.Resources
		nodeSelector = p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.NodeSelector
		tolerations = p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.Tolerations
		affinity = p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.Affinity
		err = p.updateCSPCOperator(deploy)

	case types.CVCOperatorNameKey:
		replicas = p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Replicas
		resources = p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Resources
		nodeSelector = p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.NodeSelector
		tolerations = p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Tolerations
		affinity = p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Affinity
		err = p.updateCVCOperator(deploy)

	case types.CStorAdmissionServerNameKey:
		replicas = p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer.Replicas
		resources = p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer.Resources
		nodeSelector = p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer.NodeSelector
		tolerations = p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer.Tolerations
		affinity = p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer.Affinity
		err = p.updateCStorAdmissionServer(deploy)

	case types.MoacDeploymentNameKey:
		replicas = p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Replicas
		resources = p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Resources
		nodeSelector = p.ObservedOpenEBS.Spec.MayastorConfig.Moac.NodeSelector
		tolerations = p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Tolerations
		affinity = p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Affinity
		err = p.updateMoac(deploy)
	}
	if err != nil {
		return deploy, err
	}
	// update the replica count only if it is greater than 1 since the
	// default value itself is 1.
	// TODO: Validate the replica count value and throw error or take
	// some action based on that.
	if *replicas > 1 {
		err = unstructured.SetNestedField(deploy.Object, int64(*replicas), "spec", "replicas")
		if err != nil {
			return deploy, err
		}
	}
	containers, err := unstruct.GetNestedSliceOrError(deploy, "spec",
		"template", "spec", "containers")
	if err != nil {
		return deploy, err
	}
	updateContainer := func(obj *unstructured.Unstructured) error {
		err = unstructured.SetNestedField(obj.Object,
			p.ObservedOpenEBS.Spec.ImagePullPolicy, "spec", "imagePullPolicy")
		if err != nil {
			return err
		}
		if resources != nil {
			err = unstructured.SetNestedField(obj.Object, resources, "spec", "resources")
		} else if p.ObservedOpenEBS.Spec.Resources != nil {
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.Resources, "spec", "resources")
		}
		if err != nil {
			return err
		}
		return nil
	}
	err = unstruct.SliceIterator(containers).ForEachUpdate(updateContainer)
	if err != nil {
		return deploy, err
	}

	err = unstructured.SetNestedSlice(deploy.Object, containers, "spec", "template", "spec", "containers")
	if err != nil {
		return deploy, err
	}
	// update the nodeSelector value
	if nodeSelector != nil {
		err = unstructured.SetNestedStringMap(deploy.Object, nodeSelector, "spec",
			"template", "spec", "nodeSelector")
		if err != nil {
			return deploy, err
		}
	}
	// update the tolerations if any
	if len(tolerations) > 0 {
		err = unstructured.SetNestedSlice(deploy.Object, tolerations, "spec",
			"template", "spec", "tolerations")
		if err != nil {
			return deploy, err
		}
	}
	// update affinity if set
	if affinity != nil {
		err = unstructured.SetNestedField(deploy.Object, affinity, "spec",
			"template", "spec", "affinity")
		if err != nil {
			return deploy, err
		}
	}
	// update pod version label
	err = p.updatePodTemplateVersionLabel(deploy)
	if err != nil {
		return deploy, err
	}
	// create annotations that refers to the instance which
	// triggered creation of this deployment
	deploy.SetAnnotations(
		map[string]string{
			types.AnnKeyOpenEBSUID: string(p.ObservedOpenEBS.GetUID()),
		},
	)
	return deploy, nil
}

// getDesiredConfigmap updates the configmap manifest as per the given configuration.
func (p *Planner) getDesiredConfigmap(configmap *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	var err error

	configmap.SetNamespace(p.ObservedOpenEBS.Namespace)
	switch configmap.GetName() {
	case types.NDMConfigNameKey:
		err = p.updateNDMConfig(configmap)
	}
	if err != nil {
		return configmap, err
	}
	// create annotations that refers to the instance which
	// triggered creation of this ConfigMap
	configmap.SetAnnotations(
		map[string]string{
			types.AnnKeyOpenEBSUID: string(p.ObservedOpenEBS.GetUID()),
		},
	)
	return configmap, nil
}

// getDesiredService updates the service manifest as per the given configuration.
func (p *Planner) getDesiredService(svc *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	var err error
	svc.SetNamespace(p.ObservedOpenEBS.Namespace)
	switch svc.GetName() {
	case types.MayaAPIServerServiceNameKey:
		err = p.updateMayaAPIServerService(svc)
	case types.MoacServiceNameKey:
		err = p.updateMoacService(svc)
	case types.CVCOperatorServiceNameKey:
		err = p.updateCVCOperatorService(svc)
	}
	if err != nil {
		return svc, err
	}
	// create annotations that refers to the instance which
	// triggered creation of this Service
	svc.SetAnnotations(
		map[string]string{
			types.AnnKeyOpenEBSUID: string(p.ObservedOpenEBS.GetUID()),
		},
	)
	return svc, nil
}

// getDesiredDaemonSet updates the daemonset manifest as per the given configuration.
func (p *Planner) getDesiredDaemonSet(daemon *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	var (
		err error
	)
	nodeSelector := make(map[string]string)
	tolerations := make([]interface{}, 0)
	affinity := make(map[string]interface{})

	daemon.SetNamespace(p.ObservedOpenEBS.Namespace)
	switch daemon.GetName() {
	case types.NDMNameKey:
		nodeSelector = p.ObservedOpenEBS.Spec.NDMDaemon.NodeSelector
		tolerations = p.ObservedOpenEBS.Spec.NDMDaemon.Tolerations
		affinity = p.ObservedOpenEBS.Spec.NDMDaemon.Affinity
		err = p.updateNDM(daemon)
	case types.CStorCSINodeNameKey:
		err = p.updateOpenEBSCStorCSINode(daemon)
	case types.MayastorDaemonsetNameKey:
		err = p.updateMayastor(daemon)
	}
	if err != nil {
		return daemon, err
	}
	// update the daemonset containers with the images and imagePullPolicy
	containers, err := unstruct.GetNestedSliceOrError(daemon, "spec", "template", "spec", "containers")
	if err != nil {
		return daemon, err
	}
	updateContainer := func(obj *unstructured.Unstructured) error {
		err = unstructured.SetNestedField(obj.Object,
			p.ObservedOpenEBS.Spec.ImagePullPolicy, "spec", "imagePullPolicy")
		if err != nil {
			return err
		}
		return nil
	}
	err = unstruct.SliceIterator(containers).ForEachUpdate(updateContainer)
	if err != nil {
		return daemon, err
	}
	err = unstructured.SetNestedSlice(daemon.Object, containers, "spec",
		"template", "spec", "containers")
	if err != nil {
		return daemon, err
	}
	// update the nodeSelector value
	if nodeSelector != nil {
		err = unstructured.SetNestedStringMap(daemon.Object, nodeSelector, "spec",
			"template", "spec", "nodeSelector")
		if err != nil {
			return daemon, err
		}
	}
	// update the tolerations if any
	if len(tolerations) > 0 {
		err = unstructured.SetNestedSlice(daemon.Object, tolerations, "spec", "template", "spec",
			"tolerations")
		if err != nil {
			return daemon, err
		}
	}
	// update affinity if set
	if affinity != nil {
		err = unstructured.SetNestedField(daemon.Object, affinity, "spec", "template", "spec",
			"affinity")
		if err != nil {
			return daemon, err
		}
	}
	// update pod version label
	err = p.updatePodTemplateVersionLabel(daemon)
	if err != nil {
		return daemon, err
	}
	// create annotations that refers to the instance which
	// triggered creation of this DaemonSet
	daemon.SetAnnotations(
		map[string]string{
			types.AnnKeyOpenEBSUID: string(p.ObservedOpenEBS.GetUID()),
		},
	)
	return daemon, nil
}

// getDesiredStatefulSet updates the statefulset manifest as per the given configuration.
func (p *Planner) getDesiredStatefulSet(statefulset *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	var err error
	switch statefulset.GetName() {
	case types.CStorCSIControllerNameKey:
		err = p.updateOpenEBSCStorCSIController(statefulset)
		if err != nil {
			return statefulset, err
		}
	}
	// update pod version label
	err = p.updatePodTemplateVersionLabel(statefulset)
	if err != nil {
		return statefulset, err
	}
	// create annotations that refers to the instance which
	// triggered creation of this StatefulSet
	statefulset.SetAnnotations(
		map[string]string{
			types.AnnKeyOpenEBSUID: string(p.ObservedOpenEBS.GetUID()),
		},
	)

	return statefulset, nil
}

// getDesiredCSIDriver updates the csidrivers manifest as per the given configuration.
func (p *Planner) getDesiredCSIDriver(driver *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := driver.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for CSIDriver controller:
	// 1. openebs-upgrade.dao.mayadata.io/component-group: cstor-csi
	// 2. openebs-upgrade.dao.mayadata.io/component-name: cstor.csi.openebs.io
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSCStorCSIComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] =
		types.CStorCSIDriverNameKey
	// set the desired labels
	driver.SetLabels(desiredLabels)
	// create annotations that refers to the instance which
	// triggered creation of this CSIDriver
	driver.SetAnnotations(
		map[string]string{
			types.AnnKeyOpenEBSUID: string(p.ObservedOpenEBS.GetUID()),
		},
	)

	return driver, nil
}

// compareVersion compares given version i.e v1 and v2.
// It returns -1 if v1 is less than v2, 0 if v1 equal to v2, 1 if v1 is greater than v2.
// It returns -2 in case of any error with error.
func compareVersion(v1, v2 string) (int, error) {

	// removes alphabets from the version.
	reg := regexp.MustCompile("[^\\d.]")
	v1 = reg.ReplaceAllString(v1, "")
	v2 = reg.ReplaceAllString(v2, "")

	v1Array := strings.Split(v1, ".")
	v2Array := strings.Split(v2, ".")

	for i := 0; i < len(v1Array) || i < len(v2Array); i++ {
		if i < len(v1Array) && i < len(v2Array) {
			v1, err := strconv.Atoi(v1Array[i])
			v2, err := strconv.Atoi(v2Array[i])
			if err != nil {
				return -2, err
			}
			if v1 < v2 {
				return -1, nil
			} else if v1 > v2 {
				return 1, nil
			}
		} else if i < len(v1Array) {
			v1, err := strconv.Atoi(v1Array[i])
			if err != nil {
				return -2, err
			}
			if v1 != 0 {
				return 1, nil
			}
		} else if i < len(v2Array) {
			v2, err := strconv.Atoi(v2Array[i])
			if err != nil {
				return -2, err
			}
			if v2 != 0 {
				return -1, nil
			}
		}
	}

	return 0, nil
}

// updatePodTemplateVersionLabel updates the version label of pod template present in the deployment,
// daemonset, statefulset, etc with the given version and imageTagSuffix.
func (p *Planner) updatePodTemplateVersionLabel(resource *unstructured.Unstructured) error {
	// update pod version label
	podLabels, exist, err := unstructured.NestedStringMap(resource.Object, "spec",
		"template", "metadata", "labels")
	if err != nil {
		return err
	}
	if exist {
		if _, exist := podLabels[types.OpenEBSVersionLabelKey]; exist {
			// update the version label as per the given OpenEBS version and imageTagSuffix
			// if given.
			podLabels[types.OpenEBSVersionLabelKey] =
				p.ObservedOpenEBS.Spec.Version + p.ObservedOpenEBS.Spec.ImageTagSuffix
			err = unstructured.SetNestedStringMap(resource.Object, podLabels, "spec",
				"template", "metadata", "labels")
			if err != nil {
				return err
			}
		}
	}
	return nil
}
