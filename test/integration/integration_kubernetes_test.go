//go:build integration_test

package integration

import (
	"fmt"
	"io/ioutil"
	"net"
	"path"
	"strconv"
	"testing"
	"time"

	"github.com/newrelic-forks/newrelic-prometheus/configurator"
	"github.com/newrelic-forks/newrelic-prometheus/test/integration/mocks"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

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

	// Create initial pod
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

	terminationGracePeriodSeconds := int64(1)

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

	// Create service
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testservice",
			Namespace: k8sEnv.testNamespace.Name,
			Labels: map[string]string{
				"k8s.io/app": "myApp",
				"test.label": "test.value",
			},
			Annotations: map[string]string{
				"prometheus.io/scrape":     "true",
				"prometheus.io/scheme":     "https",
				"prometheus.io/path":       "/custom-path",
				"prometheus.io/port":       "8001",
				"prometheus.io/param_test": "test-param",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:     80,
					Protocol: "TCP",
				},
			},
			Selector: map[string]string{
				"k8s.io/app": "myApp",
			},
		},
	}

	pod = k8sEnv.addPodAndWaitOnPhase(t, pod, corev1.PodRunning)
	failedPod = k8sEnv.addPodAndWaitOnPhase(t, failedPod, corev1.PodFailed)
	svc = k8sEnv.addService(t, svc)

	input := configurator.Input{
		Common: configurator.GlobalConfig{
			ScrapeInterval: 1 * time.Second,
		},
		Kubernetes: configurator.KubernetesInput{
			Jobs: []configurator.KubernetesJob{
				{
					JobInput:      configurator.JobInput{},
					JobNamePrefix: "test-k8s",
					TargetDiscovery: configurator.KubernetesTargetDiscovery{
						Pod:       false,
						Endpoints: true,
						Filter:    nil,
						AdditionalConfig: &configurator.AdditionalConfig{
							KubeconfigFile: k8sEnv.kubeconfigFullPath,
							Namespaces: &configurator.KubernetesSdNamespace{
								Names: []string{k8sEnv.testNamespace.Name},
							},
						},
					},
				},
			},
		},
	}

	output, err := configurator.BuildOutput(&input)
	require.NoError(t, err)

	file, err := yaml.Marshal(output)
	require.NoError(t, err)

	promConfig := path.Join(t.TempDir(), "test-config.yml")
	err = ioutil.WriteFile(promConfig, file, 0o444)
	require.NoError(t, err)

	ps := newPrometheusServer(t)
	ps.start(t, promConfig)

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
	asserter.activeTargetField(t, scrapeURLKey, scrapeURL)
	asserter.activeTargetLabels(t, map[string]string{"namespace": k8sEnv.testNamespace.Name})
	asserter.activeTargetLabels(t, map[string]string{"service": svc.Name})
	asserter.activeTargetLabels(t, map[string]string{"node": pod.Spec.NodeName})
	asserter.activeTargetLabels(t, map[string]string{"test_label": svc.Labels["test.label"]})

	// Dropped targets
	asserter.droppedTargetLabels(t, map[string]string{"__meta_kubernetes_pod_name": failedPod.Name})
}
