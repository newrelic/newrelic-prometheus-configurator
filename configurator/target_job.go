package configurator

// ScrapeJobInput represents a target job config for configurator.
type ScrapeJobInput struct {
	TargetJob                 TargetJob               `yaml:",inline"`
	Targets                   []string                `yaml:"targets"`
	Labels                    map[string]string       `yaml:"labels"`
	ExtraRelabelConfigs       []PrometheusExtraConfig `yaml:"extra_relabel_config"`
	ExtraMetricRelabelConfigs []PrometheusExtraConfig `yaml:"extra_metric_relabel_config"`
}

// TargetJobOutput represents a prometheus scrape_config Job config with static_configs which can be obtained from input.
type TargetJobOutput struct {
	TargetJob            TargetJob               `yaml:",inline"`
	StaticConfigs        []StaticConfig          `yaml:"static_configs,omitempty"`
	RelabelConfigs       []PrometheusExtraConfig `yaml:"relabel_configs,omitempty"`
	MetricRelabelConfigs []PrometheusExtraConfig `yaml:"metric_relabel_configs,omitempty"`
	KubernetesSdConfigs  []map[string]string     `yaml:"kubernetes_sd_config,omitempty"`
}

func BuildTargetJob(job ScrapeJobInput) TargetJobOutput {
	return TargetJobOutput{
		TargetJob: job.TargetJob,
		StaticConfigs: []StaticConfig{
			{
				Targets: job.Targets,
				Labels:  job.Labels,
			},
		},
		RelabelConfigs:       job.ExtraRelabelConfigs,
		MetricRelabelConfigs: job.ExtraMetricRelabelConfigs,
	}
}
