package adoptopenebs

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/types"
	"mayadata.io/openebs-upgrade/unstruct"
	"mayadata.io/openebs-upgrade/util"
)

// formMayastorMOACConfig forms the desired OpenEBS CR config for Mayastor's moac deployment.
func (p *Planner) formMayastorMOACConfig(moac *unstructured.Unstructured) error {
	mayastorConfig := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	moacConfig := make(map[string]interface{}, 0)
	if p.MayastorConfig != nil {
		mayastorConfig = p.MayastorConfig
		moacExistingConfig, _, _ := unstructured.NestedMap(mayastorConfig.Object, "moac")
		if moacExistingConfig != nil {
			moacConfig = moacExistingConfig
		}
	}
	// moacDeployDetails will store the details for Mayastor moac deployment.
	moacDeployDetails, err := p.getResourceCommonDetails(moac, moacConfig)
	if err != nil {
		return err
	}
	// get the containers and required info of Mayastor moac deployment.
	containers, err := unstruct.GetNestedSliceOrError(moac, "spec",
		"template", "spec", "containers")
	if err != nil {
		return err
	}
	getContainerDetails := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}
		if containerName == types.MoacContainerKey {
			moacDeployDetails[types.KeyResources], _, err = unstructured.NestedMap(obj.Object,
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
				moacDeployDetails[types.KeyImageTag] = imageTag
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
	mayastorConfig.Object["moac"] = moacDeployDetails
	p.MayastorConfig = mayastorConfig

	return nil
}

func (p *Planner) formMOACServiceConfig(moacSVC *unstructured.Unstructured) error {
	// moacService config is part of moac config which is part of mayastor config.
	mayastorConfig := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	moacConfig := make(map[string]interface{}, 0)
	if p.MayastorConfig != nil {
		mayastorConfig = p.MayastorConfig
		moacExistingConfig, _, _ := unstructured.NestedMap(mayastorConfig.Object, "moac")
		if moacExistingConfig != nil {
			moacConfig = moacExistingConfig
		}
	}
	moacConfig["service"] = map[string]interface{}{
		"name": moacSVC.GetName(),
	}
	mayastorConfig.Object["moac"] = moacConfig
	p.MayastorConfig = mayastorConfig

	return nil
}

// formMayastorDaemonConfig forms the desired OpenEBS CR config for Mayastor's mayastor daemonset.
func (p *Planner) formMayastorDaemonConfig(mayastorDaemon *unstructured.Unstructured) error {
	var (
		err error
	)
	// mayastorDaemon Config is part of mayastorConfig.
	mayastorConfig := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	mayastorDaemonConfig := make(map[string]interface{}, 0)
	if p.MayastorConfig != nil {
		mayastorConfig = p.MayastorConfig
		mayastorDaemonExistingConfig, _, _ := unstructured.NestedMap(mayastorConfig.Object, "mayastor")
		if mayastorDaemonExistingConfig != nil {
			mayastorDaemonConfig = mayastorDaemonExistingConfig
		}
	}
	// mayastorDaemonDetails will store the details for mayastor daemonset.
	mayastorDaemonDetails, err := p.getResourceCommonDetails(mayastorDaemon, mayastorDaemonConfig)
	if err != nil {
		return err
	}
	// get the containers and required info of mayastor daemonset.
	containers, err := unstruct.GetNestedSliceOrError(mayastorDaemon, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}
	getContainerDetails := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
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
		if containerName == types.MayastorContainerKey {
			mayastorContainerDetails := make(map[string]interface{}, 0)
			mayastorContainerDetails[types.KeyResources], _, err = unstructured.NestedMap(obj.Object, "spec", "resources")
			if err != nil {
				return err
			}
			if imageTag != p.OpenEBSVersion {
				mayastorContainerDetails[types.KeyImageTag] = imageTag
			}
			mayastorDaemonDetails["mayastor"] = mayastorContainerDetails
		} else if containerName == types.MayastorGRPCContainerKey {
			mayastorGRPCContainerDetails := make(map[string]interface{}, 0)
			mayastorGRPCContainerDetails[types.KeyResources], _, err = unstructured.NestedMap(obj.Object, "spec", "resources")
			if err != nil {
				return err
			}
			if imageTag != p.OpenEBSVersion {
				mayastorGRPCContainerDetails[types.KeyImageTag] = imageTag
			}
			mayastorDaemonDetails["mayastorGrpc"] = mayastorGRPCContainerDetails
		}
		return nil
	}
	err = unstruct.SliceIterator(containers).ForEach(getContainerDetails)
	if err != nil {
		return err
	}
	mayastorConfig.Object["mayastor"] = mayastorDaemonDetails
	p.MayastorConfig = mayastorConfig

	return nil
}
