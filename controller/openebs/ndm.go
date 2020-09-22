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
	"strconv"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/types"
	"mayadata.io/openebs-upgrade/unstruct"
)

// SupportedNDMVersionForOpenEBSVersion stores the mapping for
// NDM to OpenEBS version i.e., a NDM version for each of the
// supported OpenEBS versions.
// Note: This will be referred to form the container images in
// order to install/update NDM components for a particular OpenEBS
// version.
var SupportedNDMVersionForOpenEBSVersion = map[string]string{
	types.OpenEBSVersion150:    types.NDMVersion045,
	types.OpenEBSVersion160:    types.NDMVersion046,
	types.OpenEBSVersion170:    types.NDMVersion047,
	types.OpenEBSVersion180:    types.NDMVersion048,
	types.OpenEBSVersion190:    types.NDMVersion049,
	types.OpenEBSVersion190EE:  types.NDMVersion049EE,
	types.OpenEBSVersion1100:   types.NDMVersion050,
	types.OpenEBSVersion1100EE: types.NDMVersion050EE,
	types.OpenEBSVersion1110:   types.NDMVersion060,
	types.OpenEBSVersion1110EE: types.NDMVersion060EE,
	types.OpenEBSVersion1120:   types.NDMVersion070,
	types.OpenEBSVersion1120EE: types.NDMVersion070EE,
	types.OpenEBSVersion200:    types.NDMVersion080,
	types.OpenEBSVersion200EE:  types.NDMVersion080EE,
}

// add/update NDM defaults if not already provided
func (p *Planner) setNDMDefaultsIfNotSet() error {
	// Check if NDM field is set or not, if not then
	// initialize it.
	if p.ObservedOpenEBS.Spec.NDMDaemon == nil {
		p.ObservedOpenEBS.Spec.NDMDaemon = &types.NDMDaemon{}
	}
	// Check if NDM enabling field is set or not,
	// if not then set it to true.
	// Note: By default, NDM will always be enabled i.e.
	// will be installed alongwith other OpenEBS components.
	if p.ObservedOpenEBS.Spec.NDMDaemon.Enabled == nil {
		p.ObservedOpenEBS.Spec.NDMDaemon.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.NDMDaemon.Enabled = true
	}
	// set the name with which ndm daemon will be deployed
	if len(p.ObservedOpenEBS.Spec.NDMDaemon.Name) == 0 {
		p.ObservedOpenEBS.Spec.NDMDaemon.Name = types.NDMNameKey
	}
	// Check if imageTag fiels is set or not, if not
	// then set the NDM image tag as per the OpenEBS version
	// given.
	if p.ObservedOpenEBS.Spec.NDMDaemon.ImageTag == "" {
		if ndmVersion, exist :=
			SupportedNDMVersionForOpenEBSVersion[p.ObservedOpenEBS.Spec.Version]; exist {
			p.ObservedOpenEBS.Spec.NDMDaemon.ImageTag = ndmVersion +
				p.ObservedOpenEBS.Spec.ImageTagSuffix
		} else {
			return errors.Errorf("Failed to get NDM version for the given OpenEBS version: %s",
				p.ObservedOpenEBS.Spec.Version)
		}
	}
	// Form the container image for NDM components based on the image prefix
	// and image tag.
	p.ObservedOpenEBS.Spec.NDMDaemon.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"node-disk-manager-amd64:" + p.ObservedOpenEBS.Spec.NDMDaemon.ImageTag
	// set the default values for NDM probes if not already set.
	err := p.setNDMProbeDefaultsIfNotSet()
	if err != nil {
		return errors.Errorf("Error setting NDM probe default values: %+v", err)
	}
	// set the default values for NDM filters if not already set.
	err = p.setNDMFilterDefaultsIfNotSet()
	if err != nil {
		return errors.Errorf("Error setting NDM filter default values: %+v", err)
	}
	// Initialize sparse field if not already set.
	if p.ObservedOpenEBS.Spec.NDMDaemon.Sparse == nil {
		p.ObservedOpenEBS.Spec.NDMDaemon.Sparse = &types.Sparse{}
	}
	// set the sparse path if not already set as per the default storage path
	// given.
	if p.ObservedOpenEBS.Spec.NDMDaemon.Sparse.Path == "" {
		p.ObservedOpenEBS.Spec.NDMDaemon.Sparse.Path = p.ObservedOpenEBS.Spec.DefaultStoragePath + "/sparse"
	}
	// set the default sparse size if not already set
	// TODO: See if this needs to be handled based on OpenEBS version.
	if p.ObservedOpenEBS.Spec.NDMDaemon.Sparse.Size == "" {
		p.ObservedOpenEBS.Spec.NDMDaemon.Sparse.Size = types.DefaultNDMSparseSize
	}
	// set the default sparse count if not already set
	// TODO: See if this needs to be handled based on OpenEBS version.
	if p.ObservedOpenEBS.Spec.NDMDaemon.Sparse.Count == "" {
		p.ObservedOpenEBS.Spec.NDMDaemon.Sparse.Count = types.DefaultNDMSparseCount
	}
	return nil
}

// Set the NDM probes default if not already set
func (p *Planner) setNDMProbeDefaultsIfNotSet() error {
	// Check if NDM probes field is set or not, if not
	// then initialize it in order to fill the defaults.
	if p.ObservedOpenEBS.Spec.NDMDaemon.Probes == nil {
		p.ObservedOpenEBS.Spec.NDMDaemon.Probes = &types.NDMProbes{}
	}
	// Initialize Udev probe if not set
	if p.ObservedOpenEBS.Spec.NDMDaemon.Probes.Udev == nil {
		p.ObservedOpenEBS.Spec.NDMDaemon.Probes.Udev = &types.ProbeState{}
	}
	// Enable Udev probe if the field is not set i.e. set the
	// value to true.
	// TODO: Validate the values that can be provided for this
	// field.
	if p.ObservedOpenEBS.Spec.NDMDaemon.Probes.Udev.Enabled == nil {
		p.ObservedOpenEBS.Spec.NDMDaemon.Probes.Udev.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.NDMDaemon.Probes.Udev.Enabled = true
	}
	// Initialize Smart probe if not set
	if p.ObservedOpenEBS.Spec.NDMDaemon.Probes.Smart == nil {
		p.ObservedOpenEBS.Spec.NDMDaemon.Probes.Smart = &types.ProbeState{}
	}
	// Enable Smart probe if the field is not set i.e. set the
	// value to true.
	// TODO: Validate the values that can be provided for this
	// field.
	if p.ObservedOpenEBS.Spec.NDMDaemon.Probes.Smart.Enabled == nil {
		p.ObservedOpenEBS.Spec.NDMDaemon.Probes.Smart.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.NDMDaemon.Probes.Smart.Enabled = true
	}
	// Initialize Seachest probe if not set
	if p.ObservedOpenEBS.Spec.NDMDaemon.Probes.Seachest == nil {
		p.ObservedOpenEBS.Spec.NDMDaemon.Probes.Seachest = &types.ProbeState{}
	}
	// Disable Seachest probe if the field is not set i.e. set the
	// value to false.
	// TODO: Validate the values that can be provided for this
	// field.
	if p.ObservedOpenEBS.Spec.NDMDaemon.Probes.Seachest.Enabled == nil {
		p.ObservedOpenEBS.Spec.NDMDaemon.Probes.Seachest.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.NDMDaemon.Probes.Seachest.Enabled = false
	}

	return nil
}

// Set the NDM filters default if not already set
func (p *Planner) setNDMFilterDefaultsIfNotSet() error {
	// Initialize NDM filters field if not already set
	if p.ObservedOpenEBS.Spec.NDMDaemon.Filters == nil {
		p.ObservedOpenEBS.Spec.NDMDaemon.Filters = &types.NDMFilters{}
	}
	// Initialize OS Disk filter if not already set
	if p.ObservedOpenEBS.Spec.NDMDaemon.Filters.OSDisk == nil {
		p.ObservedOpenEBS.Spec.NDMDaemon.Filters.OSDisk = &types.FilterConfigs{}
	}
	// Enable OS disk filter if the field is not set i.e. set the
	// value to true.
	// TODO: Validate the values that can be provided for this
	// field.
	if p.ObservedOpenEBS.Spec.NDMDaemon.Filters.OSDisk.Enabled == nil {
		p.ObservedOpenEBS.Spec.NDMDaemon.Filters.OSDisk.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.NDMDaemon.Filters.OSDisk.Enabled = true
	}
	// Set the OS disk filter's exclude value if not set
	if p.ObservedOpenEBS.Spec.NDMDaemon.Filters.OSDisk.Exclude == nil {
		p.ObservedOpenEBS.Spec.NDMDaemon.Filters.OSDisk.Exclude = new(string)
		*p.ObservedOpenEBS.Spec.NDMDaemon.Filters.OSDisk.Exclude = "/,/etc/hosts,/boot"
	}
	// Initialize Vendor filter if not already set
	if p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Vendor == nil {
		p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Vendor = &types.FilterConfigs{}
	}
	// Enable vendor filter if the field is not set i.e. set the
	// value to true.
	// TODO: Validate the values that can be provided for this
	// field.
	if p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Vendor.Enabled == nil {
		p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Vendor.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Vendor.Enabled = true
	}
	// Set the vendor filter's exclude value if not set
	if p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Vendor.Exclude == nil {
		p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Vendor.Exclude = new(string)
		*p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Vendor.Exclude = "CLOUDBYT,OpenEBS"
	}
	// Set the vendor filter's include value if not set
	if p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Vendor.Include == nil {
		p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Vendor.Include = new(string)
		*p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Vendor.Include = ""
	}
	// Initialize Path filter if not already set
	if p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Path == nil {
		p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Path = &types.FilterConfigs{}
	}
	// Enable path filter if the field is not set i.e. set the
	// value to true.
	// TODO: Validate the values that can be provided for this
	// field.
	if p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Path.Enabled == nil {
		p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Path.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Path.Enabled = true
	}
	// Set the path filter's exclude value if not set
	if p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Path.Exclude == nil {
		p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Path.Exclude = new(string)
		*p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Path.Exclude =
			"loop,/dev/fd0,/dev/sr0,/dev/ram,/dev/dm-,/dev/md"
	}
	// Set the path filter's include value if not set
	if p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Path.Include == nil {
		p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Path.Include = new(string)
		*p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Path.Include = ""
	}
	return nil
}

// Set the NDM operator default values if not set
func (p *Planner) setNDMOperatorDefaultsIfNotSet() error {
	// Initialize NDM operator field if not already set
	if p.ObservedOpenEBS.Spec.NDMOperator == nil {
		p.ObservedOpenEBS.Spec.NDMOperator = &types.NDMOperator{}
	}
	// Enable NDM operator if the field is not set i.e. set the
	// value to true
	if p.ObservedOpenEBS.Spec.NDMOperator.Enabled == nil {
		p.ObservedOpenEBS.Spec.NDMOperator.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.NDMOperator.Enabled = true
	}
	// set the name with which ndm operator will be deployed
	if len(p.ObservedOpenEBS.Spec.NDMOperator.Name) == 0 {
		p.ObservedOpenEBS.Spec.NDMOperator.Name = types.NDMOperatorNameKey
	}
	// set the NDM operator image as per the given config values
	if p.ObservedOpenEBS.Spec.NDMOperator.ImageTag == "" {
		if ndmOperatorVersion, exist :=
			SupportedNDMVersionForOpenEBSVersion[p.ObservedOpenEBS.Spec.Version]; exist {
			p.ObservedOpenEBS.Spec.NDMOperator.ImageTag = ndmOperatorVersion +
				p.ObservedOpenEBS.Spec.ImageTagSuffix
		} else {
			return errors.Errorf(
				"Failed to get NDM Operator version for the given OpenEBS version: %s",
				p.ObservedOpenEBS.Spec.Version)
		}
	}
	// Form the NDM operator image as per the image prefix and image tag.
	p.ObservedOpenEBS.Spec.NDMOperator.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"node-disk-operator-amd64:" + p.ObservedOpenEBS.Spec.NDMOperator.ImageTag
	// set the replicas value for NDM operator to 1
	if p.ObservedOpenEBS.Spec.NDMOperator.Replicas == nil {
		p.ObservedOpenEBS.Spec.NDMOperator.Replicas = new(int32)
		*p.ObservedOpenEBS.Spec.NDMOperator.Replicas = types.DefaultNDMOperatorReplicaCount
	}
	return nil
}

// Set the NDM configmap default values if not set
func (p *Planner) setNDMConfigMapDefaultsIfNotSet() error {
	// Initialize NDM configmap field if not already set
	if p.ObservedOpenEBS.Spec.NDMConfigMap == nil {
		p.ObservedOpenEBS.Spec.NDMConfigMap = &types.NDMConfigMap{}
	}
	// set the name of NDM configMap.
	if len(p.ObservedOpenEBS.Spec.NDMConfigMap.Name) == 0 {
		p.ObservedOpenEBS.Spec.NDMConfigMap.Name = types.NDMConfigNameKey
	}
	return nil
}

// updateNDMConfig updates/sets the default values for ndm configmap
// as per the values provided in the OpenEBS CR.
func (p *Planner) updateNDMConfig(configmap *unstructured.Unstructured) error {
	configmap.SetName(p.ObservedOpenEBS.Spec.NDMConfigMap.Name)
	// Initialize NDM config data structure i.e., data field of the configmap
	// in order to form the data to be put in the configmap with the updated values.
	ndmConfigData := &types.NDMConfig{}
	// get the configmap template which we will use as a structure to fill
	// in the given/default values.
	dataMap, _, err := unstructured.NestedMap(configmap.Object, "data")
	if err != nil {
		return err
	}
	ndmConfigDataTemplate := dataMap["node-disk-manager.config"]
	err = yaml.Unmarshal([]byte(ndmConfigDataTemplate.(string)), ndmConfigData)
	if err != nil {
		return errors.Errorf("Error unmarshalling NDM config data: %+v, Error: %+v", ndmConfigDataTemplate, err)
	}
	// Enable/disable probes as per the given values in OpenEBS CR or the
	// default values
	for i, probeConfig := range ndmConfigData.ProbeConfigs {
		if probeConfig.Key == types.UdevProbeKey {
			probeConfig.State = strconv.FormatBool(
				*p.ObservedOpenEBS.Spec.NDMDaemon.Probes.Udev.Enabled)
		} else if probeConfig.Key == types.SmartProbeKey {
			probeConfig.State = strconv.FormatBool(
				*p.ObservedOpenEBS.Spec.NDMDaemon.Probes.Smart.Enabled)
		} else if probeConfig.Key == types.SeachestProbeKey {
			probeConfig.State = strconv.FormatBool(
				*p.ObservedOpenEBS.Spec.NDMDaemon.Probes.Seachest.Enabled)
		}
		// update the updates probes in the NDM configmap
		ndmConfigData.ProbeConfigs[i] = probeConfig
	}
	// Enable/disable filters as per the given values in OpenEBS CR or the
	// default values
	for i, filterConfig := range ndmConfigData.FilterConfigs {
		if filterConfig.Key == types.OSDiskFilterKey {
			filterConfig.State = strconv.FormatBool(
				*p.ObservedOpenEBS.Spec.NDMDaemon.Filters.OSDisk.Enabled)
			filterConfig.Exclude = p.ObservedOpenEBS.Spec.NDMDaemon.Filters.OSDisk.Exclude

		} else if filterConfig.Key == types.VendorFilterKey {
			filterConfig.State = strconv.FormatBool(
				*p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Vendor.Enabled)
			filterConfig.Exclude = p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Vendor.Exclude
			filterConfig.Include = p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Vendor.Include
		} else if filterConfig.Key == types.PathFilterKey {
			filterConfig.State = strconv.FormatBool(
				*p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Path.Enabled)
			filterConfig.Exclude = p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Path.Exclude
			filterConfig.Include = p.ObservedOpenEBS.Spec.NDMDaemon.Filters.Path.Include
		}
		// update the updated filters in the NDM configmap
		ndmConfigData.FilterConfigs[i] = filterConfig
	}
	ndmConfigDataString, err := yaml.Marshal(ndmConfigData)
	if err != nil {
		return errors.Errorf("Error marshalling configmap data: %+v", err)
	}
	dataMap["node-disk-manager.config"] = string(ndmConfigDataString)
	unstructured.SetNestedMap(configmap.Object, dataMap, "data")

	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := configmap.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for openebs-ndm-config configmap
	// 1. openebs-upgrade.dao.mayadata.io/component-group: ndm
	// 2. openebs-upgrade.dao.mayadata.io/component-name: openebs-ndm-config
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSNDMComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.NDMConfigNameKey
	// set the desired labels
	configmap.SetLabels(desiredLabels)

	return nil
}

// updateNDM updates the NDM structure as per the provided values otherwise
// default values.
func (p *Planner) updateNDM(daemonset *unstructured.Unstructured) error {
	daemonset.SetName(p.ObservedOpenEBS.Spec.NDMDaemon.Name)
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := daemonset.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for openebs-ndm-operator deploy
	// 1. openebs-upgrade.dao.mayadata.io/component-group: ndm
	// 2. openebs-upgrade.dao.mayadata.io/component-subgroup: daemon
	// 3. openebs-upgrade.dao.mayadata.io/component-name: openebs-ndm
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSNDMComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentSubGroupLabelKey] =
		types.OpenEBSDaemonComponentSubGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.NDMNameKey
	// set the desired labels
	daemonset.SetLabels(desiredLabels)

	// add hostPID field if provided, this field is used in case if feature-gate "APIService"
	// is enabled for NDM in order to check ISCSI service status.
	ndmTemplateSpec, err := unstruct.GetNestedMapOrError(daemonset, "spec", "template", "spec")
	if err != nil {
		return err
	}
	if p.ObservedOpenEBS.Spec.NDMDaemon.EnableHostPID != nil {
		ndmTemplateSpec["hostPID"] = p.ObservedOpenEBS.Spec.NDMDaemon.EnableHostPID
	}
	err = unstructured.SetNestedMap(daemonset.Object, ndmTemplateSpec, "spec", "template", "spec")
	if err != nil {
		return err
	}
	volumes, err := unstruct.GetNestedSliceOrError(daemonset, "spec", "template", "spec", "volumes")
	if err != nil {
		return err
	}
	updateVolume := func(obj *unstructured.Unstructured) error {
		volumeName, err := unstruct.GetString(obj, "spec", "name")
		if err != nil {
			return err
		}
		if volumeName == "sparsepath" {
			err = unstructured.SetNestedField(obj.Object,
				p.ObservedOpenEBS.Spec.NDMDaemon.Sparse.Path, "spec", "hostPath", "path")
			if err != nil {
				return err
			}
		} else if volumeName == "basepath" {
			err = unstructured.SetNestedField(obj.Object,
				p.ObservedOpenEBS.Spec.DefaultStoragePath+"/ndm", "spec", "hostPath", "path")
			if err != nil {
				return err
			}
		}
		return nil
	}
	err = unstruct.SliceIterator(volumes).ForEachUpdate(updateVolume)
	if err != nil {
		return err
	}

	err = unstructured.SetNestedSlice(daemonset.Object, volumes,
		"spec", "template", "spec", "volumes")
	if err != nil {
		return err
	}

	// update the daemonset containers
	containers, err := unstruct.GetNestedSliceOrError(daemonset, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}
	// this function updates the envs of node-disk-manager container
	updateNodeDiskManagerEnv := func(env *unstructured.Unstructured) error {
		envName, _, err := unstructured.NestedString(env.Object, "spec", "name")
		if err != nil {
			return err
		}
		if envName == types.SparseFileDirectoryEnv {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.NDMDaemon.Sparse.Path, "spec", "value")
		} else if envName == types.SparseFileSizeEnv {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.NDMDaemon.Sparse.Size, "spec", "value")
		} else if envName == types.SparseFileCountEnv {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.NDMDaemon.Sparse.Count, "spec", "value")
		}
		if err != nil {
			return err
		}

		return nil
	}
	// this function updates the volumeMounts of node-disk-manager container
	updateNodeDiskManagerVolumeMount := func(vm *unstructured.Unstructured) error {
		vmName, _, err := unstructured.NestedString(vm.Object, "spec", "name")
		if err != nil {
			return err
		}
		if vmName == "sparsepath" {
			err = unstructured.SetNestedField(vm.Object,
				p.ObservedOpenEBS.Spec.NDMDaemon.Sparse.Path, "spec", "mountPath")
			if err != nil {
				return err
			}
		} else if vmName == "basepath" {
			err = unstructured.SetNestedField(vm.Object,
				p.ObservedOpenEBS.Spec.DefaultStoragePath+"/ndm", "spec", "mountPath")
			if err != nil {
				return err
			}
		}
		return nil
	}
	updateContainer := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}
		envs, _, err := unstruct.GetSlice(obj, "spec", "env")
		if err != nil {
			return err
		}
		volumeMounts, _, err := unstruct.GetSlice(obj, "spec", "volumeMounts")
		if err != nil {
			return err
		}
		args, _, err := unstruct.GetSlice(obj, "spec", "args")
		if err != nil {
			return err
		}
		// update the envs and volume mounts per container i.e., envs and volume mounts
		// could be different for each container and should be updated as such.
		if containerName == "node-disk-manager" {
			// update the container name if not same.
			if len(p.ObservedOpenEBS.Spec.NDMDaemon.ContainerName) != 0 {
				err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.NDMDaemon.ContainerName, "spec", "name")
				if err != nil {
					return err
				}
			}
			// Set the image of the container.
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.NDMDaemon.Image, "spec", "image")
			if err != nil {
				return err
			}
			err = unstruct.SliceIterator(envs).ForEachUpdate(updateNodeDiskManagerEnv)
			if err != nil {
				return err
			}
			err = unstruct.SliceIterator(volumeMounts).ForEachUpdate(updateNodeDiskManagerVolumeMount)
			if err != nil {
				return err
			}
			// add the required args for node-disk-manager
			args = p.addNodeDiskManagerArgs(args)
		}
		err = unstructured.SetNestedSlice(obj.Object, envs, "spec", "env")
		if err != nil {
			return err
		}
		err = unstructured.SetNestedSlice(obj.Object, volumeMounts, "spec", "volumeMounts")
		if err != nil {
			return err
		}
		err = unstructured.SetNestedSlice(obj.Object, args, "spec", "args")
		if err != nil {
			return err
		}

		// Set the resource of the containers.
		if p.ObservedOpenEBS.Spec.NDMDaemon.Resources != nil {
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.NDMDaemon.Resources,
				"spec", "resources")
		} else if p.ObservedOpenEBS.Spec.Resources != nil {
			err = unstructured.SetNestedField(obj.Object,
				p.ObservedOpenEBS.Spec.Resources, "spec", "resources")
		}
		if err != nil {
			return err
		}
		return nil
	}
	err = unstruct.SliceIterator(containers).ForEachUpdate(updateContainer)
	if err != nil {
		return err
	}
	err = unstructured.SetNestedSlice(daemonset.Object, containers,
		"spec", "template", "spec", "containers")
	if err != nil {
		return err
	}

	return nil
}

// addNodeDiskManagerArgs adds args to the node-disk-manager container such as features gates, etc.
func (p *Planner) addNodeDiskManagerArgs(args []interface{}) []interface{} {
	if p.ObservedOpenEBS.Spec.NDMDaemon.FeatureGates != nil ||
		len(p.ObservedOpenEBS.Spec.NDMDaemon.FeatureGates) != 0 {
		for _, featureGate := range p.ObservedOpenEBS.Spec.NDMDaemon.FeatureGates {
			// if an empty feature gate is provided then do not do anything
			if len(featureGate) == 0 {
				continue
			}
			args = append(args, "--feature-gates="+featureGate)
		}
	}
	return args
}

// updateNDMOperator updates the NDM Operator structure as per the provided values otherwise
// default values.
func (p *Planner) updateNDMOperator(deploy *unstructured.Unstructured) error {
	deploy.SetName(p.ObservedOpenEBS.Spec.NDMOperator.Name)
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := deploy.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for openebs-ndm-operator deploy
	// 1. openebs-upgrade.dao.mayadata.io/component-group: ndm
	// 2. openebs-upgrade.dao.mayadata.io/component-subgroup: operator
	// 3. openebs-upgrade.dao.mayadata.io/component-name: openebs-ndm-operator
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSNDMComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentSubGroupLabelKey] =
		types.OpenEBSOperatorComponentSubGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.NDMOperatorNameKey
	// set the desired labels
	deploy.SetLabels(desiredLabels)

	// update the daemonset containers
	containers, err := unstruct.GetNestedSliceOrError(deploy, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}
	updateNDMOperatorEnv := func(env *unstructured.Unstructured) error {
		envName, _, err := unstructured.NestedString(env.Object, "spec", "name")
		if err != nil {
			return err
		}
		if envName == types.CleanupJobImageEnv {
			err = unstructured.SetNestedField(env.Object, p.ObservedOpenEBS.Spec.Helper.Image,
				"spec", "value")
			if err != nil {
				return err
			}
		}
		return nil
	}
	updateContainer := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}
		envs, _, err := unstruct.GetSlice(obj, "spec", "env")
		if err != nil {
			return err
		}
		if containerName == "node-disk-operator" {
			// update the container name if not same.
			if len(p.ObservedOpenEBS.Spec.NDMOperator.ContainerName) != 0 {
				err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.NDMOperator.ContainerName, "spec", "name")
				if err != nil {
					return err
				}
			}
			// Set the image of the container.
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.NDMOperator.Image,
				"spec", "image")
			if err != nil {
				return err
			}
			err = unstruct.SliceIterator(envs).ForEachUpdate(updateNDMOperatorEnv)
			if err != nil {
				return err
			}
		} else {
			return nil
		}
		err = unstructured.SetNestedSlice(obj.Object, envs, "spec", "env")
		if err != nil {
			return err
		}
		return nil
	}
	err = unstruct.SliceIterator(containers).ForEachUpdate(updateContainer)
	if err != nil {
		return err
	}

	err = unstructured.SetNestedSlice(deploy.Object, containers,
		"spec", "template", "spec", "containers")
	if err != nil {
		return err
	}

	return nil
}

func (p *Planner) fillNDMOperatorExistingValues(observedComponentDetails ObservedComponentDesiredDetails) error {
	var (
		containerName string
		err           error
	)
	p.ObservedOpenEBS.Spec.NDMOperator.MatchLabels = observedComponentDetails.MatchLabels
	p.ObservedOpenEBS.Spec.NDMOperator.PodTemplateLabels = observedComponentDetails.PodTemplateLabels
	if len(p.ObservedOpenEBS.Spec.NDMOperator.ContainerName) > 0 {
		containerName = p.ObservedOpenEBS.Spec.NDMOperator.ContainerName
	} else {
		containerName = types.NodeDiskOperatorContainerKey
	}
	p.ObservedOpenEBS.Spec.NDMOperator.ENV, err = fetchExistingContainerEnvs(
		observedComponentDetails.Containers, containerName)
	if err != nil {
		return err
	}

	return nil
}

func (p *Planner) fillNDMDaemonExistingValues(observedComponentDetails ObservedComponentDesiredDetails) error {
	var (
		containerName string
		err           error
	)
	p.ObservedOpenEBS.Spec.NDMDaemon.MatchLabels = observedComponentDetails.MatchLabels
	p.ObservedOpenEBS.Spec.NDMDaemon.PodTemplateLabels = observedComponentDetails.PodTemplateLabels
	if len(p.ObservedOpenEBS.Spec.NDMDaemon.ContainerName) > 0 {
		containerName = p.ObservedOpenEBS.Spec.NDMDaemon.ContainerName
	} else {
		containerName = types.NDMDaemonContainerKey
	}
	p.ObservedOpenEBS.Spec.NDMDaemon.ENV, err = fetchExistingContainerEnvs(
		observedComponentDetails.Containers, containerName)
	if err != nil {
		return err
	}

	return nil
}
