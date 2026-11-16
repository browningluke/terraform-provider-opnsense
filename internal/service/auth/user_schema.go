package auth

import (
	"strings"

	"github.com/browningluke/opnsense-go/pkg/api"
	opnauth "github.com/browningluke/opnsense-go/pkg/auth"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type userResourceModel struct {
	Enabled          types.Bool   `tfsdk:"enabled"`
	Name             types.String `tfsdk:"name"`
	FullName         types.String `tfsdk:"full_name"`
	Email            types.String `tfsdk:"email"`
	Comment          types.String `tfsdk:"comment"`
	Shell            types.String `tfsdk:"shell"`
	Groups           types.Set    `tfsdk:"groups"`
	Privileges       types.Set    `tfsdk:"privileges"`
	AuthorizedKeys   types.Set    `tfsdk:"authorized_keys"`
	Password         types.String `tfsdk:"password"`
	ScramblePassword types.Bool   `tfsdk:"scramble_password"`

	Id types.String `tfsdk:"id"`
}

func userResourceSchema() schema.Schema {
	return schema.Schema{
		Version:             1,
		MarkdownDescription: "Local OPNsense users (`System > Access > Users`), including SSH authorized keys and group memberships.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this user. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Login name (e.g. `alice`).",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"full_name": schema.StringAttribute{
				MarkdownDescription: "Display name / description shown in the UI.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "Email address for the user.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"comment": schema.StringAttribute{
				MarkdownDescription: "Optional comment.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"shell": schema.StringAttribute{
				MarkdownDescription: "Login shell. Empty means the OPNsense default (no shell except for root). Use `/bin/sh` for SSH shell access.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"groups": schema.SetAttribute{
				MarkdownDescription: "Group memberships by **name** (e.g. `admins`). Names only: the API stores GIDs but Terraform always reads names back.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"privileges": schema.SetAttribute{
				MarkdownDescription: "Direct privilege IDs (e.g. `page-all`). Prefer assigning privileges via groups when possible.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"authorized_keys": schema.SetAttribute{
				MarkdownDescription: "SSH public keys (one OpenSSH public key line per element).",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password for GUI/login. Sensitive. Do not set together with `scramble_password = true` (set `scramble_password = false` when providing a password). Not returned by the API after create.",
				Optional:            true,
				Sensitive:           true,
			},
			"scramble_password": schema.BoolAttribute{
				MarkdownDescription: "When `true` (default), OPNsense generates a random password. Prefer this for SSH-key-only users. Ignored when `password` is set (the configured password is used instead).",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the user.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func convertUserSchemaToStruct(d *userResourceModel, groupGIDs []string) (*opnauth.User, error) {
	keys := tools.SetToStringSlice(d.AuthorizedKeys)
	privs := tools.SetToStringSlice(d.Privileges)

	password := ""
	if !d.Password.IsNull() && !d.Password.IsUnknown() {
		password = d.Password.ValueString()
	}

	// Never send scramble together with an explicit password: OPNsense would
	// generate a random password and ignore the configured one.
	scramble := "0"
	if password == "" && d.ScramblePassword.ValueBool() {
		scramble = "1"
	}

	return &opnauth.User{
		Name:              d.Name.ValueString(),
		Disabled:          tools.BoolToString(!d.Enabled.ValueBool()),
		Fullname:          d.FullName.ValueString(),
		Email:             d.Email.ValueString(),
		Comment:           d.Comment.ValueString(),
		Shell:             api.SelectedMap(d.Shell.ValueString()),
		GroupMemberships:  api.SelectedMapList(groupGIDs),
		Privilege:         api.SelectedMapList(privs),
		AuthorizedKeys:    strings.Join(keys, "\n"),
		Password:          password,
		ScrambledPassword: scramble,
	}, nil
}

func convertUserStructToSchema(d *opnauth.User, groupNames []string) (*userResourceModel, error) {
	keys := []string{}
	for _, line := range strings.Split(d.AuthorizedKeys, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			keys = append(keys, line)
		}
	}

	return &userResourceModel{
		Enabled:          types.BoolValue(!tools.StringToBool(d.Disabled)),
		Name:             types.StringValue(d.Name),
		FullName:         types.StringValue(d.Fullname),
		Email:            types.StringValue(d.Email),
		Comment:          types.StringValue(d.Comment),
		Shell:            types.StringValue(d.Shell.String()),
		Groups:           tools.StringSliceToSet(groupNames),
		Privileges:       tools.StringSliceToSet([]string(d.Privilege)),
		AuthorizedKeys:   tools.StringSliceToSet(keys),
		ScramblePassword: types.BoolValue(false),
		Password:         types.StringNull(),
	}, nil
}
