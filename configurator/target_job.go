package configurator

import (
	"net/url"
	"time"

	"github.com/alecthomas/units"
)

// TargetJob holds fields which do not change from input and output jobs.
type TargetJob struct {
	JobName               string           `yaml:"job_name"`
	HonorLabels           bool             `yaml:"honor_labels,omitempty"`
	HonorTimestamps       bool             `yaml:"honor_timestamps"`
	Params                url.Values       `yaml:"params,omitempty"`
	Scheme                string           `yaml:"scheme,omitempty"`
	BodySizeLimit         units.Base2Bytes `yaml:"body_size_limit,omitempty"`
	SampleLimit           uint             `yaml:"sample_limit,omitempty"`
	TargetLimit           uint             `yaml:"target_limit,omitempty"`
	LabelLimit            uint             `yaml:"label_limit,omitempty"`
	LabelNameLengthLimit  uint             `yaml:"label_name_length_limit,omitempty"`
	LabelValueLengthLimit uint             `yaml:"label_value_length_limit,omitempty"`
	MetricsPath           string           `yaml:"metrics_path,omitempty"`
	ScrapeInterval        time.Duration    `yaml:"scrape_interval,omitempty"`
	ScrapeTimeout         time.Duration    `yaml:"scrape_timeout,omitempty"`
	TLSConfig             *TLSConfig       `yaml:"tls_config,omitempty"`
	BasicAuth             *BasicAuth       `yaml:"basic_auth,omitempty"`
	Authorization         Authorization    `yaml:"authorization,omitempty"`
	OAuth2                OAuth2           `yaml:"oauth2,omitempty"`
}

// TargetJobInput represents a target job config for configurator.
type TargetJobInput struct {
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
}

func BuildTargetJob(job TargetJobInput) TargetJobOutput {
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
