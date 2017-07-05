# Deploying BOSH for KUBO on GCP

1. Return to the root of the kubo-deployment repo

    ```bash
    cd ~/kubo-deployment
    ```
1. Create environment config:

    ```bash
    ./bin/generate_env_config <ENV_PATH> <ENV_NAME> gcp
    ```
    > Run `bin/generate_env_config --help` for more detailed information.

    This will generate couple of `yml` files in `<ENV_PATH>/<ENV_NAME>` (it is
    called `<KUBO_ENV>` in this guide). Follow the comments in
    `<KUBO_ENV>/director.yml` to fill in the values. You might need to fill in
    the values in `<KUBO_ENV>/director-secrets.yml` as well.

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
