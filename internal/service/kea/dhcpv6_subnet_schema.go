package kea

import (
	"strings"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/kea"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// dhcpv6SubnetResourceModel describes the resource data model.
type dhcpv6SubnetResourceModel struct {
	Subnet      types.String `tfsdk:"subnet"`
	Allocator   types.String `tfsdk:"allocator"`
	PDAllocator types.String `tfsdk:"pd_allocator"`
	Pools       types.Set    `tfsdk:"pools"`
	Interface   types.String `tfsdk:"interface"`

	DnsServers   types.Set `tfsdk:"dns_servers"`
	DomainSearch types.Set `tfsdk:"domain_search"`

	Description types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

func dhcpv6SubnetResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Configure DHCPv6 subnets for Kea.",

		Attributes: map[string]schema.Attribute{
			"subnet": schema.StringAttribute{
				MarkdownDescription: "IPv6 Subnet (e.g. `\"2001:db8::/64\"`).",
				Required:            true,
			},
			"allocator": schema.StringAttribute{
				MarkdownDescription: "Address allocator to use. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"pd_allocator": schema.StringAttribute{
				MarkdownDescription: "Prefix delegation allocator. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"pools": schema.SetAttribute{
				MarkdownDescription: "Set of address pools (e.g. `\"2001:db8::100 - 2001:db8::200\"`). Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"interface": schema.StringAttribute{
				MarkdownDescription: "Interface to listen on. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"dns_servers": schema.SetAttribute{
				MarkdownDescription: "DNS servers to offer. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"domain_search": schema.SetAttribute{
				MarkdownDescription: "Domain search list. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the subnet.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func dhcpv6SubnetDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure DHCPv6 subnets for Kea.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the subnet.",
				Required:            true,
			},
			"subnet": dschema.StringAttribute{
				MarkdownDescription: "IPv6 Subnet.",
				Computed:            true,
			},
			"allocator": dschema.StringAttribute{
				MarkdownDescription: "Address allocator in use.",
				Computed:            true,
			},
			"pd_allocator": dschema.StringAttribute{
				MarkdownDescription: "Prefix delegation allocator in use.",
				Computed:            true,
			},
			"pools": dschema.SetAttribute{
				MarkdownDescription: "Set of address pools.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"interface": dschema.StringAttribute{
				MarkdownDescription: "Interface to listen on.",
				Computed:            true,
			},
			"dns_servers": dschema.SetAttribute{
				MarkdownDescription: "DNS servers offered.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"domain_search": dschema.SetAttribute{
				MarkdownDescription: "Domain search list.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "Optional description.",
				Computed:            true,
			},
		},
	}
}

func convertDhcpv6SubnetSchemaToStruct(d *dhcpv6SubnetResourceModel) (*kea.SubnetV6, error) {
	return &kea.SubnetV6{
		Subnet:      d.Subnet.ValueString(),
		Allocator:   api.SelectedMap(d.Allocator.ValueString()),
		PDAllocator: api.SelectedMap(d.PDAllocator.ValueString()),
		Pools:       tools.SetToString(d.Pools, "\n"),
		Interface:   api.SelectedMap(d.Interface.ValueString()),
		OptionData: kea.OptionDataV6{
			DomainNameServers: tools.SetToStringSlice(d.DnsServers),
			DomainSearch:      tools.SetToStringSlice(d.DomainSearch),
		},
		Description: d.Description.ValueString(),
	}, nil
}

func convertDhcpv6SubnetStructToSchema(d *kea.SubnetV6) (*dhcpv6SubnetResourceModel, error) {
	return &dhcpv6SubnetResourceModel{
		Subnet:       types.StringValue(d.Subnet),
		Allocator:    types.StringValue(d.Allocator.String()),
		PDAllocator:  types.StringValue(d.PDAllocator.String()),
		Pools:        tools.StringSliceToSet(strings.Split(d.Pools, "\n")),
		Interface:    types.StringValue(d.Interface.String()),
		DnsServers:   tools.StringSliceToSet(d.OptionData.DomainNameServers),
		DomainSearch: tools.StringSliceToSet(d.OptionData.DomainSearch),
		Description:  types.StringValue(d.Description),
	}, nil
}
