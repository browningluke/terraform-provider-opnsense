package haproxy

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/haproxy"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type serverResourceModel struct {
	Enabled     types.Bool   `tfsdk:"enabled"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Address     types.String `tfsdk:"address"`
	Port        types.Int64  `tfsdk:"port"`
	Mode        types.String `tfsdk:"mode"`
	SSL         types.Bool   `tfsdk:"ssl"`
	SSLVerify   types.Bool   `tfsdk:"ssl_verify"`

	Id types.String `tfsdk:"id"`
}

func serverResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "A HAProxy real server: one target host that backend pools send traffic to. Requires the `os-haproxy` plugin.",
		Version:             1,

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this server. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name to identify this server. Letters, digits, `.`, `_` and `-` only.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description for this server. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"address": schema.StringAttribute{
				MarkdownDescription: "FQDN or IP address of this server.",
				Required:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port of this server. `-1` leaves it empty, which makes HAProxy use the port the client connected to. Defaults to `-1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"mode": schema.StringAttribute{
				MarkdownDescription: "Server mode: `active`, `backup` (only used when all active servers are down) or `disabled`. Defaults to `active`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("active"),
				Validators:          []validator.String{stringvalidator.OneOf("active", "backup", "disabled")},
			},
			"ssl": schema.BoolAttribute{
				MarkdownDescription: "When enabled, HAProxy connects to this server over SSL/TLS. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"ssl_verify": schema.BoolAttribute{
				MarkdownDescription: "When enabled, the server certificate is verified when `ssl` is enabled. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the server.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func convertServerSchemaToStruct(d *serverResourceModel) (*haproxy.Server, error) {
	return &haproxy.Server{
		Enabled:     tools.BoolToString(d.Enabled.ValueBool()),
		Name:        d.Name.ValueString(),
		Description: d.Description.ValueString(),
		Address:     d.Address.ValueString(),
		Port:        tools.Int64ToStringNegative(d.Port.ValueInt64()),
		Mode:        api.SelectedMap(d.Mode.ValueString()),
		SSL:         tools.BoolToString(d.SSL.ValueBool()),
		SSLVerify:   tools.BoolToString(d.SSLVerify.ValueBool()),
	}, nil
}

func convertServerStructToSchema(d *haproxy.Server) (*serverResourceModel, error) {
	return &serverResourceModel{
		Enabled:     types.BoolValue(tools.StringToBool(d.Enabled)),
		Name:        types.StringValue(d.Name),
		Description: types.StringValue(d.Description),
		Address:     types.StringValue(d.Address),
		Port:        types.Int64Value(tools.StringToInt64(d.Port)),
		Mode:        types.StringValue(d.Mode.String()),
		SSL:         types.BoolValue(tools.StringToBool(d.SSL)),
		SSLVerify:   types.BoolValue(tools.StringToBool(d.SSLVerify)),
	}, nil
}
