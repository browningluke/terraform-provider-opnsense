package unbound

import (
	"context"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/unbound"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/browningluke/terraform-provider-opnsense/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// aclResourceModel describes the resource data model.
type aclResourceModel struct {
	Enabled     types.Bool   `tfsdk:"enabled"`
	Name        types.String `tfsdk:"name"`
	Action      types.String `tfsdk:"action"`
	Networks    types.Set    `tfsdk:"networks"`
	Description types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

func aclResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Access Control List entries control which source networks are permitted to query the Unbound resolver.",
		Version:             1,

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "When enabled, this ACL entry is active. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Descriptive name for this ACL entry.",
				Required:            true,
			},
			"action": schema.StringAttribute{
				MarkdownDescription: "Action to take for queries from the listed networks. One of: `allow`, `deny`, `refuse`, `allow_snoop`, `deny_non_local`, `refuse_non_local`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"allow",
						"deny",
						"refuse",
						"allow_snoop",
						"deny_non_local",
						"refuse_non_local",
					),
				},
			},
			"networks": schema.SetAttribute{
				MarkdownDescription: "One or more CIDR blocks to match (e.g. `10.0.0.0/24`).",
				Required:            true,
				ElementType:         types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(validators.CIDR()),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the ACL entry.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func aclDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Access Control List entries control which source networks are permitted to query the Unbound resolver.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Whether this ACL entry is enabled.",
				Computed:            true,
			},
			"name": dschema.StringAttribute{
				MarkdownDescription: "Descriptive name for this ACL entry.",
				Computed:            true,
			},
			"action": dschema.StringAttribute{
				MarkdownDescription: "Action to take for queries from the listed networks.",
				Computed:            true,
			},
			"networks": dschema.SetAttribute{
				MarkdownDescription: "One or more CIDR blocks to match.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Computed:            true,
			},
		},
	}
}

func convertAclSchemaToStruct(d *aclResourceModel) (*unbound.Acl, error) {
	var networksList []string
	d.Networks.ElementsAs(context.Background(), &networksList, false)

	return &unbound.Acl{
		Enabled:     tools.BoolToString(d.Enabled.ValueBool()),
		Name:        d.Name.ValueString(),
		Action:      api.SelectedMap(d.Action.ValueString()),
		Networks:    api.SelectedMapList(networksList),
		Description: d.Description.ValueString(),
	}, nil
}

func convertAclStructToSchema(d *unbound.Acl) (*aclResourceModel, error) {
	return &aclResourceModel{
		Enabled:     types.BoolValue(tools.StringToBool(d.Enabled)),
		Name:        types.StringValue(d.Name),
		Action:      types.StringValue(d.Action.String()),
		Networks:    tools.StringSliceToSet(d.Networks),
		Description: tools.StringOrNull(d.Description),
	}, nil
}
