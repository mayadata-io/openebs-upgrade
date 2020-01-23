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
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceAccount is a wrapper over k8s ServiceAccount
type ServiceAccount struct {
	Object *corev1.ServiceAccount `json:"object"`
}

// createOrUpdate checks if the resource provided is present or not, if
// not present then it creates the resource otherwise updates it.
func (sa *ServiceAccount) createOrUpdate() error {
	existingSa, err := Clientset.CoreV1().ServiceAccounts(sa.Object.Namespace).Get(sa.Object.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			_, err = Clientset.CoreV1().ServiceAccounts(sa.Object.Namespace).Create(sa.Object)
			if err != nil {
				return errors.Errorf("Error creating service account: %s: %+v", sa.Object.Name, err)
			}
		} else {
			return errors.Errorf("Error getting service account: %s: %+v", sa.Object.Name, err)
		}
	}
	// Set the resource version of the object to be updated
	sa.Object.SetResourceVersion(existingSa.GetResourceVersion())
	_, err = Clientset.CoreV1().ServiceAccounts(sa.Object.Namespace).Update(sa.Object)
	if err != nil {
		return errors.Errorf("Error updating service account: %s: %+v", sa.Object.Name, err)
	}
	return nil
}

// DeployServiceAccount creates/updates a given serviceAccount based on
// the given YAML.
func DeployServiceAccount(YAML string) error {
	sa := &corev1.ServiceAccount{}
	err := yaml.Unmarshal([]byte(YAML), sa)
	if err != nil {
		return errors.Errorf(
			"Error unmarshalling serviceAccount YAML: %+v", err)
	}
	serviceAccount := &ServiceAccount{
		Object: sa,
	}
	err = serviceAccount.createOrUpdate()
	if err != nil {
		return err
	}
	return nil
}
