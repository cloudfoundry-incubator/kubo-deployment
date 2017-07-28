# Configuring IaaS routing for Kubo

On platforms that support native load-balancers kubo can be configured to leverage the
IaaS load balancers. Currently, only the GCP and AWS platforms are supported.

## Prerequisites

Before deploying and configuring Kubo, you need to carry out the following steps to 
setup the Load balancers:
   
1. This guide expects to be run in the same bash session as the [BOSH install](../../platforms/aws/install-bosh.md).
   If, for some reason, that is not the case, please set the `kubo_env_name` variable to the name
   of the Kubo environment before running the rest of the scripts.
   

1. On the BOSH bastion `cd` into the guide directory

   ```bash
   cd /share/kubo-deployment/docs/user-guide/routing/aws
   ```

1. Export these values. If you haven't tweaked any settings then use these defaults:

   ```bash
   export state_dir=~/kubo-env/${kubo_env_name}
   export kubo_terraform_state=${state_dir}/terraform.tfstate
   ``` 

1. Create the resources
   
   ```bash
   terraform apply \
      -var region=${region} \
      -var vpc_id=${vpc_id} \
      -var node_security_group_id=${default_security_groups} \
      -var public_subnet_id=${public_subnet_id} \
      -state=${kubo_terraform_state}
   ```

1. To get the outputs for the configuration files, run the following snippet:
   
   ```bash
   export master_target_pool=$(terraform output -state=${kubo_terraform_state} kubo_master_target_pool) # master_target_pool                                                                             
   export worker_target_pool=$(terraform output -state=${kubo_terraform_state} kubo_worker_target_pool) # worker_target_pool
   export kubernetes_master_host=$(terraform output -state=${kubo_terraform_state} master_lb_ip_address) # kubernetes_master_host
   ```

1. Update the Kubo environment using the following snippet:

   ```bash
   . ~/kubo-deployment/docs/user-guide/platforms/aws/setup_helpers
   set_iaas_routing "${state_dir}/director.yml"
   ```
   
   > It is also possible to set the configuration manually by editing the <KUBO_ENV>/director.yml  
