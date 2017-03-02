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

// Subnet for Kubo 
resource "google_compute_subnetwork" "kubo-subnet" {
  name          = "${var.prefix}kubo-${var.kubo_region}"
  region        = "${var.kubo_region}"
  ip_cidr_range = "10.0.1.0/24"
  network       = "https://www.googleapis.com/compute/v1/projects/${var.projectid}/global/networks/${var.network}"
}

output "kubo_subnet" {
   value = "${google_compute_subnetwork.kubo-subnet.name}"
}

