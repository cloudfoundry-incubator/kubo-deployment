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

#### Without internet access

Perform the section above on your workstation and copy the blobs and releases to your `bosh-bastion` via SCP. Then update your director.yml config file to set the `etcd_release_url` and `kubo_release_url` fields to link to the local copies of those tarballs. The following example is for GCP and will be different on other providers.

```bash
# from your workstation
gcloud compute copy-files "~/workspace/kubo-release/blobs/*" "$BOSH_ENV-bosh-bastion:kubo-release/blobs"
gcloud compute copy-files "~/workspace/kubo-service-adapter-release/blobs/*" "$BOSH_ENV-bosh-bastion:kubo-release/blobs"
gcloud compute copy-files "~/workspace/etcd.tgz" "$BOSH_ENV-bosh-bastion:/home/username/workspace/"
gcloud compute copy-files "~/workspace/kubo-release.tgz" "$BOSH_ENV-bosh-bastion:/home/username/workspace/"
```
