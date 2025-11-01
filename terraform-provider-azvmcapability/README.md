# Terraform Custom Provider - azvmcapability
Checks if Azure VM SKU supports Encryption-At-Host. Uses a Terraform custom Provider to encapsulate Azure REST API calls.

## How to Run
1. Set your Azure token:
    ```
    export AZURE_ACCESS_TOKEN=$(az account get-access-token --query accessToken -o tsv)
    ```
1. Build the provider:
    ```
    go build -o terraform-provider-azvmcapability
    ```
1. Place it in your Terraform plugin directory.
1. Run:
    ```
    terraform init
    terraform apply
    ```


