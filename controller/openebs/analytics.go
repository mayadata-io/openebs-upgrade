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
func (r *Reconciler) setAnalyticsDefaultsIfNotSet() error {
	if r.OpenEBS.Spec.Analytics == nil {
		r.OpenEBS.Spec.Analytics = &types.Analytics{}
	}
	if r.OpenEBS.Spec.Analytics.Enabled == nil {
		r.OpenEBS.Spec.Analytics.Enabled = new(bool)
		*r.OpenEBS.Spec.Analytics.Enabled = true
	}
	return nil
}
