package adoptopenebs

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/types"
)

// formPoliciesConfig forms the desired OpenEBS CR config for policies.
func (p *Planner) formPoliciesConfig() error {
	policiesConfig := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	if p.PoliciesConfig != nil {
		policiesConfig = p.PoliciesConfig
	}
	monitoringConfig := make(map[string]interface{}, 0)
	monitoringConfig[types.KeyEnabled] = true
	if p.PoolExporterImageTag != p.OpenEBSVersion {
		monitoringConfig[types.KeyImageTag] = p.PoolExporterImageTag
	}
	policiesConfig.Object["monitoring"] = monitoringConfig
	p.PoliciesConfig = policiesConfig

	return nil
}
