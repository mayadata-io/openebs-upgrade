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

	"github.com/mayadata-io/openebs-operator/types"
	appsv1 "k8s.io/api/apps/v1"
)

const (
	// DefaultAPIServerReplicaCount is the default value of replica for
	// API server.
	DefaultAPIServerReplicaCount int32 = 1
)

// MayaAPIServer is a wrapper over k8s deployment structure.
type MayaAPIServer struct {
	Object *appsv1.Deployment `json:"object"`
}

// updateManifest updates the MayaAPIServer manifest as per the reconciler.OpenEBS values.
func (apiServer *MayaAPIServer) updateManifest(r *Reconciler) (*MayaAPIServer, error) {
	// Update the container values
	for i, container := range apiServer.Object.Spec.Template.Spec.Containers {
		// Update the container's ENVs.
		for i, env := range container.Env {
			if env.Name == "OPENEBS_IO_INSTALL_DEFAULT_CSTOR_SPARSE_POOL" {
				env.Value = r.OpenEBS.Spec.APIServer.Sparse.Enabled
			} else if env.Name == "OPENEBS_IO_CREATE_DEFAULT_STORAGE_CONFIG" {
				env.Value = r.OpenEBS.Spec.CreateDefaultStorageConfig
			} else if env.Name == "OPENEBS_IO_JIVA_CONTROLLER_IMAGE" {
				env.Value = r.OpenEBS.Spec.Jiva.Image
			} else if env.Name == "OPENEBS_IO_JIVA_REPLICA_IMAGE" {
				env.Value = r.OpenEBS.Spec.Jiva.Image
			} else if env.Name == "OPENEBS_IO_JIVA_REPLICA_COUNT" {
				env.Value = strconv.FormatInt(int64(*r.OpenEBS.Spec.Jiva.Replicas), 10)
			} else if env.Name == "OPENEBS_IO_CSTOR_TARGET_IMAGE" {
				env.Value = r.OpenEBS.Spec.Cstor.Target.Image
			} else if env.Name == "OPENEBS_IO_CSTOR_POOL_IMAGE" {
				env.Value = r.OpenEBS.Spec.Cstor.Pool.Image
			} else if env.Name == "OPENEBS_IO_CSTOR_POOL_MGMT_IMAGE" {
				env.Value = r.OpenEBS.Spec.Cstor.PoolMgmt.Image
			} else if env.Name == "OPENEBS_IO_CSTOR_VOLUME_MGMT_IMAGE" {
				env.Value = r.OpenEBS.Spec.Cstor.VolumeMgmt.Image
			} else if env.Name == "OPENEBS_IO_VOLUME_MONITOR_IMAGE" {
				env.Value = r.OpenEBS.Spec.Policies.Monitoring.Image
			} else if env.Name == "OPENEBS_IO_CSTOR_POOL_EXPORTER_IMAGE" {
				env.Value = r.OpenEBS.Spec.Policies.Monitoring.Image
			} else if env.Name == "OPENEBS_IO_HELPER_IMAGE" {
				env.Value = r.OpenEBS.Spec.Helper.Image
			} else if env.Name == "OPENEBS_IO_ENABLE_ANALYTICS" {
				env.Value = r.OpenEBS.Spec.Analytics.Enabled
			}
			// update container with the updated ENV
			container.Env[i] = env
		}
		// Update the APIServer container with updated container values
		apiServer.Object.Spec.Template.Spec.Containers[i] = container
	}
	return apiServer, nil
}

// setAPIServerDefaultsIfNotSet sets the default values for APIServer if not
// set.
func (r *Reconciler) setAPIServerDefaultsIfNotSet() error {
	if r.OpenEBS.Spec.APIServer == nil {
		r.OpenEBS.Spec.APIServer = &types.APIServer{}
	}
	if r.OpenEBS.Spec.APIServer.Enabled == "" {
		r.OpenEBS.Spec.APIServer.Enabled = types.True
	}
	if r.OpenEBS.Spec.APIServer.ImageTag == "" {
		r.OpenEBS.Spec.APIServer.ImageTag = r.OpenEBS.Spec.Version
	}
	// form the container image as per the image prefix and image tag.
	r.OpenEBS.Spec.APIServer.Image = r.OpenEBS.Spec.ImagePrefix + "m-apiserver:" +
		r.OpenEBS.Spec.APIServer.ImageTag

	if r.OpenEBS.Spec.APIServer.Sparse == nil {
		r.OpenEBS.Spec.APIServer.Sparse = &types.SparsePools{}
	}
	// Sparse pools will be disabled by default.
	if r.OpenEBS.Spec.APIServer.Sparse.Enabled == "" {
		r.OpenEBS.Spec.APIServer.Sparse.Enabled = types.False
	}
	if r.OpenEBS.Spec.APIServer.Replicas == nil {
		r.OpenEBS.Spec.APIServer.Replicas = new(int32)
		*r.OpenEBS.Spec.APIServer.Replicas = DefaultAPIServerReplicaCount
	}
	return nil
}
