# Terraform Azure App

This project sets up an Azure infrastructure using Terraform. It creates an Azure Resource Group, a Storage Account, a Free App Service Plan, and an App Service. Additionally, it configures a Storage Account for storing the Terraform state.

## Prerequisites

- [Terraform](https://www.terraform.io/downloads.html) installed on your machine.
- An Azure account. You can create a free account [here](https://azure.microsoft.com/free/).

## Getting Started

1. **Clone the repository:**

   ```bash
   git clone <repository-url>
   cd terraform-azure-app
   ```

2. **Create storage account for terraform state**

   If you don't have resource group, you can create it using the following command:

   ```bash
   az group create --name <resource-group-name> --location <location>
   ```

   Create a storage account in Azure to store the Terraform state. You can do this manually in the Azure portal or use the following command:

   ```bash
   az storage account create --name <storage-account-name> --resource-group <resource-group-name> --location <location> --sku Standard_LRS
   ```

   Replace `<storage-account-name>`, `<resource-group-name>`, and `<location>` with your desired values.

   Create container in the storage account to store the Terraform state:

   Replace backend configuraiton in state.config file with the storage account name and resource group name.

3. **Initialize Terraform:**

   This command initializes the Terraform configuration and downloads the necessary provider plugins.

   ```bash
   terraform init -backend-config="./state.config"`
   ```

4. **Plan the deployment:**

   This command creates an execution plan, showing what actions Terraform will take to achieve the desired state.

   ```bash
   terraform plan
   ```

5. **Apply the configuration:**

   This command applies the changes required to reach the desired state of the configuration.

   ```bash
   terraform apply
   ```

   Confirm the action by typing `yes` when prompted.


## Cleanup

To remove all resources created by this Terraform configuration, run:

```bash
terraform destroy
```

Confirm the action by typing `yes` when prompted.