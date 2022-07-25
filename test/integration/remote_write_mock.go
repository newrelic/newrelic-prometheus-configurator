package integration

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/go-kit/log"
	"github.com/prometheus/prometheus/model/exemplar"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/storage/remote"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

// startRemoteWriteEndpoint start an remote write endpoint with proxy.
func startRemoteWriteEndpoint(t *testing.T, appendable storage.Appendable) *httptest.Server {
	t.Helper()

	handler := remote.NewWriteHandler(log.NewJSONLogger(os.Stderr), appendable)

	url := ""
	remoteWriteServer := httptest.NewTLSServer(handlerWithProxy(t, handler, &url))
	// There is a small race condition. If the first request comes before this is set the dial fails.
	url = strings.Replace(remoteWriteServer.URL, "https://", "", 1)

	t.Cleanup(func() {
		remoteWriteServer.Close()
	})

	return remoteWriteServer
}

// handlerWithProxy injects a proxy to a handler.
// If the method of the request is not Connect everything works as usual.
// On the other hand with Connect the handler Hijack the connection and creates two pipes connecting the source (promtheus server)
// and the destination (in this case it is the proxy itself).
// This is needed whenever we connect through an HTTPS proxy against an HTTPS server.
func handlerWithProxy(t *testing.T, handler http.Handler, url *string) http.Handler {
	t.Helper()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodConnect {
			handler.ServeHTTP(w, r)
			return
		}
		defer r.Body.Close()
		conn, err := net.Dial("tcp", *url)
		require.NoError(t, err)

		w.WriteHeader(http.StatusOK)

		hj, ok := w.(http.Hijacker)
		if !ok {
			t.Fatalf("Unable to hijack connection")
		}

		reqConn, wbuf, err := hj.Hijack()
		require.NoError(t, err)

		defer reqConn.Close()
		defer wbuf.Flush()

		g := errgroup.Group{}
		g.Go(func() error { return pipe(t, reqConn, conn) })
		g.Go(func() error { return pipe(t, conn, reqConn) })

		if err := g.Wait(); err != nil {
			require.NoError(t, err)
		}
	})
}

func pipe(t *testing.T, from net.Conn, to net.Conn) error {
	t.Helper()

	defer from.Close()
	_, err := io.Copy(from, to)
	if err != nil && !strings.Contains(err.Error(), "closed network") {
		return fmt.Errorf("error in pipe: %w", err)
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

func (m *mockAppendable) GetMetric(metricName string) (mockSample, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	s, ok := m.latestSamples[metricName]

	return s, ok
}
