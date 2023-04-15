package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-opnsense/internal/opnsense"
)

// Ensure OPNsenseProvider satisfies various provider interfaces.
var _ provider.Provider = &OPNsenseProvider{}

// OPNsenseProvider defines the provider implementation.
type OPNsenseProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// OPNsenseProviderModel describes the provider data model.
type OPNsenseProviderModel struct {
	Uri       types.String `tfsdk:"uri"`
	APIKey    types.String `tfsdk:"api_key"`
	APISecret types.String `tfsdk:"api_secret"`
}

func (p *OPNsenseProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "opnsense"
	resp.Version = p.version
}

func (p *OPNsenseProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"uri": schema.StringAttribute{
				MarkdownDescription: "The URI to an OPNsense host. Alternatively, can be configured using the `OPNSENSE_URI` environment variable.",
				Required:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The API key for a user. Alternatively, can be configured using the `OPNSENSE_API_KEY` environment variable.",
				Required:            true,
			},
			"api_secret": schema.StringAttribute{
				MarkdownDescription: "The API secret for a user. Alternatively, can be configured using the `OPNSENSE_API_SECRET` environment variable.",
				Required:            true,
			},
		},
	}
}

func (p *OPNsenseProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data OPNsenseProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	opnOptions := opnsense.Options{
		Uri:       data.Uri.ValueString(),
		APIKey:    data.APIKey.ValueString(),
		APISecret: data.APISecret.ValueString(),
	}

	client := opnsense.NewClient(opnOptions)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *OPNsenseProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *OPNsenseProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OPNsenseProvider{
			version: version,
		}
	}
}
