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

import "github.com/mayadata-io/openebs-operator/types"

const (
	// DefaultSnapshotOperatorReplicaCount is the default value of replica for
	// snapshot operator.
	DefaultSnapshotOperatorReplicaCount int32 = 1
)

// setSnapshotOperatorDefaultsIfNotSet sets the default values for snapshot
// operator.
func (r *Reconciler) setSnapshotOperatorDefaultsIfNotSet() error {
	if r.OpenEBS.Spec.SnapshotOperator == nil {
		r.OpenEBS.Spec.SnapshotOperator = &types.SnapshotOperator{}
	}
	if r.OpenEBS.Spec.SnapshotOperator.Enabled == "" {
		r.OpenEBS.Spec.SnapshotOperator.Enabled = types.True
	}
	// form the snapshot-provisioner image
	if r.OpenEBS.Spec.SnapshotOperator.Provisioner.ImageTag == "" {
		r.OpenEBS.Spec.SnapshotOperator.Provisioner.ImageTag = r.OpenEBS.Spec.Version
	}
	r.OpenEBS.Spec.SnapshotOperator.Provisioner.Image = r.OpenEBS.Spec.ImagePrefix +
		"snapshot-provisioner:" + r.OpenEBS.Spec.SnapshotOperator.Provisioner.ImageTag

	// form the snapshot-controller image
	if r.OpenEBS.Spec.SnapshotOperator.Controller.ImageTag == "" {
		r.OpenEBS.Spec.SnapshotOperator.Controller.ImageTag = r.OpenEBS.Spec.Version
	}
	r.OpenEBS.Spec.SnapshotOperator.Controller.Image = r.OpenEBS.Spec.ImagePrefix +
		"snapshot-controller:" + r.OpenEBS.Spec.SnapshotOperator.Controller.ImageTag

	if r.OpenEBS.Spec.SnapshotOperator.Replicas == nil {
		r.OpenEBS.Spec.SnapshotOperator.Replicas = new(int32)
		*r.OpenEBS.Spec.SnapshotOperator.Replicas = DefaultSnapshotOperatorReplicaCount
	}
	return nil
}
