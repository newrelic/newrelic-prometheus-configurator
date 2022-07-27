package configurator

// podSettingsBuilder returns a copy of `tg` including the specific settings for when pods kind is set.
func podSettingsBuilder(tg TargetJobOutput, _ KubernetesJob) TargetJobOutput {
	tg.KubernetesSdConfigs = []map[string]string{{"role": "pod"}}
	tg.RelabelConfigs = append(tg.RelabelConfigs)
	return tg
}
