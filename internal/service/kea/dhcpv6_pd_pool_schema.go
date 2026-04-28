package kea

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/kea"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// dhcpv6PdPoolResourceModel describes the resource data model.
type dhcpv6PdPoolResourceModel struct {
	SubnetId     types.String `tfsdk:"subnet_id"`
	Prefix       types.String `tfsdk:"prefix"`
	PrefixLen    types.String `tfsdk:"prefix_len"`
	DelegatedLen types.String `tfsdk:"delegated_len"`
	Description  types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

func dhcpv6PdPoolResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Configure DHCPv6 prefix delegation pools for Kea.",

		Attributes: map[string]schema.Attribute{
			"subnet_id": schema.StringAttribute{
				MarkdownDescription: "Subnet ID the PD pool belongs to.",
				Required:            true,
			},
			"prefix": schema.StringAttribute{
				MarkdownDescription: "IPv6 prefix for the PD pool (e.g. `\"2001:db8::/48\"`).",
				Required:            true,
			},
			"prefix_len": schema.StringAttribute{
				MarkdownDescription: "Prefix length of the PD pool prefix.",
				Required:            true,
			},
			"delegated_len": schema.StringAttribute{
				MarkdownDescription: "Length of the delegated prefix.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the PD pool.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func dhcpv6PdPoolDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure DHCPv6 prefix delegation pools for Kea.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the PD pool.",
				Required:            true,
			},
			"subnet_id": dschema.StringAttribute{
				MarkdownDescription: "Subnet ID the PD pool belongs to.",
				Computed:            true,
			},
			"prefix": dschema.StringAttribute{
				MarkdownDescription: "IPv6 prefix for the PD pool.",
				Computed:            true,
			},
			"prefix_len": dschema.StringAttribute{
				MarkdownDescription: "Prefix length of the PD pool prefix.",
				Computed:            true,
			},
			"delegated_len": dschema.StringAttribute{
				MarkdownDescription: "Length of the delegated prefix.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "Optional description.",
				Computed:            true,
			},
		},
	}
}

func convertDhcpv6PdPoolSchemaToStruct(d *dhcpv6PdPoolResourceModel) (*kea.PDPool, error) {
	return &kea.PDPool{
		Subnet:       api.SelectedMap(d.SubnetId.ValueString()),
		Prefix:       d.Prefix.ValueString(),
		PrefixLen:    d.PrefixLen.ValueString(),
		DelegatedLen: d.DelegatedLen.ValueString(),
		Description:  d.Description.ValueString(),
	}, nil
}

func convertDhcpv6PdPoolStructToSchema(d *kea.PDPool) (*dhcpv6PdPoolResourceModel, error) {
	return &dhcpv6PdPoolResourceModel{
		SubnetId:     types.StringValue(d.Subnet.String()),
		Prefix:       types.StringValue(d.Prefix),
		PrefixLen:    types.StringValue(d.PrefixLen),
		DelegatedLen: types.StringValue(d.DelegatedLen),
		Description:  types.StringValue(d.Description),
	}, nil
}
