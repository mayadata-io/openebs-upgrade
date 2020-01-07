package openebs

import (
	"github.com/mayadata-io/openebs-operator/types"
)

const (
	// DefaultJivaReplicaCount is the default value of jiva replicas
	DefaultJivaReplicaCount int32 = 3
)

// Set the default values for JIVA.
func (r *Reconciler) setJIVADefaultsIfNotSet() error {
	if r.OpenEBS.Spec.Jiva == nil {
		r.OpenEBS.Spec.Jiva = &types.Jiva{}
	}
	// form the jiva image being used by jiva-controller and
	// replica.
	if r.OpenEBS.Spec.Jiva.ImageTag == "" {
		r.OpenEBS.Spec.Jiva.ImageTag = r.OpenEBS.Spec.Version
	}
	r.OpenEBS.Spec.Jiva.Image = r.OpenEBS.Spec.ImagePrefix +
		"jiva:" + r.OpenEBS.Spec.Jiva.ImageTag

	// Set the default replica count for Jiva which is 3.
	if r.OpenEBS.Spec.Jiva.Replicas == nil {
		r.OpenEBS.Spec.Jiva.Replicas = new(int32)
		*r.OpenEBS.Spec.Jiva.Replicas = DefaultJivaReplicaCount
	}
	return nil
}
