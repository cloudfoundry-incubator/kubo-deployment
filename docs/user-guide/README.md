# Kubo user guide

## Infrastructure paving
Kubo can be deployed on various infrastructure providers. The currently supported ones 
are GCP and vSphere. vSphere installation is not yet covered in this guide. 

Please follow the link for infrastructure paving on your platform:
- [Google Cloud Platform](platforms/gcp/paving.md)

## Deploying BOSH

### Deploy BOSH using kubo-deployment

- [Google Cloud Platform](platforms/gcp/install-bosh.md)

### Custom BOSH deployment

In order to deploy Kubo you need a BOSH Director with the following 
releases installed: UAA, CredHub 1.0+, PowerDNS. 

## Configure Kubo

### Basic configuration

If the bosh director has been deployed using `kubo-deployment`,
skip to the [Routing options](#routing-options) section.

Create a [kubo environment](create-kubo-env.md) and fill
in the network and IaaS-related properties, as well as all the details
for the running BOSH.

### Routing options

Kubo can leverage different routing options to improve security and high
availability of the cluster. Please configure the kubo environment according
to one of the options below:

- [CF Routing](routing/cf.md)
- [IaaS Load-Balancing](routing/gcp/load-balancing.md)

## Deploying Kubo

## Accessing Kubernetes               