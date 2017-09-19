# Deploying BOSH for Kubo on GCP

1. From the GCP Cloud shell, SSH onto the bastion created during the [paving step](paving.md):

    ```bash
    gcloud compute ssh "${prefix}bosh-bastion" --zone ${zone}
    ```
    
1. Change directory to the root of the `kubo-deployment` repo, which was copied into the bastion during the paving step:

    ```bash
    cd /share/kubo-deployment
    ```
    
1. Generate a Kubo configuration template:

    ```bash
    export kubo_envs=~/kubo-env
    export kubo_env_name=kubo
    export kubo_env_path="${kubo_envs}/${kubo_env_name}"

    mkdir -p "${kubo_envs}"
    ./bin/generate_env_config "${kubo_envs}" ${kubo_env_name} gcp
    ```

    `kubo_env_path` points to the directory containing the Kubo configuration, and will be referred to throughout this guide as `KUBO_ENV`.

1. The `update_gcp_env` command knows about the default network settings configured during the paving step. Execute it to apply those settings onto the template:

    ```bash
    /usr/bin/update_gcp_env "${kubo_env_path}/director.yml"
    ```

    > **Note:** you can directly edit the configuration file located at `${kubo_env_path}/director.yml`

1. Deploy the BOSH director for Kubo:
    
    ```bash
    ./bin/deploy_bosh "${kubo_env_path}" ~/terraform.key.json
    ```

    `deploy_bosh` uses the Kubo configuration we generated and our GCP key from the paving step to deploy BOSH. After it's done, you will have a BOSH director with all the necessary components to deploy Kubo. During its execution, `deploy_bosh` generates credentials and SSL certificates for that BOSH director and saves it in `creds.yml` located in `KUBO_ENV`. This file contains sensitive information and should not be stored in VCS. Also, the [state for the BOSH environment](https://bosh.io/docs/cli-envs.html#deployment-state) is stored in `state.json`. Subsequent runs of `deploy_bosh` will use `creds.yml` and `state.json` to intelligently apply changes to the BOSH environment.
