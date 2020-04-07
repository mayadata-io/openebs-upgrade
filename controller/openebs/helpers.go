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

// Set the default values for helpers used by OpenEBS components
// such as linux-utils, etc.
func (p *Planner) setHelperDefaultsIfNotSet() error {
	if p.ObservedOpenEBS.Spec.Helper == nil {
		p.ObservedOpenEBS.Spec.Helper = &types.Helper{}
	}
	// form the linux-utils image
	if p.ObservedOpenEBS.Spec.Helper.ImageTag == "" {
		p.ObservedOpenEBS.Spec.Helper.ImageTag = p.ObservedOpenEBS.Spec.Version
	}
	p.ObservedOpenEBS.Spec.Helper.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"linux-utils:" + p.ObservedOpenEBS.Spec.Helper.ImageTag
	return nil
}
