# Contributing to Kubo

## Contributor License Agreement
If you have not previously done so, please fill out and submit an [Individual Contributor License Agreement](https://www.cloudfoundry.org/governance/cff_individual_cla/) or a [Corporate Contributor License Agreement](https://www.cloudfoundry.org/governance/cff_corporate_cla/). 

## Prerequisites
Review the following to understand Kubernetes
 1. https://docs.cloudfoundry.org/concepts/architecture/
 1. https://thenewstack.io/kubernetes-an-overview/
 1. https://github.com/kubernetes/community/blob/master/contributors/design-proposals/architecture.md

## Developer Workflow
1. Fork the project on [GitHub](https://github.com/cloudfoundry-incubator/kubo-deployment)
1. Create a feature branch.
1. Make your feature addition or bug fix. Please make sure there is appropriate test coverage.
1. Run tests.
1. Rebase on top of master.
1. Send a pull request.

Before making significant changes it's best to communicate with the maintainers of the project through [GitHub Issues](https://github.com/cloudfoundry-incubator/kubo-deployment/issues).

## Running Integration Tests

Please make sure to run all tests before submitting a pull request.

### Shell script tests

This repo provides a test runner for running integration tests against the shell scripts. It requires [`ginkgo`](https://github.com/onsi/ginkgo) binary to be installed locally. To run the tests, execute `bin/run_tests` from the repository directory.

### Deployment tests

The sequence to run deployment tests includes the following steps:

>  1. (re-)deploy KuBOSH
>  1. deploy a kubernetes cluster on the new KuBOSH
>  1. deploy a workload on the cluster and make sure it is working

Optionally, it is possible to tear down the service by running `bosh -e <KUBO_ENV> -d <CLUSTER_DEPLOYMENT_NAME> delete-deployment` followed by the `bin/destroy_bosh` command.

## Additional BOSH configuration
We support only basic BOSH configuration. If you have some additional ops-file that will be useful for community, add them to https://github.com/cloudfoundry/bosh-deployment We have included this repo as a subtree and update it periodically.
