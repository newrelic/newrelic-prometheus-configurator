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
	Job                  Job                 `yaml:",inline"`
	StaticConfigs        []StaticConfig      `yaml:"static_configs,omitempty"`
	RelabelConfigs       []any               `yaml:"relabel_configs,omitempty"`
	MetricRelabelConfigs []any               `yaml:"metric_relabel_configs,omitempty"`
	KubernetesSdConfigs  []map[string]string `yaml:"kubernetes_sd_configs,omitempty"`
}

func BuildJobOutput(job JobInput) JobOutput {
	jobOutput := JobOutput{
		Job: job.Job,
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

func (o JobOutput) WithExtraConfigs(i JobInput) JobOutput {
	for _, c := range i.ExtraRelabelConfigs {
		o.RelabelConfigs = append(o.RelabelConfigs, c)
	}
	for _, c := range i.ExtraMetricRelabelConfigs {
		o.MetricRelabelConfigs = append(o.MetricRelabelConfigs, c)
	}
	return o
}
