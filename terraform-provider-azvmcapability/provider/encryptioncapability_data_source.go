package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource = &encryptioncapabilityDataSource{}
)

type encryptioncapabilityDataSource struct{}

func NewEncryptioncapabilityDataSource() datasource.DataSource {
	return &encryptioncapabilityDataSource{}
}

type encryptioncapabilityModel struct {
	SubscriptionID types.String `tfsdk:"subscription_id"`
	Region         types.String `tfsdk:"region"`
	SKUName        types.String `tfsdk:"sku_name"`
	Supported      types.Bool   `tfsdk:"supported"`
}

func (d *encryptioncapabilityDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_encryptioncapability"
}

func (d *encryptioncapabilityDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"subscription_id": schema.StringAttribute{Required: true},
			"region":          schema.StringAttribute{Required: true},
			"sku_name":        schema.StringAttribute{Required: true},
			"supported":       schema.BoolAttribute{Computed: true},
		},
	}
}

func (d *encryptioncapabilityDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data encryptioncapabilityModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	token := strings.TrimRight(os.Getenv("AZURE_ACCESS_TOKEN"), "\r\n")
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
		headerString := ""
		for key, values := range reqHttp.Header {
			headerString += fmt.Sprintf("%s: %s\n", key, values)
		}
		resp.Diagnostics.AddError("API call failed", fmt.Sprintf("%s\nHeaders: %s", err.Error(), headerString))
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		resp.Diagnostics.AddError("Read response failed", err.Error())
		return
	}

	// Check HTTP status
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		resp.Diagnostics.AddError("API returned error", fmt.Sprintf("status: %d\nbody: %s\nheaders: %v", res.StatusCode, string(body), reqHttp.Header))
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		resp.Diagnostics.AddError("Invalid JSON response", err.Error())
		return
	}

	val, ok := result["value"]
	if !ok || val == nil {
		resp.Diagnostics.AddError("Unexpected API response", fmt.Sprintf("missing or null 'value' field in response: %s", string(body)))
		return
	}

	items, ok := val.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Unexpected API response", fmt.Sprintf("'value' is not an array (type %T): %s", val, string(body)))
		return
	}

	supported := false
	for _, item := range items {
		vm, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		name, _ := vm["name"].(string)
		resourceType, _ := vm["resourceType"].(string)
		if name == data.SKUName.ValueString() && resourceType == "virtualMachines" {
			locations, _ := vm["locations"].([]interface{})
			for _, locRaw := range locations {
				loc, _ := locRaw.(string)
				if loc == data.Region.ValueString() {
					capabilities, _ := vm["capabilities"].([]interface{})
					for _, capRaw := range capabilities {
						c, _ := capRaw.(map[string]interface{})
						if c == nil {
							continue
						}
						if n, _ := c["name"].(string); n == "EncryptionAtHostSupported" {
							if v, _ := c["value"].(string); v == "True" {
								supported = true
								break
							}
						}
					}
				}
				if supported {
					break
				}
			}
		}
		if supported {
			break
		}
	}

	data.Supported = types.BoolValue(supported)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
