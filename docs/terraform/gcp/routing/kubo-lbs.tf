variable "projectid" {
    type = "string"
}

variable "region" {
	type = "string"
	default = "us-west1"
}

variable "ip_cidr_range" {
    type = "string"
    default = "10.0.1.0/24"
}

variable "network" {
    type = "string"
}

variable "prefix" {
    type = "string"
    default = ""
}

provider "google" {
    credentials = ""
    project = "${var.projectid}"
    region = "${var.region}"
}

// Static IP address for HTTP forwarding rule
resource "google_compute_address" "cfcr-tcp" {
  name = "${var.prefix}cfcr"
}

// TCP Load Balancer
resource "google_compute_target_pool" "cfcr-tcp-public" {
    region = "${var.region}"
    name = "${var.prefix}cfcr-tcp-public"
}

resource "google_compute_forwarding_rule" "cfcr-tcp" {
  name        = "${var.prefix}cfcr-tcp"
  target      = "${google_compute_target_pool.cfcr-tcp-public.self_link}"
  port_range  = "8443"
  ip_protocol = "TCP"
  ip_address  = "${google_compute_address.cfcr-tcp.address}"
}

resource "google_compute_firewall" "cfcr-tcp-public" {
  name    = "${var.prefix}cfcr-tcp-public"
  network       = "${var.network}"

  allow {
    protocol = "tcp"
    ports    = ["8443"]
  }

  target_tags = ["master"]
}

output "cfcr_master_target_pool" {
   value = "${google_compute_target_pool.cfcr-tcp-public.name}"
}

output "master_lb_ip_address" {
  value = "${google_compute_address.cfcr-tcp.address}"
}
