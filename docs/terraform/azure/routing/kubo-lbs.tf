
variable "location" {
	type = "string"
	default = "eastus2"
}

variable "prefix" {
    type = "string"
}

variable "kubernetes_master_port" {
	type = "string"
	default = "8443"
}

variable "subscription_id" {}

variable "tenant_id" {}

variable "client_id" {}

variable "client_secret" {}


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
    environment = "${var.prefix}"
  }
}

// Static IP address for HTTP forwarding rule
resource "azurerm_public_ip" "cfcr-tcp" {
  name = "${var.prefix}-cfcr"
  location = "${var.location}"
  resource_group_name = "${azurerm_resource_group.bosh.name}"
   public_ip_address_allocation = "static"

  tags {
    environment = "${var.prefix}-cfcr"
  }
}

resource "azurerm_lb" "cfcr-tcp" {
  name                = "${var.prefix}-master-lb"
  location            = "${var.location}"
  resource_group_name = "${azurerm_resource_group.bosh.name}"
  depends_on = ["azurerm_public_ip.cfcr-tcp"]
  frontend_ip_configuration {
    name                 = "PublicIPAddress"
    public_ip_address_id = "${azurerm_public_ip.cfcr-tcp.id}"
  }
}


resource "azurerm_lb_backend_address_pool" "cfcr-pool" {
  resource_group_name = "${azurerm_resource_group.bosh.name}"
  location            = "${var.location}"
  depends_on = ["azurerm_lb.cfcr-tcp"]
  loadbalancer_id     = "${azurerm_lb.cfcr-tcp.id}"
  name                = "BackEndAddressPool"
}

resource "azurerm_lb_probe" "cfcr-probe" {
  resource_group_name = "${azurerm_resource_group.bosh.name}"
  location            = "${var.location}"
  loadbalancer_id     = "${azurerm_lb.cfcr-tcp.id}"
    depends_on = ["azurerm_lb.cfcr-tcp"]
  name                = "probe-master-api"
  port                = 8443
  protocol            = "Tcp"
}

resource "azurerm_lb_rule" "cfcr-tcp" {
  name        = "api-access"
  location            = "${var.location}"
  resource_group_name = "${azurerm_resource_group.bosh.name}"
  depends_on = [ "azurerm_lb_probe.cfcr-probe", "azurerm_lb_backend_address_pool.cfcr-pool"  ]
  loadbalancer_id      = "${azurerm_lb.cfcr-tcp.id}"
  frontend_port  = "${var.kubernetes_master_port}"
  backend_address_pool_id = "${azurerm_lb_backend_address_pool.cfcr-pool.id}"
  backend_port  = "${var.kubernetes_master_port}"
  protocol = "Tcp"
  frontend_ip_configuration_name  = "PublicIPAddress"
  probe_id = "${azurerm_lb_probe.cfcr-probe.id}"
}

output "cfcr_master_target_pool" {
   value = "${azurerm_lb.cfcr-tcp.name}"
}

output "master_lb_ip_address" {
  value = "${azurerm_public_ip.cfcr-tcp.ip_address}"
}
