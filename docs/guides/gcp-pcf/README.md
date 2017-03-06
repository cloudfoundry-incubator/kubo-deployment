# Example: Existing Pivotal Cloud Foundry and Kubo on GCP

## Prerequisites

1. These instructions assume that you have an installation of PCF running on GCP with the [Routing Release](https://github.com/cloudfoundry-incubator/routing-release) enabled

## Prepare Infrastructure

1. In your existing Google Cloud Platform project, enable the following APIs:
- [GCE API](https://console.developers.google.com/apis/api/compute_component/overview)
- [IAM API](https://console.cloud.google.com/apis/api/iam.googleapis.com/overview)
- [Cloud Resource Manager API](https://console.cloud.google.com/apis/api/cloudresourcemanager.googleapis.com/overview)

1. In your existing Google Cloud Platform project, open Cloud Shell (the small `>_` prompt icon in the web console menu bar).

1.  Configure a few environment variables:

  ```
  export project_id=$(gcloud config list 2>/dev/null | grep project | sed -e 's/project = //g')
  export network=<Network that your Cloud Foundry installation is running in>
  export subnet_ip_prefix="10.0.1" # Create new subnet for deployment in $subnet_ip_prefix.0/24
  export region=us-east1
  export zone=us-east1-d
  export terraform_state_dir=~/kubo-env
  export service_account_email=terraform@${project_id}.iam.gserviceaccount.com
  ```

1. Configure `gcloud` to use your preferred region and zone:

  ```
  gcloud config set compute/zone ${zone}
  gcloud config set compute/region ${region}
  ```

1. Create a service account and key:

  ```
  gcloud iam service-accounts create terraform
  gcloud iam service-accounts keys create ~/terraform.key.json \
      --iam-account ${service_account_email}
  ```

1. Grant the new service account editor access to your project:

  ```
  gcloud projects add-iam-policy-binding ${project_id} \
    --member serviceAccount:${service_account_email} \
    --role roles/owner
  ```

1. Make your service account's key available in an environment variable to be used by `terraform`:

  ```
  export GOOGLE_CREDENTIALS=$(cat ~/terraform.key.json)
  ```

## Deploy supporting infrastructure

### Steps
> **Note:** All of these steps should be performed inside the Cloud Shell in your browser.

1. Clone this repository and go into the installation docs directory:

  ```
  git clone https://github.com/pivotal-cf-experimental/kubo-deployment.git
  cd kubo-deployment/docs/guides/gcp-pcf
  ```

1. Create a folder to store the terraform output
   ```
   mkdir -p ${terraform_state_dir}
   ```

1. View the Terraform execution plan to see the resources that will be created:
  ```
  docker run -i -t \
    -e "GOOGLE_CREDENTIALS=${GOOGLE_CREDENTIALS}" \
    -v `pwd`:/$(basename `pwd`) \
    -w /$(basename `pwd`) \
    hashicorp/terraform:light plan \
      -var service_account_email=${service_account_email} \
      -var network=${network} \
      -var projectid=${project_id} \
      -var region=${region} \
      -var kubo_region=${kubo_region} \
      -var zone=${zone} \
      -var subnet_ip_prefix=${subnet_ip_prefix} \
      -state=${kubo_terraform_state}
  ```

1. Create the resources (should take between 60-90 seconds):

  ```
  docker run -i -t \
    -e "GOOGLE_CREDENTIALS=${GOOGLE_CREDENTIALS}" \
    -v `pwd`:/$(basename `pwd`) \
    -w /$(basename `pwd`) \
    hashicorp/terraform:light apply \
      -var service_account_email=${service_account_email} \
      -var projectid=${project_id} \
      -var network=${network} \
      -var region=${region} \
      -var kubo_region=${kubo_region} \
      -var zone=${zone} \
      -var subnet_ip_prefix=${subnet_ip_prefix} \
      -state=${kubo_terraform_state}
  ```

Now you have the infrastructure ready to deploy a BOSH director.

## Configure Kubo

1. SSH to the `bosh-bastion` VM. You can SSH form Cloud Shell or any workstation that has `gcloud` installed:

  ```
  gcloud compute ssh bosh-bastion
  ```

1.  Configure a few environment variables:

  ```
  export project_id=$(gcloud config list 2>/dev/null | grep project | sed -e 's/project = //g')
  export network=<Network that your Cloud Foundry installation is running in>
  export kubo_region=us-east1
  export kubo_zone=us-east1-d
  export kubo_env=kube
  export state_dir=~/kubo-env/${kubo_env}
  export common_secret=c1oudc0w
  export subnet_ip_prefix="10.0.1" # Create new subnet for deployment in $subnet_ip_prefix.0/24

  export tcp_router_domain=[domain of TCP router]
  export cf_system_domain=[Cloud Foundry system domain]
  export routing_cf_client_id=[client id for that can access TCP router]
  export kubo_subnet=[name of subnet created by Terraform]
  ```

1. Create a folder to store the environment configuration

  ```
  mkdir -p ${state_dir}
  ```

1. `cd` to the `kubo-deployment` root
   ```
   cd ~/kubo-deployment
   ```

1. Generate the environment configuration
   ```
   bin/generate_env_config ~/kubo-env ${kubo_env} gcp
   ```

1. Populate the director configurations
   ```
   erb docs/guides/gcp-pcf/director.yml.erb > ${state_dir}/director.yml
   ```

1. Fill in the values in ${state_dir}/director-secrets.yml

1. Generate a service account key for the bosh-user
   ```
   export service_account=bosh-user
   export service_account_creds=${state_dir}/service_account.json
   export service_account_email=${service_account}@${project_id}.iam.gserviceaccount.com
   gcloud iam service-accounts create ${service_account}
   gcloud iam service-accounts keys create ${service_account_creds} --iam-account ${service_account_email}
   gcloud projects add-iam-policy-binding ${project_id} \
     --member serviceAccount:${service_account_email=} --role roles/compute.instanceAdmin
   gcloud projects add-iam-policy-binding ${project_id} \
     --member serviceAccount:${service_account_email=} --role roles/compute.storageAdmin
   gcloud projects add-iam-policy-binding ${project_id} \
     --member serviceAccount:${service_account_email=} --role roles/storage.admin
   gcloud projects add-iam-policy-binding ${project_id} \
     --member serviceAccount:${service_account_email=} --role  roles/compute.networkAdmin
   gcloud projects add-iam-policy-binding ${project_id} \
     --member serviceAccount:${service_account_email=} --role roles/iam.serviceAccountActor
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

1. Setup kubectl and access your new Kubernetes cluster
   ```
   bin/set_kubeconfig ${state_dir} kube
   kubectl get pods --namespace=kube-system
   ```
