package adoptopenebs

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/types"
)

// OpenEBSIdentifier helps in identifying the OpenEBS component type.
type OpenEBSIdentifier struct {
	Object *unstructured.Unstructured
}

// IdentifyOpenEBSVersion identifies the OpenEBS version which is installed or throws
// error if some error or unable to determine the OpenEBS version.
func (p *Planner) IdentifyOpenEBSVersion() error {
	var initFuncs = []func() error{
		// Add various functions or techniques to identify OpenEBS version.
		// All the functions will be run one by one until one of the functions
		// succeed in identifying the OpenEBS version or all fails.
		p.identifyOpenEBSVersionUsingLabels,
	}
	for _, fn := range initFuncs {
		err := fn()
		if err != nil {
			return err
		}
		if p.OpenEBSVersion != "" {
			return nil
		}
	}
	// If none of the above functions are able to find the OpenEBS version, this fn
	// will error out.
	if p.OpenEBSVersion == "" {
		return errors.Errorf("Unable to determine the OpenEBS version via the available methods")
	}
	return nil
}

func (p *Planner) identifyOpenEBSVersionUsingLabels() error {
	for _, component := range p.ObservedAdoptOpenEBSComponents {
		// Make use of `openebs.io/version` OpenEBS label in order to identify
		// the OpenEBS version.
		componentLabels := component.GetLabels()
		if openEBSVersionLabelValue, exist := componentLabels[types.OpenEBSVersionLabelKey]; exist {
			p.OpenEBSVersion = openEBSVersionLabelValue
			break
		}
	}
	return nil
}

// IdentifyOpenEBSComponentType identifies the type of OpenEBS component i.e.,
// whether it is a NDM daemon, NDM operator, maya-apiserver, provisioner,
// localProvisioner, etc.
func (oi *OpenEBSIdentifier) IdentifyOpenEBSComponentType() (string, error) {
	var initFuncs = []func() (string, error){
		// Add various functions or techniques to identify an OpenEBS component.
		// All the functions will be run one by one until one of the functions
		// succeed in identifying the type of OpenEBS component here or all fails.
		oi.identifyOpenEBSComponentUsingLabels,
	}
	for _, fn := range initFuncs {
		componentType, err := fn()
		if err != nil {
			return "", err
		}
		if componentType != "" {
			return componentType, nil
		}
	}
	return "", nil
}

// identifyOpenEBSComponentUsingLabels tries to identify the type of OpenEBS component
// using various predefined labels.
func (oi *OpenEBSIdentifier) identifyOpenEBSComponentUsingLabels() (string, error) {
	var (
		componentIdentifier string
		componentType       string
	)
	// Make use of a number of predefined OpenEBS labels in order to identify
	// the component type such as `openebs.io/component-name`.
	componentLabels := oi.Object.GetLabels()
	if componentNameLabelValue, exist := componentLabels[types.ComponentNameLabelKey]; exist {
		componentIdentifier = componentNameLabelValue
	}

	switch componentIdentifier {
	case types.MayaAPIServerComponentNameLabelValue:
		componentType = types.MayaAPIServerNameKey
	case types.MayaAPIServerSVCComponentNameLabelValue:
		componentType = types.MayaAPIServerServiceNameKey
	case types.NDMComponentNameLabelValue:
		componentType = types.NDMNameKey
	case types.NDMOperatorComponentNameLabelValue:
		componentType = types.NDMOperatorNameKey
	case types.NDMConfigComponentNameLabelValue:
		componentType = types.NDMConfigNameKey
	case types.OpenEBSProvisionerComponentNameLabelValue:
		componentType = types.ProvisionerNameKey
	case types.LocalPVProvisionerComponentNameLabelValue:
		componentType = types.LocalProvisionerNameKey
	case types.SnapshotOperatorComponentNameLabelValue:
		componentType = types.SnapshotOperatorNameKey
	case types.AdmissionServerComponentNameLabelValue:
		componentType = types.AdmissionServerNameKey
	case types.AdmissionServerSVCComponentNameLabelValue:
		componentType = types.AdmissionServerSVCNameKey
	case types.CSPCOperatorComponentNameLabelValue:
		componentType = types.CSPCOperatorNameKey
	case types.CVCOperatorComponentNameLabelValue:
		componentType = types.CVCOperatorNameKey
	case types.CStorCSINodeComponentNameLabelValue:
		componentType = types.CStorCSINodeNameKey
	case types.CStorCSIControllerComponentNameLabelValue:
		componentType = types.CStorCSIControllerNameKey
	}

	return componentType, nil
}
