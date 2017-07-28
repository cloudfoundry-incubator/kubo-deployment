# Deploying BOSH for KUBO on AWS

1. Find the Public DNS name of the bastion instance on AWS

1. SSH onto the bastion created during the [paving step](paving.md)

    ```bash
    ssh -i ~/deployer.pem ubuntu@<public DNS name>
    ```
    
1. Change directory to the root of the kubo-deployment repo

    ```bash
    cd /share/kubo-deployment
    ```
    
1. Create a kubo environment to set the configuration for BOSH and Kubo.

    ```bash
    export kubo_envs=~/kubo-env
    export kubo_env_name=kubo
    export kubo_env_path="${kubo_envs}/${kubo_env_name}"
 
    mkdir -p "${kubo_envs}"
    ./bin/generate_env_config "${kubo_envs}" ${kubo_env_name} aws
    ```

1. Apply the default networking settings by running the following line:

    ```bash
    . docs/user-guide/platforms/aws/setup_helpers
    update_aws_env "${kubo_env_path}/director.yml" 
    ```

    The `kubo_env_path` will point to the folder containing the kubo settings,
    and will be referred to throughout this guide as `KUBO_ENV`.
    
    > Alternatively, it is possible to directly edit the file located at `${kubo_env_path}/director.yml`

1. Fill in `${kubo_env_path}/director-secrets.yml` with your access key id and access key secret

1. Deploy a BOSH director for Kubo
    
    ```bash
    ./bin/deploy_bosh "${kubo_env_path}" ~/deployer.pem
    ```
    Credentials and SSL certificates for the environment will be generated and
    saved in a file called `creds.yml` located in `KUBO_ENV`. This file
    contains sensitive information and should not be stored in VCS. The file
    `state.json` contains 
    [environment state for the BOSH environment](https://bosh.io/docs/cli-envs.html#deployment-state).

    Subsequent runs of `bin/bosh_deploy` will apply changes made to
    the configuration to an already existing BOSH installation, reusing
    the credentials stored in the `creds.yml`.
