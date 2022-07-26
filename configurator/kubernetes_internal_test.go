package configurator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKubernetesTargetBuilder_KubernetesNotEnabled(t *testing.T) {
	t.Parallel()

	i := &Input{Kubernetes: KubernetesInput{Enabled: false}}
	builder := &kubernetesTargetBuilder{}
	assert.Nil(t, builder.Build(i))
}

//nolint: funlen
func TestKubernetesTargetBuilder(t *testing.T) {
	t.Parallel()

	podSettingsMock := func(tg TargetJobOutput, job KubernetesJob) TargetJobOutput {
		if tg.StaticConfigs[0].Labels == nil {
			tg.StaticConfigs[0].Labels = map[string]string{}
		}
		tg.StaticConfigs[0].Labels["pods"] = "pods"
		return tg
	}

	endpointSettingsMock := func(tg TargetJobOutput, job KubernetesJob) TargetJobOutput {
		if tg.StaticConfigs[0].Labels == nil {
			tg.StaticConfigs[0].Labels = map[string]string{}
		}
		tg.StaticConfigs[0].Labels["endpoints"] = "endpoints"
		return tg
	}

	selectorSettingsMock := func(tg TargetJobOutput, job KubernetesJob) TargetJobOutput {
		if tg.StaticConfigs[0].Labels == nil {
			tg.StaticConfigs[0].Labels = map[string]string{}
		}
		tg.StaticConfigs[0].Labels["selector"] = "selector"
		return tg
	}

	cases := []struct {
		Name     string
		Input    *Input
		Expected []TargetJobOutput
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
							TargetJobInput: TargetJobInput{
								TargetJob: TargetJob{
									JobName: "job-pods",
								},
							},
							TargetKind: []string{"pods"},
						},
					},
				},
			},
			Expected: []TargetJobOutput{
				{
					TargetJob:     TargetJob{JobName: "job-pods"},
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
							TargetJobInput: TargetJobInput{
								TargetJob: TargetJob{
									JobName: "job-endpoints",
								},
							},
							TargetKind: []string{"endpoints"},
						},
					},
				},
			},
			Expected: []TargetJobOutput{
				{
					TargetJob:     TargetJob{JobName: "job-endpoints"},
					StaticConfigs: []StaticConfig{{Labels: map[string]string{"endpoints": "endpoints"}}},
				},
			},
		},
		{
			Name: "Selector defined",
			Input: &Input{
				Kubernetes: KubernetesInput{
					Enabled: true,
					Jobs: []KubernetesJob{
						{
							TargetJobInput: TargetJobInput{
								TargetJob: TargetJob{
									JobName: "job-selector",
								},
							},
							Selector: &KubernetesSelector{},
						},
					},
				},
			},
			Expected: []TargetJobOutput{
				{
					TargetJob:     TargetJob{JobName: "job-selector"},
					StaticConfigs: []StaticConfig{{Labels: map[string]string{"selector": "selector"}}},
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
							TargetJobInput: TargetJobInput{
								TargetJob: TargetJob{
									JobName: "job-complete",
								},
							},
							TargetKind: []string{"pods", "endpoints"},
							Selector:   &KubernetesSelector{},
						},
					},
				},
			},
			Expected: []TargetJobOutput{
				{
					TargetJob: TargetJob{JobName: "job-complete"},
					StaticConfigs: []StaticConfig{
						{
							Labels: map[string]string{"selector": "selector", "pods": "pods", "endpoints": "endpoints"},
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

			builder := newKubernetesTargetBuilder(podSettingsMock, endpointSettingsMock, selectorSettingsMock)
			targets := builder.Build((c.Input))
			assert.Equal(t, c.Expected, targets)
		})
	}
}
