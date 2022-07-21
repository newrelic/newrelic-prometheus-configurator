package integration

import (
	"encoding/json"
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

// activeTargetsURLs checks that Prometheus has active targets for each expectedTargetUrl.
func (a *asserter) activeTargetsURLs(t *testing.T, expectedTargetURL ...string) {
	t.Helper()

	promResp := &response{}

	err := retryUntilTrue(a.defaultTimeout, a.defaultBackoff, func() bool {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/v1/targets?state=active", a.prometheusPort))
		if err != nil {
			return false
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(promResp)
		require.NoError(t, err)

		targets := &targetDiscovery{}
		err = json.Unmarshal(promResp.Data, targets)
		require.NoError(t, err)

		for _, eu := range expectedTargetURL {
			for i, at := range targets.ActiveTargets {
				if at.ScrapeURL == eu {
					break
				}

				// fail if there is no active target with the expectedTargetUrl.
				if i == (len(targets.ActiveTargets) - 1) {
					return false
				}
			}
		}

		return true
	})

	require.NoError(t, err, "expected targets URL: ", expectedTargetURL, "Prometheus targets response: ", promResp)
}

// droppedTargetLabels checks that Prometheus has dropped target with labels.
func (a *asserter) droppedTargetLabels(t *testing.T, expectedLabels map[string]string) {
	t.Helper()

	promResp := &response{}

	err := retryUntilTrue(a.defaultTimeout, a.defaultBackoff, func() bool {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/v1/targets?state=dropped", a.prometheusPort))
		if err != nil {
			return false
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(promResp)
		require.NoError(t, err)

		targets := &targetDiscovery{}
		err = json.Unmarshal(promResp.Data, targets)
		require.NoError(t, err)

		// Iterates over the dropped targets finding the first one that contains all
		// the expected labels.
		for _, dt := range targets.DroppedTargets {
			i := 0
			for k, v := range expectedLabels {
				if val, ok := dt.DiscoveredLabels[k]; !ok || val != v {
					break
				}

				i++
				// Returns true if all labels are in one of the dropped targets.
				if i == len(expectedLabels) {
					return true
				}
			}
		}

		return false
	})

	require.NoError(t, err)
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

type response struct {
	Status    string          `json:"status"`
	Data      json.RawMessage `json:"data,omitempty"`
	ErrorType string          `json:"errorType,omitempty"`
	Error     string          `json:"error,omitempty"`
	Warnings  []string        `json:"warnings,omitempty"`
}

// target has the information for one target.
type target struct {
	// Labels before any processing.
	DiscoveredLabels map[string]string `json:"discoveredLabels"`
	// Any labels that are added to this target and its metrics.
	Labels map[string]string `json:"labels"`

	ScrapePool string `json:"scrapePool"`
	ScrapeURL  string `json:"scrapeUrl"`
	GlobalURL  string `json:"globalUrl"`

	LastError          string    `json:"lastError"`
	LastScrape         time.Time `json:"lastScrape"`
	LastScrapeDuration float64   `json:"lastScrapeDuration"`
	Health             string    `json:"health"`

	ScrapeInterval string `json:"scrapeInterval"`
	ScrapeTimeout  string `json:"scrapeTimeout"`
}

// droppedTarget has the information for one target that was dropped during relabelling.
type droppedTarget struct {
	// Labels before any processing.
	DiscoveredLabels map[string]string `json:"discoveredLabels"`
}

// targetDiscovery has all the active targets.
type targetDiscovery struct {
	ActiveTargets  []*target        `json:"activeTargets"`
	DroppedTargets []*droppedTarget `json:"droppedTargets"`
}
