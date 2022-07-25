package integration

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/go-kit/log"
	"github.com/prometheus/prometheus/model/exemplar"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/storage/remote"
	"github.com/stretchr/testify/require"
)

var ErrTimeout = errors.New("timeout Exceeded")

type asserter struct {
	appendable     *mockAppendable
	defaultTimeout time.Duration
	defaultBackoff time.Duration
	prometheusPort string
}

func newAsserter(prometheusPort string) *asserter {
	a := &asserter{}

	a.appendable = &mockAppendable{
		latestSamples: make(map[string]mockSample),
	}

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
			if !a.appendable.HasMetric(mn) {
				lastNotFound = mn
				return false
			}
		}

		return true
	})
	require.NoError(t, err, "metric not found: ", lastNotFound)
}

func (a *asserter) metricWithLabels(t *testing.T, expectedMetricName string, expectedlabels []string) {
	t.Helper()

	err := retryUntilTrue(a.defaultTimeout, a.defaultBackoff, func() bool {
		return a.appendable.HasMetricWithLabels(expectedMetricName, expectedlabels)
	})
	require.NoError(t, err, "metric with labels not found: ", expectedMetricName, expectedlabels)
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

func startRemoteWriteEndpoint(t *testing.T, appendable storage.Appendable) *httptest.Server {
	t.Helper()

	handler := remote.NewWriteHandler(log.NewNopLogger(), appendable)

	remoteWriteServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}))

	t.Cleanup(func() {
		remoteWriteServer.Close()
	})

	return remoteWriteServer
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

// mockAppendable implements the github.com/prometheus/prometheus/storage.Appendable interface
// which is used by the remote write server to store the received samples.
type mockAppendable struct {
	latestSamples map[string]mockSample
	lock          sync.Mutex
}

type mockSample struct {
	labels    labels.Labels
	timestamp int64
	value     float64
}

func (m *mockAppendable) Appender(_ context.Context) storage.Appender { //nolint: ireturn // External interface.
	return m
}

func (m *mockAppendable) Append(_ storage.SeriesRef, l labels.Labels, t int64, v float64) (storage.SeriesRef, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.latestSamples[l.Get("__name__")] = mockSample{l, t, v}

	return 0, nil
}

func (m *mockAppendable) Commit() error {
	return nil
}

func (*mockAppendable) Rollback() error {
	return nil
}

func (m *mockAppendable) AppendExemplar(_ storage.SeriesRef, _ labels.Labels, _ exemplar.Exemplar) (storage.SeriesRef, error) {
	return 0, nil
}

func (m *mockAppendable) HasMetric(metricName string) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	_, ok := m.latestSamples[metricName]

	return ok
}

func (m *mockAppendable) HasMetricWithLabels(metricName string, expectedLabels []string) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	sample, ok := m.latestSamples[metricName]
	if !ok {
		return false
	}

	for _, expected := range expectedLabels {
		for i, label := range sample.labels {
			if label.Name == expected {
				break
			}
			// expected label missing.
			if i == (len(sample.labels) - 1) {
				return false
			}
		}
	}

	return true
}
