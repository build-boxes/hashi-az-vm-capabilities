# HashiCorp Az VM Capability

1. Bash Script - AZ VM Capability
2. Terraform Custom Provider - azvmcapability

## Bash Script - AZ VM Capability
Checks if Azure VM SKU supports Encryption-At-Host. Uses direct calls to Azure REST API.

## Terraform Custom Provider - azvmcapability
Checks if Azure VM SKU supports Encryption-At-Host. Uses a Terraform custom Provider to encapsulate Azure REST API calls.

### How to Run
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

### How to Build
1. Initialize Go Module.
    Run this in the root of your provider directory:
    ```
    go mod init github.com/yourname/terraform-provider-azvmcapability
    go mod init github.com/build-boxes/hashi-az-vm-capabilities
    go mod init github.com/build-boxes/hashi-az-vm-capabilities/terraform-provider-azvmcapability
    go mod init terraform-provider-azvmcapability
    ```
    Replace "yourname" with your Github or local namespace.

1. After Any Old Updates - Tidy Up
    Run:
    ```
    go mod tidy
    ```

1. Add Required Dependencies
    Install the Terraform Plugin Framework:
    ```
    go get github.com/hashicorp/terraform-plugin-framework@latest
    ```
    You can also pin a specific version if needed:
    ```
    go get github.com/hashicorp/terraform-plugin-framework@v1.0.0
    ```
  
    Install some more dependencies:
    ```
    go get github.com/hashicorp/terraform-plugin-framework/internal/fwserver@v1.16.1
    go get github.com/hashicorp/terraform-plugin-framework/internal/logging@v1.16.1
    go get github.com/hashicorp/terraform-plugin-framework/function@v1.16.1
    go get github.com/hashicorp/terraform-plugin-framework/providerserver@v1.16.1
    ```
1. Build the Provider Binary
    Compile the provider:
    ```
    go build -o terraform-provider-azvmcapability
    ```
    This will generate a binary named "terraform-provider-azvmcapability" in your current directory.
1. Install Locally for Terraform
    Terraform looks for custom providers in a specific directory structure. Create this under your home directory:
    ```
    mkdir -p ~/.terraform.d/plugins/example/azvmcapability/0.1.0/linux_amd64
    mv terraform-provider-azvmcapability ~/.terraform.d/plugins/example/azvmcapability/0.1.0/linux_amd64/terraform-provider-example_azvmcapability
    ```
    The path must match the source and version you declare in your Terraform config:
    ```Hcl
    terraform {
      required_providers {
        azvmcapability = {
          source  = "example/azvmcapability"
          version = "0.1.0"
        }
      }
    }
    ```
1. Test It
    Navigate to your example/ directory:
    ```
    cd example
    terraform init
    terraform apply
    ```
    Make sure you have exported youre Azure token:
    ```
    export AZURE_ACCESS_TOKEN=$(az account get-access-token --query accessToken -o tsv)
    ```
1. Update .terraformrc For Local Development
    The Terraform Init script should contain something like:
    ```
    provider_installation {
        dev_overrides {
            // "registry.terraform.io/example/azvmcapability" =  "/home/<USERNAME>/.terraform.d/plugins/example/azvmcapability/0.1.0/linux_amd64"
            "registry.terraform.io/example/azvmcapability" = "/home/<USERNAME>/go/bin"
        }

        direct {}
    }

    ```
1. (Optional:) Add a Makefile
    To simplify builds:
    ```Makefile
    build:
    	go build -o terraform-provider-azvmcapability

    install:
    	mkdir -p ~/.terraform.d/plugins/example/azvmcapability/0.1.0/linux_amd64
    	mv terraform-provider-azvmcapability ~/.terraform.d/plugins/example/azvmcapability/0.1.0/linux_amd64/terraform-provider-example_azvmcapability
    ```
    Then run:
    ```Bash
    make build install
    ```

## Reference Links
1. [https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework)


