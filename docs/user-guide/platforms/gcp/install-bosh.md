# Deploying BOSH for KUBO on GCP

1. Return to the root of the kubo-deployment repo

    ```bash
    cd ~/kubo-deployment
    ```
1. [Create a kubo environment](../../create-kubo-env.md). Please make sure
to fill in all the networking and IaaS-specific options.

1. Deploy a BOSH director for Kubo
    ```bash
    bin/deploy_bosh <KUBO_ENV> ${service_account_key}
    ```
    Credentials and SSL certificates for the environment will be generated and
    saved into the configuration path in a file called `creds.yml`. This file
    contains sensitive information and should not be stored in VCS. The file
    `state.json` contains [environment state for the BOSH environment](https://bosh.io/docs/cli-envs.html#deployment-state).

    Subsequent runs of `bin/bosh_deploy` will apply changes made to
    the configuration to an already existing KuBOSH installation, reusing
    the credentials stored in the `creds.yml`.
