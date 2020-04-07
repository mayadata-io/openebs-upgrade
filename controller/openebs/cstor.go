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
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/types"
	"mayadata.io/openebs-upgrade/unstruct"
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

	// Set the default values for cstor csi controller statefulset in configuration.
	if p.ObservedOpenEBS.Spec.CstorConfig.CStorCSI.CStorCSIController.Enabled == nil {
		p.ObservedOpenEBS.Spec.CstorConfig.CStorCSI.CStorCSIController.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.CstorConfig.CStorCSI.CStorCSIController.Enabled = true
	}

	if *p.ObservedOpenEBS.Spec.CstorConfig.CStorCSI.CStorCSIController.Enabled == true {
		if p.ObservedOpenEBS.Spec.CstorConfig.CStorCSI.CStorCSIController.ImageTag == "" {
			p.ObservedOpenEBS.Spec.CstorConfig.CStorCSI.CStorCSIController.ImageTag = p.ObservedOpenEBS.Spec.Version
		}
		p.ObservedOpenEBS.Spec.CstorConfig.CStorCSI.CStorCSIController.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
			"cstor-csi-driver:" + p.ObservedOpenEBS.Spec.CstorConfig.CStorCSI.CStorCSIController.ImageTag
	}

	// Set the default values for cstor csi node daemonset in configuration.
	if p.ObservedOpenEBS.Spec.CstorConfig.CStorCSI.CStorCSINode.Enabled == nil {
		p.ObservedOpenEBS.Spec.CstorConfig.CStorCSI.CStorCSINode.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.CstorConfig.CStorCSI.CStorCSINode.Enabled = true
	}

	if *p.ObservedOpenEBS.Spec.CstorConfig.CStorCSI.CStorCSINode.Enabled == true {
		if p.ObservedOpenEBS.Spec.CstorConfig.CStorCSI.CStorCSINode.ImageTag == "" {
			p.ObservedOpenEBS.Spec.CstorConfig.CStorCSI.CStorCSINode.ImageTag = p.ObservedOpenEBS.Spec.Version
		}
		p.ObservedOpenEBS.Spec.CstorConfig.CStorCSI.CStorCSINode.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
			"cstor-csi-driver:" + p.ObservedOpenEBS.Spec.CstorConfig.CStorCSI.CStorCSINode.ImageTag

		if p.ObservedOpenEBS.Spec.CstorConfig.CStorCSI.CStorCSINode.ISCSIPath == "" {
			p.ObservedOpenEBS.Spec.CstorConfig.CStorCSI.CStorCSINode.ISCSIPath = "/sbin/iscsiadm"
		}
	}

	// form the cstor-pool-manager image(CSPI_MGMT)
	if p.ObservedOpenEBS.Spec.CstorConfig.CSPIMgmt.ImageTag == "" {
		p.ObservedOpenEBS.Spec.CstorConfig.CSPIMgmt.ImageTag = p.ObservedOpenEBS.Spec.Version
	}
	p.ObservedOpenEBS.Spec.CstorConfig.CSPIMgmt.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"cstor-pool-manager:" + p.ObservedOpenEBS.Spec.CstorConfig.CSPIMgmt.ImageTag

	// form the cstor-volume-manager image
	if p.ObservedOpenEBS.Spec.CstorConfig.VolumeManager.ImageTag == "" {
		p.ObservedOpenEBS.Spec.CstorConfig.VolumeManager.ImageTag = p.ObservedOpenEBS.Spec.Version
	}
	p.ObservedOpenEBS.Spec.CstorConfig.VolumeManager.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"cstor-volume-manager:" + p.ObservedOpenEBS.Spec.CstorConfig.VolumeManager.ImageTag

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

	return nil
}

// updateOpenEBSCStorCSINode updates the values of openebs-cstor-csi-node daemonset as per given configuration.
func (p *Planner) updateOpenEBSCStorCSINode(daemonset *unstructured.Unstructured) error {
	daemonset.SetNamespace(types.NamespaceKubeSystem)
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
				p.ObservedOpenEBS.Spec.CstorConfig.CStorCSI.CStorCSINode.ISCSIPath, "spec", "hostPath", "path")
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
				p.ObservedOpenEBS.Spec.CstorConfig.CStorCSI.CStorCSINode.ISCSIPath, "spec", "mountPath")
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
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.CstorConfig.CStorCSI.CStorCSINode.Image,
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

// updateOpenEBSCStorCSIController updates the values of openebs-cstor-csi-controller statefulset as per given configuration.
func (p *Planner) updateOpenEBSCStorCSIController(statefulset *unstructured.Unstructured) error {
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
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.CstorConfig.CStorCSI.CStorCSIController.Image,
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
				p.ObservedOpenEBS.Spec.CstorConfig.VolumeManager.Image, "spec", "value")
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
