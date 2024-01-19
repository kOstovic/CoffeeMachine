terraform {
  required_version = ">= 1.6"

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.87"
    }
  }

  #backend "azurerm" {
  #  resource_group_name  = "CoffeeMachine"
  #  storage_account_name = "kostovic"
  #  container_name       = "terraformstate-{environment}"
  #  key                  = "CoffeeMachine-{environment}-{clientName}.tfstate"
  #}

}
