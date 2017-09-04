# Configuring GCP IaaS routing for Kubo

This document describes how to configure native GCP load balancers for Kubo.

## Prerequisites

Before deploying and configuring Kubo, you need to carry out the following steps to 
setup the Load balancer:
   
1. This guide expects to be run in the same bash session as the [BOSH install](../../platforms/gcp/install-bosh.md).
   If, for some reason, that is not the case, please set the `kubo_env_name` variable to the name
   of the Kubo environment before running the rest of the scripts.
   

1. On the BOSH bastion, change directory into the guide:

    ```bash
    cd /share/kubo-deployment/docs/user-guide/routing/gcp
    ```

1. Export environment variables that will be needed later. If you haven't tweaked any settings then use these defaults:

    ```bash
    export state_dir=~/kubo-env/${kubo_env_name}
    export kubo_terraform_state=${state_dir}/terraform.tfstate
    ```

1. Use Terraform to create the GCP resources for Kubo:

    ```bash
    terraform apply \
        -var network=${network} \
        -var projectid=${project_id} \
        -var region=${region} \
        -var prefix=${prefix} \
        -var ip_cidr_range="${subnet_ip_prefix}.0/24" \
        -state=${kubo_terraform_state}
    ```

1. To get the outputs for the configuration files, run the following snippet:
   
   ```bash
   export master_target_pool=$(terraform output -state=${kubo_terraform_state} kubo_master_target_pool) # master_target_pool                                                                             
   export kubernetes_master_host=$(terraform output -state=${kubo_terraform_state} master_lb_ip_address) # kubernetes_master_host
   ```

1. Update the Kubo environment:

    ```bash
    /usr/bin/set_iaas_routing "${state_dir}/director.yml"
    ```

    > **Note:** You can also set the configuration manually by editing <KUBO_ENV>/director.yml
