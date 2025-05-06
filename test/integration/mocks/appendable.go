//go:build integration_test

package mocks

import (
	"context"
	"sync"

	"github.com/prometheus/prometheus/model/exemplar"
	"github.com/prometheus/prometheus/model/histogram"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/model/metadata"
	"github.com/prometheus/prometheus/storage"
)

// AppendableSample represent a sample that the appendable.go would store.
type AppendableSample struct {
	Labels    labels.Labels
	Timestamp int64
	Value     float64
}

// Appendable implements the github.com/prometheus/prometheus/storage.Appendable interface
// which is used by the remote write server to store the received samples.
type Appendable struct {
	latestSamples map[string]AppendableSample
	lock          sync.Mutex
}

func NewAppendable() *Appendable {
	return &Appendable{
		latestSamples: make(map[string]AppendableSample),
	}
}

func (m *Appendable) Appender(_ context.Context) storage.Appender { //nolint: ireturn // External interface.
	return m
}

func (m *Appendable) Append(_ storage.SeriesRef, l labels.Labels, t int64, v float64) (storage.SeriesRef, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.latestSamples[l.Get("__name__")] = AppendableSample{l, t, v}

	return 0, nil
}

func (m *Appendable) AppendHistogram(ref storage.SeriesRef, l labels.Labels, t int64, h *histogram.Histogram, fh *histogram.FloatHistogram) (storage.SeriesRef, error) {
	return 0, nil
}

func (m *Appendable) AppendCTZeroSample(ref storage.SeriesRef, l labels.Labels, _, ct int64) (storage.SeriesRef, error) {
	return m.Append(ref, l, ct, 0.0)
}

func (m *Appendable) AppendHistogramCTZeroSample(ref storage.SeriesRef, l labels.Labels, _, ct int64, h *histogram.Histogram, _ *histogram.FloatHistogram) (storage.SeriesRef, error) {
	if h != nil {
		return m.AppendHistogram(ref, l, ct, &histogram.Histogram{}, nil)
	}
	return m.AppendHistogram(ref, l, ct, nil, &histogram.FloatHistogram{})
}

func (m *Appendable) Commit() error {
	return nil
}

func (*Appendable) Rollback() error {
	return nil
}

func (m *Appendable) SetOptions(_ *storage.AppendOptions) {
	panic("unimplemented")
}

func (m *Appendable) UpdateMetadata(_ storage.SeriesRef, l labels.Labels, mp metadata.Metadata) (storage.SeriesRef, error) {
	return 0, nil
}

func (m *Appendable) AppendExemplar(_ storage.SeriesRef, _ labels.Labels, _ exemplar.Exemplar) (storage.SeriesRef, error) {
	return 0, nil
}

// GetMetric returns the last sample stored by the appendable.
func (m *Appendable) GetMetric(metricName string) (AppendableSample, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	s, ok := m.latestSamples[metricName]

	return s, ok
}
