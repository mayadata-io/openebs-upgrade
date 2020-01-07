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
	"github.com/mayadata-io/openebs-operator/types"
)

const (
	// DefaultLocalProvisionerReplicaCount is the default replica count for
	// openebs local provisioner.
	DefaultLocalProvisionerReplicaCount int32 = 1
)

// setLocalProvisionerDefaultsIfNotSet sets the default values for local-provisioner
func (r *Reconciler) setLocalProvisionerDefaultsIfNotSet() error {
	if r.OpenEBS.Spec.LocalProvisioner == nil {
		r.OpenEBS.Spec.LocalProvisioner = &types.LocalProvisioner{}
	}
	if r.OpenEBS.Spec.LocalProvisioner.Enabled == "" {
		r.OpenEBS.Spec.LocalProvisioner.Enabled = types.True
	}
	if r.OpenEBS.Spec.LocalProvisioner.ImageTag == "" {
		r.OpenEBS.Spec.LocalProvisioner.ImageTag = r.OpenEBS.Spec.Version
	}
	r.OpenEBS.Spec.LocalProvisioner.Image = r.OpenEBS.Spec.ImagePrefix +
		"provisioner-localpv:" + r.OpenEBS.Spec.LocalProvisioner.ImageTag

	if r.OpenEBS.Spec.LocalProvisioner.Replicas == nil {
		r.OpenEBS.Spec.LocalProvisioner.Replicas = new(int32)
		*r.OpenEBS.Spec.LocalProvisioner.Replicas = DefaultLocalProvisionerReplicaCount
	}
	return nil
}
