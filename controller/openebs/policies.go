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
