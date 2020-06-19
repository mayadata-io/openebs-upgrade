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
	case types.CSPCCRDV1alpha1NameKey:
		err = p.updateCSPCCRDV1alpha1(crd)
	case types.CSPCCRDV1NameKey:
		err = p.updateCSPCCRDV1(crd)
	case types.CSPICRDV1NameKey:
		err = p.updateCSPICRDV1(crd)
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
	case types.CStorVolumeAttachmentsCRDNameKey:
		err = p.updateCStorVolumeAttachmentsCRD(crd)
	case types.CStorVolumesCRDV1NameKey:
		err = p.updateCStorVolumesCRDV1(crd)
	case types.CStorVolumeConfigsCRDV1NameKey:
		err = p.updateCStorVolumeConfigsCRDV1(crd)
	case types.CStorVolumeReplicasCRDV1NameKey:
		err = p.updateCStorVolumeReplicasCRDV1(crd)
	case types.CStorVolumePoliciesCRDV1NameKey:
		err = p.updateCStorVolumePoliciesCRDV1(crd)
	case types.CStorBackupsCRDV1alpha1NameKey:
		err = p.updateCStorBackupCRDV1alpha1(crd)
	case types.CStorCompletedBackupsCRDV1alpha1NameKey:
		err = p.updateCStorCompletedBackupCRDV1alpha1(crd)
	case types.CStorRestoresCRDV1alpha1NameKey:
		err = p.updateCStorRestoresCRDV1alpha1(crd)
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
	// 1. openebs-upgrade.dao.mayadata.io/component-group: csi
	// 2. openebs-upgrade.dao.mayadata.io/component-name: csinodeinfos.csi.storage.k8s.io
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
	// 1. openebs-upgrade.dao.mayadata.io/component-group: csi
	// 2. openebs-upgrade.dao.mayadata.io/component-name: csivolumes.openebs.io
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
	// 1. openebs-upgrade.dao.mayadata.io/component-name: volumesnapshots.snapshot.storage.k8s.io
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
	// 1. openebs-upgrade.dao.mayadata.io/component-name: volumesnapshotclasses.snapshot.storage.k8s.io
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
	// 1. openebs-upgrade.dao.mayadata.io/component-name: volumesnapshotcontents.snapshot.storage.k8s.io
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.VolumeSnapshotContentCRDNameKey
	// set the desired labels
	crd.SetLabels(desiredLabels)

	return nil
}

// updateCStorVolumeAttachmentsCRD updates the CStor volume attachments CRD manifest as per the
// reconcile.ObservedOpenEBS values.
func (p *Planner) updateCStorVolumeAttachmentsCRD(crd *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := crd.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for CStor volume attachments CRD
	// 1. openebs-upgrade.dao.mayadata.io/component-name: cstorvolumeattachments.cstor.openebs.io
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.CStorVolumeAttachmentsCRDNameKey
	// set the desired labels
	crd.SetLabels(desiredLabels)

	return nil
}

// updateCSPCCRDV1 updates the CSPC V1 CRD manifest as per the reconcile.ObservedOpenEBS values.
func (p *Planner) updateCSPCCRDV1(crd *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := crd.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for CSPC CRD
	// 1. openebs-upgrade.dao.mayadata.io/component-group: cspc
	// 2. openebs-upgrade.dao.mayadata.io/component-name: cstorpoolclusters.cstor.openebs.io
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.CSPCComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.CSPCCRDV1NameKey
	// set the desired labels
	crd.SetLabels(desiredLabels)

	return nil
}

// updateCSPCCRDV1alpha1 updates the CSPC V1alpha1 CRD manifest as per the reconcile.ObservedOpenEBS values.
func (p *Planner) updateCSPCCRDV1alpha1(crd *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := crd.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for CSPC CRD
	// 1. openebs-upgrade.dao.mayadata.io/component-group: cspc
	// 2. openebs-upgrade.dao.mayadata.io/component-name: cstorpoolclusters.openebs.io
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.CSPCComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.CSPCCRDV1alpha1NameKey
	// set the desired labels
	crd.SetLabels(desiredLabels)

	return nil
}

// updateCSPICRDV1 updates the CSPI V1 CRD manifest as per the reconcile.ObservedOpenEBS values.
func (p *Planner) updateCSPICRDV1(crd *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := crd.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for CSPI CRD
	// 1. openebs-upgrade.dao.mayadata.io/component-group: cspi
	// 2. openebs-upgrade.dao.mayadata.io/component-name: cstorpoolinstances.cstor.openebs.io
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.CSPIComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.CSPICRDV1NameKey
	// set the desired labels
	crd.SetLabels(desiredLabels)

	return nil
}

// updateCStorVolumesCRDV1 updates the CStor volumes CRD(V1) manifest as per the
// reconcile.ObservedOpenEBS values.
func (p *Planner) updateCStorVolumesCRDV1(crd *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := crd.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for CStor volumes CRD
	// 1. openebs-upgrade.dao.mayadata.io/component-name: cstorvolumes.cstor.openebs.io
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.CStorVolumesCRDV1NameKey
	// set the desired labels
	crd.SetLabels(desiredLabels)

	return nil
}

// updateCStorVolumeConfigsCRDV1 updates the CStor volume configs CRD(V1) manifest as per the
// reconcile.ObservedOpenEBS values.
func (p *Planner) updateCStorVolumeConfigsCRDV1(crd *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := crd.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for CStor volume configs CRD
	// 1. openebs-upgrade.dao.mayadata.io/component-name: cstorvolumeconfigs.cstor.openebs.io
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.CStorVolumeConfigsCRDV1NameKey
	// set the desired labels
	crd.SetLabels(desiredLabels)

	return nil
}

// updateCStorVolumePoliciesCRDV1 updates the CStor volume policies CRD(V1) manifest as per the
// reconcile.ObservedOpenEBS values.
func (p *Planner) updateCStorVolumePoliciesCRDV1(crd *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := crd.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for CStor volume policies CRD
	// 1. openebs-upgrade.dao.mayadata.io/component-name: cstorvolumepolicies.cstor.openebs.io
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.CStorVolumePoliciesCRDV1NameKey
	// set the desired labels
	crd.SetLabels(desiredLabels)

	return nil
}

// updateCStorVolumeReplicasCRDV1 updates the CStor volume replicas CRD(V1) manifest as per the
// reconcile.ObservedOpenEBS values.
func (p *Planner) updateCStorVolumeReplicasCRDV1(crd *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := crd.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for CStor volume replicas CRD
	// 1. openebs-upgrade.dao.mayadata.io/component-name: cstorvolumereplicas.cstor.openebs.io
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.CStorVolumeReplicasCRDV1NameKey
	// set the desired labels
	crd.SetLabels(desiredLabels)

	return nil
}

// updateCStorBackupsCRDV1alpha1 updates the CStor backup CRD(v1alpha1) manifest as per the
// reconcile.ObservedOpenEBS values.
func (p *Planner) updateCStorBackupCRDV1alpha1(crd *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := crd.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for CStor backup CRD
	// 1. openebs-upgrade.dao.mayadata.io/component-name: cstorbackups.openebs.io
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.CStorBackupsCRDV1alpha1NameKey
	// set the desired labels
	crd.SetLabels(desiredLabels)

	return nil
}

// updateCStorCompletedBackupCRDV1alpha1 updates the CStor completed backup CRD(v1alpha1)
// manifest as per the reconcile.ObservedOpenEBS values.
func (p *Planner) updateCStorCompletedBackupCRDV1alpha1(crd *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := crd.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for CStor completed backups CRD
	// 1. openebs-upgrade.dao.mayadata.io/component-name: cstorcompletedbackups.openebs.io
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.CStorCompletedBackupsCRDV1alpha1NameKey
	// set the desired labels
	crd.SetLabels(desiredLabels)

	return nil
}

// updateCStorRestoresCRDV1alpha1 updates the CStor restores CRD(v1alpha1) manifest as per the
// reconcile.ObservedOpenEBS values.
func (p *Planner) updateCStorRestoresCRDV1alpha1(crd *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := crd.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for CStor restores CRD
	// 1. openebs-upgrade.dao.mayadata.io/component-name: cstorrestores.openebs.io
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.CStorRestoresCRDV1alpha1NameKey
	// set the desired labels
	crd.SetLabels(desiredLabels)

	return nil
}
