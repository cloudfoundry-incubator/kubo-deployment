# Paving the Infrastructure for Kubo on GCP

## Setup the Shell Environment

1. In your existing Google Cloud Platform project, open Cloud Shell (the small `>_` prompt icon in the web console menu bar).

1. When deploying kubo more than once, it is required to set a unique prefix
  for every installation. Please use letters and dashes only.

    ```bash
    export prefix=my-kubo # This prefix should be unique for every install
    ```

1. Create a GCP network using the [GCP dashboard](https://console.cloud.google.com). Configure a CIDR range with a mask exactly 24-bits large for a GCP subnet (it will be automatically created later). In the example below, we're using a subnet with CIDR range `10.0.1.0/24`.

    Export environment variables for these resources:

    ```bash
    export network=<your GCP network name>
    export subnet_ip_prefix="10.0.1" # Your subnet prefix
    ```

    > **Note:** When using the [Cloud Foundry routing mode](../../routing/cf.md) the GCP network above needs to be the same network that CloudFoundry is using.

1. Configure other environment variables that will be used in this guide:

    ```bash
    export project_id=$(gcloud config get-value project)
    export region=us-east1 # region that you will deploy Kubo in
    export zone=us-east1-d # zone that you will deploy Kubo in
    export service_account_email=${prefix}terraform@${project_id}.iam.gserviceaccount.com
    ```

1. Configure `gcloud` to use the region and zone specified above:

    ```bash
    gcloud config set compute/zone ${zone}
    gcloud config set compute/region ${region}
    ```

## Setup a GCP Account for Terraform

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

## Deploy Supporting Infrastructure

This step sets up a subnetwork with a bastion VM and a set of firewall
rules to secure access to the kubo deployment.

### Steps

1. Get latest version of kubo-deployment:

    ```bash
    cd ~
    wget https://storage.googleapis.com/kubo-public/kubo-deployment-latest.tgz
    tar -xvf kubo-deployment-latest.tgz
    cd ~/kubo-deployment/docs/user-guide/platforms/gcp
    ```
1. Initialize the terraform cloud provider:

    ```bash
    docker run -i -t \
      -v $(pwd):/$(basename $(pwd)) \
      -w /$(basename $(pwd)) \
      hashicorp/terraform:light init
    ```

1. Create the bastion and other resources (should take between 60-90 seconds):

    > **Note:** It's possible to preview the terraform execution plan before applying it by running `plan` instead of `apply`.

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

1. Copy the service account key to the newly created bastion box:

    ```bash
    gcloud compute scp ~/terraform.key.json "${prefix}bosh-bastion":./ --zone ${zone}
    ```

    This will be used later when we SSH into the bastion to deploy BOSH for Kubo.
