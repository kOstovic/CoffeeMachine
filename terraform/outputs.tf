output "azure_instance_public_dns" {
  value       = azurerm_public_ip.ip[*].fqdn
  description = "Public DNS for VMs"
}