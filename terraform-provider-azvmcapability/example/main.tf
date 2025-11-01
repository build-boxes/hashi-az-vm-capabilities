terraform {
  required_providers {
    azvmcapability = {
      source  = "example/azvmcapability"
      version = "0.1.0"
    }
  }
}

provider "azvmcapability" {}

data "azvm_encryption_capability" "check" {
  subscription_id = "00000000-0000-0000-0000-000000000000"
  region          = "eastus"
  sku_name        = "Standard_D2s_v3"
}

output "encryption_supported" {
  value = data.azvm_encryption_capability.check.supported
}