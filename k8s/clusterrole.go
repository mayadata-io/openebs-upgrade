package k8s

import (
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterRole is a wrapper over k8s ClusterRole
type ClusterRole struct {
	Object *rbacv1beta1.ClusterRole `json:"clusterRole"`
}

// createOrUpdate checks if the resource provided is present or not, if
// not present then it creates the resource otherwise updates it.
func (cr *ClusterRole) createOrUpdate() error {
	existingCr, err := Clientset.RbacV1beta1().ClusterRoles().Get(cr.Object.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			_, err = Clientset.RbacV1beta1().ClusterRoles().Create(cr.Object)
			if err != nil {
				return errors.Errorf("Error creating cluster role: %s: %+v", cr.Object.Name, err)
			}
		} else {
			return errors.Errorf("Error getting cluster role: %s: %+v", cr.Object.Name, err)
		}

	}
	// Set the resource version of the object to be updated
	cr.Object.SetResourceVersion(existingCr.GetResourceVersion())
	_, err = Clientset.RbacV1beta1().ClusterRoles().Update(cr.Object)
	if err != nil {
		return errors.Errorf("Error updating cluster role: %s: %+v", cr.Object.Name, err)
	}
	return nil
}

// DeployClusterRole creates/updates a given cluster role based on
// the given YAML.
func DeployClusterRole(YAML string) error {
	cr := &rbacv1beta1.ClusterRole{}
	err := yaml.Unmarshal([]byte(YAML), cr)
	if err != nil {
		return errors.Errorf(
			"Error unmarshalling clusterRole YAML: %+v", err)
	}
	clusterRole := &ClusterRole{
		Object: cr,
	}
	err = clusterRole.createOrUpdate()
	if err != nil {
		return err
	}
	return nil
}
