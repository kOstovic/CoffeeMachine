variable "project" {
  type        = string
  description = "project name"
  default     = "CoffeeMachine"
}

variable "client" {
  type        = string
  description = "client name"
}

variable "environment" {
  type        = string
  description = "environment name"
}

variable "sufix" {
  type        = string
  description = "sufix if needed - optional"
  default     = ""
}

variable "resource_group" {
  type        = string
  description = "resource_group name"
}

variable "username" {
  type        = string
  description = "username that will be added to VM as default admin"
  validation {
    condition     = length(var.username) >= 1 && length(var.username) <= 64
    error_message = "username must be between 1 and 64 characters"
  }
}

variable "instance_type_VM" {
  type        = string
  description = "Type for VM Instnace"
  default     = "Standard_B1s"
}

variable "location" {
  type        = string
  description = "Azure location"
  default     = "North Europe"
}