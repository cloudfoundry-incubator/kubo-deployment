variable "region" {
    type = "string"
}

variable "vpc_id" {
    type = "string"
}

variable "node_security_group_id" {
    type = "string"
}

variable "public_subnet_id" {
    type = "string"
}

variable "prefix" {
    type = "string"
}

provider "aws" {
    region = "${var.region}"
}

resource "aws_security_group" "api" {
    name        = "${var.prefix}api-access"
    vpc_id = "${var.vpc_id}"

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

resource "aws_security_group_rule" "api" {
    type            = "ingress"
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    source_security_group_id = "${aws_security_group.api.id}"

    security_group_id = "${var.node_security_group_id}"
}

resource "aws_elb" "api" {
    name               = "${var.prefix}kubo-api"
    subnets = ["${var.public_subnet_id}"]
    security_groups = ["${aws_security_group.api.id}"]

    listener {
      instance_port      = 8443
      instance_protocol  = "tcp"
      lb_port            = 8443
      lb_protocol        = "tcp"
    }

    health_check {
      healthy_threshold   = 2
      unhealthy_threshold = 2
      timeout             = 2
      target              = "TCP:8443/"
      interval            = 5
    }
}

resource "aws_security_group" "apps" {
    name        = "${var.prefix}apps-access"
    vpc_id = "${var.vpc_id}"

    ingress {
      from_port   = 30000
      to_port     = 32000
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

resource "aws_security_group_rule" "apps" {
    type            = "ingress"
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    source_security_group_id = "${aws_security_group.apps.id}"

    security_group_id = "${var.node_security_group_id}"
}

resource "aws_elb" "apps" {
    name               = "${var.prefix}kubo-apps"
    subnets = ["${var.public_subnet_id}"]
    security_groups = ["${aws_security_group.apps.id}"]

    listener {
      instance_port      = 31000
      instance_protocol  = "tcp"
      lb_port            = 31000
      lb_protocol        = "tcp"
    }

    health_check {
      healthy_threshold   = 2
      unhealthy_threshold = 2
      timeout             = 2
      target              = "TCP:22/"
      interval            = 5
    }
}

output "kubo_master_target_pool" {
   value = "${aws_elb.api.name}"
}

output "master_lb_ip_address" {
  value = "${aws_elb.api.dns_name}"
}

output "kubo_worker_target_pool" {
   value = "${aws_elb.apps.name}"
}
