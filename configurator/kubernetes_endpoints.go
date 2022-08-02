package configurator

// endpointSettingsBuilder returns a copy of `tg` including the specific settings for when endpoints kind is set.
func endpointSettingsBuilder(job JobOutput, input KubernetesJob) JobOutput {
	job.Job.HonorLabels = true

	kubernetesConfig := KubernetesSdConfig{
		Role: "endpoints",
	}

	if input.SdConfig != nil {
		kubernetesConfig.KubeconfigFile = input.SdConfig.KubeconfigFile
		kubernetesConfig.Namespaces = input.SdConfig.Namespaces
		kubernetesConfig.Selectors = input.SdConfig.Selectors
		if input.SdConfig.Node != nil {
			kubernetesConfig.AttachMetadata = &AttachMetadata{Node: input.SdConfig.Node}
		}
	}

	job.KubernetesSdConfigs = []KubernetesSdConfig{kubernetesConfig}

	job.RelabelConfigs = append(job.RelabelConfigs,
		RelabelConfig{
			SourceLabels: []string{"__meta_kubernetes_pod_phase"},
			Action:       "drop",
			Regex:        "Pending|Succeeded|Failed|Completed",
		},
		RelabelConfig{
			SourceLabels: []string{"__meta_kubernetes_service_annotation_prometheus_io_scheme"},
			Action:       "replace",
			TargetLabel:  "__scheme__",
			Regex:        `(https?)`,
		},
		RelabelConfig{
			SourceLabels: []string{"__meta_kubernetes_service_annotation_prometheus_io_path"},
			Action:       "replace",
			Regex:        `(.+)`,
			TargetLabel:  "__metrics_path__",
		},
		RelabelConfig{
			SourceLabels: []string{"__address__", "__meta_kubernetes_service_annotation_prometheus_io_port"},
			Action:       "replace",
			TargetLabel:  "__address__",
			Regex:        `(.+?)(?::\d+)?;(\d+)`,
			Replacement:  "$1:$2",
		},
		RelabelConfig{
			Action:      "labelmap",
			Regex:       `__meta_kubernetes_service_annotation_prometheus_io_param_(.+)`,
			Replacement: "__param_$1",
		},
		RelabelConfig{
			Action: "labelmap",
			Regex:  `__meta_kubernetes_service_label_(.+)`,
		},
		RelabelConfig{
			SourceLabels: []string{"__meta_kubernetes_namespace"},
			Action:       "replace",
			TargetLabel:  "namespace",
		},
		RelabelConfig{
			SourceLabels: []string{"__meta_kubernetes_service_name"},
			Action:       "replace",
			TargetLabel:  "service",
		},
		RelabelConfig{
			SourceLabels: []string{"__meta_kubernetes_pod_node_name"},
			Action:       "replace",
			TargetLabel:  "node",
		},
	)

	return job
}
