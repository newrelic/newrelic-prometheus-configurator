//go:build integration_test

package integration

import (
	"fmt"
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func Test_TargetDiscoveryFilter(t *testing.T) {
	t.Parallel()

	type metadata struct {
		annotations map[string]string
		labels      map[string]string
	}

	testsCases := []struct {
		name   string
		filter string

		// When matchPod/Endpoint is added, the test will generate a pod/endpoint
		// with specified metadata and check if is an Active Target.
		matchPod      *metadata
		matchEndpoint *metadata

		// When dropPod/Endpoint is added, the test will generate a pod/endpoint
		// with specified metadata and check if is a Dropped Target.
		dropPod      *metadata
		dropEndpoint *metadata
	}{
		{
			name: "no filter match any pod",
			matchPod: &metadata{
				annotations: map[string]string{
					"foo": "bar",
				},
			},
		},
		{
			name:   "drops on single annotation",
			filter: filter(`"single.annotation": true`, ""),
			dropPod: &metadata{
				annotations: map[string]string{
					"single.annotation": "value not matching the filter",
				},
			},
			dropEndpoint: &metadata{
				annotations: map[string]string{
					"single.annotation": "value not matching the filter",
				},
			},
		},
		{
			name:   "match a single annotation",
			filter: filter(`"single.annotation": true`, ""),
			matchPod: &metadata{
				annotations: map[string]string{
					"single.annotation": "true",
				},
			},
			matchEndpoint: &metadata{
				annotations: map[string]string{
					"single.annotation": "true",
				},
			},
		},
		{
			name:   "match annotation and label filter",
			filter: filter(`"prometheus.io/scrape": "true"`, `"k8s.io/app": "foo"`),
			matchPod: &metadata{
				annotations: map[string]string{
					"prometheus.io/scrape": "true",
				},
				labels: map[string]string{
					"k8s.io/app": "foo",
				},
			},
			matchEndpoint: &metadata{
				annotations: map[string]string{
					"prometheus.io/scrape": "true",
				},
				labels: map[string]string{
					"k8s.io/app": "foo",
				},
			},
		},
		{
			name:   "match annotation exists",
			filter: filter(`"empty.annotation": ""`, ""),
			matchPod: &metadata{
				annotations: map[string]string{
					"empty.annotation": "Value doesn't care",
				},
			},
			matchEndpoint: &metadata{
				annotations: map[string]string{
					"empty.annotation": "Value doesn't care",
				},
			},
		},
	}

	for _, testCase := range testsCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			k8sEnv := newK8sEnvironment(t)

			var matchPod *corev1.Pod
			if testCase.matchPod != nil {
				matchPod = k8sEnv.addPod(
					t,
					fakePod("match-pod-", testCase.matchPod.annotations, testCase.matchPod.labels),
				)
			}

			var dropPod *corev1.Pod
			if testCase.dropPod != nil {
				dropPod = k8sEnv.addPod(
					t,
					fakePod("drop-pod-", testCase.dropPod.annotations, testCase.dropPod.labels),
				)
			}

			var matchEndpointService *corev1.Service
			if testCase.matchEndpoint != nil {
				selector := map[string]string{"k8s.io/app": "myApp"}
				k8sEnv.addPod(
					t,
					fakePod("match-endpoint-pod-", map[string]string{}, selector),
				)

				matchEndpointService = k8sEnv.addService(
					t,
					fakeService("test", selector, testCase.matchEndpoint.annotations, testCase.matchEndpoint.labels),
				)
			}

			var dropEndpointService *corev1.Service
			if testCase.dropEndpoint != nil {
				selector := map[string]string{"k8s.io/app": "myDroppedApp"}
				k8sEnv.addPod(
					t,
					fakePod("drop-endpoint-pod-", map[string]string{}, selector),
				)

				dropEndpointService = k8sEnv.addService(
					t,
					fakeService("test", selector, testCase.dropEndpoint.annotations, testCase.dropEndpoint.labels),
				)

			}

			inputConfig := fmt.Sprintf(`
newrelic_remote_write:
  license_key: nrLicenseKey
common:
  scrape_interval: 1s
kubernetes:
  jobs:
    - job_name_prefix: test-k8s
      target_discovery:
        endpoints: true
        pod: true
        additional_config:
          kubeconfig_file: %s
          namespaces:
            names:
            - %s
        %s
`, k8sEnv.kubeconfigFullPath, k8sEnv.testNamespace.Name, testCase.filter)

			ps := newPrometheusServer(t)
			ps.start(t, runConfigurator(t, inputConfig))

			asserter := newAsserter(ps)

			targetCount := 0

			if testCase.matchPod != nil {
				asserter.activeTargetLabels(t, map[string]string{
					"pod": matchPod.Name,
				})
				targetCount++
			}

			if testCase.matchEndpoint != nil {
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

			if testCase.dropEndpoint != nil {
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

func filter(annotation, label string) string {
	return fmt.Sprintf(`
        filter:
          annotations:
            %s
          labels:
            %s
`, annotation, label)
}
