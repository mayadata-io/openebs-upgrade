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

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/types"
	"mayadata.io/openebs-upgrade/unstruct"
)

const (
	// DefaultAPIServerReplicaCount is the default value of replica for
	// API server.
	DefaultAPIServerReplicaCount int32 = 1
)

// updateMayaAPIServer updates the MayaAPIServer manifest as per the reconcile.ObservedOpenEBS values.
func (p *Planner) updateMayaAPIServer(deploy *unstructured.Unstructured) error {
	deploy.SetName(p.ObservedOpenEBS.Spec.APIServer.Name)
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := deploy.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for maya-apiserver deploy:
	// 1. openebs-upgrade.dao.mayadata.io/component-group: maya-apiserver
	// 2. openebs-upgrade.dao.mayadata.io/component-name: maya-apiserver
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSMayaAPIServerComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.MayaAPIServerNameKey
	// set the desired labels
	deploy.SetLabels(desiredLabels)

	// get the containers of the maya-apiserver and update the desired fields
	containers, err := unstruct.GetNestedSliceOrError(deploy, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}
	// update the env value of containers
	updateMayaAPIServerEnv := func(env *unstructured.Unstructured) error {
		envName, _, err := unstructured.NestedString(env.Object, "spec", "name")
		if err != nil {
			return err
		}
		if envName == "OPENEBS_IO_INSTALL_DEFAULT_CSTOR_SPARSE_POOL" {
			err = unstructured.SetNestedField(env.Object, strconv.FormatBool(
				*p.ObservedOpenEBS.Spec.APIServer.CstorSparsePool.Enabled), "spec", "value")
		} else if envName == "OPENEBS_IO_CREATE_DEFAULT_STORAGE_CONFIG" {
			err = unstructured.SetNestedField(env.Object, strconv.FormatBool(
				*p.ObservedOpenEBS.Spec.CreateDefaultStorageConfig), "spec", "value")
		} else if envName == "OPENEBS_IO_BASE_DIR" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.DefaultStoragePath, "spec", "value")
		} else if envName == "OPENEBS_IO_CSTOR_TARGET_DIR" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.DefaultStoragePath+"/sparse", "spec", "value")
		} else if envName == "OPENEBS_IO_CSTOR_POOL_SPARSE_DIR" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.DefaultStoragePath+"/sparse", "spec", "value")
		} else if envName == "OPENEBS_IO_JIVA_POOL_DIR" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.DefaultStoragePath, "spec", "value")
		} else if envName == "OPENEBS_IO_LOCALPV_HOSTPATH_DIR" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.DefaultStoragePath+"/local", "spec", "value")
		} else if envName == "OPENEBS_IO_JIVA_CONTROLLER_IMAGE" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.JivaConfig.Image, "spec", "value")
		} else if envName == "OPENEBS_IO_JIVA_REPLICA_IMAGE" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.JivaConfig.Image, "spec", "value")
		} else if envName == "OPENEBS_IO_JIVA_REPLICA_COUNT" {
			envValue := strconv.FormatInt(int64(*p.ObservedOpenEBS.Spec.JivaConfig.Replicas), 10)
			err = unstructured.SetNestedField(env.Object, envValue, "spec", "value")
		} else if envName == "OPENEBS_IO_CSTOR_TARGET_IMAGE" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.CstorConfig.Target.Image, "spec", "value")
		} else if envName == "OPENEBS_IO_CSTOR_POOL_IMAGE" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.CstorConfig.Pool.Image, "spec", "value")
		} else if envName == "OPENEBS_IO_CSTOR_POOL_MGMT_IMAGE" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.CstorConfig.PoolMgmt.Image, "spec", "value")
		} else if envName == "OPENEBS_IO_CSTOR_VOLUME_MGMT_IMAGE" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.CstorConfig.VolumeMgmt.Image, "spec", "value")
		} else if envName == "OPENEBS_IO_VOLUME_MONITOR_IMAGE" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.Policies.Monitoring.Image, "spec", "value")
		} else if envName == "OPENEBS_IO_CSTOR_POOL_EXPORTER_IMAGE" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.Policies.Monitoring.Image, "spec", "value")
		} else if envName == "OPENEBS_IO_HELPER_IMAGE" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.Helper.Image, "spec", "value")
		} else if envName == "OPENEBS_IO_ENABLE_ANALYTICS" {
			envValue := strconv.FormatBool(*p.ObservedOpenEBS.Spec.Analytics.Enabled)
			err = unstructured.SetNestedField(env.Object, envValue, "spec", "value")
		}
		if err != nil {
			return err
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
		// update the envs of maya-apiserver container
		// In order to update envs of other containers, just write an updateEnv
		// function for specific containers.
		if containerName == "maya-apiserver" {
			// Set the image of the container.
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.APIServer.Image,
				"spec", "image")
			if err != nil {
				return err
			}
			err = unstruct.SliceIterator(envs).ForEachUpdate(updateMayaAPIServerEnv)
			if err != nil {
				return err
			}
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
	err = unstructured.SetNestedSlice(deploy.Object,
		containers, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}

	return nil
}

// setAPIServerDefaultsIfNotSet sets the default values for APIServer if not
// set.
func (p *Planner) setAPIServerDefaultsIfNotSet() error {
	if p.ObservedOpenEBS.Spec.APIServer == nil {
		p.ObservedOpenEBS.Spec.APIServer = &types.APIServer{}
	}
	if p.ObservedOpenEBS.Spec.APIServer.Service == nil {
		p.ObservedOpenEBS.Spec.APIServer.Service = &types.APIServerService{}
	}
	if p.ObservedOpenEBS.Spec.APIServer.Enabled == nil {
		p.ObservedOpenEBS.Spec.APIServer.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.APIServer.Enabled = true
	}
	// set the name with which apiserver and apiserver service will be deployed
	if len(p.ObservedOpenEBS.Spec.APIServer.Name) == 0 {
		p.ObservedOpenEBS.Spec.APIServer.Name = types.MayaAPIServerNameKey
	}
	if len(p.ObservedOpenEBS.Spec.APIServer.Service.Name) == 0 {
		p.ObservedOpenEBS.Spec.APIServer.Service.Name = types.MayaAPIServerServiceNameKey
	}
	if p.ObservedOpenEBS.Spec.APIServer.ImageTag == "" {
		p.ObservedOpenEBS.Spec.APIServer.ImageTag = p.ObservedOpenEBS.Spec.Version +
			p.ObservedOpenEBS.Spec.ImageTagSuffix
	}
	// form the container image as per the image prefix and image tag.
	p.ObservedOpenEBS.Spec.APIServer.Image = p.ObservedOpenEBS.Spec.ImagePrefix + "m-apiserver:" +
		p.ObservedOpenEBS.Spec.APIServer.ImageTag

	if p.ObservedOpenEBS.Spec.APIServer.CstorSparsePool == nil {
		p.ObservedOpenEBS.Spec.APIServer.CstorSparsePool = &types.CstorSparsePool{}
	}
	// Sparse pools will be disabled by default.
	if p.ObservedOpenEBS.Spec.APIServer.CstorSparsePool.Enabled == nil {
		p.ObservedOpenEBS.Spec.APIServer.CstorSparsePool.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.APIServer.CstorSparsePool.Enabled = false
	}
	if p.ObservedOpenEBS.Spec.APIServer.Replicas == nil {
		p.ObservedOpenEBS.Spec.APIServer.Replicas = new(int32)
		*p.ObservedOpenEBS.Spec.APIServer.Replicas = DefaultAPIServerReplicaCount
	}
	return nil
}

// updateMayaAPIServerService updates the maya-apiserver-service manifest as per the
// reconcile.ObservedOpenEBS values.
func (p *Planner) updateMayaAPIServerService(svc *unstructured.Unstructured) error {
	svc.SetName(p.ObservedOpenEBS.Spec.APIServer.Service.Name)
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := svc.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for maya-apiserver-service service
	// 1. openebs-upgrade.dao.mayadata.io/component-group: maya-apiserver
	// 2. openebs-upgrade.dao.mayadata.io/component-name: maya-apiserver-service
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSMayaAPIServerComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.MayaAPIServerServiceNameKey
	// set the desired labels
	svc.SetLabels(desiredLabels)

	return nil
}
