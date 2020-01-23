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

package k8s

import (
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterRoleBinding is a wrapper over k8s ClusterRoleBinding
type ClusterRoleBinding struct {
	Object *rbacv1beta1.ClusterRoleBinding `json:"object"`
}

// createOrUpdate checks if the resource provided is present or not, if
// not present then it creates the resource otherwise updates it.
func (crb *ClusterRoleBinding) createOrUpdate() error {
	existingCrb, err := Clientset.RbacV1beta1().ClusterRoleBindings().Get(crb.Object.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			_, err = Clientset.RbacV1beta1().ClusterRoleBindings().Create(crb.Object)
			if err != nil {
				return errors.Errorf("Error creating cluster role binding: %s: %+v", crb.Object.Name, err)
			}
		} else {
			return errors.Errorf("Error getting cluster role binding: %s: %+v", crb.Object.Name, err)
		}

	}
	// Set the resource version of the object to be updated
	crb.Object.SetResourceVersion(existingCrb.GetResourceVersion())
	_, err = Clientset.RbacV1beta1().ClusterRoleBindings().Update(crb.Object)
	if err != nil {
		return errors.Errorf("Error updating cluster role binding: %s: %+v", crb.Object.Name, err)
	}

	return nil
}

// DeployClusterRoleBinding creates/updates a given clusterRoleBinding based on
// the given YAML.
func DeployClusterRoleBinding(YAML string) error {
	crb := &rbacv1beta1.ClusterRoleBinding{}
	err := yaml.Unmarshal([]byte(YAML), crb)
	if err != nil {
		return errors.Errorf(
			"Error unmarshalling clusterRoleBinding YAML: %+v", err)
	}
	clusterRoleBinding := &ClusterRoleBinding{
		Object: crb,
	}
	err = clusterRoleBinding.createOrUpdate()
	if err != nil {
		return err
	}
	return nil
}
