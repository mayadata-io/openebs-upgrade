package adoptopenebs

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/types"
	"mayadata.io/openebs-upgrade/unstruct"
	"mayadata.io/openebs-upgrade/util"
	"strings"
)

// formAdmissionServerConfig forms the desired OpenEBS CR config for admission-server.
func (p *Planner) formAdmissionServerConfig(admissionServer *unstructured.Unstructured) error {
	admissionServerConfig := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	if p.AdmissionServerConfig != nil {
		admissionServerConfig = p.AdmissionServerConfig
	}
	// admissionServerDetails will store the details for admissionServer
	admissionServerDetails, err := p.getResourceCommonDetails(admissionServer, admissionServerConfig.Object)
	if err != nil {
		return err
	}
	// get the containers and required info of admission server
	containers, err := unstruct.GetNestedSliceOrError(admissionServer, "spec",
		"template", "spec", "containers")
	if err != nil {
		return err
	}
	getContainerDetails := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}
		if containerName == types.AdmissionServerContainerKey ||
			strings.Contains(containerName, types.AdmissionServerContainerKey) {
			admissionServerDetails[types.KeyResources], _, err = unstructured.NestedMap(obj.Object,
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
				admissionServerDetails[types.KeyImageTag] = imageTag
			}
			if containerName != types.AdmissionServerContainerKey {
				admissionServerDetails[types.KeyContainerName] = containerName
			}
		}
		return nil
	}
	err = unstruct.SliceIterator(containers).ForEach(getContainerDetails)
	if err != nil {
		return err
	}
	admissionServerConfig.Object = admissionServerDetails
	p.AdmissionServerConfig = admissionServerConfig

	return nil
}
