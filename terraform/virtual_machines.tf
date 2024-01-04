resource "azurerm_linux_virtual_machine" "vms" {
  name                = "${var.resource_group}-${var.environment}${var.sufix != "" ? "-${var.sufix}" : ""}-${count.index}"
  resource_group_name = var.resource_group
  location            = azurerm_resource_group.group.location
  size                = var.instance_type_VM
  count               = 1
  admin_username      = var.username
  network_interface_ids = [
    azurerm_network_interface.this[count.index].id,
  ]

  admin_ssh_key {
    username   = var.username
    public_key = file("~/.ssh/id_rsa.pub")
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Premium_LRS"
    disk_size_gb         = "20"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts-gen2"
    version   = "latest"
  }
  vtpm_enabled        = true
  secure_boot_enabled = true

  provisioner "file" {
    source      = "../deployments/DockerCompose"
    destination = "/app"
  }

  provisioner "remote-exec" {
    inline = [
      "echo 'Running CoffeeMachine app over docker-compose'",
      "cd /app",
      "docker-compose up"
    ]
    connection {
      type        = "ssh"
      user        = var.username
      host        = azurerm_public_ip.ip[count.index].ip_address
      private_key = file("~/.ssh/id_rsa.pub")
    }
  }

  tags = local.common_tags
}
