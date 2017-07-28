# Paving the infrastructure for Kubo on AWS

## Setup the shell environment

1. Create an EC2 [key pair](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html) named `deployer`

1. Save the key file on your machine at ~/deployer.pem

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
  export key_name=deployer
  export private_key="$(cat ~/${key_name}.pem)"
  export region=us-west-2 # region that you will deploy Kubo in
  export zone=us-west-2a # zone that you will deploy Kubo in
  export public_subnet_ip_prefix="10.0.1"
  export private_subnet_ip_prefix="10.0.2"
  export kubo_terraform_state=~/terraform.tfstate
  ```
  
  > When using the [CloudFoundry routing mode](../../routing/cf.md) the VPC above 
  > needs to be the same network that CloudFoundry is using 

## Deploy supporting infrastructure

This step sets up a subnetwork with a bastion VM and security group
rules to secure access to the kubo deployment.

### Steps

1. Clone this repository and go into the installation docs directory:

  ```bash
  git clone https://github.com/cloudfoundry-incubator/kubo-deployment.git
  cd kubo-deployment/docs/user-guide/platforms/aws
  ```

1. Create the resources (should take between 60-90 seconds):

  ```bash
  terraform apply \
    -var region="${region}" \
    -var zone="${zone}" \
    -var vpc_id="${vpc_id}"
    -var prefix="${prefix}" \
    -var public_subnet_ip_prefix="${public_subnet_ip_prefix}" \
    -var private_subnet_ip_prefix="${private_subnet_ip_prefix}" \
    -var private_key="${private_key}" \
    -var key_name="${key_name}"
    -state=${kubo_terraform_state}
  ```

> _Note_: It's possible to preview the terraform execution plan by running the 
> same command, using `plan` in place of `apply`
