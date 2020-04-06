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
	// DefaultAdmissionServerReplicaCount is the default value of replica for
	// Admission server.
	DefaultAdmissionServerReplicaCount int32 = 1
)

// Set the admission server default values if not already set
func (p *Planner) setAdmissionServerDefaultsIfNotSet() error {
	// Initialize admissionserver field if not set
	if p.ObservedOpenEBS.Spec.AdmissionServer == nil {
		p.ObservedOpenEBS.Spec.AdmissionServer = &types.AdmissionServer{}
	}
	// Enable admission server if the field is not set i.e. set the
	// value to true.
	// TODO: Validate the values that can be provided for this
	// field.
	if p.ObservedOpenEBS.Spec.AdmissionServer.Enabled == nil {
		p.ObservedOpenEBS.Spec.AdmissionServer.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.AdmissionServer.Enabled = true
	}
	if p.ObservedOpenEBS.Spec.AdmissionServer.ImageTag == "" {
		p.ObservedOpenEBS.Spec.AdmissionServer.ImageTag = p.ObservedOpenEBS.Spec.Version
	}
	p.ObservedOpenEBS.Spec.AdmissionServer.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"admission-server:" + p.ObservedOpenEBS.Spec.AdmissionServer.ImageTag

	if p.ObservedOpenEBS.Spec.AdmissionServer.Replicas == nil {
		p.ObservedOpenEBS.Spec.AdmissionServer.Replicas = new(int32)
		*p.ObservedOpenEBS.Spec.AdmissionServer.Replicas = DefaultAdmissionServerReplicaCount
	}
	return nil
}
