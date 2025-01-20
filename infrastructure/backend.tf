terraform {
  backend "azurerm" {
    resource_group_name  = ""
    storage_account_name = ""
    container_name       = ""
    key                  = "terraform.tfstate"
  }
}


provider "azurerm" {
  subscription_id = "4192f9b9-7bbf-4e11-8156-67f5431de563"
  features {}
}