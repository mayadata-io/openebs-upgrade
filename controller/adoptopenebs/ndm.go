package adoptopenebs

import (
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/controller/openebs"
	"mayadata.io/openebs-upgrade/types"
	"mayadata.io/openebs-upgrade/unstruct"
	"mayadata.io/openebs-upgrade/util"
	"strconv"
)

// formNDMOperatorConfig forms the desired OpenEBS CR config for NDM operator.
func (p *Planner) formNDMOperatorConfig(ndmOperator *unstructured.Unstructured) error {
	ndmOperatorConfig := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	// ndmOperatorDetails will store the details for ndmOperator
	ndmOperatorDetails, err := p.getResourceCommonDetails(ndmOperator)
	if err != nil {
		return err
	}
	// get the containers and required info of ndmOperator
	containers, err := unstruct.GetNestedSliceOrError(ndmOperator, "spec",
		"template", "spec", "containers")
	if err != nil {
		return err
	}
	getContainerDetails := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}
		if containerName == "node-disk-operator" {
			ndmOperatorDetails[types.KeyResources], _, err = unstructured.NestedMap(obj.Object,
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
			if imageTag != openebs.SupportedNDMVersionForOpenEBSVersion[p.OpenEBSVersion] {
				ndmOperatorDetails[types.KeyImageTag] = imageTag
			}
		}
		return nil
	}
	err = unstruct.SliceIterator(containers).ForEach(getContainerDetails)
	if err != nil {
		return err
	}
	ndmOperatorConfig.Object = ndmOperatorDetails
	p.NDMOperatorConfig = ndmOperatorConfig

	return nil
}

// formNDMDaemonConfig forms the desired OpenEBS CR config for NDM daemon.
func (p *Planner) formNDMDaemonConfig(ndm *unstructured.Unstructured) error {
	ndmConfig := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	if p.NDMDaemonConfig != nil {
		ndmConfig = p.NDMDaemonConfig
	}
	var (
		// NDM sparse related details
		sparsePath  string
		sparseCount string
		sparseSize  string
	)
	// ndmDaemonDetails will store the details for ndmDaemon
	ndmDaemonDetails, err := p.getResourceCommonDetails(ndm)
	if err != nil {
		return err
	}
	// get the containers and required info of ndmDaemon
	containers, err := unstruct.GetNestedSliceOrError(ndm, "spec",
		"template", "spec", "containers")
	if err != nil {
		return err
	}
	getNDMDaemonENVs := func(env *unstructured.Unstructured) error {
		envName, _, err := unstructured.NestedString(env.Object, "spec", "name")
		if err != nil {
			return err
		}
		envValue, _, err := unstructured.NestedString(env.Object, "spec", "value")
		if err != nil {
			return err
		}
		if envName == types.SparseFileDirectoryEnv {
			sparsePath = envValue
		} else if envName == types.SparseFileSizeEnv {
			sparseSize = envValue
		} else if envName == types.SparseFileCountEnv {
			sparseCount = envValue
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
		if containerName == types.NDMDaemonContainerKey {
			ndmDaemonDetails[types.KeyResources], _, err = unstructured.NestedMap(obj.Object,
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
			if imageTag != openebs.SupportedNDMVersionForOpenEBSVersion[p.OpenEBSVersion] {
				ndmDaemonDetails[types.KeyImageTag] = imageTag
			}
			// get the environments of the container.
			err = unstruct.SliceIterator(envs).ForEach(getNDMDaemonENVs)
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
	ndmDaemonDetails["sparse"] = map[string]interface{}{
		"path":  sparsePath,
		"size":  sparseSize,
		"count": sparseCount,
	}
	// update ndmConfig with the formed details as NDMDaemonConfig can be updated
	// by multiple functions.
	for key, value := range ndmDaemonDetails {
		ndmConfig.Object[key] = value
	}
	p.NDMDaemonConfig = ndmConfig

	return nil
}

// formNDMConfigMapConfig forms the desired OpenEBS CR config for NDM configmap.
func (p *Planner) formNDMConfigMapConfig(ndmCm *unstructured.Unstructured) error {
	ndmConfigMapConfig := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	ndmConfigMap := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	if p.NDMDaemonConfig != nil {
		ndmConfigMapConfig = p.NDMDaemonConfig
	}

	var (
		udevProbeEnabled     bool
		smartProbeEnabled    bool
		seachestProbeEnabled bool
		OSDiskFilterEnabled  bool
		OSDiskFilterExclude  string
		vendorFilterEnabled  bool
		vendorFilterExclude  string
		vendorFilterInclude  string
		pathFilterEnabled    bool
		pathFilterExclude    string
		pathFilterInclude    string
	)
	// Initialize NDM config data structure i.e., data field of the configmap
	// in order to get the data from the configmap.
	ndmConfigData := &types.NDMConfig{}
	// get the configmap template which we will use as a structure to fill
	// in the given/default values.
	dataMap, _, err := unstructured.NestedMap(ndmCm.Object, "data")
	if err != nil {
		return err
	}
	ndmConfigDataTemplate := dataMap["node-disk-manager.config"]
	err = yaml.Unmarshal([]byte(ndmConfigDataTemplate.(string)), ndmConfigData)
	if err != nil {
		return errors.Errorf("Error unmarshalling NDM config data: %+v, Error: %+v", ndmConfigDataTemplate, err)
	}
	// get the value of probes if they are enabled or disabled.
	for _, probeConfig := range ndmConfigData.ProbeConfigs {
		// if the value is empty then do not parse it.
		if len(probeConfig.State) <= 0 {
			continue
		}
		if probeConfig.Key == types.UdevProbeKey {
			udevProbeEnabled, err = strconv.ParseBool(probeConfig.State)
		} else if probeConfig.Key == types.SmartProbeKey {
			smartProbeEnabled, err = strconv.ParseBool(probeConfig.State)
		} else if probeConfig.Key == types.SeachestProbeKey {
			seachestProbeEnabled, err = strconv.ParseBool(probeConfig.State)
		}
		if err != nil {
			return err
		}
	}
	// Get the data of NDM filters.
	for _, filterConfig := range ndmConfigData.FilterConfigs {
		if len(filterConfig.State) <= 0 {
			continue
		}
		if filterConfig.Key == types.OSDiskFilterKey {
			OSDiskFilterEnabled, err = strconv.ParseBool(filterConfig.State)
			OSDiskFilterExclude = *filterConfig.Exclude

		} else if filterConfig.Key == types.VendorFilterKey {
			vendorFilterEnabled, err = strconv.ParseBool(filterConfig.State)
			vendorFilterExclude = *filterConfig.Exclude
			vendorFilterInclude = *filterConfig.Include
		} else if filterConfig.Key == types.PathFilterKey {
			pathFilterEnabled, err = strconv.ParseBool(filterConfig.State)
			pathFilterExclude = *filterConfig.Exclude
			pathFilterInclude = *filterConfig.Include
		}
		if err != nil {
			return err
		}
	}

	// form the NDMConfigMap config which will be filled in OpenEBS CR.
	ndmConfigMapConfig.Object["filters"] = map[string]interface{}{
		"osDisk": map[string]interface{}{
			types.KeyEnabled: OSDiskFilterEnabled,
			"exclude":        OSDiskFilterExclude,
		},
		"vendor": map[string]interface{}{
			types.KeyEnabled: vendorFilterEnabled,
			"exclude":        vendorFilterExclude,
			"include":        vendorFilterInclude,
		},
		"path": map[string]interface{}{
			types.KeyEnabled: pathFilterEnabled,
			"exclude":        pathFilterExclude,
			"include":        pathFilterInclude,
		},
	}
	ndmConfigMapConfig.Object["probes"] = map[string]interface{}{
		"udev": map[string]interface{}{
			types.KeyEnabled: udevProbeEnabled,
		},
		"smart": map[string]interface{}{
			types.KeyEnabled: smartProbeEnabled,
		},
		"seachest": map[string]interface{}{
			types.KeyEnabled: seachestProbeEnabled,
		},
	}
	// set the configMap related details to ndmDaemon field
	p.NDMDaemonConfig = ndmConfigMapConfig
	// set the ndmConfigMap values wrt ndmConfigMap field also
	ndmConfigMap.SetName(ndmCm.GetName())
	p.NDMConfigMapConfig = ndmConfigMap

	return nil
}
