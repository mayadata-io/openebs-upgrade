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
)

const (
	// DefaultSnapshotOperatorReplicaCount is the default value of replica for
	// snapshot operator.
	DefaultSnapshotOperatorReplicaCount int32 = 1
)

// setSnapshotOperatorDefaultsIfNotSet sets the default values for snapshot
// operator.
func (p *Planner) setSnapshotOperatorDefaultsIfNotSet() error {
	if p.ObservedOpenEBS.Spec.SnapshotOperator == nil {
		p.ObservedOpenEBS.Spec.SnapshotOperator = &types.SnapshotOperator{}
	}
	if p.ObservedOpenEBS.Spec.SnapshotOperator.Enabled == nil {
		p.ObservedOpenEBS.Spec.SnapshotOperator.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.SnapshotOperator.Enabled = true
	}
	// form the snapshot-provisioner image
	if p.ObservedOpenEBS.Spec.SnapshotOperator.Provisioner.ImageTag == "" {
		p.ObservedOpenEBS.Spec.SnapshotOperator.Provisioner.ImageTag = p.ObservedOpenEBS.Spec.Version
	}
	p.ObservedOpenEBS.Spec.SnapshotOperator.Provisioner.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"snapshot-provisioner:" + p.ObservedOpenEBS.Spec.SnapshotOperator.Provisioner.ImageTag

	// form the snapshot-controller image
	if p.ObservedOpenEBS.Spec.SnapshotOperator.Controller.ImageTag == "" {
		p.ObservedOpenEBS.Spec.SnapshotOperator.Controller.ImageTag = p.ObservedOpenEBS.Spec.Version
	}
	p.ObservedOpenEBS.Spec.SnapshotOperator.Controller.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"snapshot-controller:" + p.ObservedOpenEBS.Spec.SnapshotOperator.Controller.ImageTag

	if p.ObservedOpenEBS.Spec.SnapshotOperator.Replicas == nil {
		p.ObservedOpenEBS.Spec.SnapshotOperator.Replicas = new(int32)
		*p.ObservedOpenEBS.Spec.SnapshotOperator.Replicas = DefaultSnapshotOperatorReplicaCount
	}
	return nil
}

// updateSnapshotOperator updates the openebs-snapshot-operator manifest as per the
// reconcile.ObservedOpenEBS values.
func (p *Planner) updateSnapshotOperator(deploy *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := deploy.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for openebs-snapshot-operator deploy
	// 1. openebs-upgrade.dao.mayadata.io/component-type: deployment
	// 2. openebs-upgrade.dao.mayadata.io/component-name: openebs-snapshot-operator
	desiredLabels[types.OpenEBSComponentTypeLabelKey] =
		types.OpenEBSDeploymentComponentTypeLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.SnapshotOperatorNameKey
	// set the desired labels
	deploy.SetLabels(desiredLabels)

	return nil
}
