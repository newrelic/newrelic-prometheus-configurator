//go:build integration_test

package integration

import (
	"fmt"
	"net"
	"strconv"
	"strings"
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

	mod := shardingHashMod(pod.Status.PodIP, uint64(numberOfShards))

	checkPrometheusShards(t, numberOfShards, func(ps *prometheusServer, asserter *asserter, shardIndex int) {
		nrConfig := k8sShardingNRConfig(numberOfShards, shardIndex, k8sEnv, true, false)
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

// If you specify scraping annotation for a kubernetes pod, typically for each container port a target is added by kubernetes/pod.go.
// If you have a second container in the pod, that does not expose any container ports, then kubernetes/pod.go will add the IP address
// to the target, but no port value. So the format would be 1.2.3.4 instead of 1.2.3.4:5432.
// In the scenario that there are two containers in a pod, one exports a port and the other does not export port. Two target addresses
// will be generated (1.2.3.4:5432 and 1.2.3.4). The current sharding relabel config uses target address and hash function to calculate
// the sharding assignment. Suppose there are two Prometheus agents, it is possible that prometheus agent 0 hashes the 1.2.3.4:5432 to 0
// and prometheus agent 0 hashes the 1.2.3.4 to 1. That means both agents accept the pod as its own target. This results in duplicate metrics from agents.
// To fix it, we extract the IP as the only source of hash function. In this way, only one agent will accept the pod as the agent.
// This test case is to test this scenario.
func Test_Sharding_Pod_With_Two_Containers(t *testing.T) {
	t.Parallel()

	numberOfShards := 2

	k8sEnv := newK8sEnvironment(t)

	pod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "test-pod"}, Spec: fakePodSpecWithTwoContainers()}
	pod = k8sEnv.addPodAndWaitOnPhase(t, pod, corev1.PodRunning)

	mod := shardingHashMod(pod.Status.PodIP, uint64(numberOfShards))

	checkPrometheusShards(t, numberOfShards, func(ps *prometheusServer, asserter *asserter, shardIndex int) {
		nrConfig := k8sShardingNRConfig(numberOfShards, shardIndex, k8sEnv, true, false)
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
		mod := shardingHashMod(pod.Status.PodIP, uint64(numberOfShards))
		scrapeURLHashMod[fmt.Sprintf("http://%s/metrics", address)] = mod
	}

	// Start a prometheus server for each shard.
	checkPrometheusShards(t, numberOfShards, func(ps *prometheusServer, asserter *asserter, shardIndex int) {
		nrConfig := k8sShardingNRConfig(numberOfShards, shardIndex, k8sEnv, false, true)
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

	mod := shardingHashMod(pod.Status.PodIP, uint64(numberOfShards))

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

	checkPrometheusShards(t, numberOfShards, func(ps *prometheusServer, asserter *asserter, shardIndex int) {
		nrConfig := k8sShardingNRConfig(numberOfShards, shardIndex, k8sEnv, false, true)
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
	// only use IP as the sharding hash function input
	addressComponents := strings.Split(address, ":")
	mod := shardingHashMod(addressComponents[0], uint64(numberOfShards))

	checkPrometheusShards(t, numberOfShards, func(ps *prometheusServer, asserter *asserter, shardIndex int) {
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
`, numberOfShards, shardIndex, address)

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

func Test_Sharding_Skip_Sharding(t *testing.T) {
	t.Parallel()

	numberOfShards := 2

	// Create a mock for each static target.
	exSkip := mocks.StartExporter(t)
	addressSkip := exSkip.Listener.Addr().String()
	scrapeURLSkip := fmt.Sprintf("http://%s/metrics", addressSkip)
	exRegular := mocks.StartExporter(t)
	addressRegular := exRegular.Listener.Addr().String()
	scrapeURLRegular := fmt.Sprintf("http://%s/metrics-a", addressRegular)
	// only use IP as the sharding hash function input
	addressComponents := strings.Split(addressRegular, ":")
	mod := shardingHashMod(addressComponents[0], uint64(numberOfShards))

	checkPrometheusShards(t, numberOfShards, func(ps *prometheusServer, asserter *asserter, shardIndex int) {
		nrConfig := fmt.Sprintf(`
sharding:
  total_shards_count: %d
  shard_index: %d

static_targets:
  jobs:
    - job_name: skip-sharding-job
      skip_sharding: true
      targets:
        - "%s"
    - job_name: regular-job
      targets:
        - "%s"
      metrics_path: /metrics-a

newrelic_remote_write:
  license_key: nrLicenseKey
`, numberOfShards, shardIndex, addressSkip, addressRegular)

		prometheusConfigPath := runConfigurator(t, nrConfig)
		ps.start(t, prometheusConfigPath)

		// The skip-sharding-job should be scrapped by all servers
		asserter.activeTargetWithScrapeURL(t, scrapeURLSkip)

		// The regular job should be scrapped only for the server whose shard-index is equal to the corresponding mod.
		if mod == uint64(shardIndex) {
			asserter.activeTargetWithScrapeURL(t, scrapeURLRegular)
		} else {
			asserter.droppedTargetLabels(t, map[string]string{"__metrics_path__": "/metrics-a"})
		}
	})
}
