package haproxy

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/haproxy"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type backendResourceModel struct {
	Enabled            types.Bool   `tfsdk:"enabled"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	Mode               types.String `tfsdk:"mode"`
	Algorithm          types.String `tfsdk:"algorithm"`
	LinkedServers      types.Set    `tfsdk:"linked_servers"`
	HealthCheckEnabled types.Bool   `tfsdk:"health_check_enabled"`

	Id types.String `tfsdk:"id"`
}

func backendResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "A HAProxy backend pool: a group of servers that receives the traffic a rule or frontend sends to it. Requires the `os-haproxy` plugin.",
		Version:             1,

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this backend pool. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name to identify this backend pool. Letters, digits, `.`, `_` and `-` only.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description for this backend pool. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"mode": schema.StringAttribute{
				MarkdownDescription: "Mode of the backend pool, `http` (layer 7) or `tcp` (layer 4). Defaults to `http`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("http"),
				Validators:          []validator.String{stringvalidator.OneOf("http", "tcp")},
			},
			"algorithm": schema.StringAttribute{
				MarkdownDescription: "Balancing algorithm: `source`, `roundrobin`, `static-rr`, `leastconn`, `uri` or `random`. Defaults to `source`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("source"),
				Validators:          []validator.String{stringvalidator.OneOf("source", "roundrobin", "static-rr", "leastconn", "uri", "random")},
			},
			"linked_servers": schema.SetAttribute{
				MarkdownDescription: "UUIDs of the servers in this pool. Defaults to `[]`.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"health_check_enabled": schema.BoolAttribute{
				MarkdownDescription: "When enabled, HAProxy health-checks the servers in this pool. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the backend pool.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func convertBackendSchemaToStruct(d *backendResourceModel) (*haproxy.Backend, error) {
	return &haproxy.Backend{
		Enabled:            tools.BoolToString(d.Enabled.ValueBool()),
		Name:               d.Name.ValueString(),
		Description:        d.Description.ValueString(),
		Mode:               api.SelectedMap(d.Mode.ValueString()),
		Algorithm:          api.SelectedMap(d.Algorithm.ValueString()),
		LinkedServers:      tools.SetToStringSlice(d.LinkedServers),
		HealthCheckEnabled: tools.BoolToString(d.HealthCheckEnabled.ValueBool()),
	}, nil
}

func convertBackendStructToSchema(d *haproxy.Backend) (*backendResourceModel, error) {
	return &backendResourceModel{
		Enabled:            types.BoolValue(tools.StringToBool(d.Enabled)),
		Name:               types.StringValue(d.Name),
		Description:        types.StringValue(d.Description),
		Mode:               types.StringValue(d.Mode.String()),
		Algorithm:          types.StringValue(d.Algorithm.String()),
		LinkedServers:      tools.StringSliceToSet(d.LinkedServers),
		HealthCheckEnabled: types.BoolValue(tools.StringToBool(d.HealthCheckEnabled)),
	}, nil
}
