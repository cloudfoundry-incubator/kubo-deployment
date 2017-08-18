# Prerequisites for AWS Kubo deployment

1. An existing AWS account
1. A shell environment with the [Terraform CLI](https://www.terraform.io/docs/commands/index.html) installed
1. Access key ID and access key secret for an IAM user with administrator access.
1. A VPC in the zone that you will deploy kubo into, with a CIDR with a /22 netmask. **Please make sure, 
when using automatic paving, that the VPC is not created via the "Start VPC Wizard"**.
1. DNS hostnames enabled for the VPC.
    > **Note:** Select your VPC and click Actions > Edit DNS Hostnames.
