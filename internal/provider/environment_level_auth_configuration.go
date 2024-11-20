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
var _ resource.Resource = &environmentLevelAuthConfigurationResource{}
var _ resource.ResourceWithConfigure = &environmentLevelAuthConfigurationResource{}

func NewEnvironmentLevelAuthConfigurationResource() resource.Resource {
	return &environmentLevelAuthConfigurationResource{}
}

// environmentLevelAuthConfigurationResource defines the resource implementation.
type environmentLevelAuthConfigurationResource struct {
	client *propelauth.PropelAuthClient
}

// environmentLevelAuthConfigurationResourceModel describes the resource data model.
type environmentLevelAuthConfigurationResourceModel struct {
	Environment                           types.String `tfsdk:"environment"`
	AllowPublicSignups                    types.Bool   `tfsdk:"allow_public_signups"`
	RequireEmailConfirmation              types.Bool   `tfsdk:"require_email_confirmation"`
	WaitlistUsersRequireEmailConfirmation types.Bool   `tfsdk:"waitlist_users_require_email_confirmation"`
}

func (r *environmentLevelAuthConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment_level_auth_configuration"
}

func (r *environmentLevelAuthConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Environment-level Auth Configuration. This is for configuring elements of the signup and login experience " +
			"in PropelAuth that you may want to differ between test and production environments.",
		Attributes: map[string]schema.Attribute{
			"environment": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("Test", "Staging", "Prod"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The environment for which you are configuring the login and signup experience. " +
					"Accepted values are `Test`, `Staging`, and `Prod`.",
			},
			"allow_public_signups": schema.BoolAttribute{
				Optional: true,
				Description: "If true, new users will be able to sign up for your product directly in the PropelAuth hosted pages." +
					"The default setting is true for all environments.",
			},
			"require_email_confirmation": schema.BoolAttribute{
				Optional: true,
				Description: "If true, all users are required to have confirmed email addresses. Whenever PropelAuth doesn't know for " +
					"certain whether a user's email adderss is in fact owned by them, PropelAuth will trigger an email confirmation flow. " +
					"The default setting is true for `Prod` and `Staging` environments but is false for `Test` for ease of development.",
			},
			"waitlist_users_require_email_confirmation": schema.BoolAttribute{
				Optional: true,
				Description: "If true, all waitlisted users are required to have confirmed email addresses. Whenever PropelAuth doesn't know for " +
					"certain whether a waitlisted user's email adderss is in fact owned by them, PropelAuth will trigger an email confirmation flow. " +
					"The default setting is false for all environments.",
			},
		},
	}
}

func (r *environmentLevelAuthConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *environmentLevelAuthConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan environmentLevelAuthConfigurationResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the configuration in PropelAuth
	environment := plan.Environment.ValueString()
	realmConfigUpdate := propelauth.RealmConfigUpdate{
		AllowPublicSignups:                    plan.AllowPublicSignups.ValueBoolPointer(),
		AutoConfirmEmails:                     propelauth.FlipBoolRef(plan.RequireEmailConfirmation.ValueBoolPointer()),
		WaitlistUsersRequireEmailConfirmation: plan.WaitlistUsersRequireEmailConfirmation.ValueBoolPointer(),
	}

	realmConfigResponse, err := r.client.UpdateRealmConfig(environment, realmConfigUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting environment-level auth configuration",
			"Could not set environment-level auth configuration, unexpected error: "+err.Error(),
		)
		return
	}

	// Check that all field were updated to the new value if not empty
	if plan.AllowPublicSignups.ValueBoolPointer() != nil &&
		plan.AllowPublicSignups.ValueBool() != realmConfigResponse.AllowPublicSignups {
		resp.Diagnostics.AddError(
			"Error updating environment-level auth configuration",
			"AllowPublicSignups failed to update. The `allow_public_signups` is instead "+fmt.Sprintf("%t", realmConfigResponse.AllowPublicSignups),
		)
		return
	}
	if plan.RequireEmailConfirmation.ValueBoolPointer() != nil &&
		plan.RequireEmailConfirmation.ValueBool() == realmConfigResponse.AutoConfirmEmails {
		resp.Diagnostics.AddError(
			"Error updating environment-level auth configuration",
			"RequireEmailConfirmation failed to update. The `require_email_confirmation` is instead "+fmt.Sprintf("%t", realmConfigResponse.AutoConfirmEmails),
		)
		return
	}
	if plan.WaitlistUsersRequireEmailConfirmation.ValueBoolPointer() != nil &&
		plan.WaitlistUsersRequireEmailConfirmation.ValueBool() != realmConfigResponse.WaitlistUsersRequireEmailConfirmation {
		resp.Diagnostics.AddError(
			"Error updating environment-level auth configuration",
			"WaitlistUsersRequireEmailConfirmation failed to update. The `waitlist_users_require_email_confirmation` is instead "+fmt.Sprintf("%t", realmConfigResponse.WaitlistUsersRequireEmailConfirmation),
		)
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a propelauth_environment_level_auth_configuration resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *environmentLevelAuthConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state and read it into the model
	var state environmentLevelAuthConfigurationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// retrieve the environment config from PropelAuth
	realmConfigResponse, err := r.client.GetRealmConfig(state.Environment.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading PropelAuth environment-level auth configuration",
			"Could not read PropelAuth environment-level auth configuration: "+err.Error(),
		)
		return
	}

	// Save into the Terraform state only if the value is not null in Terraform.
	// Null, or unset values, in Terraform are left to be manually managed in the dashboard.
	if state.AllowPublicSignups.ValueBoolPointer() != nil {
		state.AllowPublicSignups = types.BoolValue(realmConfigResponse.AllowPublicSignups)
	}
	if state.RequireEmailConfirmation.ValueBoolPointer() != nil {
		state.RequireEmailConfirmation = types.BoolValue(!realmConfigResponse.AutoConfirmEmails)
	}
	if state.WaitlistUsersRequireEmailConfirmation.ValueBoolPointer() != nil {
		state.WaitlistUsersRequireEmailConfirmation = types.BoolValue(realmConfigResponse.WaitlistUsersRequireEmailConfirmation)
	}

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *environmentLevelAuthConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan environmentLevelAuthConfigurationResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the configuration in PropelAuth
	realmConfigUpdate := propelauth.RealmConfigUpdate{
		AllowPublicSignups:                    plan.AllowPublicSignups.ValueBoolPointer(),
		AutoConfirmEmails:                     propelauth.FlipBoolRef(plan.RequireEmailConfirmation.ValueBoolPointer()),
		WaitlistUsersRequireEmailConfirmation: plan.WaitlistUsersRequireEmailConfirmation.ValueBoolPointer(),
	}

	realmConfigResponse, err := r.client.UpdateRealmConfig(plan.Environment.ValueString(), realmConfigUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting environment-level auth configuration",
			"Could not set environment-level auth configuration, unexpected error: "+err.Error(),
		)
		return
	}

	// Check that all field were updated to the new value if not empty
	if plan.AllowPublicSignups.ValueBoolPointer() != nil &&
		plan.AllowPublicSignups.ValueBool() != realmConfigResponse.AllowPublicSignups {
		resp.Diagnostics.AddError(
			"Error updating environment-level auth configuration",
			"AllowPublicSignups failed to update. The `allow_public_signups` is instead "+fmt.Sprintf("%t", realmConfigResponse.AllowPublicSignups),
		)
		return
	}
	if plan.RequireEmailConfirmation.ValueBoolPointer() != nil &&
		plan.RequireEmailConfirmation.ValueBool() == realmConfigResponse.AutoConfirmEmails {
		resp.Diagnostics.AddError(
			"Error updating environment-level auth configuration",
			"RequireEmailConfirmation failed to update. The `require_email_confirmation` is instead "+fmt.Sprintf("%t", realmConfigResponse.AutoConfirmEmails),
		)
		return
	}
	if plan.WaitlistUsersRequireEmailConfirmation.ValueBoolPointer() != nil &&
		plan.WaitlistUsersRequireEmailConfirmation.ValueBool() != realmConfigResponse.WaitlistUsersRequireEmailConfirmation {
		resp.Diagnostics.AddError(
			"Error updating environment-level auth configuration",
			"WaitlistUsersRequireEmailConfirmation failed to update. The `waitlist_users_require_email_confirmation` is instead "+fmt.Sprintf("%t", realmConfigResponse.WaitlistUsersRequireEmailConfirmation),
		)
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "updated a propelauth_environment_level_auth_configuration resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *environmentLevelAuthConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "deleted a propelauth_environment_level_auth_configuration resource")
}
