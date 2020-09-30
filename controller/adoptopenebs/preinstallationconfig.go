package adoptopenebs

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

// formPreInstallationConfig forms the desired OpenEBS CR config for components/tools/dependencies
// which are required to be installed prior to OpenEBS installation.
func (p *Planner) formPreInstallationConfig() error {
	preInstallationConfig := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	if p.PreInstallationConfig != nil {
		preInstallationConfig = p.PreInstallationConfig
	}
	// form the config for ISCSI client installation, for adoption cases,
	// it will be false by default since OpenEBS is already up and running.
	iscsiClientConfig := make(map[string]interface{})
	iscsiClientConfig["enabled"] = false
	iscsiClientConfig["isSetupDone"] = true
	preInstallationConfig.Object["iscsiClient"] = iscsiClientConfig
	p.PreInstallationConfig = preInstallationConfig

	return nil
}
