package cron

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/cron"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type jobResourceModel struct {
	Enabled     types.Bool   `tfsdk:"enabled"`
	Minutes     types.String `tfsdk:"minutes"`
	Hours       types.String `tfsdk:"hours"`
	Days        types.String `tfsdk:"days"`
	Months      types.String `tfsdk:"months"`
	Weekdays    types.String `tfsdk:"weekdays"`
	Who         types.String `tfsdk:"who"`
	Command     types.String `tfsdk:"command"`
	Parameters  types.String `tfsdk:"parameters"`
	Description types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

func jobResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Cron jobs allow scheduling commands to run at specified times.",
		Version:             1,

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this cron job. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"minutes": schema.StringAttribute{
				MarkdownDescription: "Minute(s) to run (cron expression, e.g. `0` or `*/15`).",
				Required:            true,
			},
			"hours": schema.StringAttribute{
				MarkdownDescription: "Hour(s) to run (cron expression, e.g. `0` or `*/6`).",
				Required:            true,
			},
			"days": schema.StringAttribute{
				MarkdownDescription: "Day(s) of month to run (cron expression, e.g. `*` or `1`).",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("*"),
			},
			"months": schema.StringAttribute{
				MarkdownDescription: "Month(s) to run (cron expression, e.g. `*` or `1,6`).",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("*"),
			},
			"weekdays": schema.StringAttribute{
				MarkdownDescription: "Day(s) of week to run (cron expression, e.g. `*` or `1-5`).",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("*"),
			},
			"who": schema.StringAttribute{
				MarkdownDescription: "User to run the command as. Defaults to `root`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("root"),
			},
			"command": schema.StringAttribute{
				MarkdownDescription: "Command to execute. Must be a command key known to OPNsense (e.g. `firmware poll`, `unbound restart`).",
				Required:            true,
			},
			"parameters": schema.StringAttribute{
				MarkdownDescription: "Optional parameters to pass to the command.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description for this cron job.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the cron job.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func jobDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Cron jobs allow scheduling commands to run at specified times.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the cron job.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Whether this cron job is enabled.",
				Computed:            true,
			},
			"minutes": dschema.StringAttribute{
				MarkdownDescription: "Minute(s) to run.",
				Computed:            true,
			},
			"hours": dschema.StringAttribute{
				MarkdownDescription: "Hour(s) to run.",
				Computed:            true,
			},
			"days": dschema.StringAttribute{
				MarkdownDescription: "Day(s) of month to run.",
				Computed:            true,
			},
			"months": dschema.StringAttribute{
				MarkdownDescription: "Month(s) to run.",
				Computed:            true,
			},
			"weekdays": dschema.StringAttribute{
				MarkdownDescription: "Day(s) of week to run.",
				Computed:            true,
			},
			"who": dschema.StringAttribute{
				MarkdownDescription: "User the command runs as.",
				Computed:            true,
			},
			"command": dschema.StringAttribute{
				MarkdownDescription: "Command key.",
				Computed:            true,
			},
			"parameters": dschema.StringAttribute{
				MarkdownDescription: "Parameters passed to the command.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "Description of the cron job.",
				Computed:            true,
			},
		},
	}
}

func convertJobSchemaToStruct(d *jobResourceModel) (*cron.Job, error) {
	return &cron.Job{
		Enabled:     tools.BoolToString(d.Enabled.ValueBool()),
		Minutes:     d.Minutes.ValueString(),
		Hours:       d.Hours.ValueString(),
		Days:        d.Days.ValueString(),
		Months:      d.Months.ValueString(),
		Weekdays:    d.Weekdays.ValueString(),
		Who:         d.Who.ValueString(),
		Command:     api.SelectedMap(d.Command.ValueString()),
		Parameters:  d.Parameters.ValueString(),
		Description: d.Description.ValueString(),
	}, nil
}

func convertJobStructToSchema(d *cron.Job) (*jobResourceModel, error) {
	return &jobResourceModel{
		Enabled:     types.BoolValue(tools.StringToBool(d.Enabled)),
		Minutes:     types.StringValue(d.Minutes),
		Hours:       types.StringValue(d.Hours),
		Days:        types.StringValue(d.Days),
		Months:      types.StringValue(d.Months),
		Weekdays:    types.StringValue(d.Weekdays),
		Who:         types.StringValue(d.Who),
		Command:     types.StringValue(d.Command.String()),
		Parameters:  tools.StringOrNull(d.Parameters),
		Description: types.StringValue(d.Description),
	}, nil
}
