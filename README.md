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