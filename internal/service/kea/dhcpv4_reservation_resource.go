package kea

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

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &dhcpv4ReservationResource{}
var _ resource.ResourceWithConfigure = &dhcpv4ReservationResource{}
var _ resource.ResourceWithImportState = &dhcpv4ReservationResource{}

func newDhcpv4ReservationResource() resource.Resource {
	return &dhcpv4ReservationResource{}
}

// dhcpv4ReservationResource defines the resource implementation.
type dhcpv4ReservationResource struct {
	client opnsense.Client
}

func (r *dhcpv4ReservationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kea_dhcpv4_reservation"
}

func (r *dhcpv4ReservationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = dhcpv4ReservationResourceSchema()
}

func (r *dhcpv4ReservationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *dhcpv4ReservationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *dhcpv4ReservationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reservation, err := convertDhcpv4ReservationSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse reservation, got error: %s", err))
		return
	}

	id, err := r.client.Kea().AddReservationV4(ctx, reservation)
	if err != nil {
		if id != "" {
			data.Id = types.StringValue(id)

			// Read back so state captures API-normalised values (defaults,
			// sorting, trimming); fall back to plan-only state if the
			// read-back fails so the upstream resource isn't orphaned.
			if readStruct, readErr := r.client.Kea().GetReservationV4(ctx, id); readErr == nil {
				if readModel, convErr := convertDhcpv4ReservationStructToSchema(readStruct); convErr == nil {
					readModel.Id = data.Id
					data = readModel
				}
			}

			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create reservation, got error: %s", err))
		return
	}

	data.Id = types.StringValue(id)
	tflog.Trace(ctx, "created a resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dhcpv4ReservationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *dhcpv4ReservationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reservation, err := r.client.Kea().GetReservationV4(ctx, data.Id.ValueString())
	if err != nil {
		var notFoundError *errs.NotFoundError
		if errors.As(err, &notFoundError) {
			tflog.Warn(ctx, fmt.Sprintf("reservation not present in remote, removing from state"))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read reservation, got error: %s", err))
		return
	}

	resModel, err := convertDhcpv4ReservationStructToSchema(reservation)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read reservation, got error: %s", err))
		return
	}

	resModel.Id = data.Id
	resp.Diagnostics.Append(resp.State.Set(ctx, &resModel)...)
}

func (r *dhcpv4ReservationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *dhcpv4ReservationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := convertDhcpv4ReservationSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse reservation, got error: %s", err))
		return
	}

	err = r.client.Kea().UpdateReservationV4(ctx, data.Id.ValueString(), res)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to update reservation, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dhcpv4ReservationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *dhcpv4ReservationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Kea().DeleteReservationV4(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to delete reservation, got error: %s", err))
		return
	}
}

func (r *dhcpv4ReservationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
