# kubo-deployment

Kubo is a [BOSH](https://bosh.io/) release for Kubernetes. It provides a solution for deploying and managing Kubernetes with BOSH.

This repository contains the documentation and manifests for deploying [kubo-release](https://github.com/cloudfoundry-incubator/kubo-release) with BOSH.


**Slack**: #cfcr on https://slack.cloudfoundry.org
**Pivotal Tracker**: https://www.pivotaltracker.com/n/projects/2093412

## CI Status

Build Kubo Release status [![Build Kubo Release Badge](https://ci.kubo.sh/api/v1/pipelines/kubo-deployment/jobs/build-kubo-release/badge)](https://ci.kubo.sh/pipelines/kubo-deployment)

### IaaS specific jobs

| Job | GCP with CF routing pipeline Status |GCP with load balancer status|vSphere status|
|---------|--------|--------|--------|
| Install BOSH | [![BOSH GCP Badge](https://ci.kubo.sh/api/v1/pipelines/kubo-deployment/jobs/install-bosh-gcp/badge)](https://ci.kubo.sh/pipelines/kubo-deployment) | [![BOSH GCP LB Badge](https://ci.kubo.sh/api/v1/pipelines/kubo-deployment/jobs/install-bosh-gcp-lb/badge)](https://ci.kubo.sh/pipelines/kubo-deployment) | [![BOSH vSphere Badge](https://ci.kubo.sh/api/v1/pipelines/kubo-deployment/jobs/install-bosh-vsphere/badge)](https://ci.kubo.sh/pipelines/kubo-deployment) |
| Deploy K8s | [![Deploy K8s GCP Badge](https://ci.kubo.sh/api/v1/pipelines/kubo-deployment/jobs/deploy-k8s-gcp/badge)](https://ci.kubo.sh/pipelines/kubo-deployment) | [![Deploy K8s GCP LB Badge](https://ci.kubo.sh/api/v1/pipelines/kubo-deployment/jobs/deploy-k8s-gcp-lb/badge)](https://ci.kubo.sh/pipelines/kubo-deployment) | [![Deploy K8s vSphere Badge](https://ci.kubo.sh/api/v1/pipelines/kubo-deployment/jobs/deploy-k8s-vsphere/badge)](https://ci.kubo.sh/pipelines/kubo-deployment) |
| Run smoke tests | [![Run smoke tests GCP Badge](https://ci.kubo.sh/api/v1/pipelines/kubo-deployment/jobs/run-k8s-integration-tests-gcp/badge)](https://ci.kubo.sh/pipelines/kubo-deployment) | [![Run smoke tests GCP LB Badge](https://ci.kubo.sh/api/v1/pipelines/kubo-deployment/jobs/run-k8s-integration-tests-gcp-lb/badge)](https://ci.kubo.sh/pipelines/kubo-deployment) | [![Run smoke tests vSphere Badge](https://ci.kubo.sh/api/v1/pipelines/kubo-deployment/jobs/run-k8s-integration-tests-vsphere/badge)](https://ci.kubo.sh/pipelines/kubo-deployment) |

See the [complete pipeline](https://ci.kubo.sh/pipelines/kubo-deployment) for more details. The CI pipeline definitions are stored in the [kubo-ci](https://github.com/pivotal-cf-experimental/kubo-ci) repository.

## Documentation
To deploy CFCR go [here](https://github.com/cloudfoundry-incubator/kubo-release/#deploying-cfcr)

## Contributing

For instructions on contributing to this project, please see [CONTRIBUTING.md](CONTRIBUTING.md).

## Troubleshooting

Please refer to the [troubleshooting guide](https://docs-cfcr.cfapps.io/managing/troubleshooting/) to look for solutions to the most common issues.

## Design

### Components

A specialized BOSH director manages the virtual machines for the Kubo instance. This involves VM creation, health checking, and resurrection of missing or unhealthy VMs. The BOSH director includes CredHub and PowerDNS to handle certificate generation within the kubo clusters. Additionally, Credhub is used to store the auto-generated passwords.

### Networking Topology - using IaaS Load Balancers

![Diagram describing how traffic is routed to Kubo](docs/images/cfcr-architecture.png)

The nodes that run the Kubernetes API (master nodes) are exposed through an IaaS specific load balancer. The load balancer will have an external static IP address that is used as a public and internal endpoint for traffic to the Kubernetes API.

Kubernetes services can be exposed using a second IaaS specific load balancer which forwards traffic to the Kubernetes worker nodes.

### Networking Topology - using Cloud Foundry routing

![Diagram describing how traffic is routed to Kubo using CF](docs/images/kubo-network-cf.png)

The nodes that run the Kubernetes API (master nodes) register themselves with the Cloud Foundry TCP router. The TCP Router acts as both public and internal endpoint for the Kubernetes API to route traffic to the master nodes of a Kubo instance. All traffic to the API goes through the Cloud Foundry TCP router and then to a healthy node. 

The Cloud Foundry subnet must be able to route traffic directly to the Kubo subnet. It is recommended to keep them in separate subnets when possible to avoid the BOSH directors from trying to provision the same addresses. This diagram specifies CIDR ranges for demonstration purposes as well as a public router in front of the Cloud Foundry gorouter and tcp-router which is typical.
