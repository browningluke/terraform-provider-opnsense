package service

import (
	"github.com/browningluke/opnsense-go/pkg/ipsec"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IpsecVtiResourceModel describes the resource data model.
type IpsecVtiResourceModel struct {
	Enabled         types.String `tfsdk:"enabled"`
	RequestID       types.String `tfsdk:"request_id"`
	LocalIP         types.String `tfsdk:"local_ip"`
	RemoteIP        types.String `tfsdk:"remote_ip"`
	TunnelLocalIP   types.String `tfsdk:"tunnel_local_ip"`
	TunnelRemoteIP  types.String `tfsdk:"tunnel_remote_ip"`
	TunnelLocalIP2  types.String `tfsdk:"tunnel_local_ip2"`
	TunnelRemoteIP2 types.String `tfsdk:"tunnel_remote_ip2"`
	Description     types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

func IpsecVtiResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "IPsec Virtual Tunnel Interfaces (VTIs) are used by routed IPsec VPN connections.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.StringAttribute{
				MarkdownDescription: "Enable or disable the VTI.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("1"),
			},
			"request_id": schema.StringAttribute{
				MarkdownDescription: "Request ID for the VTI.",
				Required:            true,
			},
			"local_ip": schema.StringAttribute{
				MarkdownDescription: "Local IP address for the VTI.",
				Required:            true,
			},
			"remote_ip": schema.StringAttribute{
				MarkdownDescription: "Remote IP address for the VTI.",
				Required:            true,
			},
			"tunnel_local_ip": schema.StringAttribute{
				MarkdownDescription: "Local tunnel IP address for the VTI.",
				Required:            true,
			},
			"tunnel_remote_ip": schema.StringAttribute{
				MarkdownDescription: "Remote tunnel IP address for the VTI.",
				Required:            true,
			},
			"tunnel_local_ip2": schema.StringAttribute{
				MarkdownDescription: "Second local tunnel IP address for the VTI.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"tunnel_remote_ip2": schema.StringAttribute{
				MarkdownDescription: "Second remote tunnel IP address for the VTI.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description for the VTI.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"id": schema.StringAttribute{
				Description: "UUID of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func IpsecVtiDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "IPsec Virtual Tunnel Interfaces (VTIs) are used by routed IPsec VPN connections.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.StringAttribute{
				MarkdownDescription: "Enable or disable the VTI.",
				Computed:            true,
			},
			"request_id": dschema.StringAttribute{
				MarkdownDescription: "Request ID for the VTI.",
				Computed:            true,
			},
			"local_ip": dschema.StringAttribute{
				MarkdownDescription: "Local IP address for the VTI.",
				Computed:            true,
			},
			"remote_ip": dschema.StringAttribute{
				MarkdownDescription: "Remote IP address for the VTI.",
				Computed:            true,
			},
			"tunnel_local_ip": dschema.StringAttribute{
				MarkdownDescription: "Local tunnel IP address for the VTI.",
				Computed:            true,
			},
			"tunnel_remote_ip": dschema.StringAttribute{
				MarkdownDescription: "Remote tunnel IP address for the VTI.",
				Computed:            true,
			},
			"tunnel_local_ip2": dschema.StringAttribute{
				MarkdownDescription: "Second local tunnel IP address for the VTI.",
				Computed:            true,
			},
			"tunnel_remote_ip2": dschema.StringAttribute{
				MarkdownDescription: "Second remote tunnel IP address for the VTI.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "Optional description for the VTI.",
				Computed:            true,
			},
		},
	}
}

func convertIpsecVtiSchemaToStruct(d *IpsecVtiResourceModel) (*ipsec.IPsecVTI, error) {
	return &ipsec.IPsecVTI{
		Enabled:         d.Enabled.ValueString(),
		RequestID:       d.RequestID.ValueString(),
		LocalIP:         d.LocalIP.ValueString(),
		RemoteIP:        d.RemoteIP.ValueString(),
		TunnelLocalIP:   d.TunnelLocalIP.ValueString(),
		TunnelRemoteIP:  d.TunnelRemoteIP.ValueString(),
		TunnelLocalIP2:  d.TunnelLocalIP2.ValueString(),
		TunnelRemoteIP2: d.TunnelRemoteIP2.ValueString(),
		Description:     d.Description.ValueString(),
	}, nil
}

func convertIpsecVtiStructToSchema(d *ipsec.IPsecVTI) (*IpsecVtiResourceModel, error) {
	return &IpsecVtiResourceModel{
		Enabled:         types.StringValue(d.Enabled),
		RequestID:       types.StringValue(d.RequestID),
		LocalIP:         types.StringValue(d.LocalIP),
		RemoteIP:        types.StringValue(d.RemoteIP),
		TunnelLocalIP:   types.StringValue(d.TunnelLocalIP),
		TunnelRemoteIP:  types.StringValue(d.TunnelRemoteIP),
		TunnelLocalIP2:  types.StringValue(d.TunnelLocalIP2),
		TunnelRemoteIP2: types.StringValue(d.TunnelRemoteIP2),
		Description:     types.StringValue(d.Description),
	}, nil
}
