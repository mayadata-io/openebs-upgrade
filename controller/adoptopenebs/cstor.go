package adoptopenebs

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/types"
	"mayadata.io/openebs-upgrade/unstruct"
	"mayadata.io/openebs-upgrade/util"
)

// formCSPCOperatorConfig forms the desired OpenEBS CR config for CSPC operator.
func (p *Planner) formCSPCOperatorConfig(cspcOperator *unstructured.Unstructured) error {
	// CSPCOperator config is part of CStor config.
	cstorConfig := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	if p.CstorConfig != nil {
		cstorConfig = p.CstorConfig
	}
	// check if CSPCOperatorDetails are already filled, if yes then fill in
	// other details for cspc-operator.
	// Note: We are not handling error since we are interested in getting the value only
	// if it exists otherwise nil will be passed to the getResourceCommonDetails fn.
	cspcOperatorExistingDetails, _, _ := unstructured.NestedMap(cstorConfig.Object, "cspcOperator")
	// cspcOperatorDetails will store the details for cspcOperator
	cspcOperatorDetails, err := p.getResourceCommonDetails(cspcOperator, cspcOperatorExistingDetails)
	if err != nil {
		return err
	}
	// get the containers and required info of cspc operator
	containers, err := unstruct.GetNestedSliceOrError(cspcOperator, "spec",
		"template", "spec", "containers")
	if err != nil {
		return err
	}
	// // get the CSPI mgmt image tag
	getCSPCOperatorENVs := func(env *unstructured.Unstructured) error {
		envName, _, err := unstructured.NestedString(env.Object, "spec", "name")
		if err != nil {
			return err
		}
		if envName == "OPENEBS_IO_CSPI_MGMT_IMAGE" {
			cspiMgmtImage, _, err := unstructured.NestedString(env.Object, "spec", "value")
			if err != nil {
				return err
			}
			p.CSPIMgmtImageTag, err = util.GetImageTagFromContainerImage(cspiMgmtImage)
			if err != nil {
				return err
			}
		}
		return nil
	}
	getContainerDetails := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}
		envs, _, err := unstruct.GetSlice(obj, "spec", "env")
		if err != nil {
			return err
		}
		if containerName == types.CSPCOperatorContainerKey {
			cspcOperatorDetails[types.KeyResources], _, err = unstructured.NestedMap(obj.Object,
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
				cspcOperatorDetails[types.KeyImageTag] = imageTag
			}
			// get the environmets of the container.
			err = unstruct.SliceIterator(envs).ForEach(getCSPCOperatorENVs)
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
	cstorConfig.Object["cspcOperator"] = cspcOperatorDetails
	p.CstorConfig = cstorConfig

	return nil
}

// formCVCOperatorConfig forms the desired OpenEBS CR config for CVC operator.
func (p *Planner) formCVCOperatorConfig(cvcOperator *unstructured.Unstructured) error {
	// CVCOperator config is part of CStor config.
	cstorConfig := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	if p.CstorConfig != nil {
		cstorConfig = p.CstorConfig
	}
	// cvcOperatorDetails will store the details for cvcOperator
	cvcOperatorExistingDetails, _, _ := unstructured.NestedMap(cstorConfig.Object, "cvcOperator")
	cvcOperatorDetails, err := p.getResourceCommonDetails(cvcOperator, cvcOperatorExistingDetails)
	if err != nil {
		return err
	}
	// get the containers and required info of cvc operator
	containers, err := unstruct.GetNestedSliceOrError(cvcOperator, "spec",
		"template", "spec", "containers")
	if err != nil {
		return err
	}
	// get the values of required CVC operator environment variables
	getCVCOperatorENVs := func(env *unstructured.Unstructured) error {
		envName, _, err := unstructured.NestedString(env.Object, "spec", "name")
		if err != nil {
			return err
		}
		envValue, _, err := unstructured.NestedString(env.Object, "spec", "value")
		if err != nil {
			return err
		}
		if envName == "OPENEBS_IO_CSTOR_VOLUME_MGMT_IMAGE" {
			p.CStorVolumeManagerImageTag, err = util.GetImageTagFromContainerImage(envValue)
		}
		if err != nil {
			return err
		}
		return nil
	}
	getContainerDetails := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}
		envs, _, err := unstruct.GetSlice(obj, "spec", "env")
		if err != nil {
			return err
		}
		if containerName == types.CVCOperatorContainerKey {
			cvcOperatorDetails[types.KeyResources], _, err = unstructured.NestedMap(obj.Object,
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
				cvcOperatorDetails[types.KeyImageTag] = imageTag
			}
			// get the environments of the container.
			err = unstruct.SliceIterator(envs).ForEach(getCVCOperatorENVs)
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
	cstorConfig.Object["cvcOperator"] = cvcOperatorDetails
	p.CstorConfig = cstorConfig

	return nil
}

func (p *Planner) formCVCOperatorServiceConfig(cvcOperatorSVC *unstructured.Unstructured) error {
	// CVCOperatorService config is part of CVC operator config which is part of CStor config.
	cstorConfig := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	cvcOperatorConfig := make(map[string]interface{}, 0)
	if p.CstorConfig != nil {
		cstorConfig = p.CstorConfig
		cvcOperatorExistingConfig, _, _ := unstructured.NestedMap(cstorConfig.Object, "cvcOperator")
		if cvcOperatorExistingConfig != nil {
			cvcOperatorConfig = cvcOperatorExistingConfig
		}
	}
	cvcOperatorConfig["service"] = map[string]interface{}{
		"name": cvcOperatorSVC.GetName(),
	}
	cstorConfig.Object["cvcOperator"] = cvcOperatorConfig
	p.CstorConfig = cstorConfig

	return nil
}

// formCstorAdmissionServer forms the desired OpenEBS CR config for CStor admission server.
func (p *Planner) formCstorAdmissionServer(admissionServer *unstructured.Unstructured) error {
	// CStor admission server config is part of CStor config.
	cstorConfig := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	if p.CstorConfig != nil {
		cstorConfig = p.CstorConfig
	}
	// admissionServerDetails will store the details for admissionServer
	admissionServerDetails, err := p.getResourceCommonDetails(admissionServer, nil)
	if err != nil {
		return err
	}
	// get the containers and required info of admissionServer
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
		if containerName == "admission-webhook" {
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
		}
		return nil
	}
	err = unstruct.SliceIterator(containers).ForEach(getContainerDetails)
	if err != nil {
		return err
	}
	cstorConfig.Object["admissionServer"] = admissionServerDetails
	p.CstorConfig = cstorConfig

	return nil
}

// formCstorCSIControllerConfig forms the desired OpenEBS CR config for CStor CSI controller.
func (p *Planner) formCstorCSIController(cstorCSIController *unstructured.Unstructured) error {
	// CstorCSIController config is part of CStor config.
	cstorConfig := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	csiConfig := make(map[string]interface{}, 0)
	if p.CstorConfig != nil {
		cstorConfig = p.CstorConfig
		csi, exist, err := unstructured.NestedMap(cstorConfig.Object, "csi")
		if err != nil {
			return errors.Errorf(
				"Error forming CStorCSIController: error fetching CSI config from cstorConfig: %+v", err)
		}
		if exist {
			csiConfig = csi
		}
	}
	// csiControllerDetails will store the details for CSI controller component
	csiControllerDetails, err := p.getResourceCommonDetails(cstorCSIController, nil)
	if err != nil {
		return err
	}
	// get the containers and required info of CStor CSI controller
	containers, err := unstruct.GetNestedSliceOrError(cstorCSIController, "spec",
		"template", "spec", "containers")
	if err != nil {
		return err
	}
	getContainerDetails := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}
		if containerName == types.CStorCSIControllerCSIPluginContainerKey {
			csiControllerDetails[types.KeyResources], _, err = unstructured.NestedMap(obj.Object,
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
				csiControllerDetails[types.KeyImageTag] = imageTag
			}
		}
		return nil
	}
	err = unstruct.SliceIterator(containers).ForEach(getContainerDetails)
	if err != nil {
		return err
	}
	csiConfig["csiController"] = csiControllerDetails
	cstorConfig.Object["csi"] = csiConfig
	p.CstorConfig = cstorConfig

	return nil
}

// formCstorCSINodeConfig forms the desired OpenEBS CR config for CStor CSI node.
func (p *Planner) formCstorCSINode(cstorCSINode *unstructured.Unstructured) error {
	var (
		err error
	)
	// CstorCSINode config is part of CStor config.
	cstorConfig := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	csiConfig := make(map[string]interface{}, 0)
	if p.CstorConfig != nil {
		cstorConfig = p.CstorConfig
		csi, exist, err := unstructured.NestedMap(cstorConfig.Object, "csi")
		if err != nil {
			return errors.Errorf(
				"Error forming CStorCSINode: error fetching CSI config from cstorConfig: %+v", err)
		}
		if exist {
			csiConfig = csi
		}
	}
	// csiNodeDetails will store the details for CSI node component
	csiNodeDetails, err := p.getResourceCommonDetails(cstorCSINode, nil)
	if err != nil {
		return err
	}
	// getOpenEBSCSIPluginVolumeMountDetails fetches the volumeMounts path of openebs-csi-plugin container.
	getOpenEBSCSIPluginVolumeMountDetails := func(vm *unstructured.Unstructured) error {
		vmName, _, err := unstructured.NestedString(vm.Object, "spec", "name")
		if err != nil {
			return err
		}
		if vmName == "iscsiadm-bin" {
			vmPath, _, err := unstructured.NestedString(vm.Object, "spec", "mountPath")
			if err != nil {
				return err
			}
			csiNodeDetails[types.KeyISCSIPath] = vmPath
		}
		return nil
	}
	// get the containers and required info of CStor CSI node
	containers, err := unstruct.GetNestedSliceOrError(cstorCSINode, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}
	getContainerDetails := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}
		if containerName == types.CStorCSINodeCSIPluginContainerKey {
			csiNodeDetails[types.KeyResources], _, err = unstructured.NestedMap(obj.Object, "spec", "resources")
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
				csiNodeDetails[types.KeyImageTag] = imageTag
			}
			volumeMounts, _, err := unstruct.GetSlice(obj, "spec", "volumeMounts")
			if err != nil {
				return err
			}
			err = unstruct.SliceIterator(volumeMounts).ForEach(getOpenEBSCSIPluginVolumeMountDetails)
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
	csiConfig["csiNode"] = csiNodeDetails
	cstorConfig.Object["csi"] = csiConfig
	p.CstorConfig = cstorConfig

	return nil
}

// formCStorConfig forms the desired OpenEBS CR config for CStor.
func (p *Planner) formCStorConfig() error {
	cstorConfig := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	if p.CstorConfig != nil {
		cstorConfig = p.CstorConfig
	}
	poolConfig := make(map[string]interface{}, 0)
	poolMgmtConfig := make(map[string]interface{}, 0)
	targetConfig := make(map[string]interface{}, 0)
	cspiMgmtConfig := make(map[string]interface{}, 0)
	volumeMgmtConfig := make(map[string]interface{}, 0)
	volumeManagerConfig := make(map[string]interface{}, 0)

	if p.CStorPoolImageTag != p.OpenEBSVersion {
		poolConfig[types.KeyImageTag] = p.CStorPoolImageTag
	}
	if p.CStorPoolMgmtImageTag != p.OpenEBSVersion {
		poolMgmtConfig[types.KeyImageTag] = p.CStorPoolMgmtImageTag
	}
	if p.CStorTargetImageTag != p.OpenEBSVersion {
		targetConfig[types.KeyImageTag] = p.CStorTargetImageTag
	}
	if p.CSPIMgmtImageTag != p.OpenEBSVersion {
		cspiMgmtConfig[types.KeyImageTag] = p.CSPIMgmtImageTag
	}
	if p.CStorVolumeMgmtImageTag != p.OpenEBSVersion {
		volumeMgmtConfig[types.KeyImageTag] = p.CStorVolumeMgmtImageTag
	}
	if p.CStorVolumeManagerImageTag != p.OpenEBSVersion {
		volumeManagerConfig[types.KeyImageTag] = p.CStorVolumeManagerImageTag
	}
	cstorConfig.Object["pool"] = poolConfig
	cstorConfig.Object["poolMgmt"] = poolMgmtConfig
	cstorConfig.Object["target"] = targetConfig
	cstorConfig.Object["volumeMgmt"] = volumeMgmtConfig
	cstorConfig.Object["cspiMgmt"] = cspiMgmtConfig
	cstorConfig.Object["volumeManager"] = volumeManagerConfig
	p.CstorConfig = cstorConfig

	return nil
}

// formCStorCSIISCSIADMConfigmapConfig forms the desired OpenEBS CR config for openebs-cstor-csi-iscsiadm
// configmap.
func (p *Planner) formCStorCSIISCSIADMConfigmapConfig(iscsiadmConfigmap *unstructured.Unstructured) error {
	ISCSIADMConfigmapConfig := make(map[string]interface{}, 0)
	// CstorCSIISCSIADM config is part of CStor config.
	cstorConfig := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	csiConfig := make(map[string]interface{}, 0)
	if p.CstorConfig != nil {
		cstorConfig = p.CstorConfig
		csi, exist, err := unstructured.NestedMap(cstorConfig.Object, "csi")
		if err != nil {
			return errors.Errorf(
				"Error forming CStorCSIISCSIADMConfig: error fetching CSI config from cstorConfig: %+v", err)
		}
		if exist {
			csiConfig = csi
		}
	}
	// set the ISCSIADMConfigmapConfig values wrt ISCSIADMConfigmap field
	ISCSIADMConfigmapConfig[types.KeyName] = iscsiadmConfigmap.GetName()
	csiConfig["iscsiadmConfigmap"] = ISCSIADMConfigmapConfig
	cstorConfig.Object["csi"] = csiConfig
	p.CstorConfig = cstorConfig

	return nil
}
