package service

import (
	"github.com/browningluke/opnsense-go/pkg/unbound"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-opnsense/internal/tools"
)

// UnboundForwardResourceModel describes the resource data model.
type UnboundForwardResourceModel struct {
	Enabled    types.Bool   `tfsdk:"enabled"`
	Domain     types.String `tfsdk:"domain"`
	ServerIP   types.String `tfsdk:"server_ip"`
	ServerPort types.Int64  `tfsdk:"server_port"`
	VerifyCN   types.String `tfsdk:"verify_cn"`

	Id types.String `tfsdk:"id"`
}

func unboundForwardResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Query Forwarding section allows for entering arbitrary nameservers to forward queries to. Can forward queries normally, or over TLS.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this query forward.  Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "If a domain is entered here, queries for this specific domain will be forwarded to the specified server. Set to `\"\"` to forward all queries to the specified server.",
				Required:            true,
			},
			"server_ip": schema.StringAttribute{
				MarkdownDescription: "IP address of DNS server to forward all requests.",
				Required:            true,
			},
			"server_port": schema.Int64Attribute{
				MarkdownDescription: "Port of DNS server, for usual DNS use `53`, if you use DoT set it to `853`. Defaults to `53`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(53),
			},
			"verify_cn": schema.StringAttribute{
				MarkdownDescription: "The Common Name of the DNS server (e.g. `dns.example.com`). This field is required to verify its TLS certificate. DNS-over-TLS is susceptible to man-in-the-middle attacks unless certificates can be verified. Set to `\"\"` to accept self-signed yet also potentially fraudulent certificates. Must be set when `type` is `dot`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the forward.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func UnboundForwardDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Query Forwarding section allows for entering arbitrary nameservers to forward queries to. Can forward queries normally, or over TLS.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Whether this route is enabled.",
				Computed:            true,
			},
			"domain": dschema.StringAttribute{
				MarkdownDescription: "If a domain is entered here, queries for this specific domain will be forwarded to the specified server.",
				Computed:            true,
			},
			"server_ip": dschema.StringAttribute{
				MarkdownDescription: "IP address of DNS server to forward all requests.",
				Computed:            true,
			},
			"server_port": dschema.Int64Attribute{
				MarkdownDescription: "Port of DNS server, for usual DNS use `53`, if you use DoT set it to `853`.",
				Computed:            true,
			},
			"verify_cn": dschema.StringAttribute{
				MarkdownDescription: "The Common Name of the DNS server (e.g. `dns.example.com`). This field is required to verify its TLS certificate. DNS-over-TLS is susceptible to man-in-the-middle attacks unless certificates can be verified.",
				Computed:            true,
			},
		},
	}
}

func convertUnboundForwardSchemaToStruct(d *UnboundForwardResourceModel) (*unbound.Forward, error) {
	return &unbound.Forward{
		Enabled:  tools.BoolToString(d.Enabled.ValueBool()),
		Domain:   d.Domain.ValueString(),
		Server:   d.ServerIP.ValueString(),
		Port:     tools.Int64ToString(d.ServerPort.ValueInt64()),
		VerifyCN: d.VerifyCN.ValueString(),
	}, nil
}

func convertUnboundForwardStructToSchema(d *unbound.Forward) (*UnboundForwardResourceModel, error) {
	return &UnboundForwardResourceModel{
		Enabled:    types.BoolValue(tools.StringToBool(d.Enabled)),
		Domain:     types.StringValue(d.Domain),
		ServerIP:   types.StringValue(d.Server),
		ServerPort: types.Int64Value(tools.StringToInt64(d.Port)),
		VerifyCN:   types.StringValue(d.VerifyCN),
	}, nil
}
