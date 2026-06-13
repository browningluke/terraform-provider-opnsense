package quagga

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/quagga"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ospfResourceModel describes the resource data model.
// This is a SINGLETON resource — it manages existing upstream configuration
// that cannot be created or destroyed via Terraform.
type ospfResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	Enabled             types.Bool   `tfsdk:"enabled"`
	CARPDemote          types.Bool   `tfsdk:"carp_demote"`
	RouterID            types.String `tfsdk:"router_id"`
	CostReference       types.Int64  `tfsdk:"cost_reference"`
	LogAdjacencyChanges types.Bool   `tfsdk:"log_adjacency_changes"`
	Originate           types.Bool   `tfsdk:"originate"`
	OriginateAlways     types.Bool   `tfsdk:"originate_always"`
	OriginateMetric     types.Int64  `tfsdk:"originate_metric"`
	PassiveInterfaces   types.Set    `tfsdk:"passive_interfaces"`
}

func ospfResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages Quagga OSPF (Open Shortest Path First) global settings. This is a singleton resource that manages existing upstream configuration.\n\n" +
			"**Important:** This resource must be imported before it can be managed:\n" +
			"```bash\n" +
			"terraform import opnsense_quagga_ospf.ospf quagga_ospf\n" +
			"```\n\n" +
			"After importing, you can manage the configuration with `terraform apply`. " +
			"Running `terraform destroy` will remove the resource from state but will NOT modify the upstream configuration.",

		Version: 1,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Always set to `quagga_ospf`. Use this value when importing: `terraform import opnsense_quagga_ospf.ospf quagga_ospf`",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable the OSPF routing daemon. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"carp_demote": schema.BoolAttribute{
				MarkdownDescription: "Increase the CARP demotion counter while OSPF is not fully converged, causing CARP to prefer the peer until OSPF is ready. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"router_id": schema.StringAttribute{
				MarkdownDescription: "OSPF router ID in IPv4 dotted-decimal notation (e.g. `\"10.0.0.1\"`). Leave empty to use the highest interface IP address.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"cost_reference": schema.Int64Attribute{
				MarkdownDescription: "Reference bandwidth in Mbps used to calculate interface cost. Interfaces with a speed of `cost_reference` Mbps receive a cost of 1. Use `-1` to leave unset.",
				Optional:            true,
				Computed:            true,
			},
			"log_adjacency_changes": schema.BoolAttribute{
				MarkdownDescription: "Log a message when an OSPF adjacency changes state. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"originate": schema.BoolAttribute{
				MarkdownDescription: "Originate and advertise an OSPF default route (type-5 LSA) into the OSPF domain. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"originate_always": schema.BoolAttribute{
				MarkdownDescription: "Always originate the default route even when no default route exists in the routing table. Requires `originate` to be `true`. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"originate_metric": schema.Int64Attribute{
				MarkdownDescription: "Metric value assigned to the originated default route. Use `-1` to leave unset.",
				Optional:            true,
				Computed:            true,
			},
			"passive_interfaces": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Set of interface names to make passive. Passive interfaces do not send or receive OSPF Hello packets but still advertise their connected networks.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func convertOSPFSchemaToStruct(d *ospfResourceModel) (*quagga.QuaggaOSPF, error) {
	return &quagga.QuaggaOSPF{
		Enabled:             tools.BoolToString(d.Enabled.ValueBool()),
		CARPDemote:          tools.BoolToString(d.CARPDemote.ValueBool()),
		RouterID:            d.RouterID.ValueString(),
		CostReference:       tools.Int64ToStringNegative(d.CostReference.ValueInt64()),
		LogAdjacencyChanges: tools.BoolToString(d.LogAdjacencyChanges.ValueBool()),
		Originate:           tools.BoolToString(d.Originate.ValueBool()),
		OriginateAlways:     tools.BoolToString(d.OriginateAlways.ValueBool()),
		OriginateMetric:     tools.Int64ToStringNegative(d.OriginateMetric.ValueInt64()),
		PassiveInterfaces:   api.SelectedMapList(tools.SetToStringSlice(d.PassiveInterfaces)),
	}, nil
}

func convertOSPFStructToSchema(d *quagga.QuaggaOSPF) (*ospfResourceModel, error) {
	return &ospfResourceModel{
		Enabled:             types.BoolValue(tools.StringToBool(d.Enabled)),
		CARPDemote:          types.BoolValue(tools.StringToBool(d.CARPDemote)),
		RouterID:            types.StringValue(d.RouterID),
		CostReference:       types.Int64Value(tools.StringToInt64(d.CostReference)),
		LogAdjacencyChanges: types.BoolValue(tools.StringToBool(d.LogAdjacencyChanges)),
		Originate:           types.BoolValue(tools.StringToBool(d.Originate)),
		OriginateAlways:     types.BoolValue(tools.StringToBool(d.OriginateAlways)),
		OriginateMetric:     types.Int64Value(tools.StringToInt64(d.OriginateMetric)),
		PassiveInterfaces:   tools.StringSliceToSet([]string(d.PassiveInterfaces)),
	}, nil
}
