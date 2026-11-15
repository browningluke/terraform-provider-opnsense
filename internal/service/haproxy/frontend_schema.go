package haproxy

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/haproxy"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/browningluke/terraform-provider-opnsense/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type frontendResourceModel struct {
	Enabled        types.Bool   `tfsdk:"enabled"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Bind           types.Set    `tfsdk:"bind"`
	Mode           types.String `tfsdk:"mode"`
	DefaultBackend types.String `tfsdk:"default_backend"`

	Id types.String `tfsdk:"id"`
}

func frontendResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "A HAProxy frontend (OPNsense: Public Service) without TLS. " +
			"Settings not listed here keep their OPNsense defaults on create and their current values on update. " +
			"This resource never reads or writes the frontend's rule list; manage that with `opnsense_haproxy_frontend_actions`, " +
			"which in turn touches nothing else, so both can point at the same frontend. Requires the `os-haproxy` plugin.",
		Version: 1,

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this frontend. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name to identify this frontend. Letters, digits, `.`, `_` and `-` only.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description for this frontend. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"bind": schema.SetAttribute{
				MarkdownDescription: "Listen addresses as `address:port`, e.g. `0.0.0.0:80` or `[::]:80`.",
				ElementType:         types.StringType,
				Required:            true,
				Validators:          []validator.Set{setvalidator.SizeAtLeast(1)},
			},
			"mode": schema.StringAttribute{
				MarkdownDescription: "`http` (layer 7) or `tcp` (layer 4). The `ssl` mode is not supported because its certificate settings are not modelled. Defaults to `http`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("http"),
				Validators:          []validator.String{stringvalidator.OneOf("http", "tcp")},
			},
			"default_backend": schema.StringAttribute{
				MarkdownDescription: "UUID of the backend pool (`opnsense_haproxy_backend`) that receives requests no rule matched. Empty means none. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Validators:          []validator.String{validators.IsUUIDv4()},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the frontend.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func convertFrontendSchemaToStruct(d *frontendResourceModel) (*haproxy.Frontend, error) {
	return &haproxy.Frontend{
		Enabled:        tools.BoolToString(d.Enabled.ValueBool()),
		Name:           d.Name.ValueString(),
		Description:    d.Description.ValueString(),
		Bind:           tools.SetToStringSlice(d.Bind),
		Mode:           api.SelectedMap(d.Mode.ValueString()),
		DefaultBackend: api.SelectedMap(d.DefaultBackend.ValueString()),
	}, nil
}

func convertFrontendStructToSchema(d *haproxy.Frontend) (*frontendResourceModel, error) {
	return &frontendResourceModel{
		Enabled:        types.BoolValue(tools.StringToBool(d.Enabled)),
		Name:           types.StringValue(d.Name),
		Description:    types.StringValue(d.Description),
		Bind:           tools.StringSliceToSet(d.Bind),
		Mode:           types.StringValue(d.Mode.String()),
		DefaultBackend: types.StringValue(d.DefaultBackend.String()),
	}, nil
}
