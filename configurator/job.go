package configurator

// JobInput represents job config for configurator.
type JobInput struct {
	Job                       Job                     `yaml:",inline"`
	Targets                   []string                `yaml:"targets"`
	Labels                    map[string]string       `yaml:"labels"`
	ExtraRelabelConfigs       []PrometheusExtraConfig `yaml:"extra_relabel_config"`
	ExtraMetricRelabelConfigs []PrometheusExtraConfig `yaml:"extra_metric_relabel_config"`
}

// JobOutput represents a prometheus scrape_config Job config with static_configs which can be obtained from input.
type JobOutput struct {
	Job                  Job                     `yaml:",inline"`
	StaticConfigs        []StaticConfig          `yaml:"static_configs,omitempty"`
	RelabelConfigs       []PrometheusExtraConfig `yaml:"relabel_configs,omitempty"`
	MetricRelabelConfigs []PrometheusExtraConfig `yaml:"metric_relabel_configs,omitempty"`
	KubernetesSdConfigs  []map[string]string     `yaml:"kubernetes_sd_configs,omitempty"`
}

func BuildJobOutput(job JobInput) JobOutput {
	jobOutput := JobOutput{
		Job:                  job.Job,
		RelabelConfigs:       job.ExtraRelabelConfigs,
		MetricRelabelConfigs: job.ExtraMetricRelabelConfigs,
	}

	if (len(job.Targets) > 0) || (len(job.Labels) > 0) {
		jobOutput.StaticConfigs = []StaticConfig{
			{
				Targets: job.Targets,
				Labels:  job.Labels,
			},
		}
	}

	return jobOutput
}
