# Paving the infrastructure for Kubo on AWS

## Setup the shell environment

1. Create an EC2 [key pair](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html)
named `deployer` and save the key file on your machine at ~/deployer.pem.

    Make sure that the key has the proper permissions by running the following command:
    ```bash
    chmod 600 ~/deployer.pem
    ```
1. When deploying kubo more than once, it is required to set a unique prefix
  for every installation. Please use letters and dashes only.

    ```bash
    export prefix=my-kubo # This prefix should be unique for every install
    ```
1. Configure the following environment variables:

    ```bash
    export AWS_ACCESS_KEY_ID=<Your AWS access key ID>
    export AWS_SECRET_ACCESS_KEY=<Your AWS secret access key>
    export vpc_id=<An existing VPC for deploying kubo>
    export key_name=deployer # name of private key to use on Kubo VMs
    export private_key_filename="~/${key_name}.pem" # private key to use on Kubo VMs
    export region=us-west-2 # region that you will deploy Kubo in
    export zone=us-west-2a # zone that you will deploy Kubo in
    export public_subnet_ip_prefix="10.0.1" # subnet that will be used for bastion VM, NAT Gateway and load balancers
    export private_subnet_ip_prefix="10.0.2" # subnet that will be used for Kubo VMs and BOSH director
    export kubo_terraform_state=~/terraform.tfstate
    ```

    > When using the [CloudFoundry routing mode](../../routing/cf.md) the VPC above
    > needs to be the same network that CloudFoundry is using

## Deploy supporting infrastructure

This step sets up a subnetwork with a bastion VM and security group
rules to secure access to the kubo deployment.

### Steps

1. Get latest version of kubo-deployment:

    ```bash
    cd ~
    wget https://storage.googleapis.com/kubo-public/kubo-deployment-latest.tgz
    tar -xvf kubo-deployment-latest.tgz
    cd kubo-deployment/docs/user-guide/platforms/aws
    ```

1. Initialise terraform working directory

    ```bash
    terraform init
    ```

1. Create the resources (should take between 60-90 seconds):

    > **Note:** It's possible to preview the terraform execution plan by running the same command, using `plan` in place of `apply`

    ```bash
    terraform apply \
      -var region="${region}" \
      -var zone="${zone}" \
      -var vpc_id="${vpc_id}" \
      -var prefix="${prefix}" \
      -var public_subnet_ip_prefix="${public_subnet_ip_prefix}" \
      -var private_subnet_ip_prefix="${private_subnet_ip_prefix}" \
      -var private_key_filename="${private_key_filename}" \
      -var key_name="${key_name}" \
      -state=${kubo_terraform_state}
    ```
    > **Note:** The previously created bastion box will be deleted by subsequent runs of the `terraform apply`
