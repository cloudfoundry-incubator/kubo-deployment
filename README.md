# Kubo deployment

This repository contains automation scripts to deploy [Kubernetes](https://kubernetes.io/) cluster using 
[BOSH](https://bosh.io/) within a [CloudFoundry](https://cloudfoundry.org) infrastructure.

# Required software

You'll need the following tools on your machine to continue:

- [bosh-cli](https://bosh.io/docs/cli-v2.html) for interacting with BOSH. Version 0.142 and above.
- [credhub cli](https://github.com/pivotal-cf/credhub-cli) for interacting with CredHub. Version 0.4 only.

# Deploy

The scripts require custom config. 

## Environment preparation

The environment preparation is out of scope. Check [here](https://bosh.io/docs/init.html) for more documentation.

## CloudFoundry

CloudFoundry deployment is out of scope. 

## Generate configuration

The configuration can be generated using `bin/generate_bosh_config`.
It requires a path to a directory where the configuration wil be stored, the name of the environment,
and the IaaS name.

It will create a directory with the same name as the environment at the specified path, containing three files: 
- `iaas` which contains IaaS name
- `director.yml` which contains public BOSH director, IaaS and network configurations 
- `director-secrets.yml` which contains secrets, such as OpenStack password

## Fill in configuration

The format for all configuration files is YAML. The main configuration is stored in a `director.yml` file. All 
the properties have comments explaining their possible values and their purpose.

## Deploy BOSH++

### What is BOSH++? 

BOSH with integrated PowerDNS and CredHub. It is used to auto-generate certificates for kubelets. 
You can generate certificate for kubelets manually and use them as part of deployment. In that case you can use 
regular BOSH.

### Deployment process

When the environment preparation and configuration is completed, `BOSH++` can be
easily deployed with a single command:

```bash
bin/deploy_bosh <path to configuration>
```

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

## Deploy Kubo

### Deployment Process

Once BOSH++ is deployed, the Kubernetes BOSH release can be built and deployed with this command:

```bash
bin/deploy_k8s <BOSH_ENV> <DEPLOYMENT_NAME> <RELEASE_SOURCE>
```

The `RELEASE_SOURCE` parameter allows you to either build and deploy a local copy of the repository, or deploy our pre-built kubo-release tarball.

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

### Deployment using separate scripts

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
