# Paving the infrastructure for Kubo on GCP

## Deploy a BOSH bastion

Configure a GCP project and deploy a BOSH bastion by following the "Configure your Google Cloud Platform environment" and "Deploy supporting infrastructure" steps in
  [these instructions](https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/blob/c2cdba4f2ac8944ce7eb9749f053d45588932e3b/docs/bosh/README.md).
  
## Set up infrastructure for Kubo

1. ssh into the bosh-bastion created in the prerequisites
    ```bash
    gcloud compute ssh bosh-bastion
    ```

1. Clone the [kubo-deployment](https://github.com/cloudfoundry-incubator/kubo-deployment) repo
    ```bash
    git clone https://github.com/cloudfoundry-incubator/kubo-deployment.git ~/kubo-deployment
    cd ~/kubo-deployment/docs/templates/gcp
    ```

1. Export the project id:

    ```bash
    export project_id=$(gcloud config get-value project)
    ```
      
1. _(Optional)_ You can override the defaults by specifying any of the following 
  environment variables:
   
    ```bash
    export network=custom-kubo # The default value is `bosh`
    export ip_cidr_range=10.0.0.1/24 # The default value is 10.0.1.0/24
    export kubo_region=us-east1 # The default value is us-west1
    export kubo_prefix=custom # The default value is an empty string
    ```

1. Create a folder to store the environment configuration
    ```bash
    export state_dir=~/kubo-env
    export kubo_terraform_state=${state_dir}/gcp.tfstate
    mkdir -p ${state_dir}
    ```

1. Create the resources
    ```bash
    terraform apply \
      -var network=${network} \
      -var projectid=${project_id} \
      -var kubo_region=${kubo_region} \
      -var ip_cidr_range=${ip_cidr_range} \
      -var prefix=${kubo_prefix} \
      -state=${kubo_terraform_state}
    ```
    
    > You can also preview your deployment before applying using `terraform plan`.
   
1. Get the name of the kubo subnet by running:

    ```bash
    terraform output -state=${kubo_terraform_state} kubo_subnet
    ```