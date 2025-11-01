package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type encryptionCapabilityDataSource struct{}

func NewEncryptionCapabilityDataSource() datasource.DataSource {
	return &encryptionCapabilityDataSource{}
}

type encryptionCapabilityModel struct {
	SubscriptionID types.String `tfsdk:"subscription_id"`
	Region         types.String `tfsdk:"region"`
	SKUName        types.String `tfsdk:"sku_name"`
	Supported      types.Bool   `tfsdk:"supported"`
}

func (d *encryptionCapabilityDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "azvm_encryption_capability"
}

func (d *encryptionCapabilityDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"subscription_id": schema.StringAttribute{Required: true},
			"region":          schema.StringAttribute{Required: true},
			"sku_name":        schema.StringAttribute{Required: true},
			"supported":       schema.BoolAttribute{Computed: true},
		},
	}
}

func (d *encryptionCapabilityDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data encryptionCapabilityModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	token := os.Getenv("AZURE_ACCESS_TOKEN")
	if token == "" {
		resp.Diagnostics.AddError("Missing token", "Set AZURE_ACCESS_TOKEN env variable")
		return
	}

	url := fmt.Sprintf("https://management.azure.com/subscriptions/%s/providers/Microsoft.Compute/skus?api-version=2022-08-01", data.SubscriptionID.ValueString())
	reqHttp, _ := http.NewRequest("GET", url, nil)
	reqHttp.Header.Add("Authorization", "Bearer "+token)
	reqHttp.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(reqHttp)
	if err != nil {
		resp.Diagnostics.AddError("API call failed", err.Error())
		return
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	supported := false
	for _, item := range result["value"].([]interface{}) {
		vm := item.(map[string]interface{})
		if vm["name"] == data.SKUName.ValueString() &&
			vm["resourceType"] == "virtualMachines" {
			for _, loc := range vm["locations"].([]interface{}) {
				if loc.(string) == data.Region.ValueString() {
					for _, cap := range vm["capabilities"].([]interface{}) {
						c := cap.(map[string]interface{})
						if c["name"] == "EncryptionAtHostSupported" && c["value"] == "True" {
							supported = true
						}
					}
				}
			}
		}
	}

	data.Supported = types.BoolValue(supported)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}