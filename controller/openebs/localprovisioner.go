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
	// DefaultLocalProvisionerReplicaCount is the default replica count for
	// openebs local provisioner.
	DefaultLocalProvisionerReplicaCount int32 = 1
)

// setLocalProvisionerDefaultsIfNotSet sets the default values for local-provisioner
func (p *Planner) setLocalProvisionerDefaultsIfNotSet() error {
	if p.ObservedOpenEBS.Spec.LocalProvisioner == nil {
		p.ObservedOpenEBS.Spec.LocalProvisioner = &types.LocalProvisioner{}
	}
	if p.ObservedOpenEBS.Spec.LocalProvisioner.Enabled == nil {
		p.ObservedOpenEBS.Spec.LocalProvisioner.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.LocalProvisioner.Enabled = true
	}
	// set the name with which localpv-provisioner will be deployed
	if len(p.ObservedOpenEBS.Spec.LocalProvisioner.Name) == 0 {
		p.ObservedOpenEBS.Spec.LocalProvisioner.Name = types.LocalProvisionerNameKey
	}
	if p.ObservedOpenEBS.Spec.LocalProvisioner.ImageTag == "" {
		p.ObservedOpenEBS.Spec.LocalProvisioner.ImageTag = p.ObservedOpenEBS.Spec.Version +
			p.ObservedOpenEBS.Spec.ImageTagSuffix
	}
	p.ObservedOpenEBS.Spec.LocalProvisioner.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"provisioner-localpv:" + p.ObservedOpenEBS.Spec.LocalProvisioner.ImageTag

	if p.ObservedOpenEBS.Spec.LocalProvisioner.Replicas == nil {
		p.ObservedOpenEBS.Spec.LocalProvisioner.Replicas = new(int32)
		*p.ObservedOpenEBS.Spec.LocalProvisioner.Replicas = DefaultLocalProvisionerReplicaCount
	}
	return nil
}

// updateLocalProvisioner updates the localProvisioner structure as per the provided
// values otherwise default values.
func (p *Planner) updateLocalProvisioner(deploy *unstructured.Unstructured) error {
	deploy.SetName(p.ObservedOpenEBS.Spec.LocalProvisioner.Name)
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := deploy.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for openebs-ndm-operator deploy
	// 1. openebs-upgrade.dao.mayadata.io/component-name: openebs-localpv-provisioner
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.LocalProvisionerNameKey
	// set the desired labels
	deploy.SetLabels(desiredLabels)

	// update the daemonset containers
	containers, err := unstruct.GetNestedSliceOrError(deploy, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}
	updateProvisionerHostPathEnv := func(env *unstructured.Unstructured) error {
		envName, _, err := unstructured.NestedString(env.Object, "spec", "name")
		if err != nil {
			return err
		}
		if envName == "OPENEBS_IO_HELPER_IMAGE" {
			err = unstructured.SetNestedField(env.Object, p.ObservedOpenEBS.Spec.Helper.Image,
				"spec", "value")
			if err != nil {
				return err
			}
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
		if containerName == "openebs-provisioner-hostpath" {
			// Set the image of the container.
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.LocalProvisioner.Image,
				"spec", "image")
			if err != nil {
				return err
			}
			err = unstruct.SliceIterator(envs).ForEachUpdate(updateProvisionerHostPathEnv)
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
	err = unstructured.SetNestedSlice(deploy.Object, containers,
		"spec", "template", "spec", "containers")
	if err != nil {
		return err
	}

	return nil
}
