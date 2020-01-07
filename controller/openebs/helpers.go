package openebs

import (
	"github.com/mayadata-io/openebs-operator/types"
)

// Set the default values for helpers used by OpenEBS components
// such as linux-utils, etc.
func (r *Reconciler) setHelperDefaultsIfNotSet() error {
	if r.OpenEBS.Spec.Helper == nil {
		r.OpenEBS.Spec.Helper = &types.Helper{}
	}
	// form the linux-utils image
	if r.OpenEBS.Spec.Helper.ImageTag == "" {
		r.OpenEBS.Spec.Helper.ImageTag = r.OpenEBS.Spec.Version
	}
	r.OpenEBS.Spec.Helper.Image = r.OpenEBS.Spec.ImagePrefix +
		"linux-utils:" + r.OpenEBS.Spec.Helper.ImageTag
	return nil
}
