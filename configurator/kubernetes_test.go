package configurator_test

import (
	"testing"

	"github.com/newrelic-forks/newrelic-prometheus/configurator"
	"github.com/stretchr/testify/require"
)

func TestBuildFailWhen(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		k8sConfig configurator.KubernetesInput
		want      error
	}{
		{
			name: "JobNamePrefix is empty",
			k8sConfig: configurator.KubernetesInput{
				Jobs: []configurator.KubernetesJob{
					{
						JobNamePrefix:   "",
						TargetDiscovery: configurator.KubernetesTargetDiscovery{Pod: true},
					},
				},
			},
			want: configurator.ErrInvalidK8sJobPrefix,
		},
		{
			name: "All TargetKind are disabled",
			k8sConfig: configurator.KubernetesInput{
				Jobs: []configurator.KubernetesJob{
					{
						JobNamePrefix: "test",
					},
				},
			},
			want: configurator.ErrInvalidK8sJobKinds,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := configurator.NewKubernetesJobBuilder().Build(&configurator.Input{
				Kubernetes: tt.k8sConfig,
			})
			require.ErrorIs(t, err, tt.want)
		})
	}
}
