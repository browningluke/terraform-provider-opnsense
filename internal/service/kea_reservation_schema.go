package service

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

// KeaReservationResourceModel describes the resource data model.
type KeaReservationResourceModel struct {
	SubnetId types.String `tfsdk:"subnet_id"`

	IpAddress  types.String `tfsdk:"ip_address"`
	MacAddress types.String `tfsdk:"mac_address"`
	Hostname   types.String `tfsdk:"hostname"`

	Description types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

func KeaReservationResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Configure DHCP reservations for Kea.",

		Attributes: map[string]schema.Attribute{
			"subnet_id": schema.StringAttribute{
				MarkdownDescription: "Subnet ID the reservation belongs to.",
				Required:            true,
			},
			"ip_address": schema.StringAttribute{
				MarkdownDescription: "IP address to offer to the client.",
				Required:            true,
			},
			"mac_address": schema.StringAttribute{
				MarkdownDescription: "MAC/Ether address of the client in question.",
				Required:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "Hostname to offer to the client. Defaults to `\"\"`..",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the reservation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func KeaReservationDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure DHCP reservations for Kea.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the reservation.",
				Required:            true,
			},
			"subnet_id": schema.StringAttribute{
				MarkdownDescription: "Subnet ID the reservation belongs to.",
				Computed:            true,
			},
			"ip_address": schema.StringAttribute{
				MarkdownDescription: "IP address to offer to the client.",
				Computed:            true,
			},
			"mac_address": schema.StringAttribute{
				MarkdownDescription: "MAC/Ether address of the client in question.",
				Computed:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "Hostname to offer to the client.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Computed:            true,
			},
		},
	}
}

func convertKeaReservationSchemaToStruct(d *KeaReservationResourceModel) (*kea.Reservation, error) {
	return &kea.Reservation{
		Subnet:      api.SelectedMap(d.SubnetId.ValueString()),
		IpAddress:   d.IpAddress.ValueString(),
		HwAddress:   d.MacAddress.ValueString(),
		Hostname:    d.Hostname.ValueString(),
		Description: d.Description.ValueString(),
	}, nil
}

func convertKeaReservationStructToSchema(d *kea.Reservation) (*KeaReservationResourceModel, error) {
	return &KeaReservationResourceModel{
		SubnetId:    types.StringValue(d.Subnet.String()),
		IpAddress:   types.StringValue(d.IpAddress),
		MacAddress:  types.StringValue(d.HwAddress),
		Hostname:    types.StringValue(d.Hostname),
		Description: types.StringValue(d.Description),
	}, nil
}
