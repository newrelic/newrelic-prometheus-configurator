//go:build integration_test

package integration

import (
	"crypto/md5"
	"fmt"
	"sync"
	"testing"
)

// checkPrometheus shards sets up as many prometheus servers as defined in `numberOfShards` in parallel and it
// executes the `checkFn` for each of them. Note the `checkFn` should start the prometheus server before performing
// any check.
func checkPrometheusShards(
	t *testing.T, numberOfShards int,
	checkFn func(ps *prometheusServer, asserter *asserter, shardIndex int),
) {
	t.Helper()

	var wg sync.WaitGroup

	for shardIndex := 0; shardIndex < numberOfShards; shardIndex++ {
		wg.Add(1)

		go func(shardIndex int) {
			defer wg.Done()

			ps := newPrometheusServer(t)
			asserter := newAsserter(ps)

			checkFn(ps, asserter, shardIndex)
		}(shardIndex)
	}

	wg.Wait()
}

// k8sShardingNRConfig is a helper to provide NR config to sharding tests involving k8s
func k8sShardingNRConfig(nShards int, shardIndex int, k8sEnv k8sEnvironment, pod, endpoints bool) string {
	return fmt.Sprintf(`
newrelic_remote_write:
  license_key: nrLicenseKey
common:
  scrape_interval: 1s
sharding:
  total_shards_count: %d
  shard_index: %d
kubernetes:
  jobs:
    - job_name_prefix: test-k8s
      target_discovery:
        pod: %t
        endpoints: %t
        additional_config:
         kubeconfig_file: %s
         namespaces:
          names:
          - %s
`, nShards, shardIndex, pod, endpoints, k8sEnv.kubeconfigFullPath, k8sEnv.testNamespace.Name)
}

// shardingHashMod returns the modulus of hash sum, as it is done in prometheus.
// Check: <https://github.com/prometheus/prometheus/blob/8b863c42dd956d35d18a7a0b39c89c86adf7cebf/model/relabel/relabel.go#L250>
func shardingHashMod(value string, hashMod uint64) uint64 {
	return shardingSum64(md5.Sum([]byte(value))) % uint64(hashMod)
}

// shardingSum64 sums the md5 hash to an uint64. Taken from prometheus relabel implementation.
func shardingSum64(hash [md5.Size]byte) uint64 {
	var s uint64

	for i, b := range hash {
		shift := uint64((md5.Size - i - 1) * 8)

		s |= uint64(b) << shift
	}
	return s
}
