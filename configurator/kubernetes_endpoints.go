package configurator

// endpointSettingsBuilder returns a copy of `tg` including the specific settings for when endpoints kind is set.
func endpointSettingsBuilder(job JobOutput, _ KubernetesJob) JobOutput {
	job.Job.HonorLabels = true
	job.KubernetesSdConfigs = []map[string]string{
		{"role": "endpoints"},
	}
	job.RelabelConfigs = append(job.RelabelConfigs,
		RelabelConfig{
			SourceLabels: []string{"__meta_kubernetes_service_annotation_prometheus_io_scrape"},
			Action:       "keep",
			Regex:        "true",
		},
		RelabelConfig{
			SourceLabels: []string{"__meta_kubernetes_service_annotation_prometheus_io_scrape_slow"},
			Action:       "drop",
			Regex:        "true",
		},
		RelabelConfig{
			SourceLabels: []string{"__meta_kubernetes_service_annotation_prometheus_io_scheme"},
			Action:       "replace",
			TargetLabel:  "__metrics_path__",
			Regex:        `(.+)`,
		},
		RelabelConfig{
			SourceLabels: []string{"__meta_kubernetes_pod_annotation_prometheus_io_path"},
			Action:       "replace",
			Regex:        "(.+)",
			TargetLabel:  "__metrics_path__",
		},
		RelabelConfig{
			SourceLabels: []string{"__address__", "__meta_kubernetes_service_annotation_prometheus_io_port"},
			Action:       "replace",
			TargetLabel:  "__address__",
			Regex:        `([^:]+)(?::\d+)?;(\d+)`,
			Replacement:  "$1:$2",
		},
		RelabelConfig{
			SourceLabels: []string{"__scheme__", "__address__", "__metrics_path__"},
			Action:       "replace",
			TargetLabel:  "scrapedTargetURL",
			Regex:        `(.+);(.+);(.+)`,
			Replacement:  "$1://$2$3",
		},
		RelabelConfig{
			SourceLabels: []string{"__meta_kubernetes_endpoints_name"},
			Action:       "replace",
			TargetLabel:  "targetName",
		},
		RelabelConfig{
			SourceLabels: []string{"job"},
			Action:       "replace",
			TargetLabel:  "scrappedTargetKind",
		},
		RelabelConfig{
			Action: "labeldrop",
			Regex:  "job",
		},
		RelabelConfig{
			Action:      "labelmap",
			Regex:       "__meta_kubernetes_pod_label_(.+)",
			Replacement: "label_$1",
		},
		RelabelConfig{
			SourceLabels: []string{"__meta_kubernetes_namespace"},
			Action:       "replace",
			TargetLabel:  "namespaceName",
		},
		RelabelConfig{
			SourceLabels: []string{"__meta_kubernetes_service_name"},
			Action:       "replace",
			TargetLabel:  "serviceName",
		},
		RelabelConfig{
			SourceLabels: []string{"__meta_kubernetes_pod_phase"},
			Action:       "drop",
			Regex:        "Pending|Succeeded|Failed|Completed",
		},
	)

	return job
}
