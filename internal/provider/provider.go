package provider

import (
	"context"
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-opnsense/internal/service"
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
	Uri           types.String `tfsdk:"uri"`
	APIKey        types.String `tfsdk:"api_key"`
	APISecret     types.String `tfsdk:"api_secret"`
	AllowInsecure types.Bool   `tfsdk:"allow_insecure"`
	MaxBackoff    types.Int64  `tfsdk:"max_backoff"`
	MinBackoff    types.Int64  `tfsdk:"min_backoff"`
	MaxRetries    types.Int64  `tfsdk:"retries"`
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
			"allow_insecure": schema.BoolAttribute{
				MarkdownDescription: "Allow insecure TLS connections. Alternatively, can be configured using the `OPNSENSE_ALLOW_INSECURE` environment variable. Defaults to `false`.",
				Optional:            true,
			},
			"max_backoff": schema.Int64Attribute{
				MarkdownDescription: "Maximum backoff period in seconds after failed API calls. Alternatively, can be configured using the `OPNSENSE_MAX_BACKOFF` environment variable.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"min_backoff": schema.Int64Attribute{
				MarkdownDescription: "Minimum backoff period in seconds after failed API calls. Alternatively, can be configured using the `OPNSENSE_MIN_BACKOFF` environment variable.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"retries": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of retries to perform when an API request fails. Alternatively, can be configured using the `OPNSENSE_RETRIES` environment variable.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 2147483647), // Since we convert the int64 to an int(32), set an upper bound.
				},
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

	opnOptions := api.Options{
		Uri:           data.Uri.ValueString(),
		APIKey:        data.APIKey.ValueString(),
		APISecret:     data.APISecret.ValueString(),
		AllowInsecure: data.AllowInsecure.ValueBool(),
		MaxBackoff:    data.MaxBackoff.ValueInt64(),
		MinBackoff:    data.MinBackoff.ValueInt64(),
		MaxRetries:    data.MaxRetries.ValueInt64(),
	}

	client := api.NewClient(opnOptions)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *OPNsenseProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Interfaces
		service.NewInterfacesVlanResource,
		// Routes
		service.NewRouteResource,
		// Unbound
		service.NewUnboundHostOverrideResource,
		service.NewUnboundHostAliasResource,
		service.NewUnboundDomainOverrideResource,
		service.NewUnboundForwardResource,
		// Firewall
		service.NewFirewallFilterResource,
		service.NewFirewallNATResource,
		service.NewFirewallCategoryResource,
	}
}

func (p *OPNsenseProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Interfaces
		service.NewInterfacesVlanDataSource,
		// Routes
		service.NewRouteDataSource,
		// Unbound
		service.NewUnboundHostOverrideDataSource,
		service.NewUnboundHostAliasDataSource,
		service.NewUnboundDomainOverrideDataSource,
		service.NewUnboundForwardDataSource,
		// Firewall
		service.NewFirewallFilterDataSource,
		service.NewFirewallNATDataSource,
		service.NewFirewallCategoryDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OPNsenseProvider{
			version: version,
		}
	}
}
