package kea

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/kea"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// dhcpv4PeerResourceModel describes the resource data model.
type dhcpv4PeerResourceModel struct {
	Name types.String `tfsdk:"name"`
	Url  types.String `tfsdk:"url"`
	Role types.String `tfsdk:"role"`

	Id types.String `tfsdk:"id"`
}

func dhcpv4PeerResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Configure HA Peers for Kea DHCPv4.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Peer name, there should be one entry matching this machine's \"This server name\".",
				Required:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL of the server instance, which should use a different port than the control agent (e.g. `http://192.0.2.1:8001/`).",
				Required:            true,
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "Peer's role. Defaults to `\"primary\"`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("primary", "standby"),
				},
				Default: stringdefault.StaticString("primary"),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the peer.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func dhcpv4PeerDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure HA Peers for Kea DHCPv4.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the peer.",
				Required:            true,
			},
			"name": dschema.StringAttribute{
				MarkdownDescription: "Peer name, there should be one entry matching this machine's \"This server name\".",
				Computed:            true,
			},
			"url": dschema.StringAttribute{
				MarkdownDescription: "URL of the server instance.",
				Computed:            true,
			},
			"role": dschema.StringAttribute{
				MarkdownDescription: "Peer's role.",
				Computed:            true,
			},
		},
	}
}

func convertDhcpv4PeerSchemaToStruct(d *dhcpv4PeerResourceModel) (*kea.PeerV4, error) {
	return &kea.PeerV4{
		Name: d.Name.ValueString(),
		Url:  d.Url.ValueString(),
		Role: api.SelectedMap(d.Role.ValueString()),
	}, nil
}

func convertDhcpv4PeerStructToSchema(d *kea.PeerV4) (*dhcpv4PeerResourceModel, error) {
	return &dhcpv4PeerResourceModel{
		Name: types.StringValue(d.Name),
		Url:  types.StringValue(d.Url),
		Role: types.StringValue(d.Role.String()),
	}, nil
}
