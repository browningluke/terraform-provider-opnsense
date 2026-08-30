package interfaces

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/interfaces"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// bridgeResourceModel describes the resource data model.
type bridgeResourceModel struct {
	Members        types.Set    `tfsdk:"members"`
	LinkLocal      types.Bool   `tfsdk:"link_local"`
	EnableStp      types.Bool   `tfsdk:"enable_stp"`
	StpProto       types.String `tfsdk:"stp_proto"`
	StpInterfaces  types.Set    `tfsdk:"stp_interfaces"`
	MaxAge         types.Int64  `tfsdk:"max_age"`
	ForwardDelay   types.Int64  `tfsdk:"forward_delay"`
	HoldCount      types.Int64  `tfsdk:"hold_count"`
	MaxAddresses   types.Int64  `tfsdk:"max_addresses"`
	AddressTimeout types.Int64  `tfsdk:"address_timeout"`
	SpanInterface  types.String `tfsdk:"span_interface"`
	EdgePorts      types.Set    `tfsdk:"edge_ports"`
	AutoEdgePorts  types.Set    `tfsdk:"auto_edge_ports"`
	PtpPorts       types.Set    `tfsdk:"ptp_ports"`
	AutoPtpPorts   types.Set    `tfsdk:"auto_ptp_ports"`
	StaticPorts    types.Set    `tfsdk:"static_ports"`
	PrivatePorts   types.Set    `tfsdk:"private_ports"`
	Description    types.String `tfsdk:"description"`
	Device         types.String `tfsdk:"device"`

	Id types.String `tfsdk:"id"`
}

func bridgeResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Bridges connect two or more network interfaces into a single layer 2 segment.",

		Attributes: map[string]schema.Attribute{
			"members": schema.SetAttribute{
				MarkdownDescription: "Assigned interfaces to add as members of the bridge, e.g. `[\"opt1\", \"opt2\"]`. At least one member is required.",
				ElementType:         types.StringType,
				Required:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"link_local": schema.BoolAttribute{
				MarkdownDescription: "Enable link-local (IPv6) addresses on the bridge and its members. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"enable_stp": schema.BoolAttribute{
				MarkdownDescription: "Enable spanning tree options for this bridge. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"stp_proto": schema.StringAttribute{
				MarkdownDescription: "Spanning tree protocol to use. Available values: `rstp`, `stp`. Defaults to `rstp`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("rstp"),
				Validators: []validator.String{
					stringvalidator.OneOf("rstp", "stp"),
				},
			},
			"stp_interfaces": schema.SetAttribute{
				MarkdownDescription: "Interfaces to enable spanning tree on. Defaults to `[]`.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"max_age": schema.Int64Attribute{
				MarkdownDescription: "Time (in seconds) that a spanning tree configuration is valid. Must be between `6` and `40`.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(6, 40),
				},
			},
			"forward_delay": schema.Int64Attribute{
				MarkdownDescription: "Time (in seconds) that must pass before an interface begins forwarding packets when spanning tree is enabled. Must be between `4` and `30`.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(4, 30),
				},
			},
			"hold_count": schema.Int64Attribute{
				MarkdownDescription: "Number of packets transmitted before being rate limited when spanning tree is enabled. Must be between `1` and `10`.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 10),
				},
			},
			"max_addresses": schema.Int64Attribute{
				MarkdownDescription: "Maximum size of the bridge address cache. Must be at least `1`.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"address_timeout": schema.Int64Attribute{
				MarkdownDescription: "Timeout (in seconds) of bridge address cache entries. `0` disables expiry.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"span_interface": schema.StringAttribute{
				MarkdownDescription: "Interface to transmit a copy of every frame received by the bridge on. Set to `\"\"` for none. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"edge_ports": schema.SetAttribute{
				MarkdownDescription: "Interfaces to set as edge ports (always connected to an end of a LAN, never to another bridge). Defaults to `[]`.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"auto_edge_ports": schema.SetAttribute{
				MarkdownDescription: "Interfaces to automatically detect edge status on. Defaults to `[]`.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"ptp_ports": schema.SetAttribute{
				MarkdownDescription: "Interfaces to set as point-to-point links. Defaults to `[]`.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"auto_ptp_ports": schema.SetAttribute{
				MarkdownDescription: "Interfaces to automatically detect point-to-point status on. Defaults to `[]`.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"static_ports": schema.SetAttribute{
				MarkdownDescription: "Interfaces to mark as static (no address learning). Defaults to `[]`.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"private_ports": schema.SetAttribute{
				MarkdownDescription: "Interfaces to mark as private (do not forward traffic to other private ports). Defaults to `[]`.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Optional:            true,
			},
			"device": schema.StringAttribute{
				MarkdownDescription: "Custom bridge name, e.g. `bridge0`. Set to `\"\"` to generate a device name. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the bridge.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func bridgeDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Bridges connect two or more network interfaces into a single layer 2 segment.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"members": dschema.SetAttribute{
				MarkdownDescription: "Assigned interfaces that are members of the bridge.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"link_local": dschema.BoolAttribute{
				MarkdownDescription: "Whether link-local (IPv6) addresses are enabled on the bridge and its members.",
				Computed:            true,
			},
			"enable_stp": dschema.BoolAttribute{
				MarkdownDescription: "Whether spanning tree options are enabled for this bridge.",
				Computed:            true,
			},
			"stp_proto": dschema.StringAttribute{
				MarkdownDescription: "Spanning tree protocol in use.",
				Computed:            true,
			},
			"stp_interfaces": dschema.SetAttribute{
				MarkdownDescription: "Interfaces spanning tree is enabled on.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"max_age": dschema.Int64Attribute{
				MarkdownDescription: "Time (in seconds) that a spanning tree configuration is valid.",
				Computed:            true,
			},
			"forward_delay": dschema.Int64Attribute{
				MarkdownDescription: "Time (in seconds) that must pass before an interface begins forwarding packets when spanning tree is enabled.",
				Computed:            true,
			},
			"hold_count": dschema.Int64Attribute{
				MarkdownDescription: "Number of packets transmitted before being rate limited when spanning tree is enabled.",
				Computed:            true,
			},
			"max_addresses": dschema.Int64Attribute{
				MarkdownDescription: "Maximum size of the bridge address cache.",
				Computed:            true,
			},
			"address_timeout": dschema.Int64Attribute{
				MarkdownDescription: "Timeout (in seconds) of bridge address cache entries.",
				Computed:            true,
			},
			"span_interface": dschema.StringAttribute{
				MarkdownDescription: "Interface a copy of every frame received by the bridge is transmitted on.",
				Computed:            true,
			},
			"edge_ports": dschema.SetAttribute{
				MarkdownDescription: "Interfaces set as edge ports.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"auto_edge_ports": dschema.SetAttribute{
				MarkdownDescription: "Interfaces automatically detecting edge status.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"ptp_ports": dschema.SetAttribute{
				MarkdownDescription: "Interfaces set as point-to-point links.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"auto_ptp_ports": dschema.SetAttribute{
				MarkdownDescription: "Interfaces automatically detecting point-to-point status.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"static_ports": dschema.SetAttribute{
				MarkdownDescription: "Interfaces marked as static.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"private_ports": dschema.SetAttribute{
				MarkdownDescription: "Interfaces marked as private.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Computed:            true,
			},
			"device": dschema.StringAttribute{
				MarkdownDescription: "Bridge device name, e.g. `bridge0`.",
				Computed:            true,
			},
		},
	}
}

func convertBridgeSchemaToStruct(d *bridgeResourceModel) (*interfaces.Bridge, error) {
	return &interfaces.Bridge{
		Device:         d.Device.ValueString(),
		Members:        api.SelectedMapList(tools.SetToStringSlice(d.Members)),
		LinkLocal:      tools.BoolToString(d.LinkLocal.ValueBool()),
		EnableStp:      tools.BoolToString(d.EnableStp.ValueBool()),
		StpProto:       api.SelectedMap(d.StpProto.ValueString()),
		StpInterfaces:  api.SelectedMapList(tools.SetToStringSlice(d.StpInterfaces)),
		MaxAge:         int64ToStringEmpty(d.MaxAge),
		ForwardDelay:   int64ToStringEmpty(d.ForwardDelay),
		HoldCount:      int64ToStringEmpty(d.HoldCount),
		MaxAddresses:   int64ToStringEmpty(d.MaxAddresses),
		AddressTimeout: int64ToStringEmpty(d.AddressTimeout),
		SpanInterface:  api.SelectedMap(d.SpanInterface.ValueString()),
		EdgePorts:      api.SelectedMapList(tools.SetToStringSlice(d.EdgePorts)),
		AutoEdgePorts:  api.SelectedMapList(tools.SetToStringSlice(d.AutoEdgePorts)),
		PtpPorts:       api.SelectedMapList(tools.SetToStringSlice(d.PtpPorts)),
		AutoPtpPorts:   api.SelectedMapList(tools.SetToStringSlice(d.AutoPtpPorts)),
		StaticPorts:    api.SelectedMapList(tools.SetToStringSlice(d.StaticPorts)),
		PrivatePorts:   api.SelectedMapList(tools.SetToStringSlice(d.PrivatePorts)),
		Description:    d.Description.ValueString(),
	}, nil
}

func convertBridgeStructToSchema(d *interfaces.Bridge) (*bridgeResourceModel, error) {
	// The API reports the span interface's "None" option as an empty selection,
	// which SelectedMap renders as "". The schema default is also "".
	spanInterface := d.SpanInterface.String()

	return &bridgeResourceModel{
		Device:         types.StringValue(d.Device),
		Members:        tools.StringSliceToSet(d.Members),
		LinkLocal:      types.BoolValue(tools.StringToBool(d.LinkLocal)),
		EnableStp:      types.BoolValue(tools.StringToBool(d.EnableStp)),
		StpProto:       types.StringValue(d.StpProto.String()),
		StpInterfaces:  tools.StringSliceToSet(d.StpInterfaces),
		MaxAge:         tools.StringToInt64Null(d.MaxAge),
		ForwardDelay:   tools.StringToInt64Null(d.ForwardDelay),
		HoldCount:      tools.StringToInt64Null(d.HoldCount),
		MaxAddresses:   tools.StringToInt64Null(d.MaxAddresses),
		AddressTimeout: tools.StringToInt64Null(d.AddressTimeout),
		SpanInterface:  types.StringValue(spanInterface),
		EdgePorts:      tools.StringSliceToSet(d.EdgePorts),
		AutoEdgePorts:  tools.StringSliceToSet(d.AutoEdgePorts),
		PtpPorts:       tools.StringSliceToSet(d.PtpPorts),
		AutoPtpPorts:   tools.StringSliceToSet(d.AutoPtpPorts),
		StaticPorts:    tools.StringSliceToSet(d.StaticPorts),
		PrivatePorts:   tools.StringSliceToSet(d.PrivatePorts),
		Description:    tools.StringOrNull(d.Description),
	}, nil
}

// int64ToStringEmpty renders a null Int64 as the empty string the OPNsense
// API expects for unset integer fields.
func int64ToStringEmpty(v types.Int64) string {
	if v.IsNull() {
		return ""
	}
	return tools.Int64ToString(v.ValueInt64())
}
