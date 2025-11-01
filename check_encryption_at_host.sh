#!/bin/bash

# Description:
#   Checks if a given Azure VM SKU supports EncryptionAtHost in a specific region using Azure REST API.
# Usage:
#   ./check_encryption_at_host.sh <subscription_id> <region> <vm_sku>
# Example:
#   ./check_encryption_at_host.sh 00000000-0000-0000-0000-000000000000 eastus Standard_D2s_v3

set -euo pipefail

# Input parameters
SUBSCRIPTION_ID="${1:-}"
REGION="${2:-}"
VM_SKU="${3:-}"

# Azure REST API version
API_VERSION="2022-08-01"

# Validate input
if [[ -z "$SUBSCRIPTION_ID" || -z "$REGION" || -z "$VM_SKU" ]]; then
  echo "Usage: $0 <subscription_id> <region> <vm_sku>"
  exit 1
fi

# Function: Get Azure access token
get_access_token() {
  az account get-access-token --query accessToken -o tsv
}

# Function: Query VM SKUs from Azure REST API
query_vm_skus() {
  local token="$1"
  curl -sS -H "Authorization: Bearer $token" \
       -H "Content-Type: application/json" \
       "https://management.azure.com/subscriptions/$SUBSCRIPTION_ID/providers/Microsoft.Compute/skus?api-version=$API_VERSION"
}

# Function: Check EncryptionAtHostSupported capability
check_encryption_capability() {
  local json="$1"
  local result
  result=$(echo "$json" | jq -r --arg sku "$VM_SKU" --arg region "$REGION" '
    .value[]
    | select(.resourceType == "virtualMachines")
    | select(.name == $sku)
    | select(.locations[] | ascii_downcase == ($region | ascii_downcase))
    | .capabilities[]
    | select(.name == "EncryptionAtHostSupported")
    | .value
  ')

  if [[ "$result" == "True" ]]; then
    echo "✅ $VM_SKU supports EncryptionAtHost in $REGION."
  elif [[ "$result" == "False" ]]; then
    echo "❌ $VM_SKU does NOT support EncryptionAtHost in $REGION."
  else
    echo "⚠️ Capability not listed for $VM_SKU in $REGION. Possibly unsupported or unavailable."
  fi
}

# Main execution
main() {
  echo "🔍 Checking EncryptionAtHost support for SKU '$VM_SKU' in region '$REGION'..."
  local token
  token=$(get_access_token)
  local json
  json=$(query_vm_skus "$token")
  check_encryption_capability "$json"
}

main