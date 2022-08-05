package configurator

// podSettingsBuilder includes the specific settings for pods in the provided JobOutput and returns it.
func podSettingsBuilder(job *JobOutput, kj KubernetesJob) {
	job.Job.HonorLabels = true

	kubernetesSdConfig := setK8sSdConfigFromJob("pod", kj)

	job.KubernetesSdConfigs = []KubernetesSdConfig{kubernetesSdConfig}
	job.RelabelConfigs = append(job.RelabelConfigs,
		RelabelConfig{
			SourceLabels: []string{"__meta_kubernetes_pod_phase"},
			Regex:        "Pending|Succeeded|Failed|Completed",
			Action:       "drop",
		},
		RelabelConfig{
			SourceLabels: []string{"__meta_kubernetes_pod_annotation_prometheus_io_scheme"},
			Action:       "replace",
			Regex:        "(https?)",
			TargetLabel:  "__scheme__",
		},
		RelabelConfig{
			SourceLabels: []string{"__meta_kubernetes_pod_annotation_prometheus_io_path"},
			Action:       "replace",
			Regex:        "(.+)",
			TargetLabel:  "__metrics_path__",
		},
		RelabelConfig{
			SourceLabels: []string{"__address__", "__meta_kubernetes_pod_annotation_prometheus_io_port"},
			Action:       "replace",
			Regex:        `(.+?)(?::\d+)?;(\d+)`,
			TargetLabel:  "__address__",
			Replacement:  "$1:$2",
		},
		RelabelConfig{
			Action:      "labelmap",
			Regex:       "__meta_kubernetes_pod_annotation_prometheus_io_param_(.+)",
			Replacement: "__param_$1",
		},
		RelabelConfig{
			Action: "labelmap",
			Regex:  "__meta_kubernetes_pod_label_(.+)",
		},
		RelabelConfig{
			SourceLabels: []string{"__meta_kubernetes_namespace"},
			Action:       "replace",
			TargetLabel:  "namespace",
		},
		RelabelConfig{
			SourceLabels: []string{"__meta_kubernetes_pod_node_name"},
			Action:       "replace",
			TargetLabel:  "node",
		},
		RelabelConfig{
			SourceLabels: []string{"__meta_kubernetes_pod_name"},
			Action:       "replace",
			TargetLabel:  "pod",
		},
	)
}
