package integration

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/newrelic-forks/newrelic-prometheus/test/integration/mocks"

	"github.com/stretchr/testify/require"
)

var ErrTimeout = errors.New("timeout Exceeded")

type asserter struct {
	appendable     *mocks.Appendable
	defaultTimeout time.Duration
	defaultBackoff time.Duration
	prometheusPort string
}

func newAsserter(prometheusPort string) *asserter {
	a := &asserter{}

	a.appendable = mocks.NewAppendable()
	a.defaultBackoff = time.Second
	a.defaultTimeout = time.Second * 20
	a.prometheusPort = prometheusPort

	return a
}

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

func (a *asserter) prometheusServerReady(t *testing.T) {
	t.Helper()

	err := retryUntilTrue(a.defaultTimeout, a.defaultBackoff, func() bool {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/-/healthy", a.prometheusPort))
		if err != nil {
			return false
		}

		if resp.StatusCode != http.StatusOK {
			return false
		}

		return true
	})
	require.NoError(t, err, "readiness probe failed")
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
