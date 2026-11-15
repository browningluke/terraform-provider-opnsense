package haproxy

import (
	"context"
	"errors"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/errs"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &frontendResource{}
var _ resource.ResourceWithConfigure = &frontendResource{}
var _ resource.ResourceWithImportState = &frontendResource{}

func newFrontendResource() resource.Resource {
	return &frontendResource{}
}

type frontendResource struct {
	client opnsense.Client
}

func (r *frontendResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_haproxy_frontend"
}

func (r *frontendResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = frontendResourceSchema()
}

func (r *frontendResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	apiClient, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *opnsense.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = opnsense.NewClient(apiClient)
}

func (r *frontendResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *frontendResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := convertFrontendSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse haproxy frontend, got error: %s", err))
		return
	}

	id, err := r.client.Haproxy().AddFrontend(ctx, item)
	if err != nil {
		if id != "" {
			// Created but reconfigure failed: keep it in state so it is not orphaned.
			data.Id = types.StringValue(id)
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		}
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create haproxy frontend, got error: %s", err))
		return
	}

	data.Id = types.StringValue(id)
	tflog.Trace(ctx, "created a resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *frontendResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *frontendResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.Haproxy().GetFrontend(ctx, data.Id.ValueString())
	if err != nil {
		var notFoundError *errs.NotFoundError
		if errors.As(err, &notFoundError) {
			tflog.Warn(ctx, "haproxy frontend not present in remote, removing from state")
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read haproxy frontend, got error: %s", err))
		return
	}

	frontendModel, err := convertFrontendStructToSchema(item)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read haproxy frontend, got error: %s", err))
		return
	}

	frontendModel.Id = data.Id
	resp.Diagnostics.Append(resp.State.Set(ctx, &frontendModel)...)
}

func (r *frontendResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *frontendResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := convertFrontendSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse haproxy frontend, got error: %s", err))
		return
	}

	err = r.client.Haproxy().UpdateFrontend(ctx, data.Id.ValueString(), item)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to update haproxy frontend, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *frontendResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *frontendResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Haproxy().DeleteFrontend(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to delete haproxy frontend, got error: %s", err))
		return
	}
}

func (r *frontendResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
