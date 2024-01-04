locals {
  instance_prefix = ""
  common_tags = {
    Part-of = "${var.project}-${var.environment}${var.sufix != "" ? "-${var.sufix}" : ""}"
    Product = "CoffeeMachine"
    Client  = var.client
  }
}