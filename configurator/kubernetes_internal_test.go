package configurator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKubernetesTargetBuilder_KubernetesNotEnabled(t *testing.T) {
	t.Parallel()

	i := &Input{Kubernetes: KubernetesInput{Enabled: false}}
	builder := &kubernetesJobBuilder{}
	result, err := builder.Build(i)
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestKubernetesJobBuilder_InvalidSettings(t *testing.T) {
	t.Parallel()

	i := &Input{
		Kubernetes: KubernetesInput{
			Enabled: true,
			Jobs: []KubernetesJob{
				{JobNamePrefix: "job"}, // No kind defined.
			},
		},
	}
	builder := &kubernetesJobBuilder{}
	_, err := builder.Build(i)
	require.Error(t, err)
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
		addPodSettings:      withStaticConfigReady(podSettingsMock),
		addEndpointSettings: withStaticConfigReady(endpointSettingsMock),
		addSelectorSettings: withStaticConfigReady(selectorSettingsMock),
	}

	cases := []struct {
		Name     string
		Input    *Input
		Expected []JobOutput
	}{
		{
			Name:     "Kubernetes not enabled",
			Input:    &Input{Kubernetes: KubernetesInput{Enabled: false}},
			Expected: nil,
		},
		{
			Name: "Kind pods",
			Input: &Input{
				Kubernetes: KubernetesInput{
					Enabled: true,
					Jobs: []KubernetesJob{
						{
							JobNamePrefix: "job",
							TargetKind:    KubernetesTargetKind{Pod: true},
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
					Enabled: true,
					Jobs: []KubernetesJob{
						{
							JobNamePrefix: "job",
							TargetKind:    KubernetesTargetKind{Endpoint: true},
						},
					},
				},
			},
			Expected: []JobOutput{
				{
					Job:           Job{JobName: "job-endpoint"},
					StaticConfigs: []StaticConfig{{Labels: map[string]string{"endpoints": "endpoints"}}},
				},
			},
		},
		{
			Name: "Selector defined and pod",
			Input: &Input{
				Kubernetes: KubernetesInput{
					Enabled: true,
					Jobs: []KubernetesJob{
						{
							JobNamePrefix: "job",
							TargetKind:    KubernetesTargetKind{Pod: true},
							Selector:      &KubernetesSelector{},
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
					Enabled: true,
					Jobs: []KubernetesJob{
						{
							JobNamePrefix: "job",
							TargetKind:    KubernetesTargetKind{Endpoint: true},
							Selector:      &KubernetesSelector{},
						},
					},
				},
			},
			Expected: []JobOutput{
				{
					Job:           Job{JobName: "job-endpoint"},
					StaticConfigs: []StaticConfig{{Labels: map[string]string{"selector": "selector", "endpoints": "endpoints"}}},
				},
			},
		},
		{
			Name: "Pods, endpoints and selector defined",
			Input: &Input{
				Kubernetes: KubernetesInput{
					Enabled: true,
					Jobs: []KubernetesJob{
						{
							JobNamePrefix: "job",
							TargetKind:    KubernetesTargetKind{Pod: true, Endpoint: true},
							Selector:      &KubernetesSelector{},
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
					Job: Job{JobName: "job-endpoint"},
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
