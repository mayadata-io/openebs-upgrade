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

// setPoliciesDefaultsIfNotSet sets the default values for various policies
// being used in OpenEBS such as monitoring, etc.
func (p *Planner) setPoliciesDefaultsIfNotSet() error {
	if p.ObservedOpenEBS.Spec.Policies == nil {
		p.ObservedOpenEBS.Spec.Policies = &types.Policies{
			Monitoring: &types.Monitoring{},
		}
	}
	if p.ObservedOpenEBS.Spec.Policies.Monitoring.Enabled == nil {
		p.ObservedOpenEBS.Spec.Policies.Monitoring.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.Policies.Monitoring.Enabled = true
	}
	// form the monitoring image which is used by cstor pool exporter
	// and volume monitor containers.
	if p.ObservedOpenEBS.Spec.Policies.Monitoring.ImageTag == "" {
		p.ObservedOpenEBS.Spec.Policies.Monitoring.ImageTag = p.ObservedOpenEBS.Spec.Version +
			p.ObservedOpenEBS.Spec.ImageTagSuffix
	}
	p.ObservedOpenEBS.Spec.Policies.Monitoring.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"m-exporter:" + p.ObservedOpenEBS.Spec.Policies.Monitoring.ImageTag
	return nil
}
