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
	if r.OpenEBS.Spec.Cstor == nil {
		r.OpenEBS.Spec.Cstor = &types.Cstor{}
	}
	// form the cstor-pool image
	if r.OpenEBS.Spec.Cstor.Pool.ImageTag == "" {
		r.OpenEBS.Spec.Cstor.Pool.ImageTag = r.OpenEBS.Spec.Version
	}
	r.OpenEBS.Spec.Cstor.Pool.Image = r.OpenEBS.Spec.ImagePrefix +
		"cstor-pool:" + r.OpenEBS.Spec.Cstor.Pool.ImageTag

	// form the cstor-pool-mgmt image
	if r.OpenEBS.Spec.Cstor.PoolMgmt.ImageTag == "" {
		r.OpenEBS.Spec.Cstor.PoolMgmt.ImageTag = r.OpenEBS.Spec.Version
	}
	r.OpenEBS.Spec.Cstor.PoolMgmt.Image = r.OpenEBS.Spec.ImagePrefix +
		"cstor-pool-mgmt:" + r.OpenEBS.Spec.Cstor.PoolMgmt.ImageTag

	// form the cstor-istgt image
	if r.OpenEBS.Spec.Cstor.Target.ImageTag == "" {
		r.OpenEBS.Spec.Cstor.Target.ImageTag = r.OpenEBS.Spec.Version
	}
	r.OpenEBS.Spec.Cstor.Target.Image = r.OpenEBS.Spec.ImagePrefix +
		"cstor-istgt:" + r.OpenEBS.Spec.Cstor.Target.ImageTag

	// form the cstor-volume-mgmt image
	if r.OpenEBS.Spec.Cstor.VolumeMgmt.ImageTag == "" {
		r.OpenEBS.Spec.Cstor.VolumeMgmt.ImageTag = r.OpenEBS.Spec.Version
	}
	r.OpenEBS.Spec.Cstor.VolumeMgmt.Image = r.OpenEBS.Spec.ImagePrefix +
		"cstor-volume-mgmt:" + r.OpenEBS.Spec.Cstor.VolumeMgmt.ImageTag

	return nil
}
