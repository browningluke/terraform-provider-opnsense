package haproxy

import (
	"context"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/haproxy"
	"github.com/browningluke/terraform-provider-opnsense/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type frontendActionsResourceModel struct {
	FrontendId types.String `tfsdk:"frontend_id"`
	ActionIds  types.List   `tfsdk:"action_ids"`

	Id types.String `tfsdk:"id"`
}

func frontendActionsResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "The ordered rule list (`linkedActions`) of an existing HAProxy frontend. " +
			"This resource owns only that list. The frontend itself is managed with `opnsense_haproxy_frontend`, which never touches the list, or outside Terraform. " +
			"The list is authoritative, so every rule the frontend should run must be listed, including rules other tools add " +
			"(the ACME client plugin appends its `redirect_acme_challenges` rule if it is missing). " +
			"Destroying this resource empties the frontend's rule list. Requires the `os-haproxy` plugin.",
		Version: 1,

		Attributes: map[string]schema.Attribute{
			"frontend_id": schema.StringAttribute{
				MarkdownDescription: "UUID of the frontend (OPNsense: Public Service) whose rules are managed.",
				Required:            true,
				Validators:          []validator.String{validators.IsUUIDv4()},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"action_ids": schema.ListAttribute{
				MarkdownDescription: "UUIDs of the rules (`opnsense_haproxy_action`) in evaluation order. HAProxy uses the first `use_backend` rule that matches.",
				ElementType:         types.StringType,
				Required:            true,
				Validators: []validator.List{
					listvalidator.UniqueValues(),
					listvalidator.ValueStringsAre(validators.IsUUIDv4()),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Same as `frontend_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func convertFrontendActionsSchemaToStruct(ctx context.Context, d *frontendActionsResourceModel) (*haproxy.FrontendActions, diag.Diagnostics) {
	ids := api.SelectedMapListOrdered{}
	diags := d.ActionIds.ElementsAs(ctx, &ids, false)
	return &haproxy.FrontendActions{LinkedActions: ids}, diags
}

func convertFrontendActionsStructToSchema(ctx context.Context, id string, d *haproxy.FrontendActions) (*frontendActionsResourceModel, diag.Diagnostics) {
	ids, diags := types.ListValueFrom(ctx, types.StringType, []string(d.LinkedActions))
	return &frontendActionsResourceModel{
		FrontendId: types.StringValue(id),
		ActionIds:  ids,
		Id:         types.StringValue(id),
	}, diags
}
