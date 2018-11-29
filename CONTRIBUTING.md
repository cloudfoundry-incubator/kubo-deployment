# Contributing to Kubo Deployment

As a potential contributor, your changes and ideas are always welcome.

## Contributor License Agreement
If you have not previously done so, please fill out and submit an [Individual Contributor License Agreement](https://www.cloudfoundry.org/governance/cff_individual_cla/) or a [Corporate Contributor License Agreement](https://www.cloudfoundry.org/governance/cff_corporate_cla/).

## Contributing Instructions

Please refer to the [Contributing guide](https://github.com/cloudfoundry-incubator/kubo-release/blob/master/CONTRIBUTING.md) in the Kubo-Release repository.

## Running Unit Tests for Kubo-Deployment
### Pre-requisites

1. Install [Ruby](https://www.ruby-lang.org/en/documentation/installation/)
1. Install Bundler: `gem install bundler`

## Running Tests

### Manifest Operation file tests

CFCR contains special tests to verify if Operation files were documented and that it is possible to apply each operation file to the manifest.
The tests require only Bosh CLI.

1. Run `./bin/run_tests`

### Smoke test

Follow the steps below to deploy and test kubernetes BOSH deployment. This test verifies the basic functionality is not broken.

1. Deploy a Kubernetes cluster.
1. Run the `apply-addons` errand.
1. Run `smoke-tests` errand.

### Conformance tests

Follow the steps below to test your cluster against the [certified Kubernetes](https://github.com/cncf/k8s-conformance) conformance tests.  The instructions differ from the official kubernetes instructions and allow the tests to be run in parallel.  In order to submit to the Certified Kubernetes programme, you will have to follow the official instructions.

#### Running tests using CFCR Docker file

You can run Conformance tests using Docker file, it is much faster, but can be flaky.

1. `docker run -it --rm --mount type=bind,source="${HOME}/.kube/config",target="/kubeconfig" pcfkubo/conformance:stable /bin/bash`
1. `ginkgo -p -progress -focus  "\[Conformance\]" -skip "\[Serial\]" /e2e.test`
1. `ginkgo -focus="\[Serial\].*\[Conformance\]" /e2e.test`