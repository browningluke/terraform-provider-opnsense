package wireguard

import (
	"context"

	"github.com/browningluke/opnsense-go/pkg/wireguard"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// serverResourceModel describes the resource data model.
type serverResourceModel struct {
	Enabled       types.Bool   `tfsdk:"enabled"`
	Name          types.String `tfsdk:"name"`
	PublicKey     types.String `tfsdk:"public_key"`
	PrivateKey    types.String `tfsdk:"private_key"`
	Port          types.Int64  `tfsdk:"port"`
	MTU           types.Int64  `tfsdk:"mtu"`
	DNS           types.Set    `tfsdk:"dns"`
	TunnelAddress types.Set    `tfsdk:"tunnel_address"`
	Peers         types.Set    `tfsdk:"peers"`
	DisableRoutes types.Bool   `tfsdk:"disable_routes"`
	Gateway       types.String `tfsdk:"gateway"`

	Id       types.String `tfsdk:"id"`
	Instance types.String `tfsdk:"instance"`
}

func serverResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Server resources can be used to setup Wireguard servers.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this server. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the server.",
				Required:            true,
			},
			"public_key": schema.StringAttribute{
				MarkdownDescription: "Public key of this server. Must be a 256-bit base64 string.",
				Required:            true,
			},
			"private_key": schema.StringAttribute{
				MarkdownDescription: "Private key of this server. Must be a 256-bit base64 string.",
				Required:            true,
				Sensitive:           true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "The fixed port for this instance to listen on. The standard port range starts at 51820. Defaults to `-1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"mtu": schema.Int64Attribute{
				MarkdownDescription: "The interface MTU for this interface. Set to `-1` to use the MTU from main interface. Defaults to `-1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"dns": schema.SetAttribute{
				MarkdownDescription: "The interface specific DNS servers. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
				ElementType:         types.StringType,
			},
			"tunnel_address": schema.SetAttribute{
				MarkdownDescription: "List of addresses to configure on the tunnel adapter. Please use CIDR notation like `\"10.0.0.1/24\"`. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
				ElementType:         types.StringType,
			},
			"peers": schema.SetAttribute{
				MarkdownDescription: "List of peer IDs for this server. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
				ElementType:         types.StringType,
			},
			"disable_routes": schema.BoolAttribute{
				MarkdownDescription: "Disables installation of routes. Usually you only enable this to do own routing decisions via a local gateway and gateway rules. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"gateway": schema.StringAttribute{
				MarkdownDescription: "The gateway IP here when using Disable Routes feature. You also have to add this as a gateway in OPNsense. Must be set when `disable_routes` is `true`. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("disable_routes"),
					}...),
				},
			},
			"instance": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The instance number to give the wg interface a unique name (wgX).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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

func serverDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Server resources can be used to setup Wireguard servers.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Whether this server is enabled.",
				Computed:            true,
			},
			"name": dschema.StringAttribute{
				MarkdownDescription: "Name of the server.",
				Computed:            true,
			},

			"public_key": schema.StringAttribute{
				MarkdownDescription: "Public key of this server.",
				Computed:            true,
			},
			"private_key": schema.StringAttribute{
				MarkdownDescription: "Private key of this server.",
				Computed:            true,
				Sensitive:           true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "The fixed port for this instance to listen on. The standard port range starts at 51820.",
				Computed:            true,
			},
			"mtu": schema.Int64Attribute{
				MarkdownDescription: "The interface MTU for this interface. Set to `-1` to use the MTU from main interface.",
				Computed:            true,
			},
			"dns": schema.SetAttribute{
				MarkdownDescription: "The interface specific DNS servers.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"tunnel_address": schema.SetAttribute{
				MarkdownDescription: "List of addresses to configure on the tunnel adapter.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"peers": schema.SetAttribute{
				MarkdownDescription: "List of peer IDs for this server.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"disable_routes": schema.BoolAttribute{
				MarkdownDescription: "Disables installation of routes.",
				Computed:            true,
			},
			"gateway": schema.StringAttribute{
				MarkdownDescription: "The gateway IP here when using Disable Routes feature.",
				Computed:            true,
			},
			"instance": dschema.StringAttribute{
				MarkdownDescription: "The instance number to give the wg interface a unique name (wgX).",
				Computed:            true,
			},
		},
	}
}

func convertServerSchemaToStruct(d *serverResourceModel) (*wireguard.Server, error) {
	// Parse 'DNS'
	var dnsList []string
	d.DNS.ElementsAs(context.Background(), &dnsList, false)

	// Parse 'TunnelAddress'
	var tunnelAddressList []string
	d.TunnelAddress.ElementsAs(context.Background(), &tunnelAddressList, false)

	// Parse 'Peers'
	var peersList []string
	d.Peers.ElementsAs(context.Background(), &peersList, false)

	return &wireguard.Server{
		Enabled:       tools.BoolToString(d.Enabled.ValueBool()),
		Name:          d.Name.ValueString(),
		PublicKey:     d.PublicKey.ValueString(),
		PrivateKey:    d.PrivateKey.ValueString(),
		Port:          tools.Int64ToStringNegative(d.Port.ValueInt64()),
		MTU:           tools.Int64ToStringNegative(d.MTU.ValueInt64()),
		DNS:           dnsList,
		TunnelAddress: tunnelAddressList,
		Peers:         peersList,
		DisableRoutes: tools.BoolToString(d.DisableRoutes.ValueBool()),
		Gateway:       d.Gateway.ValueString(),
		Instance:      "", // Instance must be set, but always to an empty string (plugin requirement)
	}, nil
}

func convertServerStructToSchema(d *wireguard.Server) (*serverResourceModel, error) {
	model := &serverResourceModel{
		Enabled:       types.BoolValue(tools.StringToBool(d.Enabled)),
		Name:          types.StringValue(d.Name),
		PublicKey:     types.StringValue(d.PublicKey),
		PrivateKey:    types.StringValue(d.PrivateKey),
		Port:          types.Int64Value(tools.StringToInt64(d.Port)),
		MTU:           types.Int64Value(tools.StringToInt64(d.MTU)),
		DNS:           types.SetNull(types.StringType),
		TunnelAddress: types.SetNull(types.StringType),
		Peers:         types.SetNull(types.StringType),
		DisableRoutes: types.BoolValue(tools.StringToBool(d.DisableRoutes)),
		Gateway:       types.StringValue(d.Gateway),
		Instance:      types.StringValue(d.Instance),
	}

	// Parse 'DNS'
	model.DNS = tools.StringSliceToSet(d.DNS)

	// Parse 'TunnelAddress'
	model.TunnelAddress = tools.StringSliceToSet(d.TunnelAddress)

	// Parse 'Peers'
	model.Peers = tools.StringSliceToSet(d.Peers)

	return model, nil
}
