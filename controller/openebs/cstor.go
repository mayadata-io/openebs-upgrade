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
	appsv1 "k8s.io/api/apps/v1"
	"mayadata.io/openebs-upgrade/types"
)

const (
	// ContainerOpenEBSCSIPlugin is the name of the container openebs csi plugin
	ContainerOpenEBSCSIPluginName string = "openebs-csi-plugin"
	// EnvOpenEBSNamespaceKey is the env key for openebs namespace
	EnvOpenEBSNamespaceKey string = "OPENEBS_NAMESPACE"
	// NamespaceKubeSystem is the value of kube-system namespace
	NamespaceKubeSystem string = "kube-system"
)

// Set the default values for Cstor if not already given.
func (p *Planner) setCStorDefaultsIfNotSet() error {
	if p.ObservedOpenEBS.Spec.CstorConfig == nil {
		p.ObservedOpenEBS.Spec.CstorConfig = &types.CstorConfig{}
	}
	// form the cstor-pool image
	if p.ObservedOpenEBS.Spec.CstorConfig.Pool.ImageTag == "" {
		p.ObservedOpenEBS.Spec.CstorConfig.Pool.ImageTag = p.ObservedOpenEBS.Spec.Version
	}
	p.ObservedOpenEBS.Spec.CstorConfig.Pool.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"cstor-pool:" + p.ObservedOpenEBS.Spec.CstorConfig.Pool.ImageTag

	// form the cstor-pool-mgmt image
	if p.ObservedOpenEBS.Spec.CstorConfig.PoolMgmt.ImageTag == "" {
		p.ObservedOpenEBS.Spec.CstorConfig.PoolMgmt.ImageTag = p.ObservedOpenEBS.Spec.Version
	}
	p.ObservedOpenEBS.Spec.CstorConfig.PoolMgmt.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"cstor-pool-mgmt:" + p.ObservedOpenEBS.Spec.CstorConfig.PoolMgmt.ImageTag

	// form the cstor-istgt image
	if p.ObservedOpenEBS.Spec.CstorConfig.Target.ImageTag == "" {
		p.ObservedOpenEBS.Spec.CstorConfig.Target.ImageTag = p.ObservedOpenEBS.Spec.Version
	}
	p.ObservedOpenEBS.Spec.CstorConfig.Target.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"cstor-istgt:" + p.ObservedOpenEBS.Spec.CstorConfig.Target.ImageTag

	// form the cstor-volume-mgmt image
	if p.ObservedOpenEBS.Spec.CstorConfig.VolumeMgmt.ImageTag == "" {
		p.ObservedOpenEBS.Spec.CstorConfig.VolumeMgmt.ImageTag = p.ObservedOpenEBS.Spec.Version
	}
	p.ObservedOpenEBS.Spec.CstorConfig.VolumeMgmt.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"cstor-volume-mgmt:" + p.ObservedOpenEBS.Spec.CstorConfig.VolumeMgmt.ImageTag

	if r.OpenEBS.Spec.CstorConfig.CStorCSI.Enabled == nil {
		r.OpenEBS.Spec.CstorConfig.CStorCSI.Enabled = new(bool)
		*r.OpenEBS.Spec.CstorConfig.CStorCSI.Enabled = true
	}

	return nil
}

func (r *Reconciler) updateOpenEBSCStorCSINode(daemonset *appsv1.DaemonSet) {
	daemonset.Namespace = NamespaceKubeSystem

	for i, container := range daemonset.Spec.Template.Spec.Containers {

		if container.Name == ContainerOpenEBSCSIPluginName {
			for j, env := range container.Env {
				if env.Name == EnvOpenEBSNamespaceKey {
					env.Value = r.OpenEBS.Namespace
				}
				container.Env[j] = env
			}
		}
		daemonset.Spec.Template.Spec.Containers[i] = container
	}
}

func (r *Reconciler) updateOpenEBSCStorCSIController(statefulset *appsv1.StatefulSet) {
	statefulset.Namespace = NamespaceKubeSystem

	for i, container := range statefulset.Spec.Template.Spec.Containers {

		if container.Name == ContainerOpenEBSCSIPluginName {
			for j, env := range container.Env {
				if env.Name == EnvOpenEBSNamespaceKey {
					env.Value = r.OpenEBS.Namespace
				}
				container.Env[j] = env
			}
		}
		statefulset.Spec.Template.Spec.Containers[i] = container
	}
}
