package service

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func testProvider() provider.Provider {
	return &testOPNsenseProvider{}
}

type testOPNsenseProvider struct{}

func (p *testOPNsenseProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "opnsense"
	resp.Version = "test"
}

func (p *testOPNsenseProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"uri": schema.StringAttribute{
				Optional: true,
			},
			"api_key": schema.StringAttribute{
				Optional: true,
			},
			"api_secret": schema.StringAttribute{
				Optional: true,
			},
			"allow_insecure": schema.BoolAttribute{
				Optional: true,
			},
		},
	}
}

func (p *testOPNsenseProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	uri := os.Getenv("OPNSENSE_URI")
	apiKey := os.Getenv("OPNSENSE_API_KEY")
	apiSecret := os.Getenv("OPNSENSE_API_SECRET")

	if uri == "" || apiKey == "" || apiSecret == "" {
		return
	}

	opnOptions := api.Options{
		Uri:           uri,
		APIKey:        apiKey,
		APISecret:     apiSecret,
		AllowInsecure: true,
	}
	client := api.NewClient(opnOptions)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *testOPNsenseProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewFirewallAliasResource,
		NewRouteResource,
		NewUnboundHostOverrideResource,
		NewInterfacesVlanResource,
		NewIpsecPskResource,
		NewIpsecConnectionResource,
		NewIpsecVtiResource,
		NewIpsecAuthLocalResource,
		NewIpsecAuthRemoteResource,
		NewIpsecChildResource,
	}
}

func (p *testOPNsenseProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewFirewallAliasDataSource,
		NewInterfaceDataSource,
		NewInterfaceAllDataSource,
		NewRouteDataSource,
	}
}

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"opnsense": providerserver.NewProtocol6WithError(testProvider()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("OPNSENSE_URI"); v == "" {
		t.Fatal("OPNSENSE_URI must be set for acceptance tests")
	}
	if v := os.Getenv("OPNSENSE_API_KEY"); v == "" {
		t.Fatal("OPNSENSE_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("OPNSENSE_API_SECRET"); v == "" {
		t.Fatal("OPNSENSE_API_SECRET must be set for acceptance tests")
	}
}
