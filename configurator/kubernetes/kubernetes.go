package kubernetes

import (
	"errors"

	"github.com/newrelic-forks/newrelic-prometheus/configurator/promcfg"
)

const (
	podKind       = "pod"
	endpointsKind = "endpoints"
)

var (
	ErrInvalidK8sJobKinds  = errors.New("at least one kind should be set in target_kinds field")
	ErrInvalidK8sJobPrefix = errors.New("prefix cannot be empty in kubernetes jobs")
)

// Input defines all fields to set up prometheus.
type Input struct {
	Jobs []Job `yaml:"jobs"`
}

func (i Input) Build() ([]promcfg.Job, error) {
	var jobs []promcfg.Job

	for _, job := range i.Jobs {
		if err := i.validate(job); err != nil {
			return nil, err
		}

		if job.TargetDiscovery.Pod {
			podJob := job.Job

			podJob.JobName = job.JobNamePrefix + "-" + podKind

			podJob.KubernetesSdConfigs = append(podJob.KubernetesSdConfigs, buildSdConfig(podKind, job.TargetDiscovery.AdditionalConfig))

			podJob.RelabelConfigs = append(podJob.RelabelConfigs, podRelableConfigs(job)...)

			for _, c := range job.ExtraMetricRelabelConfigs {
				podJob.MetricRelabelConfigs = append(podJob.MetricRelabelConfigs, c)
			}

			jobs = append(jobs, podJob)
		}

		if job.TargetDiscovery.Endpoints {
			endpointsJob := job.Job

			endpointsJob.JobName = job.JobNamePrefix + "-" + endpointsKind

			endpointsJob.KubernetesSdConfigs = append(endpointsJob.KubernetesSdConfigs, buildSdConfig(endpointsKind, job.TargetDiscovery.AdditionalConfig))

			endpointsJob.RelabelConfigs = append(endpointsJob.RelabelConfigs, endpointsRelableConfigs(job)...)

			for _, c := range job.ExtraMetricRelabelConfigs {
				endpointsJob.MetricRelabelConfigs = append(endpointsJob.MetricRelabelConfigs, c)
			}

			jobs = append(jobs, endpointsJob)
		}
	}

	return jobs, nil
}

func (i Input) validate(job Job) error {
	if !job.TargetDiscovery.Valid() {
		return ErrInvalidK8sJobKinds
	}

	if job.JobNamePrefix == "" {
		return ErrInvalidK8sJobPrefix
	}

	return nil
}

// Job holds the configuration which will parsed to a prometheus scrape job including the
// specific rules needed.
type Job struct {
	Job                       promcfg.Job                     `yaml:",inline"`
	JobNamePrefix             string                          `yaml:"job_name_prefix"`
	TargetDiscovery           TargetDiscovery                 `yaml:"target_discovery"`
	ExtraRelabelConfigs       []promcfg.PrometheusExtraConfig `yaml:"extra_relabel_config"`
	ExtraMetricRelabelConfigs []promcfg.PrometheusExtraConfig `yaml:"extra_metric_relabel_config"`
}

type TargetDiscovery struct {
	Pod              bool              `yaml:"pod"`
	Endpoints        bool              `yaml:"endpoints"`
	Filter           Filter            `yaml:"filter,omitempty"`
	AdditionalConfig *AdditionalConfig `yaml:"additional_config,omitempty"`
}

// Valid returns true when the defined configuration is valid.
func (k *TargetDiscovery) Valid() bool {
	return k.Pod || k.Endpoints
}

// AdditionalConfig holds additional config for the service discovery.
type AdditionalConfig struct {
	KubeconfigFile string                          `yaml:"kubeconfig_file,omitempty"`
	Namespaces     *promcfg.KubernetesSdNamespace  `yaml:"namespaces,omitempty"`
	Selectors      *[]promcfg.KubernetesSdSelector `yaml:"selectors,omitempty"`
	AttachMetadata *promcfg.AttachMetadata         `yaml:"attach_metadata,omitempty"`
}

func buildSdConfig(kind string, ac *AdditionalConfig) promcfg.KubernetesSdConfig {
	k8sSdConfig := promcfg.KubernetesSdConfig{
		Role: kind,
	}

	if ac == nil {
		return k8sSdConfig
	}

	k8sSdConfig.KubeconfigFile = ac.KubeconfigFile

	if ac.Namespaces != nil {
		k8sSdConfig.Namespaces = ac.Namespaces
	}

	if ac.Selectors != nil {
		k8sSdConfig.Selectors = ac.Selectors
	}

	if ac.AttachMetadata != nil &&
		ac.AttachMetadata.Node != nil {
		k8sSdConfig.AttachMetadata = &promcfg.AttachMetadata{Node: ac.AttachMetadata.Node}
	}

	return k8sSdConfig
}
