package quagga

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/quagga"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// bgpPeerGroupResourceModel describes the resource data model.
type bgpPeerGroupResourceModel struct {
	Enabled          types.Bool   `tfsdk:"enabled"`
	Name             types.String `tfsdk:"name"`
	RemoteASMode     types.String `tfsdk:"remote_as_mode"`
	RemoteAS         types.String `tfsdk:"remote_as"`
	Family           types.String `tfsdk:"family"`
	ListenRanges     types.Set    `tfsdk:"listen_ranges"`
	UpdateSource     types.String `tfsdk:"update_source"`
	NextHopSelf      types.Bool   `tfsdk:"next_hop_self"`
	DefaultOriginate types.Bool   `tfsdk:"default_originate"`
	PrefixListIn     types.String `tfsdk:"prefix_list_in"`
	PrefixListOut    types.String `tfsdk:"prefix_list_out"`
	RouteMapIn       types.String `tfsdk:"route_map_in"`
	RouteMapOut      types.String `tfsdk:"route_map_out"`

	Id types.String `tfsdk:"id"`
}

func bgpPeerGroupResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Configure peer groups for BGP.",
		Version:             1,

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this peer group. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the peer group.",
				Required:            true,
			},
			"remote_as_mode": schema.StringAttribute{
				MarkdownDescription: "The remote AS mode for this peer group. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.OneOf("", "internal", "external"),
				},
			},
			"remote_as": schema.StringAttribute{
				MarkdownDescription: "The remote AS number for this peer group. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"family": schema.StringAttribute{
				MarkdownDescription: "The address family for this peer group. Defaults to `\"ipv4\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("ipv4"),
				Validators: []validator.String{
					stringvalidator.OneOf("ipv4", "ipv6"),
				},
			},
			"listen_ranges": schema.SetAttribute{
				MarkdownDescription: "The listen ranges for this peer group.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"update_source": schema.StringAttribute{
				MarkdownDescription: "The update source interface for this peer group. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"next_hop_self": schema.BoolAttribute{
				MarkdownDescription: "Enable next-hop-self for this peer group. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"default_originate": schema.BoolAttribute{
				MarkdownDescription: "Enable sending the default route to this peer group. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"prefix_list_in": schema.StringAttribute{
				MarkdownDescription: "The prefix list UUID for inbound direction. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"prefix_list_out": schema.StringAttribute{
				MarkdownDescription: "The prefix list UUID for outbound direction. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"route_map_in": schema.StringAttribute{
				MarkdownDescription: "The route map UUID for inbound direction. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"route_map_out": schema.StringAttribute{
				MarkdownDescription: "The route map UUID for outbound direction. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the peer group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func bgpPeerGroupDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure peer groups for BGP.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable this peer group.",
				Computed:            true,
			},
			"name": dschema.StringAttribute{
				MarkdownDescription: "The name of the peer group.",
				Computed:            true,
			},
			"remote_as_mode": dschema.StringAttribute{
				MarkdownDescription: "The remote AS mode for this peer group.",
				Computed:            true,
			},
			"remote_as": dschema.StringAttribute{
				MarkdownDescription: "The remote AS number for this peer group.",
				Computed:            true,
			},
			"family": dschema.StringAttribute{
				MarkdownDescription: "The address family for this peer group.",
				Computed:            true,
			},
			"listen_ranges": dschema.SetAttribute{
				MarkdownDescription: "The listen ranges for this peer group.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"update_source": dschema.StringAttribute{
				MarkdownDescription: "The update source interface for this peer group.",
				Computed:            true,
			},
			"next_hop_self": dschema.BoolAttribute{
				MarkdownDescription: "Enable next-hop-self for this peer group.",
				Computed:            true,
			},
			"default_originate": dschema.BoolAttribute{
				MarkdownDescription: "Enable sending the default route to this peer group.",
				Computed:            true,
			},
			"prefix_list_in": dschema.StringAttribute{
				MarkdownDescription: "The prefix list UUID for inbound direction.",
				Computed:            true,
			},
			"prefix_list_out": dschema.StringAttribute{
				MarkdownDescription: "The prefix list UUID for outbound direction.",
				Computed:            true,
			},
			"route_map_in": dschema.StringAttribute{
				MarkdownDescription: "The route map UUID for inbound direction.",
				Computed:            true,
			},
			"route_map_out": dschema.StringAttribute{
				MarkdownDescription: "The route map UUID for outbound direction.",
				Computed:            true,
			},
		},
	}
}

func convertBGPPeerGroupSchemaToStruct(d *bgpPeerGroupResourceModel) (*quagga.BGPPeerGroup, error) {
	return &quagga.BGPPeerGroup{
		Enabled:          tools.BoolToString(d.Enabled.ValueBool()),
		Name:             d.Name.ValueString(),
		RemoteASMode:     api.SelectedMap(d.RemoteASMode.ValueString()),
		RemoteAS:         d.RemoteAS.ValueString(),
		Family:           api.SelectedMap(d.Family.ValueString()),
		ListenRanges:     api.SelectedMapList(tools.SetToStringSlice(d.ListenRanges)),
		UpdateSource:     api.SelectedMap(d.UpdateSource.ValueString()),
		NextHopSelf:      tools.BoolToString(d.NextHopSelf.ValueBool()),
		DefaultOriginate: tools.BoolToString(d.DefaultOriginate.ValueBool()),
		PrefixListIn:     api.SelectedMap(d.PrefixListIn.ValueString()),
		PrefixListOut:    api.SelectedMap(d.PrefixListOut.ValueString()),
		RouteMapIn:       api.SelectedMap(d.RouteMapIn.ValueString()),
		RouteMapOut:      api.SelectedMap(d.RouteMapOut.ValueString()),
	}, nil
}

func convertBGPPeerGroupStructToSchema(d *quagga.BGPPeerGroup) (*bgpPeerGroupResourceModel, error) {
	return &bgpPeerGroupResourceModel{
		Enabled:          types.BoolValue(tools.StringToBool(d.Enabled)),
		Name:             types.StringValue(d.Name),
		RemoteASMode:     types.StringValue(d.RemoteASMode.String()),
		RemoteAS:         types.StringValue(d.RemoteAS),
		Family:           types.StringValue(d.Family.String()),
		ListenRanges:     tools.StringSliceToSet([]string(d.ListenRanges)),
		UpdateSource:     types.StringValue(d.UpdateSource.String()),
		NextHopSelf:      types.BoolValue(tools.StringToBool(d.NextHopSelf)),
		DefaultOriginate: types.BoolValue(tools.StringToBool(d.DefaultOriginate)),
		PrefixListIn:     types.StringValue(d.PrefixListIn.String()),
		PrefixListOut:    types.StringValue(d.PrefixListOut.String()),
		RouteMapIn:       types.StringValue(d.RouteMapIn.String()),
		RouteMapOut:      types.StringValue(d.RouteMapOut.String()),
	}, nil
}
