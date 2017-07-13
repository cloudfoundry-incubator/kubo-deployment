# Deploying BOSH for KUBO on GCP

1. From the GCP Cloud shell, SSH onto the bastion created during the [paving step](paving.md)

    ```bash
    gcloud compute ssh "${prefix}bosh-bastion" --zone ${zone}
    ```
    
1. Change directory to the root of the kubo-deployment repo

    ```bash
    cd /share/kubo-deployment
    ```
    
1. Create a kubo environment to set the configuration for BOSH and Kubo.

    ```bash
    export kubo_env=~/kubo-env
    export kubo_env_name=kubo
    export kubo_env_path="${kubo_env}/${kubo_env_name}"
    ./bin/generate_env_config "${kubo_env}" ${kubo_env_name} gcp
    ```

1.  Apply the default networking settings by running the following line:

    ```bash
    /usr/bin/update_gcp_env "${kubo_env_path}/director.yml" 
    ```

    The `kubo_env_path` will point to the folder containing the kubo settings,
    and will be referred to throughout this guide as `KUBO_ENV`.
    
    > Alternatively, it is possible to directly edit the file located at `${kubo_env_path}/director.yml`

1. Deploy a BOSH director for Kubo
    
    ```bash
    ./bin/deploy_bosh "${kubo_env_path}" ~/terraform.key.json
    ```
    Credentials and SSL certificates for the environment will be generated and
    saved into the configuration path in a file called `creds.yml`. This file
    contains sensitive information and should not be stored in VCS. The file
    `state.json` contains 
    [environment state for the BOSH environment](https://bosh.io/docs/cli-envs.html#deployment-state).

    Subsequent runs of `bin/bosh_deploy` will apply changes made to
    the configuration to an already existing KuBOSH installation, reusing
    the credentials stored in the `creds.yml`.