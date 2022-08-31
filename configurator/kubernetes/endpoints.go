package kubernetes

import (
	"fmt"

	"github.com/newrelic-forks/newrelic-prometheus/configurator/promcfg"
)

// endpointsRelabelConfigs returns all relabel configs for an Endpoints job.
func endpointsRelabelConfigs(job K8sJob, sharding promcfg.Sharding) []promcfg.RelabelConfig {
	rc := []promcfg.RelabelConfig{}

	if job.TargetDiscovery.Filter.Valid() {
		rc = append(rc, job.TargetDiscovery.Filter.Endpoints())
	}

	if sharding.TotalShardsCount > 1 {
		rc = append(rc, shardingEndpointsRelabelConfigs(sharding)...)
	}

	rc = append(rc, endpointsDefaultRelabelConfigs()...)

	rc = append(rc, job.ExtraRelabelConfigs...)

	return rc
}

func endpointsDefaultRelabelConfigs() []promcfg.RelabelConfig {
	return []promcfg.RelabelConfig{
		{
			SourceLabels: []string{"__meta_kubernetes_pod_phase"},
			Action:       "drop",
			Regex:        "Pending|Succeeded|Failed|Completed",
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
		{
			SourceLabels: []string{"__meta_kubernetes_pod_node_name"},
			Action:       "replace",
			TargetLabel:  "node",
		},
	}
}

func shardingEndpointsRelabelConfigs(sharding promcfg.Sharding) []promcfg.RelabelConfig {
	if sharding.ShardIndex == "" {
		return []promcfg.RelabelConfig{}
	}

	return []promcfg.RelabelConfig{
		{
			SourceLabels: []string{"__address__", "_meta_kubernetes_service_name"},
			Modulus:      sharding.TotalShardsCount,
			Action:       "hashmod",
			TargetLabel:  "__tmp_hash",
		},
		{
			SourceLabels: []string{"__tmp_hash"},
			Regex:        fmt.Sprintf("^%v$", sharding.ShardIndex),
			Action:       "keep",
		},
	}
}
