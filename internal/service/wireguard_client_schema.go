package service

import (
	"context"
	"github.com/browningluke/opnsense-go/pkg/wireguard"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-opnsense/internal/tools"
)

// WireguardClientResourceModel describes the resource data model.
type WireguardClientResourceModel struct {
	Enabled       types.Bool   `tfsdk:"enabled"`
	Name          types.String `tfsdk:"name"`
	PublicKey     types.String `tfsdk:"public_key"`
	PSK           types.String `tfsdk:"psk"`
	ServerAddress types.String `tfsdk:"server_address"`
	ServerPort    types.Int64  `tfsdk:"server_port"`
	TunnelAddress types.Set    `tfsdk:"tunnel_address"`
	KeepAlive     types.Int64  `tfsdk:"keep_alive"`

	Id types.String `tfsdk:"id"`
}

func wireguardClientResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Client resources can be used to setup Wireguard clients.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this client config. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the client config.",
				Required:            true,
			},
			"public_key": schema.StringAttribute{
				MarkdownDescription: "Public key of this client config. Must be a 256-bit base64 string.",
				Required:            true,
			},
			"psk": schema.StringAttribute{
				MarkdownDescription: "Shared secret (PSK) for this peer. You can generate a key using `wg genpsk` on a client with WireGuard installed. Must be a 256-bit base64 string. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"server_address": schema.StringAttribute{
				MarkdownDescription: "The public IP address the endpoint listens to. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"server_port": schema.Int64Attribute{
				MarkdownDescription: "The port the endpoint listens to. Defaults to `-1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"tunnel_address": schema.SetAttribute{
				MarkdownDescription: "List of addresses allowed to pass trough the tunnel adapter. Please use CIDR notation like `\"10.0.0.1/24\"`. Defaults to `[]`.",
				Required:            true,
				ElementType:         types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"keep_alive": schema.Int64Attribute{
				MarkdownDescription: "The persistent keepalive interval in seconds. Defaults to `-1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the client.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func WireguardClientDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Client resources can be used to setup Wireguard clients.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Whether this client config is enabled.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the client config.",
				Computed:            true,
			},
			"public_key": schema.StringAttribute{
				MarkdownDescription: "Public key of this client config.",
				Computed:            true,
			},
			"psk": schema.StringAttribute{
				MarkdownDescription: "Shared secret (PSK) for this peer.",
				Computed:            true,
			},
			"server_address": schema.StringAttribute{
				MarkdownDescription: "The public IP address the endpoint listens to.",
				Computed:            true,
			},
			"server_port": schema.Int64Attribute{
				MarkdownDescription: "The port the endpoint listens to.",
				Computed:            true,
			},
			"tunnel_address": schema.SetAttribute{
				MarkdownDescription: "List of addresses allowed to pass trough the tunnel adapter.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"keep_alive": schema.Int64Attribute{
				MarkdownDescription: "The persistent keepalive interval in seconds.",
				Computed:            true,
			},
		},
	}
}

func convertWireguardClientSchemaToStruct(d *WireguardClientResourceModel) (*wireguard.Client, error) {
	// Parse 'TunnelAddress'
	var tunnelAddressList []string
	d.TunnelAddress.ElementsAs(context.Background(), &tunnelAddressList, false)

	return &wireguard.Client{
		Enabled:       tools.BoolToString(d.Enabled.ValueBool()),
		Name:          d.Name.ValueString(),
		PublicKey:     d.PublicKey.ValueString(),
		PSK:           d.PSK.ValueString(),
		ServerAddress: d.ServerAddress.ValueString(),
		ServerPort:    tools.Int64ToStringNegative(d.ServerPort.ValueInt64()),
		TunnelAddress: tunnelAddressList,
		KeepAlive:     tools.Int64ToStringNegative(d.KeepAlive.ValueInt64()),
	}, nil
}

func convertWireguardClientStructToSchema(d *wireguard.Client) (*WireguardClientResourceModel, error) {
	model := &WireguardClientResourceModel{
		Enabled:       types.BoolValue(tools.StringToBool(d.Enabled)),
		Name:          types.StringValue(d.Name),
		PublicKey:     types.StringValue(d.PublicKey),
		PSK:           types.StringValue(d.PSK),
		ServerAddress: types.StringValue(d.ServerAddress),
		ServerPort:    types.Int64Value(tools.StringToInt64(d.ServerPort)),
		TunnelAddress: types.SetNull(types.StringType),
		KeepAlive:     types.Int64Value(tools.StringToInt64(d.KeepAlive)),
	}

	// Parse 'TunnelAddress'
	model.TunnelAddress = tools.StringSliceToSet(d.TunnelAddress)

	return model, nil
}
