package configurator

// KubernetesSelector defines the field needed to provided filtering capabilities to a kubernetes scrape job.
type KubernetesSelector struct {
	// TODO: define selector when this is implemented
}

// selectorSettingsBuilder returns a copy of `tg` including the specific settings for when selectors are defined.
func selectorSettingsBuilder(tg TargetJobOutput, _ KubernetesJob) TargetJobOutput {
	// TODO: include the specific settings
	return tg
}
