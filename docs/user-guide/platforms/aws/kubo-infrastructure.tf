variable "region" {
    type = "string"
}

variable "zone" {
    type = "string"
}

variable "vpc_id" {
    type = "string"
}

variable "public_subnet_ip_prefix" {
    type = "string"
    default = "10.0.1"
}

variable "private_subnet_ip_prefix" {
    type = "string"
    default = "10.0.2"
}

variable "key_name" {
    type = "string"
}

variable "private_key_filename" {
    type = "string"
}

variable "prefix" {
    type = "string"
    default = ""
}

provider "aws" {
    region = "${var.region}"
}

resource "random_id" "kubernetes-cluster-tag" {
  byte_length = 16
}

resource "aws_internet_gateway" "gateway" {
    vpc_id = "${var.vpc_id}"
}

resource "aws_subnet" "public" {
    vpc_id     = "${var.vpc_id}"
    cidr_block = "${var.public_subnet_ip_prefix}.0/24"
    availability_zone = "${var.zone}"

    tags {
      Name = "${var.prefix}kubo-public"
      KubernetesCluster = "${random_id.kubernetes-cluster-tag.b64}"
    }
}

resource "aws_route_table" "public" {
    vpc_id = "${var.vpc_id}"

    tags {
      Name = "${var.prefix}public-route-table"
    }

    route {
      cidr_block = "0.0.0.0/0"
      gateway_id = "${aws_internet_gateway.gateway.id}"
    }
}

resource "aws_route_table_association" "public" {
    subnet_id      = "${aws_subnet.public.id}"
    route_table_id = "${aws_route_table.public.id}"
}

resource "aws_eip" "nat" {
}

resource "aws_nat_gateway" "nat" {
    allocation_id = "${aws_eip.nat.id}"
    subnet_id     = "${aws_subnet.public.id}"
}


resource "aws_subnet" "private" {
    vpc_id     = "${var.vpc_id}"
    cidr_block = "${var.private_subnet_ip_prefix}.0/24"
    availability_zone = "${var.zone}"

    tags {
      Name = "${var.prefix}kubo-private"
      KubernetesCluster = "${random_id.kubernetes-cluster-tag.b64}"
    }
}

resource "aws_route_table" "private" {
    vpc_id = "${var.vpc_id}"

    tags {
      Name = "${var.prefix}private-route-table"
    }

    route {
      cidr_block = "0.0.0.0/0"
      gateway_id = "${aws_nat_gateway.nat.id}"
    }
}

resource "aws_route_table_association" "private" {
    subnet_id      = "${aws_subnet.private.id}"
    route_table_id = "${aws_route_table.private.id}"
}

resource "aws_security_group" "nodes" {
    name        = "${var.prefix}node-access"
    vpc_id      = "${var.vpc_id}"
}

resource "aws_security_group_rule" "outbound" {
    type            = "egress"
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    cidr_blocks     = ["0.0.0.0/0"]

    security_group_id = "${aws_security_group.nodes.id}"
}

resource "aws_security_group_rule" "UAA" {
    type        = "ingress"
    from_port   = 8443
    to_port     = 8443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]

    security_group_id = "${aws_security_group.nodes.id}"
}

resource "aws_security_group_rule" "ssh" {
    type            = "ingress"
    from_port       = 22
    to_port         = 22
    protocol        = "tcp"
    cidr_blocks     = ["0.0.0.0/0"]

    security_group_id = "${aws_security_group.nodes.id}"
}

resource "aws_security_group_rule" "node-to-node" {
    type            = "ingress"
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    source_security_group_id = "${aws_security_group.nodes.id}"

    security_group_id = "${aws_security_group.nodes.id}"
}

data "aws_ami" "ubuntu" {
    most_recent = true

    filter {
      name   = "name"
      values = ["ubuntu/images/hvm-ssd/ubuntu-trusty-14.04-amd64-server-*"]
    }

    filter {
      name   = "virtualization-type"
      values = ["hvm"]
    }

    owners = ["099720109477"] # Canonical
}

resource "aws_iam_role_policy" "kubo-master" {
    name = "${var.prefix}kubo-master"
    role = "${aws_iam_role.kubo-master.id}"

    policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "ec2:*",
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": "elasticloadbalancing:*",
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_instance_profile" "kubo-master" {
    name = "${var.prefix}kubo-master"
    role = "${aws_iam_role.kubo-master.name}"
}

resource "aws_iam_role" "kubo-master" {
    name = "${var.prefix}kubo-master"
    assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "kubo-worker" {
    name = "${var.prefix}kubo-worker"
    role = "${aws_iam_role.kubo-worker.id}"

    policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "ec2:Describe*",
      "Effect": "Allow",
      "Resource": "*"
    },
    {
      "Action": "ec2:AttachVolume",
      "Effect": "Allow",
      "Resource": "*"
    },
    {
      "Action": "ec2:DetachVolume",
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_instance_profile" "kubo-worker" {
    name = "${var.prefix}kubo-worker"
    role = "${aws_iam_role.kubo-worker.name}"
}

resource "aws_iam_role" "kubo-worker" {
    name = "${var.prefix}kubo-worker"
    assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_instance" "bastion" {
    ami           = "${data.aws_ami.ubuntu.id}"
    instance_type = "t2.micro"
    subnet_id     = "${aws_subnet.public.id}"
    availability_zone = "${var.zone}"
    key_name      = "${var.key_name}"
    vpc_security_group_ids = ["${aws_security_group.nodes.id}"]
    associate_public_ip_address = true
    tags {
      Name = "${var.prefix}bosh-bastion"
    }
    provisioner "remote-exec" {
        inline = [
            "set -eu",
            "sudo apt-get update",
            "sudo apt-get install -y build-essential zlibc zlib1g-dev ruby ruby-dev openssl libxslt-dev libxml2-dev libssl-dev libreadline6 libreadline6-dev libyaml-dev libsqlite3-dev sqlite3",
            "sudo apt-get install -y git",
            "sudo apt-get install -y unzip",
            "curl -L https://github.com/cloudfoundry-incubator/credhub-cli/releases/download/1.3.0/credhub-linux-1.3.0.tgz | tar zxv && sudo chmod a+x credhub && sudo mv credhub /usr/bin",
            "sudo curl -L https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl -o /usr/bin/kubectl && sudo chmod a+x /usr/bin/kubectl",
            "sudo curl https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-2.0.27-linux-amd64 -o /usr/bin/bosh-cli && sudo chmod a+x /usr/bin/bosh-cli",
            "sudo wget https://releases.hashicorp.com/terraform/0.10.2/terraform_0.10.2_linux_amd64.zip",
            "sudo unzip terraform*.zip -d /usr/local/bin",
            "sudo sh -c 'sudo cat > /etc/profile.d/bosh.sh <<'EOF'",
            "#!/bin/bash",
            "export private_subnet_id=${aws_subnet.private.id}",
            "export public_subnet_id=${aws_subnet.public.id}",
            "export vpc_id=${var.vpc_id}",
            "export default_security_groups=${aws_security_group.nodes.id}",
            "export private_subnet_ip_prefix=${var.private_subnet_ip_prefix}",
            "export prefix=${var.prefix}",
            "export default_key_name=${var.key_name}",
            "export region=${var.region}",
            "export zone=${var.zone}",
            "export kubernetes_cluster_tag=${random_id.kubernetes-cluster-tag.b64}",
            "EOF'",
            "sudo mkdir /share",
            "sudo chown ubuntu:ubuntu /share",
            "wget https://storage.googleapis.com/kubo-public/kubo-deployment-latest.tgz",
            "tar -xvf kubo-deployment-latest.tgz -C /share",
            "echo \"${file(var.private_key_filename)}\" > /home/ubuntu/deployer.pem",
            "chmod 600 /home/ubuntu/deployer.pem"
	]

        connection {
            type     = "ssh"
            user = "ubuntu"
            private_key = "${file(var.private_key_filename)}"
        }
    }

}

output "bosh-bastion-ip" {
    value = "${aws_instance.bastion.public_ip}"
}
