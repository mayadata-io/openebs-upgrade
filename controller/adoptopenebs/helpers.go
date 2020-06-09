package adoptopenebs

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/types"
)

// formHelperConfig forms the desired OpenEBS CR config for helper tools.
func (p *Planner) formHelperConfig() error {
	helperConfig := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	if p.HelperConfig != nil {
		helperConfig = p.HelperConfig
	}
	if p.HelperImageTag != p.OpenEBSVersion {
		helperConfig.Object[types.KeyImageTag] = p.HelperImageTag
	}
	p.HelperConfig = helperConfig

	return nil
}
