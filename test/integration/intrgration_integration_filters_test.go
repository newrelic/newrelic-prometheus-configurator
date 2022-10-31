//go:build integration_test

package integration

import (
	"fmt"
	"testing"
)

func Test_IntegrationFilter(t *testing.T) {
	t.Parallel()

	testsCases := []struct {
		name string

		// When matchPod/Endpoint is added, the test will generate a pod/endpoint
		// with specified metadata and check if is an Active Target.
		matchPod       *metadata
		matchEndpoints *metadata

		// When dropPod/Endpoint is added, the test will generate a pod/endpoint
		// with specified metadata and check if is a Dropped Target.
		dropPod       *metadata
		dropEndpoints *metadata
	}{
		{
			name: "matching labels",
			matchPod: &metadata{
				annotations: map[string]string{"prometheus.io/scrape": "true"},
				labels:      map[string]string{"app.kubernetes.io/name": "asdTesT1asd"},
			},
			matchEndpoints: &metadata{
				annotations: map[string]string{"prometheus.io/scrape": "true"},
				labels:      map[string]string{"app.kubernetes.io/name": "asdTesT1asd"},
			},
		},
		{
			name: "matching labels",
			matchPod: &metadata{
				annotations: map[string]string{"prometheus.io/scrape": "true"},
				labels:      map[string]string{"something-different2": "asdtest2"},
			},
			matchEndpoints: &metadata{
				annotations: map[string]string{"prometheus.io/scrape": "true"},
				labels:      map[string]string{"something-different2": "asdtest2"},
			},
		},
		{
			name: "matching annotation without integrations filters",
			matchPod: &metadata{
				annotations: map[string]string{"newrelic.io/scrape": "true"},
			},
			matchEndpoints: &metadata{
				annotations: map[string]string{"newrelic.io/scrape": "true"},
			},
		},
		{
			name: "drop due to integration filters",
			dropPod: &metadata{
				annotations: map[string]string{"prometheus.io/scrape": "true"},
			},
			dropEndpoints: &metadata{
				annotations: map[string]string{"prometheus.io/scrape": "true"},
			},
		},
	}

	for _, testCase := range testsCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			k8sEnv := newK8sEnvironment(t)

			// Resources are only generated when metadata is specified in the test.
			matchPod := addPodToEnv(t, testCase.matchPod, "match-", k8sEnv)
			dropPod := addPodToEnv(t, testCase.dropPod, "drop-", k8sEnv)
			matchEndpointService := addEndpointsToEnv(t, testCase.matchEndpoints, "match-endpoint-", k8sEnv)
			dropEndpointService := addEndpointsToEnv(t, testCase.dropEndpoints, "drop-endpoint-", k8sEnv)

			nrConfigConfig := fmt.Sprintf(`
newrelic_remote_write:
  license_key: nrLicenseKey
common:
  scrape_interval: 1s
kubernetes:
  jobs:
    - job_name_prefix: test-k8s
      integrations_filter:
        enabled: true
        app_values:
          - test1
          - test2
      target_discovery:
        endpoints: true
        pod: true
        additional_config:
         kubeconfig_file: %s
         namespaces:
          names:
          - %s
    - job_name_prefix: newrelic-k8s
      integrations_filter:
        enabled: false
      target_discovery:
        endpoints: true
        pod: true
        additional_config:
         kubeconfig_file: %s
         namespaces:
          names:
          - %s
        filter:
          annotations:
            newrelic.io/scrape: true
  integrations_filter:
    enabled: true
    source_labels:
      - something-different1
      - app.kubernetes.io/name
      - something-different2
`, k8sEnv.kubeconfigFullPath, k8sEnv.testNamespace.Name, k8sEnv.kubeconfigFullPath, k8sEnv.testNamespace.Name)

			ps := newPrometheusServer(t)
			ps.start(t, runConfigurator(t, nrConfigConfig))

			asserter := newAsserter(ps)

			targetCount := 0

			if testCase.matchPod != nil {
				asserter.activeTargetLabels(t, map[string]string{
					"pod": matchPod.Name,
				})
				targetCount++
			}

			if testCase.matchEndpoints != nil {
				asserter.activeTargetLabels(t, map[string]string{
					"service": matchEndpointService.Name,
				})
				targetCount++
			}

			asserter.activeTargetCount(t, targetCount)

			if testCase.dropPod != nil {
				asserter.droppedTargetLabels(t, map[string]string{
					// dropped pods labels are not processed so check discovered labels.
					"__meta_kubernetes_pod_name": dropPod.Name,
				})
			}

			if testCase.dropEndpoints != nil {
				asserter.droppedTargetLabels(t, map[string]string{
					"__meta_kubernetes_service_name": dropEndpointService.Name,
				})
			}

			// no counting on dropped targets since there are more than the ones specified.
			// For instance the endpoints pod discovered by the Pod target kind will be dropped
			// since doesn't have any label/annotation.
		})

	}
}
