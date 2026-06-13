package quagga

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/quagga"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
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

// ripResourceModel describes the resource data model.
// This is a SINGLETON resource — it manages existing upstream configuration
// that cannot be created or destroyed via Terraform.
type ripResourceModel struct {
	Id                types.String `tfsdk:"id"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	Version           types.String `tfsdk:"version"`
	Networks          types.Set    `tfsdk:"networks"`
	PassiveInterfaces types.Set    `tfsdk:"passive_interfaces"`
	Redistribute      types.Set    `tfsdk:"redistribute"`
	DefaultMetric     types.Int64  `tfsdk:"default_metric"`
}

func ripResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages Quagga RIP (Routing Information Protocol) settings. This is a singleton resource that manages existing upstream configuration.\n\n" +
			"**Important:** This resource must be imported before it can be managed:\n" +
			"```bash\n" +
			"terraform import opnsense_quagga_rip.rip quagga_rip\n" +
			"```\n\n" +
			"After importing, you can manage the configuration with `terraform apply`. " +
			"Running `terraform destroy` will remove the resource from state but will NOT modify the upstream configuration.",

		Version: 1,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Always set to `quagga_rip`. Use this value when importing: `terraform import opnsense_quagga_rip.rip quagga_rip`",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable the RIP routing daemon. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "RIP protocol version to use. `\"1\"` uses broadcast updates; `\"2\"` uses multicast updates and supports CIDR. Defaults to `\"2\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("2"),
				Validators: []validator.String{
					stringvalidator.OneOf("1", "2"),
				},
			},
			"networks": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Set of network prefixes (CIDR notation) to enable RIP on. RIP will send and receive updates on interfaces that belong to these networks.",
				Optional:            true,
				Computed:            true,
			},
			"passive_interfaces": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Set of interface names to make passive. Passive interfaces do not send RIP updates but still advertise their connected networks.",
				Optional:            true,
				Computed:            true,
			},
			"redistribute": schema.SetAttribute{
				ElementType: types.StringType,
				MarkdownDescription: "Set of routing protocols whose routes should be redistributed into RIP. " +
					"Valid values: `\"bgp\"`, `\"connected\"`, `\"kernel\"`, `\"ospf\"`, `\"static\"`.",
				Optional: true,
				Computed: true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf("bgp", "connected", "kernel", "ospf", "static"),
					),
				},
			},
			"default_metric": schema.Int64Attribute{
				MarkdownDescription: "Default metric assigned to redistributed routes. Use `-1` to leave unset.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func convertRIPSchemaToStruct(d *ripResourceModel) (*quagga.QuaggaRIP, error) {
	return &quagga.QuaggaRIP{
		Enabled:           tools.BoolToString(d.Enabled.ValueBool()),
		Version:           d.Version.ValueString(),
		Networks:          api.SelectedMapList(tools.SetToStringSlice(d.Networks)),
		PassiveInterfaces: api.SelectedMapList(tools.SetToStringSlice(d.PassiveInterfaces)),
		Redistribute:      api.SelectedMapList(tools.SetToStringSlice(d.Redistribute)),
		DefaultMetric:     tools.Int64ToStringNegative(d.DefaultMetric.ValueInt64()),
	}, nil
}

func convertRIPStructToSchema(d *quagga.QuaggaRIP) (*ripResourceModel, error) {
	return &ripResourceModel{
		Enabled:           types.BoolValue(tools.StringToBool(d.Enabled)),
		Version:           types.StringValue(d.Version),
		Networks:          tools.StringSliceToSet([]string(d.Networks)),
		PassiveInterfaces: tools.StringSliceToSet([]string(d.PassiveInterfaces)),
		Redistribute:      tools.StringSliceToSet([]string(d.Redistribute)),
		DefaultMetric:     types.Int64Value(tools.StringToInt64(d.DefaultMetric)),
	}, nil
}
