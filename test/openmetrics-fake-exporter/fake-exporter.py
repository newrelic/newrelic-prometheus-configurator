import random
import time
import prometheus_client as prometheus

# Create a metric to track time spent and requests made.
METRIC_HISTOGRAM = prometheus.Histogram('fake_openmetrics_histogram', 'Histogram metric name')
METRIC_SUMMARY = prometheus.Summary('fake_openmetrics_summary', 'Summary metric name')
METRIC_GAUGE = prometheus.Gauge('fake_openmetrics_gauge', 'Gauge metric name')
METRIC_COUNTER = prometheus.Counter('fake_openmetrics_counter_total', 'Counter metric name')
METRIC_INFO = prometheus.Info('fake_openmetrics', 'Info metric name')
METRIC_STATESET = prometheus.Enum('fake_openmetrics_stateset', 'Description of enum',labelnames=['foo'], states=['starting', 'running', 'stopped'])

# Decorate function with metric.
@METRIC_HISTOGRAM.time()
@METRIC_SUMMARY.time()
def process_request(t):
    """A dummy function that takes some time."""
    time.sleep(t)

if __name__ == '__main__':
    # Start up the server to expose the metrics.
    prometheus.start_http_server(8000)
    # Generate some requests.

    METRIC_INFO.info({'version': '1.2.3', 'buildhost': 'foo@bar'})

    METRIC_STATESET.labels('bar').state('running')

    while True:
        process_request(random.random())
        METRIC_GAUGE.set(random.random())
        METRIC_COUNTER.inc()
