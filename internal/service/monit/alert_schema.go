package monit

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/monit"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type alertResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	Recipient   types.String `tfsdk:"recipient"`
	NotOn       types.Bool   `tfsdk:"not_on"`
	Events      types.Set    `tfsdk:"events"`
	Format      types.String `tfsdk:"format"`
	Reminder    types.String `tfsdk:"reminder"`
	Description types.String `tfsdk:"description"`
}

func alertResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages a Monit alert recipient.",

		Version: 1,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the alert.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this alert. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"recipient": schema.StringAttribute{
				MarkdownDescription: "Email address to send alerts to.",
				Required:            true,
			},
			"not_on": schema.BoolAttribute{
				MarkdownDescription: "When enabled, negate the selected `events` (alert on all events except the selected ones). Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"events": schema.SetAttribute{
				MarkdownDescription: "Set of events to alert on, e.g. `action`, `checksum`, `connection`, `content`, `data`, `exec`, `fsflag`, `gid`, `icmp`, `instance`, `invalid`, `link`, `nonexist`, `permission`, `pid`, `ppid`, `resource`, `saturation`, `size`, `status`, `timeout`, `timestamp`, `uid`, `uptime`. When empty, alerts on all events.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"format": schema.StringAttribute{
				MarkdownDescription: "Custom message template for the alert email.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"reminder": schema.StringAttribute{
				MarkdownDescription: "Resend the alert every N cycles while the condition persists. `0` disables reminders. Defaults to `0`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("0"),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description for this alert.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
		},
	}
}

func alertDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Retrieves a Monit alert recipient by UUID.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "UUID of the alert.",
			},
			"enabled": dschema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether this alert is enabled.",
			},
			"recipient": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Email address alerts are sent to.",
			},
			"not_on": dschema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the selected events are negated.",
			},
			"events": dschema.SetAttribute{
				Computed:            true,
				MarkdownDescription: "Set of events this alert triggers on.",
				ElementType:         types.StringType,
			},
			"format": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Custom message template for the alert email.",
			},
			"reminder": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Reminder interval, in cycles.",
			},
			"description": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Description of the alert.",
			},
		},
	}
}

func convertAlertSchemaToStruct(d *alertResourceModel) (*monit.Alert, error) {
	return &monit.Alert{
		Enabled:     tools.BoolToString(d.Enabled.ValueBool()),
		Recipient:   d.Recipient.ValueString(),
		Noton:       tools.BoolToString(d.NotOn.ValueBool()),
		Events:      api.SelectedMapList(tools.SetToStringSlice(d.Events)),
		Format:      d.Format.ValueString(),
		Reminder:    d.Reminder.ValueString(),
		Description: d.Description.ValueString(),
	}, nil
}

func convertAlertStructToSchema(d *monit.Alert) (*alertResourceModel, error) {
	return &alertResourceModel{
		Enabled:     types.BoolValue(tools.StringToBool(d.Enabled)),
		Recipient:   types.StringValue(d.Recipient),
		NotOn:       types.BoolValue(tools.StringToBool(d.Noton)),
		Events:      tools.StringSliceToSet(d.Events),
		Format:      types.StringValue(d.Format),
		Reminder:    types.StringValue(d.Reminder),
		Description: types.StringValue(d.Description),
	}, nil
}
