package openebs

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/types"
)

// setPreInstallationDefaultsIfNotSet sets the default values for the dependencies
// which are mandatory to be installed prior to OpenEBS installation such as ISCSI client
// if not already given.
func (p *Planner) setPreInstallationDefaultsIfNotSet() error {
	if p.ObservedOpenEBS.Spec.PreInstallation.ISCSIClient.Enabled == nil {
		p.ObservedOpenEBS.Spec.PreInstallation.ISCSIClient.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.PreInstallation.ISCSIClient.Enabled = true
	}
	return nil
}

// getPreInstallationManifests returns a mapping of component's "name_kind" to YAML of
// the components that are required to be installed prior to OpenEBS.
// Note: This method makes use of the various preinstallation components YAMLs to form this
// mapping such as ISCSI installation YAML.
func (p *Planner) getPreInstallationManifests() error {
	var err error
	// initialize component manifests field.
	if p.ComponentManifests == nil {
		p.ComponentManifests = make(map[string]*unstructured.Unstructured, 0)
	}
	// set the pre-installation defaults.
	err = p.setPreInstallationDefaultsIfNotSet()
	if err != nil {
		return err
	}
	if *p.ObservedOpenEBS.Spec.PreInstallation.ISCSIClient.Enabled &&
		!p.ObservedOpenEBS.Spec.PreInstallation.ISCSIClient.IsSetupDone {
		isISCSISetupComponentsRunning, err := p.getISCSISetupComponentsStatus()
		if err != nil {
			return err
		}
		if isISCSISetupComponentsRunning {
			return nil
		}
		// get the ISCSI installation related YAMLs
		iscsiYAMLMap, err := p.getISCSIInstallationManifest()
		if err != nil {
			return err
		}
		for key, value := range iscsiYAMLMap {
			if key == "_" {
				continue
			}
			p.ComponentManifests[key] = value
		}
		// do not delete the existing components i.e., components which are already running in case of
		// setting up ISCSI client.
		// This will be the case where user wants to run ISCSI client setup when OpenEBS components
		// are already running.
		for _, component := range p.observedOpenEBSComponents {
			if component.GetKind() == types.KindDaemonSet {
				if component.GetName() == types.OpenEBSNodeSetupDaemonsetNameKey {
					continue
				}
			} else if component.GetKind() == types.KindConfigMap {
				if component.GetName() == types.OpenEBSNodeSetupConfigmapNameKey {
					continue
				}
			}
			p.ComponentManifests[component.GetName()+"_"+component.GetKind()] = component
		}
	}

	return nil
}
