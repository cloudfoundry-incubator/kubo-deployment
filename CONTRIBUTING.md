# Contributing to Kubo Deployment

As a potential contributor, your changes and ideas are always welcome. Please do not hesitate to ask a question using GitHub issues or send a pull request to contribute changes.

## Contributor License Agreement
If you have not previously done so, please fill out and submit an [Individual Contributor License Agreement](https://www.cloudfoundry.org/governance/cff_individual_cla/) or a [Corporate Contributor License Agreement](https://www.cloudfoundry.org/governance/cff_corporate_cla/).

## Contributor Workflow
We encourage contributors to have discussion around design and implmentation with team members before making significant changes to the project through [GitHub Issues](https://github.com/cloudfoundry-incubator/kubo-deployment/issues). The project manager will prioritize where the feature will fit into the project's road map.

1. Fork the project on [GitHub](https://github.com/cloudfoundry-incubator/kubo-deployment)
1. Create a feature branch.
1. Make your feature addition or bug fix. Please make sure there is appropriate test coverage.
1. Run integration [tests](#running-integration-tests).
1. Ensure your feature branch is up to date with `master` branch.
1. Submit a pull request with clear description of intended change.
1. The team will triage the pull request and assign it to a team member.
1. A team member will approve the pull request or make suggestions for changes.

## Writing the tests

Unit Testing is the responsibility of all contributors.

### Shell script tests

If the PR incudes any change to bash script consider adding appropriate shell tests. Shell tests are located under `./src/kubo-deployment-tests` directory.

## Running Tests
### Pre-requisites

1. Install [Go](https://golang.org/doc/install)
1. Install [Ginkgo](https://onsi.github.io/ginkgo/)

Please make sure to run all tests before submitting a pull request.

### Shell script tests

This repository provides a test runner for running integration tests against the shell scripts. To run the tests, execute `./bin/run_tests` from the repository directory.

### Smoke test

Follow the steps below to deploy and test kubernetes BOSH deployment. This test verifies the basic functionality is not broken.

1. (Re-)deploy a BOSH director using `bin/deploy_bosh`.
1. Deploy a Kubernetes cluster using `bin/deploy_k8s`.
1. [Deploy a workload](https://kubernetes.io/docs/tasks/run-application/run-stateless-application-deployment/) on the cluster.

If you want to start fresh, it is possible to tear down the service by running `bosh -e <KUBO_ENV> -d <CLUSTER_DEPLOYMENT_NAME> delete-deployment`. You can then destroy the BOSH director using the `bin/destroy_bosh` command.

### Conformance tests

Follow the steps below to test your cluster against the [certified Kubernetes](https://github.com/cncf/k8s-conformance) conformance tests.  The instructions differ from the official kubernetes instructions and allow the tests to be run in parallel.  In order to submit to the Certified Kubernetes programme, you will have to follow the official instructions.

#### Prerequisites

Ensure you have a CFCR cluster and a Kubeconfig file.  The following instructions assume your Kubeconfig file is located at `${HOME}/.kube/config`.

#### Running the conformance tests

1. `docker run -it --rm --mount type=bind,source="${HOME}/.kube/config",target="/kubeconfig" pcfkubo/conformance:stable /bin/bash`
1. `ginkgo -p -progress -focus  "\[Conformance\]" -skip "\[Serial\]" /e2e.test`
1. `ginkgo -focus="\[Serial\].*\[Conformance\]" /e2e.test`
