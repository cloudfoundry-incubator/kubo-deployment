variable "region" {
    type = "string"
    default = "ca-central-1"
}

variable "zone" {
    type = "string"
    default = "ca-central-1b"
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

variable "private_key" {
    type = "string"
}

variable "prefix" {
    type = "string"
    default = ""
}

provider "aws" {
    region = "${var.region}"
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
  enable_dns_hostnames = true
}

resource "aws_internet_gateway" "gateway" {
    vpc_id = "${aws_vpc.main.id}"
}

resource "aws_subnet" "public" {
    vpc_id     = "${aws_vpc.main.id}"
    cidr_block = "${var.public_subnet_ip_prefix}.0/24"
    availability_zone = "${var.zone}"

    tags {
      Name = "Kubo Public"
    }
}

resource "aws_route_table" "public" {
    vpc_id = "${aws_vpc.main.id}"

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

resource "aws_nat_gateway" "gateway" {
    allocation_id = "${aws_eip.nat.id}"
    subnet_id     = "${aws_subnet.public.id}"
}


resource "aws_subnet" "private" {
    vpc_id     = "${aws_vpc.main.id}"
    cidr_block = "${var.private_subnet_ip_prefix}.0/24"

    tags {
      Name = "Kubo Private"
    }
}

resource "aws_route_table" "private" {
    vpc_id = "${aws_vpc.main.id}"

    route {
      cidr_block = "0.0.0.0/0"
      gateway_id = "${aws_nat_gateway.gateway.id}"
    }
}

resource "aws_route_table_association" "private" {
    subnet_id      = "${aws_subnet.private.id}"
    route_table_id = "${aws_route_table.private.id}"
}

resource "aws_security_group" "nodes" {
    name        = "node access"
    vpc_id = "${aws_vpc.main.id}"

    ingress {
      from_port   = 8443
      to_port     = 8443
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    }

    egress {
      from_port       = 0
      to_port         = 0
      protocol        = "-1"
      cidr_blocks     = ["0.0.0.0/0"]
    }
}

resource "aws_security_group_rule" "ssh" {
    type            = "ingress"
    from_port       = 22
    to_port         = 22
    protocol        = "tcp"
    cidr_blocks     = ["0.0.0.0/0"]

    security_group_id = "${aws_security_group.nodes.id}"
}

resource "aws_security_group_rule" "nodes" {
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

resource "aws_instance" "bastion" {
    ami           = "${data.aws_ami.ubuntu.id}"
    instance_type = "t2.micro"
    subnet_id     = "${aws_subnet.public.id}"
    availability_zone = "${var.zone}"
    key_name      = "${var.key_name}"
    vpc_security_group_ids = ["${aws_security_group.nodes.id}"]
    associate_public_ip_address = true
    provisioner "remote-exec" {
        inline = [
            "sudo apt-get update",
            "sudo apt-get install -y build-essential zlibc zlib1g-dev ruby ruby-dev openssl libxslt-dev libxml2-dev libssl-dev libreadline6 libreadline6-dev libyaml-dev libsqlite3-dev sqlite3",
            "sudo apt-get install -y git",
            "curl -L https://github.com/cloudfoundry-incubator/credhub-cli/releases/download/1.0.0/credhub-linux-1.0.0.tgz | tar zxv && chmod a+x credhub && sudo mv credhub /usr/bin",
            "sudo curl -L https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl -o /usr/bin/kubectl && sudo chmod a+x /usr/bin/kubectl",
            "sudo curl https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-2.0.27-linux-amd64 -o /usr/bin/bosh-cli && sudo chmod a+x /usr/bin/bosh-cli",
            "sudo sh -c 'sudo cat > /etc/profile.d/bosh.sh <<'EOF'",
            "#!/bin/bash",
            "export subnet_id=${aws_subnet.private.id}",
            "export default_security_groups=${aws_security_group.nodes.id}",
            "export subnet_ip_prefix=${var.private_subnet_ip_prefix}",
            "export prefix=${var.prefix}",
	    "export region=${var.region}",
            "export zone=${var.zone}",
            "EOF'",
            "git clone https://www.github.com/cloudfoundry-incubator/kubo-deployment.git /home/ubuntu/kubo-deployment",
            "echo \"${var.private_key}\" > /home/ubuntu/deployer.pem",
	]

        connection {
            type     = "ssh"
            user = "ubuntu"
            private_key = "${var.private_key}"
        }
    }

}
