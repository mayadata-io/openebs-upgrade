package adoptopenebs

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/types"
	"mayadata.io/openebs-upgrade/unstruct"
	"mayadata.io/openebs-upgrade/util"
	"strings"
)

// formLocalProvisionerConfig forms the desired OpenEBS CR config for localpv-provisioner.
func (p *Planner) formLocalProvisionerConfig(localProvisioner *unstructured.Unstructured) error {
	localProvisionerConfig := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	// localProvisionerDetails will store the details for localProvisioner
	localProvisionerDetails, err := p.getResourceCommonDetails(localProvisioner, nil)
	if err != nil {
		return err
	}
	// get the containers and required info of localProvisioner
	containers, err := unstruct.GetNestedSliceOrError(localProvisioner, "spec",
		"template", "spec", "containers")
	if err != nil {
		return err
	}
	getContainerDetails := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}
		// contains is mainly being used to cover helm, since for helm, some of the container
		// names are also different.
		if containerName == types.LocalPVProvisionerContainerKey ||
			strings.Contains(containerName, "localpv") {
			localProvisionerDetails[types.KeyResources], _, err = unstructured.NestedMap(obj.Object,
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
				localProvisionerDetails[types.KeyImageTag] = imageTag
			}
			// update the container name if different from what this controller wants i.e.,
			// openebs-provisioner-hostpath
			if containerName != types.LocalPVProvisionerContainerKey {
				localProvisionerDetails[types.KeyContainerName] = containerName
			}
		}
		return nil
	}
	err = unstruct.SliceIterator(containers).ForEach(getContainerDetails)
	if err != nil {
		return err
	}
	localProvisionerConfig.Object = localProvisionerDetails
	p.LocalProvisionerConfig = localProvisionerConfig

	return nil
}
