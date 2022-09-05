//go:build integration_test

package integration

import (
	"fmt"
	"net"
	"strconv"
	"testing"

	"github.com/newrelic/newrelic-prometheus-configurator/test/integration/mocks"

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
				"test.label": "test.value",
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

	rw := mocks.StartRemoteWriteEndpoint(t, asserter.appendable)

	ex := mocks.StartExporter(t)

	exporterURL := ex.Listener.Addr().String()

	nrConfigConfig := fmt.Sprintf(`
newrelic_remote_write:
  license_key: nrLicenseKey
  proxy_url: %s
  tls_config:
    insecure_skip_verify: true
common:
  scrape_interval: 1s
kubernetes:
  jobs:
    - job_name_prefix: test-k8s
      proxy_url: http://%s
      target_discovery:
        pod: true
        additional_config:
         kubeconfig_file: %s
         namespaces:
          names:
          - %s
`, rw.URL, exporterURL, k8sEnv.kubeconfigFullPath, k8sEnv.testNamespace.Name)

	t.Log(nrConfigConfig)
	prometheusConfigConfigPath := runConfigurator(t, nrConfigConfig)

	ps.start(t, prometheusConfigConfigPath)

	instance := net.JoinHostPort(pod.Status.PodIP, strconv.Itoa(defaultPodPort))
	expectedLabels := map[string]string{
		"test_label": "test.value",
		"pod":        pod.Name,
		"namespace":  pod.Namespace,
		"instance":   instance,
		"job":        "test-k8s-pod",
	}
	asserter.metricLabels(t, expectedLabels, "mock_gauge_metric")
}

func Test_PodPhaseDropRule(t *testing.T) {
	t.Parallel()

	k8sEnv := newK8sEnvironment(t)
	terminationGracePeriodSeconds := int64(1)

	// Create running pod
	runningPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "testpod-running",
			Labels: map[string]string{
				"k8s.io/app": "myApp",
			},
			Annotations: map[string]string{
				"prometheus.io/scrape": "true",
			},
		},
		Spec: fakePodSpec(),
	}

	// Create failing pod
	failedPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "testpod-failed",
			Labels: map[string]string{
				"k8s.io/app": "myApp",
			},
			Annotations: map[string]string{
				"prometheus.io/scrape": "true",
			},
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

	// Create a succeeded pod which will be added to dropped targets
	succeededPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "testpod-succeeded",
			Labels: map[string]string{
				"k8s.io/app": "myApp",
			},
			Annotations: map[string]string{
				"prometheus.io/scrape": "true",
			},
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyNever,
			Containers: []corev1.Container{
				{
					Name:  "busybox",
					Image: "busybox:latest",
					Command: []string{
						"sh",
						"-c",
						"exit 0",
					},
				},
			},
		},
	}

	runningPod = k8sEnv.addPodAndWaitOnPhase(t, runningPod, corev1.PodRunning)
	failedPod = k8sEnv.addPodAndWaitOnPhase(t, failedPod, corev1.PodFailed)
	succeededPod = k8sEnv.addPodAndWaitOnPhase(t, succeededPod, corev1.PodSucceeded)

	nrConfigConfig := fmt.Sprintf(`
newrelic_remote_write:
  license_key: nrLicenseKey
common:
  scrape_interval: 1s
kubernetes:
  jobs:
    - job_name_prefix: test-k8s
      target_discovery:
        pod: true
        additional_config:
         kubeconfig_file: %s
         namespaces:
          names:
          - %s
`, k8sEnv.kubeconfigFullPath, k8sEnv.testNamespace.Name)

	prometheusConfigConfigPath := runConfigurator(t, nrConfigConfig)

	ps := newPrometheusServer(t)
	ps.start(t, prometheusConfigConfigPath)

	asserter := newAsserter(ps)

	asserter.activeTargetLabels(t, map[string]string{
		"__meta_kubernetes_pod_name": runningPod.Name,
	})
	asserter.droppedTargetLabels(t, map[string]string{
		"__meta_kubernetes_pod_name": failedPod.Name,
	})
	asserter.droppedTargetLabels(t, map[string]string{
		"__meta_kubernetes_pod_name": succeededPod.Name,
	})
}

func Test_PodRelabelRules(t *testing.T) {
	t.Parallel()

	k8sEnv := newK8sEnvironment(t)

	// Create running pod
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "testpod-running",
			Labels: map[string]string{
				"k8s.io/app": "myApp",
			},
			Annotations: map[string]string{
				"prometheus.io/scrape": "true",
				"prometheus.io/scheme": "https",
				"prometheus.io/port":   "8001",
				"prometheus.io/path":   "/custom-path",
			},
		},
		Spec: fakePodSpec(),
	}

	pod = k8sEnv.addPodAndWaitOnPhase(t, pod, corev1.PodRunning)

	nrConfigConfig := fmt.Sprintf(`
newrelic_remote_write:
  license_key: nrLicenseKey
common:
  scrape_interval: 1s
kubernetes:
  jobs:
    - job_name_prefix: test-k8s
      target_discovery:
        pod: true
        additional_config:
         kubeconfig_file: %s
         namespaces:
          names:
          - %s
`, k8sEnv.kubeconfigFullPath, k8sEnv.testNamespace.Name)

	prometheusConfigConfigPath := runConfigurator(t, nrConfigConfig)

	ps := newPrometheusServer(t)
	ps.start(t, prometheusConfigConfigPath)

	scrapeURL := fmt.Sprintf("%s://%s:%s%s",
		pod.Annotations["prometheus.io/scheme"],
		pod.Status.PodIP,
		pod.Annotations["prometheus.io/port"],
		pod.Annotations["prometheus.io/path"],
	)

	asserter := newAsserter(ps)
	asserter.activeTargetField(t, scrapeURLKey, scrapeURL)
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

	nrConfigConfig := fmt.Sprintf(`
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

	prometheusConfigConfigPath := runConfigurator(t, nrConfigConfig)

	ps := newPrometheusServer(t)
	ps.start(t, prometheusConfigConfigPath)

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
