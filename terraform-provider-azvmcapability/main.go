package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"terraform-provider-azvmcapability/provider"
)

func main() {
	providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/example/azvmcapability",
	})
}