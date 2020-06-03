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

// getDesiredNamespace updates the namespace manifest as per the given configuration
// in OpenEBS CR.
func (p *Planner) getDesiredNamespace(namespace *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	namespace.SetName(p.ObservedOpenEBS.Namespace)
	// create annotations that refers to the instance which
	// triggered creation of this namespace
	namespace.SetAnnotations(
		map[string]string{
			types.AnnKeyOpenEBSUID: string(p.ObservedOpenEBS.GetUID()),
		},
	)

	return namespace, nil
}

// getDesiredServiceAccount updates the service account manifest as per the
// given configuration in OpenEBS CR.
func (p *Planner) getDesiredServiceAccount(sa *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	var err error
	sa.SetNamespace(p.ObservedOpenEBS.Namespace)
	// create annotations that refers to the instance which
	// triggered creation of this ServiceAccount
	sa.SetAnnotations(
		map[string]string{
			types.AnnKeyOpenEBSUID: string(p.ObservedOpenEBS.GetUID()),
		},
	)

	switch sa.GetName() {
	case types.OpenEBSMayaOperatorSANameKey:
		err = p.updateOpenEBSServiceAccount(sa)
	case types.CStorCSIControllerSANameKey:
		// Overwrite the namespace to kube-system for csi based components.
		// Note: csi based components will be installed only in kube-system namespace only.
		sa.SetNamespace(types.NamespaceKubeSystem)
		err = p.updateCStorCSIControllerServiceAccount(sa)
	case types.CStorCSINodeSANameKey:
		// Overwrite the namespace to kube-system for csi based components.
		// Note: csi based components will be installed only in kube-system namespace only.
		sa.SetNamespace(types.NamespaceKubeSystem)
		err = p.updateCStorCSINodeServiceAccount(sa)
	case types.MoacSANameKey:
		err = p.updateMoacServiceAccount(sa)
	}
	if err != nil {
		return sa, err
	}

	return sa, nil
}

// updateOpenEBSServiceAccount updates the openebs-maya-operator service account
// structure as per the provided values otherwise default values.
func (p *Planner) updateOpenEBSServiceAccount(sa *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := sa.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Set some component specific labels in order to identify specific components.
	// These labels will be only set by openebs-upgrade and will help the end-users
	// identify a particular or a set of OpenEBS components.
	//
	// Component specific labels for openebs-maya-operator service account:
	// 1. openebs-upgrade.dao.mayadata.io/component-name: openebs-maya-operator
	desiredLabels[types.OpenEBSComponentNameLabelKey] =
		types.OpenEBSSAComponentNameLabelValue

	// set the desired labels
	sa.SetLabels(desiredLabels)

	return nil
}

// updateMoacServiceAccount updates the moac service account
// structure as per the provided values otherwise default values.
func (p *Planner) updateMoacServiceAccount(sa *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := sa.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Set some component specific labels in order to identify specific components.
	// These labels will be only set by openebs-upgrade and will help the end-users
	// identify a particular or a set of OpenEBS components.
	//
	// Component specific labels for moac service account:
	// 1. openebs-upgrade.dao.mayadata.io/component-group: mayastor
	// 2. openebs-upgrade.dao.mayadata.io/component-name: moac
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSMayastorComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] =
		types.MoacSANameKey

	// set the desired labels
	sa.SetLabels(desiredLabels)

	return nil
}

// updateCStorCSIControllerServiceAccount updates the CStor CSI controller service account
// structure as per the provided values otherwise default values.
func (p *Planner) updateCStorCSIControllerServiceAccount(sa *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := sa.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Set some component specific labels in order to identify specific components.
	// These labels will be only set by openebs-upgrade and will help the end-users
	// identify a particular or a set of OpenEBS components.
	//
	// Component specific labels for openebs-maya-operator service account:
	// 1. openebs-upgrade.dao.mayadata.io/component-group: cstor-csi
	// 2. openebs-upgrade.dao.mayadata.io/component-name: cstor-csi-controller
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSCStorCSIComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] =
		types.CStorCSICtrlSAComponentNameLabelValue

	// set the desired labels
	sa.SetLabels(desiredLabels)

	return nil
}

// updateCStorCSINodeServiceAccount updates the CStor CSI node service account
// structure as per the provided values otherwise default values.
func (p *Planner) updateCStorCSINodeServiceAccount(sa *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := sa.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Set some component specific labels in order to identify specific components.
	// These labels will be only set by openebs-upgrade and will help the end-users
	// identify a particular or a set of OpenEBS components.
	//
	// Component specific labels for openebs-maya-operator service account:
	// 1. openebs-upgrade.dao.mayadata.io/component-group: cstor-csi
	// 2. openebs-upgrade.dao.mayadata.io/component-name: cstor-csi-node
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSCStorCSIComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] =
		types.CStorCSINodeSAComponentNameLabelValue

	// set the desired labels
	sa.SetLabels(desiredLabels)

	return nil
}

// getDesiredClusterRole updates the cluster role manifest as per the
// given configuration in OpenEBS CR.
func (p *Planner) getDesiredClusterRole(cr *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	var err error

	switch cr.GetName() {
	case types.OpenEBSMayaOperatorRoleNameKey:
		err = p.updateOpenEBSClusterRole(cr)
	case types.CStorCSISnapshottterRoleNameKey:
		err = p.updateCStorCSISnapshotterRole(cr)
	case types.CStorCSIProvisionerRoleNameKey:
		err = p.updateCStorCSIProvisionerRole(cr)
	case types.CStorCSIAttacherRoleNameKey:
		err = p.updateCStorCSIAttacherRole(cr)
	case types.CStorCSIClusterRegistrarRoleNameKey:
		err = p.updateCStorCSIClusterRegistrarRole(cr)
	case types.CStorCSIRegistrarRoleNameKey:
		err = p.updateCStorCSIRegistrarRole(cr)
	case types.MoacClusterRoleNameKey:
		err = p.updateMoacClusterRole(cr)
	}
	if err != nil {
		return cr, err
	}
	// create annotations that refers to the instance which
	// triggered creation of this ClusterRole
	cr.SetAnnotations(
		map[string]string{
			types.AnnKeyOpenEBSUID: string(p.ObservedOpenEBS.GetUID()),
		},
	)
	return cr, nil
}

// updateOpenEBSClusterRole updates the OpenEBS maya-operator cluster role
// structure as per the provided values otherwise default values.
func (p *Planner) updateOpenEBSClusterRole(sa *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := sa.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Set some component specific labels in order to identify specific components.
	// These labels will be only set by openebs-upgrade and will help the end-users
	// identify a particular or a set of OpenEBS components.
	//
	// 1. openebs-upgrade.dao.mayadata.io/component-name: openebs-maya-operator
	desiredLabels[types.OpenEBSComponentNameLabelKey] =
		types.OpenEBSRoleComponentNameLabelValue

	// set the desired labels
	sa.SetLabels(desiredLabels)

	return nil
}

// updateMoacClusterRole updates the moac cluster role
// structure as per the provided values otherwise default values.
func (p *Planner) updateMoacClusterRole(sa *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := sa.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Set some component specific labels in order to identify specific components.
	// These labels will be only set by openebs-upgrade and will help the end-users
	// identify a particular or a set of OpenEBS components.
	//
	// Component specific labels for moac cluster role:
	// 1. openebs-upgrade.dao.mayadata.io/component-group: mayastor
	// 2. openebs-upgrade.dao.mayadata.io/component-name: moac
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSMayastorComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] =
		types.MoacClusterRoleNameKey

	// set the desired labels
	sa.SetLabels(desiredLabels)

	return nil
}

// updateCStorCSISnapshotterRole updates the CStor CSI snapshotter cluster role
// structure as per the provided values otherwise default values.
func (p *Planner) updateCStorCSISnapshotterRole(sa *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := sa.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Set some component specific labels in order to identify specific components.
	// These labels will be only set by openebs-upgrade and will help the end-users
	// identify a particular or a set of OpenEBS components.
	//
	// Component specific labels for cstor CSI snapshotter cluster role:
	// 1. openebs-upgrade.dao.mayadata.io/component-group: cstor-csi
	// 2. openebs-upgrade.dao.mayadata.io/component-name: openebs-cstor-csi-snapshotter-role
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSCStorCSIComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] =
		types.CStorCSISnapshottterRoleNameKey

	// set the desired labels
	sa.SetLabels(desiredLabels)

	return nil
}

// updateCStorCSIProvisionerRole updates the CStor CSI provisioner cluster role
// structure as per the provided values otherwise default values.
func (p *Planner) updateCStorCSIProvisionerRole(sa *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := sa.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Set some component specific labels in order to identify specific components.
	// These labels will be only set by openebs-upgrade and will help the end-users
	// identify a particular or a set of OpenEBS components.
	//
	// Component specific labels for cstor CSI provisioner cluster role:
	// 1. openebs-upgrade.dao.mayadata.io/component-group: cstor-csi
	// 2. openebs-upgrade.dao.mayadata.io/component-name: openebs-cstor-csi-provisioner-role
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSCStorCSIComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] =
		types.CStorCSIProvisionerRoleNameKey

	// set the desired labels
	sa.SetLabels(desiredLabels)

	return nil
}

// updateCStorCSIAttacherRole updates the CStor CSI attacher cluster role
// structure as per the provided values otherwise default values.
func (p *Planner) updateCStorCSIAttacherRole(sa *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := sa.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Set some component specific labels in order to identify specific components.
	// These labels will be only set by openebs-upgrade and will help the end-users
	// identify a particular or a set of OpenEBS components.
	//
	// Component specific labels for openebs-cstor-csi-attacher-role cluster role:
	// 1. openebs-upgrade.dao.mayadata.io/component-group: cstor-csi
	// 2. openebs-upgrade.dao.mayadata.io/component-name: openebs-cstor-csi-attacher-role
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSCStorCSIComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] =
		types.CStorCSIAttacherRoleNameKey

	// set the desired labels
	sa.SetLabels(desiredLabels)

	return nil
}

// updateCStorCSIClusterRegistrarRole updates the CStor CSI cluster registrar cluster role
// structure as per the provided values otherwise default values.
func (p *Planner) updateCStorCSIClusterRegistrarRole(sa *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := sa.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Set some component specific labels in order to identify specific components.
	// These labels will be only set by openebs-upgrade and will help the end-users
	// identify a particular or a set of OpenEBS components.
	//
	// Component specific labels for openebs-cstor-csi-cluster-registrar-role cluster role:
	// 1. openebs-upgrade.dao.mayadata.io/component-group: cstor-csi
	// 2. openebs-upgrade.dao.mayadata.io/component-name: openebs-cstor-csi-cluster-registrar-role
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSCStorCSIComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] =
		types.CStorCSIClusterRegistrarRoleNameKey

	// set the desired labels
	sa.SetLabels(desiredLabels)

	return nil
}

// updateCStorCSIRegistrarRole updates the CStor CSI registrar cluster role
// structure as per the provided values otherwise default values.
func (p *Planner) updateCStorCSIRegistrarRole(sa *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := sa.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Set some component specific labels in order to identify specific components.
	// These labels will be only set by openebs-upgrade and will help the end-users
	// identify a particular or a set of OpenEBS components.
	//
	// Component specific labels for openebs-cstor-csi-registrar-role cluster role:
	// 1. openebs-upgrade.dao.mayadata.io/component-group: cstor-csi
	// 2. openebs-upgrade.dao.mayadata.io/component-name: openebs-cstor-csi-registrar-role
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSCStorCSIComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] =
		types.CStorCSIRegistrarRoleNameKey

	// set the desired labels
	sa.SetLabels(desiredLabels)

	return nil
}

// getDesiredClusterRoleBinding updates the clusterRoleBinding manifest as per the
// given configuration in OpenEBS CR.
func (p *Planner) getDesiredClusterRoleBinding(crb *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	var err error

	switch crb.GetName() {
	case types.OpenEBSMayaOperatorBindingNameKey:
		err = p.updateOpenEBSClusterRoleBinding(crb)
	case types.CStorCSISnapshottterBindingNameKey:
		err = p.updateCStorCSISnapshotterBinding(crb)
	case types.CStorCSIProvisionerBindingNameKey:
		err = p.updateCStorCSIProvisionerBinding(crb)
	case types.CStorCSIAttacherBindingNameKey:
		err = p.updateCStorCSIAttacherBinding(crb)
	case types.CStorCSIClusterRegistrarBindingNameKey:
		err = p.updateCStorCSIClusterRegistrarBinding(crb)
	case types.CStorCSIRegistrarBindingNameKey:
		err = p.updateCStorCSIRegistrarBinding(crb)
	case types.MoacClusterRoleBindingNameKey:
		err = p.updateMoacClusterRoleBinding(crb)
	}
	if err != nil {
		return crb, err
	}
	setNamespaceOfEachSubject := func(obj *unstructured.Unstructured) error {
		err := unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Namespace, "spec", "namespace")
		if err != nil {
			return err
		}

		// get the name of the subject i.e service account name
		objName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}
		// Overwrite the namespace to kube-system for csi based components.
		// Note: csi based components will be installed only in kube-system namespace only.
		if objName == types.CStorCSINodeSANameKey || objName == types.CStorCSIControllerSANameKey {
			err := unstructured.SetNestedField(obj.Object, types.NamespaceKubeSystem, "spec", "namespace")
			if err != nil {
				return err
			}
		}

		return nil
	}
	crbSubjects, _, err := unstruct.GetSlice(crb, "subjects")
	if err != nil {
		return crb, err
	}
	err = unstruct.SliceIterator(crbSubjects).ForEachUpdate(setNamespaceOfEachSubject)
	if err != nil {
		return crb, err
	}
	err = unstructured.SetNestedSlice(crb.Object, crbSubjects, "subjects")
	if err != nil {
		return crb, err
	}

	// create annotations that refers to the instance which
	// triggered creation of this ClusterRoleBinding
	crb.SetAnnotations(
		map[string]string{
			types.AnnKeyOpenEBSUID: string(p.ObservedOpenEBS.GetUID()),
		},
	)
	return crb, nil
}

// updateOpenEBSClusterRoleBinding updates the OpenEBS maya-operator cluster role
// binding structure as per the provided values otherwise default values.
func (p *Planner) updateOpenEBSClusterRoleBinding(sa *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := sa.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Set some component specific labels in order to identify specific components.
	// These labels will be only set by openebs-upgrade and will help the end-users
	// identify a particular or a set of OpenEBS components.
	//
	// Component specific labels for openebs-maya-operator cluster role binding:
	// 1. openebs-upgrade.dao.mayadata.io/component-name: openebs-maya-operator
	desiredLabels[types.OpenEBSComponentNameLabelKey] =
		types.OpenEBSRoleBindingComponentNameLabelValue

	// set the desired labels
	sa.SetLabels(desiredLabels)

	return nil
}

// updateMoacClusterRoleBinding updates the moac cluster role
// binding structure as per the provided values otherwise default values.
func (p *Planner) updateMoacClusterRoleBinding(sa *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := sa.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Set some component specific labels in order to identify specific components.
	// These labels will be only set by openebs-upgrade and will help the end-users
	// identify a particular or a set of OpenEBS components.
	//
	// Component specific labels for moac cluster role binding:
	// 1. openebs-upgrade.dao.mayadata.io/component-group: mayastor
	// 2. openebs-upgrade.dao.mayadata.io/component-name: moac
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSMayastorComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] =
		types.MoacClusterRoleBindingNameKey

	// set the desired labels
	sa.SetLabels(desiredLabels)

	return nil
}

// updateCStorCSISnapshotterBinding updates the CStor CSI snapshotter cluster role
// binding structure as per the provided values otherwise default values.
func (p *Planner) updateCStorCSISnapshotterBinding(sa *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := sa.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Set some component specific labels in order to identify specific components.
	// These labels will be only set by openebs-upgrade and will help the end-users
	// identify a particular or a set of OpenEBS components.
	//
	// Component specific labels for openebs-cstor-csi-snapshotter-binding cluster role binding:
	// 1. openebs-upgrade.dao.mayadata.io/component-subtype: cstor-csi
	// 2. openebs-upgrade.dao.mayadata.io/component-name: openebs-cstor-csi-snapshotter-binding
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSCStorCSIComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] =
		types.CStorCSISnapshottterBindingNameKey

	// set the desired labels
	sa.SetLabels(desiredLabels)

	return nil
}

// updateCStorCSIProvisionerBinding updates the CStor CSI provisioner cluster role
// binding structure as per the provided values otherwise default values.
func (p *Planner) updateCStorCSIProvisionerBinding(sa *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := sa.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Set some component specific labels in order to identify specific components.
	// These labels will be only set by openebs-upgrade and will help the end-users
	// identify a particular or a set of OpenEBS components.
	//
	// Component specific labels for openebs-cstor-csi-provisioner-binding cluster role binding:
	// 1. openebs-upgrade.dao.mayadata.io/component-subtype: cstor-csi
	// 2. openebs-upgrade.dao.mayadata.io/component-name: openebs-cstor-csi-provisioner-binding
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSCStorCSIComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] =
		types.CStorCSIProvisionerBindingNameKey

	// set the desired labels
	sa.SetLabels(desiredLabels)

	return nil
}

// updateCStorCSIAttacherBinding updates the CStor CSI attacher cluster role binding
// structure as per the provided values otherwise default values.
func (p *Planner) updateCStorCSIAttacherBinding(sa *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := sa.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Set some component specific labels in order to identify specific components.
	// These labels will be only set by openebs-upgrade and will help the end-users
	// identify a particular or a set of OpenEBS components.
	//
	// Component specific labels for openebs-cstor-csi-attacher-binding cluster role binding:
	// 1. openebs-upgrade.dao.mayadata.io/component-subtype: cstor-csi
	// 2. openebs-upgrade.dao.mayadata.io/component-name: openebs-cstor-csi-attacher-binding
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSCStorCSIComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] =
		types.CStorCSIAttacherBindingNameKey

	// set the desired labels
	sa.SetLabels(desiredLabels)

	return nil
}

// updateCStorCSIClusterRegistrarBinding updates the CStor CSI cluster registrar cluster role
// binding structure as per the provided values otherwise default values.
func (p *Planner) updateCStorCSIClusterRegistrarBinding(sa *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := sa.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Set some component specific labels in order to identify specific components.
	// These labels will be only set by openebs-upgrade and will help the end-users
	// identify a particular or a set of OpenEBS components.
	//
	// Component specific labels for openebs-cstor-csi-cluster-registrar-binding cluster role:
	// 1. openebs-upgrade.dao.mayadata.io/component-subtype: cstor-csi
	// 2. openebs-upgrade.dao.mayadata.io/component-name: openebs-cstor-csi-cluster-registrar-binding
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSCStorCSIComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] =
		types.CStorCSIClusterRegistrarBindingNameKey

	// set the desired labels
	sa.SetLabels(desiredLabels)

	return nil
}

// updateCStorCSIRegistrarBinding updates the CStor CSI registrar cluster role binding
// structure as per the provided values otherwise default values.
func (p *Planner) updateCStorCSIRegistrarBinding(sa *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := sa.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Set some component specific labels in order to identify specific components.
	// These labels will be only set by openebs-upgrade and will help the end-users
	// identify a particular or a set of OpenEBS components.
	//
	// Component specific labels for openebs-cstor-csi-registrar-binding cluster role binding:
	// 1. openebs-upgrade.dao.mayadata.io/component-subtype: cstor-csi
	// 2. openebs-upgrade.dao.mayadata.io/component-name: openebs-cstor-csi-registrar-binding
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSCStorCSIComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] =
		types.CStorCSIRegistrarBindingNameKey

	// set the desired labels
	sa.SetLabels(desiredLabels)

	return nil
}
