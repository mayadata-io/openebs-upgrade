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

// Set the analytics default values if not already set
func (p *Planner) setAnalyticsDefaultsIfNotSet() error {
	if p.ObservedOpenEBS.Spec.Analytics == nil {
		p.ObservedOpenEBS.Spec.Analytics = &types.Analytics{}
	}
	if p.ObservedOpenEBS.Spec.Analytics.Enabled == nil {
		p.ObservedOpenEBS.Spec.Analytics.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.Analytics.Enabled = true
	}
	return nil
}
