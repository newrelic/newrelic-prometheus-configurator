package kubernetes

import (
	"errors"
	"fmt"
	"strings"

	"github.com/newrelic/newrelic-prometheus-configurator/internal/promcfg"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/scrapejob"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/sharding"
)

const (
	podKind       = "pod"
	endpointsKind = "endpoints"
)

var ErrIntegrationFilterConfig = errors.New("neither default or config specified")

var (
	ErrInvalidK8sJobKinds      = errors.New("at least one kind should be set in target_kinds field")
	ErrInvalidK8sJobPrefix     = errors.New("prefix cannot be empty in kubernetes jobs")
	ErrInvalidSkipShardingFlag = errors.New("kubernetes jobs do not support skip_sharding flag")
)

// Config defines all fields to set up prometheus to scrape k8s targets.
type Config struct {
	K8sJobs           []K8sJob          `yaml:"jobs"`
	IntegrationFilter IntegrationFilter `yaml:"integrations_filter"`
}

// IntegrationFilter holds the configuration for the IntegrationFilter filtering.
type IntegrationFilter struct {
	Enabled      *bool    `yaml:"enabled"`
	SourceLabels []string `yaml:"source_labels"`
	AppValues    []string `yaml:"app_values"`
}

// This struct is used internally to improve readability of function signatures.
type jobRelabelConfig struct {
	endpoints []promcfg.RelabelConfig
	pods      []promcfg.RelabelConfig
}

// Build will create a Prometheus Job list based on the kubernetes configuration.
func (c Config) Build(shardingConfig sharding.Config) ([]promcfg.Job, error) {
	var promScrapeJobs []promcfg.Job

	for _, k8sJob := range c.K8sJobs {
		if err := c.validate(k8sJob); err != nil {
			return nil, err
		}

		jrc, err := c.buildRelabelConfig(k8sJob)
		if err != nil {
			return nil, fmt.Errorf("building relabel configs: %w", err)
		}

		if k8sJob.TargetDiscovery.Pod {
			promJob := buildPromJob(shardingConfig, k8sJob, podKind, jrc.pods)
			promScrapeJobs = append(promScrapeJobs, promJob)
		}

		if k8sJob.TargetDiscovery.Endpoints {
			promJob := buildPromJob(shardingConfig, k8sJob, endpointsKind, jrc.endpoints)
			promScrapeJobs = append(promScrapeJobs, promJob)
		}
	}

	return promScrapeJobs, nil
}

func buildPromJob(shardingConfig sharding.Config, k8sJob K8sJob, objPrefix string, relabelConfig []promcfg.RelabelConfig) promcfg.Job {
	jobName := k8sJob.JobNamePrefix + "-" + objPrefix
	promJob := k8sJob.ScrapeJob.
		WithName(jobName).
		WithRelabelConfigs(relabelConfig).
		BuildPrometheusJob(shardingConfig)
	promJob.KubernetesSdConfigs = append(promJob.KubernetesSdConfigs, buildSdConfig(objPrefix, k8sJob.TargetDiscovery.AdditionalConfig))

	return promJob
}

func (c Config) buildRelabelConfig(k8sJob K8sJob) (jobRelabelConfig, error) {
	jrc := jobRelabelConfig{pods: podRelabelConfigs(k8sJob), endpoints: endpointsRelabelConfigs(k8sJob)}

	if !integrationFilterToBeApplied(c.IntegrationFilter, k8sJob.IntegrationFilter) {
		return jrc, nil
	}

	crRelabelConfig, err := buildIntegrationFilter(c.IntegrationFilter, k8sJob.IntegrationFilter)
	if err != nil {
		return jobRelabelConfig{}, fmt.Errorf("building relabel configs for integration filters: %w", err)
	}

	jrc.endpoints = append(jrc.endpoints, crRelabelConfig.endpoints...)
	jrc.pods = append(jrc.pods, crRelabelConfig.pods...)

	return jrc, nil
}

func integrationFilterToBeApplied(filters IntegrationFilter, jobFilterfilters IntegrationFilter) bool {
	if jobFilterfilters.Enabled != nil {
		return *jobFilterfilters.Enabled
	} else if filters.Enabled != nil {
		return *filters.Enabled
	}

	return false
}

func buildIntegrationFilter(filters IntegrationFilter, jobFilters IntegrationFilter) (jobRelabelConfig, error) {
	filterLabels, err := getConfigWithFallback(filters.SourceLabels, jobFilters.SourceLabels)
	if err != nil {
		return jobRelabelConfig{}, fmt.Errorf("source labels are empty for both the default and the job integration filters: %w", err)
	}

	filterAppValues, err := getConfigWithFallback(filters.AppValues, jobFilters.AppValues)
	if err != nil {
		return jobRelabelConfig{}, fmt.Errorf("filter app values are empty for both the default and the job integration filters: %w", err)
	}

	sourceLabelsPod := make([]string, 0, len(filterLabels))
	sourceLabelsEndpoint := make([]string, 0, len(filterLabels))
	for _, fL := range filterLabels {
		sanitizedLabel := invalidPrometheusLabelCharRegex.ReplaceAllString(fL, "_")

		sourceLabelsPod = append(sourceLabelsPod, fmt.Sprintf("%s%s_%s", podMetadata, labelMetadata, sanitizedLabel))
		sourceLabelsEndpoint = append(sourceLabelsEndpoint, fmt.Sprintf("%s%s_%s", serviceMetadata, labelMetadata, sanitizedLabel))
	}

	regex := strings.Join(filterAppValues, "|")
	caseInsensitiveRegex := fmt.Sprintf("(?i)(%s)", regex)
	unanchoredRegex := fmt.Sprintf(".*%s.*", caseInsensitiveRegex)

	return jobRelabelConfig{
		endpoints: []promcfg.RelabelConfig{
			{
				SourceLabels: sourceLabelsEndpoint,
				Separator:    ";",
				Regex:        unanchoredRegex,
				Action:       "keep",
			},
		},
		pods: []promcfg.RelabelConfig{
			{
				SourceLabels: sourceLabelsPod,
				Separator:    ";",
				Regex:        unanchoredRegex,
				Action:       "keep",
			},
		},
	}, nil
}

func getConfigWithFallback(defaultConfig []string, config []string) ([]string, error) {
	if len(config) != 0 {
		return config, nil
	} else if len(defaultConfig) != 0 {
		return defaultConfig, nil
	}

	return nil, ErrIntegrationFilterConfig
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
	ScrapeJob         scrapejob.Job     `yaml:",inline"`
	JobNamePrefix     string            `yaml:"job_name_prefix"`
	TargetDiscovery   TargetDiscovery   `yaml:"target_discovery"`
	IntegrationFilter IntegrationFilter `yaml:"integrations_filter"`
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
