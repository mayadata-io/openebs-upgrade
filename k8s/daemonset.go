package k8s

import (
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DaemonSet is a wrapper over k8s DaemonSet
type DaemonSet struct {
	Object *appsv1.DaemonSet `json:"object"`
}

// createOrUpdate checks if the resource provided is present or not, if
// not present then it creates the resource otherwise updates it.
func (ds *DaemonSet) createOrUpdate() error {
	existingDs, err := Clientset.AppsV1().DaemonSets(ds.Object.Namespace).Get(ds.Object.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			_, err = Clientset.AppsV1().DaemonSets(ds.Object.Namespace).Create(ds.Object)
			if err != nil {
				return errors.Errorf("Error creating daemonset: %s: %+v", ds.Object.Name, err)
			}
		} else {
			return errors.Errorf("Error getting daemonset: %s: %+v", ds.Object.Name, err)
		}
	}
	// Set the resource version of the object to be updated
	ds.Object.SetResourceVersion(existingDs.GetResourceVersion())
	_, err = Clientset.AppsV1().DaemonSets(ds.Object.Namespace).Update(ds.Object)
	if err != nil {
		return errors.Errorf("Error updating daemonset: %s: %+v", ds.Object.Name, err)
	}
	return nil
}

// DeployDaemonSet creates/updates a given daemonset based on
// the given YAML.
func DeployDaemonSet(YAML string) error {
	ds := &appsv1.DaemonSet{}
	err := yaml.Unmarshal([]byte(YAML), ds)
	if err != nil {
		return errors.Errorf(
			"Error unmarshalling daemonSet YAML: %+v", err)
	}
	daemonSet := &DaemonSet{
		Object: ds,
	}
	err = daemonSet.createOrUpdate()
	if err != nil {
		return err
	}
	return nil
}
