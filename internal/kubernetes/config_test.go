package kubernetes_test

import (
	"testing"

	"github.com/newrelic/newrelic-prometheus-configurator/internal/kubernetes"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/scrapejob"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/sharding"
	"github.com/stretchr/testify/require"
)

func TestBuildFailWhen(t *testing.T) { //nolint: funlen
	t.Parallel()

	tests := []struct {
		name      string
		k8sConfig kubernetes.Config
		want      error
	}{
		{
			name: "JobNamePrefix is empty",
			k8sConfig: kubernetes.Config{
				K8sJobs: []kubernetes.K8sJob{
					{
						JobNamePrefix:   "",
						TargetDiscovery: kubernetes.TargetDiscovery{Pod: true},
					},
				},
			},
			want: kubernetes.ErrInvalidK8sJobPrefix,
		},
		{
			name: "All TargetKind are disabled",
			k8sConfig: kubernetes.Config{
				K8sJobs: []kubernetes.K8sJob{
					{
						JobNamePrefix: "test",
					},
				},
			},
			want: kubernetes.ErrInvalidK8sJobKinds,
		},
		{
			name: "skip_sharding flag is set",
			k8sConfig: kubernetes.Config{
				K8sJobs: []kubernetes.K8sJob{
					{
						ScrapeJob:       scrapejob.Job{SkipSharding: true},
						JobNamePrefix:   "test",
						TargetDiscovery: kubernetes.TargetDiscovery{Pod: true},
					},
				},
			},
			want: kubernetes.ErrInvalidSkipShardingFlag,
		},
		{
			name: "two labels only",
			k8sConfig: kubernetes.Config{
				K8sJobs: []kubernetes.K8sJob{
					{
						JobNamePrefix: "test-pod",
						TargetDiscovery: kubernetes.TargetDiscovery{
							Pod: true,
						},
					},
				},
				IntegrationFilter: kubernetes.IntegrationFilter{
					SourceLabels: []string{"label1", "label2"},
					AppValues:    []string{},
					Enabled:      boolPtr(true),
				},
			},
			want: kubernetes.ErrIntegrationFilterConfig,
		},
		{
			name: "two labels only, enable in job definition",
			k8sConfig: kubernetes.Config{
				K8sJobs: []kubernetes.K8sJob{
					{
						JobNamePrefix: "test-pod",
						TargetDiscovery: kubernetes.TargetDiscovery{
							Pod: true,
						},
						IntegrationFilter: kubernetes.IntegrationFilter{
							Enabled: boolPtr(true),
						},
					},
				},
				IntegrationFilter: kubernetes.IntegrationFilter{
					SourceLabels: []string{"label1", "label2"},
					AppValues:    []string{},
					Enabled:      boolPtr(false),
				},
			},
			want: kubernetes.ErrIntegrationFilterConfig,
		},
		{
			name: "two regexes only",
			k8sConfig: kubernetes.Config{
				K8sJobs: []kubernetes.K8sJob{
					{
						JobNamePrefix: "test-pod",
						TargetDiscovery: kubernetes.TargetDiscovery{
							Pod: true,
						},
					},
				},
				IntegrationFilter: kubernetes.IntegrationFilter{
					SourceLabels: []string{},
					AppValues:    []string{"regex1", "regex2"},
					Enabled:      boolPtr(true),
				},
			},
			want: kubernetes.ErrIntegrationFilterConfig,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := tt.k8sConfig.Build(sharding.Config{})
			require.ErrorIs(t, err, tt.want)
		})
	}
}

func TestBuildFilter(t *testing.T) { //nolint: funlen
	t.Parallel()

	annotationsFilter := kubernetes.Filter{
		Annotations: map[string]string{"prometheus.io/scrape": "true"},
	}

	combinedFilter := kubernetes.Filter{
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

	emptyLabelFilter := kubernetes.Filter{
		Labels: map[string]string{
			"check.if.present": "",
		},
	}

	type regexBySourceLabel map[string]string

	tests := []struct {
		name     string
		nrConfig kubernetes.Config
		want     regexBySourceLabel
	}{
		{
			name: "annotation pod filter",
			nrConfig: kubernetes.Config{
				K8sJobs: []kubernetes.K8sJob{
					{
						JobNamePrefix: "test-pod",
						TargetDiscovery: kubernetes.TargetDiscovery{
							Pod:    true,
							Filter: annotationsFilter,
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
			nrConfig: kubernetes.Config{
				K8sJobs: []kubernetes.K8sJob{
					{
						JobNamePrefix: "test-endpoints",
						TargetDiscovery: kubernetes.TargetDiscovery{
							Pod:    true,
							Filter: emptyLabelFilter,
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
			nrConfig: kubernetes.Config{
				K8sJobs: []kubernetes.K8sJob{
					{
						JobNamePrefix: "test-pod",
						TargetDiscovery: kubernetes.TargetDiscovery{
							Pod:    true,
							Filter: combinedFilter,
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
			nrConfig: kubernetes.Config{
				K8sJobs: []kubernetes.K8sJob{
					{
						JobNamePrefix: "test-endpoints",
						TargetDiscovery: kubernetes.TargetDiscovery{
							Endpoints: true,
							Filter:    annotationsFilter,
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
			nrConfig: kubernetes.Config{
				K8sJobs: []kubernetes.K8sJob{
					{
						JobNamePrefix: "test-endpoints",
						TargetDiscovery: kubernetes.TargetDiscovery{
							Endpoints: true,
							Filter:    emptyLabelFilter,
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
			nrConfig: kubernetes.Config{
				K8sJobs: []kubernetes.K8sJob{
					{
						JobNamePrefix: "test-pod",
						TargetDiscovery: kubernetes.TargetDiscovery{
							Endpoints: true,
							Filter:    combinedFilter,
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
			job, err := tt.nrConfig.Build(sharding.Config{})
			require.NoError(t, err)

			// tests should be independent and contain just one job entry.
			require.Len(t, job, 1)
			require.GreaterOrEqual(t, len(job[0].RelabelConfigs), 1)

			// we expect the filter relabel config as first one.
			actualRelabelConfig := job[0].RelabelConfigs[0]

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

func TestBuildIntegrationFilter(t *testing.T) { //nolint: funlen
	t.Parallel()

	annotationsFilter := kubernetes.Filter{
		Annotations: map[string]string{"prometheus.io/scrape": "true"},
	}

	type regexBySourceLabel map[string]string

	tests := []struct {
		name     string
		nrConfig kubernetes.Config
		want     *regexBySourceLabel
		len      int
	}{
		{
			name: "two labels two app values",
			nrConfig: kubernetes.Config{
				K8sJobs: []kubernetes.K8sJob{
					{
						JobNamePrefix: "test-pod",
						TargetDiscovery: kubernetes.TargetDiscovery{
							Pod:    true,
							Filter: annotationsFilter,
						},
					},
				},
				IntegrationFilter: kubernetes.IntegrationFilter{
					SourceLabels: []string{"label1", "label2"},
					AppValues:    []string{"regex1", "regex2"},
					Enabled:      boolPtr(true),
				},
			},
			want: &regexBySourceLabel{
				"__meta_kubernetes_pod_label_label1": ".*(?i)(regex1|regex2).*",
				"__meta_kubernetes_pod_label_label2": ".*(?i)(regex1|regex2).*",
			},
			len: 11,
		},
		{
			name: "two labels two app values from default even if disabled",
			nrConfig: kubernetes.Config{
				K8sJobs: []kubernetes.K8sJob{
					{
						JobNamePrefix: "test-pod",
						TargetDiscovery: kubernetes.TargetDiscovery{
							Pod:    true,
							Filter: annotationsFilter,
						},
						IntegrationFilter: kubernetes.IntegrationFilter{
							Enabled: boolPtr(true),
						},
					},
				},
				IntegrationFilter: kubernetes.IntegrationFilter{
					SourceLabels: []string{"label1", "label2"},
					AppValues:    []string{"regex1", "regex2"},
					Enabled:      boolPtr(false),
				},
			},
			want: &regexBySourceLabel{
				"__meta_kubernetes_pod_label_label1": ".*(?i)(regex1|regex2).*",
				"__meta_kubernetes_pod_label_label2": ".*(?i)(regex1|regex2).*",
			},
			len: 11,
		},
		{
			name: "two labels two regexes at different levels",
			nrConfig: kubernetes.Config{
				K8sJobs: []kubernetes.K8sJob{
					{
						JobNamePrefix: "test-pod",
						TargetDiscovery: kubernetes.TargetDiscovery{
							Pod:    true,
							Filter: annotationsFilter,
						},
						IntegrationFilter: kubernetes.IntegrationFilter{
							SourceLabels: []string{"label1", "label2"},
							AppValues:    []string{"regex1", "regex2"},
							Enabled:      boolPtr(true),
						},
					},
				},
				IntegrationFilter: kubernetes.IntegrationFilter{
					SourceLabels: []string{"different1", "different2"},
					AppValues:    []string{"regexDifferent1", "regexDifferent2"},
					Enabled:      boolPtr(true),
				},
			},
			want: &regexBySourceLabel{
				"__meta_kubernetes_pod_label_label1": ".*(?i)(regex1|regex2).*",
				"__meta_kubernetes_pod_label_label2": ".*(?i)(regex1|regex2).*",
			},
			len: 11,
		},
		{
			name: "no labels no regex but disabled at job label",
			nrConfig: kubernetes.Config{
				K8sJobs: []kubernetes.K8sJob{
					{
						JobNamePrefix: "test-pod",
						TargetDiscovery: kubernetes.TargetDiscovery{
							Pod:    true,
							Filter: annotationsFilter,
						},
						IntegrationFilter: kubernetes.IntegrationFilter{
							Enabled: boolPtr(false),
						},
					},
				},
				IntegrationFilter: kubernetes.IntegrationFilter{
					Enabled: boolPtr(true),
				},
			},
			len: 10,
		},
		{
			name: "No labels no regex but disabled at default label",
			nrConfig: kubernetes.Config{
				K8sJobs: []kubernetes.K8sJob{
					{
						JobNamePrefix: "test-pod",
						TargetDiscovery: kubernetes.TargetDiscovery{
							Pod:    true,
							Filter: annotationsFilter,
						},
					},
				},
				IntegrationFilter: kubernetes.IntegrationFilter{
					Enabled: boolPtr(false),
				},
			},
			len: 10,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			job, err := tt.nrConfig.Build(sharding.Config{})
			require.NoError(t, err)

			// we expect the filter relabel config as last one.
			l := len(job[0].RelabelConfigs)
			actualRelabelConfig := job[0].RelabelConfigs[l-1]

			require.Len(t, job[0].RelabelConfigs, tt.len)

			for _, actualSourceLabel := range actualRelabelConfig.SourceLabels {
				if tt.want != nil {
					regexBySourceLabel := *tt.want
					val, ok := regexBySourceLabel[actualSourceLabel]
					require.True(t, ok, "source label not expected: ", actualSourceLabel)
					require.Equal(t, actualRelabelConfig.Regex, val)
				}
			}
		})
	}
}

func TestBuildIntegrationFilterDifferentConfig(t *testing.T) {
	t.Parallel()

	annotationsFilter := kubernetes.Filter{
		Annotations: map[string]string{"prometheus.io/scrape": "true"},
	}

	nrConfig := kubernetes.Config{
		K8sJobs: []kubernetes.K8sJob{
			{
				JobNamePrefix: "job-with-default-filtering",
				TargetDiscovery: kubernetes.TargetDiscovery{
					Pod:    true,
					Filter: annotationsFilter,
				},
			},
			{
				JobNamePrefix: "job-with-custom-filtering",
				TargetDiscovery: kubernetes.TargetDiscovery{
					Pod:    true,
					Filter: annotationsFilter,
				},
				IntegrationFilter: kubernetes.IntegrationFilter{
					SourceLabels: []string{"different-label"},
					AppValues:    []string{"different-regex"},
				},
			},
			{
				JobNamePrefix: "job-with-no-filtering",
				TargetDiscovery: kubernetes.TargetDiscovery{
					Pod:    true,
					Filter: annotationsFilter,
				},
				IntegrationFilter: kubernetes.IntegrationFilter{
					Enabled: boolPtr(false),
				},
			},
		},
		IntegrationFilter: kubernetes.IntegrationFilter{
			SourceLabels: []string{"label1", "label2"},
			AppValues:    []string{"regex1", "regex2"},
			Enabled:      boolPtr(true),
		},
	}

	expectedJobSpecs := map[int]struct {
		want *map[string]string
		len  int
	}{
		0: {
			want: &map[string]string{
				"__meta_kubernetes_pod_label_label1": ".*(?i)(regex1|regex2).*",
				"__meta_kubernetes_pod_label_label2": ".*(?i)(regex1|regex2).*",
			},
			len: 11,
		},
		1: {
			want: &map[string]string{
				"__meta_kubernetes_pod_label_different_label": ".*(?i)(different-regex).*",
			},
			len: 11,
		},
		2: {
			len: 10,
		},
	}

	job, err := nrConfig.Build(sharding.Config{})
	require.NoError(t, err)

	for indexJob, js := range expectedJobSpecs {
		// we expect the filter relabel config as last one.
		l := len(job[indexJob].RelabelConfigs)
		actualRelabelConfig := job[indexJob].RelabelConfigs[l-1]

		require.Len(t, job[indexJob].RelabelConfigs, js.len)

		for _, actualSourceLabel := range actualRelabelConfig.SourceLabels {
			if js.want != nil {
				regexBySourceLabel := *js.want
				val, ok := regexBySourceLabel[actualSourceLabel]
				require.True(t, ok, "source label not expected: ", actualSourceLabel)
				require.Equal(t, actualRelabelConfig.Regex, val)
			}
		}
	}
}

func boolPtr(b bool) *bool {
	return &b
}
