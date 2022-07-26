package configurator

const (
	podKind      = "pods"
	endpointKind = "endpoints"
)

// KubernetesInput defines all fields to set up prometheus.
type KubernetesInput struct {
	Enabled bool            `yaml:"enabled"`
	Jobs    []KubernetesJob `yaml:"jobs"`
}

// KubernetesJob holds the configuration which will parsed to a prometheus scrape job including the
// specific rules needed.
type KubernetesJob struct {
	TargetJobInput `yaml:",inline"`

	Selector *KubernetesSelector `yaml:"selector,omitempty"`
	// TargetKind currently supports 'pods' and 'services'.
	TargetKind []string `yaml:"target_kind"`
}

func (j *KubernetesJob) checkKind(kind string) bool {
	for _, tk := range j.TargetKind {
		if tk == kind {
			return true
		}
	}
	return false
}

func (j *KubernetesJob) Pods() bool {
	return j.checkKind(podKind)
}

func (j *KubernetesJob) Endpoints() bool {
	return j.checkKind(endpointKind)
}

// KubernetesSettingsBuilders defines a functions which returns a copy of the provided `TargetJobOutput` with specific
// added (considering the specified `*KubernetesJob`).
type kubernetesSettingsBuilder func(tg TargetJobOutput, k8sJob KubernetesJob) TargetJobOutput

// kubernetesTargetBuilder holds the the specific settings to add to a TargetJobOutput given the corresponding
// KubernetesJob definition.
type kubernetesTargetBuilder struct {
	pod      kubernetesSettingsBuilder
	endpoint kubernetesSettingsBuilder
	selector kubernetesSettingsBuilder
}

func newKubernetesTargetBuilder(pod, endpoint, selector kubernetesSettingsBuilder) *kubernetesTargetBuilder {
	return &kubernetesTargetBuilder{
		pod:      pod,
		endpoint: endpoint,
		selector: selector,
	}
}

// BuildKubernetesTargets builds the prometheus targets corresponding to the Kubernetes configuration in input.
func (b *kubernetesTargetBuilder) Build(i *Input) []TargetJobOutput {
	if !i.Kubernetes.Enabled {
		return nil
	}
	targetJobs := make([]TargetJobOutput, 0, len(i.Kubernetes.Jobs))
	for _, k8sJob := range i.Kubernetes.Jobs {
		tg := BuildTargetJob(k8sJob.TargetJobInput)
		if b.pod != nil && k8sJob.Pods() {
			tg = b.pod(tg, k8sJob)
		}
		if b.endpoint != nil && k8sJob.Endpoints() {
			tg = b.endpoint(tg, k8sJob)
		}
		if b.selector != nil && k8sJob.Selector != nil {
			tg = b.selector(tg, k8sJob)
		}
		targetJobs = append(targetJobs, tg)
	}
	return targetJobs
}
