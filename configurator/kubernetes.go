package configurator

import "errors"

const (
	podKind      = "pods"
	endpointKind = "endpoints"
)

var ErrInvalidK8sJobKinds = errors.New("at least one kind should be set in target_kinds field")

// KubernetesInput defines all fields to set up prometheus.
type KubernetesInput struct {
	Enabled bool            `yaml:"enabled"`
	Jobs    []KubernetesJob `yaml:"jobs"`
}

// KubernetesJob holds the configuration which will parsed to a prometheus scrape job including the
// specific rules needed.
type KubernetesJob struct {
	JobInput `yaml:",inline"`

	JobNamePrefix string               `yaml:"job_name_prefix"`
	Selector      *KubernetesSelector  `yaml:"selector,omitempty"`
	TargetKinds   KubernetesTargetKind `yaml:"target_kinds"`
}

type KubernetesTargetKind struct {
	Pods      bool `yaml:"pods"`
	Endpoints bool `yaml:"endpoints"`
}

// Valid returns true when the defined configuration is valid.
func (k *KubernetesTargetKind) Valid() bool {
	return k.Pods || k.Endpoints
}

// KubernetesSettingsBuilders defines a functions which updates and returns a `TargetJobOutput` with specific settings
// added (considering the specified `KubernetesJob`).
type kubernetesSettingsBuilder func(job JobOutput, k8sJob KubernetesJob) JobOutput

// kubernetesJobBuilder holds the the specific settings to add to a TargetJobOutput given the corresponding
// KubernetesJob definition.
type kubernetesJobBuilder struct {
	setPodRules      kubernetesSettingsBuilder
	setEndpointRules kubernetesSettingsBuilder
	setSelectorRules kubernetesSettingsBuilder
}

func newKubernetesJobBuilder(pod, endpoint, selector kubernetesSettingsBuilder) *kubernetesJobBuilder {
	return &kubernetesJobBuilder{
		setPodRules:      pod,
		setEndpointRules: endpoint,
		setSelectorRules: selector,
	}
}

// BuildKubernetesTargets builds the prometheus targets corresponding to the Kubernetes configuration in input.
func (b *kubernetesJobBuilder) Build(i *Input) ([]JobOutput, error) {
	if !i.Kubernetes.Enabled {
		return nil, nil
	}
	jobs := make([]JobOutput, 0, len(i.Kubernetes.Jobs))
	for _, k8sJob := range i.Kubernetes.Jobs {
		if !k8sJob.TargetKinds.Valid() {
			return nil, ErrInvalidK8sJobKinds
		}

		if k8sJob.TargetKinds.Pods && b.setPodRules != nil {
			job := b.buildJob(k8sJob, podKind)
			job = b.setPodRules(job, k8sJob)
			jobs = append(jobs, job)
		}

		if k8sJob.TargetKinds.Endpoints && b.setEndpointRules != nil {
			job := b.buildJob(k8sJob, endpointKind)
			job = b.setEndpointRules(job, k8sJob)
			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}

func (b *kubernetesJobBuilder) buildJob(k8sJob KubernetesJob, kind string) JobOutput {
	// build base job
	job := BuildJobOutput(k8sJob.JobInput)
	// apply selector rules if defined
	if b.setSelectorRules != nil && k8sJob.Selector != nil {
		job = b.setSelectorRules(job, k8sJob)
	}
	// build its name based on the prefix
	job.Job.JobName = b.buildJobName(k8sJob.JobNamePrefix, kind)
	return job
}

func (b *kubernetesJobBuilder) buildJobName(prefix string, kind string) string {
	return prefix + "-" + kind
}
