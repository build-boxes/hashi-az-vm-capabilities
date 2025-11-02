package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/build-boxes/hashi-az-vm-capabilities/terraform-provider-azvmcapability/provider"
)

func main() {
	providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/example/azvmcapability",
	})
}