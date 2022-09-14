# Horizontal Scaling

Horizontal Scaling is supported by setting up a configuration parameter which allows running several prometheus servers in agent mode to gather your data.

If you define `sharding.total_shards_count` value, the deployed StatefulSet will include as many replicas as you defined there. When it is used, the _configurator_ component will automatically include some additional relabel rules so each target will only be scraped by one prometheus server. Those
rules rely on the target's address [hash-mod](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#relabel_config).

For example, if `custom-values.yaml` includes:

```yaml
# (...)
sharding:
  total_shards_count: 5
# (...)
```

And then, the release is upgraded:

```shell
helm upgrade my-prometheus-release newrelic-prometheus-configurator/newrelic-prometheus -f custom-values.yaml
```

Then five prometheus servers will be executed and each target will only be scraped by only one of them.

## Self metrics

Usually, prometheus server self-metrics should be gathered from all prometheus servers, so the additional rules when sharding is configured should not apply to the job gathering the prometheus self-metrics. This is possible because the _configurator_ accepts the flag `skip_sharding` in the static_target jobs. This parameter is already set up in the default self-metrics job.


## Limitations

If additional scrape jobs are included in configuration as `extra_scrape_configs`, as that field will hold the raw definition of prometheus jobs, the _configurator_ will not include the rules corresponding to sharding configuration and, as a result, the corresponding targets will be scrapped by all
prometheus servers.

Currently auto-scaling is not supported. If increasing or decreasing the number of shards is required, an update of the cluster settings (which implies restarting the prometheus servers) is needed.
