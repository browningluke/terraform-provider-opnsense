package service

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/bind"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-opnsense/internal/tools"
)

// BindRecordResourceModel describes the resource data model.
type BindRecordResourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled"`
	Domain  types.String `tfsdk:"domain"`
	Name    types.String `tfsdk:"name"`
	Type    types.String `tfsdk:"type"`
	Value   types.String `tfsdk:"value"`

	Id types.String `tfsdk:"id"`
}

func BindRecordResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "A single DNS entry for a domain.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this Record. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "The domain this record belongs to. Must be the UUID of the OpnSense resource.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(36, 36),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The domain this record belongs to. Must be the UUID of the OpnSense resource.",
				Optional:            true,
				Validators: []validator.String{
					// TODO check maximum length again
					// TODO test if this validator does not trigger if the value is left out
					stringvalidator.LengthBetween(1, 63),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the record.",
				Required:            true,
				Validators: []validator.String{
					// see https://github.com/opnsense/plugins/blob/7a7a3138a3170fc9a70d95ee9ae6f383eb2a0644/dns/bind/src/opnsense/mvc/app/models/OPNsense/Bind/Record.xml#L26
					stringvalidator.OneOf(
						"A",
						"AAAA",
						"CAA",
						"CNAME",
						"DNAME",
						"DNSKEY",
						"DS",
						"MX",
						"NS",
						"PTR",
						"RP",
						"RRSIG",
						"SRV",
						"SSHFP",
						"TLSA",
						"TXT",
					),
				},
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "The value of the record. The expected format depends on the record type.",
				Required:            true,
				Validators: []validator.String{
					// TODO check maximum length again
					// TODO test if this validator does not trigger if the value is left out
					stringvalidator.LengthBetween(1, 63),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func BindRecordDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Records are named lists of networks that can be used to configure who can access which resources in Bind.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable this Record.",
				Computed:            true,
			},
			"domain": dschema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"name": dschema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"type": dschema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"value": dschema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
		},
	}
}

func convertBindRecordSchemaToStruct(d *BindRecordResourceModel) (*bind.Record, error) {
	return &bind.Record{
		Enabled: tools.BoolToString(d.Enabled.ValueBool()),
		Domain:  api.SelectedMap(d.Domain.ValueString()),
		Name:    d.Name.ValueString(),
		Type:    api.SelectedMap(d.Type.ValueString()),
		Value:   d.Value.ValueString(),
	}, nil
}

func convertBindRecordStructToSchema(d *bind.Record) (*BindRecordResourceModel, error) {
	model := &BindRecordResourceModel{
		Enabled: types.BoolValue(tools.StringToBool(d.Enabled)),
		Domain:  types.StringValue(d.Domain.String()),
		Name:    types.StringValue(d.Name),
		Type:    types.StringValue(d.Type.String()),
		Value:   types.StringValue(d.Value),
	}

	return model, nil
}
