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
	corev1 "k8s.io/api/core/v1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
)

// updateNamespace updates the namespace manifest as per the given configuration
// in OpenEBS CR.
func (r *Reconciler) updateNamespace(YAML string) (string, error) {
	ns := &corev1.Namespace{}
	err := yaml.Unmarshal([]byte(YAML), ns)
	if err != nil {
		return "", errors.Errorf("Error unmarshalling namespace YAML: %+v", err)
	}
	ns.Name = r.OpenEBS.Namespace

	rawNamespace, err := yaml.Marshal(ns)
	if err != nil {
		return "", errors.Errorf("Error marshalling namespace: %+v", err)
	}
	return string(rawNamespace), nil
}

// updateServiceAccount updates the service account manifest as per the
// given configuration in OpenEBS CR.
func (r *Reconciler) updateServiceAccount(YAML string) (string, error) {
	serviceAccount := &corev1.ServiceAccount{}
	err := yaml.Unmarshal([]byte(YAML), serviceAccount)
	if err != nil {
		return "", errors.Errorf("Error unmarshalling service account YAML: %+v", err)
	}
	serviceAccount.Namespace = r.OpenEBS.Namespace

	rawServiceAccount, err := yaml.Marshal(serviceAccount)
	if err != nil {
		return "", errors.Errorf("Error marshalling serviceAccount struct: %+v", err)
	}
	return string(rawServiceAccount), nil
}

// updateClusterRoleBinding updates the clusterRoleBinding manifest as per the
// given configuration in OpenEBS CR.
func (r *Reconciler) updateClusterRoleBinding(YAML string) (string, error) {
	clusterRoleBinding := &rbacv1beta1.ClusterRoleBinding{}
	err := yaml.Unmarshal([]byte(YAML), clusterRoleBinding)
	if err != nil {
		return "", errors.Errorf("Error unmarshalling clusterRoleBinding YAML: %+v", err)
	}
	for i, subject := range clusterRoleBinding.Subjects {
		subject.Namespace = r.OpenEBS.Namespace
		clusterRoleBinding.Subjects[i] = subject
	}

	rawClusterRoleBinding, err := yaml.Marshal(clusterRoleBinding)
	if err != nil {
		return "", errors.Errorf("Error marshalling clusterRoleBinding struct: %+v", err)
	}
	return string(rawClusterRoleBinding), nil
}
