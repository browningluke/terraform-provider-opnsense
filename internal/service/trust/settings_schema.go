package trust

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/trust"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type settingsResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	StoreIntermediateCerts  types.Bool   `tfsdk:"store_intermediate_certs"`
	InstallCrls             types.Bool   `tfsdk:"install_crls"`
	FetchCrls               types.Bool   `tfsdk:"fetch_crls"`
	EnableLegacySect        types.Bool   `tfsdk:"enable_legacy_sect"`
	EnableConfigConstraints types.Bool   `tfsdk:"enable_config_constraints"`
	CipherString            types.Set    `tfsdk:"cipher_string"`
}

func settingsResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages Trust global settings (TLS cipher policy, CRL handling, OpenSSL legacy mode). This is a singleton resource that manages existing upstream configuration.\n\n" +
			"**Important:** This resource must be imported before it can be managed:\n" +
			"```bash\n" +
			"terraform import opnsense_trust_settings.settings trust_settings\n" +
			"```\n\n" +
			"After importing, you can manage the configuration with `terraform apply`. " +
			"Running `terraform destroy` will remove the resource from state but will NOT modify the upstream configuration.",

		Version: 1,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Always set to `trust_settings`. Use this value when importing: `terraform import opnsense_trust_settings.settings trust_settings`",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"store_intermediate_certs": schema.BoolAttribute{
				MarkdownDescription: "When enabled, intermediate CA certificates are stored in the system trust store. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"install_crls": schema.BoolAttribute{
				MarkdownDescription: "When enabled, fetched CRLs are automatically installed into the system trust store. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"fetch_crls": schema.BoolAttribute{
				MarkdownDescription: "When enabled, a cron job periodically fetches CRLs from distribution points embedded in certificates. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"enable_legacy_sect": schema.BoolAttribute{
				MarkdownDescription: "When enabled, the OpenSSL legacy provider section is active (enables older algorithms such as MD4, DES). Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"enable_config_constraints": schema.BoolAttribute{
				MarkdownDescription: "When enabled, OpenSSL policy constraints are enforced. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"cipher_string": schema.SetAttribute{
				MarkdownDescription: "Set of TLS cipher names to allow (OpenSSL cipher-suite identifiers, e.g. `TLS_AES_256_GCM_SHA384`). When empty, OPNsense applies its built-in safe default set.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func settingsDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Reads Trust global settings from the upstream system.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Always set to `trust_settings`.",
			},
			"store_intermediate_certs": dschema.BoolAttribute{
				MarkdownDescription: "Whether intermediate CA certificates are stored in the system trust store.",
				Computed:            true,
			},
			"install_crls": dschema.BoolAttribute{
				MarkdownDescription: "Whether fetched CRLs are automatically installed.",
				Computed:            true,
			},
			"fetch_crls": dschema.BoolAttribute{
				MarkdownDescription: "Whether the periodic CRL fetch cron job is enabled.",
				Computed:            true,
			},
			"enable_legacy_sect": dschema.BoolAttribute{
				MarkdownDescription: "Whether the OpenSSL legacy provider section is active.",
				Computed:            true,
			},
			"enable_config_constraints": dschema.BoolAttribute{
				MarkdownDescription: "Whether OpenSSL policy constraints are enforced.",
				Computed:            true,
			},
			"cipher_string": dschema.SetAttribute{
				MarkdownDescription: "Set of allowed TLS cipher names.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func convertSettingsSchemaToStruct(d *settingsResourceModel) (*trust.TrustSettings, error) {
	return &trust.TrustSettings{
		StoreIntermediateCerts:  tools.BoolToString(d.StoreIntermediateCerts.ValueBool()),
		InstallCrls:             tools.BoolToString(d.InstallCrls.ValueBool()),
		FetchCrls:               tools.BoolToString(d.FetchCrls.ValueBool()),
		EnableLegacySect:        tools.BoolToString(d.EnableLegacySect.ValueBool()),
		EnableConfigConstraints: tools.BoolToString(d.EnableConfigConstraints.ValueBool()),
		CipherString:            api.SelectedMapList(tools.SetToStringSlice(d.CipherString)),
	}, nil
}

func convertSettingsStructToSchema(d *trust.TrustSettings) (*settingsResourceModel, error) {
	return &settingsResourceModel{
		StoreIntermediateCerts:  types.BoolValue(tools.StringToBool(d.StoreIntermediateCerts)),
		InstallCrls:             types.BoolValue(tools.StringToBool(d.InstallCrls)),
		FetchCrls:               types.BoolValue(tools.StringToBool(d.FetchCrls)),
		EnableLegacySect:        types.BoolValue(tools.StringToBool(d.EnableLegacySect)),
		EnableConfigConstraints: types.BoolValue(tools.StringToBool(d.EnableConfigConstraints)),
		CipherString:            tools.StringSliceToSet(d.CipherString),
	}, nil
}
