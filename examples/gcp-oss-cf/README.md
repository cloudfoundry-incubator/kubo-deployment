# Example: Open Source Cloud Foundry and Kubo on GCP

## Prerequisites

1. Deploy a [bosh-bastion and BOSH director](https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/tree/master/docs/bosh#deploy-bosh-on-google-cloud-platform)

2. Deploy [Cloud Foundry](https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/tree/master/docs/cloudfoundry#deploying-cloud-foundry-on-google-compute-engine) 

## Prepare Infrastructure

The rest of the document assumes you are logged into the `bosh-bastion` you deployed above. The remaining steps should all be done in succession from a single session to retain required environment variables.

1. Start from the home directory of the bosh-bastion:
   ```
   cd
   ```

1. Install deployment dependencies:
   ```
   sudo curl https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-2.0.1-linux-amd64 -o /usr/bin/bosh-cli
   sudo chmod a+x /usr/bin/bosh-cli
   curl -L https://github.com/cloudfoundry-incubator/credhub-cli/releases/download/0.4.0/credhub-linux-0.4.0.tgz | tar zxv
   sudo mv credhub /usr/bin
   ```


1. Clone the [kubo-deployment](https://github.com/pivotal-cf-experimental/kubo-deployment) repo
   ```
   git clone https://github.com/pivotal-cf-experimental/kubo-deployment.git
   ```

1. `cd` into the examples directory

   ```
   cd ~/kubo-deployment/examples/gcp-oss-cf
   ```

1. Export these values. If you haven't tweaked any settings then use these defaults:

   ```
   export kubo_region=us-west1
   export kubo_zone=us-west1-a
   export cf_terraform_state=/share/docs/cloudfoundry/terraform.tfstate
   export network=$(terraform output -state=${cf_terraform_state} network)
   export kubo_env=kube
   export state_dir=~/kubo-env/${kubo_env}
   export kubo_terraform_state=${state_dir}/terraform.tfstate
   ``` 

1. Create a folder to store the environment configuration
   ```
   mkdir -p ${state_dir} 
   ```

1. View the Terraform execution plan to see the resources that will be created:
   ```
   terraform plan \
      -var network=${network} \
      -var projectid=${project_id} \
      -var kubo_region=${kubo_region} \
      -state=${kubo_terraform_state}
   ```

1. Create the resources
   ```
   terraform apply \
      -var network=${network} \
      -var projectid=${project_id} \
      -var kubo_region=${kubo_region} \
      -state=${kubo_terraform_state}
   ```

## Configure Kubo

1. `cd` to the `kubo-deployment` root
   ```
   cd ~/kubo-deployment
   ```

1. Generate the environment configuration
   ```
   bin/generate_env_config ~/kubo-env ${kubo_env} gcp
   ```

1. Retrieve the outputs of your Terraform run to be used in your Kubo deployment

   ```
   export kubo_subnet=$(terraform output -state=${kubo_terraform_state} kubo_subnet)
   export tcp_router_domain=tcp.$(terraform output -state=${cf_terraform_state} tcp_ip).xip.io
   export cf_system_domain=$(terraform output -state=${cf_terraform_state} ip).xip.io
   export common_secret=c1oudwc0w
   ```

1. Populate the director configurations
   ```
   erb examples/gcp-oss-cf/director.yml.erb > ${state_dir}/director.yml
   erb examples/gcp-oss-cf/director-secrets.yml.erb > ${state_dir}/director-secrets.yml
   ```

1. Generate a service account key for the bosh-user
   ```
   export service_account=bosh-user
   export service_account_creds=${state_dir}/service_account.json
   export service_account_email=${service_account}@${project_id}.iam.gserviceaccount.com
   gcloud iam service-accounts keys create ${service_account_creds} --iam-account ${service_account_email}
   ```

## Deploy Kubo

1. Deploy a BOSH director for Kubo
   ```
   bin/deploy_bosh ${state_dir} ${service_account_creds} 
   ```

1. Find the password for your BOSH director
   ```
   bosh-cli int ${state_dir}/creds.yml --path=/admin_password
   ```

1. Login with that password
   ```
   bosh-cli login -e ${kubo_env}

   Email: admin
   Password: [output of previous command]
   ```

1. Deploy Kubo
   ```
   bin/deploy_k8s ${state_dir} kube public
   ```

