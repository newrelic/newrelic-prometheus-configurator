package configurator

import "errors"

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
	JobInput `yaml:",inline"`

	JobNamePrefix   string                    `yaml:"job_name_prefix"`
	Selector        *KubernetesSelector       `yaml:"selector,omitempty"`
	TargetDiscovery KubernetesTargetDiscovery `yaml:"target_discovery"`
}

type KubernetesTargetDiscovery struct {
	Pod       bool `yaml:"pod"`
	Endpoints bool `yaml:"endpoints"`
	// TODO: Filter Filter `yaml:"filter,omitempty"`. Implemented in Guillermo's PR.
	AdditionalConfig *AdditionalConfig `yaml:"additional_config,omitempty"`
}

// AdditionalConfig holds additional config for the service discovery.
type AdditionalConfig struct {
	KubeconfigFile string                  `yaml:"kubeconfig_file,omitempty"`
	Namespaces     *KubernetesSdNamespace  `yaml:"namespaces,omitempty"`
	Selectors      *[]KubernetesSdSelector `yaml:"selectors,omitempty"`
	AttachMetadata *AttachMetadata         `yaml:"attach_metadata,omitempty"`
}

// Valid returns true when the defined configuration is valid.
func (k *KubernetesTargetDiscovery) Valid() bool {
	return k.Pod || k.Endpoints
}

// KubernetesSettingsBuilders defines a functions which updates and returns a `TargetJobOutput` with specific settings
// added (considering the specified `KubernetesJob`).
type kubernetesSettingsBuilder func(job JobOutput, k8sJob KubernetesJob) JobOutput

// kubernetesJobBuilder holds the specific settings to add to a TargetJobOutput given the corresponding
// KubernetesJob definition.
type kubernetesJobBuilder struct {
	addPodSettings       kubernetesSettingsBuilder
	addEndpointsSettings kubernetesSettingsBuilder
	addSelectorSettings  kubernetesSettingsBuilder
}

// newKubernetesJobBuilder creates a builder using the default settings builders.
func newKubernetesJobBuilder() *kubernetesJobBuilder {
	return &kubernetesJobBuilder{
		addPodSettings:       podSettingsBuilder,
		addEndpointsSettings: endpointSettingsBuilder,
		addSelectorSettings:  selectorSettingsBuilder,
	}
}

// Build builds the prometheus targets corresponding to the Kubernetes configuration in input.
func (b *kubernetesJobBuilder) Build(i *Input) ([]JobOutput, error) {
	var jobs []JobOutput
	for _, k8sJob := range i.Kubernetes.Jobs {
		if !k8sJob.TargetDiscovery.Valid() {
			return nil, ErrInvalidK8sJobKinds
		}

		if err := b.checkJob(k8sJob); err != nil {
			return nil, err
		}

		if k8sJob.TargetDiscovery.Pod && b.addPodSettings != nil {
			job := b.buildJob(k8sJob, podKind)
			job = b.addPodSettings(job, k8sJob).WithExtraConfigs(k8sJob.JobInput)
			jobs = append(jobs, job)
		}

		if k8sJob.TargetDiscovery.Endpoints && b.addEndpointsSettings != nil {
			job := b.buildJob(k8sJob, endpointsKind)
			job = b.addEndpointsSettings(job, k8sJob).WithExtraConfigs(k8sJob.JobInput)
			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}

// buildJob creates a base JobOutput given the kubernetes settings and the target kind. It add
// the selector settings (if any) and builds the job name.
func (b *kubernetesJobBuilder) buildJob(k8sJob KubernetesJob, targetKind string) JobOutput {
	// build base job
	job := BuildJobOutput(k8sJob.JobInput)
	// apply selector rules if defined
	if b.addSelectorSettings != nil && k8sJob.Selector != nil {
		job = b.addSelectorSettings(job, k8sJob)
	}
	// build its name based on the prefix
	job.Job.JobName = b.buildJobName(k8sJob.JobNamePrefix, targetKind)
	return job
}

func (b *kubernetesJobBuilder) buildJobName(prefix string, kind string) string {
	return prefix + "-" + kind
}

func (b *kubernetesJobBuilder) checkJob(k8sJob KubernetesJob) error {
	if k8sJob.JobNamePrefix == "" {
		return ErrInvalidK8sJobPrefix
	}
	return nil
}
