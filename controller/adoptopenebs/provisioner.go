package adoptopenebs

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/types"
	"mayadata.io/openebs-upgrade/unstruct"
	"mayadata.io/openebs-upgrade/util"
)

// formOpenEBSProvisionerConfig forms the desired OpenEBS CR config for openebs-provisioner.
func (p *Planner) formOpenEBSProvisionerConfig(provisioner *unstructured.Unstructured) error {
	provisionerConfig := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	// provisionerDetails will store the details for OpenEBS provisioner
	provisionerDetails, err := p.getResourceCommonDetails(provisioner)
	if err != nil {
		return err
	}
	// get the containers and required info of OpenEBS Provisioner
	containers, err := unstruct.GetNestedSliceOrError(provisioner, "spec",
		"template", "spec", "containers")
	if err != nil {
		return err
	}
	getContainerDetails := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}
		if containerName == types.OpenEBSProvisionerContainerKey {
			provisionerDetails[types.KeyResources], _, err = unstructured.NestedMap(obj.Object,
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
				provisionerDetails[types.KeyImageTag] = imageTag
			}
			// set the imagePrefix or registry for the container images
			//
			// NOTE: It is assumed that the registry being used is same for the containers
			// across the components.
			p.ImagePrefix, err = util.GetImagePrefixFromContainerImage(image)
			if err != nil {
				return err
			}
		}
		return nil
	}
	err = unstruct.SliceIterator(containers).ForEach(getContainerDetails)
	if err != nil {
		return err
	}
	provisionerConfig.Object = provisionerDetails
	p.ProvisionerConfig = provisionerConfig

	return nil
}
