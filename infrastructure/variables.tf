variable "resource_group_name" {
  description = "The name of the resource group"
  type        = string
  default     = "wishlists-rg"
}

variable "location" {
  description = "The Azure region where resources will be created"
  type        = string
  default     = "Poland Central"
}

variable "storage_account_name" {
  description = "The name of the storage account"
  type        = string
  default     = "wishlistsstorageaccount"
}

variable "app_service_plan_name" {
  description = "The name of the App Service Plan"
  type        = string
  default     = "wishlists-app-service-plan"
}

variable "app_service_name" {
  description = "The name of the App Service"
  type        = string
  default     = "wishlists-app"
}

variable "state_resource_group" {
  description = "The name of the resource group"
  type        = string
  default     = "wishlists-rg"
}

variable "state_storage_account" {
  description = "The name of the storage account"
  type        = string
  default     = "wishlistsstorageaccount"
}

variable "state_container" {
  description = "The name of the storage container"
  type        = string
  default     = "wishlists-tfstate"
}

variable "application_image" {
  type = string
  default = "creeston/wishlist-app:latest"
}