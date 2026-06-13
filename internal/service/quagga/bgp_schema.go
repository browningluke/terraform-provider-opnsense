package quagga

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/quagga"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// bgpResourceModel describes the resource data model.
// This is a SINGLETON resource — it manages existing upstream configuration
// that cannot be created or destroyed via Terraform.
type bgpResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	Enabled            types.Bool   `tfsdk:"enabled"`
	ASNumber           types.String `tfsdk:"as_number"`
	Distance           types.Int64  `tfsdk:"distance"`
	RouterID           types.String `tfsdk:"router_id"`
	Graceful           types.Bool   `tfsdk:"graceful"`
	NetworkImportCheck types.Bool   `tfsdk:"network_import_check"`
	EnforceFirstAS     types.Bool   `tfsdk:"enforce_first_as"`
	LogNeighborChanges types.Bool   `tfsdk:"log_neighbor_changes"`
	Networks           types.Set    `tfsdk:"networks"`
	BestPath           types.Set    `tfsdk:"best_path"`
	MaximumPaths       types.Int64  `tfsdk:"maximum_paths"`
	MaximumPathsIBGP   types.Int64  `tfsdk:"maximum_paths_ibgp"`
}

func bgpResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages Quagga BGP (Border Gateway Protocol) global settings. This is a singleton resource that manages existing upstream configuration.\n\n" +
			"**Important:** This resource must be imported before it can be managed:\n" +
			"```bash\n" +
			"terraform import opnsense_quagga_bgp.bgp quagga_bgp\n" +
			"```\n\n" +
			"After importing, you can manage the configuration with `terraform apply`. " +
			"Running `terraform destroy` will remove the resource from state but will NOT modify the upstream configuration.",

		Version: 1,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Always set to `quagga_bgp`. Use this value when importing: `terraform import opnsense_quagga_bgp.bgp quagga_bgp`",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "When enabled, the BGP routing daemon is active. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"as_number": schema.StringAttribute{
				MarkdownDescription: "The local Autonomous System (AS) number for this BGP instance (e.g. `\"65000\"`). Defaults to `\"65551\"`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"distance": schema.Int64Attribute{
				MarkdownDescription: "Administrative distance applied to BGP routes. Use `-1` to leave unset. Valid range is 1–255.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"router_id": schema.StringAttribute{
				MarkdownDescription: "BGP router ID in IPv4 dotted-decimal notation (e.g. `\"10.0.0.1\"`). Leave empty to use the highest interface IP address.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"graceful": schema.BoolAttribute{
				MarkdownDescription: "When enabled, BGP peers are notified of a planned restart and routes are preserved during the restart period. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"network_import_check": schema.BoolAttribute{
				MarkdownDescription: "When enabled, verifies that the networks configured via `networks` are present in the routing table before advertising them. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"enforce_first_as": schema.BoolAttribute{
				MarkdownDescription: "When enabled, requires that the first AS in the AS_PATH of received BGP updates matches the peer AS number. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"log_neighbor_changes": schema.BoolAttribute{
				MarkdownDescription: "When enabled, logs a message whenever a BGP neighbor transitions between established and idle states. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"networks": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of network prefixes (CIDR notation) to originate and advertise via BGP.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"best_path": schema.SetAttribute{
				ElementType: types.StringType,
				MarkdownDescription: "BGP best-path selection options to enable. Valid values: " +
					"`\"as-path confed\"`, `\"as-path multipath-relax\"`, `\"compare-routerid\"`, " +
					"`\"peer-type multipath-relax\"`, `\"aigp\"`, `\"med missing-as-worst\"`.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf(
							"as-path confed",
							"as-path multipath-relax",
							"compare-routerid",
							"peer-type multipath-relax",
							"aigp",
							"med missing-as-worst",
						),
					),
				},
			},
			"maximum_paths": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of equal-cost EBGP paths to install in the routing table for multipath load balancing. Use `-1` to leave unset.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"maximum_paths_ibgp": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of equal-cost IBGP paths to install in the routing table for multipath load balancing. Use `-1` to leave unset.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func convertBGPSchemaToStruct(d *bgpResourceModel) (*quagga.QuaggaBGP, error) {
	return &quagga.QuaggaBGP{
		Enabled:            tools.BoolToString(d.Enabled.ValueBool()),
		ASNumber:           d.ASNumber.ValueString(),
		Distance:           tools.Int64ToStringNegative(d.Distance.ValueInt64()),
		RouterID:           d.RouterID.ValueString(),
		Graceful:           tools.BoolToString(d.Graceful.ValueBool()),
		NetworkImportCheck: tools.BoolToString(d.NetworkImportCheck.ValueBool()),
		EnforceFirstAS:     tools.BoolToString(d.EnforceFirstAS.ValueBool()),
		LogNeighborChanges: tools.BoolToString(d.LogNeighborChanges.ValueBool()),
		Networks:           api.SelectedMapList(tools.SetToStringSlice(d.Networks)),
		BestPath:           api.SelectedMapList(tools.SetToStringSlice(d.BestPath)),
		MaximumPaths:       tools.Int64ToStringNegative(d.MaximumPaths.ValueInt64()),
		MaximumPathsIBGP:   tools.Int64ToStringNegative(d.MaximumPathsIBGP.ValueInt64()),
	}, nil
}

func convertBGPStructToSchema(d *quagga.QuaggaBGP) (*bgpResourceModel, error) {
	return &bgpResourceModel{
		Enabled:            types.BoolValue(tools.StringToBool(d.Enabled)),
		ASNumber:           types.StringValue(d.ASNumber),
		Distance:           types.Int64Value(tools.StringToInt64(d.Distance)),
		RouterID:           types.StringValue(d.RouterID),
		Graceful:           types.BoolValue(tools.StringToBool(d.Graceful)),
		NetworkImportCheck: types.BoolValue(tools.StringToBool(d.NetworkImportCheck)),
		EnforceFirstAS:     types.BoolValue(tools.StringToBool(d.EnforceFirstAS)),
		LogNeighborChanges: types.BoolValue(tools.StringToBool(d.LogNeighborChanges)),
		Networks:           tools.StringSliceToSet([]string(d.Networks)),
		BestPath:           tools.StringSliceToSet([]string(d.BestPath)),
		MaximumPaths:       types.Int64Value(tools.StringToInt64(d.MaximumPaths)),
		MaximumPathsIBGP:   types.Int64Value(tools.StringToInt64(d.MaximumPathsIBGP)),
	}, nil
}
