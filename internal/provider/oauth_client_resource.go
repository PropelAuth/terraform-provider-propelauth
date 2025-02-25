package provider

import (
	"context"
	"fmt"

	"terraform-provider-propelauth/internal/propelauth"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &oauthClientResource{}
var _ resource.ResourceWithConfigure = &oauthClientResource{}

func NewOauthClientResource() resource.Resource {
	return &oauthClientResource{}
}

// oauthClientResource defines the resource implementation.
type oauthClientResource struct {
	client *propelauth.PropelAuthClient
}

// oauthClientResourceModel describes the resource data model.
type oauthClientResourceModel struct {
	Environment  types.String `tfsdk:"environment"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	RedirectUris types.List   `tfsdk:"redirect_uris"`
}

func (r *oauthClientResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_client"
}

func (r *oauthClientResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Oauth Client Resource. This is for configuring the basic BE API key information in PropelAuth.",
		Attributes: map[string]schema.Attribute{
			"environment": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("Test", "Staging", "Prod"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The environment for which you are integrating. Accepted values are `Test`, `Staging`, and `Prod`.",
			},
			"client_id": schema.StringAttribute{
				Computed:    true,
				Description: "The client ID set by PropelAuth.",
			},
			"client_secret": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The client secret set by PropelAuth.",
			},
			"redirect_uris": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "A list of redirect URIs that are whitelisted for this client. Must be a valid URL including a " +
					"scheme and hostname. You may specify a wildcard (*) in the hostname to allow any subdomain.",
			},
		},
	}
}

func (r *oauthClientResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *oauthClientResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan oauthClientResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// convert the redirect URIs to a string slice
	redirectUris := make([]types.String, 0, len(plan.RedirectUris.Elements()))
	diags = plan.RedirectUris.ElementsAs(ctx, &redirectUris, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	convertedRedirectUris := make([]string, len(redirectUris))
	for i, uri := range redirectUris {
		convertedRedirectUris[i] = uri.ValueString()
	}

	// create the oauth client
	oauthClientInfo, err := r.client.CreateOauthClient(plan.Environment.ValueString(), convertedRedirectUris)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating PropelAuth Oauth Client",
			"Could not create oauth client, unexpected error: "+err.Error(),
		)
		return
	}

	// set the Terraform state.
	plan.ClientId = types.StringValue(oauthClientInfo.ClientId)
	plan.ClientSecret = types.StringValue(oauthClientInfo.ClientSecret)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a propelauth_oauth_client resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *oauthClientResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state and read it into the model
	var state oauthClientResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// retrieve the oauth client from PropelAuth
	oauthClientInfo, err := r.client.GetOauthClientInfo(state.Environment.ValueString(), state.ClientId.ValueString())
	if err != nil {
		// If error is "not_found", it indicates that the resource should be deleted.
		if propelauth.IsPropelAuthNotFoundError(err) {
			tflog.Trace(ctx, "deleting a propelauth_oauth_client resource because it was not found in PropelAuth")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading PropelAuth Oauth Client",
			"Could not read PropelAuth Oauth Client: "+err.Error(),
		)
		return
	}

	// update the state for the oauth client
	redirectUris := make([]attr.Value, len(oauthClientInfo.RedirectUris))
	for i, uri := range oauthClientInfo.RedirectUris {
		redirectUris[i] = types.StringValue(uri)
	}
	convertedRedirectUris, diags := types.ListValue(types.StringType, redirectUris)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.RedirectUris = convertedRedirectUris

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *oauthClientResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan data into the model
	var plan oauthClientResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// convert the redirect URIs to a string slice
	redirectUris := make([]types.String, 0, len(plan.RedirectUris.Elements()))
	diags := plan.RedirectUris.ElementsAs(ctx, &redirectUris, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	convertedRedirectUris := make([]string, len(redirectUris))
	for i, uri := range redirectUris {
		convertedRedirectUris[i] = uri.ValueString()
	}

	// Update the oauth client
	err := r.client.UpdateOauthClient(plan.Environment.ValueString(), plan.ClientId.ValueString(), convertedRedirectUris)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating oauth client",
			"Could not update the oauth client, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "updated a propelauth_oauth_client resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *oauthClientResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state oauthClientResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing oauth client
	err := r.client.DeleteOauthClient(state.Environment.ValueString(), state.ClientId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting PropelAuth Oauth Client",
			"Could not delete oauth client, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "deleted a propelauth_oauth_client resource")
}
