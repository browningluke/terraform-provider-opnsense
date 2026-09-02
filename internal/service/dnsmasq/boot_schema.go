package dnsmasq

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/dnsmasq"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type bootResourceModel struct {
	Interface     types.String `tfsdk:"interface"`
	Tag           types.Set    `tfsdk:"tag"`
	FileName      types.String `tfsdk:"file_name"`
	ServerName    types.String `tfsdk:"server_name"`
	ServerAddress types.String `tfsdk:"server_address"`
	Description   types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

func bootResourceSchema() schema.Schema {
	return schema.Schema{
		Version:             1,
		MarkdownDescription: "Configure DHCP boot information for dnsmasq.",
		Attributes: map[string]schema.Attribute{
			"interface": schema.StringAttribute{
				MarkdownDescription: "This adds a single interface as tag so this DHCP boot option can match the interface of a DHCP range.",
				Required:            true,
			},
			"tag": schema.SetAttribute{
				MarkdownDescription: "If the optional tags are given then this option is only sent when all the tags are matched. Can be optionally combined with an interface tag.",
				Required:            false,
				Optional:            true,
				ElementType:         types.StringType,
			},
			"file_name": schema.StringAttribute{
				MarkdownDescription: "The filename to send for the client.",
				Required:            true,
			},
			"server_name": schema.StringAttribute{
				MarkdownDescription: "The server name to send for the client.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"server_address": schema.StringAttribute{
				MarkdownDescription: "The sever address to send for the client.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "You may enter a description here for your reference (not parsed).",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "UUID of the boot object.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func bootDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure DHCP boot information for dnsmasq.",
		Attributes: map[string]dschema.Attribute{
			"interface": schema.StringAttribute{
				MarkdownDescription: "This adds a single interface as tag so this DHCP boot option can match the interface of a DHCP range.",
				Computed:            true,
			},
			"tag": schema.SetAttribute{
				MarkdownDescription: "If the optional tags are given then this option is only sent when all the tags are matched. Can be optionally combined with an interface tag.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"file_name": schema.StringAttribute{
				MarkdownDescription: "The filename to send for the client.",
				Computed:            true,
			},
			"server_name": schema.StringAttribute{
				MarkdownDescription: "The server name to send for the client.",
				Computed:            true,
			},
			"server_address": schema.StringAttribute{
				MarkdownDescription: "The sever address to send for the client.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "You may enter a description here for your reference (not parsed).",
				Computed:            true,
			},
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the boot object.",
				Required:            true,
			},
		},
	}
}

func convertBootSchemaToStruct(d *bootResourceModel) (*dnsmasq.Boot, error) {
	boot := &dnsmasq.Boot{
		Interface:     api.SelectedMap(d.Interface.ValueString()),
		Filename:      d.FileName.ValueString(),
		Servername:    d.ServerName.ValueString(),
		ServerAddress: d.ServerAddress.ValueString(),
		Description:   d.Description.ValueString(),
	}
	if !d.Tag.IsNull() {
		var tags []string

		diags := d.Tag.ElementsAs(
			context.Background(),
			&tags,
			false,
		)

		if diags.HasError() {
			return nil, fmt.Errorf("failed to parse tags")
		}

		boot.Tag = api.SelectedMapList(tags)
	}
	return boot, nil
}

func convertBootStructToSchema(d *dnsmasq.Boot) (*bootResourceModel, error) {
	tagSet := types.SetNull(types.StringType)
	if len(d.Tag) > 0 {
		convertedTagSet, diags := types.SetValueFrom(
			context.Background(),
			types.StringType,
			d.Tag,
		)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert tags")
		}
		tagSet = convertedTagSet
	}

	boot := &bootResourceModel{
		FileName:      types.StringValue(d.Filename),
		ServerName:    types.StringValue(d.Servername),
		ServerAddress: types.StringValue(d.ServerAddress),
		Description:   types.StringValue(d.Description),
	}

	if d.Interface.String() == "" {
		boot.Interface = types.StringNull()
	} else {
		boot.Interface = types.StringValue(d.Interface.String())
	}
	boot.Tag = tagSet

	return boot, nil
}
