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
