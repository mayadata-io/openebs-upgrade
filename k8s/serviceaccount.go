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
