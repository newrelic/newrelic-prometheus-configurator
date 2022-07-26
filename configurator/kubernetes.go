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
	TargetJobInput `yaml:",inline"`

	JobNamePrefix string              `yaml:"job_name_prefix"`
	Selector      *KubernetesSelector `yaml:"selector,omitempty"`
	TargetKinds   KubernetesJobKind   `yaml:"target_kinds"`
}

type KubernetesJobKind struct {
	Pods      bool `yaml:"pods"`
	Endpoints bool `yaml:"endpoints"`
}

// Valid returns true when the defined configuration is valid.
func (k *KubernetesJobKind) Valid() bool {
	return k.Pods || k.Endpoints
}

// KubernetesSettingsBuilders defines a functions which updates and returns a `TargetJobOutput` with specific settings
// added (considering the specified `KubernetesJob`).
type kubernetesSettingsBuilder func(targetJob TargetJobOutput, k8sJob KubernetesJob) TargetJobOutput

// kubernetesTargetBuilder holds the the specific settings to add to a TargetJobOutput given the corresponding
// KubernetesJob definition.
type kubernetesTargetBuilder struct {
	setPodRules      kubernetesSettingsBuilder
	setEndpointRules kubernetesSettingsBuilder
	setSelectorRules kubernetesSettingsBuilder
}

func newKubernetesTargetBuilder(pod, endpoint, selector kubernetesSettingsBuilder) *kubernetesTargetBuilder {
	return &kubernetesTargetBuilder{
		setPodRules:      pod,
		setEndpointRules: endpoint,
		setSelectorRules: selector,
	}
}

// BuildKubernetesTargets builds the prometheus targets corresponding to the Kubernetes configuration in input.
func (b *kubernetesTargetBuilder) Build(i *Input) ([]TargetJobOutput, error) {
	if !i.Kubernetes.Enabled {
		return nil, nil
	}
	targetJobs := make([]TargetJobOutput, 0, len(i.Kubernetes.Jobs))
	for _, k8sJob := range i.Kubernetes.Jobs {
		if !k8sJob.TargetKinds.Valid() {
			return nil, ErrInvalidK8sJobKinds
		}

		if k8sJob.TargetKinds.Pods && b.setPodRules != nil {
			targetJob := b.buildTargetJob(k8sJob, podKind)
			targetJob = b.setPodRules(targetJob, k8sJob)
			targetJobs = append(targetJobs, targetJob)
		}

		if k8sJob.TargetKinds.Endpoints && b.setEndpointRules != nil {
			targetJob := b.buildTargetJob(k8sJob, endpointKind)
			targetJob = b.setEndpointRules(targetJob, k8sJob)
			targetJobs = append(targetJobs, targetJob)
		}
	}

	return targetJobs, nil
}

func (b *kubernetesTargetBuilder) buildTargetJob(job KubernetesJob, kind string) TargetJobOutput {
	// build base target job
	targetJob := BuildTargetJob(job.TargetJobInput)
	// apply selector rules if defined
	if b.setSelectorRules != nil && job.Selector != nil {
		targetJob = b.setSelectorRules(targetJob, job)
	}
	// build its name based on the prefix
	targetJob.TargetJob.JobName = b.buildJobName(job.JobNamePrefix, kind)
	return targetJob
}

func (b *kubernetesTargetBuilder) buildJobName(prefix string, kind string) string {
	return prefix + "-" + kind
}
