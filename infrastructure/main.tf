locals {
  mount_path = "/home/site/wwwroot/database"
}

resource "azurerm_resource_group" "main" {
  name     = var.resource_group_name
  location = var.location
}

resource "azurerm_storage_account" "state" {
  name                     = var.storage_account_name
  resource_group_name      = azurerm_resource_group.main.name
  location                 = azurerm_resource_group.main.location
  account_tier            = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "wishlist_share" {
  name                 = "wishlist-share"
  storage_account_id = azurerm_storage_account.state.id
  quota                = 1  
}

resource "azurerm_service_plan" "main" {
  name                = var.app_service_plan_name
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  sku_name = "F1"
  os_type = "Linux"

}

resource "azurerm_linux_web_app" "main" {
  name                = var.app_service_name
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  service_plan_id = azurerm_service_plan.main.id

  storage_account {
    name = azurerm_storage_account.state.name
    type = "AzureFiles"
    share_name = azurerm_storage_share.wishlist_share.name
    account_name = azurerm_storage_account.state.name
    access_key = azurerm_storage_account.state.primary_access_key
    mount_path = local.mount_path
  }
  
  site_config {
    always_on = false
    application_stack {
     docker_image_name = var.application_image
     docker_registry_url = "https://index.docker.io"
    }
  }

  app_settings = {
    "BASE_URL" = "https://${var.app_service_name}.azurewebsites.net/"
    "PORT" = "80"
    "MAX_ITEMS_COUNT" = "10"
    "MAX_BODY_SIZE" = "10K"
    "MAX_WISHLISTS_PER_DAY" = "3"
    "MAX_ITEM_LENGTH" = "20"
    "USE_IN_MEMORY_DB": "False"
    "SQLITE_DB_NAME": "${local.mount_path}/sqlite.db"
  }
}
