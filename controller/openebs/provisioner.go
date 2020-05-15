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
	// DefaultProvisionerReplicaCount is the default replica count for
	// openebs-k8s-provisioner.
	DefaultProvisionerReplicaCount int32 = 1
)

// setProvisionerDefaultsIfNotSet sets the default values for openebs-k8s-provisioner.
func (p *Planner) setProvisionerDefaultsIfNotSet() error {
	if p.ObservedOpenEBS.Spec.Provisioner == nil {
		p.ObservedOpenEBS.Spec.Provisioner = &types.Provisioner{}
	}
	if p.ObservedOpenEBS.Spec.Provisioner.Enabled == nil {
		p.ObservedOpenEBS.Spec.Provisioner.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.Provisioner.Enabled = true
	}
	// set the name with which openebs-provisioner will be deployed
	if len(p.ObservedOpenEBS.Spec.Provisioner.Name) == 0 {
		p.ObservedOpenEBS.Spec.Provisioner.Name = types.ProvisionerNameKey
	}
	if p.ObservedOpenEBS.Spec.Provisioner.ImageTag == "" {
		p.ObservedOpenEBS.Spec.Provisioner.ImageTag = p.ObservedOpenEBS.Spec.Version +
			p.ObservedOpenEBS.Spec.ImageTagSuffix
	}
	// form the image for openebs-k8s-provisioner.
	p.ObservedOpenEBS.Spec.Provisioner.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"openebs-k8s-provisioner:" + p.ObservedOpenEBS.Spec.Provisioner.ImageTag

	if p.ObservedOpenEBS.Spec.Provisioner.Replicas == nil {
		p.ObservedOpenEBS.Spec.Provisioner.Replicas = new(int32)
		*p.ObservedOpenEBS.Spec.Provisioner.Replicas = DefaultProvisionerReplicaCount
	}
	return nil
}

// updateOpenEBSProvisioner updates the openebs-provisioner manifest as per the
// reconcile.ObservedOpenEBS values.
func (p *Planner) updateOpenEBSProvisioner(deploy *unstructured.Unstructured) error {
	deploy.SetName(p.ObservedOpenEBS.Spec.Provisioner.Name)
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := deploy.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for openebs-provisioner deploy
	// 1. openebs-upgrade.dao.mayadata.io/component-name: openebs-provisioner
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.ProvisionerNameKey
	// set the desired labels
	deploy.SetLabels(desiredLabels)

	containers, err := unstruct.GetNestedSliceOrError(deploy, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}
	// update the containers
	updateContainer := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}
		if containerName == types.OpenEBSProvisionerContainerKey {
			// Set the image of the container.
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.Provisioner.Image,
				"spec", "image")
			if err != nil {
				return err
			}
		}
		return nil
	}
	// Update the containers.
	err = unstruct.SliceIterator(containers).ForEachUpdate(updateContainer)
	if err != nil {
		return err
	}
	// Set back the value of the containers.
	err = unstructured.SetNestedSlice(deploy.Object,
		containers, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}

	return nil
}
