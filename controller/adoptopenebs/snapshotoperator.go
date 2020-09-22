package adoptopenebs

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/types"
	"mayadata.io/openebs-upgrade/unstruct"
	"mayadata.io/openebs-upgrade/util"
	"strings"
)

// formSnapshotOperatorConfig forms the desired OpenEBS CR config for Snapshot operator.
func (p *Planner) formSnapshotOperatorConfig(snapshotOperator *unstructured.Unstructured) error {
	snapshotOperatorConfig := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	if p.SnapshotOperatorConfig != nil {
		snapshotOperatorConfig = p.SnapshotOperatorConfig
	}
	var (
		controllerImageTag       string
		provisionerImageTag      string
		controllerContainerName  string
		provisionerContainerName string
	)
	// snapshotOperatorDetails will store the details for snapshotOperator
	snapshotOperatorDetails, err := p.getResourceCommonDetails(snapshotOperator, nil)
	if err != nil {
		return err
	}
	// get the containers and required info of snapshot operator
	containers, err := unstruct.GetNestedSliceOrError(snapshotOperator, "spec",
		"template", "spec", "containers")
	if err != nil {
		return err
	}
	getContainerDetails := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}
		if containerName == types.SnapshotControllerContainerKey ||
			strings.Contains(containerName, "snapshot-controller") {
			snapshotOperatorDetails[types.KeyResources], _, err = unstructured.NestedMap(obj.Object,
				"spec", "resources")
			if err != nil {
				return err
			}
			image, _, err := unstructured.NestedString(obj.Object, "spec", "image")
			if err != nil {
				return err
			}
			// fill the imageTag only if it is different from OpenEBS version
			imageTag, err := util.GetImageTagFromContainerImage(image)
			if err != nil {
				return err
			}
			if imageTag != p.OpenEBSVersion {
				controllerImageTag = imageTag
			}
			if containerName != types.SnapshotControllerContainerKey {
				controllerContainerName = containerName
			}
		}
		if containerName == types.SnapshotProvisionerContainerKey ||
			strings.Contains(containerName, "snapshot-provisioner") {
			image, _, err := unstructured.NestedString(obj.Object, "spec", "image")
			if err != nil {
				return err
			}
			// fill the imageTag only if it is different from OpenEBS version
			imageTag, err := util.GetImageTagFromContainerImage(image)
			if err != nil {
				return err
			}
			if imageTag != p.OpenEBSVersion {
				provisionerImageTag = imageTag
			}
			if containerName != types.SnapshotProvisionerContainerKey {
				provisionerContainerName = containerName
			}
		}
		return nil
	}
	err = unstruct.SliceIterator(containers).ForEach(getContainerDetails)
	if err != nil {
		return err
	}
	snapshotOperatorConfig.Object = snapshotOperatorDetails
	snapshotOperatorConfig.Object["controller"] = map[string]interface{}{
		types.KeyImageTag:      controllerImageTag,
		types.KeyContainerName: controllerContainerName,
	}
	snapshotOperatorConfig.Object["provisioner"] = map[string]interface{}{
		types.KeyImageTag:      provisionerImageTag,
		types.KeyContainerName: provisionerContainerName,
	}
	p.SnapshotOperatorConfig = snapshotOperatorConfig

	return nil
}
