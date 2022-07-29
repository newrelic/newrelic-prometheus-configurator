package configurator

// endpointSettingsBuilder returns a copy of `tg` including the specific settings for when endpoints kind is set.
func endpointSettingsBuilder(tg JobOutput, _ KubernetesJob) JobOutput {
	// TODO: include the specific settings
	return tg
}
