variable "subscription_id" {
  description = "Azure subscription ID"
  type        = string
}

variable "region" {
  description = "Azure region (location)"
  type        = string
  default     = "eastus"
}

variable "sku_name" {
  description = "VM SKU name"
  type        = string
  default     = "Standard_DS2_v2"
}