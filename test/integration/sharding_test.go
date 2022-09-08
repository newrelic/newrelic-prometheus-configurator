//go:build integration_test

package integration

import (
	"fmt"
	"net"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/newrelic/newrelic-prometheus-configurator/test/integration/mocks"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_Sharding_Pod(t *testing.T) {
	t.Parallel()

	numberOfShards := 2

	k8sEnv := newK8sEnvironment(t)

	pod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "test-pod"}, Spec: fakePodSpec()}
	pod = k8sEnv.addPodAndWaitOnPhase(t, pod, corev1.PodRunning)

	podAddress := net.JoinHostPort(pod.Status.PodIP, strconv.Itoa(defaultPodPort))
	mod := shardingHashMod(podAddress, uint64(numberOfShards))

	checkPrometheusShards(t, numberOfShards, func(ps *prometheusServer, asserter *asserter, rw *httptest.Server, shardIndex int) {
		nrConfig := k8sShardingNRConfig(rw.URL, numberOfShards, shardIndex, k8sEnv, true, false)
		prometheusConfigPath := runConfigurator(t, nrConfig)
		ps.start(t, prometheusConfigPath)

		// Only the server whose shardIndex is equal to the address hash-mod should scrape the pod.
		if mod == uint64(shardIndex) {
			asserter.activeTargetLabels(t, map[string]string{"__meta_kubernetes_pod_name": pod.Name})
		} else {
			asserter.droppedTargetLabels(t, map[string]string{"__meta_kubernetes_pod_name": pod.Name})
		}
	})
}

func Test_Sharding_Endpoints(t *testing.T) {
	t.Parallel()

	numberOfShards := 3

	k8sEnv := newK8sEnvironment(t)

	serviceSelector := map[string]string{"k8s.io/app": "myApp"}

	// Create a service using the service previously defined service selector.
	svc := fakeService(
		"test-service",
		serviceSelector,
		map[string]string{},
		map[string]string{"k8s.io/app": "myApp"},
	)
	svc = k8sEnv.addService(t, svc)

	pods := k8sEnv.addManyPodsWaitingOnPhase(t, 5, corev1.PodRunning, func(index int) *corev1.Pod {
		return fakePod(fmt.Sprintf("test-pod-%d", index), nil, serviceSelector)
	})
	scrapeURLHashMod := map[string]uint64{}
	for pod := range pods {
		address := net.JoinHostPort(pod.Status.PodIP, strconv.Itoa(defaultPodPort))
		mod := shardingHashMod(address, uint64(numberOfShards))
		scrapeURLHashMod[fmt.Sprintf("http://%s/metrics", address)] = mod
	}

	// Start a prometheus server for each shard.
	checkPrometheusShards(t, numberOfShards, func(ps *prometheusServer, asserter *asserter, rw *httptest.Server, shardIndex int) {
		nrConfig := k8sShardingNRConfig(rw.URL, numberOfShards, shardIndex, k8sEnv, false, true)
		prometheusConfigPath := runConfigurator(t, nrConfig)
		ps.start(t, prometheusConfigPath)

		shouldScrapeAny := false
		for scrapeURL, mod := range scrapeURLHashMod {
			// The server whose shardIndex corresponds to an address hash-mod should scrape the corresponding target url.
			if mod == uint64(shardIndex) {
				shouldScrapeAny = true
				asserter.activeTargetLabels(t, map[string]string{"service": svc.Name})
				asserter.activeTargetWithScrapeURL(t, scrapeURL)
			}
		}
		// If the index does not correspond to any address hash-mod, it should not have any active-targets.
		if !shouldScrapeAny {
			asserter.activeTargetCount(t, 0)
			// The droppedTargetLabel is '__meta_kubernetes_service_name' instead of 'service' because the has-mod
			// rule should be applied first.
			asserter.droppedTargetLabels(t, map[string]string{"__meta_kubernetes_service_name": svc.Name})
		}
	})
}

func Test_Sharding_Endpoints_Sharing_Services(t *testing.T) {
	t.Parallel()

	numberOfShards := 2

	k8sEnv := newK8sEnvironment(t)

	serviceSelector := map[string]string{"k8s.io/app": "myApp"}

	// Create only one pod with the service selector
	pod := fakePod(fmt.Sprintf("test-pod"), nil, serviceSelector)
	pod = k8sEnv.addPodAndWaitOnPhase(t, pod, corev1.PodRunning)

	podAddress := net.JoinHostPort(pod.Status.PodIP, strconv.Itoa(defaultPodPort))
	mod := shardingHashMod(podAddress, uint64(numberOfShards))

	// Create many services sharing the same selector
	numberOfServices := 10
	for i := 0; i < numberOfServices; i++ {
		svc := fakeService(
			fmt.Sprintf("test-service-%d", i),
			serviceSelector,
			map[string]string{},
			map[string]string{"k8s.io/app": "myApp"},
		)
		k8sEnv.addService(t, svc)
	}

	checkPrometheusShards(t, numberOfShards, func(ps *prometheusServer, asserter *asserter, rw *httptest.Server, shardIndex int) {
		nrConfig := k8sShardingNRConfig(rw.URL, numberOfShards, shardIndex, k8sEnv, false, true)
		prometheusConfigPath := runConfigurator(t, nrConfig)
		ps.start(t, prometheusConfigPath)

		// Only the server whose shardIndex is equal to any address will scrape all the services because we use
		// __address__ to get the shard-mod.
		// This scenario is not expected to be common outside testing environments, if a different behavior was
		// needed we would need to updated the labels used to obtain the shard-mod.
		if mod == uint64(shardIndex) {
			asserter.activeTargetCount(t, numberOfServices)
		} else {
			asserter.activeTargetCount(t, 0)
		}
	})
}

func Test_Sharding_Static_Targets(t *testing.T) {
	t.Parallel()

	numberOfShards := 2

	// Create a mock for the static target.
	ex := mocks.StartExporter(t)
	address := ex.Listener.Addr().String()
	mod := shardingHashMod(address, uint64(numberOfShards))

	checkPrometheusShards(t, numberOfShards, func(ps *prometheusServer, asserter *asserter, rw *httptest.Server, shardIndex int) {
		nrConfig := fmt.Sprintf(`
sharding:
  total_shards_count: %d
  shard_index: %d

static_targets:
  jobs:
    - job_name: metrics-a
      scrape_interval: 1s
      targets:
        - "%s"
      metrics_path: /metrics-a
      labels:
        custom_label: foo

newrelic_remote_write:
  license_key: nrLicenseKey
  proxy_url: %s
  tls_config:
    insecure_skip_verify: true
`, numberOfShards, shardIndex, address, rw.URL)

		prometheusConfigPath := runConfigurator(t, nrConfig)
		ps.start(t, prometheusConfigPath)

		// Only the server whose shard-index is equal to the static-target's address hash-mod should scrape it.
		if mod == uint64(shardIndex) {
			asserter.activeTargetCount(t, 1)
		} else {
			asserter.activeTargetCount(t, 0)
		}
	})
}
