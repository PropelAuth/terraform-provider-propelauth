package provider

import (
	"context"
	"fmt"

	"terraform-provider-propelauth/internal/propelauth"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &basicAuthConfigurationResource{}
var _ resource.ResourceWithConfigure = &basicAuthConfigurationResource{}
var _ resource.ResourceWithValidateConfig = &basicAuthConfigurationResource{}

func NewBasicAuthConfigurationResource() resource.Resource {
	return &basicAuthConfigurationResource{}
}

// basicAuthConfigurationResource defines the resource implementation.
type basicAuthConfigurationResource struct {
	client *propelauth.PropelAuthClient
}

// basicAuthConfigurationResourceModel describes the resource data model.
type basicAuthConfigurationResourceModel struct {
	AllowUsersToSignupWithPersonalEmail types.Bool     `tfsdk:"allow_users_to_signup_with_personal_email"`
	SignupDomainAllowlist               []types.String `tfsdk:"signup_domain_allowlist"`
	SignupDomainBlocklist               []types.String `tfsdk:"signup_domain_blocklist"`
	HasPasswordLogin                    types.Bool     `tfsdk:"has_password_login"`
	HasPasswordlessLogin                types.Bool     `tfsdk:"has_passwordless_login"`
	WaitlistUsersEnabled                types.Bool     `tfsdk:"waitlist_users_enabled"`
	UserAutologoutSeconds               types.Int64    `tfsdk:"user_autologout_seconds"`
	UserAutologoutType                  types.String   `tfsdk:"user_autologout_type"`
	UsersCanDeleteOwnAccount            types.Bool     `tfsdk:"users_can_delete_own_account"`
	UsersCanChangeEmail                 types.Bool     `tfsdk:"users_can_change_email"`
	IncludeLoginMethod                  types.Bool     `tfsdk:"include_login_method"`
}

func (r *basicAuthConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_basic_auth_configuration"
}

func (r *basicAuthConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Basic Auth Configuration. This is for configuring basic authentication, signup, and " +
			"user-account-management settings in PropelAuth.",
		Attributes: map[string]schema.Attribute{
			"allow_users_to_signup_with_personal_email": schema.BoolAttribute{
				Optional: true,
				Description: "If true, your users will be able to sign up using personal email domains (@gmail.com, @yahoo.com, etc.)." +
					"The default setting is true. This is only enabled if `signup_domain_allowlist` is empty.",
			},
			// "signup_domain_allowlist_enabled": schema.BoolAttribute{
			// 	Optional:    true,
			// 	Description: "If true, only users with email domains in the allowlist will be able to sign up. The default setting is false.",
			// },
			"signup_domain_allowlist": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "A list of email domains that are allowed to sign up. This is only used if `signup_domain_allowlist_enabled` is true.",
			},
			"signup_domain_blocklist": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "A list of email domains that are blocked from signing up. This is only used if `signup_domain_blocklist_enabled` is true and `signup_domain_allowlist` is empty.",
			},
			// "signup_domain_blocklist_enabled": schema.BoolAttribute{
			// 	Optional:    true,
			// 	Description: "If true, users with email domains in the blocklist will not be able to sign up. The default setting is false. This is only used if `signup_domain_allowlist_enabled` is false.",
			// },
			"has_password_login": schema.BoolAttribute{
				Optional:    true,
				Description: "If true, your users will be able to log in using their email and password. The default setting is true.",
			},
			"has_passwordless_login": schema.BoolAttribute{
				Optional:    true,
				Description: "If true, your users will be able to log in using a magic link sent to their email. The default setting is false.",
			},
			"waitlist_users_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "If true, you will be able to use the APIs to collect emails before launching. The default setting is false.",
			},
			"user_autologout_seconds": schema.Int64Attribute{
				Optional: true,
				Description: "The number of seconds before a user is automatically logged out. The default setting is 1209600 (14 days)." +
					"See also \"user_autologout_type\" for more information.",
			},
			"user_autologout_type": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("AfterInactivity", "AfterLogin"),
				},
				Description: "This sets the behavior for when the counting for \"user_autologout_seconds\" starts. " +
					"Valid values are \"AfterInactivity\" and the stricter \"AfterLogin\". The default setting is \"AfterInactivity\".",
			},
			"users_can_delete_own_account": schema.BoolAttribute{
				Optional:    true,
				Description: "If true, your users will be able to delete their own account. The default setting is false.",
			},
			"users_can_change_email": schema.BoolAttribute{
				Optional:    true,
				Description: "If true, your users will be able to change their email address. The default setting is true.",
			},
			"include_login_method": schema.BoolAttribute{
				Optional: true,
				Description: "If true, the login method will be included in the access token. The default setting is false." +
					"See `https://docs.propelauth.com/overview/user-management/user-properties#login-method-property` for more information.",
			},
		},
	}
}

func (r *basicAuthConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *basicAuthConfigurationResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var plan basicAuthConfigurationResourceModel

	// Read Terraform plan data into the model
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate the plan data
	if plan.SignupDomainAllowlist != nil && plan.SignupDomainBlocklist != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("signup_domain_allowlist"),
			"Invalid `signup_domain_allowlist`",
			"`signup_domain_allowlist` and `signup_domain_blocklist` cannot both be set",
		)
		return
	}
}

func (r *basicAuthConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan basicAuthConfigurationResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the configuration in PropelAuth
	environmentConfigUpdate := propelauth.EnvironmentConfigUpdate{
		AllowUsersToSignupWithPersonalEmail: plan.AllowUsersToSignupWithPersonalEmail.ValueBoolPointer(),
		HasPasswordLogin:                    plan.HasPasswordLogin.ValueBoolPointer(),
		HasPasswordlessLogin:                plan.HasPasswordlessLogin.ValueBoolPointer(),
		WaitlistUsersEnabled:                plan.WaitlistUsersEnabled.ValueBoolPointer(),
		UserAutologoutSeconds:               plan.UserAutologoutSeconds.ValueInt64Pointer(),
		UserAutologoutType:                  plan.UserAutologoutType.ValueString(),
		UsersCanDeleteOwnAccount:            plan.UsersCanDeleteOwnAccount.ValueBoolPointer(),
		UsersCanChangeEmail:                 plan.UsersCanChangeEmail.ValueBoolPointer(),
		IncludeLoginMethod:                  plan.IncludeLoginMethod.ValueBoolPointer(),
	}

	var signupDomainAllowlistEnabled bool
	if plan.SignupDomainAllowlist != nil {
		signupDomainAllowlistEnabled = true
		environmentConfigUpdate.SignupDomainAllowlist = make([]string, len(plan.SignupDomainAllowlist))
		for i, domain := range plan.SignupDomainAllowlist {
			environmentConfigUpdate.SignupDomainAllowlist[i] = domain.ValueString()
		}
	} else {
		signupDomainAllowlistEnabled = false
	}
	environmentConfigUpdate.SignupDomainAllowlistEnabled = &signupDomainAllowlistEnabled

	var signupDomainBlocklistEnabled bool
	if plan.SignupDomainBlocklist != nil {
		signupDomainBlocklistEnabled = true
		environmentConfigUpdate.SignupDomainBlocklist = make([]string, len(plan.SignupDomainBlocklist))
		for i, domain := range plan.SignupDomainBlocklist {
			environmentConfigUpdate.SignupDomainBlocklist[i] = domain.ValueString()
		}
	} else {
		signupDomainBlocklistEnabled = false
	}
	environmentConfigUpdate.SignupDomainBlocklistEnabled = &signupDomainBlocklistEnabled

	environmentConfigResponse, err := r.client.UpdateEnvironmentConfig(&environmentConfigUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting basic auth configuration",
			"Could not set basic auth configuration, unexpected error: "+err.Error(),
		)
		return
	}

	// Check that all field were updated to the new value if not empty
	if plan.AllowUsersToSignupWithPersonalEmail.ValueBoolPointer() != nil &&
		plan.AllowUsersToSignupWithPersonalEmail.ValueBool() != environmentConfigResponse.AllowUsersToSignupWithPersonalEmail {
		resp.Diagnostics.AddError(
			"Error updating basic auth configuration",
			"AllowUsersToSignupWithPersonalEmail failed to update. The `allow_users_to_signup_with_personal_email` is instead "+fmt.Sprintf("%t", environmentConfigResponse.AllowUsersToSignupWithPersonalEmail),
		)
		return
	}
	if plan.HasPasswordLogin.ValueBoolPointer() != nil &&
		plan.HasPasswordLogin.ValueBool() != environmentConfigResponse.HasPasswordLogin {
		resp.Diagnostics.AddError(
			"Error updating basic auth configuration",
			"HasPasswordLogin failed to update. The `has_password_login` is instead "+fmt.Sprintf("%t", environmentConfigResponse.HasPasswordLogin),
		)
		return
	}
	if plan.HasPasswordlessLogin.ValueBoolPointer() != nil &&
		plan.HasPasswordlessLogin.ValueBool() != environmentConfigResponse.HasPasswordlessLogin {
		resp.Diagnostics.AddError(
			"Error updating basic auth configuration",
			"HasPasswordlessLogin failed to update. The `has_passwordless_login` is instead "+fmt.Sprintf("%t", environmentConfigResponse.HasPasswordlessLogin),
		)
		return
	}
	if plan.WaitlistUsersEnabled.ValueBoolPointer() != nil &&
		plan.WaitlistUsersEnabled.ValueBool() != environmentConfigResponse.WaitlistUsersEnabled {
		resp.Diagnostics.AddError(
			"Error updating basic auth configuration",
			"WaitlistUsersEnabled failed to update. The `waitlist_users_enabled` is instead "+fmt.Sprintf("%t", environmentConfigResponse.WaitlistUsersEnabled),
		)
		return
	}
	if plan.UserAutologoutSeconds.ValueInt64Pointer() != nil &&
		plan.UserAutologoutSeconds.ValueInt64() != environmentConfigResponse.UserAutologoutSeconds {
		resp.Diagnostics.AddError(
			"Error updating basic auth configuration",
			"UserAutologoutSeconds failed to update. The `user_autologout_seconds` is instead "+fmt.Sprintf("%d", environmentConfigResponse.UserAutologoutSeconds),
		)
		return
	}
	if plan.UserAutologoutType.ValueString() != "" &&
		plan.UserAutologoutType.ValueString() != environmentConfigResponse.UserAutologoutType {
		resp.Diagnostics.AddError(
			"Error updating basic auth configuration",
			"UserAutologoutType failed to update. The `user_autologout_type` is instead "+environmentConfigResponse.UserAutologoutType,
		)
		return
	}
	if plan.UsersCanDeleteOwnAccount.ValueBoolPointer() != nil &&
		plan.UsersCanDeleteOwnAccount.ValueBool() != environmentConfigResponse.UsersCanDeleteOwnAccount {
		resp.Diagnostics.AddError(
			"Error updating basic auth configuration",
			"UsersCanDeleteOwnAccount failed to update. The `users_can_delete_own_account` is instead "+fmt.Sprintf("%t", environmentConfigResponse.UsersCanDeleteOwnAccount),
		)
		return
	}
	if plan.UsersCanChangeEmail.ValueBoolPointer() != nil &&
		plan.UsersCanChangeEmail.ValueBool() != environmentConfigResponse.UsersCanChangeEmail {
		resp.Diagnostics.AddError(
			"Error updating basic auth configuration",
			"UsersCanChangeEmail failed to update. The `users_can_change_email` is instead "+fmt.Sprintf("%t", environmentConfigResponse.UsersCanChangeEmail),
		)
		return
	}
	if plan.IncludeLoginMethod.ValueBoolPointer() != nil &&
		plan.IncludeLoginMethod.ValueBool() != environmentConfigResponse.IncludeLoginMethod {
		resp.Diagnostics.AddError(
			"Error updating basic auth configuration",
			"IncludeLoginMethod failed to update. The `include_login_method` is instead "+fmt.Sprintf("%t", environmentConfigResponse.IncludeLoginMethod),
		)
		return
	}
	if plan.SignupDomainAllowlist != nil {
		if len(plan.SignupDomainAllowlist) != len(environmentConfigResponse.SignupDomainAllowlist) {
			resp.Diagnostics.AddError(
				"Error updating basic auth configuration",
				"SignupDomainAllowlist failed to update. The `signup_domain_allowlist` is instead "+fmt.Sprintf("%v", environmentConfigResponse.SignupDomainAllowlist),
			)
			return
		}
		for i, domain := range plan.SignupDomainAllowlist {
			if domain.ValueString() != environmentConfigResponse.SignupDomainAllowlist[i] {
				resp.Diagnostics.AddError(
					"Error updating basic auth configuration",
					"SignupDomainAllowlist failed to update. The `signup_domain_allowlist` is instead "+fmt.Sprintf("%v", environmentConfigResponse.SignupDomainAllowlist),
				)
				return
			}
		}
	}
	if plan.SignupDomainBlocklist != nil {
		if len(plan.SignupDomainBlocklist) != len(environmentConfigResponse.SignupDomainBlocklist) {
			resp.Diagnostics.AddError(
				"Error updating basic auth configuration",
				"SignupDomainBlocklist failed to update. The `signup_domain_blocklist` is instead "+fmt.Sprintf("%v", environmentConfigResponse.SignupDomainBlocklist),
			)
			return
		}
		for i, domain := range plan.SignupDomainBlocklist {
			if domain.ValueString() != environmentConfigResponse.SignupDomainBlocklist[i] {
				resp.Diagnostics.AddError(
					"Error updating basic auth configuration",
					"SignupDomainBlocklist failed to update. The `signup_domain_blocklist` is instead "+fmt.Sprintf("%v", environmentConfigResponse.SignupDomainBlocklist),
				)
				return
			}
		}
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a propelauth_basic_auth_configuration resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *basicAuthConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state and read it into the model
	var state basicAuthConfigurationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// retrieve the environment config from PropelAuth
	environmentConfigResponse, err := r.client.GetEnvironmentConfig()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading PropelAuth basic auth configuration",
			"Could not read PropelAuth basic auth configuration: "+err.Error(),
		)
		return
	}

	// Save into the Terraform state only if the value is not null in Terraform.
	// Null, or unset values, in Terraform are left to be manually managed in the dashboard.
	if state.AllowUsersToSignupWithPersonalEmail.ValueBoolPointer() != nil {
		state.AllowUsersToSignupWithPersonalEmail = types.BoolValue(environmentConfigResponse.AllowUsersToSignupWithPersonalEmail)
	}
	if state.HasPasswordLogin.ValueBoolPointer() != nil {
		state.HasPasswordLogin = types.BoolValue(environmentConfigResponse.HasPasswordLogin)
	}
	if state.HasPasswordlessLogin.ValueBoolPointer() != nil {
		state.HasPasswordlessLogin = types.BoolValue(environmentConfigResponse.HasPasswordlessLogin)
	}
	if state.WaitlistUsersEnabled.ValueBoolPointer() != nil {
		state.WaitlistUsersEnabled = types.BoolValue(environmentConfigResponse.WaitlistUsersEnabled)
	}
	if state.UserAutologoutSeconds.ValueInt64Pointer() != nil {
		state.UserAutologoutSeconds = types.Int64Value(environmentConfigResponse.UserAutologoutSeconds)
	}
	if state.UserAutologoutType.ValueString() != "" {
		state.UserAutologoutType = types.StringValue(environmentConfigResponse.UserAutologoutType)
	}
	if state.UsersCanDeleteOwnAccount.ValueBoolPointer() != nil {
		state.UsersCanDeleteOwnAccount = types.BoolValue(environmentConfigResponse.UsersCanDeleteOwnAccount)
	}
	if state.UsersCanChangeEmail.ValueBoolPointer() != nil {
		state.UsersCanChangeEmail = types.BoolValue(environmentConfigResponse.UsersCanChangeEmail)
	}
	if state.IncludeLoginMethod.ValueBoolPointer() != nil {
		state.IncludeLoginMethod = types.BoolValue(environmentConfigResponse.IncludeLoginMethod)
	}

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *basicAuthConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan basicAuthConfigurationResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the configuration in PropelAuth
	environmentConfigUpdate := propelauth.EnvironmentConfigUpdate{
		AllowUsersToSignupWithPersonalEmail: plan.AllowUsersToSignupWithPersonalEmail.ValueBoolPointer(),
		HasPasswordLogin:                    plan.HasPasswordLogin.ValueBoolPointer(),
		HasPasswordlessLogin:                plan.HasPasswordlessLogin.ValueBoolPointer(),
		WaitlistUsersEnabled:                plan.WaitlistUsersEnabled.ValueBoolPointer(),
		UserAutologoutSeconds:               plan.UserAutologoutSeconds.ValueInt64Pointer(),
		UserAutologoutType:                  plan.UserAutologoutType.ValueString(),
		UsersCanDeleteOwnAccount:            plan.UsersCanDeleteOwnAccount.ValueBoolPointer(),
		UsersCanChangeEmail:                 plan.UsersCanChangeEmail.ValueBoolPointer(),
		IncludeLoginMethod:                  plan.IncludeLoginMethod.ValueBoolPointer(),
	}

	var signupDomainAllowlistEnabled bool
	if plan.SignupDomainAllowlist != nil {
		signupDomainAllowlistEnabled = true
		environmentConfigUpdate.SignupDomainAllowlist = make([]string, len(plan.SignupDomainAllowlist))
		for i, domain := range plan.SignupDomainAllowlist {
			environmentConfigUpdate.SignupDomainAllowlist[i] = domain.ValueString()
		}
	} else {
		signupDomainAllowlistEnabled = false
	}
	environmentConfigUpdate.SignupDomainAllowlistEnabled = &signupDomainAllowlistEnabled

	var signupDomainBlocklistEnabled bool
	if plan.SignupDomainBlocklist != nil {
		signupDomainBlocklistEnabled = true
		environmentConfigUpdate.SignupDomainBlocklist = make([]string, len(plan.SignupDomainBlocklist))
		for i, domain := range plan.SignupDomainBlocklist {
			environmentConfigUpdate.SignupDomainBlocklist[i] = domain.ValueString()
		}
	} else {
		signupDomainBlocklistEnabled = false
	}
	environmentConfigUpdate.SignupDomainBlocklistEnabled = &signupDomainBlocklistEnabled

	environmentConfigResponse, err := r.client.UpdateEnvironmentConfig(&environmentConfigUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting basic auth configuration",
			"Could not set basic auth configuration, unexpected error: "+err.Error(),
		)
		return
	}

	// Check that all field were updated to the new value if not empty
	if plan.AllowUsersToSignupWithPersonalEmail.ValueBoolPointer() != nil &&
		plan.AllowUsersToSignupWithPersonalEmail.ValueBool() != environmentConfigResponse.AllowUsersToSignupWithPersonalEmail {
		resp.Diagnostics.AddError(
			"Error updating basic auth configuration",
			"AllowUsersToSignupWithPersonalEmail failed to update. The `allow_users_to_signup_with_personal_email` is instead "+fmt.Sprintf("%t", environmentConfigResponse.AllowUsersToSignupWithPersonalEmail),
		)
		return
	}
	if plan.HasPasswordLogin.ValueBoolPointer() != nil &&
		plan.HasPasswordLogin.ValueBool() != environmentConfigResponse.HasPasswordLogin {
		resp.Diagnostics.AddError(
			"Error updating basic auth configuration",
			"HasPasswordLogin failed to update. The `has_password_login` is instead "+fmt.Sprintf("%t", environmentConfigResponse.HasPasswordLogin),
		)
		return
	}
	if plan.HasPasswordlessLogin.ValueBoolPointer() != nil &&
		plan.HasPasswordlessLogin.ValueBool() != environmentConfigResponse.HasPasswordlessLogin {
		resp.Diagnostics.AddError(
			"Error updating basic auth configuration",
			"HasPasswordlessLogin failed to update. The `has_passwordless_login` is instead "+fmt.Sprintf("%t", environmentConfigResponse.HasPasswordlessLogin),
		)
		return
	}
	if plan.WaitlistUsersEnabled.ValueBoolPointer() != nil &&
		plan.WaitlistUsersEnabled.ValueBool() != environmentConfigResponse.WaitlistUsersEnabled {
		resp.Diagnostics.AddError(
			"Error updating basic auth configuration",
			"WaitlistUsersEnabled failed to update. The `waitlist_users_enabled` is instead "+fmt.Sprintf("%t", environmentConfigResponse.WaitlistUsersEnabled),
		)
		return
	}
	if plan.UserAutologoutSeconds.ValueInt64Pointer() != nil &&
		plan.UserAutologoutSeconds.ValueInt64() != environmentConfigResponse.UserAutologoutSeconds {
		resp.Diagnostics.AddError(
			"Error updating basic auth configuration",
			"UserAutologoutSeconds failed to update. The `user_autologout_seconds` is instead "+fmt.Sprintf("%d", environmentConfigResponse.UserAutologoutSeconds),
		)
		return
	}
	if plan.UserAutologoutType.ValueString() != "" &&
		plan.UserAutologoutType.ValueString() != environmentConfigResponse.UserAutologoutType {
		resp.Diagnostics.AddError(
			"Error updating basic auth configuration",
			"UserAutologoutType failed to update. The `user_autologout_type` is instead "+environmentConfigResponse.UserAutologoutType,
		)
		return
	}
	if plan.UsersCanDeleteOwnAccount.ValueBoolPointer() != nil &&
		plan.UsersCanDeleteOwnAccount.ValueBool() != environmentConfigResponse.UsersCanDeleteOwnAccount {
		resp.Diagnostics.AddError(
			"Error updating basic auth configuration",
			"UsersCanDeleteOwnAccount failed to update. The `users_can_delete_own_account` is instead "+fmt.Sprintf("%t", environmentConfigResponse.UsersCanDeleteOwnAccount),
		)
		return
	}
	if plan.UsersCanChangeEmail.ValueBoolPointer() != nil &&
		plan.UsersCanChangeEmail.ValueBool() != environmentConfigResponse.UsersCanChangeEmail {
		resp.Diagnostics.AddError(
			"Error updating basic auth configuration",
			"UsersCanChangeEmail failed to update. The `users_can_change_email` is instead "+fmt.Sprintf("%t", environmentConfigResponse.UsersCanChangeEmail),
		)
		return
	}
	if plan.IncludeLoginMethod.ValueBoolPointer() != nil &&
		plan.IncludeLoginMethod.ValueBool() != environmentConfigResponse.IncludeLoginMethod {
		resp.Diagnostics.AddError(
			"Error updating basic auth configuration",
			"IncludeLoginMethod failed to update. The `include_login_method` is instead "+fmt.Sprintf("%t", environmentConfigResponse.IncludeLoginMethod),
		)
		return
	}
	if plan.SignupDomainAllowlist != nil {
		if len(plan.SignupDomainAllowlist) != len(environmentConfigResponse.SignupDomainAllowlist) {
			resp.Diagnostics.AddError(
				"Error updating basic auth configuration",
				"SignupDomainAllowlist failed to update. The `signup_domain_allowlist` is instead "+fmt.Sprintf("%v", environmentConfigResponse.SignupDomainAllowlist),
			)
			return
		}
		for i, domain := range plan.SignupDomainAllowlist {
			if domain.ValueString() != environmentConfigResponse.SignupDomainAllowlist[i] {
				resp.Diagnostics.AddError(
					"Error updating basic auth configuration",
					"SignupDomainAllowlist failed to update. The `signup_domain_allowlist` is instead "+fmt.Sprintf("%v", environmentConfigResponse.SignupDomainAllowlist),
				)
				return
			}
		}
	}
	if plan.SignupDomainBlocklist != nil {
		if len(plan.SignupDomainBlocklist) != len(environmentConfigResponse.SignupDomainBlocklist) {
			resp.Diagnostics.AddError(
				"Error updating basic auth configuration",
				"SignupDomainBlocklist failed to update. The `signup_domain_blocklist` is instead "+fmt.Sprintf("%v", environmentConfigResponse.SignupDomainBlocklist),
			)
			return
		}
		for i, domain := range plan.SignupDomainBlocklist {
			if domain.ValueString() != environmentConfigResponse.SignupDomainBlocklist[i] {
				resp.Diagnostics.AddError(
					"Error updating basic auth configuration",
					"SignupDomainBlocklist failed to update. The `signup_domain_blocklist` is instead "+fmt.Sprintf("%v", environmentConfigResponse.SignupDomainBlocklist),
				)
				return
			}
		}
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "updated a propelauth_basic_auth_configuration resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *basicAuthConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "deleted a propelauth_basic_auth_configuration resource")
}
