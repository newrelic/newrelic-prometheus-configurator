package mocks

import (
	"context"
	"sync"

	"github.com/prometheus/prometheus/model/exemplar"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
)

// MockSample represent a sample that the appendable.go would stores.
type MockSample struct {
	Labels    labels.Labels
	Timestamp int64
	Value     float64
}

// MockAppendable implements the github.com/prometheus/prometheus/storage.Appendable interface
// which is used by the remote write server to store the received samples.
type MockAppendable struct {
	latestSamples map[string]MockSample
	lock          sync.Mutex
}

func NewMockAppendable() *MockAppendable {
	return &MockAppendable{
		latestSamples: make(map[string]MockSample),
	}
}

func (m *MockAppendable) Appender(_ context.Context) storage.Appender { //nolint: ireturn // External interface.
	return m
}

func (m *MockAppendable) Append(_ storage.SeriesRef, l labels.Labels, t int64, v float64) (storage.SeriesRef, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.latestSamples[l.Get("__name__")] = MockSample{l, t, v}

	return 0, nil
}

func (m *MockAppendable) Commit() error {
	return nil
}

func (*MockAppendable) Rollback() error {
	return nil
}

func (m *MockAppendable) AppendExemplar(_ storage.SeriesRef, _ labels.Labels, _ exemplar.Exemplar) (storage.SeriesRef, error) {
	return 0, nil
}

func (m *MockAppendable) GetMetric(metricName string) (MockSample, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	s, ok := m.latestSamples[metricName]

	return s, ok
}
