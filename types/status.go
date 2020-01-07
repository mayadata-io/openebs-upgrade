/*
Copyright 2019 The MayaData Authors.
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

// ConditionType is a custom datatype that refers to
// various conditions supported by this operator.
type ConditionType string

const (
	// OpenEBSReconcileErrorCondition is used to indicate
	// presence or absence of error while reconciling
	// OpenEBS.
	OpenEBSReconcileErrorCondition ConditionType = "OpenEBSReconcileError"
)

// ConditionState is a custom datatype that
// refers to presence or absence of any condition
type ConditionState string

const (
	// ConditionIsPresent refers to presence of any condition
	ConditionIsPresent ConditionState = "True"

	// ConditionIsAbsent refers to absence of any condition
	ConditionIsAbsent ConditionState = "False"
)

// StatusPhase refers to various phases found in a resource
// status
type StatusPhase string

const (
	// StatusPhaseError refers to a generic error status phase
	StatusPhaseError StatusPhase = "Error"
)

// now returns the current time in following format
// 2006-01-02 15:04:05.000000
func now() string {
	return metav1.Now().Format("2006-01-02 15:04:05.000000")
}

// MakeOpenEBSReconcileErrCond builds a new
// OpenEBSReconcileError condition
// suitable to be used in API status.conditions
func MakeOpenEBSReconcileErrCond(err error) map[string]interface{} {
	return map[string]interface{}{
		"type":             OpenEBSReconcileErrorCondition,
		"status":           ConditionIsPresent,
		"reason":           err.Error(),
		"lastObservedTime": now(),
	}
}

// MergeNoReconcileErrorOnOpenEBS sets
// OpenEBSConditionReconcileError condition to false.
func MergeNoReconcileErrorOnOpenEBS(obj *OpenEBS) {
	noErrCond := OpenEBSStatusCondition{
		Type:             OpenEBSReconcileErrorCondition,
		Status:           ConditionIsAbsent,
		LastObservedTime: now(),
	}
	var newConds []OpenEBSStatusCondition
	for _, old := range obj.Status.Conditions {
		if old.Type == OpenEBSReconcileErrorCondition {
			// ignore previous occurrence of ReconcileError
			continue
		}
		newConds = append(newConds, old)
	}
	newConds = append(newConds, noErrCond)
	obj.Status.Conditions = newConds
}
