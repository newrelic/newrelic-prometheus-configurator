[![Community Project header](https://github.com/newrelic/opensource-website/raw/master/src/images/categories/Community_Project.png)](https://opensource.newrelic.com/oss-category/#community-project)

# [Name of Project] [build badges go here when available]

>[Brief description - what is the project and value does it provide? How often should users expect to get releases? How is versioning set up? Where does this project want to go?]

## Installation

> [Include a step-by-step procedure on how to get your code installed. Be sure to include any third-party dependencies that need to be installed separately]

## Getting Started
>[Simple steps to start working with the software similar to a "Hello World"]

## Usage
>[**Optional** - Include more thorough instructions on how to use the software. This section might not be needed if the Getting Started section is enough. Remove this section if it's not needed.]


## Building

>[**Optional** - Include this section if users will need to follow specific instructions to build the software from source. Be sure to include any third party build dependencies that need to be installed separately. Remove this section if it's not needed.]

## Testing

### Running unit tests

```bash
make test
```

### Running integration tests

Detailed info about integration tests [here](./test/integration/README.md).

## Run local environment

We use minikube and Tilt to launch a local cluster and deploy the [main chart](charts/newrelic-prometheus/) and a set of testing endpoints from the [test-resource](charts/internal/test-resources/).

Make sure you have these tools or install them:
- [Install minikube](https://minikube.sigs.k8s.io/docs/start/)
- [Install Tilt](https://docs.tilt.dev/install.html)
- [Install Helm](https://helm.sh/docs/intro/install/)

A license key and cluster name are required to run the environment. Configure them by the environment variables `NR_PROM_CLUSTER` and `NR_PROM_LICENSE_KEY`.

Start the local environment:
```shell
export NR_PROM_CLUSTER=<cluster name>
export NR_PROM_LICENSE_KEY=<NewRelic ingest key>
make start-local-cluster
make tilt-up
```

Notice that local images are build and pushed to docker running inside the minikube cluster since we are running `eval $(minikube docker-env)` before launching Tilt.

### Running e2e tests

This test are based on the `newrelic-integration-e2e`. This tool will start the local environment and check if the expected metrics has been reach the NewRelic platform.

To run it locally you should already be able to run the local cluster and apart from that you must install the `newrelic-integration-e2e` binary:
```bash
git clone https://github.com/newrelic/newrelic-integration-e2e-action
cd newrelic-integration-e2e-action/newrelic-integration-e2e
go build -o  $GOPATH/bin/newrelic-integration-e2e ./cmd/...
```

Then run:
```bash
export ACCOUNT_ID=<NewRelic account number>
export API_REST_KEY=<NewRelic api rest key>
export LICENSE_KEY=<NewRelic ingest key>
make start-local-cluster
make e2e-test
```

## Support

New Relic hosts and moderates an online forum where customers can interact with New Relic employees as well as other customers to get help and share best practices. Like all official New Relic open source projects, there's a related Community topic in the New Relic Explorers Hub. You can find this project's topic/threads here:

>Add the url for the support thread here: discuss.newrelic.com

## Contribute

We encourage your contributions to improve [project name]! Keep in mind that when you submit your pull request, you'll need to sign the CLA via the click-through using CLA-Assistant. You only have to sign the CLA one time per project.

If you have any questions, or to execute our corporate CLA (which is required if your contribution is on behalf of a company), drop us an email at opensource@newrelic.com.

**A note about vulnerabilities**

As noted in our [security policy](../../security/policy), New Relic is committed to the privacy and security of our customers and their data. We believe that providing coordinated disclosure by security researchers and engaging with the security community are important means to achieve our security goals.

If you believe you have found a security vulnerability in this project or any of New Relic's products or websites, we welcome and greatly appreciate you reporting it to New Relic through [HackerOne](https://hackerone.com/newrelic).

If you would like to contribute to this project, review [these guidelines](./CONTRIBUTING.md).

To all contributors, we thank you!  Without your contribution, this project would not be what it is today.  We also host a community project page dedicated to [Project Name](<LINK TO https://opensource.newrelic.com/projects/... PAGE>).

## License
[Project Name] is licensed under the [Apache 2.0](http://apache.org/licenses/LICENSE-2.0.txt) License.
>[If applicable: The [project name] also uses source code from third-party libraries. You can find full details on which libraries are used and the terms under which they are licensed in the third-party notices document.]
