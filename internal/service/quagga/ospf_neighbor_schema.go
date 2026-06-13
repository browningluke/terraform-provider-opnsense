package quagga

import (
	"github.com/browningluke/opnsense-go/pkg/quagga"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ospfNeighborResourceModel describes the resource data model.
type ospfNeighborResourceModel struct {
	Enabled      types.Bool   `tfsdk:"enabled"`
	Description  types.String `tfsdk:"description"`
	Address      types.String `tfsdk:"address"`
	PollInterval types.Int64  `tfsdk:"poll_interval"`
	Priority     types.Int64  `tfsdk:"priority"`

	Id types.String `tfsdk:"id"`
}

func ospfNeighborResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Configure OSPF neighbors.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this OSPF neighbor. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "An optional description for this neighbor. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"address": schema.StringAttribute{
				MarkdownDescription: "The neighbor IP address.",
				Required:            true,
			},
			"poll_interval": schema.Int64Attribute{
				MarkdownDescription: "Poll interval in seconds. Use `-1` for unset. Defaults to `-1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Neighbor priority. Use `-1` for unset. Defaults to `-1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the OSPF neighbor.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func ospfNeighborDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure OSPF neighbors.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable this OSPF neighbor.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "An optional description for this neighbor.",
				Computed:            true,
			},
			"address": dschema.StringAttribute{
				MarkdownDescription: "The neighbor IP address.",
				Computed:            true,
			},
			"poll_interval": dschema.Int64Attribute{
				MarkdownDescription: "Poll interval in seconds.",
				Computed:            true,
			},
			"priority": dschema.Int64Attribute{
				MarkdownDescription: "Neighbor priority.",
				Computed:            true,
			},
		},
	}
}

func convertOSPFNeighborSchemaToStruct(d *ospfNeighborResourceModel) (*quagga.OSPFNeighbor, error) {
	return &quagga.OSPFNeighbor{
		Enabled:      tools.BoolToString(d.Enabled.ValueBool()),
		Description:  d.Description.ValueString(),
		Address:      d.Address.ValueString(),
		PollInterval: tools.Int64ToStringNegative(d.PollInterval.ValueInt64()),
		Priority:     tools.Int64ToStringNegative(d.Priority.ValueInt64()),
	}, nil
}

func convertOSPFNeighborStructToSchema(d *quagga.OSPFNeighbor) (*ospfNeighborResourceModel, error) {
	return &ospfNeighborResourceModel{
		Enabled:      types.BoolValue(tools.StringToBool(d.Enabled)),
		Description:  types.StringValue(d.Description),
		Address:      types.StringValue(d.Address),
		PollInterval: types.Int64Value(tools.StringToInt64(d.PollInterval)),
		Priority:     types.Int64Value(tools.StringToInt64(d.Priority)),
	}, nil
}
