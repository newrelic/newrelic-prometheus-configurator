# Example Instructions

The yaml files in this example folder can be used to test out New Relic's Prometheus-based Kubernetes monitoring quickstart.

## If you are using the newrelic-bundle chart:

If you're installing `newrelic-prometheus-configurator` as part of the Kubernetes integration package with the [newrelic-bundle](https://github.com/newrelic/helm-charts/tree/master/charts/nri-bundle) chart, use the values from [newrelic-bundle-values.yaml](newrelic-bundle-values.yaml) to configure the chart bundle. For example, to upgrade an existing installation with just the required values:

```
curl -O https://raw.githubusercontent.com/newrelic/newrelic-prometheus-configurator/main/examples/kubernetes/newrelic-bundle-values.yaml
helm upgrade --reuse-values newrelic-bundle newrelic/nri-bundle -n newrelic -f newrelic-bundle-values.yaml
```

To add these values to a new installation, you can include them in your complete `values.yaml` file for the `newrelic-bundle` Helm chart. Or, if you're using the Helm command from the guided install wizard, you can add `-f newrelic-bundle-values.yaml` to the command to include those chart values with the new installation.


## If you are using the newrelic-prometheus-configurator chart only:

1. Deploy the following yamls. Note that if you already have KSM setup in your cluster, not all of these yamls will be necessary.

```
kubectl apply -f cluster-role.yaml
kubectl apply -f ksm-cluster-role-binding.yaml
kubectl apply -f ksm-cluster-role.yaml
kubectl apply -f ksm-deployment.yaml
kubectl apply -f ksm-service-account.yaml
kubectl apply -f ksm-service.yaml
```

2. Modify the `configurator-values.yaml` file to add your cluster name and license key.

3. Install the helm chart:

```
helm repo add newrelic-prometheus https://newrelic.github.io/newrelic-prometheus-configurator
helm upgrade --install newrelic newrelic-prometheus/newrelic-prometheus-agent -f configurator-values.yaml
```

4. Install the [Kubernetes Prometheus Quickstart dashboard](https://newrelic.com/instant-observability/kubernetes-prometheus) and watch for data to appear.
