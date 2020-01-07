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
	Version                    string `json:"version"`
	DefaultStoragePath         string `json:"defaultStoragePath"`
	CreateDefaultStorageConfig string `json:"createDefaultStorageConfig"`
	ImagePrefix                string `json:"imagePrefix"`
	ImagePullPolicy            string `json:"imagePullPolicy"`
	Components                 `json:",inline"`
}

// Components stores all the OpenEBS components.
type Components struct {
	APIServer        APIServer        `json:"apiServer"`
	Provisioner      Provisioner      `json:"provisioner"`
	LocalProvisioner LocalProvisioner `json:"localProvisioner"`
	SnapshotOperator SnapshotOperator `json:"snapshotOperator"`
	NDM              NDM              `json:"ndm"`
	NDMOperator      NDMOperator      `json:"ndmOperator"`
	Jiva             Jiva             `json:"jiva"`
	Cstor            Cstor            `json:"cstor"`
}

// ComponentCommonConfig stores the configuration which are common/same
// for most of the components.
type ComponentCommonConfig struct {
	Enabled      string            `json:"enabled"`
	Replicas     int               `json:"replicas"`
	NodeSelector map[string]string `json:"nodeSelector"`
}

// APIServer stores the configuration of maya-apiserver
type APIServer struct {
	ComponentCommonConfig `json:",inline"`
	Container             `json:",inline"`
	Sparse                SparsePools `json:"sparse"`
}

// SparsePools stores the configuration for sparse pools i.e. whether sparse
// pools should be installed by default or not
type SparsePools struct {
	Enabled string `json:"enabled"`
}

// Provisioner stores the configuration of OpenEBS provisioner
type Provisioner struct {
	ComponentCommonConfig `json:",inline"`
	Container             `json:",inline"`
}

// LocalProvisioner stores the configuration of OpenEBS local provisioner
type LocalProvisioner struct {
	ComponentCommonConfig `json:",inline"`
	Container             `json:",inline"`
}

// SnapshotOperator stores the configuration of Snapshot Operator
type SnapshotOperator struct {
	ComponentCommonConfig `json:",inline"`
	Controller            Container `json:"controller"`
	Provisioner           Container `json:"provisioner"`
}

// NDM stores the configuration of Node disk Manager
type NDM struct {
	ComponentCommonConfig `json:",inline"`
	Container             `json:",inline"`
	Sparse                Sparse     `json:"sparse"`
	Filters               NDMFilters `json:"filters"`
	Probes                NDMProbes  `json:"probes"`
}

// Sparse file configuration for NDM.
type Sparse struct {
	Path  string `json:"path"`
	Size  string `json:"size"`
	Count string `json:"count"`
}

// NDMFilters stores the configuration for filters being used by NDM
// i.e., filters contain the config for excluding or including vendors,
// paths, etc.
type NDMFilters struct {
	Exclude Exclude `json:"exclude"`
	Include Include `json:"include"`
}

// Exclude can be used to exclude something from a group i.e., it
// could be vendor, path, etc in case of NDM.
type Exclude struct {
	Vendors string `json:"vendors"`
	Paths   string `json:"paths"`
}

// Include can be used to include something to a group or set of things.
type Include struct {
	Paths string `json:"paths"`
}

// NDMProbes can be used to configure NDM probes i.e., it can be used to
// enable/disable various probes used by NDM such as seachest, smart,
// capacity, etc.
type NDMProbes struct {
	Enable []string `json:"enable"`
}

// NDMOperator stores the configuration of NDM operator
type NDMOperator struct {
	ComponentCommonConfig `json:",inline"`
	Container             `json:",inline"`
}

// Jiva stores the configuration of jiva.
type Jiva struct {
	ComponentCommonConfig `json:",inline"`
	Container             `json:",inline"`
}

// Cstor stores the configuration of cstor.
type Cstor struct {
	Pool       Container `json:"pool"`
	PoolMgmt   Container `json:"poolMgmt"`
	Target     Container `json:"target"`
	VolumeMgmt Container `json:"volumeMgmt"`
}

// Container stores the details of a container
type Container struct {
	ImageTag string `json:"imageTag"`
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
