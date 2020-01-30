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
	if r.OpenEBS.Spec.JivaConfig == nil {
		r.OpenEBS.Spec.JivaConfig = &types.JivaConfig{}
	}
	// form the jiva image being used by jiva-controller and
	// replica.
	if r.OpenEBS.Spec.JivaConfig.ImageTag == "" {
		r.OpenEBS.Spec.JivaConfig.ImageTag = r.OpenEBS.Spec.Version
	}
	r.OpenEBS.Spec.JivaConfig.Image = r.OpenEBS.Spec.ImagePrefix +
		"jiva:" + r.OpenEBS.Spec.JivaConfig.ImageTag

	// Set the default replica count for Jiva which is 3.
	if r.OpenEBS.Spec.JivaConfig.Replicas == nil {
		r.OpenEBS.Spec.JivaConfig.Replicas = new(int32)
		*r.OpenEBS.Spec.JivaConfig.Replicas = DefaultJivaReplicaCount
	}
	return nil
}
