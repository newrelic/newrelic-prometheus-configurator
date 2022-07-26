package configurator

// endpointSettingsBuilder returns a copy of `tg` including the specific settings for when endpoints kind is set.
func endpointSettingsBuilder(tg TargetJobOutput, _ KubernetesJob) TargetJobOutput {
	// TODO: include the specific settings
	return tg
}
