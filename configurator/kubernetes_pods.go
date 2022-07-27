package configurator

// podSettingsBuilder returns a copy of `tg` including the specific settings for when pods kind is set.
func podSettingsBuilder(tg JobOutput, _ KubernetesJob) JobOutput {
	return tg
}
