package kubernetes

import (
	"errors"

	"github.com/newrelic/newrelic-prometheus-configurator/internal/promcfg"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/scrapejobs"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/sharding"
)

const (
	podKind       = "pod"
	endpointsKind = "endpoints"
)

var (
	ErrInvalidK8sJobKinds      = errors.New("at least one kind should be set in target_kinds field")
	ErrInvalidK8sJobPrefix     = errors.New("prefix cannot be empty in kubernetes jobs")
	ErrInvalidSkipShardingFlag = errors.New("kubernetes jobs do not support skip_sharding flag")
)

// Config defines all fields to set up prometheus to scrape k8s targets.
type Config struct {
	K8sJobs []K8sJob `yaml:"jobs"`
}

// Build will create a Prometheus Job list based on the kubernetes configuration.
func (c Config) Build(shardingConfig sharding.Config) ([]promcfg.Job, error) {
	var promScrapeJobs []promcfg.Job

	for _, k8sJob := range c.K8sJobs {
		if err := c.validate(k8sJob); err != nil {
			return nil, err
		}

		if k8sJob.TargetDiscovery.Pod {
			podJob := k8sJob.ScrapeJob.
				WithName(k8sJob.JobNamePrefix + "-" + podKind).
				WithRelabelConfigs(podRelabelConfigs(k8sJob)).
				BuildPrometheusJob(shardingConfig)

			podJob.KubernetesSdConfigs = append(podJob.KubernetesSdConfigs, buildSdConfig(podKind, k8sJob.TargetDiscovery.AdditionalConfig))

			promScrapeJobs = append(promScrapeJobs, podJob)
		}

		if k8sJob.TargetDiscovery.Endpoints {
			endpointsJob := k8sJob.ScrapeJob.
				WithName(k8sJob.JobNamePrefix + "-" + endpointsKind).
				WithRelabelConfigs(endpointsRelabelConfigs(k8sJob)).
				BuildPrometheusJob(shardingConfig)

			endpointsJob.KubernetesSdConfigs = append(endpointsJob.KubernetesSdConfigs, buildSdConfig(endpointsKind, k8sJob.TargetDiscovery.AdditionalConfig))

			promScrapeJobs = append(promScrapeJobs, endpointsJob)
		}
	}

	return promScrapeJobs, nil
}

func (c Config) validate(k8sJob K8sJob) error {
	if !k8sJob.TargetDiscovery.Valid() {
		return ErrInvalidK8sJobKinds
	}

	if k8sJob.JobNamePrefix == "" {
		return ErrInvalidK8sJobPrefix
	}

	if k8sJob.ScrapeJob.SkipSharding {
		return ErrInvalidSkipShardingFlag
	}

	return nil
}

// K8sJob holds the configuration which will parsed to a prometheus scrape job including the
// specific rules needed.
type K8sJob struct {
	ScrapeJob       scrapejobs.Job  `yaml:",inline"`
	JobNamePrefix   string          `yaml:"job_name_prefix"`
	TargetDiscovery TargetDiscovery `yaml:"target_discovery"`
}

type TargetDiscovery struct {
	Pod              bool              `yaml:"pod"`
	Endpoints        bool              `yaml:"endpoints"`
	Filter           Filter            `yaml:"filter,omitempty"`
	AdditionalConfig *AdditionalConfig `yaml:"additional_config,omitempty"`
}

// Valid returns true when the defined configuration is valid.
func (td TargetDiscovery) Valid() bool {
	return td.Pod || td.Endpoints
}

// AdditionalConfig holds additional config for the service discovery.
type AdditionalConfig struct {
	KubeconfigFile string                          `yaml:"kubeconfig_file,omitempty"`
	Namespaces     *promcfg.KubernetesSdNamespace  `yaml:"namespaces,omitempty"`
	Selectors      *[]promcfg.KubernetesSdSelector `yaml:"selectors,omitempty"`
	AttachMetadata *promcfg.AttachMetadata         `yaml:"attach_metadata,omitempty"`
}

func buildSdConfig(jobKind string, ac *AdditionalConfig) promcfg.KubernetesSdConfig {
	k8sSdConfig := promcfg.KubernetesSdConfig{
		Role: jobKind,
	}

	// Check if Additional configs has been set in the config.
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
