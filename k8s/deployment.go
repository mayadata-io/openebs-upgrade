package k8s

import (
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Deployment is a wrapper over k8s Deployment
type Deployment struct {
	Object *appsv1.Deployment `json:"object"`
}

// createOrUpdate checks if the resource provided is present or not, if
// not present then it creates the resource otherwise updates it.
func (deploy *Deployment) createOrUpdate() error {
	existingDeploy, err := Clientset.AppsV1().Deployments(deploy.Object.Namespace).Get(deploy.Object.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			_, err = Clientset.AppsV1().Deployments(deploy.Object.Namespace).Create(deploy.Object)
			if err != nil {
				return errors.Errorf("Error creating deployment: %s: %+v", deploy.Object.Name, err)
			}
		} else {
			return errors.Errorf("Error getting deployment: %s: %+v", deploy.Object.Name, err)
		}
	}
	// Set the resource version of the object to be updated
	deploy.Object.SetResourceVersion(existingDeploy.GetResourceVersion())
	_, err = Clientset.AppsV1().Deployments(deploy.Object.Namespace).Update(deploy.Object)
	if err != nil {
		return errors.Errorf("Error updating deployment: %s: %+v", deploy.Object.Name, err)
	}
	return nil
}

// DeployDeployment creates/updates a given deployment based on
// the given YAML.
func DeployDeployment(YAML string) error {
	deploy := &appsv1.Deployment{}
	err := yaml.Unmarshal([]byte(YAML), deploy)
	if err != nil {
		return errors.Errorf(
			"Error unmarshalling deployment YAML: %+v", err)
	}
	deployment := &Deployment{
		Object: deploy,
	}
	err = deployment.createOrUpdate()
	if err != nil {
		return err
	}
	return nil
}
