package openebs

import (
	"github.com/mayadata-io/openebs-operator/types"
)

// setPoliciesDefaultsIfNotSet sets the default values for various policies
// being used in OpenEBS such as monitoring, etc.
func (r *Reconciler) setPoliciesDefaultsIfNotSet() error {
	if r.OpenEBS.Spec.Policies == nil {
		r.OpenEBS.Spec.Policies = &types.Policies{
			Monitoring: &types.Monitoring{},
		}
	}
	if r.OpenEBS.Spec.Policies.Monitoring.Enabled == "" {
		r.OpenEBS.Spec.Policies.Monitoring.Enabled = types.True
	}
	// form the monitoring image which is used by cstor pool exporter
	// and volume monitor containers.
	if r.OpenEBS.Spec.Policies.Monitoring.ImageTag == "" {
		r.OpenEBS.Spec.Policies.Monitoring.ImageTag = r.OpenEBS.Spec.Version
	}
	r.OpenEBS.Spec.Policies.Monitoring.Image = r.OpenEBS.Spec.ImagePrefix +
		"m-exporter:" + r.OpenEBS.Spec.Policies.Monitoring.ImageTag
	return nil
}
