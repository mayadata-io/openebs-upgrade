/*
Copyright 2019 The MayaData Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=openebsoperator

// OpenEBSOperator defines the intent to get
// OpenEBS deployed on a Kubernetes setup
type OpenEBSOperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec OpenEBSOperatorSpec `json:"spec"`

	Status OpenEBSOperatorStatus `json:"status"`
}

// OpenEBSOperatorSpec defines the specifications
// that determines what OpenEBS (e.g. version, components, etc)
// get deployed on a Kubernetes setup
type OpenEBSOperatorSpec struct {
	Version string `json:"version"`
}

// OpenEBSOperatorStatus defines the current status of
// openebs operator
type OpenEBSOperatorStatus struct {
	// Phase is the current state of openebs operator
	Phase OpenEBSOperatorPhase `json:"phase"`

	// Conditions are various states that this openebs operator
	// is currently passing through.
	//
	// NOTE:
	//  There can be cases when this operator can be identified by
	// multiple conditions
	Conditions []OpenEBSOperatorCondition `json:"conditions"`
}

// OpenEBSOperatorPhase is the resulting state of this operator
// that will provide some insight to external resources / cluster
type OpenEBSOperatorPhase string

const (
	// OpenEBSOperatorStarted indicates that this operator specifications
	// are being worked upon
	OpenEBSOperatorStarted OpenEBSOperatorPhase = "Started"

	// OpenEBSOperatorFailed indicates that this operator specifications
	// met with some failure during reconcile
	OpenEBSOperatorFailed OpenEBSOperatorPhase = "Failed"

	// OpenEBSOperatorCompleted indicates if this operator specifications
	// was completed successfully. In other words, the reconcile attempt
	// was successful in getting the desired state to actual state.
	OpenEBSOperatorCompleted OpenEBSOperatorPhase = "Completed"

	// OpenEBSOperatorIgnored indicates if this operator specifications
	// was ignored due to some reason. This reason might be understood
	// by looking at the operator status' conditions.
	OpenEBSOperatorIgnored OpenEBSOperatorPhase = "Ignored"
)

// OpenEBSOperatorCondition defines the current state of this
// operator
type OpenEBSOperatorCondition struct {
	// Type is a unique identification of operator condition
	Type OpenEBSOperatorConditionType `json:"type"`

	// Message provides a descriptive message associated to this
	// condition. This can be an error, warning or info message
	Message string `json:"message"`

	// LastUpdatedTime is the last time this condition was found
	LastUpdatedTime metav1.Time `json:"lastUpdatedTime"`

	// Frequency is the number of consecutive times this condition
	// was found
	Frequency int64 `json:"frequency"`
}

// OpenEBSOperatorConditionType is a unique identification of
// operator condition
type OpenEBSOperatorConditionType string
