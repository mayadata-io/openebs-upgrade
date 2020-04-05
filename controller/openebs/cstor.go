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
	// NamespaceKubeSystem is the value of kube-system namespace
	NamespaceKubeSystem string = "kube-system"
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

	if p.ObservedOpenEBS.Spec.CstorConfig.CStorCSI.Enabled == nil {
		p.ObservedOpenEBS.Spec.CstorConfig.CStorCSI.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.CstorConfig.CStorCSI.Enabled = true
	}

	return nil
}

func (p *Planner) updateOpenEBSCStorCSINode(daemonset *unstructured.Unstructured) error {
	daemonset.SetNamespace(NamespaceKubeSystem)

	containers, err := unstruct.GetNestedSliceOrError(daemonset, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}

	// update the env value of openebs-csi-plugin container
	updateEnv := func(env *unstructured.Unstructured) error {
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
		if containerName == ContainerOpenEBSCSIPluginName {
			envs, _, err := unstruct.GetSlice(obj, "spec", "env")
			if err != nil {
				return err
			}
			err = unstruct.SliceIterator(envs).ForEachUpdate(updateEnv)
			if err != nil {
				return err
			}
			err = unstructured.SetNestedSlice(obj.Object, envs, "spec", "env")
			if err != nil {
				return err
			}
		}

		return nil
	}

	err = unstruct.SliceIterator(containers).ForEachUpdate(updateContainer)
	if err != nil {
		return err
	}

	err = unstructured.SetNestedSlice(daemonset.Object,
		containers, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}

	return nil
}

func (p *Planner) updateOpenEBSCStorCSIController(statefulset *unstructured.Unstructured) error {
	statefulset.SetNamespace(NamespaceKubeSystem)

	containers, err := unstruct.GetNestedSliceOrError(statefulset, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}

	// update the env value of openebs-csi-plugin container
	updateEnv := func(env *unstructured.Unstructured) error {
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
		if containerName == ContainerOpenEBSCSIPluginName {
			envs, _, err := unstruct.GetSlice(obj, "spec", "env")
			if err != nil {
				return err
			}
			err = unstruct.SliceIterator(envs).ForEachUpdate(updateEnv)
			if err != nil {
				return err
			}
			err = unstructured.SetNestedSlice(obj.Object, envs, "spec", "env")
			if err != nil {
				return err
			}
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
