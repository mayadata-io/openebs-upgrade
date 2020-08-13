package adoptopenebs

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/types"
	"mayadata.io/openebs-upgrade/unstruct"
	"mayadata.io/openebs-upgrade/util"
	"strconv"
)

// formMayaAPIServerConfig forms the desired OpenEBS CR config for MayaAPIServer.
func (p *Planner) formMayaAPIServerConfig(mayaAPIServer *unstructured.Unstructured) error {
	mayaAPIServerConfig := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	if p.APIServerConfig != nil {
		mayaAPIServerConfig = p.APIServerConfig
	}
	var (
		cstorSparsePoolEnabled     bool
		createDefaultStorageConfig bool
		enableAnalytics            bool
		jivaReplicaCount           int64
	)
	// mayaAPIServerDetails will store the details for mayaAPIServer
	mayaAPIServerDetails, err := p.getResourceCommonDetails(mayaAPIServer, mayaAPIServerConfig.Object)
	if err != nil {
		return err
	}
	// get the containers and required info of mayaAPIServer
	containers, err := unstruct.GetNestedSliceOrError(mayaAPIServer, "spec",
		"template", "spec", "containers")
	if err != nil {
		return err
	}
	getMayaAPIServerENVs := func(env *unstructured.Unstructured) error {
		envName, _, err := unstructured.NestedString(env.Object, "spec", "name")
		if err != nil {
			return err
		}
		envValue, _, err := unstructured.NestedString(env.Object, "spec", "value")
		if err != nil {
			return err
		}
		// get the value of cstorSparsePool enabled or not in order to fill the
		// OpenEBS CR.
		if envName == "OPENEBS_IO_INSTALL_DEFAULT_CSTOR_SPARSE_POOL" {
			if len(envValue) > 0 {
				cstorSparsePoolEnabled, err = strconv.ParseBool(envValue)
			}
		} else if envName == "OPENEBS_IO_CREATE_DEFAULT_STORAGE_CONFIG" {
			if len(envValue) > 0 {
				createDefaultStorageConfig, err = strconv.ParseBool(envValue)
			}
			p.CreateDefaultStorageConfig = createDefaultStorageConfig
		} else if envName == "OPENEBS_IO_JIVA_CONTROLLER_IMAGE" {
			p.JivaCtrlImageTag, err = util.GetImageTagFromContainerImage(envValue)
		} else if envName == "OPENEBS_IO_JIVA_REPLICA_IMAGE" {
			p.JivaReplicaImageTag, err = util.GetImageTagFromContainerImage(envValue)
		} else if envName == "OPENEBS_IO_JIVA_REPLICA_COUNT" {
			if len(envValue) > 0 {
				jivaReplicaCount, err = strconv.ParseInt(envValue, 10, 32)
			}
			p.JivaReplicaCount = jivaReplicaCount
		} else if envName == "OPENEBS_IO_CSTOR_TARGET_IMAGE" {
			p.CStorTargetImageTag, err = util.GetImageTagFromContainerImage(envValue)
		} else if envName == "OPENEBS_IO_CSTOR_POOL_IMAGE" {
			p.CStorPoolImageTag, err = util.GetImageTagFromContainerImage(envValue)
		} else if envName == "OPENEBS_IO_CSTOR_POOL_MGMT_IMAGE" {
			p.CStorPoolMgmtImageTag, err = util.GetImageTagFromContainerImage(envValue)
		} else if envName == "OPENEBS_IO_CSTOR_VOLUME_MGMT_IMAGE" {
			p.CStorVolumeMgmtImageTag, err = util.GetImageTagFromContainerImage(envValue)
		} else if envName == "OPENEBS_IO_VOLUME_MONITOR_IMAGE" {
			p.VolumeMonitorImageTag, err = util.GetImageTagFromContainerImage(envValue)
		} else if envName == "OPENEBS_IO_CSTOR_POOL_EXPORTER_IMAGE" {
			p.PoolExporterImageTag, err = util.GetImageTagFromContainerImage(envValue)
		} else if envName == "OPENEBS_IO_HELPER_IMAGE" {
			p.HelperImageTag, err = util.GetImageTagFromContainerImage(envValue)
		} else if envName == "OPENEBS_IO_ENABLE_ANALYTICS" {
			if len(envValue) > 0 {
				enableAnalytics, err = strconv.ParseBool(envValue)
			}
			p.EnableAnalytics = enableAnalytics
		} else if envName == "OPENEBS_IO_BASE_DIR" {
			// set the default storage path to the value of this env variable if provided.
			p.DefaultStoragePath = envValue
		}
		if err != nil {
			return errors.Errorf("Error forming mayaAPIServerConfig: %+v", err)
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
		if containerName == "maya-apiserver" {
			mayaAPIServerDetails[types.KeyResources], _, err = unstructured.NestedMap(obj.Object,
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
				mayaAPIServerDetails[types.KeyImageTag] = imageTag
			}
			// get the environmets of the container.
			err = unstruct.SliceIterator(envs).ForEach(getMayaAPIServerENVs)
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
	mayaAPIServerDetails["cstorSparsePool"] = map[string]interface{}{
		"enabled": cstorSparsePoolEnabled,
	}
	mayaAPIServerConfig.Object = mayaAPIServerDetails
	p.APIServerConfig = mayaAPIServerConfig
	return nil
}

func (p *Planner) formMayaAPIServerServiceConfig(apiService *unstructured.Unstructured) error {
	// mayaAPIServiceConfig is a part of mayaAPIServer config only in OpenEBS CR.
	mayaAPIServiceConfig := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	if p.APIServerConfig != nil {
		mayaAPIServiceConfig = p.APIServerConfig
	}
	mayaAPIServiceConfig.Object["service"] = map[string]interface{}{
		"name": apiService.GetName(),
	}
	p.APIServerConfig = mayaAPIServiceConfig

	return nil
}
