package k8s

import (
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"strings"
)

// GetNodes returns the list of nodes.
func GetNodes() (*v1.NodeList, error) {
	nodes, err := Clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return &v1.NodeList{}, err
	}

	return nodes, nil
}

// GetOSImage return the value of the OS Image of a Node.
func GetOSImage() (string, error) {
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

// GetUbuntuVersion returns the ubuntu version of a Node.
// For example: 18.04, 16.04
func GetUbuntuVersion() (float64, error) {
	var version float64

	osImage, err := GetOSImage()
	if err != nil {
		return version, errors.Errorf("Error getting OS Image. Error: %v", err)
	}

	if !strings.Contains(strings.ToLower(osImage), strings.ToLower("Ubuntu")) {
		return version, nil
	}

	versionString := strings.Split(osImage, " ")[1]
	// Take the version upto first decimal.
	versionString = strings.Join(strings.Split(versionString, ".")[0:2], ".")

	version, err = strconv.ParseFloat(versionString, 64)
	if err != nil {
		return version, errors.Errorf("Error parsing string to float. Error: %v", err)
	}

	return version, nil
}

// GetK8sVersion returns the k8s version.
func GetK8sVersion() (string, error) {
	versionInfo, err := Clientset.Discovery().ServerVersion()
	if err != nil {
		return "", err
	}

	return versionInfo.GitVersion, nil
}
