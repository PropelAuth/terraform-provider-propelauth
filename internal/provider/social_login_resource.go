package provider

import (
	"context"
	"fmt"

	"terraform-provider-propelauth/internal/propelauth"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &socialLoginResource{}
var _ resource.ResourceWithConfigure = &socialLoginResource{}
var _ resource.ResourceWithImportState = &socialLoginResource{}

func NewSocialLoginResource() resource.Resource {
	return &socialLoginResource{}
}

// socialLoginResource defines the resource implementation.
type socialLoginResource struct {
	client *propelauth.PropelAuthClient
}

// socialLoginResourceModel describes the resource data model.
type socialLoginResourceModel struct {
	SocialProvider types.String `tfsdk:"social_provider"`
	ClientId       types.String `tfsdk:"client_id"`
	ClientSecret   types.String `tfsdk:"client_secret"`
}

func (r *socialLoginResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_social_login"
}

func (r *socialLoginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Backend API Key resource. This is for configuring the basic BE API key information in PropelAuth.",
		Attributes: map[string]schema.Attribute{
			"social_provider": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("Google", "Microsoft", "GitHub", "Slack", "LinkedIn", "Atlassian", "Apple", "Salesforce", "QuickBooks", "Xero", "Salesloft", "Outreach"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The OIDC provider for the Social Login you're configuring. This is only for internal dislay purposes." +
					"Accepted values are `Google`, `Microsoft`, `GitHub`, `Slack`, `LinkedIn`, `Atlassian`, `Apple`, " +
					"`Salesforce`, `QuickBooks`, `Xero`, `Salesloft`, and `Outreach`.",
			},
			"client_id": schema.StringAttribute{
				Required: true,
				Description: "The client ID. This is a unique identifier for the oauth client that can be retrieved from the " +
					"OIDC provider.",
			},
			"client_secret": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The client secret for the oauth client that can be retrieved from the OIDC provider.",
			},
		},
	}
}

func (r *socialLoginResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *socialLoginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan socialLoginResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// upsert the client credentials for the social login
	err := r.client.UpsertSocialLoginInfo(plan.SocialProvider.ValueString(), plan.ClientId.ValueString(), plan.ClientSecret.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating a Social Login in PropelAuth",
			"Could not upsert social login's client credentials, unexpected error: "+err.Error(),
		)
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a propelauth_social_login resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *socialLoginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state and read it into the model
	var state socialLoginResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// retrieve the social login from PropelAuth
	socialLoginInfo, err := r.client.GetSocialLoginInfo(state.SocialProvider.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading PropelAuth Social Login Info",
			"Could not read PropelAuth Social Login Info: "+err.Error(),
		)
		return
	}

	// update the state for social login client's credentials
	state.ClientId = types.StringValue(socialLoginInfo.ClientId)
	// state.ClientSecret = types.BoolValue(socialLoginInfo.ClientSecret)

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *socialLoginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan data into the model
	var plan socialLoginResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// upsert the client credentials for the social login
	err := r.client.UpsertSocialLoginInfo(plan.SocialProvider.ValueString(), plan.ClientId.ValueString(), plan.ClientSecret.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating social login",
			"Could not update the social login, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "updated a propelauth_social_login resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *socialLoginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state socialLoginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing social login
	err := r.client.DeleteSocialLogin(state.SocialProvider.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting PropelAuth Social Login",
			"Could not delete social login, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "deleted a propelauth_social_login resource")
}

func (r *socialLoginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("social_provider"), req, resp)
}
