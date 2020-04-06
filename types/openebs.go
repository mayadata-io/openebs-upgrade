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
	// OpenEBS Version to be installed or updated to.
	Version string `json:"version"`

	// DefaultStoragePath is the directory which will be used by
	// default for various OpenEBS operations i.e.,it can be used to
	// specify the hostpath to be used for default Jiva StoragePool
	// loaded by OpenEBS.
	//
	// Defaults to /var/openebs
	DefaultStoragePath string `json:"defaultStoragePath"`

	// If createDefaultStorageConfig is false then OpenEBS default
	// storage class and storage pool will not be created.
	//
	// Defaults to true
	CreateDefaultStorageConfig *bool `json:"createDefaultStorageConfig"`

	// A custom registry could be specified for pulling the container
	// images.
	// Note: This field should be used when user has pulled and pushed
	// the images to a custom registry.
	ImagePrefix string `json:"imagePrefix"`

	// Defaults to IfNotPresent
	// Note: This policy will be applicable to all the images being used
	// for OpenEBS components.
	ImagePullPolicy string `json:"imagePullPolicy"`

	// Resources can be used to specify the resource requests of the containers
	// of the OpenEBS components in terms of CPU and Memory
	//
	// resources provided at this level i.e., .spec.resources will be applicable
	// to all the containers of all the components.
	//
	// This can be overrided by providing it for a particular component in the
	// component's specified section, for example, inside apiServer.
	Resources map[string]interface{} `json:"resources"`

	// All the OpenEBS components that will get installed/updated.
	Components `json:",inline"`
}

// Components stores all the OpenEBS components.
type Components struct {
	APIServer        *APIServer        `json:"apiServer"`
	Provisioner      *Provisioner      `json:"provisioner"`
	LocalProvisioner *LocalProvisioner `json:"localProvisioner"`
	SnapshotOperator *SnapshotOperator `json:"snapshotOperator"`
	AdmissionServer  *AdmissionServer  `json:"admissionServer"`
	NDMDaemon        *NDMDaemon        `json:"ndmDaemon"`
	NDMOperator      *NDMOperator      `json:"ndmOperator"`
	JivaConfig       *JivaConfig       `json:"jivaConfig"`
	CstorConfig      *CstorConfig      `json:"cstorConfig"`
	Helper           *Helper           `json:"helper"`
	Policies         *Policies         `json:"policies"`
	Analytics        *Analytics        `json:"analytics"`
}

// Helper consists of alpine based linux utils docker image used for
// launching helper jobs.
type Helper struct {
	Container `json:",inline"`
}

// Policies consists of the various policies supported by OpenEBS such as
// monitoring.
//
// It stores the config such as which all policies are enabled and what are the
// image tags that should be used for deploying the containers in the k8s cluster.
//
// Currently, only monitoring policy is supported which is deployed as m-exporter
// container.
type Policies struct {
	Monitoring *Monitoring `json:"monitoring"`
}

// Monitoring contains the details of monitoring container
type Monitoring struct {
	Component `json:",inline"`
	Container `json:",inline"`
}

// Analytics is used for enabling/disabling google analytics. If set to true, it
// sends anonymous usage events to Google Analytics
//
// It is set to true by default.
type Analytics struct {
	Enabled      *bool  `json:"enabled"`
	PingInterval string `json:"pingInterval"`
}

// Component stores the configuration of a particular
// component such as it it is enabled or not, no of
// replicas, nodeselector, etc.
type Component struct {
	Enabled      *bool                  `json:"enabled"`
	Replicas     *int32                 `json:"replicas"`
	Resources    map[string]interface{} `json:"resources"`
	NodeSelector map[string]string      `json:"nodeSelector"`
	Tolerations  []interface{}          `json:"tolerations"`
	Affinity     map[string]interface{} `json:"affinity"`
}

// APIServer store the configuration for maya-apiserver
//
// Maya-apiserver helps with the creation of CAS Volumes and provides
// API endpoints to manage those volumes. It can also be considered as
// a template engine that can be easily extended to support any kind of
// CAS storage solutions.
//
// It is deployed as a deployment in the k8s cluster.
type APIServer struct {
	Component       `json:",inline"`
	Container       `json:",inline"`
	CstorSparsePool *CstorSparsePool `json:"cstorSparsePool"`
}

// CstorSparsePool stores the configuration for sparse pools i.e. whether sparse
// pools should be installed by default or not
type CstorSparsePool struct {
	Enabled *bool `json:"enabled"`
}

// Provisioner stores the configuration for OpenEBS provisioner
//
// Provisioner is an implementation of Kubernetes Dynamic Provisioner
// that processes the PVC requests by interacting with maya-apiserver.
//
// It is deployed as a deployment in the k8s cluster.
type Provisioner struct {
	Component `json:",inline"`
	Container `json:",inline"`
}

// LocalProvisioner stores the configuration for OpenEBS local
// provisioner.
//
// LocalProvisioner is responsible for processing the PVC requests
// for provisioning local persistent volumes
//
// It is deployed as a deployment in the k8s cluster.
type LocalProvisioner struct {
	Component `json:",inline"`
	Container `json:",inline"`
}

// AdmissionServer is an implementation of kubernetes validation admission webhook.
//
// It is used for validating various operations before proceeding with them like
// PVC delete operation, etc.
//
// It is deployed as a deployment in k8s cluster.
type AdmissionServer struct {
	Component `json:",inline"`
	Container `json:",inline"`
}

// SnapshotOperator stores the configuration for snapshot operator.
//
// Operator for the snapshot controller and provisioner.
// It consists of the snapshot controller and provisioner containers.
//
// It is deployed as a deployment in the k8s cluster.
type SnapshotOperator struct {
	Component   `json:",inline"`
	Controller  Container `json:"controller"`
	Provisioner Container `json:"provisioner"`
}

// NDMDaemon stores the configuration for node-disk-manager daemonset.
//
// It is a daemonset that helps to manage the disks attached to
// the Kubernetes Nodes.
// It can be used to extend the capabilities of Kubernetes to provide
// access to disk inventory across cluster.
//
// It is deployed as a daemonset in the k8s cluster.
type NDMDaemon struct {
	Component `json:",inline"`
	Container `json:",inline"`
	Sparse    *Sparse     `json:"sparse"`
	Filters   *NDMFilters `json:"filters"`
	Probes    *NDMProbes  `json:"probes"`
}

// Sparse stores the configuration for sparse files.
//
// Sparse File help simulate disk objects that can be used for testing and
// proto typing solutions built using node-disk-manager(NDM).Sparse files
// will be created if NDM is provided with the location where sparse files
// should be located.
type Sparse struct {
	// Path defines a sparse directory for creating a sparse file
	// at the specified directory and an associated BlockDevice CR
	// gets added to Kubernetes.
	Path string `json:"path"`
	// Size define the size of created sparse file
	Size string `json:"size"`
	// Count defines the number of sparse files to be created
	Count string `json:"count"`
}

// NDMFilters stores the configuration for filters being used by NDM
// i.e., filters contain the config for excluding or including vendors,
// paths, etc.
type NDMFilters struct {
	OSDisk *FilterConfigs `json:"osDisk"`
	Vendor *FilterConfigs `json:"vendor"`
	Path   *FilterConfigs `json:"path"`
}

// FilterConfigs contains the config for NDM filters
type FilterConfigs struct {
	Enabled *bool   `json:"enabled"`
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
	Enabled *bool `json:"enabled"`
}

// NDMOperator stores the configuration for ndm operator
//
// NDMOperator is responsible for installation, upgrade and
// lifecycle-management of node-disk-manager.
//
// It is deployed as a deployment in the k8s cluster.
type NDMOperator struct {
	Component `json:",inline"`
	Container `json:",inline"`
}

// JivaConfig stores the configuration for Jiva: CAS Data Engine.
// Jiva provides highly available iSCSI block storage Persistent Volumes
// for Kubernetes Stateful Applications, by making use of the host filesystem.
//
// It consists of a target (or a Storage Controller) that exposes iSCSI,
// while synchronously replicating the data to one or more Replicas and a set of
// replicas that a Target uses to read/write data.
type JivaConfig struct {
	Component `json:",inline"`
	Container `json:",inline"`
}

// CstorConfig stores the configuration for Cstor: CAS Data Engine.
// The primary function of cStor is to serve the iSCSI block storage using
// the underlying disks in a cloud native way.cStor is a very light weight
// and feature rich storage engine. It provides enterprise grade features
// such as synchronous data replication, snapshots, clones, thin provisioning
// of data, high resiliency of data, etc.
//
// It has two main components: cStor pool pods and cStor target pods.
//
// pool, poolMgmt, target and volumeMgmt are the containers which are deployed
// in the k8s cluster.
type CstorConfig struct {
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
	// Phase is the current state of OpenEBS, it can be either Online
	// or Error.
	Phase OpenEBSStatusPhase `json:"phase"`

	// Reason is a brief CamelCase string that describes any failure and is meant
	// for machine parsing and tidy display in the CLI.
	Reason string `json:"reason,omitempty"`

	// Conditions are various states that OpenEBS
	// is currently passing through.
	Conditions []OpenEBSStatusCondition `json:"conditions"`
}

// OpenEBSStatusPhase reports the current phase of OpenEBS
type OpenEBSStatusPhase string

const (
	// OpenEBSStatusPhaseFailed indicates error in
	// OpenEBS
	OpenEBSStatusPhaseFailed OpenEBSStatusPhase = "Failed"

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
