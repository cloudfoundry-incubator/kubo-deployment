# Kubo User Guide

**The Kubo User Guide is in the process of being deprecated in favor of the new Kubo documentation.**

## View the New Docs

To view the new Kubo documentation, visit https://docs-kubo.cfapps.io/, or jump to one of the following sections:

  * [Overview](https://docs-kubo.cfapps.io/): Learn what Kubo is and how it works.
  * [Installing and Configuring](https://docs-kubo.cfapps.io/installing/): Install Kubo on Google Cloud Platform (GCP), vSphere, Amazon Web Services (AWS), or OpenStack. 
  * [Managing and Troubleshooting](https://docs-kubo.cfapps.io/managing/): Manage and troubleshoot your Kubo deployment.

## Contribute to the New Docs

To learn how to contribute to the new Kubo documentation, see the README in https://github.com/cloudfoundry/docs-kubo.

## View the Old Docs

Because the new Kubo documentation is currently in development, it does not yet contain all of the information required to install, manage, and troubleshoot Kubo. The following topics are still included in this repo:

* [CF Routing](routing/cf.md): Configure Cloud Foundry to handle routing for your Kubo deployment.
* [Customized Kubo installation](customized-kubo-installation.md): Perform a customized Kubo installation that modifies the default behavior. 

See below for more information that will soon be incorporated into the Kubo documentation.

### Destroying Kubo 

#### Destroying the Cluster

Use the BOSH CLI if you want to destroy the cluster:

```bash
bosh-cli -e <KUBO_ENV> login
bosh-cli -e <KUBO_ENV_NAME> -d ${cluster_name} delete-deployment
```

Your username is admin and your password is the `admin_password` field in `<KUBO_ENV>/creds.yml`.

`KUBO_ENV_NAME` was set up in the Install BOSH step.

#### Destroying the BOSH Environment

To destroy your BOSH environment, follow the guide for your specific platform:

* [GCP](platforms/gcp/destroy-bosh.md)

### Other Tips

A Kubo environment is a set of configuration files used to deploy and update both BOSH and Kubo. 

**Note**: The [Kubo documentation](https://docs-kubo.cfapps.io/installing/) provides IaaS-specific procedures for deploying BOSH for Kubo.

Run `./bin/generate_env_config <ENV_PATH> <ENV_NAME> <PLATFORM_NAME>`
to generate a Kubo environment. The environment will be referred to as `KUBO_ENV`
in the Kubo documentation and in this guide, and will be located at `<ENV_PATH>/<ENV_NAME>`.

> Run `bin/generate_env_config --help` for more detailed information.

To view an adjust the available configuration options, edit the `<KUBO_ENV>/director.yml` file that
was created by the `generate_env_config` script.


