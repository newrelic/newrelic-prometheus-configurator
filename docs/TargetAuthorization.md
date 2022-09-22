# Target authorization configuration.

All authorization methods supported by Prometheus can be configured on `static_targets` and `kubernetes` jobs.

The supported methods are:
- [TLS](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#tls_config)
- [OAuth2](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#oauth2)
- [Authorization Header](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config)
  - Basic Auth

Below there are some examples:
``` yaml
kubernetes:
  jobs:
  - job_name_prefix: skip-verify-on-https-targets
    target_discovery: 
      pod: true
      filter:
        annotation:
          newrelic.io/scrape: 'true'
  - job_name_prefix: bearer-token
    target_discovery: 
      pod: true
      filter:
        label:
          k8s.io/app: my-app-with-token
    authorization:
      type: Bearer
      credentials_file: '/etc/my-app/token'

startic_targets:
  jobs:
  - job_name: mtls-target
    scheme: https
    targets:
    - 'my-mtls-target:8181'
    tls_config:
      ca_file: '/etc/my-app/client-ca.crt'
      cert_file: '/etc/my-app/client.crt'
      key_file: '/etc/my-app/client.key'
  
  - job_name: basic-auth-target
    targets:
    - 'my-basic-auth-static:8181'
    basic_auth:
      password_file: '/etc/my-app/pass.htpasswd'
```

