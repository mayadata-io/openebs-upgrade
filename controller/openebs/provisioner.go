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
	"fmt"

	"github.com/mayadata-io/openebs-operator/types"
)

const (
	// DefaultProvisionerReplicaCount is the default replica count for
	// openebs-k8s-provisioner.
	DefaultProvisionerReplicaCount int32 = 1
)

// setProvisionerDefaultsIfNotSet sets the default values for openebs-k8s-provisioner.
func (r *Reconciler) setProvisionerDefaultsIfNotSet() error {
	if r.OpenEBS.Spec.Provisioner == nil {
		fmt.Println("nil provisioner")
		r.OpenEBS.Spec.Provisioner = &types.Provisioner{}
	}
	if r.OpenEBS.Spec.Provisioner.Enabled == "" {
		fmt.Println("provisioner not set")
		r.OpenEBS.Spec.Provisioner.Enabled = types.True
	}
	if r.OpenEBS.Spec.Provisioner.ImageTag == "" {
		r.OpenEBS.Spec.Provisioner.ImageTag = r.OpenEBS.Spec.Version
	}
	// form the image for openebs-k8s-provisioner.
	r.OpenEBS.Spec.Provisioner.Image = r.OpenEBS.Spec.ImagePrefix +
		"openebs-k8s-provisioner:" + r.OpenEBS.Spec.Provisioner.ImageTag

	if r.OpenEBS.Spec.Provisioner.Replicas == nil {
		r.OpenEBS.Spec.Provisioner.Replicas = new(int32)
		*r.OpenEBS.Spec.Provisioner.Replicas = DefaultProvisionerReplicaCount
	}
	return nil
}
