# Contributing to Kubo

## Contributor License Agreement
If you have not previously done so, please fill out and submit the [Contributor License Agreement](https://cla.pivotal.io/sign/pivotal). 

## Developer Workflow
1. Fork the project on [GitHub](https://github.com/pivotal-cf-experimental/kubo-deployment)
1. Create a feature branch.
1. Make your feature addition or bug fix. Please make sure there is appropriate test coverage.
1. Run tests.
1. Rebase on top of master.
1. Send a pull request.

Before making significant changes it's best to communicate with the maintainers of the project through [GitHub Issues](https://github.com/pivotal-cf-experimental/kubo-deployment/issues).

## Running Integration Tests
Before submitting pull request please redeploy KuBosh, deploy service on it and try workload.

## Additional BOSH configuration
We support only basic BOSH configuration. If you have some additional ops-file that will be useful for community, add them to https://github.com/cloudfoundry/bosh-deployment We have included this repo as a subtree and update it periodically.
