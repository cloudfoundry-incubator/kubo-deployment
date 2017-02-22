# Kubo deployment

This repository contains automation scripts to deploy [Kubernetes](https://kubernetes.io/) cluster using
[BOSH](https://bosh.io/) within a [CloudFoundry](https://cloudfoundry.org) infrastructure.

# Required software

You'll need the following tools on your machine to continue:

- [bosh-cli](https://bosh.io/docs/cli-v2.html) for interacting with BOSH. Version 2.0.1 and above. Please ensure the binary is installed as `bosh-cli` and not `bosh`.
- [credhub cli](https://github.com/pivotal-cf/credhub-cli/releases/tag/0.4.0) for interacting with CredHub. Version 0.4 only.
- [Ruby 2.3+](https://www.ruby-lang.org/en/downloads) required by the bosh-cli to deploy BOSH++
- [make](https://www.gnu.org/software/make) required by the bosh-cli to deploy BOSH++

# Installation steps

Here is an overview of the steps you would need to take in order to get a running cluster

1. [Setup IAAS](#Setup-IAAS)
1. [Setup BOSH configuration files](#Setup-BOSH-configuration-files)
1. [Deploy BOSH++](#Deploy-BOSH++)
1. [Deploy CF](#Cloud-Foundry)
1. [Deploy K8s](#Deploy-Kubo)
1. Push your app!

## Setup IAAS

You will need to setup your IAAS before proceeding to installing BOSH. Check [here for instructions on how to setup GCP](https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/tree/master/docs/bosh#configure-your-google-cloud-platform-environment) and [here for all other IAASes](https://bosh.io/docs/init.html). Please only follow the steps to pave your IAAS.

## Setup BOSH configuration files

### Generate configuration

The configuration can be generated using `bin/generate_bosh_config`.
It requires a path to a directory where the configuration wil be stored, the name of the environment,
and the IaaS name.

It will create a directory with the same name as the environment at the specified path, containing three files:
- `iaas` which contains IaaS name
- `director.yml` which contains public BOSH director, IaaS and network configurations
- `director-secrets.yml` which contains secrets, such as OpenStack password

### Fill in configuration

The format for all configuration files is YAML. The main configuration is stored in a `director.yml` file. All
the properties have comments explaining their possible values and their purpose.

Refer to `ci/environments/gcp/director.yml` as an example.

## Deploy BOSH++

### What is BOSH++?

BOSH with integrated PowerDNS and CredHub. It is used to auto-generate certificates for kubelets.
You can generate certificate for kubelets manually and use them as part of deployment. In that case you can use
regular BOSH.

### Deployment process

When the environment preparation and configuration is completed, `BOSH++` can be
easily deployed with a single command:

```bash
bin/deploy_bosh <path to configuration> <private or service account key filename for BOSH to use for deployments>

```

As a pre-requisite, you will need to ensure that the machine you are running the command from has network access to the BOSH director. Else you may get the error

```
```

There are multiple ways to ensure access. A couple of examples
1. Run the command on a Bastion VM
1. Use [sshuttle](https://github.com/apenwarr/sshuttle) to get access to VMs on the network

During the deployment, all the passwords and SSL certificates will be automatically
generated. Most of them will be saved into the configuration path in a file called
`creds.yml`. Because this file will contain sensitive information, it is not recommended
to store it in a VCS. This file is also required to successfully deploy the kubernetes
service or the on-demand broker.

Subsequent runs of `bin/bosh_deploy` will apply changes made to the configuration
to an already existing BOSH++ installation, reusing the secrets stored in the `creds.yml`.

Another file that gets created during initial deployment is called `state.json`. It contains
the BOSH state identical to the one used by [bosh-init](https://bosh.io/docs/using-bosh-init.html).

Additionally, the deployment script creates the `default` CA certificate within CredHub.

## Cloud Foundry

Cloud Foundry deployment is not covered in this document. You can use an existing CF deployment you may have, or deploy using instructions.

Please ensure that instances of your CF deployment are able to communicate with other deployments deployed via BOSH++. This could simply mean deploying CF in the same network as K8s, or whitelisting traffic.

TCP router and the Routing API should be enabled on the Cloud Foundry installation. See [OSS CF](https://docs.cloudfoundry.org/adminguide/enabling-tcp-routing.html) or
[PCF](http://docs.pivotal.io/pivotalcf/1-8/opsguide/tcp-routing-ert-config.html) documentation for
further details.

A UAA client with [appropriate authorities](https://github.com/cloudfoundry-incubator/routing-api#configure-oauth-clients-manually-using-uaac-cli-for-uaa)
is required in order to register the TCP routes.

## Deploy Kubo

### The Easy Way: Automated

Once BOSH++ is deployed, the Kubernetes BOSH release can be built and deployed with this command:

```bash
bin/deploy_k8s <BOSH_ENV> <DEPLOYMENT_NAME> <RELEASE_SOURCE>
```

The `RELEASE_SOURCE` parameter allows you to either build and deploy a local copy of the repository, or deploy our pre-built kubo-release tarball, and is optional.

Note that the scripts will:

- Generate certificate in CredHub
- upload any Cloud Config changes to the director
- create the kubo release from source and upload it if RELEASE_SOURCE is dev
or
- upload the kubo release tarball from specified location if RELEASE_SOURCE is local
- regenerate the deployment manifest
- kick off the deployment using `bosh_admin` UAA client

By default, the deployment will use the latest versions of the releases. If releases were uploaded from different machines or
used different sources, deployment might use wrong release.

### The Hard Way: Step by Step

#### Generate Cloud Config

Kubo deployment uses [BOSH 2.0 Cloud Config](https://bosh.io/docs/cloud-config.html).

The default cloud config for GCP uses n1-standard-1 VMs for supporting services and n1-standard-2 VMs
for Kubernetes workers. The network properties are pulled in from the environment configuration file
which is stored at `<BOSH_ENV>/director.yml`.

Cloud config can be generated using following command:
```bash
bin/generate_cloud_config <BOSH_ENV>
```

##### Create and upload release

To create dev release, download the [kubo-release repository](https://github.com/pivotal-cf-experimental/kubo-release)
and follow [documentation](https://bosh.io/docs/create-release.html#dev-release)

##### Generate manifest

Manifest can be generated using command `bin/generate_service_manifest <BOSH_ENV> <DEPLOYMENT_NAME>`

##### Deploy

Run deployment using the following command: `bosh-cli -e <BOSH_ALIAS> -d <DEPLOYMENT_NAME> deploy <PATH to MANIFEST>`
where `<BOSH_ALIAS>` is BOSH director name from configuration file or BOSH director address

`bosh-cli` has to be authenticated to BOSH director

## Accessing Kubo

Configure kubectl for your Kubo instance with the following command:

```bash
bin/set_kubeconfig <BOSH_ENV> <DEPLOYMENT_NAME>
```

You can now issue kubectl commands such as:
```bash
kubectl get pods --namespace=kube-system
kubectl get nodes
```

## Development

This repo uses https://github.com/cloudfoundry/bosh-deployment as git subtree. Run command ` git subtree pull --prefix bosh-deployment git@github.com:cloudfoundry/bosh-deployment.git <ref>` to update subtree.
