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
	// JivaVersion162 is the supported image tag for Jiva components
	// for OpenEBS version 1.6.0.
	JivaVersion162 string = "1.6.2"
	// JivaVersion171 is the supported image tag for Jiva components
	// for OpenEBS version 1.7.0.
	JivaVersion171 string = "1.7.1"
)

// supportedJivaVersionForOpenEBSVersion stores the mapping for
// Jiva to OpenEBS version i.e., a Jiva version for each of the
// supported OpenEBS versions.
// Note: This will be referred to form the container images in
// order to install/update jiva controller/replica components for a
// particular OpenEBS version.
//
// NOTE: There will be an entry for only those OpenEBS versions and the
// supported Jiva versions where the image tag for Jiva images are different
// from the OpenEBS version such as in case of OpenEBS version 1.6.0, the jiva
// controller/replica image that should be used is 1.6.2 instead of 1.6.0 since
// 1.6.2 have some critical fixes which 1.6.0 doesn't have.
var supportedJivaVersionForOpenEBSVersion = map[string]string{
	types.OpenEBSVersion160: JivaVersion162,
	types.OpenEBSVersion170: JivaVersion171,
}

// Set the default values for JIVA.
func (p *Planner) setJIVADefaultsIfNotSet() error {
	if p.ObservedOpenEBS.Spec.JivaConfig == nil {
		p.ObservedOpenEBS.Spec.JivaConfig = &types.JivaConfig{}
	}
	// form the jiva image being used by jiva-controller and
	// replica.
	if p.ObservedOpenEBS.Spec.JivaConfig.ImageTag == "" {
		if jivaVersion, exist := supportedJivaVersionForOpenEBSVersion[p.ObservedOpenEBS.Spec.Version]; exist {
			p.ObservedOpenEBS.Spec.JivaConfig.ImageTag = jivaVersion
		} else {
			p.ObservedOpenEBS.Spec.JivaConfig.ImageTag = p.ObservedOpenEBS.Spec.Version
		}
	}
	p.ObservedOpenEBS.Spec.JivaConfig.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"jiva:" + p.ObservedOpenEBS.Spec.JivaConfig.ImageTag

	// Set the default replica count for Jiva which is 3.
	if p.ObservedOpenEBS.Spec.JivaConfig.Replicas == nil {
		p.ObservedOpenEBS.Spec.JivaConfig.Replicas = new(int32)
		*p.ObservedOpenEBS.Spec.JivaConfig.Replicas = DefaultJivaReplicaCount
	}
	return nil
}
