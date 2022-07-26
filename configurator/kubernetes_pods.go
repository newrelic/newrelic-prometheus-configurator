package configurator

// podSettingsBuilder returns a copy of `tg` including the specific settings for when pods kind is set.
func podSettingsBuilder(tg TargetJobOutput, _ KubernetesJob) TargetJobOutput {
	// TODO: include the specific settings
	return tg
}
