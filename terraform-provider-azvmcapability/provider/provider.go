package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure the implementation satisfies the expected interfaces.
// var (
//     _ provider.Provider = &azvmcapabilityProvider{}
// )

func New() provider.Provider {
	return &azvmcapabilityProvider{}
}

type azvmcapabilityProvider struct{}

func (p *azvmcapabilityProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "azvmcapability"
}

func (p *azvmcapabilityProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *azvmcapabilityProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {}

func (p *azvmcapabilityProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewEncryptionCapabilityDataSource,
	}
}

func (p *azvmcapabilityProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}