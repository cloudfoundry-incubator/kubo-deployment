

variable "subscription_id" {}

variable "tenant_id" {}

variable "client_id" {}

variable "client_secret" {}

variable "latest_ubuntu" {
    type = "map"
    default = {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "14.04.5-LTS"
    version   = "latest"
    }
}

variable "ssh_user_username" {
   type = "string"
   default = "ubuntu"
}

variable "ssh_private_key_filename" {
    type = "string"
}

variable "ssh_public_key_filename" {
   type = "string"
}

variable "location" {
    type = "string"
    default = "eastus2"
}

variable "prefix" {
    type = "string"
}

variable "network_cidr" {
  default = "10.0.0.0/16"
}

provider "azurerm" {
  subscription_id = "${var.subscription_id}"
  tenant_id       = "${var.tenant_id}"
  client_id       = "${var.client_id}"
  client_secret   = "${var.client_secret}"
}

resource "azurerm_resource_group" "bosh" {
  name     = "${var.prefix}-cfcr"
  location = "${var.location}"

  tags {
    environment = "${var.prefix}-cfcr"
  }
}

resource "azurerm_public_ip" "bosh-bastion" {
  name                         = "${var.prefix}-cfcr-ip"
  location                     = "${var.location}"
  depends_on  = ["azurerm_resource_group.bosh"]
  resource_group_name          = "${azurerm_resource_group.bosh.name}"
  public_ip_address_allocation = "static"

  tags {
    environment = "${var.prefix}-cfcr"
  }
}

// Subnet for CFCR
resource "azurerm_virtual_network" "cfcr-vnet" {
  name          = "${var.prefix}-cfcr-vnet"
  location      = "${var.location}"
  depends_on  = ["azurerm_resource_group.bosh"]

  resource_group_name = "${azurerm_resource_group.bosh.name}"
  address_space = ["${var.network_cidr}"]
  dns_servers   = ["168.63.129.16"]
 
}

resource "azurerm_subnet" "cfcr-subnet" {
    name = "cfcr-subnet" 
    depends_on  = ["azurerm_virtual_network.cfcr-vnet"]

    resource_group_name = "${azurerm_resource_group.bosh.name}"
    virtual_network_name = "${azurerm_virtual_network.cfcr-vnet.name}"
    address_prefix = "${cidrsubnet(var.network_cidr, 8, 0)}"
}


// Allow SSH to BOSH bastion
resource "azurerm_network_security_group" "bosh-bastion" {
  name    = "${var.prefix}bosh-bastion"
  location      = "${var.location}"
  resource_group_name = "${azurerm_resource_group.bosh.name}"
  depends_on  = ["azurerm_resource_group.bosh"]
 security_rule {
    name                       = "ssh"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "22"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}

// Allow port 8443 to master
resource "azurerm_network_security_group" "cfcr-master" {
  name    = "${var.prefix}cfcr-master"
  location      = "${var.location}"
  resource_group_name = "${azurerm_resource_group.bosh.name}"
  depends_on  = ["azurerm_resource_group.bosh"]
 security_rule {
    name                       = "master"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "8443"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}


// BOSH bastion host

resource "azurerm_network_interface" "bosh-bastion" {
  name                      = "${var.prefix}bosh-bastion-nic"
  depends_on                = ["azurerm_public_ip.bosh-bastion", "azurerm_subnet.cfcr-subnet", "azurerm_network_security_group.bosh-bastion"]
  location                  = "${var.location}"
  resource_group_name       = "${azurerm_resource_group.bosh.name}"
  network_security_group_id = "${azurerm_network_security_group.bosh-bastion.id}"

  ip_configuration {
    name                          = "${var.prefix}-bosh-bastion-ip-config"
    subnet_id                     = "${azurerm_subnet.cfcr-subnet.id}"
    private_ip_address_allocation = "static"
    private_ip_address            = "${cidrhost(azurerm_subnet.cfcr-subnet.address_prefix,4)}"
    public_ip_address_id          = "${azurerm_public_ip.bosh-bastion.id}"
  }
}


resource "azurerm_virtual_machine" "bosh-bastion" {
  name         = "${var.prefix}bosh-bastion"
  depends_on   = ["azurerm_network_interface.bosh-bastion"]
  vm_size      = "Standard_D2_V2"
  location      = "${var.location}"
  resource_group_name = "${azurerm_resource_group.bosh.name}"
  network_interface_ids = ["${azurerm_network_interface.bosh-bastion.id}"]
  storage_image_reference = ["${var.latest_ubuntu}"]

  storage_os_disk {
    name              = "osdisk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
    disk_size_gb = "100"
  }

  os_profile_linux_config {
    disable_password_authentication = true
    ssh_keys = [{
      path     = "/home/${var.ssh_user_username}/.ssh/authorized_keys"
      key_data = "${file(var.ssh_public_key_filename)}"
    }]
  }

   os_profile {
    computer_name = "bosh-bastion"
    admin_username = "${var.ssh_user_username}"
    custom_data = <<EOT
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
curl -o /tmp/cf.tgz https://s3.amazonaws.com/go-cli/releases/v6.20.0/cf-cli_6.20.0_linux_x86-64.tgz
tar -zxvf /tmp/cf.tgz && mv cf /usr/bin/cf && chmod +x /usr/bin/cf

cat > /etc/profile.d/bosh.sh <<'EOF'
#!/bin/bash
# Misc vars
export prefix=${var.prefix}
export ssh_key_path=$HOME/.ssh/bosh

# Vars from Terraform
export resource_group=${azurerm_resource_group.bosh.name}
export subnetwork=${azurerm_subnet.cfcr-subnet.name}
export network=${azurerm_virtual_network.cfcr-vnet.name}
export subnet_ip_prefix=${azurerm_subnet.cfcr-subnet.address_prefix}
export subscription_id=${var.subscription_id}
export tenant_id=${var.tenant_id}
export client_id=${var.client_id}
export client_secret=${var.client_secret}1
export location=${var.location}

EOF

cat > /usr/bin/update_azure_env <<'EOF'
#!/bin/bash

if [[ ! -f "$1" ]] || [[ ! "$1" =~ director.yml$ ]]; then
  echo 'Please specify the path to director.yml'
  exit 1
fi

# Azure specific updates
sed -i -e 's/^\(resource_group_name:\).*\(#.*\)/\1 ${azurerm_resource_group.bosh.name} \2/' "$1"
sed -i -e 's/^\(vnet_resource_group_name:\).*\(#.*\)/\1 ${azurerm_resource_group.bosh.name} \2/' "$1"
sed -i -e 's/^\(vnet_name:\).*\(#.*\)/\1 ${azurerm_virtual_network.cfcr-vnet.name} \2/' "$1"
sed -i -e 's/^\(subnet_name:\).*\(#.*\)/\1 ${azurerm_subnet.cfcr-subnet.name} \2/' "$1"
sed -i -e 's/^\(location:\).*\(#.*\)/\1 ${var.location} \2/' "$1"
sed -i -e 's/^\(default_security_group:\).*\(#.*\)/\1 ${azurerm_network_security_group.cfcr-master.name} \2/' "$1"

# Generic updates
sed -i -e 's/^\(internal_ip:\).*\(#.*\)/\1 ${cidrhost(azurerm_subnet.cfcr-subnet.address_prefix, 5)} \2/' "$1"
sed -i -e 's=^\(internal_cidr:\).*\(#.*\)=\1 ${azurerm_subnet.cfcr-subnet.address_prefix} \2=' "$1"
sed -i -e 's/^\(internal_gw:\).*\(#.*\)/\1 ${cidrhost(azurerm_subnet.cfcr-subnet.address_prefix, 1)} \2/' "$1"
sed -i -e 's/^\(director_name:\).*\(#.*\)/\1 ${var.prefix}bosh \2/' "$1"

EOF
chmod a+x /usr/bin/update_azure_env

cat > /usr/bin/update_azure_secrets <<'EOF'
#!/bin/bash

if [[ ! -f "$1" ]] || [[ ! "$1" =~ director-secrets.yml$ ]]; then
  echo 'Please specify the path to director-secrets.yml'
  exit 1
fi

# Azure secrets updates
sed -i -e 's/^\(subscription_id:\).*\(#.*\)/\1 ${var.subscription_id} \2/' "$1"
sed -i -e 's=^\(tenant_id:\).*\(#.*\)=\1 ${var.tenant_id} \2=' "$1"
sed -i -e 's/^\(client_id:\).*\(#.*\)/\1 ${var.client_id} \2/' "$1"
sed -i -e 's/^\(client_secret:\).*\(#.*\)/\1 ${var.client_secret} \2/' "$1"

EOF
chmod a+x /usr/bin/update_azure_secrets


cat > /usr/bin/set_iaas_routing <<'EOF'
#!/bin/bash

if [[ ! -f "$1" ]] || [[ ! "$1" =~ director.yml$ ]]; then
  echo 'Please specify the path to director.yml'
  exit 1
fi

sed -i -e 's/^#* *\(routing_mode:.*\)$/# \1/' "$1"
sed -i -e 's/^#* *\(routing_mode:\) *\(iaas\).*$/\1 \2/' "$1"

sed -i -e "s/^\(kubernetes_master_host:\).*\(#.*\)/\1 $${kubernetes_master_host} \2/" "$1"
sed -i -e "s/^\(kubernetes_master_port:\).*\(#.*\)/\1 $${kubernetes_master_port:-8443} \2/" "$1"
sed -i -e "s/^\(master_target_pool:\).*\(#.*\).*$/\1 $${master_target_pool} \2/" "$1"

EOF
chmod a+x /usr/bin/set_iaas_routing

# Get kubo-deployment
wget https://s3.amazonaws.com/scharlton-piv/kubo-deployment-latest.tgz
mkdir /share
tar -xvf kubo-deployment-latest.tgz -C /share
chmod -R 777 /share

# Install Terraform
wget https://releases.hashicorp.com/terraform/0.7.7/terraform_0.7.7_linux_amd64.zip
unzip terraform*.zip -d /usr/local/bin
rm /etc/motd

cd
sudo curl https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-2.0.48-linux-amd64 -o /usr/bin/bosh
sudo chmod a+x /usr/bin/bosh
sudo curl -L https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl -o /usr/bin/kubectl
sudo chmod a+x /usr/bin/kubectl
curl -L https://github.com/cloudfoundry-incubator/credhub-cli/releases/download/1.4.0/credhub-linux-1.4.0.tgz | tar zxv
chmod a+x credhub
sudo mv credhub /usr/bin

EOT
 }
}

output "kubo_subnet" {
   value = "${azurerm_subnet.cfcr-subnet.name}"
}
            
