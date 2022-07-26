package configurator

const (
	PodKind      = "pods"
	EndpointKind = "endpoints"
)

// KubernetesInput defines all fields to set up prometheus.
type KubernetesInput struct {
	Enabled bool            `yaml:"enabled"`
	Jobs    []KubernetesJob `yaml:"jobs"`
}

// KubernetesJob holds the configuration which will parsed to a prometheus scrape job including the
// specific rules needed.
type KubernetesJob struct {
	InputJob `yaml:",inline"`

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
	return j.checkKind(PodKind)
}

func (j *KubernetesJob) Endpoints() bool {
	return j.checkKind(EndpointKind)
}

// KubernetesSelector defines the field needed to provided filtering capabilities to a kubernetes scrape job.
type KubernetesSelector struct {
	// TODO: define selector when this is implemented
}

// BuildKubernetesTargets builds the prometheus targets corresponding to the Kubernetes configuration in input.
func BuildKubernetesTargets(i *Input) []TargetJobOutput {
	if !i.Kubernetes.Enabled {
		return nil
	}
	targetJobs := make([]TargetJobOutput, 0, len(i.Kubernetes.Jobs))
	for _, k8sJob := range i.Kubernetes.Jobs {
		tg := BuildTargetJob(k8sJob.InputJob)
		if k8sJob.Pods() {
			tg = SetupPodsRules(tg)
		}
		if k8sJob.Endpoints() {
			tg = SetupPodsRules(tg)
		}
		if k8sJob.Selector != nil {
			tg = SetupSelectorRules(tg, k8sJob.Selector)
		}
		targetJobs = append(targetJobs, tg)
	}
	return targetJobs
}

func SetupPodsRules(tg TargetJobOutput) TargetJobOutput {
	return tg // TODO: implement it!
}

func SetupServicesRules(tg TargetJobOutput) TargetJobOutput {
	return tg // TODO: implement it!
}

func SetupSelectorRules(tg TargetJobOutput, selector *KubernetesSelector) TargetJobOutput {
	return tg // TODO: implement it!
}
