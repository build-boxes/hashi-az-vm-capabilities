package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	// "github.com/build-boxes/hashi-az-vm-capabilities/terraform-provider-azvmcapability/provider"
	"terraform-provider-azvmcapability/provider"
)

func main() {
	err := providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/example/azvmcapability",
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}
