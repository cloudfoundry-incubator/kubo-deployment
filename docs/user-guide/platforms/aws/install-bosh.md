# Deploying BOSH for KUBO on AWS

1. SSH onto the bastion created during the [paving step](paving.md)

    ```bash
    cd $(dirname $kubo_terraform_state)
    ssh -i ~/deployer.pem ubuntu@$(terraform output bosh-bastion-ip)
    ```
    
1. Change directory to the root of the kubo-deployment repo

    ```bash
    cd /share/kubo-deployment
    ```
    
1. Create a kubo environment to set the configuration for BOSH and Kubo.

    ```bash
    export kubo_envs=~/kubo-env
    export kubo_env_name=kubo
    export kubo_env_path="${kubo_envs}/${kubo_env_name}"
 
    mkdir -p "${kubo_envs}"
    ./bin/generate_env_config "${kubo_envs}" "${kubo_env_name}" aws
    ```

1. Apply the default networking settings by running the following commands:

    ```bash
    . docs/user-guide/platforms/aws/setup_helpers
    update_aws_env "${kubo_env_path}/director.yml" 
    ```

    The `kubo_env_path` will point to the folder containing the kubo settings,
    and will be referred to throughout this guide as `KUBO_ENV`.
    
    > Alternatively, it is possible to directly edit the file located at `${kubo_env_path}/director.yml`

1. Go to IAM console (or use [aws-cli](https://aws.amazon.com/cli/)) and create a user with `Programmatic access` and following policy:
    > Please note that you will need to fill in your account id (__without hyphens__) in resource for iam:PassRole.
    ```json
    {
        "Version": "2012-10-17",
        "Statement": [
            {
                "Action": [
                    "ec2:AssociateAddress",
                    "ec2:AttachVolume",
                    "ec2:CreateVolume",
                    "ec2:DeleteSnapshot",
                    "ec2:DeleteVolume",
                    "ec2:DescribeAddresses",
                    "ec2:DescribeImages",
                    "ec2:DescribeInstances",
                    "ec2:DescribeRegions",
                    "ec2:DescribeSecurityGroups",
                    "ec2:DescribeSnapshots",
                    "ec2:DescribeSubnets",
                    "ec2:DescribeVolumes",
                    "ec2:DetachVolume",
                    "ec2:CreateSnapshot",
                    "ec2:CreateTags",
                    "ec2:RunInstances",
                    "ec2:TerminateInstances",
                    "ec2:RegisterImage",
                    "ec2:DeregisterImage",
                    "elasticloadbalancing:*"
                ],
                "Effect": "Allow",
                "Resource": "*"
            },
            {
                "Effect": "Allow",
                "Action": "iam:PassRole",
                "Resource": "arn:aws:iam::<account_id>:role/*kubo*"
            }
        ]
    }
    ```

1. Fill in `${kubo_env_path}/director-secrets.yml` with access key id and access key secret for the newly created user.

1. Deploy a BOSH director for Kubo
    
    ```bash
    ./bin/deploy_bosh "${kubo_env_path}" ~/deployer.pem
    ```
    Credentials and SSL certificates for the environment will be generated and
    saved in a file called `creds.yml` located in `KUBO_ENV`. This file
    contains sensitive information and should not be stored in VCS. The file
    `state.json` contains 
    [environment state for the BOSH environment](https://bosh.io/docs/cli-envs.html#deployment-state).

    Subsequent runs of `bin/bosh_deploy` will apply changes made to
    the configuration to an already existing BOSH installation, reusing
    the credentials stored in the `creds.yml`.
