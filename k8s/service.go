package k8s

import (
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Service is a wrapper over k8s Service
type Service struct {
	Object *corev1.Service `json:"object"`
}

// createOrUpdate checks if the resource provided is present or not, if
// not present then it creates the resource otherwise updates it.
func (svc *Service) createOrUpdate() error {
	existingSvc, err := Clientset.CoreV1().Services(svc.Object.Namespace).Get(svc.Object.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			_, err = Clientset.CoreV1().Services(svc.Object.Namespace).Create(svc.Object)
			if err != nil {
				return errors.Errorf("Error creating service: %s: %+v", svc.Object.Name, err)
			}
		} else {
			return errors.Errorf("Error getting service: %s: %+v", svc.Object.Name, err)
		}

	}
	// TODO: Handle the case more elegantly where Service Type/ClusterIP is still not
	// assigned to service and it tries to update the service with empty ClusterIP and
	// it throws an error for that particular reconciliation.
	if existingSvc.Spec.Type == "" {
		return nil
	}
	// Update the cluster IP of the existing svc into this svc structure
	if existingSvc.Spec.Type == corev1.ServiceTypeClusterIP {
		svc.Object.Spec.ClusterIP = existingSvc.Spec.ClusterIP
	}
	svc.Object.SetResourceVersion(existingSvc.GetResourceVersion())
	_, err = Clientset.CoreV1().Services(svc.Object.Namespace).Update(svc.Object)
	if err != nil {
		return errors.Errorf("Error updating service: %s: %+v", svc.Object.Name, err)
	}
	return nil
}

// DeployService creates/updates a given service based on
// the given YAML.
func DeployService(YAML string) error {
	svc := &corev1.Service{}
	err := yaml.Unmarshal([]byte(YAML), svc)
	if err != nil {
		return errors.Errorf(
			"Error unmarshalling service YAML: %+v", err)
	}
	service := &Service{
		Object: svc,
	}
	err = service.createOrUpdate()
	if err != nil {
		return err
	}
	return nil
}
