//go:build integration_test

package mocks

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"syscall"
	"testing"

	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/storage/remote"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

// StartRemoteWriteEndpoint start an remote write endpoint with proxy.
func StartRemoteWriteEndpoint(t *testing.T, appendable storage.Appendable) *httptest.Server {
	t.Helper()

	handler := remote.NewWriteHandler(mockLog{t}, appendable)

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
		require.True(t, ok, "Unable to hijack connection")

		reqConn, wbuf, err := hj.Hijack()
		require.NoError(t, err)

		defer reqConn.Close()
		defer wbuf.Flush()

		g := errgroup.Group{}
		g.Go(func() error { return pipe(reqConn, conn) })
		g.Go(func() error { return pipe(conn, reqConn) })

		err = g.Wait()
		if err != nil {
			// We cannot check it with assert.NoError since at this point the test is usually completed and merely
			// releasing all resources and calling t.Fail would trigger a panic.
			t.Logf("Error while waiting for the pipe copying: %s", err.Error())
		}
	})
}

func pipe(from net.Conn, to net.Conn) error {
	defer from.Close()

	_, err := io.Copy(from, to)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, net.ErrClosed):
		return nil
	case errors.Is(err, syscall.ECONNRESET):
		return nil
	default:
		return fmt.Errorf("error in pipe: %w", err)
	}
}

type mockLog struct {
	t *testing.T
}

func (ml mockLog) Log(keys ...interface{}) error {
	ml.t.Log(keys)
	return nil
}
