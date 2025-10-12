package provider

import (
	"context"
	"os"
	"strconv"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/terraform-provider-opnsense/internal/service/diagnostics"
	"github.com/browningluke/terraform-provider-opnsense/internal/service/firewall"
	"github.com/browningluke/terraform-provider-opnsense/internal/service/interfaces"
	"github.com/browningluke/terraform-provider-opnsense/internal/service/ipsec"
	"github.com/browningluke/terraform-provider-opnsense/internal/service/routes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure OPNsenseProvider satisfies various provider interfaces.
var _ provider.Provider = &opnsenseProvider{}

// OPNsenseProvider defines the provider implementation.
type opnsenseProvider struct {
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

func (p *opnsenseProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "opnsense"
	resp.Version = p.version
}

func (p *opnsenseProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"uri": schema.StringAttribute{
				// Required, but computed from environment variable if possible
				MarkdownDescription: "The URI to an OPNsense host. Alternatively, can be configured using the `OPNSENSE_URI` environment variable.",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				// Required, but computed from environment variable if possible
				MarkdownDescription: "The API key for a user. Alternatively, can be configured using the `OPNSENSE_API_KEY` environment variable.",
				Optional:            true,
			},
			"api_secret": schema.StringAttribute{
				// Required, but computed from environment variable if possible
				MarkdownDescription: "The API secret for a user. Alternatively, can be configured using the `OPNSENSE_API_SECRET` environment variable.",
				Optional:            true,
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

func (p *opnsenseProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Get data from config
	var data OPNsenseProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If value is set, it must be known

	if data.Uri.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("uri"),
			"Unknown OPNsense API URI",
			"The provider cannot create the OPNsense API client as there is an unknown configuration value for the OPNsense API uri. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OPNSENSE_URI environment variable.",
		)
	}

	if data.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown OPNsense API Key",
			"The provider cannot create the OPNsense API client as there is an unknown configuration value for the OPNsense API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OPNSENSE_API_KEY environment variable.",
		)
	}

	if data.APISecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_secret"),
			"Unknown OPNsense API Secret",
			"The provider cannot create the OPNsense API client as there is an unknown configuration value for the OPNsense API secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OPNSENSE_API_SECRET environment variable.",
		)
	}

	if data.AllowInsecure.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("allow_insecure"),
			"Unknown OPNsense API Value: allow_insecure",
			"The provider cannot create the OPNsense API client as there is an unknown configuration value for allow_insecure. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OPNSENSE_ALLOW_INSECURE environment variable.",
		)
	}

	if data.MaxBackoff.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("max_backoff"),
			"Unknown OPNsense API Value: max_backoff",
			"The provider cannot create the OPNsense API client as there is an unknown configuration value for max_backoff. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OPNSENSE_MAX_BACKOFF environment variable.",
		)
	}

	if data.MinBackoff.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("min_backoff"),
			"Unknown OPNsense API Value: min_backoff",
			"The provider cannot create the OPNsense API client as there is an unknown configuration value for min_backoff. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OPNSENSE_MIN_BACKOFF environment variable.",
		)
	}

	if data.MaxRetries.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("retries"),
			"Unknown OPNsense API Value: retries",
			"The provider cannot create the OPNsense API client as there is an unknown configuration value for retries. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OPNSENSE_RETRIES environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Attempt to load values from environment variables

	uri := os.Getenv("OPNSENSE_URI")
	if !data.Uri.IsNull() {
		uri = data.Uri.ValueString()
	}

	apiKey := os.Getenv("OPNSENSE_API_KEY")
	if !data.APIKey.IsNull() {
		apiKey = data.APIKey.ValueString()
	}

	apiSecret := os.Getenv("OPNSENSE_API_SECRET")
	if !data.APISecret.IsNull() {
		apiSecret = data.APISecret.ValueString()
	}

	allowInsecureStr := os.Getenv("OPNSENSE_ALLOW_INSECURE")
	allowInsecure, err := strconv.ParseBool(allowInsecureStr)
	if err != nil {
		// Set to default (false) if string is unparsable
		allowInsecure = false
	}
	if !data.AllowInsecure.IsNull() {
		allowInsecure = data.AllowInsecure.ValueBool()
	}

	maxBackoffStr := os.Getenv("OPNSENSE_MAX_BACKOFF")
	maxBackoff, err := strconv.ParseInt(maxBackoffStr, 10, 64)
	if err != nil {
		// Set to 0 to use client default downstream if string is unparsable
		maxBackoff = 0
	}
	if !data.MaxBackoff.IsNull() {
		maxBackoff = data.MaxBackoff.ValueInt64()
	}

	minBackoffStr := os.Getenv("OPNSENSE_MIN_BACKOFF")
	minBackoff, err := strconv.ParseInt(minBackoffStr, 10, 64)
	if err != nil {
		// Set to 0 to use client default downstream if string is unparsable
		minBackoff = 0
	}
	if !data.MinBackoff.IsNull() {
		maxBackoff = data.MinBackoff.ValueInt64()
	}

	retriesStr := os.Getenv("OPNSENSE_RETRIES")
	retries, err := strconv.ParseInt(retriesStr, 10, 64)
	if err != nil {
		// Set to 0 to use client default downstream if string is unparsable
		retries = 0
	}
	if !data.MaxRetries.IsNull() {
		retries = data.MaxRetries.ValueInt64()
	}

	// Ensure expected variables are not empty

	if uri == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("uri"),
			"Missing OPNsense API URI",
			"The provider cannot create the OPNsense API client as there is a missing or empty value for the OPNsense API uri. "+
				"Set the host value in the configuration or use the OPNSENSE_URI environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing OPNsense API Key",
			"The provider cannot create the OPNsense API client as there is a missing or empty value for the OPNsense API key. "+
				"Set the host value in the configuration or use the OPNSENSE_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_secret"),
			"Missing OPNsense API Secret",
			"The provider cannot create the OPNsense API client as there is a missing or empty value for the OPNsense API secret. "+
				"Set the host value in the configuration or use the OPNSENSE_API_SECRET environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	// Create the OPNsense client
	opnOptions := api.Options{
		Uri:           uri,
		APIKey:        apiKey,
		APISecret:     apiSecret,
		AllowInsecure: allowInsecure,
		MaxBackoff:    maxBackoff,
		MinBackoff:    minBackoff,
		MaxRetries:    retries,
	}
	client := api.NewClient(opnOptions)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *opnsenseProvider) Resources(ctx context.Context) []func() resource.Resource {
	controllers := [][]func() resource.Resource{
		diagnostics.Resources(ctx),
		firewall.Resources(ctx),
		interfaces.Resources(ctx),
		ipsec.Resources(ctx),
		routes.Resources(ctx),
	}

	var resources []func() resource.Resource
	for _, s := range controllers {
		resources = append(resources, s...)
	}
	return resources
}

func (p *opnsenseProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	controllers := [][]func() datasource.DataSource{
		diagnostics.DataSources(ctx),
		firewall.DataSources(ctx),
		interfaces.DataSources(ctx),
		ipsec.DataSources(ctx),
		routes.DataSources(ctx),
	}

	var dataSources []func() datasource.DataSource
	for _, s := range controllers {
		dataSources = append(dataSources, s...)
	}
	return dataSources
}

func NewProvider(ctx context.Context) (provider.Provider, error) {
	return &opnsenseProvider{}, nil
}
