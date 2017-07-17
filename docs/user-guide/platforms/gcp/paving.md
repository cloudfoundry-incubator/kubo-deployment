# Paving the infrastructure for Kubo on GCP

## Setup the shell environment

1. Make sure that you have the owner role in your existing Google Cloud Platform project

1. In your existing Google Cloud Platform project, enable the [IAM API](https://console.cloud.google.com/apis/api/iam.googleapis.com/overview).

1. In your existing Google Cloud Platform project, open Cloud Shell (the small `>_` prompt icon in the web console menu bar).

1. When deploying kubo more than once, it is required to set a unique prefix
  for every installation. Please use letters and dashes only.
  
  ```bash
  export prefix=my-kubo # This prefix should be unique for every install
  ```

1.  Configure the following environment variables:

  ```bash
  export project_id=$(gcloud config get-value project)
  export subnet_ip_prefix="10.0.1" # Create new subnet for deployment in $subnet_ip_prefix.0/24
  export region=us-east1 # region that you will deploy Kubo in
  export zone=us-east1-d # zone that you will deploy Kubo in
  export state_dir=~/kubo-env # Location for storing the terraform state
  export kubo_terraform_state=${state_dir}/terraform.tfstate
  export service_account_email=${prefix}terraform@${project_id}.iam.gserviceaccount.com
  export network=<An existing GCP network for deploying kubo>
  ```
  
  > When using the [CloudFoundry routing mode](../../routing/cf.md) the GCP network above 
  > needs to be the same network that CloudFoundry is using 

1. Configure `gcloud` to use the region and zone specified above:

  ```bash
  gcloud config set compute/zone ${zone}
  gcloud config set compute/region ${region}
  ```
  
## Setup GCP account for terraform

1. Create a service account and key:
  
  ```bash
  gcloud iam service-accounts create ${prefix}terraform
  gcloud iam service-accounts keys create ~/terraform.key.json \
      --iam-account ${service_account_email}
  ```

1. Grant the new service account owner access to your project:
  
  ```bash
  gcloud projects add-iam-policy-binding ${project_id} \
    --member serviceAccount:${service_account_email} \
    --role roles/owner
  ```

1. Make your service account's key available in an environment 
  variable to be used by `terraform`:

  ```bash
  export GOOGLE_CREDENTIALS=$(cat ~/terraform.key.json)
  ```

## Deploy supporting infrastructure

This step sets up a subnetwork with a bastion VM and a set of firewall 
rules to secure access to the kubo deployment.

### Steps

1. Clone this repository and go into the installation docs directory:

  ```bash
  git clone https://github.com/cloudfoundry-incubator/kubo-deployment.git
  cd kubo-deployment/docs/user-guide/platforms/gcp
  ```

1. Create the resources (should take between 60-90 seconds):

  ```bash
  docker run -i -t \
    -e CHECKPOINT_DISABLE=1 \
    -e "GOOGLE_CREDENTIALS=${GOOGLE_CREDENTIALS}" \
    -v $(pwd):/$(basename $(pwd)) \
    -w /$(basename $(pwd)) \
    hashicorp/terraform:light apply \
      -var service_account_email=${service_account_email} \
      -var projectid=${project_id} \
      -var network=${network} \
      -var region=${region} \
      -var prefix=${prefix} \
      -var zone=${zone} \
      -var subnet_ip_prefix=${subnet_ip_prefix}
  ```

> _Note_: It's possible to preview the terraform execution plan by running the 
> same command, using `plan` in place of `apply`

1. Copy the service account key to the newly created bastion box
    
  ```bash
  gcloud compute scp ~/terraform.key.json "${prefix}bosh-bastion":./ --zone ${zone}
  ```