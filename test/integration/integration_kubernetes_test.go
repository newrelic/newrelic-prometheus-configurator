//go:build integration_test

package integration

import (
	"fmt"
	"io/ioutil"
	"net"
	"path"
	"strconv"
	"testing"

	"github.com/newrelic-forks/newrelic-prometheus/test/integration/mocks"
	"github.com/stretchr/testify/require"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_PodMetricsLabels(t *testing.T) {
	t.Parallel()

	k8sEnv := newK8sEnvironment(t)

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "testpod",
			Labels: map[string]string{
				"pod.label": "Value.of.label",
			},
			Annotations: map[string]string{
				"prometheus.io/scrape": "true",
			},
		},
		Spec: fakePodSpec(),
	}

	pod = k8sEnv.addPodAndWaitOnPhase(t, pod, corev1.PodRunning)

	ps := newPrometheusServer(t)

	asserter := newAsserter(ps)

	exporter := mocks.StartExporter(t)

	rw := mocks.StartRemoteWriteEndpoint(t, asserter.appendable)

	// TODO this test is using a Prom config directly since pods targets
	// are not implemented in the configurator yet.
	promConfig := path.Join(t.TempDir(), "test-config.yml")
	err := ioutil.WriteFile(promConfig, []byte(fmt.Sprintf(`
remote_write:
- url: https://foo:8999/write
  # not actually needed here but is needed when the newrelic_remote_write is used.
  proxy_url: %s
  tls_config:
    insecure_skip_verify: true

global:
  scrape_interval: 1s

scrape_configs:
- job_name: test-k8s
  # used to make prometheus use the test endpoint url instead of the real pod ip.
  proxy_url: %s
  kubernetes_sd_configs:
  - role: pod
    # used to connect to the testing cluster since prometheus is running outside of it.
    kubeconfig_file: %s
    # Filter only the current test namespace to make test independent.
    namespaces:
      names:
      - %s
  relabel_configs:
  - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
    action: keep
    regex: true

  - action: labelmap
    regex: __meta_kubernetes_pod_label_(.+)
  - source_labels: [__meta_kubernetes_namespace]
    action: replace
    target_label: namespace
  - source_labels: [__meta_kubernetes_pod_name]
    action: replace
    target_label: pod
`, rw.URL, exporter.URL, k8sEnv.kubeconfigFullPath, k8sEnv.testNamespace.Name)), 0o444)
	require.NoError(t, err)

	ps.start(t, promConfig)

	instance := net.JoinHostPort(pod.Status.PodIP, strconv.Itoa(defaultPodPort))
	expectedLabels := map[string]string{
		"pod_label": "Value.of.label",
		"pod":       pod.Name,
		"namespace": pod.Namespace,
		"instance":  instance,
		"job":       "test-k8s",
	}
	asserter.metricLabels(t, expectedLabels, "mock_gauge_metric")
}

func Test_PodDiscovery(t *testing.T) {
	t.Parallel()

	k8sEnv := newK8sEnvironment(t)

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "testpod",
			Labels: map[string]string{
				"k8s.io/app": "myApp",
			},
			Annotations: map[string]string{
				"prometheus.io/scrape": "true",
			},
		},
		Spec: fakePodSpec(),
	}

	podDropped := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "testpod-dropped",
			Labels: map[string]string{
				"k8s.io/app": "myOtherAppIDontWantToMonitor",
			},
			Annotations: map[string]string{
				"prometheus.io/scrape": "true",
			},
		},
		Spec: fakePodSpec(),
	}

	pod = k8sEnv.addPodAndWaitOnPhase(t, pod, corev1.PodRunning)
	podDropped = k8sEnv.addPodAndWaitOnPhase(t, podDropped, corev1.PodRunning)

	ps := newPrometheusServer(t)

	asserter := newAsserter(ps)

	// TODO this test is using a Prom config directly since pods targets
	// are not implemented in the configurator yet.
	promConfig := path.Join(t.TempDir(), "test-config.yml")
	err := ioutil.WriteFile(promConfig, []byte(fmt.Sprintf(`
remote_write:
- url: http://foo:8999/write

global:
  scrape_interval: 1s

scrape_configs:
- job_name: test-k8s
  kubernetes_sd_configs:
  - role: pod
    # used to connect to the testing cluster since prometheus is running outside of it.
    kubeconfig_file: %s
    # Filter only the current test namespace to make test independent.
    namespaces:
      names:
      - %s
  relabel_configs:
  - source_labels:
    - __meta_kubernetes_pod_annotation_prometheus_io_scrape
    - __meta_kubernetes_pod_label_k8s_io_app
    action: keep
    regex: true;%s
`, k8sEnv.kubeconfigFullPath, k8sEnv.testNamespace.Name, pod.Labels["k8s.io/app"])), 0o444)
	require.NoError(t, err)

	ps.start(t, promConfig)

	asserter.activeTargetLabels(t, map[string]string{"__meta_kubernetes_pod_label_k8s_io_app": pod.Labels["k8s.io/app"]})
	asserter.droppedTargetLabels(t, map[string]string{"__meta_kubernetes_pod_label_k8s_io_app": podDropped.Labels["k8s.io/app"]})
}

func Test_EndpointsDiscovery(t *testing.T) {
	t.Parallel()

	k8sEnv := newK8sEnvironment(t)

	serviceSelector := map[string]string{"k8s.io/app": "myApp"}
	// Create initial pod
	pod := fakePod("testpod", nil, serviceSelector)

	terminationGracePeriodSeconds := int64(1)

	// Create failing pod
	failedPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "testpod-failed",
			Labels: serviceSelector,
		},
		Spec: corev1.PodSpec{
			ActiveDeadlineSeconds: &terminationGracePeriodSeconds,
			Containers: []corev1.Container{
				{
					Name:  "fake-exporter",
					Image: "this-image-dont-exist-pod-will-fail",
				},
			},
		},
	}

	// Create service
	svc := fakeService(
		"test",
		serviceSelector,
		map[string]string{
			"prometheus.io/scheme":     "https",
			"prometheus.io/path":       "/custom-path",
			"prometheus.io/port":       "8001",
			"prometheus.io/param_test": "test-param",
		},
		map[string]string{
			"k8s.io/app": "myApp",
			"test.label": "test.value",
		},
	)

	pod = k8sEnv.addPodAndWaitOnPhase(t, pod, corev1.PodRunning)
	failedPod = k8sEnv.addPodAndWaitOnPhase(t, failedPod, corev1.PodFailed)
	svc = k8sEnv.addService(t, svc)

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
        additional_config:
         kubeconfig_file: %s
         namespaces:
          names:
          - %s
`, k8sEnv.kubeconfigFullPath, k8sEnv.testNamespace.Name)

	outputConfigPath := runConfigurator(t, inputConfig)

	ps := newPrometheusServer(t)
	ps.start(t, outputConfigPath)

	// Build scrapeURL
	instance := net.JoinHostPort(pod.Status.PodIP, svc.Annotations["prometheus.io/port"])
	params := "?test=" + svc.Annotations["prometheus.io/param_test"]

	scrapeURL := fmt.Sprintf("%s://%s%s%s",
		svc.Annotations["prometheus.io/scheme"],
		instance,
		svc.Annotations["prometheus.io/path"],
		params,
	)

	asserter := newAsserter(ps)

	// Active targets
	asserter.activeTargetCount(t, 1)
	asserter.activeTargetField(t, scrapeURLKey, scrapeURL)
	asserter.activeTargetLabels(t, map[string]string{
		"namespace":  k8sEnv.testNamespace.Name,
		"service":    svc.Name,
		"node":       pod.Spec.NodeName,
		"test_label": svc.Labels["test.label"],
	})

	// Dropped targets
	asserter.droppedTargetLabels(t, map[string]string{"__meta_kubernetes_pod_name": failedPod.Name})
}
