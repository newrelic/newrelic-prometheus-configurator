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

	"github.com/prometheus/prometheus/model/exemplar"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/storage/remote"
)

var ErrTimeout = errors.New("timeout Exceeded")

type asserter struct {
	appendable     *mockAppendable
	defaultTimeout time.Duration
	defaultBackoff time.Duration
	prometheusPort string
}

func newAsserter(options ...func(*asserter)) *asserter {
	a := &asserter{}

	a.appendable = &mockAppendable{
		latestSamples: make(map[string]mockSample),
	}

	a.defaultBackoff = time.Second
	a.defaultTimeout = time.Second * 10
	a.prometheusPort = "9090"

	for _, op := range options {
		op(a)
	}

	return a
}

func withCustomPort(prometheusPort string) func(*asserter) {
	return func(a *asserter) {
		a.prometheusPort = prometheusPort
	}
}

func (a *asserter) startRemoteWriteEndpoint(t *testing.T) *httptest.Server {
	t.Helper()

	handler := remote.NewWriteHandler(nil, a.appendable)

	remoteWriteServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}))

	t.Cleanup(func() {
		remoteWriteServer.Close()
	})

	return remoteWriteServer
}

func (a *asserter) metricName(t *testing.T, expectedMetricName ...string) {
	t.Helper()

	var lastNotFound string

	if err := retryUntilTrue(a.defaultTimeout, a.defaultBackoff, func() bool {
		for _, mn := range expectedMetricName {
			if !a.appendable.HasMetric(mn) {
				lastNotFound = mn
				return false
			}
		}

		return true
	}); err != nil {
		t.Errorf("metric not found: %s : %s", lastNotFound, err)
	}
}

func (a *asserter) prometheusServerReady(t *testing.T) {
	t.Helper()

	if err := retryUntilTrue(a.defaultTimeout, a.defaultBackoff, func() bool {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/-/healthy", a.prometheusPort))
		if err != nil {
			return false
		}

		if resp.StatusCode != http.StatusOK {
			return false
		}

		return true
	}); err != nil {
		t.Errorf("readiness probe failed")
	}
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

type mockAppendable struct {
	latestSamples map[string]mockSample
	lock          sync.Mutex
}

type mockSample struct {
	l labels.Labels
	t int64
	v float64
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

func (m *mockAppendable) AppendExemplar(_ storage.SeriesRef, l labels.Labels, e exemplar.Exemplar) (storage.SeriesRef, error) {
	return 0, nil
}

func (m *mockAppendable) HasMetric(metricName string) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	_, ok := m.latestSamples[metricName]

	return ok
}
