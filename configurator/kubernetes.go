package configurator

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

// KubernetesInput defines all fields to set up prometheus.
type KubernetesInput struct {
	Jobs []KubernetesJob `yaml:"jobs"`
}

// KubernetesJob holds the configuration which will parsed to a prometheus scrape job including the
// specific rules needed.
type KubernetesJob struct {
	Job                       promcfg.Job                     `yaml:",inline"`
	JobNamePrefix             string                          `yaml:"job_name_prefix"`
	TargetDiscovery           KubernetesTargetDiscovery       `yaml:"target_discovery"`
	ExtraRelabelConfigs       []promcfg.PrometheusExtraConfig `yaml:"extra_relabel_config"`
	ExtraMetricRelabelConfigs []promcfg.PrometheusExtraConfig `yaml:"extra_metric_relabel_config"`
}

type KubernetesTargetDiscovery struct {
	Pod              bool              `yaml:"pod"`
	Endpoints        bool              `yaml:"endpoints"`
	Filter           *Filter           `yaml:"filter,omitempty"`
	AdditionalConfig *AdditionalConfig `yaml:"additional_config,omitempty"`
}

// AdditionalConfig holds additional config for the service discovery.
type AdditionalConfig struct {
	KubeconfigFile string                          `yaml:"kubeconfig_file,omitempty"`
	Namespaces     *promcfg.KubernetesSdNamespace  `yaml:"namespaces,omitempty"`
	Selectors      *[]promcfg.KubernetesSdSelector `yaml:"selectors,omitempty"`
	AttachMetadata *promcfg.AttachMetadata         `yaml:"attach_metadata,omitempty"`
}

// Valid returns true when the defined configuration is valid.
func (k *KubernetesTargetDiscovery) Valid() bool {
	return k.Pod || k.Endpoints
}

// KubernetesSettingsBuilders defines a functions which updates and returns a `TargetJobOutput` with specific settings
// added (considering the specified `KubernetesJob`).
type kubernetesSettingsBuilder func(job promcfg.Job, k8sJob KubernetesJob) promcfg.Job

// KubernetesJobBuilder holds the specific settings to add to a TargetJobOutput given the
// corresponding KubernetesJob definition.
type KubernetesJobBuilder struct {
	addPodSettings       kubernetesSettingsBuilder
	addEndpointsSettings kubernetesSettingsBuilder
}

// NewKubernetesJobBuilder creates a builder using the default settings builders.
func NewKubernetesJobBuilder() *KubernetesJobBuilder {
	return &KubernetesJobBuilder{
		addPodSettings:       podSettingsBuilder,
		addEndpointsSettings: endpointSettingsBuilder,
	}
}

// Build builds the prometheus targets corresponding to the Kubernetes configuration in input.
func (b *KubernetesJobBuilder) Build(i *Input) ([]promcfg.Job, error) {
	var jobs []promcfg.Job

	for _, k8sJob := range i.Kubernetes.Jobs {
		if !k8sJob.TargetDiscovery.Valid() {
			return nil, ErrInvalidK8sJobKinds
		}

		if err := b.checkJob(k8sJob); err != nil {
			return nil, err
		}

		if k8sJob.TargetDiscovery.Pod && b.addPodSettings != nil {
			job := k8sJob.Job

			job.JobName = b.buildJobName(k8sJob.JobNamePrefix, podKind)

			if f := k8sJob.TargetDiscovery.Filter; f != nil {
				job.RelabelConfigs = append(job.RelabelConfigs, f.Pod())
			}

			job = b.addPodSettings(job, k8sJob)

			for _, c := range k8sJob.ExtraRelabelConfigs {
				job.RelabelConfigs = append(job.RelabelConfigs, c)
			}

			for _, c := range k8sJob.ExtraMetricRelabelConfigs {
				job.MetricRelabelConfigs = append(job.MetricRelabelConfigs, c)
			}

			jobs = append(jobs, job)
		}

		if k8sJob.TargetDiscovery.Endpoints && b.addEndpointsSettings != nil {
			job := k8sJob.Job

			job.JobName = b.buildJobName(k8sJob.JobNamePrefix, endpointsKind)

			if f := k8sJob.TargetDiscovery.Filter; f != nil {
				job.RelabelConfigs = append(job.RelabelConfigs, f.Endpoints())
			}

			job = b.addEndpointsSettings(job, k8sJob)

			for _, c := range k8sJob.ExtraRelabelConfigs {
				job.RelabelConfigs = append(job.RelabelConfigs, c)
			}

			for _, c := range k8sJob.ExtraMetricRelabelConfigs {
				job.MetricRelabelConfigs = append(job.MetricRelabelConfigs, c)
			}

			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}

func (b *KubernetesJobBuilder) buildJobName(prefix string, kind string) string {
	return prefix + "-" + kind
}

func (b *KubernetesJobBuilder) checkJob(k8sJob KubernetesJob) error {
	if k8sJob.JobNamePrefix == "" {
		return ErrInvalidK8sJobPrefix
	}

	return nil
}
