# kubo-deployment

Kubo is a [BOSH](https://bosh.io/) release for Kubernetes. It provides a solution for deploying and managing Kubernetes with BOSH

This repository contains the documentation and manifests for deploying [kubo-release](https://github.com/pivotal-cf-experimental/kubo-release) with BOSH.

## CI Status

Build Kubo Release status [![Build Kubo Release Badge](https://ci.kubo.sh/api/v1/teams/main/pipelines/kubo-deployment/jobs/build-kubo-release/badge)](https://ci.kubo.sh/teams/main/pipelines/kubo-deployment)

### IaaS specific jobs

| Job | GCP with CF routing pipeline Status |GCP with load balancer status|vSphere status|
|---------|--------|--------|--------|
| Install KuBOSH | [![KuBOSH GCP Badge](https://ci.kubo.sh/api/v1/teams/main/pipelines/kubo-deployment/jobs/install-bosh-gcp/badge)](https://ci.kubo.sh/teams/main/pipelines/kubo-deployment) | [![KuBOSH GCP LB Badge](https://ci.kubo.sh/api/v1/teams/main/pipelines/kubo-deployment/jobs/install-bosh-gcp-lb/badge)](https://ci.kubo.sh/teams/main/pipelines/kubo-deployment) | [![KuBOSH vSphere Badge](https://ci.kubo.sh/api/v1/teams/main/pipelines/kubo-deployment/jobs/install-bosh-vsphere/badge)](https://ci.kubo.sh/teams/main/pipelines/kubo-deployment) |
| Deploy K8s | [![Deploy K8s GCP Badge](https://ci.kubo.sh/api/v1/teams/main/pipelines/kubo-deployment/jobs/deploy-k8s-gcp/badge)](https://ci.kubo.sh/teams/main/pipelines/kubo-deployment) | [![Deploy K8s GCP LB Badge](https://ci.kubo.sh/api/v1/teams/main/pipelines/kubo-deployment/jobs/deploy-k8s-gcp-lb/badge)](https://ci.kubo.sh/teams/main/pipelines/kubo-deployment) | [![Deploy K8s vSphere Badge](https://ci.kubo.sh/api/v1/teams/main/pipelines/kubo-deployment/jobs/deploy-k8s-vsphere/badge)](https://ci.kubo.sh/teams/main/pipelines/kubo-deployment) |
| Run smoke tests | [![Run smoke tests GCP Badge](https://ci.kubo.sh/api/v1/teams/main/pipelines/kubo-deployment/jobs/deploy-workload-gcp/badge)](https://ci.kubo.sh/teams/main/pipelines/kubo-deployment) | [![Run smoke tests GCP LB Badge](https://ci.kubo.sh/api/v1/teams/main/pipelines/kubo-deployment/jobs/deploy-workload-gcp-lb/badge)](https://ci.kubo.sh/teams/main/pipelines/kubo-deployment) | [![Run smoke tests vSphere Badge](https://ci.kubo.sh/api/v1/teams/main/pipelines/kubo-deployment/jobs/deploy-workload-vsphere/badge)](https://ci.kubo.sh/teams/main/pipelines/kubo-deployment) |

See the [complete pipeline](https://ci.kubo.sh/teams/main/pipelines/kubo-deployment) for more details. The CI pipeline definitions are stored in the [kubo-ci](https://github.com/pivotal-cf-experimental/kubo-ci) repository.

## Table of Contents

- [Design](#design)
- [Glossary](#glossary)
- [Installation Guides](#installation)
- [Troubleshooting](#troubleshooting)
- [Contribution](#contributing)

## Design

### Components

A specialized BOSH director manages the virtual machines for the Kubo instance. This involves VM creation, health checking, and resurrection of missing or unhealthy VMs. The BOSH director includes CredHub and PowerDNS to handle certificate generation within the kubo clusters. Additionally, Credhub is used to store the auto-generated passwords.

### Networking Topology - using IaaS Load Balancers

![Diagram describing how traffic is routed to Kubo](docs/images/kubo-network.png)

The nodes that run the Kubernetes API (master nodes) are exposed through an IaaS specific load balancer. The load balancer will have an external static IP address that is used as a public and internal endpoint for traffic to the Kubernetes API.

Kubernetes services can be exposed using a second IaaS specific load balancer which forwards traffic to the Kubernetes worker nodes.

### Networking Topology - using Cloud Foundry routing

![Diagram describing how traffic is routed to Kubo using CF](docs/images/kubo-network-cf.png)

The nodes that run the Kubernetes API (master nodes) register themselves with the Cloud Foundry TCP router. The TCP Router acts as both public and internal endpoint for the Kubernetes API to route traffic to the master nodes of a Kubo instance. All traffic to the API goes through the Cloud Foundry TCP router and then to a healthy node. 

The Cloud Foundry subnet must be able to route traffic directly to the Kubo subnet. It is recommended to keep them in separate subnets when possible to avoid the BOSH directors from trying to provision the same addresses. This diagram specifies CIDR ranges for demonstration purposes as well as a public router in front of the Cloud Foundry gorouter and tcp-router which is typical.

## Glossary

- Kubo - Kubernetes on BOSH
- KuBOSH - BOSH with UAA, Credhub and PowerDNS
- [Bastion](https://en.wikipedia.org/wiki/Jump_server) - A server within the kubo network that provides secure access to kubo.
- BOSH environment Configuration - Folder that contains all configuration files needed to deploy KuBOSH and Kubo, as well as all 
  configuration files that are generated during deployment. Also called `<BOSH_ENV>`
- Creds - Credentials that are generated during KuBOSH deployment process and stored in `<BOSH_ENV>/creds.yml`
- Service - stands for [K8s service](https://kubernetes.io/docs/user-guide/services), which represents a logical collection 
  of Kubernetes pods and a way to access them without needing information about the specific pods

## Installation

Please choose the guide below that matches your requirements

1. Deploy Kubo from scratch on Google Cloud Platform - [guide](docs/guides/gcp)
1. Deploy Kubo using an existing Cloud Foundry installation for routing on Google Cloud Platform - [guide](docs/guides/gcp-cf)
1. Deploy Kubo step by step, allowing for customization - [guide](docs/guides/customized-installation.md)

Once the kubernetes is installed, read the [guide](../using-kubernetes.md) on setting up
`kubectl` and deploying services to the new cluster.

## Delete resources

### Delete Kubernetes Cluster

You can use the BOSH cli to delete your kubernetes deployment

```
bosh-cli -e kube -d kube delete-deployment
```

### Delete KuBOSH Director

Use the following script to delete your KuBOSH director

```
bin/destroy_bosh ~/kubo-env/kube ~/kubo-env/kube/service_account.json
```

## Troubleshooting

Please refer to the [troubleshooting guide](docs/troubleshooting.md) to look for solutions to the most common issues. 

## Contributing

For instructions on contributing to this project, please see [CONTRIBUTING.md](CONTRIBUTING.md).
