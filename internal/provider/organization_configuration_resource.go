package provider

import (
	"context"
	"fmt"

	"terraform-provider-propelauth/internal/propelauth"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &organizationConfigurationResource{}
var _ resource.ResourceWithConfigure = &organizationConfigurationResource{}
var _ resource.ResourceWithImportState = &organizationConfigurationResource{}

func NewOrganizationConfigurationResource() resource.Resource {
	return &organizationConfigurationResource{}
}

// organizationConfigurationResource defines the resource implementation.
type organizationConfigurationResource struct {
	client *propelauth.PropelAuthClient
}

// organizationConfigurationResourceModel describes the resource data model.
type organizationConfigurationResourceModel struct {
	HasOrgs                     types.Bool                        `tfsdk:"has_orgs"`
	MaxNumOrgsUsersCanBeIn      types.Int32                       `tfsdk:"max_num_orgs_users_can_be_in"`
	OrgsMetaname                types.String                      `tfsdk:"orgs_metaname"`
	UsersCanCreateOrgs          types.Bool                        `tfsdk:"users_can_create_orgs"`
	UsersCanDeleteTheirOwnOrgs  types.Bool                        `tfsdk:"users_can_delete_their_own_orgs"`
	UsersMustBeInAnOrganization types.Bool                        `tfsdk:"users_must_be_in_an_organization"`
	OrgsCanSetupSaml            types.Bool                        `tfsdk:"orgs_can_setup_saml"`
	UseOrgNameForSaml           types.Bool                        `tfsdk:"use_org_name_for_saml"`
	DefaultToSamlLogin          types.Bool                        `tfsdk:"default_to_saml_login"`
	SkipSamlRoleMappingStep     types.Bool                        `tfsdk:"skip_saml_role_mapping_step"`
	OrgsCanRequire2fa           types.Bool                        `tfsdk:"orgs_can_require_2fa"`
	CustomerOrgAuditLogSettings *CustomerOrgAuditLogSettingsModel `tfsdk:"customer_org_audit_log_settings"`
}

type CustomerOrgAuditLogSettingsModel struct {
	Enabled                     types.Bool `tfsdk:"enabled"`
	AllOrgsCanViewTheirAuditLog types.Bool `tfsdk:"all_orgs_can_view_their_audit_log"`
	IncludeImpersonation        types.Bool `tfsdk:"include_impersonation"`
	IncludeEmployeeActions      types.Bool `tfsdk:"include_employee_actions"`
	IncludeApiKeyActions        types.Bool `tfsdk:"include_api_key_actions"`
}

func (r *organizationConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_configuration"
}

func (r *organizationConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Organization Configuration. This is for configuring your global organization settings in PropelAuth. " +
			"Settings on specific organizations can be managed in the dashboard.",
		Attributes: map[string]schema.Attribute{
			"has_orgs": schema.BoolAttribute{
				Optional: true,
				Description: "This is the top level setting for whether organizations are in your PropelAuth integration." +
					"If false, all other organization settings are ignored. The default setting is true.",
			},
			"max_num_orgs_users_can_be_in": schema.Int32Attribute{
				Optional: true,
				Description: "This is the maximum number of organizations a user can be a member of. If a user tries to exceed this number, " +
					"they will be asked to leave an existing organization. The default setting is 10.",
			},
			"orgs_metaname": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 50),
				},
				Description: "What name do you use for organizations? This will update the copy across your hosted pages." +
					"The default setting is 'Organization'.",
			},
			"users_can_create_orgs": schema.BoolAttribute{
				Optional: true,
				Description: "If true, users have access to the 'Create Org' UI, allowing them to create their own organizations." +
					"The default setting is true.",
			},
			"users_can_delete_their_own_orgs": schema.BoolAttribute{
				Optional: true,
				Description: "If true, users with the requisite permission will be able to delete their organizations. " +
					"The default setting is false.",
			},
			"users_must_be_in_an_organization": schema.BoolAttribute{
				Optional: true,
				Description: "If true, users will be required to create or join an organization as part of the signup process. " +
					"The default setting is false.",
			},
			"orgs_can_setup_saml": schema.BoolAttribute{
				Optional: true,
				Description: "If true, your users can setup a SAML connection for their organization. This allows them to " +
					"log into your product using their existing work account managed by an Identity Provider like " +
					"Okta, Azure/Entra, Google, and more. The default setting is false. " +
					"Warning: This is only applied in prod for some billing plans",
			},
			"use_org_name_for_saml": schema.BoolAttribute{
				Optional: true,
				Description: "This is an advanced setting that only applies if SAML is enabled. If true, " +
					"users can look up and be redirected to their SSO provider using their organization's name." +
					"The default setting is false which means the SAML provider is instead inferred from their email address.",
			},
			"default_to_saml_login": schema.BoolAttribute{
				Optional: true,
				Description: "This is an advanced setting that only applies if SAML is enabled. If true, " +
					"affected users will be directed to SAML by default in the hosted pages." +
					"The default setting is false.",
			},
			"skip_saml_role_mapping_step": schema.BoolAttribute{
				Optional: true,
				Description: "This is an advanced setting that only applies if SAML is enabled. If true, " +
					"end users setting up SAML for their organization will not see the role-mapping step. " +
					"The default setting is false.",
			},
			"orgs_can_require_2fa": schema.BoolAttribute{
				Optional: true,
				Description: "If true, organizations can require their users to use 2FA." +
					"The default setting is false. " +
					"Warning: This is only applied in prod for some billing plans",
			},
			"customer_org_audit_log_settings": schema.SingleNestedAttribute{
				Optional: true,
				Description: "Settings for enabling whether and configuring how your customer organizations will have access to " +
					"their own audit log.\n\nNote: This feature is only available for use by your customer organizations in " +
					"non-test environments for some pricing plans.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Required:    true,
						Description: "If enabled, your customer organizations will have access to their own audit log.",
					},
					"all_orgs_can_view_their_audit_log": schema.BoolAttribute{
						Required: true,
						Description: "If true, all of your customer organization will automatically have access to this feature. " +
							"Otherwise, you will need to enable it for each organization individually.",
					},
					"include_impersonation": schema.BoolAttribute{
						Required: true,
						Description: "If true, the audit log will include actions that were triggered by a member of your team " +
							"impersonating one of their organization members. The impersonator will be anonymous to your customer.",
					},
					"include_api_key_actions": schema.BoolAttribute{
						Required: true,
						Description: "If true, the audit log will include actions that were triggered by your BE service utilizing " +
							"PropelAuth APIs.",
					},
					"include_employee_actions": schema.BoolAttribute{
						Required: true,
						Description: "If true, the audit log will include actions that were triggered by a member of your team " +
							"using the PropelAuth dashboard. The person who triggered the action will be anonymous to your customer.",
					},
				},
			},
		},
	}
}

func (r *organizationConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *organizationConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan organizationConfigurationResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the configuration in PropelAuth
	environmentConfigUpdate := propelauth.EnvironmentConfigUpdate{
		HasOrgs:                     plan.HasOrgs.ValueBoolPointer(),
		MaxNumOrgsUsersCanBeIn:      plan.MaxNumOrgsUsersCanBeIn.ValueInt32Pointer(),
		OrgsMetaname:                plan.OrgsMetaname.ValueString(),
		UsersCanCreateOrgs:          plan.UsersCanCreateOrgs.ValueBoolPointer(),
		UsersCanDeleteTheirOwnOrgs:  plan.UsersCanDeleteTheirOwnOrgs.ValueBoolPointer(),
		UsersMustBeInAnOrganization: plan.UsersMustBeInAnOrganization.ValueBoolPointer(),
		OrgsCanSetupSaml:            plan.OrgsCanSetupSaml.ValueBoolPointer(),
		UseOrgNameForSaml:           plan.UseOrgNameForSaml.ValueBoolPointer(),
		DefaultToSamlLogin:          plan.DefaultToSamlLogin.ValueBoolPointer(),
		SkipSamlRoleMappingStep:     plan.SkipSamlRoleMappingStep.ValueBoolPointer(),
		OrgsCanRequire2fa:           plan.OrgsCanRequire2fa.ValueBoolPointer(),
	}

	if plan.CustomerOrgAuditLogSettings != nil {
		environmentConfigUpdate.OrgsCanViewOrgAuditLog = plan.CustomerOrgAuditLogSettings.Enabled.ValueBoolPointer()
		environmentConfigUpdate.AllOrgsCanViewOrgAuditLog = plan.CustomerOrgAuditLogSettings.AllOrgsCanViewTheirAuditLog.ValueBoolPointer()
		environmentConfigUpdate.OrgAuditLogIncludesImpersonation = plan.CustomerOrgAuditLogSettings.IncludeImpersonation.ValueBoolPointer()
		environmentConfigUpdate.OrgAuditLogIncludesEmployees = plan.CustomerOrgAuditLogSettings.IncludeEmployeeActions.ValueBoolPointer()
		environmentConfigUpdate.OrgAuditLogIncludesApiKeys = plan.CustomerOrgAuditLogSettings.IncludeApiKeyActions.ValueBoolPointer()
	}

	environmentConfigResponse, err := r.client.UpdateEnvironmentConfig(&environmentConfigUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting organization configuration",
			"Could not set organization configuration, unexpected error: "+err.Error(),
		)
		return
	}

	// Check that all field were updated to the new value if not empty
	if plan.HasOrgs.ValueBoolPointer() != nil &&
		plan.HasOrgs.ValueBool() != environmentConfigResponse.HasOrgs {
		resp.Diagnostics.AddError(
			"Error updating organization configuration",
			"HasOrgs failed to update. The `allow_users_to_signup_with_personal_email` is instead "+fmt.Sprintf("%t", environmentConfigResponse.HasOrgs),
		)
		return
	}
	if plan.MaxNumOrgsUsersCanBeIn.ValueInt32Pointer() != nil &&
		plan.MaxNumOrgsUsersCanBeIn.ValueInt32() != environmentConfigResponse.MaxNumOrgsUsersCanBeIn {
		resp.Diagnostics.AddError(
			"Error updating organization configuration",
			"MaxNumOrgsUsersCanBeIn failed to update. The `max_num_orgs_users_can_be_in` is instead "+fmt.Sprintf("%d", environmentConfigResponse.MaxNumOrgsUsersCanBeIn),
		)
		return
	}
	if plan.OrgsMetaname.ValueString() != "" &&
		plan.OrgsMetaname.ValueString() != environmentConfigResponse.OrgsMetaname {
		resp.Diagnostics.AddError(
			"Error updating organization configuration",
			"OrgsMetaname failed to update. The `orgs_metaname` is instead "+environmentConfigResponse.OrgsMetaname,
		)
		return
	}
	if plan.UsersCanCreateOrgs.ValueBoolPointer() != nil &&
		plan.UsersCanCreateOrgs.ValueBool() != environmentConfigResponse.UsersCanCreateOrgs {
		resp.Diagnostics.AddError(
			"Error updating organization configuration",
			"UsersCanCreateOrgs failed to update. The `users_can_create_orgs` is instead "+fmt.Sprintf("%t", environmentConfigResponse.UsersCanCreateOrgs),
		)
		return
	}
	if plan.UsersCanDeleteTheirOwnOrgs.ValueBoolPointer() != nil &&
		plan.UsersCanDeleteTheirOwnOrgs.ValueBool() != environmentConfigResponse.UsersCanDeleteTheirOwnOrgs {
		resp.Diagnostics.AddError(
			"Error updating organization configuration",
			"UsersCanDeleteTheirOwnOrgs failed to update. The `users_can_delete_their_own_orgs` is instead "+fmt.Sprintf("%t", environmentConfigResponse.UsersCanDeleteTheirOwnOrgs),
		)
		return
	}
	if plan.UsersMustBeInAnOrganization.ValueBoolPointer() != nil &&
		plan.UsersMustBeInAnOrganization.ValueBool() != environmentConfigResponse.UsersMustBeInAnOrganization {
		resp.Diagnostics.AddError(
			"Error updating organization configuration",
			"UsersMustBeInAnOrganization failed to update. The `users_must_be_in_an_organization` is instead "+fmt.Sprintf("%t", environmentConfigResponse.UsersMustBeInAnOrganization),
		)
		return
	}
	if plan.OrgsCanSetupSaml.ValueBoolPointer() != nil &&
		plan.OrgsCanSetupSaml.ValueBool() != environmentConfigResponse.OrgsCanSetupSaml {
		resp.Diagnostics.AddError(
			"Error updating organization configuration",
			"OrgsCanSetupSaml failed to update. The `orgs_can_setup_saml` is instead "+fmt.Sprintf("%t", environmentConfigResponse.OrgsCanSetupSaml),
		)
		return
	}
	if plan.UseOrgNameForSaml.ValueBoolPointer() != nil &&
		plan.UseOrgNameForSaml.ValueBool() != environmentConfigResponse.UseOrgNameForSaml {
		resp.Diagnostics.AddError(
			"Error updating organization configuration",
			"UseOrgNameForSaml failed to update. The `use_org_name_for_saml` is instead "+fmt.Sprintf("%t", environmentConfigResponse.UseOrgNameForSaml),
		)
		return
	}
	if plan.DefaultToSamlLogin.ValueBoolPointer() != nil &&
		plan.DefaultToSamlLogin.ValueBool() != environmentConfigResponse.DefaultToSamlLogin {
		resp.Diagnostics.AddError(
			"Error updating organization configuration",
			"DefaultToSamlLogin failed to update. The `default_to_saml_login` is instead "+fmt.Sprintf("%t", environmentConfigResponse.DefaultToSamlLogin),
		)
		return
	}
	if plan.SkipSamlRoleMappingStep.ValueBoolPointer() != nil &&
		plan.SkipSamlRoleMappingStep.ValueBool() != environmentConfigResponse.SkipSamlRoleMappingStep {
		resp.Diagnostics.AddError(
			"Error updating organization configuration",
			"SkipSamlRoleMappingStep failed to update. The `skip_saml_role_mapping_step` is instead "+fmt.Sprintf("%t", environmentConfigResponse.SkipSamlRoleMappingStep),
		)
		return
	}
	if plan.OrgsCanRequire2fa.ValueBoolPointer() != nil &&
		plan.OrgsCanRequire2fa.ValueBool() != environmentConfigResponse.OrgsCanRequire2fa {
		resp.Diagnostics.AddError(
			"Error updating organization configuration",
			"OrgsCanRequire2fa failed to update. The `orgs_can_require_2fa` is instead "+fmt.Sprintf("%t", environmentConfigResponse.OrgsCanRequire2fa),
		)
		return
	}
	if plan.CustomerOrgAuditLogSettings != nil {
		if plan.CustomerOrgAuditLogSettings.Enabled.ValueBool() != environmentConfigResponse.OrgsCanViewOrgAuditLog {
			resp.Diagnostics.AddError(
				"Error updating organization configuration",
				"CustomerOrgAuditLogSettings.Enabled failed to update. The `customer_org_audit_log_settings.enabled` is instead "+fmt.Sprintf("%t", environmentConfigResponse.OrgsCanViewOrgAuditLog),
			)
			return
		}
		if plan.CustomerOrgAuditLogSettings.AllOrgsCanViewTheirAuditLog.ValueBool() != environmentConfigResponse.AllOrgsCanViewOrgAuditLog {
			resp.Diagnostics.AddError(
				"Error updating organization configuration",
				"CustomerOrgAuditLogSettings.AllOrgsCanViewTheirAuditLog failed to update. The `customer_org_audit_log_settings.all_orgs_can_view_their_audit_log` is instead "+fmt.Sprintf("%t", environmentConfigResponse.AllOrgsCanViewOrgAuditLog),
			)
			return
		}
		if plan.CustomerOrgAuditLogSettings.IncludeImpersonation.ValueBool() != environmentConfigResponse.OrgAuditLogIncludesImpersonation {
			resp.Diagnostics.AddError(
				"Error updating organization configuration",
				"CustomerOrgAuditLogSettings.IncludeImpersonation failed to update. The `customer_org_audit_log_settings.include_impersonation` is instead "+fmt.Sprintf("%t", environmentConfigResponse.OrgAuditLogIncludesImpersonation),
			)
			return
		}
		if plan.CustomerOrgAuditLogSettings.IncludeEmployeeActions.ValueBool() != environmentConfigResponse.OrgAuditLogIncludesEmployees {
			resp.Diagnostics.AddError(
				"Error updating organization configuration",
				"CustomerOrgAuditLogSettings.IncludeEmployeeActions failed to update. The `customer_org_audit_log_settings.include_employee_actions` is instead "+fmt.Sprintf("%t", environmentConfigResponse.OrgAuditLogIncludesEmployees),
			)
			return
		}
		if plan.CustomerOrgAuditLogSettings.IncludeApiKeyActions.ValueBool() != environmentConfigResponse.OrgAuditLogIncludesApiKeys {
			resp.Diagnostics.AddError(
				"Error updating organization configuration",
				"CustomerOrgAuditLogSettings.IncludeApiKeyActions failed to update. The `customer_org_audit_log_settings.include_api_key_actions` is instead "+fmt.Sprintf("%t", environmentConfigResponse.OrgAuditLogIncludesApiKeys),
			)
			return
		}
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a propelauth_organization_configuration resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *organizationConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state and read it into the model
	var state organizationConfigurationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// retrieve the environment config from PropelAuth
	environmentConfigResponse, err := r.client.GetEnvironmentConfig()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading PropelAuth organization configuration",
			"Could not read PropelAuth organization configuration: "+err.Error(),
		)
		return
	}

	// Save into the Terraform state only if the value is not null in Terraform.
	// Null, or unset values, in Terraform are left to be manually managed in the dashboard.
	if state.HasOrgs.ValueBoolPointer() != nil {
		state.HasOrgs = types.BoolValue(environmentConfigResponse.HasOrgs)
	}
	if state.OrgsMetaname.ValueString() != "" {
		state.OrgsMetaname = types.StringValue(environmentConfigResponse.OrgsMetaname)
	}
	if state.MaxNumOrgsUsersCanBeIn.ValueInt32Pointer() != nil {
		state.MaxNumOrgsUsersCanBeIn = types.Int32Value(environmentConfigResponse.MaxNumOrgsUsersCanBeIn)
	}
	if state.UsersCanCreateOrgs.ValueBoolPointer() != nil {
		state.UsersCanCreateOrgs = types.BoolValue(environmentConfigResponse.UsersCanCreateOrgs)
	}
	if state.UsersCanDeleteTheirOwnOrgs.ValueBoolPointer() != nil {
		state.UsersCanDeleteTheirOwnOrgs = types.BoolValue(environmentConfigResponse.UsersCanDeleteTheirOwnOrgs)
	}
	if state.UsersMustBeInAnOrganization.ValueBoolPointer() != nil {
		state.UsersMustBeInAnOrganization = types.BoolValue(environmentConfigResponse.UsersMustBeInAnOrganization)
	}
	if state.OrgsCanSetupSaml.ValueBoolPointer() != nil {
		state.OrgsCanSetupSaml = types.BoolValue(environmentConfigResponse.OrgsCanSetupSaml)
	}
	if state.UseOrgNameForSaml.ValueBoolPointer() != nil {
		state.UseOrgNameForSaml = types.BoolValue(environmentConfigResponse.UseOrgNameForSaml)
	}
	if state.DefaultToSamlLogin.ValueBoolPointer() != nil {
		state.DefaultToSamlLogin = types.BoolValue(environmentConfigResponse.DefaultToSamlLogin)
	}
	if state.SkipSamlRoleMappingStep.ValueBoolPointer() != nil {
		state.SkipSamlRoleMappingStep = types.BoolValue(environmentConfigResponse.SkipSamlRoleMappingStep)
	}
	if state.OrgsCanRequire2fa.ValueBoolPointer() != nil {
		state.OrgsCanRequire2fa = types.BoolValue(environmentConfigResponse.OrgsCanRequire2fa)
	}
	if state.CustomerOrgAuditLogSettings != nil {
		state.CustomerOrgAuditLogSettings.Enabled = types.BoolValue(environmentConfigResponse.OrgsCanViewOrgAuditLog)
		state.CustomerOrgAuditLogSettings.AllOrgsCanViewTheirAuditLog = types.BoolValue(environmentConfigResponse.AllOrgsCanViewOrgAuditLog)
		state.CustomerOrgAuditLogSettings.IncludeImpersonation = types.BoolValue(environmentConfigResponse.OrgAuditLogIncludesImpersonation)
		state.CustomerOrgAuditLogSettings.IncludeEmployeeActions = types.BoolValue(environmentConfigResponse.OrgAuditLogIncludesEmployees)
		state.CustomerOrgAuditLogSettings.IncludeApiKeyActions = types.BoolValue(environmentConfigResponse.OrgAuditLogIncludesApiKeys)
	}

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *organizationConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan organizationConfigurationResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the configuration in PropelAuth
	environmentConfigUpdate := propelauth.EnvironmentConfigUpdate{
		HasOrgs:                     plan.HasOrgs.ValueBoolPointer(),
		MaxNumOrgsUsersCanBeIn:      plan.MaxNumOrgsUsersCanBeIn.ValueInt32Pointer(),
		OrgsMetaname:                plan.OrgsMetaname.ValueString(),
		UsersCanCreateOrgs:          plan.UsersCanCreateOrgs.ValueBoolPointer(),
		UsersCanDeleteTheirOwnOrgs:  plan.UsersCanDeleteTheirOwnOrgs.ValueBoolPointer(),
		UsersMustBeInAnOrganization: plan.UsersMustBeInAnOrganization.ValueBoolPointer(),
		OrgsCanSetupSaml:            plan.OrgsCanSetupSaml.ValueBoolPointer(),
		UseOrgNameForSaml:           plan.UseOrgNameForSaml.ValueBoolPointer(),
		DefaultToSamlLogin:          plan.DefaultToSamlLogin.ValueBoolPointer(),
		SkipSamlRoleMappingStep:     plan.SkipSamlRoleMappingStep.ValueBoolPointer(),
		OrgsCanRequire2fa:           plan.OrgsCanRequire2fa.ValueBoolPointer(),
	}

	if plan.CustomerOrgAuditLogSettings != nil {
		environmentConfigUpdate.OrgsCanViewOrgAuditLog = plan.CustomerOrgAuditLogSettings.Enabled.ValueBoolPointer()
		environmentConfigUpdate.AllOrgsCanViewOrgAuditLog = plan.CustomerOrgAuditLogSettings.AllOrgsCanViewTheirAuditLog.ValueBoolPointer()
		environmentConfigUpdate.OrgAuditLogIncludesImpersonation = plan.CustomerOrgAuditLogSettings.IncludeImpersonation.ValueBoolPointer()
		environmentConfigUpdate.OrgAuditLogIncludesEmployees = plan.CustomerOrgAuditLogSettings.IncludeEmployeeActions.ValueBoolPointer()
		environmentConfigUpdate.OrgAuditLogIncludesApiKeys = plan.CustomerOrgAuditLogSettings.IncludeApiKeyActions.ValueBoolPointer()
	}

	environmentConfigResponse, err := r.client.UpdateEnvironmentConfig(&environmentConfigUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting organization configuration",
			"Could not set organization configuration, unexpected error: "+err.Error(),
		)
		return
	}

	// Check that all field were updated to the new value if not empty
	if plan.HasOrgs.ValueBoolPointer() != nil &&
		plan.HasOrgs.ValueBool() != environmentConfigResponse.HasOrgs {
		resp.Diagnostics.AddError(
			"Error updating organization configuration",
			"HasOrgs failed to update. The `allow_users_to_signup_with_personal_email` is instead "+fmt.Sprintf("%t", environmentConfigResponse.HasOrgs),
		)
		return
	}
	if plan.MaxNumOrgsUsersCanBeIn.ValueInt32Pointer() != nil &&
		plan.MaxNumOrgsUsersCanBeIn.ValueInt32() != environmentConfigResponse.MaxNumOrgsUsersCanBeIn {
		resp.Diagnostics.AddError(
			"Error updating organization configuration",
			"MaxNumOrgsUsersCanBeIn failed to update. The `max_num_orgs_users_can_be_in` is instead "+fmt.Sprintf("%d", environmentConfigResponse.MaxNumOrgsUsersCanBeIn),
		)
		return
	}
	if plan.OrgsMetaname.ValueString() != "" &&
		plan.OrgsMetaname.ValueString() != environmentConfigResponse.OrgsMetaname {
		resp.Diagnostics.AddError(
			"Error updating organization configuration",
			"OrgsMetaname failed to update. The `orgs_metaname` is instead "+environmentConfigResponse.OrgsMetaname,
		)
		return
	}
	if plan.UsersCanCreateOrgs.ValueBoolPointer() != nil &&
		plan.UsersCanCreateOrgs.ValueBool() != environmentConfigResponse.UsersCanCreateOrgs {
		resp.Diagnostics.AddError(
			"Error updating organization configuration",
			"UsersCanCreateOrgs failed to update. The `users_can_create_orgs` is instead "+fmt.Sprintf("%t", environmentConfigResponse.UsersCanCreateOrgs),
		)
		return
	}
	if plan.UsersCanDeleteTheirOwnOrgs.ValueBoolPointer() != nil &&
		plan.UsersCanDeleteTheirOwnOrgs.ValueBool() != environmentConfigResponse.UsersCanDeleteTheirOwnOrgs {
		resp.Diagnostics.AddError(
			"Error updating organization configuration",
			"UsersCanDeleteTheirOwnOrgs failed to update. The `users_can_delete_their_own_orgs` is instead "+fmt.Sprintf("%t", environmentConfigResponse.UsersCanDeleteTheirOwnOrgs),
		)
		return
	}
	if plan.UsersMustBeInAnOrganization.ValueBoolPointer() != nil &&
		plan.UsersMustBeInAnOrganization.ValueBool() != environmentConfigResponse.UsersMustBeInAnOrganization {
		resp.Diagnostics.AddError(
			"Error updating organization configuration",
			"UsersMustBeInAnOrganization failed to update. The `users_must_be_in_an_organization` is instead "+fmt.Sprintf("%t", environmentConfigResponse.UsersMustBeInAnOrganization),
		)
		return
	}
	if plan.OrgsCanSetupSaml.ValueBoolPointer() != nil &&
		plan.OrgsCanSetupSaml.ValueBool() != environmentConfigResponse.OrgsCanSetupSaml {
		resp.Diagnostics.AddError(
			"Error updating organization configuration",
			"OrgsCanSetupSaml failed to update. The `orgs_can_setup_saml` is instead "+fmt.Sprintf("%t", environmentConfigResponse.OrgsCanSetupSaml),
		)
		return
	}
	if plan.UseOrgNameForSaml.ValueBoolPointer() != nil &&
		plan.UseOrgNameForSaml.ValueBool() != environmentConfigResponse.UseOrgNameForSaml {
		resp.Diagnostics.AddError(
			"Error updating organization configuration",
			"UseOrgNameForSaml failed to update. The `use_org_name_for_saml` is instead "+fmt.Sprintf("%t", environmentConfigResponse.UseOrgNameForSaml),
		)
		return
	}
	if plan.DefaultToSamlLogin.ValueBoolPointer() != nil &&
		plan.DefaultToSamlLogin.ValueBool() != environmentConfigResponse.DefaultToSamlLogin {
		resp.Diagnostics.AddError(
			"Error updating organization configuration",
			"DefaultToSamlLogin failed to update. The `default_to_saml_login` is instead "+fmt.Sprintf("%t", environmentConfigResponse.DefaultToSamlLogin),
		)
		return
	}
	if plan.SkipSamlRoleMappingStep.ValueBoolPointer() != nil &&
		plan.SkipSamlRoleMappingStep.ValueBool() != environmentConfigResponse.SkipSamlRoleMappingStep {
		resp.Diagnostics.AddError(
			"Error updating organization configuration",
			"SkipSamlRoleMappingStep failed to update. The `skip_saml_role_mapping_step` is instead "+fmt.Sprintf("%t", environmentConfigResponse.SkipSamlRoleMappingStep),
		)
		return
	}
	if plan.OrgsCanRequire2fa.ValueBoolPointer() != nil &&
		plan.OrgsCanRequire2fa.ValueBool() != environmentConfigResponse.OrgsCanRequire2fa {
		resp.Diagnostics.AddError(
			"Error updating organization configuration",
			"OrgsCanRequire2fa failed to update. The `orgs_can_require_2fa` is instead "+fmt.Sprintf("%t", environmentConfigResponse.OrgsCanRequire2fa),
		)
		return
	}
	if plan.CustomerOrgAuditLogSettings != nil {
		if plan.CustomerOrgAuditLogSettings.Enabled.ValueBool() != environmentConfigResponse.OrgsCanViewOrgAuditLog {
			resp.Diagnostics.AddError(
				"Error updating organization configuration",
				"CustomerOrgAuditLogSettings.Enabled failed to update. The `customer_org_audit_log_settings.enabled` is instead "+fmt.Sprintf("%t", environmentConfigResponse.OrgsCanViewOrgAuditLog),
			)
			return
		}
		if plan.CustomerOrgAuditLogSettings.AllOrgsCanViewTheirAuditLog.ValueBool() != environmentConfigResponse.AllOrgsCanViewOrgAuditLog {
			resp.Diagnostics.AddError(
				"Error updating organization configuration",
				"CustomerOrgAuditLogSettings.AllOrgsCanViewTheirAuditLog failed to update. The `customer_org_audit_log_settings.all_orgs_can_view_their_audit_log` is instead "+fmt.Sprintf("%t", environmentConfigResponse.AllOrgsCanViewOrgAuditLog),
			)
			return
		}
		if plan.CustomerOrgAuditLogSettings.IncludeImpersonation.ValueBool() != environmentConfigResponse.OrgAuditLogIncludesImpersonation {
			resp.Diagnostics.AddError(
				"Error updating organization configuration",
				"CustomerOrgAuditLogSettings.IncludeImpersonation failed to update. The `customer_org_audit_log_settings.include_impersonation` is instead "+fmt.Sprintf("%t", environmentConfigResponse.OrgAuditLogIncludesImpersonation),
			)
			return
		}
		if plan.CustomerOrgAuditLogSettings.IncludeEmployeeActions.ValueBool() != environmentConfigResponse.OrgAuditLogIncludesEmployees {
			resp.Diagnostics.AddError(
				"Error updating organization configuration",
				"CustomerOrgAuditLogSettings.IncludeEmployeeActions failed to update. The `customer_org_audit_log_settings.include_employee_actions` is instead "+fmt.Sprintf("%t", environmentConfigResponse.OrgAuditLogIncludesEmployees),
			)
			return
		}
		if plan.CustomerOrgAuditLogSettings.IncludeApiKeyActions.ValueBool() != environmentConfigResponse.OrgAuditLogIncludesApiKeys {
			resp.Diagnostics.AddError(
				"Error updating organization configuration",
				"CustomerOrgAuditLogSettings.IncludeApiKeyActions failed to update. The `customer_org_audit_log_settings.include_api_key_actions` is instead "+fmt.Sprintf("%t", environmentConfigResponse.OrgAuditLogIncludesApiKeys),
			)
			return
		}
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a propelauth_organization_configuration resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *organizationConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "deleted a propelauth_organization_configuration resource")
}

func (r *organizationConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var state organizationConfigurationResourceModel

	// retrieve the environment config from PropelAuth
	environmentConfigResponse, err := r.client.GetEnvironmentConfig()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing PropelAuth organization configuration",
			"Could not read PropelAuth organization configuration: "+err.Error(),
		)
		return
	}

	// Save into the Terraform state all values from the dashboard.
	state.HasOrgs = types.BoolValue(environmentConfigResponse.HasOrgs)
	state.OrgsMetaname = types.StringValue(environmentConfigResponse.OrgsMetaname)
	state.MaxNumOrgsUsersCanBeIn = types.Int32Value(environmentConfigResponse.MaxNumOrgsUsersCanBeIn)
	state.UsersCanCreateOrgs = types.BoolValue(environmentConfigResponse.UsersCanCreateOrgs)
	state.UsersCanDeleteTheirOwnOrgs = types.BoolValue(environmentConfigResponse.UsersCanDeleteTheirOwnOrgs)
	state.UsersMustBeInAnOrganization = types.BoolValue(environmentConfigResponse.UsersMustBeInAnOrganization)
	state.OrgsCanSetupSaml = types.BoolValue(environmentConfigResponse.OrgsCanSetupSaml)
	state.UseOrgNameForSaml = types.BoolValue(environmentConfigResponse.UseOrgNameForSaml)
	state.DefaultToSamlLogin = types.BoolValue(environmentConfigResponse.DefaultToSamlLogin)
	state.SkipSamlRoleMappingStep = types.BoolValue(environmentConfigResponse.SkipSamlRoleMappingStep)
	state.OrgsCanRequire2fa = types.BoolValue(environmentConfigResponse.OrgsCanRequire2fa)
	state.CustomerOrgAuditLogSettings = &CustomerOrgAuditLogSettingsModel{
		Enabled:                     types.BoolValue(environmentConfigResponse.OrgsCanViewOrgAuditLog),
		AllOrgsCanViewTheirAuditLog: types.BoolValue(environmentConfigResponse.AllOrgsCanViewOrgAuditLog),
		IncludeImpersonation:        types.BoolValue(environmentConfigResponse.OrgAuditLogIncludesImpersonation),
		IncludeEmployeeActions:      types.BoolValue(environmentConfigResponse.OrgAuditLogIncludesEmployees),
		IncludeApiKeyActions:        types.BoolValue(environmentConfigResponse.OrgAuditLogIncludesApiKeys),
	}

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
