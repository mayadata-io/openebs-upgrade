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
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/k8s"
	"mayadata.io/openebs-upgrade/types"
	"mayadata.io/openebs-upgrade/unstruct"
	"strings"
)

const (
	// ContainerOpenEBSCSIPlugin is the name of the container openebs csi plugin
	ContainerOpenEBSCSIPluginName string = "openebs-csi-plugin"
	// EnvOpenEBSNamespaceKey is the env key for openebs namespace
	EnvOpenEBSNamespaceKey string = "OPENEBS_NAMESPACE"
	// DefaultCSPCOperatorReplicaCount is the default replica count for
	// cspc-operatot.
	DefaultCSPCOperatorReplicaCount int32 = 1
	// DefaultCVCOperatorReplicaCount is the default replica count for
	// cvc-operatot.
	DefaultCVCOperatorReplicaCount int32 = 1
)

// Set the default values for Cstor if not already given.
func (p *Planner) setCStorDefaultsIfNotSet() error {
	if p.ObservedOpenEBS.Spec.CstorConfig == nil {
		p.ObservedOpenEBS.Spec.CstorConfig = &types.CstorConfig{}
	}
	// form the cstor-pool image
	if p.ObservedOpenEBS.Spec.CstorConfig.Pool.ImageTag == "" {
		p.ObservedOpenEBS.Spec.CstorConfig.Pool.ImageTag = p.ObservedOpenEBS.Spec.Version
	}
	p.ObservedOpenEBS.Spec.CstorConfig.Pool.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"cstor-pool:" + p.ObservedOpenEBS.Spec.CstorConfig.Pool.ImageTag

	// form the cstor-pool-mgmt image
	if p.ObservedOpenEBS.Spec.CstorConfig.PoolMgmt.ImageTag == "" {
		p.ObservedOpenEBS.Spec.CstorConfig.PoolMgmt.ImageTag = p.ObservedOpenEBS.Spec.Version
	}
	p.ObservedOpenEBS.Spec.CstorConfig.PoolMgmt.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"cstor-pool-mgmt:" + p.ObservedOpenEBS.Spec.CstorConfig.PoolMgmt.ImageTag

	// form the cstor-istgt image
	if p.ObservedOpenEBS.Spec.CstorConfig.Target.ImageTag == "" {
		p.ObservedOpenEBS.Spec.CstorConfig.Target.ImageTag = p.ObservedOpenEBS.Spec.Version
	}
	p.ObservedOpenEBS.Spec.CstorConfig.Target.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"cstor-istgt:" + p.ObservedOpenEBS.Spec.CstorConfig.Target.ImageTag

	// form the cstor-volume-mgmt image
	if p.ObservedOpenEBS.Spec.CstorConfig.VolumeMgmt.ImageTag == "" {
		p.ObservedOpenEBS.Spec.CstorConfig.VolumeMgmt.ImageTag = p.ObservedOpenEBS.Spec.Version
	}
	p.ObservedOpenEBS.Spec.CstorConfig.VolumeMgmt.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"cstor-volume-mgmt:" + p.ObservedOpenEBS.Spec.CstorConfig.VolumeMgmt.ImageTag

	// form the cspi-mgmt image(CSPI_MGMT)
	if p.ObservedOpenEBS.Spec.CstorConfig.CSPIMgmt.ImageTag == "" {
		p.ObservedOpenEBS.Spec.CstorConfig.CSPIMgmt.ImageTag = p.ObservedOpenEBS.Spec.Version
	}
	p.ObservedOpenEBS.Spec.CstorConfig.CSPIMgmt.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"cspi-mgmt:" + p.ObservedOpenEBS.Spec.CstorConfig.CSPIMgmt.ImageTag

	// set the CSPC operator defaults
	if p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator == nil {
		p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator = &types.CSPCOperator{}
	}
	if p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.Enabled == nil {
		p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.Enabled = true
	}
	if p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.ImageTag == "" {
		p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.ImageTag = p.ObservedOpenEBS.Spec.Version
	}
	// form the container image as per the image prefix and image tag.
	p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"cspc-operator:" + p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.ImageTag
	if p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.Replicas == nil {
		p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.Replicas = new(int32)
		*p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.Replicas = DefaultCSPCOperatorReplicaCount
	}
	// set the CVC operator defaults
	if p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator == nil {
		p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator = &types.CVCOperator{}
	}
	if p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Enabled == nil {
		p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Enabled = true
	}
	if p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.ImageTag == "" {
		p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.ImageTag = p.ObservedOpenEBS.Spec.Version
	}
	// form the container image as per the image prefix and image tag.
	p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"cvc-operator:" + p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.ImageTag
	if p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Replicas == nil {
		p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Replicas = new(int32)
		*p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Replicas = DefaultCVCOperatorReplicaCount
	}

	p.setCSIDefaultsIfNotSet()

	return nil
}

func (p *Planner) setCSIDefaultsIfNotSet() {

	isCSISupported, err := p.isCSISupported()
	// Do not return the error as not to block installing other components.
	if err != nil {
		isCSISupported = false
		glog.Errorf("Failed to set CSI defaults, error: %v", err)
	}

	if !isCSISupported {
		glog.V(5).Infof("Skipping CSI installation.")
	}

	// Set the default values for cstor csi controller statefulset in configuration.
	if p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.Enabled == nil {
		p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.Enabled = true
	}

	if !isCSISupported && *p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.Enabled == true {
		*p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.Enabled = false
	}

	if *p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.Enabled == true {
		if p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.ImageTag == "" {
			p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.ImageTag = p.ObservedOpenEBS.Spec.Version
		}
		p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
			"cstor-csi-driver:" + p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.ImageTag
	}

	// Set the default values for cstor csi node daemonset in configuration.
	if p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.Enabled == nil {
		p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.Enabled = true
	}

	if !isCSISupported && *p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.Enabled == true {
		*p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.Enabled = false
	}

	if *p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.Enabled == true {
		if p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.ImageTag == "" {
			p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.ImageTag = p.ObservedOpenEBS.Spec.Version
		}
		p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
			"cstor-csi-driver:" + p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.ImageTag

		if p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.ISCSIPath == "" {
			p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.ISCSIPath = "/sbin/iscsiadm"
		}
	}
}

// isCSISupported checks if csi is supported or not in the current kubernetes cluster, if not it will
// return false else true.
func (p *Planner) isCSISupported() (bool, error) {
	// get the kubernetes version.
	k8sVersion, err := k8s.GetK8sVersion()
	if err != nil {
		return false, errors.Errorf("Unable to find kubernetes version, error: %v", err)
	}

	// compare the kubernetes version with the supported version of csi.
	comp, err := compareVersion(k8sVersion, types.CSISupportedVersion)
	if err != nil {
		return false, errors.Errorf("Error comparing versions, error: %v", err)
	}

	if comp < 0 {
		glog.Warningf("CSI is not supported in %s Kubernetes version. "+
			"CSI is supported from %s Kubernetes version.", k8sVersion, types.CSISupportedVersion)
		return false, nil
	}

	return true, nil
}

// updateOpenEBSCStorCSINode updates the values of openebs-cstor-csi-node daemonset as per given configuration.
func (p *Planner) updateOpenEBSCStorCSINode(daemonset *unstructured.Unstructured) error {
	// overwrite the namespace to kube-system as csi based components will run only in kube-system namespace.
	daemonset.SetNamespace(types.NamespaceKubeSystem)

	// this will get the extra volumes and volume mounts required to be added in the csi node daemonset
	// for the csi to work for different OS distributions/versions.
	// This volumes and volume mounts will be added in the openebs-csi-plugin container.
	extraVolumes, extraVolumeMounts, err := p.getOSSpecificVolumeMounts()
	if err != nil {
		return err
	}

	volumes, err := unstruct.GetNestedSliceOrError(daemonset, "spec", "template", "spec", "volumes")
	if err != nil {
		return err
	}
	// updateVolume updates the volume path of openebs-csi-plugin container.
	updateVolume := func(obj *unstructured.Unstructured) error {
		volumeName, err := unstruct.GetString(obj, "spec", "name")
		if err != nil {
			return err
		}
		if volumeName == "iscsiadm-bin" {
			err = unstructured.SetNestedField(obj.Object,
				p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.ISCSIPath, "spec", "hostPath", "path")
			if err != nil {
				return err
			}
		}
		return nil
	}
	err = unstruct.SliceIterator(volumes).ForEachUpdate(updateVolume)
	if err != nil {
		return err
	}

	// Append the new extra volumes with the existing volumes, required for the csi to work.
	volumes = append(volumes, extraVolumes...)

	err = unstructured.SetNestedSlice(daemonset.Object, volumes,
		"spec", "template", "spec", "volumes")
	if err != nil {
		return err
	}

	containers, err := unstruct.GetNestedSliceOrError(daemonset, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}

	// update the env value of openebs-csi-plugin container
	updateOpenEBSCSIPluginEnv := func(env *unstructured.Unstructured) error {
		envName, _, err := unstructured.NestedString(env.Object, "spec", "name")
		if err != nil {
			return err
		}
		if envName == EnvOpenEBSNamespaceKey {
			unstructured.SetNestedField(env.Object, p.ObservedOpenEBS.Namespace, "spec", "value")
		}
		return nil
	}

	// updateOpenEBSCSIPluginVolumeMount updates the volumeMounts path of openebs-csi-plugin container.
	updateOpenEBSCSIPluginVolumeMount := func(vm *unstructured.Unstructured) error {
		vmName, _, err := unstructured.NestedString(vm.Object, "spec", "name")
		if err != nil {
			return err
		}
		if vmName == "iscsiadm-bin" {
			err = unstructured.SetNestedField(vm.Object,
				p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.ISCSIPath, "spec", "mountPath")
			if err != nil {
				return err
			}
		}
		return nil
	}

	// update the containers
	updateContainer := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}
		envs, _, err := unstruct.GetSlice(obj, "spec", "env")
		if err != nil {
			return err
		}
		volumeMounts, _, err := unstruct.GetSlice(obj, "spec", "volumeMounts")
		if err != nil {
			return err
		}

		if containerName == ContainerOpenEBSCSIPluginName {
			// Set the image of the container.
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.Image,
				"spec", "image")
			if err != nil {
				return err
			}
			// Set the environmets of the container.
			err = unstruct.SliceIterator(envs).ForEachUpdate(updateOpenEBSCSIPluginEnv)
			if err != nil {
				return err
			}
			err = unstruct.SliceIterator(volumeMounts).ForEachUpdate(updateOpenEBSCSIPluginVolumeMount)
			if err != nil {
				return err
			}
		}
		err = unstructured.SetNestedSlice(obj.Object, envs, "spec", "env")
		if err != nil {
			return err
		}

		// Append the new extra volume mounts with the existing volume mounts, required for the csi to work.
		volumeMounts = append(volumeMounts, extraVolumeMounts...)
		err = unstructured.SetNestedSlice(obj.Object, volumeMounts, "spec", "volumeMounts")
		if err != nil {
			return err
		}

		return nil
	}

	// Update the containers.
	err = unstruct.SliceIterator(containers).ForEachUpdate(updateContainer)
	if err != nil {
		return err
	}

	// Set back the value of the containers.
	err = unstructured.SetNestedSlice(daemonset.Object,
		containers, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}

	return nil
}

// getOSSpecificVolumeMounts returns the volume and volume mounts based on the specific OS distribution/version.
// This volume and volume mounts are for the specific container i.e openebs-csi-plugin.
// This function will get the OS Image and the version for the ubuntu distribution and will return the
// volumes and volume mounts accordngly.
func (p *Planner) getOSSpecificVolumeMounts() ([]interface{}, []interface{}, error) {
	volumes := make([]interface{}, 0)
	volumeMounts := make([]interface{}, 0)

	osImage, err := k8s.GetOSImage()
	if err != nil {
		return volumes, volumeMounts, errors.Errorf("Error getting OS Image of a Node, error: %+v", err)
	}

	ubuntuVersion, err := k8s.GetUbuntuVersion()
	if err != nil {
		return volumes, volumeMounts, errors.Errorf("Error getting Ubuntu Version of a Node, error: %+v", err)
	}

	switch true {
	case strings.Contains(strings.ToLower(osImage), strings.ToLower(types.OSImageSLES12)):
		volumes, volumeMounts = p.getSUSE12VolumeMounts()
	case strings.Contains(strings.ToLower(osImage), strings.ToLower(types.OSImageSLES15)):
		volumes, volumeMounts = p.getSUSE15VolumeMounts()
	case strings.Contains(strings.ToLower(osImage), strings.ToLower(types.OSImageUbuntu1804)) ||
		((ubuntuVersion != 0) && ubuntuVersion >= 18.04):
		volumes, volumeMounts = p.getUbuntu1804VolumeMounts()
	}

	return volumes, volumeMounts, nil
}

// getSUSE12VolumeMounts returns the volumes and volume mounts for suse 12.
func (p *Planner) getSUSE12VolumeMounts() ([]interface{}, []interface{}) {
	volumes := make([]interface{}, 0)
	volumeMounts := make([]interface{}, 0)

	// Create new volumes for suse 12.
	libCryptoVolume := map[string]interface{}{
		"name": "iscsiadm-lib-crypto",
		"hostPath": map[string]interface{}{
			"type": "File",
			"path": "/lib64/libcrypto.so.1.0.0",
		},
	}
	libOpeniscsiusrVolume := map[string]interface{}{
		"name": "iscsiadm-lib-openiscsiusr",
		"hostPath": map[string]interface{}{
			"type": "File",
			"path": "/usr/lib64/libopeniscsiusr.so.0.2.0",
		},
	}
	volumes = append(volumes, libCryptoVolume, libOpeniscsiusrVolume)

	// Create new volume mounts for suse 12.
	libCryptoVolumeMount := map[string]interface{}{
		"name":      "iscsiadm-lib-crypto",
		"mountPath": "/lib/x86_64-linux-gnu/libcrypto.so.1.0.0",
	}
	libOpeniscsiusrVolumeMount := map[string]interface{}{
		"name":      "iscsiadm-lib-openiscsiusr",
		"mountPath": "/lib/x86_64-linux-gnu/libopeniscsiusr.so.0.2.0",
	}
	volumeMounts = append(volumeMounts, libCryptoVolumeMount, libOpeniscsiusrVolumeMount)

	return volumes, volumeMounts
}

// getSUSE15VolumeMounts returns the volumes and volume mounts for suse 15.
func (p *Planner) getSUSE15VolumeMounts() ([]interface{}, []interface{}) {
	volumes := make([]interface{}, 0)
	volumeMounts := make([]interface{}, 0)

	// Create new volumes for suse 15.
	libCryptoVolume := map[string]interface{}{
		"name": "iscsiadm-lib-crypto",
		"hostPath": map[string]interface{}{
			"type": "File",
			"path": "/usr/lib64/libcrypto.so.1.1",
		},
	}
	libOpeniscsiusrVolume := map[string]interface{}{
		"name": "iscsiadm-lib-openiscsiusr",
		"hostPath": map[string]interface{}{
			"type": "File",
			"path": "/usr/lib64/libopeniscsiusr.so.0.2.0",
		},
	}
	volumes = append(volumes, libCryptoVolume, libOpeniscsiusrVolume)

	// Create new volume mounts for suse 15.
	libCryptoVolumeMount := map[string]interface{}{
		"name":      "iscsiadm-lib-crypto",
		"mountPath": "/lib/x86_64-linux-gnu/libcrypto.so.1.1",
	}
	libOpeniscsiusrVolumeMount := map[string]interface{}{
		"name":      "iscsiadm-lib-openiscsiusr",
		"mountPath": "/lib/x86_64-linux-gnu/libopeniscsiusr.so.0.2.0",
	}
	volumeMounts = append(volumeMounts, libCryptoVolumeMount, libOpeniscsiusrVolumeMount)

	return volumes, volumeMounts
}

// getUbuntu1804VolumeMounts returns the volumes and volume mounts for ubuntu 18.04 and above.
func (p *Planner) getUbuntu1804VolumeMounts() ([]interface{}, []interface{}) {
	volumes := make([]interface{}, 0)
	volumeMounts := make([]interface{}, 0)

	// Create new volume for ubuntu 18.04 and above.
	volume := map[string]interface{}{
		"name": "iscsiadm-lib-isns-nocrypto",
		"hostPath": map[string]interface{}{
			"type": "File",
			"path": "/lib/x86_64-linux-gnu/libisns-nocrypto.so.0",
		},
	}
	volumes = append(volumes, volume)

	// Create new volume mount for ubuntu 18.04 and above.
	volumeMount := map[string]interface{}{
		"name":      "iscsiadm-lib-isns-nocrypto",
		"mountPath": "/lib/x86_64-linux-gnu/libisns-nocrypto.so.0",
	}
	volumeMounts = append(volumeMounts, volumeMount)

	return volumes, volumeMounts
}

// updateOpenEBSCStorCSIController updates the values of openebs-cstor-csi-controller statefulset as per given configuration.
func (p *Planner) updateOpenEBSCStorCSIController(statefulset *unstructured.Unstructured) error {
	// overwrite the namespace to kube-system as csi based components will run only in kube-system namespace.
	statefulset.SetNamespace(types.NamespaceKubeSystem)

	containers, err := unstruct.GetNestedSliceOrError(statefulset, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}

	// update the env value of openebs-csi-plugin container
	updateOpenEBSCSIPluginEnv := func(env *unstructured.Unstructured) error {
		envName, _, err := unstructured.NestedString(env.Object, "spec", "name")
		if err != nil {
			return err
		}
		if envName == EnvOpenEBSNamespaceKey {
			unstructured.SetNestedField(env.Object, p.ObservedOpenEBS.Namespace, "spec", "value")
		}
		return nil
	}

	// update the containers
	updateContainer := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}
		envs, _, err := unstruct.GetSlice(obj, "spec", "env")
		if err != nil {
			return err
		}

		if containerName == ContainerOpenEBSCSIPluginName {
			// Set the image of the container.
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.Image,
				"spec", "image")
			if err != nil {
				return err
			}
			// Set the environmets of the container.
			err = unstruct.SliceIterator(envs).ForEachUpdate(updateOpenEBSCSIPluginEnv)
			if err != nil {
				return err
			}
		}
		err = unstructured.SetNestedSlice(obj.Object, envs, "spec", "env")
		if err != nil {
			return err
		}

		return nil
	}

	err = unstruct.SliceIterator(containers).ForEachUpdate(updateContainer)
	if err != nil {
		return err
	}

	err = unstructured.SetNestedSlice(statefulset.Object,
		containers, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}

	return nil
}

// updateCSPCOperator updates the CSPC operator manifest as per the reconcile.ObservedOpenEBS values.
func (p *Planner) updateCSPCOperator(deploy *unstructured.Unstructured) error {
	// get the containers of the cspc-operator and update the desired fields
	containers, err := unstruct.GetNestedSliceOrError(deploy, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}
	// update the env value of cspc-operator container
	updateCSPCOperatorEnv := func(env *unstructured.Unstructured) error {
		envName, _, err := unstructured.NestedString(env.Object, "spec", "name")
		if err != nil {
			return err
		}
		if envName == "OPENEBS_IO_BASE_DIR" {
			err = unstructured.SetNestedField(env.Object, p.ObservedOpenEBS.Spec.DefaultStoragePath,
				"spec", "value")
		} else if envName == "OPENEBS_IO_CSTOR_POOL_SPARSE_DIR" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.DefaultStoragePath+"/sparse", "spec", "value")
		} else if envName == "OPENEBS_IO_CSPI_MGMT_IMAGE" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.CstorConfig.CSPIMgmt.Image, "spec", "value")
		} else if envName == "OPENEBS_IO_CSTOR_POOL_IMAGE" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.CstorConfig.Pool.Image, "spec", "value")
		} else if envName == "OPENEBS_IO_CSTOR_POOL_EXPORTER_IMAGE" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.Policies.Monitoring.Image, "spec", "value")
		}
		if err != nil {
			return err
		}

		return nil
	}
	updateContainer := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}
		envs, _, err := unstruct.GetSlice(obj, "spec", "env")
		if err != nil {
			return err
		}
		// update the envs of cspc-operator container
		// In order to update envs of other containers, just write an updateEnv
		// function for specific containers.
		if containerName == "cspc-operator" {
			err = unstruct.SliceIterator(envs).ForEachUpdate(updateCSPCOperatorEnv)
			if err != nil {
				return err
			}
		}
		err = unstructured.SetNestedSlice(obj.Object, envs, "spec", "env")
		if err != nil {
			return err
		}
		return nil
	}
	err = unstruct.SliceIterator(containers).ForEachUpdate(updateContainer)
	if err != nil {
		return err
	}
	err = unstructured.SetNestedSlice(deploy.Object,
		containers, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}

	return nil
}

// updateCVCOperator updates the CVC operator manifest as per the reconcile.ObservedOpenEBS values.
func (p *Planner) updateCVCOperator(deploy *unstructured.Unstructured) error {
	// get the containers of the cvc-operator and update the desired fields
	containers, err := unstruct.GetNestedSliceOrError(deploy, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}
	// update the env value of cvc-operator container
	updateCVCOperatorEnv := func(env *unstructured.Unstructured) error {
		envName, _, err := unstructured.NestedString(env.Object, "spec", "name")
		if err != nil {
			return err
		}
		if envName == "OPENEBS_IO_BASE_DIR" {
			err = unstructured.SetNestedField(env.Object, p.ObservedOpenEBS.Spec.DefaultStoragePath,
				"spec", "value")
		} else if envName == "OPENEBS_IO_CSTOR_TARGET_DIR" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.DefaultStoragePath+"/sparse", "spec", "value")
		} else if envName == "OPENEBS_IO_CSTOR_TARGET_IMAGE" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.CstorConfig.Target.Image, "spec", "value")
		} else if envName == "OPENEBS_IO_CSTOR_VOLUME_MGMT_IMAGE" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.CstorConfig.VolumeMgmt.Image, "spec", "value")
		} else if envName == "OPENEBS_IO_VOLUME_MONITOR_IMAGE" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.Policies.Monitoring.Image, "spec", "value")
		}
		if err != nil {
			return err
		}

		return nil
	}
	updateContainer := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}
		envs, _, err := unstruct.GetSlice(obj, "spec", "env")
		if err != nil {
			return err
		}
		// update the envs of cvc-operator container
		// In order to update envs of other containers, just write an updateEnv
		// function for specific containers.
		if containerName == "cvc-operator" {
			err = unstruct.SliceIterator(envs).ForEachUpdate(updateCVCOperatorEnv)
			if err != nil {
				return err
			}
		}
		err = unstructured.SetNestedSlice(obj.Object, envs, "spec", "env")
		if err != nil {
			return err
		}
		return nil
	}
	err = unstruct.SliceIterator(containers).ForEachUpdate(updateContainer)
	if err != nil {
		return err
	}
	err = unstructured.SetNestedSlice(deploy.Object,
		containers, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}

	return nil
}
