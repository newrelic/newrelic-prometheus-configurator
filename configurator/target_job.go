package configurator

import (
	"net/url"
	"time"

	"github.com/alecthomas/units"
)

// InputJob represents a target job config for configurator.
type InputJob struct {
	JobName                   string                  `yaml:"job_name"`
	Targets                   []string                `yaml:"targets"`
	Labels                    map[string]string       `yaml:"labels"`
	HonorLabels               bool                    `yaml:"honor_labels"`
	HonorTimestamps           bool                    `yaml:"honor_timestamps"`
	Params                    url.Values              `yaml:"params"`
	Scheme                    string                  `yaml:"scheme"`
	BodySizeLimit             units.Base2Bytes        `yaml:"body_size_limit"`
	SampleLimit               uint                    `yaml:"sample_limit"`
	TargetLimit               uint                    `yaml:"target_limit"`
	LabelLimit                uint                    `yaml:"label_limit"`
	LabelNameLengthLimit      uint                    `yaml:"label_name_length_limit"`
	LabelValueLengthLimit     uint                    `yaml:"label_value_length_limit"`
	MetricsPath               string                  `yaml:"metrics_path"`
	ScrapeInterval            time.Duration           `yaml:"scrape_interval"`
	ScrapeTimeout             time.Duration           `yaml:"scrape_timeout"`
	TLSConfig                 *TLSConfig              `yaml:"tls_config"`
	BasicAuth                 *BasicAuth              `yaml:"basic_auth"`
	Authorization             Authorization           `yaml:"authorization"`
	OAuth2                    OAuth2                  `yaml:"oauth2"`
	ExtraRelabelConfigs       []PrometheusExtraConfig `yaml:"extra_relabel_config"`
	ExtraMetricRelabelConfigs []PrometheusExtraConfig `yaml:"extra_metric_relabel_config"`
}

// StaticConfig defines each of the static_configs for the prometheus config.
type StaticConfig struct {
	Targets []string          `yaml:"targets"`
	Labels  map[string]string `yaml:"labels,omitempty"`
}

// TargetJobOutput represents a prometheus scrape_config Job config with static_configs which can be obtained from input.
type TargetJobOutput struct {
	JobName               string                  `yaml:"job_name"`
	StaticConfigs         []StaticConfig          `yaml:"static_configs,omitempty"`
	MetricsPath           string                  `yaml:"metrics_path,omitempty"`
	HonorLabels           bool                    `yaml:"honor_labels,omitempty"`
	HonorTimestamps       bool                    `yaml:"honor_timestamps"`
	Params                url.Values              `yaml:"params,omitempty"`
	Scheme                string                  `yaml:"scheme,omitempty"`
	BodySizeLimit         units.Base2Bytes        `yaml:"body_size_limit,omitempty"`
	SampleLimit           uint                    `yaml:"sample_limit,omitempty"`
	TargetLimit           uint                    `yaml:"target_limit,omitempty"`
	LabelLimit            uint                    `yaml:"label_limit,omitempty"`
	LabelNameLengthLimit  uint                    `yaml:"label_name_length_limit,omitempty"`
	LabelValueLengthLimit uint                    `yaml:"label_value_length_limit,omitempty"`
	ScrapeInterval        time.Duration           `yaml:"scrape_interval,omitempty"`
	ScrapeTimeout         time.Duration           `yaml:"scrape_timeout,omitempty"`
	TLSConfig             *TLSConfig              `yaml:"tls_config,omitempty"`
	BasicAuth             *BasicAuth              `yaml:"basic_auth,omitempty"`
	Authorization         Authorization           `yaml:"authorization,omitempty"`
	OAuth2                OAuth2                  `yaml:"oauth2,omitempty"`
	RelabelConfigs        []PrometheusExtraConfig `yaml:"relabel_configs,omitempty"`
	MetricRelabelConfigs  []PrometheusExtraConfig `yaml:"metric_relabel_configs,omitempty"`
}

func BuildTargetJob(job InputJob) TargetJobOutput {
	return TargetJobOutput{
		JobName: job.JobName,
		StaticConfigs: []StaticConfig{
			{
				Targets: job.Targets,
				Labels:  job.Labels,
			},
		},
		MetricsPath:           job.MetricsPath,
		HonorLabels:           job.HonorLabels,
		HonorTimestamps:       job.HonorTimestamps,
		Params:                job.Params,
		Scheme:                job.Scheme,
		BodySizeLimit:         job.BodySizeLimit,
		SampleLimit:           job.SampleLimit,
		TargetLimit:           job.TargetLimit,
		LabelLimit:            job.LabelLimit,
		LabelNameLengthLimit:  job.LabelNameLengthLimit,
		LabelValueLengthLimit: job.LabelValueLengthLimit,
		ScrapeInterval:        job.ScrapeInterval,
		ScrapeTimeout:         job.ScrapeTimeout,
		TLSConfig:             job.TLSConfig,
		BasicAuth:             job.BasicAuth,
		OAuth2:                job.OAuth2,
		Authorization:         job.Authorization,
		RelabelConfigs:        job.ExtraRelabelConfigs,
		MetricRelabelConfigs:  job.ExtraMetricRelabelConfigs,
	}
}
