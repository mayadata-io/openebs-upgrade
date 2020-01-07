package openebs

import (
	"github.com/mayadata-io/openebs-operator/types"
)

// Set the analytics default values if not already set
func (r *Reconciler) setAnalyticsDefaultsIfNotSet() error {
	if r.OpenEBS.Spec.Analytics == nil {
		r.OpenEBS.Spec.Analytics = &types.Analytics{}
	}
	if r.OpenEBS.Spec.Analytics.Enabled == "" {
		r.OpenEBS.Spec.Analytics.Enabled = types.True
	}
	return nil
}
