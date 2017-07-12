# Configuring IaaS routing for Kubo

On platforms that support native load-balancers kubo can be configured to leverage the
IaaS load balancers. Currently, only the GCP platform is supported.

## Prerequisites

Before deploying and configuring Kubo, you need to carry out the following steps to 
setup the Load balancers:   

1. `cd` into the guide directory

   ```bash
   cd ~/kubo-deployment/docs/templates/gcp-lb
   ```

1. Export these values. If you haven't tweaked any settings then use these defaults:

   ```bash
   export project_id=$(gcloud config get-value project)
   export kubo_region=us-west1
   export kubo_zone=us-west1-a
   export network=bosh # GCP network that BOSH resides in
   
   export kubo_env=kube
   export state_dir=~/kubo-env/${kubo_env}
   export kubo_terraform_state=${state_dir}/terraform.tfstate
   ``` 

1. Create the resources
   ```bash
   terraform apply \
      -var network=${network} \
      -var projectid=${project_id} \
      -var kubo_region=${kubo_region} \
      -state=${kubo_terraform_state}
   ```

Additionally, the terraform script accepts the following variables:
  
  - `ip_cidr_range`: the CIDR range for the kubo subnetwork. The default value is `10.0.1.0/24`
  - `prefix`: A prefix to use for all the GCP resource names. Defaults to an empty string.

To get the outputs for the configuration files, run the following snippet:
   
   ```bash
   terraform output -state=${kubo_terraform_state} kubo_subnet # subnetwork 
   terraform output -state=${kubo_terraform_state} kubo_master_target_pool # master_target_pool                                                                             
   terraform output -state=${kubo_terraform_state} kubo_worker_target_pool # worker_target_pool
   terraform output -state=${kubo_terraform_state} master_lb_ip_address # kubernetes_master_host
   ```

In order to configure kubo to use IaaS routing, modify the `<KUBO_ENV>/director.yml`:

  - Comment out all the lines grouped underneath the **CF routing mode settings** comment
  
  - Uncomment all the lines grouped underneath the **IaaS routing mode settings** comment 
    and fill in all the values, as outlined in the snippet above.
  
