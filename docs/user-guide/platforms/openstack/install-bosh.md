# Deploying BOSH for KUBO on OpenStack

**Prerequisite:** The machine executing the commands below must be able to access VMs on the OpenStack network. Depending on your network topology, a bastion host (jumpbox) may be needed.

1. Get latest version of kubo-deployment:

    ```bash
    cd ~
    wget https://storage.googleapis.com/kubo-public/kubo-deployment-latest.tgz
    tar -xvf kubo-deployment-latest.tgz
    cd kubo-deployment
    ```

1. Create a kubo environment to set the configuration for BOSH and Kubo.

    ```bash
    export kubo_env=~/kubo-env
    export kubo_env_name=kubo
    export kubo_env_path="${kubo_env}/${kubo_env_name}"

    mkdir -p "${kubo_env}"
    ./bin/generate_env_config "${kubo_env}" ${kubo_env_name} openstack
    ```

1.  Populate the environment config skeleton created at
    `${kubo_env_path}/director.yml`.

    > Inside `director.yml`, `default_key_name` should be set to the name of
      an OpenStack key pair. The private key from that pair will be required
      when we deploy BOSH. We'll refer to this private key later as `private_key.pem`.

    The `kubo_env_path` will point to the folder containing the kubo settings,
    and will be referred to throughout this guide as `KUBO_ENV`.

1. Deploy a BOSH director for Kubo

    ```bash
    ./bin/deploy_bosh "${kubo_env_path}" <private_key.pem>
    ```
    Credentials and SSL certificates for the environment will be generated and
    saved into the configuration path in a file called `creds.yml`. This file
    contains sensitive information and should not be stored in VCS. The file
    `state.json` contains
    [environment state for the BOSH environment](https://bosh.io/docs/cli-envs.html#deployment-state).

    Subsequent runs of `bin/bosh_deploy` will apply changes made to
    the configuration to an already existing BOSH installation, reusing
    the credentials stored in the `creds.yml`.
