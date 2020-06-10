package openebs

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/types"
	"mayadata.io/openebs-upgrade/unstruct"
)

const (
	MayastorVersion010EE    string = "0.1.0-ee"
	DefaultMoacReplicaCount int32  = 1
)

// supportedMayastorVersionForOpenEBSVersion stores the mapping for
// Mayastor to OpenEBS version i.e., a Mayastor version for each of the
// supported OpenEBS versions.
var supportedMayastorVersionForOpenEBSVersion = map[string]string{
	types.OpenEBSVersion1100EE: MayastorVersion010EE,
}

// Set the default values for Mayastor if not already given.
func (p *Planner) setMayastorDefaultsIfNotSet() error {
	if p.ObservedOpenEBS.Spec.MayastorConfig == nil {
		p.ObservedOpenEBS.Spec.MayastorConfig = &types.MayastorConfig{}
	}

	if p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Enabled == nil {
		p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Enabled = false
	}

	if *p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Enabled == true {
		if p.ObservedOpenEBS.Spec.MayastorConfig.Moac.ImageTag == "" {
			if moacVersion, exist :=
				supportedMayastorVersionForOpenEBSVersion[p.ObservedOpenEBS.Spec.Version]; exist {
				p.ObservedOpenEBS.Spec.MayastorConfig.Moac.ImageTag = moacVersion +
					p.ObservedOpenEBS.Spec.ImageTagSuffix
			} else {
				return errors.Errorf("Failed to get moac version for the given OpenEBS version: %s",
					p.ObservedOpenEBS.Spec.Version)
			}

			p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
				"moac:" + p.ObservedOpenEBS.Spec.MayastorConfig.Moac.ImageTag

			if p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Replicas == nil {
				p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Replicas = new(int32)
				*p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Replicas = DefaultMoacReplicaCount
			}
		}
	}

	if p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Enabled == nil {
		p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Enabled = false
	}

	if *p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Enabled == true {
		if p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Mayastor.ImageTag == "" {
			if mayastorVersion, exist :=
				supportedMayastorVersionForOpenEBSVersion[p.ObservedOpenEBS.Spec.Version]; exist {
				p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Mayastor.ImageTag = mayastorVersion +
					p.ObservedOpenEBS.Spec.ImageTagSuffix
			} else {
				return errors.Errorf("Failed to get mayastor version for the given OpenEBS version: %s",
					p.ObservedOpenEBS.Spec.Version)
			}
		}
		if p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.MayastorGRPC.ImageTag == "" {
			if mayastorVersion, exist :=
				supportedMayastorVersionForOpenEBSVersion[p.ObservedOpenEBS.Spec.Version]; exist {
				p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.MayastorGRPC.ImageTag = mayastorVersion +
					p.ObservedOpenEBS.Spec.ImageTagSuffix
			} else {
				return errors.Errorf("Failed to get mayastor grpc version for the given OpenEBS version: %s",
					p.ObservedOpenEBS.Spec.Version)
			}
		}

		p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Mayastor.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
			"mayastor:" + p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Mayastor.ImageTag
		p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.MayastorGRPC.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
			"mayastor-grpc:" + p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.MayastorGRPC.ImageTag
	}

	return nil
}

// updateMoac updates the moac manifest as per the reconcile.ObservedOpenEBS values.
func (p *Planner) updateMoac(deploy *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := deploy.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for moac deploy
	// 1. openebs-upgrade.dao.mayadata.io/component-group: mayastor
	// 2. openebs-upgrade.dao.mayadata.io/component-name: moac
	desiredLabels[types.OpenEBSComponentGroupLabelKey] = types.OpenEBSMayastorComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.MoacDeploymentNameKey
	// set the desired labels
	deploy.SetLabels(desiredLabels)

	// Overwrite the namespace to mayastor for mayastor based components.
	// Note: mayastor based components will be installed only in mayastor namespace only.
	deploy.SetNamespace(types.MayastorNamespaceNameKey)

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
		if containerName == types.MoacContainerKey {
			// Set the image of the container.
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Image,
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

// updateMoacService updates the moac service manifest as per the
// reconcile.ObservedOpenEBS values.
func (p *Planner) updateMoacService(svc *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := svc.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for moac service
	// 1. openebs-upgrade.dao.mayadata.io/component-group: mayastor
	// 2. openebs-upgrade.dao.mayadata.io/component-name: moac
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSMayastorComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.MoacServiceNameKey
	// set the desired labels
	svc.SetLabels(desiredLabels)

	// Overwrite the namespace to mayastor for mayastor based components.
	// Note: mayastor based components will be installed only in mayastor namespace only.
	svc.SetNamespace(types.MayastorNamespaceNameKey)

	return nil
}

// updateMayastor updates the values of mayastor daemonset as per given configuration.
func (p *Planner) updateMayastor(daemonset *unstructured.Unstructured) error {

	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := daemonset.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for mayastor daemonset:
	// 1. openebs-upgrade.dao.mayadata.io/component-group: mayastor
	// 2. openebs-upgrade.dao.mayadata.io/component-name: mayastor
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSMayastorComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.MayastorDaemonsetNameKey
	// set the desired labels
	daemonset.SetLabels(desiredLabels)

	// Overwrite the namespace to mayastor for mayastor based components.
	// Note: mayastor based components will be installed only in mayastor namespace only.
	daemonset.SetNamespace(types.MayastorNamespaceNameKey)

	containers, err := unstruct.GetNestedSliceOrError(daemonset, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}
	updateContainer := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}

		if containerName == types.MayastorContainerKey {
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Mayastor.Image,
				"spec", "image")
			if err != nil {
				return err
			}

			// Set the resource of the container.
			resources := map[string]interface{}{
				"limits": map[string]interface{}{
					"cpu":           "1",
					"memory":        "500Mi",
					"hugepages-2Mi": "1Gi",
				},
				"requests": map[string]interface{}{
					"cpu":           "1",
					"memory":        "500Mi",
					"hugepages-2Mi": "1Gi",
				},
			}
			err = unstructured.SetNestedField(obj.Object, resources, "spec", "resources")
		}
		if containerName == types.MayastorGRPCContainerKey {
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.MayastorGRPC.Image,
				"spec", "image")
		}
		if err != nil {
			return err
		}

		if p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Resources != nil &&
			containerName != types.MayastorContainerKey {
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Resources,
				"spec", "resources")
		} else if p.ObservedOpenEBS.Spec.Resources != nil &&
			containerName != types.MayastorContainerKey {
			err = unstructured.SetNestedField(obj.Object,
				p.ObservedOpenEBS.Spec.Resources, "spec", "resources")
		}
		if err != nil {
			return err
		}
		return nil
	}
	err = unstruct.SliceIterator(containers).ForEachUpdate(updateContainer)
	if err != nil {
		return err
	}
	err = unstructured.SetNestedSlice(daemonset.Object, containers, "spec",
		"template", "spec", "containers")
	if err != nil {
		return err
	}

	return nil
}

// updateMayastorNamespace updates the mayastor namespace
// structure as per the provided values otherwise default values.
func (p *Planner) updateMayastorNamespace(namespace *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := namespace.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Set some component specific labels in order to identify specific components.
	// These labels will be only set by openebs-upgrade and will help the end-users
	// identify a particular or a set of OpenEBS components.
	//
	// Component specific labels for mayastor daemonset:
	// 1. openebs-upgrade.dao.mayadata.io/component-group: mayastor
	// 2. openebs-upgrade.dao.mayadata.io/component-name: mayastor
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSMayastorComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.MayastorNamespaceNameKey
	// set the desired labels
	namespace.SetLabels(desiredLabels)

	namespace.SetName(types.MayastorNamespaceNameKey)

	return nil
}
