package configurator

// endpointSettingsBuilder returns a copy of `job` including the specific settings for when endpoints kind is set.
func endpointSettingsBuilder(job *JobOutput, input KubernetesJob) {
	job.Job.HonorLabels = true

	kubernetesSdConfig := setK8sSdConfigFromJob("endpoints", input)

	job.KubernetesSdConfigs = []KubernetesSdConfig{kubernetesSdConfig}

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
}

// setK8sSdConfigFromJob populates a KubernetesSdConfig from a given KubernetesJob.
func setK8sSdConfigFromJob(role string, input KubernetesJob) KubernetesSdConfig {
	k8sSdConfig := KubernetesSdConfig{
		Role: role,
	}

	if input.TargetDiscovery.AdditionalConfig == nil {
		return k8sSdConfig
	}

	k8sSdConfig.KubeconfigFile = input.TargetDiscovery.AdditionalConfig.KubeconfigFile

	if input.TargetDiscovery.AdditionalConfig.Namespaces != nil {
		k8sSdConfig.Namespaces = input.TargetDiscovery.AdditionalConfig.Namespaces
	}

	if input.TargetDiscovery.AdditionalConfig.Selectors != nil {
		k8sSdConfig.Selectors = input.TargetDiscovery.AdditionalConfig.Selectors
	}

	if input.TargetDiscovery.AdditionalConfig.AttachMetadata != nil &&
		input.TargetDiscovery.AdditionalConfig.AttachMetadata.Node != nil {
		k8sSdConfig.AttachMetadata = &AttachMetadata{Node: input.TargetDiscovery.AdditionalConfig.AttachMetadata.Node}
	}

	return k8sSdConfig
}
