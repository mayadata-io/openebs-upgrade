package adoptopenebs

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/types"
)

// formAnalyticsConfig forms the desired OpenEBS CR config for analytics.
func (p *Planner) formAnalyticsConfig() error {
	analyticsConfig := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	if p.AnalyticsConfig != nil {
		analyticsConfig = p.AnalyticsConfig
	}
	analyticsConfig.Object[types.KeyEnabled] = p.EnableAnalytics
	p.AnalyticsConfig = analyticsConfig

	return nil
}
