package configurator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKubernetesJobBuilder_InvalidSettings(t *testing.T) {
	t.Parallel()

	builder := &kubernetesJobBuilder{}

	cases := []struct {
		Name  string
		Input *Input
	}{
		{
			Name: "No kind defined",
			Input: &Input{
				Kubernetes: KubernetesInput{
					Jobs: []KubernetesJob{
						{JobNamePrefix: "job"},
					},
				},
			},
		},
		{
			Name: "No prefix defined",
			Input: &Input{
				Kubernetes: KubernetesInput{
					Jobs: []KubernetesJob{
						{TargetDiscovery: KubernetesTargetDiscovery{Pod: true}},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		c := tc
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			_, err := builder.Build(c.Input)
			require.Error(t, err)
		})
	}
}

//nolint: funlen
func TestKubernetesJobBuilder(t *testing.T) {
	t.Parallel()

	withStaticConfigReady := func(settingsBuilder kubernetesSettingsBuilder) kubernetesSettingsBuilder {
		return func(j JobOutput, k8sJob KubernetesJob) JobOutput {
			if len(j.StaticConfigs) == 0 || j.StaticConfigs[0].Labels == nil {
				j.StaticConfigs = []StaticConfig{
					{Labels: map[string]string{}},
				}
			}
			return settingsBuilder(j, k8sJob)
		}
	}

	podSettingsMock := func(j JobOutput, _ KubernetesJob) JobOutput {
		j.StaticConfigs[0].Labels["pods"] = "pods"
		return j
	}

	endpointSettingsMock := func(j JobOutput, _ KubernetesJob) JobOutput {
		j.StaticConfigs[0].Labels["endpoints"] = "endpoints"
		return j
	}

	selectorSettingsMock := func(j JobOutput, _ KubernetesJob) JobOutput {
		j.StaticConfigs[0].Labels["selector"] = "selector"
		return j
	}

	builder := &kubernetesJobBuilder{
		addPodSettings:       withStaticConfigReady(podSettingsMock),
		addEndpointsSettings: withStaticConfigReady(endpointSettingsMock),
		addSelectorSettings:  withStaticConfigReady(selectorSettingsMock),
	}

	cases := []struct {
		Name     string
		Input    *Input
		Expected []JobOutput
	}{
		{
			Name:     "Kubernetes not defined",
			Input:    &Input{},
			Expected: nil,
		},
		{
			Name:     "Kubernetes empty",
			Input:    &Input{Kubernetes: KubernetesInput{}},
			Expected: nil,
		},
		{
			Name: "Kind pods",
			Input: &Input{
				Kubernetes: KubernetesInput{
					Jobs: []KubernetesJob{
						{
							JobNamePrefix:   "job",
							TargetDiscovery: KubernetesTargetDiscovery{Pod: true},
						},
					},
				},
			},
			Expected: []JobOutput{
				{
					Job:           Job{JobName: "job-pod"},
					StaticConfigs: []StaticConfig{{Labels: map[string]string{"pods": "pods"}}},
				},
			},
		},
		{
			Name: "Kind endpoints",
			Input: &Input{
				Kubernetes: KubernetesInput{
					Jobs: []KubernetesJob{
						{
							JobNamePrefix:   "job",
							TargetDiscovery: KubernetesTargetDiscovery{Endpoints: true},
						},
					},
				},
			},
			Expected: []JobOutput{
				{
					Job:           Job{JobName: "job-endpoints"},
					StaticConfigs: []StaticConfig{{Labels: map[string]string{"endpoints": "endpoints"}}},
				},
			},
		},
		{
			Name: "Selector defined and pod",
			Input: &Input{
				Kubernetes: KubernetesInput{
					Jobs: []KubernetesJob{
						{
							JobNamePrefix:   "job",
							TargetDiscovery: KubernetesTargetDiscovery{Pod: true},
							Selector:        &KubernetesSelector{},
						},
					},
				},
			},
			Expected: []JobOutput{
				{
					Job:           Job{JobName: "job-pod"},
					StaticConfigs: []StaticConfig{{Labels: map[string]string{"selector": "selector", "pods": "pods"}}},
				},
			},
		},
		{
			Name: "Selector defined and endpoints",
			Input: &Input{
				Kubernetes: KubernetesInput{
					Jobs: []KubernetesJob{
						{
							JobNamePrefix:   "job",
							TargetDiscovery: KubernetesTargetDiscovery{Endpoints: true},
							Selector:        &KubernetesSelector{},
						},
					},
				},
			},
			Expected: []JobOutput{
				{
					Job:           Job{JobName: "job-endpoints"},
					StaticConfigs: []StaticConfig{{Labels: map[string]string{"selector": "selector", "endpoints": "endpoints"}}},
				},
			},
		},
		{
			Name: "Pods, endpoints and selector defined",
			Input: &Input{
				Kubernetes: KubernetesInput{
					Jobs: []KubernetesJob{
						{
							JobNamePrefix:   "job",
							TargetDiscovery: KubernetesTargetDiscovery{Pod: true, Endpoints: true},
							Selector:        &KubernetesSelector{},
						},
					},
				},
			},
			Expected: []JobOutput{
				{
					Job: Job{JobName: "job-pod"},
					StaticConfigs: []StaticConfig{
						{
							Labels: map[string]string{"selector": "selector", "pods": "pods"},
						},
					},
				},
				{
					Job: Job{JobName: "job-endpoints"},
					StaticConfigs: []StaticConfig{
						{
							Labels: map[string]string{"selector": "selector", "endpoints": "endpoints"},
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			targets, err := builder.Build(c.Input)
			require.NoError(t, err)
			assert.Equal(t, c.Expected, targets)
		})
	}
}
