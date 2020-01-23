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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OpenEBS defines the intent to get
// OpenEBS deployed/updated on a Kubernetes setup
type OpenEBS struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec OpenEBSSpec `json:"spec"`

	Status OpenEBSStatus `json:"status"`
}

// OpenEBSSpec defines the specifications that determines
// what OpenEBS (e.g. version, components, etc)
// get deployed on a Kubernetes setup
type OpenEBSSpec struct {
	Version                    string            `json:"version"`
	DefaultStoragePath         string            `json:"defaultStoragePath"`
	CreateDefaultStorageConfig string            `json:"createDefaultStorageConfig"`
	ImagePrefix                string            `json:"imagePrefix"`
	ImagePullPolicy            corev1.PullPolicy `json:"imagePullPolicy"`
	Components                 `json:",inline"`
}

// Components stores all the OpenEBS components.
type Components struct {
	APIServer        *APIServer        `json:"apiServer"`
	Provisioner      *Provisioner      `json:"provisioner"`
	LocalProvisioner *LocalProvisioner `json:"localProvisioner"`
	SnapshotOperator *SnapshotOperator `json:"snapshotOperator"`
	AdmissionServer  *AdmissionServer  `json:"admissionServer"`
	NDM              *NDM              `json:"ndm"`
	NDMOperator      *NDMOperator      `json:"ndmOperator"`
	Jiva             *Jiva             `json:"jiva"`
	Cstor            *Cstor            `json:"cstor"`
	Helper           *Helper           `json:"helper"`
	Policies         *Policies         `json:"policies"`
	Analytics        *Analytics        `json:"analytics"`
}

// Helper stores the details of helper containers used by OpenEBS
type Helper struct {
	Container `json:",inline"`
}

// Policies stores the details of various policies being supported
// by OpenEBS i.e. Monitoring.
type Policies struct {
	Monitoring *Monitoring `json:"monitoring"`
}

// Monitoring contains the details of monitoring container
type Monitoring struct {
	ComponentCommonConfig `json:",inline"`
	Container             `json:",inline"`
}

// Analytics contains the configuration of analytics being used
// in OpenEBS
type Analytics struct {
	Enabled      string `json:"enabled"`
	PingInterval string `json:"pingInterval"`
}

// ComponentCommonConfig stores the configuration which are common/same
// for most of the components.
type ComponentCommonConfig struct {
	Enabled      string              `json:"enabled"`
	Replicas     *int32              `json:"replicas"`
	NodeSelector map[string]string   `json:"nodeSelector"`
	Tolerations  []corev1.Toleration `json:"tolerations"`
	Affinity     *corev1.Affinity    `json:"affinity"`
}

// APIServer stores the configuration of maya-apiserver
type APIServer struct {
	ComponentCommonConfig `json:",inline"`
	Container             `json:",inline"`
	Sparse                *SparsePools `json:"sparse"`
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

// AdmissionServer stores the configuration of openebs-admission-server
type AdmissionServer struct {
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
	Sparse                *Sparse     `json:"sparse"`
	Filters               *NDMFilters `json:"filters"`
	Probes                *NDMProbes  `json:"probes"`
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
	OSDisk *FilterConfigs `json:"OSDisk"`
	Vendor *FilterConfigs `json:"vendor"`
	Path   *FilterConfigs `json:"path"`
}

// FilterConfigs contains the config for NDM filters
type FilterConfigs struct {
	Enabled string  `json:"enabled"`
	Include *string `json:"include"`
	Exclude *string `json:"exclude"`
}

// NDMProbes can be used to configure NDM probes i.e., it can be used to
// enable/disable various probes used by NDM such as seachest, smart,
// capacity, etc.
type NDMProbes struct {
	Udev     *ProbeState `json:"udev"`
	Smart    *ProbeState `json:"smart"`
	Seachest *ProbeState `json:"seachest"`
}

// ProbeState denotes the current state of a NDM probe
type ProbeState struct {
	Enabled string `json:"enabled"`
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
	Image    string `json:"image"`
}

// OpenEBSStatus defines the current status of
// OpenEBS
type OpenEBSStatus struct {
	// Phase is the current state of OpenEBS
	Phase OpenEBSStatusPhase `json:"phase"`

	// Conditions are various states that OpenEBS
	// is currently passing through.
	Conditions []OpenEBSStatusCondition `json:"conditions"`
}

// OpenEBSStatusPhase reports the current phase of OpenEBS
type OpenEBSStatusPhase string

const (
	// OpenEBSStatusPhaseError indicates error in
	// OpenEBS
	OpenEBSStatusPhaseError OpenEBSStatusPhase = "Error"

	// OpenEBSStatusPhaseOnline indicates
	// OpenEBS in Online state i.e. no error or warning
	OpenEBSStatusPhaseOnline OpenEBSStatusPhase = "Online"
)

// OpenEBSStatusCondition defines the current state of OpenEBS
type OpenEBSStatusCondition struct {
	Type             ConditionType  `json:"type"`
	Status           ConditionState `json:"status"`
	Reason           string         `json:"reason,omitempty"`
	LastObservedTime string         `json:"lastObservedTime"`
}
