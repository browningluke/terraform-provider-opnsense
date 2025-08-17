package service

import (
	"github.com/browningluke/opnsense-go/pkg/bind"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-opnsense/internal/tools"
)

// BindPrimaryDomainResourceModel describes the resource data model.
type BindPrimaryDomainResourceModel struct {
	Enabled       types.Bool   `tfsdk:"enabled"`
	DomainName    types.String `tfsdk:"domain_name"`
	AllowTransfer types.Set    `tfsdk:"allow_transfer"`
	AllowQuery    types.Set    `tfsdk:"allow_query"`
	TimeToLive    types.Int64  `tfsdk:"ttl"`
	Refresh       types.Int64  `tfsdk:"refresh"`
	Retry         types.Int64  `tfsdk:"retry"`
	Expire        types.Int64  `tfsdk:"expire"`
	Negative      types.Int64  `tfsdk:"negative"`
	MailAdmin     types.String `tfsdk:"mail_admin"`
	DnsServer     types.String `tfsdk:"dns_server"`

	Id types.String `tfsdk:"id"`
}

func BindPrimaryDomainResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Primary domains are domains fully managed by Bind.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this Primary Domain. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"domain_name": schema.StringAttribute{
				MarkdownDescription: "The domain name.",
				Required:            true,
				Validators: []validator.String{
					// TODO check maximum string length again
					stringvalidator.LengthBetween(3, 63),
				},
			},
			"allow_transfer": schema.SetAttribute{
				MarkdownDescription: "The list of ACLs to allow transfer to.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"allow_query": schema.SetAttribute{
				MarkdownDescription: "The list of ACLs to allow queries from.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"ttl": schema.Int64Attribute{
				MarkdownDescription: "Time to live for the entry (the time it may be cached). Time in seconds.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(86400),
			},
			"refresh": schema.Int64Attribute{
				MarkdownDescription: "How often a name server should check it's primary server to see if there has been any updates to the zone which it does by comparing Serial numbers.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(21600),
			},
			"retry": schema.Int64Attribute{
				MarkdownDescription: "How long a name server should wait to retry an attempt to get fresh zone data from the primary name server if the first attempt failed. Time in seconds.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(3600),
			},
			"expire": schema.Int64Attribute{
				MarkdownDescription: "Maximum time a name server will still consider itself Authoritative if it hasn't been able to refresh the zone data. Time in seconds.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(3542400),
			},
			"negative": schema.Int64Attribute{
				MarkdownDescription: "Controls negative caching time, which is how long a resolver will cache a NXDOMAIN Name Error. Time in seconds.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(3600),
			},
			"mail_admin": schema.StringAttribute{
				MarkdownDescription: "The mail admin email address. '@' must be replaced by '.'.",
				Required:            true,
			},
			"dns_server": schema.StringAttribute{
				MarkdownDescription: "The DNS server to use.",
				Required:            true,
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

func BindPrimaryDomainDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Primary Domains are named lists of networks that can be used to configure who can access which resources in Bind.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable this Primary Domain.",
				Computed:            true,
			},
			"domain_name": dschema.StringAttribute{
				MarkdownDescription: "The domain name.",
				Computed:            true,
			},
			"allow_transfer": dschema.SetAttribute{
				MarkdownDescription: "The list of ACLs to allow transfer to.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"allow_query": dschema.SetAttribute{
				MarkdownDescription: "The list of ACLs to allow queries from.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"ttl": dschema.Int64Attribute{
				MarkdownDescription: "Time to live for the entry (the time it may be cached). Time in seconds.",
				Computed:            true,
			},
			"refresh": dschema.Int64Attribute{
				MarkdownDescription: "How often a name server should check it's primary server to see if there have been any updates to the zone. This is done by comparing the (auto-generated) serial numbers. Time in seconds.",
				Computed:            true,
			},
			"retry": dschema.Int64Attribute{
				MarkdownDescription: "How long a name server should wait to retry an attempt to get fresh zone data from the primary name server if the first attempt failed. Time in seconds.",
				Computed:            true,
			},
			"expire": dschema.Int64Attribute{
				MarkdownDescription: "Maximum time a name server will still consider itself Authoritative if it hasn't been able to refresh the zone data. Time in seconds.",
				Computed:            true,
			},
			"negative": dschema.Int64Attribute{
				MarkdownDescription: "Controls negative caching time, which is how long a resolver will cache a NXDOMAIN Name Error. Time in seconds.",
				Computed:            true,
			},
			"mail_admin": dschema.StringAttribute{
				MarkdownDescription: "The mail admin email address. '@' must be replaced by '.'.",
				Computed:            true,
			},
			"dns_server": dschema.StringAttribute{
				MarkdownDescription: "The DNS server to use.",
				Computed:            true,
			},
		},
	}
}

func convertBindPrimaryDomainSchemaToStruct(d *BindPrimaryDomainResourceModel) (*bind.PrimaryDomain, error) {
	return &bind.PrimaryDomain{
		Enabled:       tools.BoolToString(d.Enabled.ValueBool()),
		DomainName:    d.DomainName.ValueString(),
		AllowTransfer: tools.SetToStringSlice(d.AllowTransfer),
		AllowQuery:    tools.SetToStringSlice(d.AllowQuery),
		TimeToLive:    tools.Int64ToString(d.TimeToLive.ValueInt64()),
		Refresh:       tools.Int64ToString(d.Refresh.ValueInt64()),
		Retry:         tools.Int64ToString(d.Retry.ValueInt64()),
		Expire:        tools.Int64ToString(d.Expire.ValueInt64()),
		Negative:      tools.Int64ToString(d.Negative.ValueInt64()),
		MailAdmin:     d.MailAdmin.ValueString(),
		DnsServer:     d.DnsServer.ValueString(),
	}, nil
}

func convertBindPrimaryDomainStructToSchema(d *bind.PrimaryDomain) (*BindPrimaryDomainResourceModel, error) {
	model := &BindPrimaryDomainResourceModel{
		Enabled:       types.BoolValue(tools.StringToBool(d.Enabled)),
		DomainName:    types.StringValue(d.DomainName),
		AllowTransfer: tools.StringSliceToSet(d.AllowTransfer),
		AllowQuery:    tools.StringSliceToSet(d.AllowQuery),
		TimeToLive:    types.Int64Value(tools.StringToInt64(d.TimeToLive)),
		Refresh:       types.Int64Value(tools.StringToInt64(d.Refresh)),
		Retry:         types.Int64Value(tools.StringToInt64(d.Retry)),
		Expire:        types.Int64Value(tools.StringToInt64(d.Expire)),
		Negative:      types.Int64Value(tools.StringToInt64(d.Negative)),
		MailAdmin:     types.StringValue(d.MailAdmin),
		DnsServer:     types.StringValue(d.DnsServer),
	}

	return model, nil
}
