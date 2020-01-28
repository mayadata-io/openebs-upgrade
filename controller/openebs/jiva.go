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
	"mayadata.io/openebs-upgrade/types"
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
