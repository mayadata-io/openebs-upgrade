package openebs

import (
	"github.com/mayadata-io/openebs-operator/types"
)

// Set the default values for Cstor if not already given.
func (r *Reconciler) setCStorDefaultsIfNotSet() error {
	if r.OpenEBS.Spec.Cstor == nil {
		r.OpenEBS.Spec.Cstor = &types.Cstor{}
	}
	// form the cstor-pool image
	if r.OpenEBS.Spec.Cstor.Pool.ImageTag == "" {
		r.OpenEBS.Spec.Cstor.Pool.ImageTag = r.OpenEBS.Spec.Version
	}
	r.OpenEBS.Spec.Cstor.Pool.Image = r.OpenEBS.Spec.ImagePrefix +
		"cstor-pool:" + r.OpenEBS.Spec.Cstor.Pool.ImageTag

	// form the cstor-pool-mgmt image
	if r.OpenEBS.Spec.Cstor.PoolMgmt.ImageTag == "" {
		r.OpenEBS.Spec.Cstor.PoolMgmt.ImageTag = r.OpenEBS.Spec.Version
	}
	r.OpenEBS.Spec.Cstor.PoolMgmt.Image = r.OpenEBS.Spec.ImagePrefix +
		"cstor-pool-mgmt:" + r.OpenEBS.Spec.Cstor.PoolMgmt.ImageTag

	// form the cstor-istgt image
	if r.OpenEBS.Spec.Cstor.Target.ImageTag == "" {
		r.OpenEBS.Spec.Cstor.Target.ImageTag = r.OpenEBS.Spec.Version
	}
	r.OpenEBS.Spec.Cstor.Target.Image = r.OpenEBS.Spec.ImagePrefix +
		"cstor-istgt:" + r.OpenEBS.Spec.Cstor.Target.ImageTag

	// form the cstor-volume-mgmt image
	if r.OpenEBS.Spec.Cstor.VolumeMgmt.ImageTag == "" {
		r.OpenEBS.Spec.Cstor.VolumeMgmt.ImageTag = r.OpenEBS.Spec.Version
	}
	r.OpenEBS.Spec.Cstor.VolumeMgmt.Image = r.OpenEBS.Spec.ImagePrefix +
		"cstor-volume-mgmt:" + r.OpenEBS.Spec.Cstor.VolumeMgmt.ImageTag

	return nil
}
