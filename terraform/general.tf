resource "azurerm_resource_group" "group" {
  name     = var.resource_group
  location = var.location
  tags     = local.common_tags
}