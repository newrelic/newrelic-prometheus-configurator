# Kubernetes target filtering

The config parameter `target_discovery.filter` should be used to filter in the targets that Prometheus will scrape. Current conditions are filtering by `label` and `annotation`. All conditions are applied using an `AND` operation.

Following example will scrape only `Pods` having the `newrelic.io/scrape: true` annotation and a label `k8s.io/app` with value `postgres` or `mysql`.

Note that if not value is specified for the label or annotation, the filter will only check that exists.

``` yaml
kubernetes:
  jobs:
  - job_name_prefix: example
    target_discovery: 
      pod: true
      filter:
        annotation:
          # <string>: <regex>
          newrelic.io/scrape: 'true'
        label:
          # <string>: <regex>
          k8s.io/app: '(postgres|mysql)'
```


