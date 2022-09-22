# Target scrape interval

`common.scrape_interval` defines the default value for all jobs in the configuration. This can be overwritten by any job if `scrape_interval` is also defined in that job.

The following example shows two Kubernetes jobs with different scrape intervals.
``` yaml
common:
  scrape_interval: 30s
kubernetes:
  jobs:
  # this job will use the default scrape_interval defined in common.
  - job_name_prefix: default-targets-with-30s-interval
    target_discovery: 
      pod: true
      filter:
        annotation:
          newrelic.io/scrape: 'true'

  - job_name_prefix: slow-targets-with-60s-interval
    scrape_interval: 60s
    target_discovery: 
      pod: true
      filter:
        annotation:
          newrelic.io/scrape_slow: 'true'
```

