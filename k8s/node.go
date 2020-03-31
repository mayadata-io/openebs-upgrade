package k8s

import (
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetNodes returns the list of nodes.
func GetNodes() (*v1.NodeList, error) {
	nodes, err := Clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return &v1.NodeList{}, err
	}

	return nodes, nil
}

// GetOSImageOfNode -
func GetOSImageOfNode() (string, error) {
	nodes, err := GetNodes()
	if err != nil {
		return "", err
	}

	if len(nodes.Items) == 0 {
		return "", errors.Errorf("No nodes found.")
	}

	// Take the first node from the list and get OS Image.
	node := nodes.Items[0]
	osImage := node.Status.NodeInfo.OSImage

	return osImage, nil
}
