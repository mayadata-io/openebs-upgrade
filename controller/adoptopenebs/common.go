package adoptopenebs

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/types"
)

func (p *Planner) formComponentOpenEBSConfig(component *unstructured.Unstructured, componentType string) error {
	var err error

	switch componentType {
	case types.MayaAPIServerNameKey:
		err = p.formMayaAPIServerConfig(component)
	case types.MayaAPIServerServiceNameKey:
		err = p.formMayaAPIServerServiceConfig(component)
	case types.NDMNameKey:
		err = p.formNDMDaemonConfig(component)
	case types.NDMOperatorNameKey:
		err = p.formNDMOperatorConfig(component)
	case types.NDMConfigNameKey:
		err = p.formNDMConfigMapConfig(component)
	case types.ProvisionerNameKey:
		err = p.formOpenEBSProvisionerConfig(component)
	case types.LocalProvisionerNameKey:
		err = p.formLocalProvisionerConfig(component)
	case types.SnapshotOperatorNameKey:
		err = p.formSnapshotOperatorConfig(component)
	case types.AdmissionServerNameKey:
		err = p.formAdmissionServerConfig(component)
	case types.AdmissionServerSVCNameKey:
		// do nothing as we don't have any OpenEBS config for admission-server-svc
	case types.CSPCOperatorNameKey:
		err = p.formCSPCOperatorConfig(component)
	case types.CVCOperatorNameKey:
		err = p.formCVCOperatorConfig(component)
	case types.CVCOperatorServiceNameKey:
		err = p.formCVCOperatorServiceConfig(component)
	case types.CStorCSINodeNameKey:
		err = p.formCstorCSINode(component)
	case types.CStorCSIControllerNameKey:
		err = p.formCstorCSIController(component)
	case types.CStorAdmissionServerNameKey:
		err = p.formCstorAdmissionServer(component)
	case types.MoacDeploymentNameKey:
		err = p.formMayastorMOACConfig(component)
	case types.MayastorMOACSVCNameKey:
		err = p.formMOACServiceConfig(component)
	case types.MayastorDaemonsetNameKey:
		err = p.formMayastorDaemonConfig(component)
	}
	if err != nil {
		return err
	}

	return nil
}

// formDefaultStoragePath forms the default OpenEBS directory.
func (p *Planner) formDefaultStoragePath() error {
	// If default storage path is not set using the environment variable OPENEBS_IO_BASE_DIR, etc.
	if len(p.DefaultStoragePath) == 0 {
		p.DefaultStoragePath = "/var/openebs"
	}
	return nil
}

func (p *Planner) formCommonOpenEBSConfig() error {
	var initFuncs = []func() error{
		p.formDefaultStoragePath,
		p.formJivaConfig,
		p.formCStorConfig,
		p.formAnalyticsConfig,
		p.formPoliciesConfig,
		p.formHelperConfig,
	}
	for _, fn := range initFuncs {
		err := fn()
		if err != nil {
			return err
		}
	}
	return nil
}

// getResourceCommonDetails gets the common details of a given resource/component such
// as name, nodeSelectors, affinity, tolerations, etc.
func (p *Planner) getResourceCommonDetails(resource *unstructured.Unstructured,
	resDetails map[string]interface{}) (map[string]interface{},
	error) {
	var err error
	resourceDetails := make(map[string]interface{}, 0)
	if resDetails != nil {
		resourceDetails = resDetails
	}
	// set this resource's enabled value to true since this fn will be called only if the resource
	// is present.
	resourceDetails[types.KeyEnabled] = true
	// get the resource name
	resourceDetails[types.KeyName] = resource.GetName()
	// get the no of replicas if its not a daemonset
	if resource.GetKind() != types.KindDaemonSet {
		resourceDetails[types.KeyReplicas], _, err = unstructured.NestedInt64(resource.Object,
			"spec", "replicas")
	}
	// get the nodeSelectors of this resource
	resourceDetails[types.KeyNodeSelector], _, err = unstructured.NestedMap(resource.Object,
		"spec", "template", "spec", "nodeSelector")
	if err != nil {
		return resourceDetails, err
	}
	// get the affinity of this resource
	resourceDetails[types.KeyAffinity], _, err = unstructured.NestedMap(resource.Object,
		"spec", "template", "spec", "affinity")
	if err != nil {
		return resourceDetails, err
	}
	// get the tolerations of this resource
	resourceDetails[types.KeyTolerations], _, err = unstructured.NestedSlice(resource.Object,
		"spec", "template", "spec", "tolerations")
	if err != nil {
		return resourceDetails, err
	}
	return resourceDetails, nil
}
