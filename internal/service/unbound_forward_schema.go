package service

import (
	"fmt"
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/unbound"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
)

// UnboundForwardResourceModel describes the resource data model.
type UnboundForwardResourceModel struct {
	Enabled    types.Bool   `tfsdk:"enabled"`
	Type       types.String `tfsdk:"type"`
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
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of forward. Available values: `forward`, `dot`. Defaults to `forward`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("forward"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"forward", "dot"}...),
				},
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

func convertUnboundForwardSchemaToStruct(d *UnboundForwardResourceModel) (*unbound.Forward, error) {
	// Parse 'Enabled'
	var enabled string
	if d.Enabled.ValueBool() {
		enabled = "1"
	} else {
		enabled = "0"
	}

	return &unbound.Forward{
		Enabled:  enabled,
		Domain:   d.Domain.ValueString(),
		Type:     api.SelectedMap(d.Type.ValueString()),
		Server:   d.ServerIP.ValueString(),
		Port:     fmt.Sprintf("%d", d.ServerPort.ValueInt64()),
		VerifyCN: d.VerifyCN.ValueString(),
	}, nil
}

func convertUnboundForwardStructToSchema(d *unbound.Forward) (*UnboundForwardResourceModel, error) {
	// Parse 'ServerPort'
	serverPort, err := strconv.ParseInt(d.Port, 10, 64)
	if err != nil {
		return nil, err
	}

	model := &UnboundForwardResourceModel{
		Enabled:    types.BoolValue(false),
		Domain:     types.StringValue(d.Domain),
		Type:       types.StringValue(d.Type.String()),
		ServerIP:   types.StringValue(d.Server),
		ServerPort: types.Int64Value(serverPort),
		VerifyCN:   types.StringValue(d.VerifyCN),
	}

	// Parse 'Enabled'
	if d.Enabled == "1" {
		model.Enabled = types.BoolValue(true)
	}

	return model, nil
}
