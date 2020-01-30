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

// Set the default values for Cstor if not already given.
func (r *Reconciler) setCStorDefaultsIfNotSet() error {
	if r.OpenEBS.Spec.CstorConfig == nil {
		r.OpenEBS.Spec.CstorConfig = &types.CstorConfig{}
	}
	// form the cstor-pool image
	if r.OpenEBS.Spec.CstorConfig.Pool.ImageTag == "" {
		r.OpenEBS.Spec.CstorConfig.Pool.ImageTag = r.OpenEBS.Spec.Version
	}
	r.OpenEBS.Spec.CstorConfig.Pool.Image = r.OpenEBS.Spec.ImagePrefix +
		"cstor-pool:" + r.OpenEBS.Spec.CstorConfig.Pool.ImageTag

	// form the cstor-pool-mgmt image
	if r.OpenEBS.Spec.CstorConfig.PoolMgmt.ImageTag == "" {
		r.OpenEBS.Spec.CstorConfig.PoolMgmt.ImageTag = r.OpenEBS.Spec.Version
	}
	r.OpenEBS.Spec.CstorConfig.PoolMgmt.Image = r.OpenEBS.Spec.ImagePrefix +
		"cstor-pool-mgmt:" + r.OpenEBS.Spec.CstorConfig.PoolMgmt.ImageTag

	// form the cstor-istgt image
	if r.OpenEBS.Spec.CstorConfig.Target.ImageTag == "" {
		r.OpenEBS.Spec.CstorConfig.Target.ImageTag = r.OpenEBS.Spec.Version
	}
	r.OpenEBS.Spec.CstorConfig.Target.Image = r.OpenEBS.Spec.ImagePrefix +
		"cstor-istgt:" + r.OpenEBS.Spec.CstorConfig.Target.ImageTag

	// form the cstor-volume-mgmt image
	if r.OpenEBS.Spec.CstorConfig.VolumeMgmt.ImageTag == "" {
		r.OpenEBS.Spec.CstorConfig.VolumeMgmt.ImageTag = r.OpenEBS.Spec.Version
	}
	r.OpenEBS.Spec.CstorConfig.VolumeMgmt.Image = r.OpenEBS.Spec.ImagePrefix +
		"cstor-volume-mgmt:" + r.OpenEBS.Spec.CstorConfig.VolumeMgmt.ImageTag

	return nil
}
