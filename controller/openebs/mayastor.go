package openebs

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/types"
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
		*p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Enabled = true
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
		*p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Enabled = true
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

	return nil
}
