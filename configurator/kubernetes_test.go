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

// nolint:funlen
func TestBuildFilter(t *testing.T) {
	t.Parallel()

	annotationsFilter := &configurator.Filter{
		Annotations: map[string]string{"prometheus.io/scrape": "true"},
	}

	combinedFilter := &configurator.Filter{
		Annotations: map[string]string{
			"prometheus.io/scrape":     "true",
			"extra.special.annotation": "yes",
			"empty":                    "",
		},
		Labels: map[string]string{
			"k8s.io/app":              "(foo|bar)",
			"my.custom.authorization": "my-auth",
		},
	}

	emptyLabelFilter := &configurator.Filter{
		Labels: map[string]string{
			"check.if.present": "",
		},
	}

	type args struct {
		Input configurator.Input
	}

	type regexBySourceLabel map[string]string

	tests := []struct {
		name string
		args args
		want regexBySourceLabel
	}{
		{
			name: "annotation pod filter",
			args: args{
				Input: configurator.Input{
					Kubernetes: configurator.KubernetesInput{
						Jobs: []configurator.KubernetesJob{
							{
								JobNamePrefix: "test-pod",
								TargetDiscovery: configurator.KubernetesTargetDiscovery{
									Pod:    true,
									Filter: annotationsFilter,
								},
							},
						},
					},
				},
			},
			want: regexBySourceLabel{
				"__meta_kubernetes_pod_annotation_prometheus_io_scrape": "true",
			},
		},
		{
			name: "check pod label is present",
			args: args{
				Input: configurator.Input{
					Kubernetes: configurator.KubernetesInput{
						Jobs: []configurator.KubernetesJob{
							{
								JobNamePrefix: "test-endpoints",
								TargetDiscovery: configurator.KubernetesTargetDiscovery{
									Pod:    true,
									Filter: emptyLabelFilter,
								},
							},
						},
					},
				},
			},
			want: regexBySourceLabel{
				"__meta_kubernetes_pod_labelpresent_check_if_present": "true",
			},
		},
		{
			name: "combined pod filter",
			args: args{
				Input: configurator.Input{
					Kubernetes: configurator.KubernetesInput{
						Jobs: []configurator.KubernetesJob{
							{
								JobNamePrefix: "test-pod",
								TargetDiscovery: configurator.KubernetesTargetDiscovery{
									Pod:    true,
									Filter: combinedFilter,
								},
							},
						},
					},
				},
			},
			want: regexBySourceLabel{
				"__meta_kubernetes_pod_annotation_prometheus_io_scrape":     "true",
				"__meta_kubernetes_pod_annotation_extra_special_annotation": "yes",
				"__meta_kubernetes_pod_annotationpresent_empty":             "true",
				"__meta_kubernetes_pod_label_k8s_io_app":                    "(foo|bar)",
				"__meta_kubernetes_pod_label_my_custom_authorization":       "my-auth",
			},
		},
		{
			name: "annotation service-endpoints filter",
			args: args{
				Input: configurator.Input{
					Kubernetes: configurator.KubernetesInput{
						Jobs: []configurator.KubernetesJob{
							{
								JobNamePrefix: "test-endpoints",
								TargetDiscovery: configurator.KubernetesTargetDiscovery{
									Endpoints: true,
									Filter:    annotationsFilter,
								},
							},
						},
					},
				},
			},
			want: regexBySourceLabel{
				"__meta_kubernetes_service_annotation_prometheus_io_scrape": "true",
			},
		},
		{
			name: "check service-endpoints label is present",
			args: args{
				Input: configurator.Input{
					Kubernetes: configurator.KubernetesInput{
						Jobs: []configurator.KubernetesJob{
							{
								JobNamePrefix: "test-endpoints",
								TargetDiscovery: configurator.KubernetesTargetDiscovery{
									Endpoints: true,
									Filter:    emptyLabelFilter,
								},
							},
						},
					},
				},
			},
			want: regexBySourceLabel{
				"__meta_kubernetes_service_labelpresent_check_if_present": "true",
			},
		},
		{
			name: "combined service-endpoints filter",
			args: args{
				Input: configurator.Input{
					Kubernetes: configurator.KubernetesInput{
						Jobs: []configurator.KubernetesJob{
							{
								JobNamePrefix: "test-pod",
								TargetDiscovery: configurator.KubernetesTargetDiscovery{
									Endpoints: true,
									Filter:    combinedFilter,
								},
							},
						},
					},
				},
			},
			want: regexBySourceLabel{
				"__meta_kubernetes_service_annotation_prometheus_io_scrape":     "true",
				"__meta_kubernetes_service_annotation_extra_special_annotation": "yes",
				"__meta_kubernetes_service_annotationpresent_empty":             "true",
				"__meta_kubernetes_service_label_k8s_io_app":                    "(foo|bar)",
				"__meta_kubernetes_service_label_my_custom_authorization":       "my-auth",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			job, err := configurator.NewKubernetesJobBuilder().Build(&tt.args.Input)
			require.NoError(t, err)

			// tests should be independent and contain just one job entry.
			require.Len(t, job, 1)
			require.GreaterOrEqual(t, len(job[0].RelabelConfigs), 1)

			// we expect the filter relabel config as first one.
			actualRelabelConfig, ok := job[0].RelabelConfigs[0].(configurator.RelabelConfig)
			require.True(t, ok)

			require.Equal(t, len(tt.want), len(actualRelabelConfig.SourceLabels))

			expectedRegex := ""

			// Order of source labels and regex is not guaranteed since they are created from
			// a map in the filter.
			// In order to test we also build an expected map and checks that all sourceLabels exist.
			for i, actualSourceLabel := range actualRelabelConfig.SourceLabels {
				val, ok := tt.want[actualSourceLabel]
				require.True(t, ok, "source label not expected: ", actualSourceLabel)

				// Since the regex depends on the position of the source labels in the array, we need to
				// build the expected one with the same order as the actual.
				expectedRegex += val

				// avoid concatenating last separator
				if i != len(actualRelabelConfig.SourceLabels)-1 {
					expectedRegex += ";"
				}
			}

			require.Equal(t, actualRelabelConfig.Regex, expectedRegex)
		})
	}
}
