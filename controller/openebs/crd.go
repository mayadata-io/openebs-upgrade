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

// getDesiredCustomResourceDefinition updates the customresourcedefinition manifest as per the given configuration.
func (p *Planner) getDesiredCustomResourceDefinition(crd *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	var err error
	switch crd.GetName() {
	case types.CSPCCRDNameKey:
		err = p.updateCSPCCRD(crd)
	case types.CSINodeInfoCRDNameKey:
		err = p.updateCSINodeInfoCRD(crd)
	case types.CSIVolumeCRDNameKey:
		err = p.updateCSIVolumeCRD(crd)
	case types.VolumeSnapshotCRDNameKey:
		err = p.updateVolumeSnapshotCRD(crd)
	case types.VolumeSnapshotClassCRDNameKey:
		err = p.updateVolumeSnapshotClassCRD(crd)
	case types.VolumeSnapshotContentCRDNameKey:
		err = p.updateVolumeSnapshotContentCRD(crd)
	}
	if err != nil {
		return crd, err
	}
	// create annotations that refers to the instance which
	// triggered creation of this CustomResourceDefinition
	crd.SetAnnotations(
		map[string]string{
			types.AnnKeyOpenEBSUID: string(p.ObservedOpenEBS.GetUID()),
		},
	)
	return crd, nil
}

// updateCSINodeInfoCRD updates the CSI node info CRD manifest as per the
// reconcile.ObservedOpenEBS values.
func (p *Planner) updateCSINodeInfoCRD(crd *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := crd.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for CSI node info CRD
	// 1. openebs-upgrade.dao.mayadata.io/component-type: customresourcedefinition
	// 2. openebs-upgrade.dao.mayadata.io/component-group: csi
	// 3. openebs-upgrade.dao.mayadata.io/component-name: csinodeinfos.csi.storage.k8s.io
	desiredLabels[types.OpenEBSComponentTypeLabelKey] =
		types.OpenEBSCRDComponentTypeLabelValue
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.CSIComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.CSINodeInfoCRDNameKey
	// set the desired labels
	crd.SetLabels(desiredLabels)

	return nil
}

// updateCSIVolumeCRD updates the CSI volume CRD manifest as per the
// reconcile.ObservedOpenEBS values.
func (p *Planner) updateCSIVolumeCRD(crd *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := crd.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for CSI volume CRD
	// 1. openebs-upgrade.dao.mayadata.io/component-type: customresourcedefinition
	// 2. openebs-upgrade.dao.mayadata.io/component-group: csi
	// 3. openebs-upgrade.dao.mayadata.io/component-name: csivolumes.openebs.io
	desiredLabels[types.OpenEBSComponentTypeLabelKey] =
		types.OpenEBSCRDComponentTypeLabelValue
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.CSIComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.CSIVolumeCRDNameKey
	// set the desired labels
	crd.SetLabels(desiredLabels)

	return nil
}

// updateVolumeSnapshotCRD updates the volume snapshot CRD manifest as per the
// reconcile.ObservedOpenEBS values.
func (p *Planner) updateVolumeSnapshotCRD(crd *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := crd.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for volume snapshot CRD
	// 1. openebs-upgrade.dao.mayadata.io/component-type: customresourcedefinition
	// 2. openebs-upgrade.dao.mayadata.io/component-name: volumesnapshots.snapshot.storage.k8s.io
	desiredLabels[types.OpenEBSComponentTypeLabelKey] =
		types.OpenEBSCRDComponentTypeLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.VolumeSnapshotCRDNameKey
	// set the desired labels
	crd.SetLabels(desiredLabels)

	return nil
}

// updateVolumeSnapshotClassCRD updates the volume snapshot class CRD manifest as per the
// reconcile.ObservedOpenEBS values.
func (p *Planner) updateVolumeSnapshotClassCRD(crd *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := crd.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for volume snapshot class CRD
	// 1. openebs-upgrade.dao.mayadata.io/component-type: customresourcedefinition
	// 2. openebs-upgrade.dao.mayadata.io/component-name: volumesnapshotclasses.snapshot.storage.k8s.io
	desiredLabels[types.OpenEBSComponentTypeLabelKey] =
		types.OpenEBSCRDComponentTypeLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.VolumeSnapshotClassCRDNameKey
	// set the desired labels
	crd.SetLabels(desiredLabels)

	return nil
}

// updateVolumeSnapshotContentCRD updates the volume snapshot content CRD manifest as per the
// reconcile.ObservedOpenEBS values.
func (p *Planner) updateVolumeSnapshotContentCRD(crd *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := crd.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for volume snapshot content CRD
	// 1. openebs-upgrade.dao.mayadata.io/component-type: customresourcedefinition
	// 2. openebs-upgrade.dao.mayadata.io/component-name: volumesnapshotcontents.snapshot.storage.k8s.io
	desiredLabels[types.OpenEBSComponentTypeLabelKey] =
		types.OpenEBSCRDComponentTypeLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.VolumeSnapshotContentCRDNameKey
	// set the desired labels
	crd.SetLabels(desiredLabels)

	return nil
}

// updateCSPCCRD updates the CSPC CRD manifest as per the reconcile.ObservedOpenEBS values.
func (p *Planner) updateCSPCCRD(crd *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := crd.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for CSPC CRD
	// 1. openebs-upgrade.dao.mayadata.io/component-type: customresourcedefinition
	// 2. openebs-upgrade.dao.mayadata.io/component-group: cspc
	// 3. openebs-upgrade.dao.mayadata.io/component-name: cstorpoolclusters.cstor.openebs.io
	desiredLabels[types.OpenEBSComponentTypeLabelKey] =
		types.OpenEBSCRDComponentTypeLabelValue
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.CSPCComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.CSPCCRDNameKey
	// set the desired labels
	crd.SetLabels(desiredLabels)

	return nil
}
