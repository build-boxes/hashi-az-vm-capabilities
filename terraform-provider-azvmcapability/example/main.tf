terraform {
  required_providers {
    azvmcapability = {
      source  = "registry.terraform.io/example/azvmcapability"
      version = "0.1.0"
    }
  }
}

provider "azvmcapability" {}

data "azvmcapability_encryptioncapability" "check" {
  subscription_id = var.subscription_id
  region          = var.region
  sku_name        = var.sku_name
}

output "encryption_supported" {
  value = data.azvmcapability_encryptioncapability.check.supported
}

