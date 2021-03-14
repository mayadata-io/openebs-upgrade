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
	// NDMVersion049 is the NDM version 0.4.9
	NDMVersion049 string = "v0.4.9"
	// NDMVersion049EE is the enterprise NDM version v0.4.9-ee
	NDMVersion049EE string = "v0.4.9-ee"
	// NDMVersion050 is the NDM version 0.5.0
	NDMVersion050 string = "0.5.0"
	// NDMVersion050EE is the NDM version 0.5.0-ee
	NDMVersion050EE string = "0.5.0-ee"
	// NDMVersion060 is the NDM version 0.6.0
	NDMVersion060 string = "0.6.0"
	// NDMVersion060EE is the NDM version 0.6.0-ee
	NDMVersion060EE string = "0.6.0-ee"
	// NDMVersion070 is the NDM version 0.7.0
	NDMVersion070 string = "0.7.0"
	// NDMVersion070EE is the NDM version 0.7.0-ee
	NDMVersion070EE string = "0.7.0-ee"
	// NDMVersion080 is the NDM version 0.8.0
	NDMVersion080 string = "0.8.0"
	// NDMVersion080EE is the NDM version 0.8.0-ee
	NDMVersion080EE string = "0.8.0-ee"
	// NDMVersion082 is the NDM version 0.8.2
	NDMVersion082 string = "0.8.2"
	// NDMVersion082EE is the NDM version 0.8.2-ee
	NDMVersion082EE string = "0.8.2-ee"
	// NDMVersion091 is the NDM version 0.9.1
	NDMVersion091 string = "0.9.1"
	// NDMVersion091EE is the NDM version 0.9.1-ee
	NDMVersion091EE string = "0.9.1-ee"
	// NDMVersion101 is the NDM version 1.0.1
	NDMVersion101 string = "1.0.1"
	// NDMVersion110 is the NDM version 1.1.0
	NDMVersion110 string = "1.1.0"
	// NDMVersion120 is the NDM version 1.2.0
	NDMVersion120 string = "1.2.0"
	// NDMVersion130 is the NDM version 1.3.0
	NDMVersion130 string = "1.3.0"
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
	// CleanupJobImageEnv is the cleanup job image env key
	CleanupJobImageEnv string = "CLEANUP_JOB_IMAGE"
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
