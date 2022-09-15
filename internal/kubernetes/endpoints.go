package kubernetes

import (
	"github.com/newrelic/newrelic-prometheus-configurator/internal/promcfg"
)

// endpointsRelabelConfigs returns all relabel configs for an Endpoints job.
func endpointsRelabelConfigs(job K8sJob) []promcfg.RelabelConfig {
	rc := []promcfg.RelabelConfig{}

	if job.TargetDiscovery.Filter.Valid() {
		rc = append(rc, job.TargetDiscovery.Filter.Endpoints())
	}

	rc = append(rc, endpointsDefaultRelabelConfigs()...)

	return rc
}

func endpointsDefaultRelabelConfigs() []promcfg.RelabelConfig {
	return []promcfg.RelabelConfig{
		{
			SourceLabels: []string{"__meta_kubernetes_pod_phase"},
			Action:       "drop",
			// Removed Pending status since #75.
			Regex: "Succeeded|Failed|Completed",
		},
		{
			SourceLabels: []string{"__meta_kubernetes_service_annotation_prometheus_io_scheme"},
			Action:       "replace",
			TargetLabel:  "__scheme__",
			Regex:        `(https?)`,
		},
		{
			SourceLabels: []string{"__meta_kubernetes_service_annotation_prometheus_io_path"},
			Action:       "replace",
			Regex:        `(.+)`,
			TargetLabel:  "__metrics_path__",
		},
		{
			SourceLabels: []string{"__address__", "__meta_kubernetes_service_annotation_prometheus_io_port"},
			Action:       "replace",
			TargetLabel:  "__address__",
			Regex:        `(.+?)(?::\d+)?;(\d+)`,
			Replacement:  "$1:$2",
		},
		{
			Action:      "labelmap",
			Regex:       `__meta_kubernetes_service_annotation_prometheus_io_param_(.+)`,
			Replacement: "__param_$1",
		},
		{
			Action: "labelmap",
			Regex:  `__meta_kubernetes_service_label_(.+)`,
		},
		{
			SourceLabels: []string{"__meta_kubernetes_namespace"},
			Action:       "replace",
			TargetLabel:  "namespace",
		},
		{
			SourceLabels: []string{"__meta_kubernetes_service_name"},
			Action:       "replace",
			TargetLabel:  "service",
		},
		// Following rule picks the node name which is not empty.
		// If both exists, use the pod_node_name.
		{
			SourceLabels: []string{
				"__meta_kubernetes_endpoint_node_name",
				"__meta_kubernetes_pod_node_name",
			},
			Separator: ";",
			Regex:     ".*;(.+)|(.+);",
			// only one of the groups will match.
			Replacement: "$1$2",
			Action:      "replace",
			TargetLabel: "node",
		},
	}
}
