package quagga

import (
	"github.com/browningluke/opnsense-go/pkg/quagga"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ospf6ResourceModel describes the resource data model.
// This is a SINGLETON resource — it manages existing upstream configuration
// that cannot be created or destroyed via Terraform.
type ospf6ResourceModel struct {
	Id              types.String `tfsdk:"id"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	CARPDemote      types.Bool   `tfsdk:"carp_demote"`
	RouterID        types.String `tfsdk:"router_id"`
	Originate       types.Bool   `tfsdk:"originate"`
	OriginateAlways types.Bool   `tfsdk:"originate_always"`
	OriginateMetric types.Int64  `tfsdk:"originate_metric"`
}

func ospf6ResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages Quagga OSPFv3 (Open Shortest Path First version 3 for IPv6) global settings. This is a singleton resource that manages existing upstream configuration.\n\n" +
			"**Important:** This resource must be imported before it can be managed:\n" +
			"```bash\n" +
			"terraform import opnsense_quagga_ospf6.ospf6 quagga_ospf6\n" +
			"```\n\n" +
			"After importing, you can manage the configuration with `terraform apply`. " +
			"Running `terraform destroy` will remove the resource from state but will NOT modify the upstream configuration.",

		Version: 1,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Always set to `quagga_ospf6`. Use this value when importing: `terraform import opnsense_quagga_ospf6.ospf6 quagga_ospf6`",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable the OSPFv3 routing daemon. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"carp_demote": schema.BoolAttribute{
				MarkdownDescription: "Increase the CARP demotion counter while OSPFv3 is not fully converged, causing CARP to prefer the peer until OSPFv3 is ready. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"router_id": schema.StringAttribute{
				MarkdownDescription: "OSPFv3 router ID in IPv4 dotted-decimal notation (e.g. `\"10.0.0.1\"`). Leave empty to use the highest interface IP address.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"originate": schema.BoolAttribute{
				MarkdownDescription: "Originate and advertise an OSPFv3 default route into the OSPFv3 domain. Defaults to `false`.",
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
		},
	}
}

func convertOSPF6SchemaToStruct(d *ospf6ResourceModel) (*quagga.QuaggaOSPF6, error) {
	return &quagga.QuaggaOSPF6{
		Enabled:         tools.BoolToString(d.Enabled.ValueBool()),
		CARPDemote:      tools.BoolToString(d.CARPDemote.ValueBool()),
		RouterID:        d.RouterID.ValueString(),
		Originate:       tools.BoolToString(d.Originate.ValueBool()),
		OriginateAlways: tools.BoolToString(d.OriginateAlways.ValueBool()),
		OriginateMetric: tools.Int64ToStringNegative(d.OriginateMetric.ValueInt64()),
	}, nil
}

func convertOSPF6StructToSchema(d *quagga.QuaggaOSPF6) (*ospf6ResourceModel, error) {
	return &ospf6ResourceModel{
		Enabled:         types.BoolValue(tools.StringToBool(d.Enabled)),
		CARPDemote:      types.BoolValue(tools.StringToBool(d.CARPDemote)),
		RouterID:        types.StringValue(d.RouterID),
		Originate:       types.BoolValue(tools.StringToBool(d.Originate)),
		OriginateAlways: types.BoolValue(tools.StringToBool(d.OriginateAlways)),
		OriginateMetric: types.Int64Value(tools.StringToInt64(d.OriginateMetric)),
	}, nil
}
