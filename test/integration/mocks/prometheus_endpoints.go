//go:build integration_test

package mocks

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// StartExporter starts a server with metrics mocked.
func StartExporter(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()

	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		response := `
# HELP mock_gauge_metric A gauge test metric.
# TYPE mock_gauge_metric gauge
mock_gauge_metric 9
`
		_, _ = fmt.Fprintln(w, response)
	})

	mux.HandleFunc("/metrics-a/", func(w http.ResponseWriter, r *http.Request) {
		response := "custom_metric_a 46"
		_, _ = fmt.Fprintln(w, response)
	})

	mux.HandleFunc("/metrics-b/", func(w http.ResponseWriter, r *http.Request) {
		response := "custom_metric_b 88"
		_, _ = fmt.Fprintln(w, response)
	})

	mockExporterServer := httptest.NewServer(mux)

	t.Cleanup(func() {
		mockExporterServer.Close()
	})

	return mockExporterServer
}
