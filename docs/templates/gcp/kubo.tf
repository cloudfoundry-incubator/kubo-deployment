variable "projectid" {
    type = "string"
}

variable "kubo_region" {
	type = "string"
	default = "us-west1"
}

variable "ip_cidr_range" {
    type = "string"
    default = "10.0.1.0/24"
}

variable "network" {
    type = "string"
    default = "bosh"
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


// Subnet for Kubo
resource "google_compute_subnetwork" "kubo-subnet" {
  name          = "${var.prefix}kubo-${var.kubo_region}"
  region        = "${var.kubo_region}"
  ip_cidr_range = "${var.ip_cidr_range}"
  network       = "https://www.googleapis.com/compute/v1/projects/${var.projectid}/global/networks/${var.network}"
}

output "kubo_subnet" {
   value = "${google_compute_subnetwork.kubo-subnet.name}"
}