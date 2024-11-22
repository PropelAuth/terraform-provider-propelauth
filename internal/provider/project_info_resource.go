package provider

import (
	"context"
	"fmt"

	"terraform-provider-propelauth/internal/propelauth"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &projectInfoResource{}
var _ resource.ResourceWithConfigure = &projectInfoResource{}
var _ resource.ResourceWithImportState = &projectInfoResource{}

func NewProjectInfoResource() resource.Resource {
	return &projectInfoResource{}
}

// projectInfoResource defines the resource implementation.
type projectInfoResource struct {
	client *propelauth.PropelAuthClient
}

// projectInfoResourceModel describes the resource data model.
type projectInfoResourceModel struct {
	Name types.String `tfsdk:"name"`
}

func (r *projectInfoResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_info"
}

func (r *projectInfoResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Project Info resource. This is for configuring the basic project information in PropelAuth.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The project's name. It will be used in emails and hosted page titles.",
			},
		},
	}
}

func (r *projectInfoResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*propelauth.PropelAuthClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *propelauth.PropelAuthClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *projectInfoResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan projectInfoResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the project info
	name := plan.Name.ValueString()
	projectInfoResponse, err := r.client.UpdateProjectInfo(&name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting project info",
			"Could not set project info, unexpected error: "+err.Error(),
		)
		return
	}

	// save into the Terraform state.
	plan.Name = types.StringValue(projectInfoResponse.Name)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a propelauth_project_info resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *projectInfoResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state and read it into the model
	var state projectInfoResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// retrieve the project info from PropelAuth
	project_info, err := r.client.GetProjectInfo()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading PropelAuth Project Info",
			"Could not read PropelAuth Project Info: "+err.Error(),
		)
		return
	}

	// overwrite project info
	state.Name = types.StringValue(project_info.Name)
	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *projectInfoResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan data into the model
	var plan projectInfoResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the project info
	name := plan.Name.ValueString()
	projectInfoResponse, err := r.client.UpdateProjectInfo(&name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting project info",
			"Could not set project info, unexpected error: "+err.Error(),
		)
		return
	}

	if name != projectInfoResponse.Name {
		resp.Diagnostics.AddError(
			"Error updating project info",
			"Project name failed to update. The `name` is instead "+projectInfoResponse.Name,
		)
		return
	}

	plan.Name = types.StringValue(projectInfoResponse.Name)
	tflog.Trace(ctx, "updated a propelauth_project_info resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *projectInfoResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "deleted a propelauth_project_info resource")
}

func (r *projectInfoResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var state projectInfoResourceModel

	// retrieve the project info from PropelAuth
	project_info, err := r.client.GetProjectInfo()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing PropelAuth Project Info",
			"Could not read PropelAuth Project Info: "+err.Error(),
		)
		return
	}

	// overwrite project info
	state.Name = types.StringValue(project_info.Name)
	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
