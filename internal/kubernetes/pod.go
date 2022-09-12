package kubernetes //nolint: dupl

import (
	"github.com/newrelic/newrelic-prometheus-configurator/internal/promcfg"
)

// podRelabelConfigs returns all relabel configs for an Pod job.
func podRelabelConfigs(job K8sJob) []promcfg.RelabelConfig {
	rc := []promcfg.RelabelConfig{}

	if job.TargetDiscovery.Filter.Valid() {
		rc = append(rc, job.TargetDiscovery.Filter.Pod())
	}

	rc = append(rc, podDefaultRelabelConfigs()...)

	return rc
}

func podDefaultRelabelConfigs() []promcfg.RelabelConfig {
	return []promcfg.RelabelConfig{
		{
			SourceLabels: []string{"__meta_kubernetes_pod_phase"},
			Regex:        "Pending|Succeeded|Failed|Completed",
			Action:       "drop",
		},
		{
			SourceLabels: []string{"__meta_kubernetes_pod_annotation_prometheus_io_scheme"},
			Action:       "replace",
			Regex:        "(https?)",
			TargetLabel:  "__scheme__",
		},
		{
			SourceLabels: []string{"__meta_kubernetes_pod_annotation_prometheus_io_path"},
			Action:       "replace",
			Regex:        "(.+)",
			TargetLabel:  "__metrics_path__",
		},
		{
			SourceLabels: []string{"__address__", "__meta_kubernetes_pod_annotation_prometheus_io_port"},
			Action:       "replace",
			Regex:        `(.+?)(?::\d+)?;(\d+)`,
			TargetLabel:  "__address__",
			Replacement:  "$1:$2",
		},
		{
			Action:      "labelmap",
			Regex:       "__meta_kubernetes_pod_annotation_prometheus_io_param_(.+)",
			Replacement: "__param_$1",
		},
		{
			Action: "labelmap",
			Regex:  "__meta_kubernetes_pod_label_(.+)",
		},
		{
			SourceLabels: []string{"__meta_kubernetes_namespace"},
			Action:       "replace",
			TargetLabel:  "namespace",
		},
		{
			SourceLabels: []string{"__meta_kubernetes_pod_node_name"},
			Action:       "replace",
			TargetLabel:  "node",
		},
		{
			SourceLabels: []string{"__meta_kubernetes_pod_name"},
			Action:       "replace",
			TargetLabel:  "pod",
		},
	}
}
