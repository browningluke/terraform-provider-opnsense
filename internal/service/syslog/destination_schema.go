package syslog

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/syslog"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type destinationResourceModel struct {
	Enabled     types.Bool   `tfsdk:"enabled"`
	Transport   types.String `tfsdk:"transport"`
	Level       types.String `tfsdk:"level"`
	Facility    types.Set    `tfsdk:"facility"`
	Program     types.Set    `tfsdk:"program"`
	Hostname    types.String `tfsdk:"hostname"`
	Certificate types.String `tfsdk:"certificate"`
	Port        types.String `tfsdk:"port"`
	Rfc5424     types.Bool   `tfsdk:"rfc5424"`
	Description types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

func destinationResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Remote syslog destinations forward OPNsense logs to an external log collector (System ‣ Settings ‣ Logging ‣ Remote).",
		Version:             1,

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this remote logging target. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"transport": schema.StringAttribute{
				MarkdownDescription: "Transport used to deliver log messages. One of `udp4`, `udp6`, `tcp4`, `tcp6`, `tls4`, `tls6`. Defaults to `udp4`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("udp4"),
			},
			"level": schema.StringAttribute{
				MarkdownDescription: "Minimum log level to forward (e.g. `info`, `notice`, `warn`, `err`, `crit`). Leave `\"\"` to forward all levels. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"facility": schema.SetAttribute{
				MarkdownDescription: "Syslog facilities to forward (e.g. `local0`, `auth`, `kern`). Leave empty to forward all facilities. Defaults to `[]`.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			},
			"program": schema.SetAttribute{
				MarkdownDescription: "Applications/programs to forward (e.g. `filterlog`, `suricata`). Leave empty to forward all applications. Defaults to `[]`.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "Hostname or IP address of the remote log collector.",
				Required:            true,
			},
			"certificate": schema.StringAttribute{
				MarkdownDescription: "ID (refid) of the certificate to use for `tls4`/`tls6` transports. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"port": schema.StringAttribute{
				MarkdownDescription: "Port of the remote log collector. Defaults to `514`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("514"),
			},
			"rfc5424": schema.BoolAttribute{
				MarkdownDescription: "When enabled, uses RFC 5424 message format instead of the legacy RFC 3164 format. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description for this remote logging target.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the remote logging target.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func destinationDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Remote syslog destinations forward OPNsense logs to an external log collector (System ‣ Settings ‣ Logging ‣ Remote).",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the remote logging target.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Whether this remote logging target is enabled.",
				Computed:            true,
			},
			"transport": dschema.StringAttribute{
				MarkdownDescription: "Transport used to deliver log messages.",
				Computed:            true,
			},
			"level": dschema.StringAttribute{
				MarkdownDescription: "Minimum log level forwarded.",
				Computed:            true,
			},
			"facility": dschema.SetAttribute{
				MarkdownDescription: "Syslog facilities forwarded.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"program": dschema.SetAttribute{
				MarkdownDescription: "Applications/programs forwarded.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"hostname": dschema.StringAttribute{
				MarkdownDescription: "Hostname or IP address of the remote log collector.",
				Computed:            true,
			},
			"certificate": dschema.StringAttribute{
				MarkdownDescription: "ID (refid) of the certificate used for TLS transports.",
				Computed:            true,
			},
			"port": dschema.StringAttribute{
				MarkdownDescription: "Port of the remote log collector.",
				Computed:            true,
			},
			"rfc5424": dschema.BoolAttribute{
				MarkdownDescription: "Whether RFC 5424 message format is used.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "Description of the remote logging target.",
				Computed:            true,
			},
		},
	}
}

func convertDestinationSchemaToStruct(d *destinationResourceModel) (*syslog.Destination, error) {
	facilityList := tools.SetToStringSlice(d.Facility)
	programList := tools.SetToStringSlice(d.Program)

	return &syslog.Destination{
		Enabled:     tools.BoolToString(d.Enabled.ValueBool()),
		Transport:   api.SelectedMap(d.Transport.ValueString()),
		Level:       api.SelectedMap(d.Level.ValueString()),
		Facility:    api.SelectedMapList(facilityList),
		Program:     api.SelectedMapList(programList),
		Hostname:    d.Hostname.ValueString(),
		Certificate: api.SelectedMap(d.Certificate.ValueString()),
		Port:        d.Port.ValueString(),
		Rfc5424:     tools.BoolToString(d.Rfc5424.ValueBool()),
		Description: d.Description.ValueString(),
	}, nil
}

func convertDestinationStructToSchema(d *syslog.Destination) (*destinationResourceModel, error) {
	return &destinationResourceModel{
		Enabled:     types.BoolValue(tools.StringToBool(d.Enabled)),
		Transport:   types.StringValue(d.Transport.String()),
		Level:       types.StringValue(d.Level.String()),
		Facility:    tools.StringSliceToSet(d.Facility),
		Program:     tools.StringSliceToSet(d.Program),
		Hostname:    types.StringValue(d.Hostname),
		Certificate: types.StringValue(d.Certificate.String()),
		Port:        types.StringValue(d.Port),
		Rfc5424:     types.BoolValue(tools.StringToBool(d.Rfc5424)),
		Description: types.StringValue(d.Description),
	}, nil
}
