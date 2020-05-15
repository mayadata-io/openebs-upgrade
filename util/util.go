package util

import (
	"github.com/pkg/errors"
	"strings"
)

// GetImageTagFromContainerImage returns the image tag for a given container image.
func GetImageTagFromContainerImage(image string) (string, error) {
	lastIndex := strings.LastIndex(image, ":")
	if lastIndex == -1 {
		return "", errors.Errorf("no version tag found on image %s", image)
	}
	return image[lastIndex+1:], nil
}

// GetImagePrefixFromContainerImage returns the image prefix or registry for a given container image.
func GetImagePrefixFromContainerImage(image string) (string, error) {
	lastIndex := strings.LastIndex(image, ":")
	if lastIndex == -1 {
		return "", errors.Errorf("no version tag found on image %s", image)
	}
	imageWithoutTag := image[:lastIndex]
	lastIndexOfForwardSlash := strings.LastIndex(imageWithoutTag, "/")
	if lastIndexOfForwardSlash == -1 {
		return "", errors.Errorf("no version tag found on image %s", image)
	}
	return imageWithoutTag[:lastIndexOfForwardSlash+1], nil
}
