package adoptopenebs

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

// formJivaConfig forms the desired OpenEBS CR config for Jiva.
func (p *Planner) formJivaConfig() error {
	jivaConfig := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	if p.JivaConfig != nil {
		jivaConfig = p.JivaConfig
	}
	if p.JivaCtrlImageTag != p.OpenEBSVersion {
		jivaConfig.Object["imageTag"] = p.JivaCtrlImageTag
	}
	jivaConfig.Object["replicas"] = p.JivaReplicaCount
	p.JivaConfig = jivaConfig

	return nil
}
