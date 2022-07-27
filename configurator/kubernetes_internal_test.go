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

func TestKubernetesTargetBuilder_InvalidSettings(t *testing.T) {
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
func TestKubernetesTargetBuilder(t *testing.T) {
	t.Parallel()

	podSettingsMock := func(tg JobOutput, job KubernetesJob) JobOutput {
		if tg.StaticConfigs[0].Labels == nil {
			tg.StaticConfigs[0].Labels = map[string]string{}
		}
		tg.StaticConfigs[0].Labels["pods"] = "pods"
		return tg
	}

	endpointSettingsMock := func(tg JobOutput, job KubernetesJob) JobOutput {
		if tg.StaticConfigs[0].Labels == nil {
			tg.StaticConfigs[0].Labels = map[string]string{}
		}
		tg.StaticConfigs[0].Labels["endpoints"] = "endpoints"
		return tg
	}

	selectorSettingsMock := func(tg JobOutput, job KubernetesJob) JobOutput {
		if tg.StaticConfigs[0].Labels == nil {
			tg.StaticConfigs[0].Labels = map[string]string{}
		}
		tg.StaticConfigs[0].Labels["selector"] = "selector"
		return tg
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
							TargetKinds:   KubernetesTargetKind{Pods: true},
						},
					},
				},
			},
			Expected: []JobOutput{
				{
					Job:           Job{JobName: "job-pods"},
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
							TargetKinds:   KubernetesTargetKind{Endpoints: true},
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
					Enabled: true,
					Jobs: []KubernetesJob{
						{
							JobNamePrefix: "job",
							TargetKinds:   KubernetesTargetKind{Pods: true},
							Selector:      &KubernetesSelector{},
						},
					},
				},
			},
			Expected: []JobOutput{
				{
					Job:           Job{JobName: "job-pods"},
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
							TargetKinds:   KubernetesTargetKind{Endpoints: true},
							Selector:      &KubernetesSelector{},
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
					Enabled: true,
					Jobs: []KubernetesJob{
						{
							JobNamePrefix: "job",
							TargetKinds:   KubernetesTargetKind{Pods: true, Endpoints: true},
							Selector:      &KubernetesSelector{},
						},
					},
				},
			},
			Expected: []JobOutput{
				{
					Job: Job{JobName: "job-pods"},
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

			builder := &kubernetesJobBuilder{
				addPodSettings:      podSettingsMock,
				addEndpointSettings: endpointSettingsMock,
				addSelectorSettings: selectorSettingsMock,
			}
			targets, err := builder.Build((c.Input))
			require.NoError(t, err)
			assert.Equal(t, c.Expected, targets)
		})
	}
}
