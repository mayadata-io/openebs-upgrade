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
	}
	return nil
}

// setDefaultImagePrefixIfNotSet sets the default registry prefix for
// all the container images if not already set.
// It also checks if the given image registry ends with a forward slash
// or not, if not then it adds one.
func (p *Planner) setDefaultImagePrefixIfNotSet() error {
	if p.ObservedOpenEBS.Spec.ImagePrefix == "" {
		p.ObservedOpenEBS.Spec.ImagePrefix = "quay.io/openebs/"
	} else if !strings.HasSuffix(p.ObservedOpenEBS.Spec.ImagePrefix, "/") {
		p.ObservedOpenEBS.Spec.ImagePrefix = p.ObservedOpenEBS.Spec.ImagePrefix + "/"
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
	default:
		return errors.Errorf(
			"Unsupported OpenEBS version provided, version: %+v", p.ObservedOpenEBS.Spec.Version)
	}
	data, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return errors.Errorf(
			"Error reading YAML file for version %s: %+v", p.ObservedOpenEBS.Spec.Version, err)
	}
	// form the mapping from component's "name_kind" as key to YAML
	// string as value using operator yaml.
	componentsYAML := strings.Split(string(data), "---")
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

	return nil
}

// getDesiredManifests updates all the component's manifest as per the provided
// or the default values.
func (p *Planner) getDesiredManifests() error {
	var err error

	for key, value := range p.ComponentManifests {
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
// TODO: Make this method modular, it is a big method which seems to be doing multiple
// things.
func (p *Planner) getDesiredDeployment(deploy *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	var (
		replicas         *int32
		image            string
		provisionerImage string
		controllerImage  string
		err              error
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
		image = p.ObservedOpenEBS.Spec.APIServer.Image
		resources = p.ObservedOpenEBS.Spec.APIServer.Resources
		nodeSelector = p.ObservedOpenEBS.Spec.APIServer.NodeSelector
		tolerations = p.ObservedOpenEBS.Spec.APIServer.Tolerations
		affinity = p.ObservedOpenEBS.Spec.APIServer.Affinity
		p.updateMayaAPIServer(deploy)

	case types.ProvisionerNameKey:
		replicas = p.ObservedOpenEBS.Spec.Provisioner.Replicas
		image = p.ObservedOpenEBS.Spec.Provisioner.Image
		resources = p.ObservedOpenEBS.Spec.Provisioner.Resources
		nodeSelector = p.ObservedOpenEBS.Spec.Provisioner.NodeSelector
		tolerations = p.ObservedOpenEBS.Spec.Provisioner.Tolerations
		affinity = p.ObservedOpenEBS.Spec.Provisioner.Affinity

	case types.SnapshotOperatorNameKey:
		replicas = p.ObservedOpenEBS.Spec.SnapshotOperator.Replicas
		provisionerImage = p.ObservedOpenEBS.Spec.SnapshotOperator.Provisioner.Image
		controllerImage = p.ObservedOpenEBS.Spec.SnapshotOperator.Controller.Image
		resources = p.ObservedOpenEBS.Spec.SnapshotOperator.Resources
		nodeSelector = p.ObservedOpenEBS.Spec.SnapshotOperator.NodeSelector
		tolerations = p.ObservedOpenEBS.Spec.SnapshotOperator.Tolerations
		affinity = p.ObservedOpenEBS.Spec.SnapshotOperator.Affinity

	case types.NDMOperatorNameKey:
		replicas = p.ObservedOpenEBS.Spec.NDMOperator.Replicas
		image = p.ObservedOpenEBS.Spec.NDMOperator.Image
		resources = p.ObservedOpenEBS.Spec.NDMOperator.Resources
		nodeSelector = p.ObservedOpenEBS.Spec.NDMOperator.NodeSelector
		tolerations = p.ObservedOpenEBS.Spec.NDMOperator.Tolerations
		affinity = p.ObservedOpenEBS.Spec.NDMOperator.Affinity
		p.updateNDMOperator(deploy)

	case types.LocalProvisionerNameKey:
		replicas = p.ObservedOpenEBS.Spec.LocalProvisioner.Replicas
		image = p.ObservedOpenEBS.Spec.LocalProvisioner.Image
		resources = p.ObservedOpenEBS.Spec.LocalProvisioner.Resources
		nodeSelector = p.ObservedOpenEBS.Spec.LocalProvisioner.NodeSelector
		tolerations = p.ObservedOpenEBS.Spec.LocalProvisioner.Tolerations
		affinity = p.ObservedOpenEBS.Spec.LocalProvisioner.Affinity
		p.updateLocalProvisioner(deploy)

	case types.AdmissionServerNameKey:
		replicas = p.ObservedOpenEBS.Spec.AdmissionServer.Replicas
		image = p.ObservedOpenEBS.Spec.AdmissionServer.Image
		resources = p.ObservedOpenEBS.Spec.AdmissionServer.Resources
		nodeSelector = p.ObservedOpenEBS.Spec.AdmissionServer.NodeSelector
		tolerations = p.ObservedOpenEBS.Spec.AdmissionServer.Tolerations
		affinity = p.ObservedOpenEBS.Spec.AdmissionServer.Affinity
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
		// Explicitly checking for openebs-snapshot-operator in order to update
		// its multiple containers.
		// TODO: handle multiple container update cases in a better way, this seems
		// to be a very naive way.
		if deploy.GetName() == types.SnapshotOperatorNameKey {
			containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
			if err != nil {
				return err
			}
			if containerName == types.SnapshotControllerContainerKey {
				err = unstructured.SetNestedField(obj.Object, controllerImage, "spec", "image")
			} else if containerName == types.SnapshotProvisionerContainerKey {
				err = unstructured.SetNestedField(obj.Object, provisionerImage, "spec", "image")
			}
		} else {
			err = unstructured.SetNestedField(obj.Object, image, "spec", "image")
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
	configmap.SetNamespace(p.ObservedOpenEBS.Namespace)
	switch configmap.GetName() {
	case types.NDMConfigNameKey:
		p.updateNDMConfig(configmap)
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
	svc.SetNamespace(p.ObservedOpenEBS.Namespace)
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
		image string
	)
	resources := make(map[string]interface{})
	nodeSelector := make(map[string]string)
	tolerations := make([]interface{}, 0)
	affinity := make(map[string]interface{})

	daemon.SetNamespace(p.ObservedOpenEBS.Namespace)
	switch daemon.GetName() {
	case types.NDMNameKey:
		image = p.ObservedOpenEBS.Spec.NDMDaemon.Image
		resources = p.ObservedOpenEBS.Spec.NDMDaemon.Resources
		nodeSelector = p.ObservedOpenEBS.Spec.NDMDaemon.NodeSelector
		tolerations = p.ObservedOpenEBS.Spec.NDMDaemon.Tolerations
		affinity = p.ObservedOpenEBS.Spec.NDMDaemon.Affinity
		p.updateNDM(daemon)
	}
	// update the daemonset containers with the images and imagePullPolicy
	containers, err := unstruct.GetNestedSliceOrError(daemon, "spec", "template", "spec", "containers")
	if err != nil {
		return daemon, err
	}
	updateContainer := func(obj *unstructured.Unstructured) error {
		err = unstructured.SetNestedField(obj.Object, image, "spec", "image")
		if err != nil {
			return err
		}
		err = unstructured.SetNestedField(obj.Object,
			p.ObservedOpenEBS.Spec.ImagePullPolicy, "spec", "imagePullPolicy")
		if err != nil {
			return err
		}
		if resources != nil {
			err = unstructured.SetNestedField(obj.Object, resources, "spec", "resources")
		} else if p.ObservedOpenEBS.Spec.Resources != nil {
			err = unstructured.SetNestedField(obj.Object,
				p.ObservedOpenEBS.Spec.Resources, "spec", "resources")
		}
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
	// create annotations that refers to the instance which
	// triggered creation of this DaemonSet
	daemon.SetAnnotations(
		map[string]string{
			types.AnnKeyOpenEBSUID: string(p.ObservedOpenEBS.GetUID()),
		},
	)
	return daemon, nil
}
