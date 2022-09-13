//go:build integration_test

package integration

import (
	"errors"
	"testing"
	"time"

	"github.com/newrelic/newrelic-prometheus-configurator/test/integration/mocks"

	"github.com/stretchr/testify/require"
)

const scrapeURLKey = "scrapeUrl"

var ErrTimeout = errors.New("timeout Exceeded")

type asserter struct {
	appendable       *mocks.Appendable
	defaultTimeout   time.Duration
	defaultBackoff   time.Duration
	prometheusServer *prometheusServer
}

func newAsserter(ps *prometheusServer) *asserter {
	a := &asserter{}

	a.appendable = mocks.NewAppendable()
	a.defaultBackoff = time.Second
	a.defaultTimeout = time.Second * 20
	a.prometheusServer = ps

	return a
}

// metricName checks that the asserter remote write receiver has received all expectedMetricName.
func (a *asserter) metricName(t *testing.T, expectedMetricName ...string) {
	t.Helper()

	var lastNotFound string

	err := retryUntilTrue(a.defaultTimeout, a.defaultBackoff, func() bool {
		for _, mn := range expectedMetricName {
			if _, ok := a.appendable.GetMetric(mn); !ok {
				lastNotFound = mn
				return false
			}
		}

		return true
	})
	require.NoError(t, err, "metric not found: ", lastNotFound)
}

func (a *asserter) metricLabels(t *testing.T, expectedMetricLabels map[string]string, expectedMetricName ...string) {
	t.Helper()

	var lastNotFound string

	err := retryUntilTrue(a.defaultTimeout, a.defaultBackoff, func() bool {
		for _, mn := range expectedMetricName {
			sample, ok := a.appendable.GetMetric(mn)
			if !ok {
				lastNotFound = mn
				return false
			}
			for k, v := range expectedMetricLabels {
				if actualValue := sample.Labels.Get(k); actualValue != v {
					t.Errorf("in the metric %s was not found the label %s %q!=%q", mn, k, v, actualValue)
				}
			}
		}
		return true
	})

	require.NoError(t, err, "metric not found: ", lastNotFound)
}

// prometheusServerReady probes the healthy endpoint of Prometheus.
func (a *asserter) prometheusServerReady(t *testing.T) {
	t.Helper()

	err := retryUntilTrue(a.defaultTimeout, a.defaultBackoff, func() bool {
		return a.prometheusServer.healthy(t)
	})
	require.NoError(t, err, "readiness probe failed")
}

// activeTargetCount check that that count of active targets match expectations.
func (a *asserter) activeTargetCount(t *testing.T, expectedLen int) {
	t.Helper()

	err := retryUntilTrue(a.defaultTimeout, a.defaultBackoff, func() bool {
		targets, ok := a.prometheusServer.targets(t)
		if !ok {
			return false
		}

		if len(targets.ActiveTargets) == expectedLen {
			return true
		}

		t.Logf("Active targets found: %d, expected:%d", len(targets.ActiveTargets), expectedLen)

		return false
	})

	require.NoError(t, err)
}

// droppedTargetCount check that that count of dropped targets match expectations.
func (a *asserter) droppedTargetCount(t *testing.T, expectedLen int) {
	t.Helper()

	err := retryUntilTrue(a.defaultTimeout, a.defaultBackoff, func() bool {
		targets, ok := a.prometheusServer.targets(t)
		if !ok {
			return false
		}

		if len(targets.DroppedTargets) == expectedLen {
			return true
		}

		t.Logf("Dropped targets found: %d, expected:%d", len(targets.DroppedTargets), expectedLen)

		return false
	})

	require.NoError(t, err)
}

// activeTargetLabels checks that Prometheus has at least one active target with all expected labels for
// discoveredLabels and labels fields.
func (a *asserter) activeTargetLabels(t *testing.T, expectedLabels map[string]string) {
	t.Helper()

	err := retryUntilTrue(a.defaultTimeout, a.defaultBackoff, func() bool {
		targets, ok := a.prometheusServer.targets(t)
		if !ok {
			return false
		}

		for _, at := range targets.ActiveTargets {
			allLabels := mergeLabels(at.DiscoveredLabels, at.Labels)
			if containsLabels(allLabels, expectedLabels) {
				return true
			}
		}

		return false
	})

	require.NoError(t, err)
}

// activeTargetWithScrapeURL checks that Prometheus has at least one active target with the given scrapeURL field value.
func (a *asserter) activeTargetWithScrapeURL(t *testing.T, value string) {
	t.Helper()

	err := retryUntilTrue(a.defaultTimeout, a.defaultBackoff, func() bool {
		targets, ok := a.prometheusServer.targets(t)
		if !ok {
			return false
		}

		for _, at := range targets.ActiveTargets {
			if at.ScrapeURL == value {
				return true
			}
		}

		return false
	})

	require.NoError(t, err)
}

// droppedTargetLabels checks that Prometheus has at least one dropped target with all expected labels.
func (a *asserter) droppedTargetLabels(t *testing.T, expectedLabels map[string]string) {
	t.Helper()

	err := retryUntilTrue(a.defaultTimeout, a.defaultBackoff, func() bool {
		targets, ok := a.prometheusServer.targets(t)
		if !ok {
			return false
		}

		for _, at := range targets.DroppedTargets {
			if containsLabels(at.DiscoveredLabels, expectedLabels) {
				return true
			}
		}

		return false
	})

	require.NoError(t, err)
}

func mergeLabels(maps ...map[string]string) (result map[string]string) {
	result = make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func containsLabels(labels, expectedLabels map[string]string) bool {
	for k, v := range expectedLabels {
		if val, ok := labels[k]; ok && val == v {
			continue
		}

		return false
	}

	return true
}

func retryUntilTrue(timeout time.Duration, backoff time.Duration, f func() bool) error {
	timeoutTicker := time.After(timeout)

	for {
		if f() {
			break
		}

		select {
		case <-timeoutTicker:
			return ErrTimeout
		default:
			time.Sleep(backoff)
		}
	}

	return nil
}
