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
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"mayadata.io/openebs-upgrade/types"
)

// supportedNDMVersionForOpenEBSVersion stores the mapping for
// NDM to OpenEBS version i.e., a NDM version for each of the
// supported OpenEBS versions.
// Note: This will be referred to form the container images in
// order to install/update NDM components for a particular OpenEBS
// version.
var supportedNDMVersionForOpenEBSVersion = map[string]string{
	types.OpenEBSVersion150: types.NDMVersion045,
	types.OpenEBSVersion160: types.NDMVersion046,
}

// add/update NDM defaults if not already provided
func (r *Reconciler) setNDMDefaultsIfNotSet() error {
	// Check if NDM field is set or not, if not then
	// initialize it.
	if r.OpenEBS.Spec.NDM == nil {
		r.OpenEBS.Spec.NDM = &types.NDM{}
	}
	// Check if NDM enabling field is set or not,
	// if not then set it to true.
	// Note: By default, NDM will always be enabled i.e.
	// will be installed alongwith other OpenEBS components.
	if r.OpenEBS.Spec.NDM.Enabled == "" {
		r.OpenEBS.Spec.NDM.Enabled = types.True
	}
	// Check if imageTag fiels is set or not, if not
	// then set the NDM image tag as per the OpenEBS version
	// given.
	if r.OpenEBS.Spec.NDM.ImageTag == "" {
		if ndmVersion, exist :=
			supportedNDMVersionForOpenEBSVersion[r.OpenEBS.Spec.Version]; exist {
			r.OpenEBS.Spec.NDM.ImageTag = ndmVersion
		} else {
			return errors.Errorf("Failed to get NDM version for the given OpenEBS version: %s",
				r.OpenEBS.Spec.Version)
		}
	}
	// Form the container image for NDM components based on the image prefix
	// and image tag.
	r.OpenEBS.Spec.NDM.Image = r.OpenEBS.Spec.ImagePrefix +
		"node-disk-manager-amd64:" + r.OpenEBS.Spec.NDM.ImageTag
	// set the default values for NDM probes if not already set.
	err := r.setNDMProbeDefaultsIfNotSet()
	if err != nil {
		return errors.Errorf("Error setting NDM probe default values: %+v", err)
	}
	// set the default values for NDM filters if not already set.
	err = r.setNDMFilterDefaultsIfNotSet()
	if err != nil {
		return errors.Errorf("Error setting NDM filter default values: %+v", err)
	}
	// Initialize sparse field if not already set.
	if r.OpenEBS.Spec.NDM.Sparse == nil {
		r.OpenEBS.Spec.NDM.Sparse = &types.Sparse{}
	}
	// set the sparse path if not already set as per the default storage path
	// given.
	if r.OpenEBS.Spec.NDM.Sparse.Path == "" {
		r.OpenEBS.Spec.NDM.Sparse.Path = r.OpenEBS.Spec.DefaultStoragePath + "/sparse"
	}
	// set the default sparse size if not already set
	// TODO: See if this needs to be handled based on OpenEBS version.
	if r.OpenEBS.Spec.NDM.Sparse.Size == "" {
		r.OpenEBS.Spec.NDM.Sparse.Size = types.DefaultNDMSparseSize
	}
	// set the default sparse count if not already set
	// TODO: See if this needs to be handled based on OpenEBS version.
	if r.OpenEBS.Spec.NDM.Sparse.Count == "" {
		r.OpenEBS.Spec.NDM.Sparse.Count = types.DefaultNDMSparseCount
	}
	return nil
}

// Set the NDM probes default if not already set
func (r *Reconciler) setNDMProbeDefaultsIfNotSet() error {
	// Check if NDM probes field is set or not, if not
	// then initialize it in order to fill the defaults.
	if r.OpenEBS.Spec.NDM.Probes == nil {
		r.OpenEBS.Spec.NDM.Probes = &types.NDMProbes{}
	}
	// Initialize Udev probe if not set
	if r.OpenEBS.Spec.NDM.Probes.Udev == nil {
		r.OpenEBS.Spec.NDM.Probes.Udev = &types.ProbeState{}
	}
	// Enable Udev probe if the field is not set i.e. set the
	// value to true.
	// TODO: Validate the values that can be provided for this
	// field.
	if r.OpenEBS.Spec.NDM.Probes.Udev.Enabled == "" {
		r.OpenEBS.Spec.NDM.Probes.Udev.Enabled = types.True
	}
	// Initialize Smart probe if not set
	if r.OpenEBS.Spec.NDM.Probes.Smart == nil {
		r.OpenEBS.Spec.NDM.Probes.Smart = &types.ProbeState{}
	}
	// Enable Smart probe if the field is not set i.e. set the
	// value to true.
	// TODO: Validate the values that can be provided for this
	// field.
	if r.OpenEBS.Spec.NDM.Probes.Smart.Enabled == "" {
		r.OpenEBS.Spec.NDM.Probes.Smart.Enabled = types.True
	}
	// Initialize Seachest probe if not set
	if r.OpenEBS.Spec.NDM.Probes.Seachest == nil {
		r.OpenEBS.Spec.NDM.Probes.Seachest = &types.ProbeState{}
	}
	// Disable Seachest probe if the field is not set i.e. set the
	// value to false.
	// TODO: Validate the values that can be provided for this
	// field.
	if r.OpenEBS.Spec.NDM.Probes.Seachest.Enabled == "" {
		r.OpenEBS.Spec.NDM.Probes.Seachest.Enabled = types.False
	}

	return nil
}

// Set the NDM filters default if not already set
func (r *Reconciler) setNDMFilterDefaultsIfNotSet() error {
	// Initialize NDM filters field if not already set
	if r.OpenEBS.Spec.NDM.Filters == nil {
		r.OpenEBS.Spec.NDM.Filters = &types.NDMFilters{}
	}
	// Initialize OS Disk filter if not already set
	if r.OpenEBS.Spec.NDM.Filters.OSDisk == nil {
		r.OpenEBS.Spec.NDM.Filters.OSDisk = &types.FilterConfigs{}
	}
	// Enable OS disk filter if the field is not set i.e. set the
	// value to true.
	// TODO: Validate the values that can be provided for this
	// field.
	if r.OpenEBS.Spec.NDM.Filters.OSDisk.Enabled == "" {
		r.OpenEBS.Spec.NDM.Filters.OSDisk.Enabled = types.True
	}
	// Set the OS disk filter's exclude value if not set
	if r.OpenEBS.Spec.NDM.Filters.OSDisk.Exclude == nil {
		r.OpenEBS.Spec.NDM.Filters.OSDisk.Exclude = new(string)
		*r.OpenEBS.Spec.NDM.Filters.OSDisk.Exclude = "/,/etc/hosts,/boot"
	}
	// Initialize Vendor filter if not already set
	if r.OpenEBS.Spec.NDM.Filters.Vendor == nil {
		r.OpenEBS.Spec.NDM.Filters.Vendor = &types.FilterConfigs{}
	}
	// Enable vendor filter if the field is not set i.e. set the
	// value to true.
	// TODO: Validate the values that can be provided for this
	// field.
	if r.OpenEBS.Spec.NDM.Filters.Vendor.Enabled == "" {
		r.OpenEBS.Spec.NDM.Filters.Vendor.Enabled = types.True
	}
	// Set the vendor filter's exclude value if not set
	if r.OpenEBS.Spec.NDM.Filters.Vendor.Exclude == nil {
		r.OpenEBS.Spec.NDM.Filters.Vendor.Exclude = new(string)
		*r.OpenEBS.Spec.NDM.Filters.Vendor.Exclude = "CLOUDBYT,OpenEBS"
	}
	// Set the vendor filter's include value if not set
	if r.OpenEBS.Spec.NDM.Filters.Vendor.Include == nil {
		r.OpenEBS.Spec.NDM.Filters.Vendor.Include = new(string)
		*r.OpenEBS.Spec.NDM.Filters.Vendor.Include = ""
	}
	// Initialize Path filter if not already set
	if r.OpenEBS.Spec.NDM.Filters.Path == nil {
		r.OpenEBS.Spec.NDM.Filters.Path = &types.FilterConfigs{}
	}
	// Enable path filter if the field is not set i.e. set the
	// value to true.
	// TODO: Validate the values that can be provided for this
	// field.
	if r.OpenEBS.Spec.NDM.Filters.Path.Enabled == "" {
		r.OpenEBS.Spec.NDM.Filters.Path.Enabled = types.True
	}
	// Set the path filter's exclude value if not set
	if r.OpenEBS.Spec.NDM.Filters.Path.Exclude == nil {
		r.OpenEBS.Spec.NDM.Filters.Path.Exclude = new(string)
		*r.OpenEBS.Spec.NDM.Filters.Path.Exclude =
			"loop,/dev/fd0,/dev/sr0,/dev/ram,/dev/dm-,/dev/md"
	}
	// Set the path filter's include value if not set
	if r.OpenEBS.Spec.NDM.Filters.Path.Include == nil {
		r.OpenEBS.Spec.NDM.Filters.Path.Include = new(string)
		*r.OpenEBS.Spec.NDM.Filters.Path.Include = ""
	}
	return nil
}

// Set the NDM operator default values if not set
func (r *Reconciler) setNDMOperatorDefaultsIfNotSet() error {
	// Initialize NDM operator field if not already set
	if r.OpenEBS.Spec.NDMOperator == nil {
		r.OpenEBS.Spec.NDMOperator = &types.NDMOperator{}
	}
	// Enable NDM operator if the field is not set i.e. set the
	// value to true
	// TODO: Validate the values that can be provided for this
	// field.
	if r.OpenEBS.Spec.NDMOperator.Enabled == "" {
		r.OpenEBS.Spec.NDMOperator.Enabled = types.True
	}
	// set the NDM operator image as per the given config values
	if r.OpenEBS.Spec.NDMOperator.ImageTag == "" {
		if ndmOperatorVersion, exist :=
			supportedNDMVersionForOpenEBSVersion[r.OpenEBS.Spec.Version]; exist {
			r.OpenEBS.Spec.NDMOperator.ImageTag = ndmOperatorVersion
		} else {
			return errors.Errorf(
				"Failed to get NDM Operator version for the given OpenEBS version: %s",
				r.OpenEBS.Spec.Version)
		}
	}
	// Form the NDM operator image as per the image prefix and image tag.
	r.OpenEBS.Spec.NDMOperator.Image = r.OpenEBS.Spec.ImagePrefix +
		"node-disk-operator-amd64:" + r.OpenEBS.Spec.NDMOperator.ImageTag
	// set the replicas value for NDM operator to 1
	if r.OpenEBS.Spec.NDMOperator.Replicas == nil {
		r.OpenEBS.Spec.NDMOperator.Replicas = new(int32)
		*r.OpenEBS.Spec.NDMOperator.Replicas = types.DefaultNDMOperatorReplicaCount
	}
	return nil
}

// updateNDMConfig updates/sets the default values for ndm configmap
// as per the values provided in the OpenEBS CR.
func (r *Reconciler) updateNDMConfig(configmap *corev1.ConfigMap) error {
	// Initialize NDM config data structure i.e., data field of the configmap
	// in order to form the data to be put in the configmap with the updated values.
	ndmConfigData := &types.NDMConfig{}
	// get the configmap template which we will use as a structure to fill
	// in the given/default values.
	ndmConfigDataTemplate := configmap.Data["node-disk-manager.config"]
	err := yaml.Unmarshal([]byte(ndmConfigDataTemplate), ndmConfigData)
	if err != nil {
		return errors.Errorf("Error unmarshalling NDM config data: %+v, Error: %+v", ndmConfigDataTemplate, err)
	}
	// Enable/disable probes as per the given values in OpenEBS CR or the
	// default values
	for i, probeConfig := range ndmConfigData.ProbeConfigs {
		if probeConfig.Key == types.UdevProbeKey {
			probeConfig.State = r.OpenEBS.Spec.NDM.Probes.Udev.Enabled
		} else if probeConfig.Key == types.SmartProbeKey {
			probeConfig.State = r.OpenEBS.Spec.NDM.Probes.Smart.Enabled
		} else if probeConfig.Key == types.SeachestProbeKey {
			probeConfig.State = r.OpenEBS.Spec.NDM.Probes.Seachest.Enabled
		}
		// update the updates probes in the NDM configmap
		ndmConfigData.ProbeConfigs[i] = probeConfig
	}
	// Enable/disable filters as per the given values in OpenEBS CR or the
	// default values
	for i, filterConfig := range ndmConfigData.FilterConfigs {
		if filterConfig.Key == types.OSDiskFilterKey {
			filterConfig.State = r.OpenEBS.Spec.NDM.Filters.OSDisk.Enabled
			filterConfig.Exclude = r.OpenEBS.Spec.NDM.Filters.OSDisk.Exclude

		} else if filterConfig.Key == types.VendorFilterKey {
			filterConfig.State = r.OpenEBS.Spec.NDM.Filters.Vendor.Enabled
			filterConfig.Exclude = r.OpenEBS.Spec.NDM.Filters.Vendor.Exclude
			filterConfig.Include = r.OpenEBS.Spec.NDM.Filters.Vendor.Include
		} else if filterConfig.Key == types.PathFilterKey {
			filterConfig.State = r.OpenEBS.Spec.NDM.Filters.Path.Enabled
			filterConfig.Exclude = r.OpenEBS.Spec.NDM.Filters.Path.Exclude
			filterConfig.Include = r.OpenEBS.Spec.NDM.Filters.Path.Include
		}
		// update the updated filters in the NDM configmap
		ndmConfigData.FilterConfigs[i] = filterConfig
	}
	ndmConfigDataString, err := yaml.Marshal(ndmConfigData)
	if err != nil {
		return errors.Errorf("Error marshalling configmap data: %+v", err)
	}
	configmap.Data["node-disk-manager.config"] = string(ndmConfigDataString)

	return nil
}

// updateNDM updates the NDM structure as per the provided values otherwise
// default values.
func (r *Reconciler) updateNDM(daemonset *appsv1.DaemonSet) error {
	// Update the volume details
	for i, volume := range daemonset.Spec.Template.Spec.Volumes {
		if volume.Name == "sparsepath" {
			volume.HostPath.Path = r.OpenEBS.Spec.NDM.Sparse.Path
		}
		daemonset.Spec.Template.Spec.Volumes[i] = volume
	}
	// update the containers as per the values given
	for i, container := range daemonset.Spec.Template.Spec.Containers {
		for _, vm := range container.VolumeMounts {
			if vm.Name == "sparsepath" {
				vm.MountPath = r.OpenEBS.Spec.NDM.Sparse.Path
			}
		}
		// update the ENVs
		for _, env := range container.Env {
			if env.Name == types.SparseFileDirectoryEnv {
				env.Value = r.OpenEBS.Spec.NDM.Sparse.Path
			} else if env.Name == types.SparseFileSizeEnv {
				env.Value = r.OpenEBS.Spec.NDM.Sparse.Size
			} else if env.Name == types.SparseFileCountEnv {
				env.Value = r.OpenEBS.Spec.NDM.Sparse.Count
			}
		}
		daemonset.Spec.Template.Spec.Containers[i] = container
	}
	return nil
}
