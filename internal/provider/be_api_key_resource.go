package provider

import (
	"context"
	"fmt"

	"terraform-provider-propelauth/internal/propelauth"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &beApiKeyResource{}
var _ resource.ResourceWithConfigure = &beApiKeyResource{}

func NewBeApiKeyResource() resource.Resource {
	return &beApiKeyResource{}
}

// beApiKeyResource defines the resource implementation.
type beApiKeyResource struct {
	client *propelauth.PropelAuthClient
}

// beApiKeyResourceModel describes the resource data model.
type beApiKeyResourceModel struct {
	Environment types.String `tfsdk:"environment"`
	Name        types.String `tfsdk:"name"`
	ReadOnly    types.Bool   `tfsdk:"read_only"`
	ApiKey      types.String `tfsdk:"api_key"`
	ApiKeyId    types.String `tfsdk:"api_key_id"`
}

func (r *beApiKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_be_api_key"
}

func (r *beApiKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Backend API Key resource. This is for configuring the basic BE API key information in PropelAuth.",
		Attributes: map[string]schema.Attribute{
			"environment": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("Test", "Staging", "Prod"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The environment for which you are configuring the backend integration. Accepted values are `Test`, `Staging`, and `Prod`.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The API key's name. This is only for internal dislay purposes.",
			},
			"read_only": schema.BoolAttribute{
				Required: true,
				Description: "If true, the API key has read-only privileges. For example, it cannot be used for " +
					"creating, editing, or deleting users/orgs. This value can only be set during the creation of the API key.",
			},
			"api_key": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The API key value. This is the secret value that is used to authenticate requests to PropelAuth.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"api_key_id": schema.StringAttribute{
				Computed:    true,
				Description: "The API key ID. This is a unique identifier for the API key.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *beApiKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *beApiKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan beApiKeyResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create the be api key
	beApiKeyInfo, err := r.client.CreateBeApiKey(plan.Environment.ValueString(), plan.Name.ValueString(), plan.ReadOnly.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating PropelAuth Backend API Key",
			"Could not create be api key, unexpected error: "+err.Error(),
		)
		return
	}

	// set the Terraform state.
	plan.Name = types.StringValue(beApiKeyInfo.Name)
	plan.ApiKey = types.StringValue(beApiKeyInfo.ApiKey)
	plan.ApiKeyId = types.StringValue(beApiKeyInfo.ApiKeyId)
	plan.ReadOnly = types.BoolValue(beApiKeyInfo.IsReadOnly)


	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a propelauth_be_api_key resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *beApiKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state and read it into the model
	var state beApiKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// retrieve the be api key from PropelAuth
	beApiKeyInfo, err := r.client.GetBeApiKeyInfo(state.Environment.ValueString(), state.ApiKeyId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading PropelAuth Backend API Key",
			"Could not read PropelAuth Backend API Key: "+err.Error(),
		)
		return
	}

	// update the state for the be api key
	state.Name = types.StringValue(beApiKeyInfo.Name)
	state.ApiKey = types.StringValue(beApiKeyInfo.ApiKey)
	state.ApiKeyId = types.StringValue(beApiKeyInfo.ApiKeyId)
	state.ReadOnly = types.BoolValue(beApiKeyInfo.IsReadOnly)

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *beApiKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan data into the model
	var plan beApiKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the be api key
	beApiKeyResponse, err := r.client.UpdateBeApiKey(plan.Environment.ValueString(), plan.ApiKeyId.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating be api key",
			"Could not update the be api key, unexpected error: "+err.Error(),
		)
		return
	}

	// Save updated state from the response into Terraform state
	plan.Name = types.StringValue(beApiKeyResponse.Name)

	tflog.Trace(ctx, "updated a propelauth_be_api_key resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *beApiKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state beApiKeyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DeleteBeApiKey(state.Environment.ValueString(), state.ApiKeyId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting PropelAuth Backend API Key",
			"Could not delete be api key, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "deleted a propelauth_be_api_key resource")
}
