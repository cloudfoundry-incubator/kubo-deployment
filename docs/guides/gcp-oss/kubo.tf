variable "projectid" {
    type = "string"
}

variable "kubo_region" {
	type = "string"
	default = "us-west1"
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
    region = "${var.kubo_region}"
}

// Static IP address for HTTP forwarding rule
resource "google_compute_address" "kubo-tcp" {
  name = "${var.prefix}kubo"
}

// TCP Load Balancer
resource "google_compute_target_pool" "kubo-tcp-public" {
    name = "${var.prefix}kubo-tcp-public"
}

resource "google_compute_forwarding_rule" "kubo-tcp" {
  name        = "${var.prefix}kubo-tcp"
  target      = "${google_compute_target_pool.kubo-tcp-public.self_link}"
  port_range  = "8443"
  ip_protocol = "TCP"
  ip_address  = "${google_compute_address.kubo-tcp.address}"
}

resource "google_compute_firewall" "kubo-tcp-public" {
  name    = "${var.prefix}kubo-tcp-public"
  network       = "${var.network}"

  allow {
    protocol = "tcp"
    ports    = ["8443"]
  }

  target_tags = ["master"]
}


// Subnet for Kubo 
resource "google_compute_subnetwork" "kubo-subnet" {
  name          = "${var.prefix}kubo-${var.kubo_region}"
  region        = "${var.kubo_region}"
  ip_cidr_range = "10.0.1.0/24"
  network       = "https://www.googleapis.com/compute/v1/projects/${var.projectid}/global/networks/${var.network}"
}

// Static IP address for HTTP forwarding rule
resource "google_compute_address" "kubo-workers-tcp" {
  name = "${var.prefix}kubo-workers"
}

// TCP Load Balancer
resource "google_compute_target_pool" "kubo-workers-tcp-public" {
    name = "${var.prefix}kubo-workers-tcp-public"
}

resource "google_compute_forwarding_rule" "kubo-workers-tcp" {
  name        = "${var.prefix}kubo-workers-tcp"
  target      = "${google_compute_target_pool.kubo-workers-tcp-public.self_link}"
  port_range  = "8443"
  ip_protocol = "TCP"
  ip_address  = "${google_compute_address.kubo-workers-tcp.address}"
}

resource "google_compute_firewall" "kubo-workers-tcp-public" {
  name    = "${var.prefix}kubo-workers-tcp-public"
  network       = "${var.network}"

  allow {
    protocol = "tcp"
    ports    = ["30000-40000"]
  }

  target_tags = ["worker"]
}

output "kubo_subnet" {
   value = "${google_compute_subnetwork.kubo-subnet.name}"
}

output "kubo_master_target_pool" {
   value = "${google_compute_target_pool.kubo-tcp-public.name}"
}

output "master_lb_ip_address" {
  value = "${google_compute_address.kubo-tcp.address}"
}

output "kubo_worker_target_pool" {
   value = "${google_compute_target_pool.kubo-workers-tcp-public.name}"
}
