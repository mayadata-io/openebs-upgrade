package openebs

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/types"
)

// getDesiredPriorityClass updates the priorityClass manifest as per the given configuration.
func (p *Planner) getDesiredPriorityClass(pc *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	var err error
	switch pc.GetName() {
	case types.OpenEBSCSIControllerPriorityClassNameKey:
		err = p.updateCSIControllerPriorityClass(pc)
	case types.OpenEBSCSINodePriorityClassNameKey:
		err = p.updateCSINodePriorityClass(pc)
	}
	if err != nil {
		return pc, err
	}
	// create annotations that refers to the instance which
	// triggered creation of this PriorityClass.
	pc.SetAnnotations(
		map[string]string{
			types.AnnKeyOpenEBSUID: string(p.ObservedOpenEBS.GetUID()),
		},
	)
	return pc, nil
}

// updateCSINodePriorityClass updates the CSI node priority class manifest as per the
// reconcile.ObservedOpenEBS values.
func (p *Planner) updateCSINodePriorityClass(crd *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := crd.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for CSI node
	// 1. openebs-upgrade.dao.mayadata.io/component-group: csi
	// 2. openebs-upgrade.dao.mayadata.io/component-name: openebs-csi-node-critical
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.CSIComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.OpenEBSCSINodePriorityClassNameKey
	// set the desired labels
	crd.SetLabels(desiredLabels)

	return nil
}

// updateCSIControllerPriorityClass updates the CSI controller priority class manifest as per the
// reconcile.ObservedOpenEBS values.
func (p *Planner) updateCSIControllerPriorityClass(crd *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := crd.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for CSI controller
	// 1. openebs-upgrade.dao.mayadata.io/component-group: csi
	// 2. openebs-upgrade.dao.mayadata.io/component-name: openebs-csi-controller-critical
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.CSIComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.OpenEBSCSIControllerPriorityClassNameKey
	// set the desired labels
	crd.SetLabels(desiredLabels)

	return nil
}
