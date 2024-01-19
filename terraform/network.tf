resource "azurerm_virtual_network" "this" {
  name                = "${var.resource_group}-virtual_network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.group.location
  resource_group_name = var.resource_group
  tags                = local.common_tags
}

resource "azurerm_subnet" "subnet" {
  name                 = "${var.resource_group}-${var.environment}${var.sufix != "" ? "-${var.sufix}" : ""}-subnet"
  resource_group_name  = azurerm_virtual_network.this.resource_group_name
  virtual_network_name = azurerm_virtual_network.this.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_network_security_group" "this" {
  name                = "${var.resource_group}-${var.environment}${var.sufix != "" ? "-${var.sufix}" : ""}-nsg"
  location            = azurerm_resource_group.group.location
  resource_group_name = var.resource_group

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

  tags = local.common_tags
}


resource "azurerm_public_ip" "ip" {
  name                    = lower("${var.resource_group}-${var.environment}${var.sufix != "" ? "-${var.sufix}" : ""}-pip-${count.index}")
  count                   = 1
  location                = azurerm_resource_group.group.location
  resource_group_name     = var.resource_group
  allocation_method       = "Static"
  sku                     = "Standard"
  idle_timeout_in_minutes = 30

  domain_name_label = lower("${var.resource_group}-${var.environment}${var.sufix != "" ? "-${var.sufix}" : ""}-${count.index}")
  tags              = local.common_tags
}

resource "azurerm_network_interface" "this" {
  name                = "${var.resource_group}-${var.environment}${var.sufix != "" ? "-${var.sufix}" : ""}-nic-${count.index}"
  count               = 1
  location            = azurerm_resource_group.group.location
  resource_group_name = var.resource_group

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.subnet.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.ip[count.index].id

  }
  tags = local.common_tags
}

resource "azurerm_network_interface_security_group_association" "this" {
  count                     = 1
  network_interface_id      = azurerm_network_interface.this[count.index].id
  network_security_group_id = azurerm_network_security_group.this.id
}