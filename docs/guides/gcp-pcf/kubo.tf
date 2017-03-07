// Easier mainteance for updating GCE image string
variable "latest_ubuntu" {
    type = "string"
    default = "ubuntu-1404-trusty-v20161109"
}

variable "projectid" {
    type = "string"
}

variable "region" {
    type = "string"
    default = "us-east1"
}

variable "zone" {
    type = "string"
    default = "us-east1-d"
}


variable "network" {
    type = "string"
}

variable "prefix" {
    type = "string"
    default = ""
}

variable "service_account_email" {
    type = "string"
    default = ""
}

variable "subnet_ip_prefix" {
    type = "string"
    default = "10.0.1"
}

provider "google" {
    project = "${var.projectid}"
    region = "${var.region}"
}

resource "google_compute_route" "nat-primary" {
  name        = "${var.prefix}nat-primary"
  dest_range  = "0.0.0.0/0"
  network       = "${var.network}"
  next_hop_instance = "${google_compute_instance.nat-instance-private-with-nat-primary.name}"
  next_hop_instance_zone = "${var.zone}"
  priority    = 800
  tags = ["no-ip"]
}

// Subnet for Kubo
resource "google_compute_subnetwork" "kubo-subnet" {
  name          = "${var.prefix}kubo-${var.region}"
  region        = "${var.region}"
  ip_cidr_range = "${var.subnet_ip_prefix}.0/24"
  network       = "https://www.googleapis.com/compute/v1/projects/${var.projectid}/global/networks/${var.network}"
}

// Allow SSH to BOSH bastion
resource "google_compute_firewall" "bosh-bastion" {
  name    = "${var.prefix}bosh-bastion"
  network = "${var.network}"

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  target_tags = ["bosh-bastion"]
}

// Allow all traffic within subnet
resource "google_compute_firewall" "intra-subnet-open" {
  name    = "${var.prefix}intra-subnet-open"
  network = "${var.network}"

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["1-65535"]
  }

  allow {
    protocol = "udp"
    ports    = ["1-65535"]
  }

  source_tags = ["internal"]
}

// BOSH bastion host
resource "google_compute_instance" "bosh-bastion" {
  name         = "${var.prefix}bosh-bastion"
  machine_type = "n1-standard-1"
  zone         = "${var.zone}"

  tags = ["bosh-bastion", "internal"]

  disk {
    image = "${var.latest_ubuntu}"
  }

  network_interface {
    subnetwork = "${google_compute_subnetwork.kubo-subnet.name}"
    access_config {
      // Ephemeral IP
    }
  }

  metadata_startup_script = <<EOT
#!/bin/bash
cat > /etc/motd <<EOF




#    #     ##     #####    #    #   #   #    #    ####
#    #    #  #    #    #   ##   #   #   ##   #   #    #
#    #   #    #   #    #   # #  #   #   # #  #   #
# ## #   ######   #####    #  # #   #   #  # #   #  ###
##  ##   #    #   #   #    #   ##   #   #   ##   #    #
#    #   #    #   #    #   #    #   #   #    #    ####

Startup scripts have not finished running, and the tools you need
are not ready yet. Please log out and log back in again in a few moments.
This warning will not appear when the system is ready.
EOF

apt-get update
apt-get install -y build-essential zlibc zlib1g-dev ruby ruby-dev openssl libxslt-dev libxml2-dev libssl-dev libreadline6 libreadline6-dev libyaml-dev libsqlite3-dev sqlite3 jq git unzip
gem install bosh_cli
curl -o /tmp/cf.tgz https://s3.amazonaws.com/go-cli/releases/v6.20.0/cf-cli_6.20.0_linux_x86-64.tgz
tar -zxvf /tmp/cf.tgz && mv cf /usr/bin/cf && chmod +x /usr/bin/cf
curl -o /usr/bin/bosh-init https://s3.amazonaws.com/bosh-init-artifacts/bosh-init-0.0.96-linux-amd64
chmod +x /usr/bin/bosh-init

cat > /etc/profile.d/bosh.sh <<'EOF'
#!/bin/bash
# Misc vars
export prefix=${var.prefix}
export ssh_key_path=$HOME/.ssh/bosh

# Vars from Terraform
export subnetwork=${google_compute_subnetwork.kubo-subnet.name}
export network=${var.network}


# Vars from metadata service
export project_id=$$(curl -s -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/project/project-id)
export zone=$$(curl -s -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/zone)
export zone=$${zone##*/}
export region=$${zone%-*}

# Configure gcloud
gcloud config set compute/zone $${zone}
gcloud config set compute/region $${region}
EOF

# Clone repo
mkdir /share
git clone https://github.com/cloudfoundry-incubator/bosh-google-cpi-release.git /share
chmod -R 777 /share

# Install Terraform
wget https://releases.hashicorp.com/terraform/0.7.7/terraform_0.7.7_linux_amd64.zip
unzip terraform*.zip -d /usr/local/bin
rm /etc/motd

cd
sudo curl https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-2.0.1-linux-amd64 -o /usr/bin/bosh-cli
sudo chmod a+x /usr/bin/bosh-cli
curl -L https://github.com/cloudfoundry-incubator/credhub-cli/releases/download/0.4.0/credhub-linux-0.4.0.tgz | tar zxv
chmod a+x credhub
sudo mv credhub /usr/bin
sudo curl -L https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl -o /usr/bin/kubectl
sudo chmod a+x /usr/bin/kubectl
EOT

  service_account {
    email = "${var.service_account_email}"
    scopes = ["cloud-platform"]
  }
}

// NAT server (primary)
resource "google_compute_instance" "nat-instance-private-with-nat-primary" {
  name         = "${var.prefix}nat-instance-primary"
  machine_type = "n1-standard-1"
  zone         = "${var.zone}"

  tags = ["nat", "internal"]

  disk {
    image = "${var.latest_ubuntu}"
  }

  network_interface {
    subnetwork = "${google_compute_subnetwork.kubo-subnet.name}"
    access_config {
      // Ephemeral IP
    }
  }

  can_ip_forward = true

  metadata_startup_script = <<EOT
#!/bin/bash
sh -c "echo 1 > /proc/sys/net/ipv4/ip_forward"
iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
EOT
}

output "kubo_subnet" {
   value = "${google_compute_subnetwork.kubo-subnet.name}"
}
