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

const (
	// NDMVersion045 is the NDM version 0.4.5
	NDMVersion045 string = "v0.4.5"
	// NDMVersion046 is the NDM version 0.4.6
	NDMVersion046 string = "v0.4.6"
	// NDMVersion047 is the NDM version 0.4.7
	NDMVersion047 string = "v0.4.7"
	// NDMVersion048 is the NDM version 0.4.8
	NDMVersion048 string = "v0.4.8"
	// DefaultNDMSparseSize is the default size for NDM Sparse
	DefaultNDMSparseSize string = "10737418240"
	// DefaultNDMSparseCount is the default count for NDM sparse
	DefaultNDMSparseCount string = "0"
	// UdevProbeKey is the key used to identify udev probe in NDM
	UdevProbeKey string = "udev-probe"
	// SmartProbeKey is the key used to identify smart probe in NDM
	SmartProbeKey string = "smart-probe"
	// SeachestProbeKey is the key used to identify seachest probe in NDM
	SeachestProbeKey string = "seachest-probe"
	// VendorFilterKey is the key used to identify vendor filter in NDM
	VendorFilterKey string = "vendor-filter"
	// PathFilterKey is the key used to identify path filter in NDM
	PathFilterKey string = "path-filter"
	// OSDiskFilterKey is the key used to identify OS disk filter in NDM
	OSDiskFilterKey string = "os-disk-exclude-filter"
	// SparseFileSizeEnv is the sparse file size env key
	SparseFileSizeEnv string = "SPARSE_FILE_SIZE"
	// SparseFileCountEnv is the sparse file count env  key
	SparseFileCountEnv string = "SPARSE_FILE_COUNT"
	// SparseFileDirectoryEnv is the sparse directory env key
	SparseFileDirectoryEnv string = "SPARSE_FILE_DIR"
	// DefaultNDMOperatorReplicaCount is the default replica count for NDM operator
	DefaultNDMOperatorReplicaCount int32 = 1
)

// NDMConfig stores the configuration for node-disk-manager configmap
type NDMConfig struct {
	ProbeConfigs  []ProbeConfig  `json:"probeconfigs"`
	FilterConfigs []FilterConfig `json:"filterconfigs"`
}

// ProbeConfig contains the configuration related to NDM probes
type ProbeConfig struct {
	Key   string `json:"key"`
	Name  string `json:"name"`
	State string `json:"state"`
}

// FilterConfig contains the configuration related to NDM filters
type FilterConfig struct {
	Key     string  `json:"key"`
	Name    string  `json:"name"`
	State   string  `json:"state"`
	Include *string `json:"include,omitempty"`
	Exclude *string `json:"exclude,omitempty"`
}
