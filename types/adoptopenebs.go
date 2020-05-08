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

package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AdoptOpenEBS defines the intent to adopt existing OpenEBS
// configuration deployed on a kubernetes setup.
type AdoptOpenEBS struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Status            AdoptOpenEBSStatus `json:"status"`
}

// AdoptOpenEBSStatus defines the current status of
// adoptOpenEBS
type AdoptOpenEBSStatus struct {
	// Phase is the current state of adoptOpenEBS, it can be either Online
	// or Failed.
	Phase AdoptOpenEBSStatusPhase `json:"phase"`

	// Reason is a brief CamelCase string that describes any failure and is meant
	// for machine parsing and tidy display in the CLI.
	Reason string `json:"reason,omitempty"`
}

// AdoptOpenEBSStatusPhase reports the current phase of AdoptOpenEBS
type AdoptOpenEBSStatusPhase string

const (
	// AdoptOpenEBSStatusPhaseFailed indicates error in
	// AdoptOpenEBS
	AdoptOpenEBSStatusPhaseFailed AdoptOpenEBSStatusPhase = "Failed"

	// AdoptOpenEBSStatusPhaseOnline indicates AdoptOpenEBS
	// in Online state i.e. no error or warning.
	AdoptOpenEBSStatusPhaseOnline AdoptOpenEBSStatusPhase = "Online"
)
