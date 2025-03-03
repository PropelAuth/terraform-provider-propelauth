package provider

import (
	"context"
	"fmt"

	"terraform-provider-propelauth/internal/propelauth"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &apiKeySettingsResource{}
var _ resource.ResourceWithConfigure = &apiKeySettingsResource{}
var _ resource.ResourceWithImportState = &apiKeySettingsResource{}
var _ resource.ResourceWithValidateConfig = &apiKeySettingsResource{}

func NewApiKeySettingsResource() resource.Resource {
	return &apiKeySettingsResource{}
}

// apiKeySettingsResource defines the resource implementation.
type apiKeySettingsResource struct {
	client *propelauth.PropelAuthClient
}

// apiKeySettingsResourceModel describes the resource data model.
type apiKeySettingsResourceModel struct {
	PersonalApiKeysEnabled             types.Bool                 `tfsdk:"personal_api_keys_enabled"`
	OrgApiKeysEnabled                  types.Bool                 `tfsdk:"org_api_keys_enabled"`
	InvalidateOrgApiKeyUponUserRemoval types.Bool                 `tfsdk:"invalidate_org_api_key_upon_user_removal"`
	ApiKeyConfig                       *apiKeyConfigResourceModel `tfsdk:"api_key_config"`
	PersonalApiKeyRateLimit            *rateLimitConfigModel      `tfsdk:"personal_api_key_rate_limit"`
	OrgApiKeyRateLimit                 *rateLimitConfigModel      `tfsdk:"org_api_key_rate_limit"`
}

type apiKeyConfigResourceModel struct {
	ExpirationOptions apiKeyExpirationOptionsResourceModel `tfsdk:"expiration_options"`
}

type apiKeyExpirationOptionsResourceModel struct {
	Options []types.String `tfsdk:"options"`
	Default types.String   `tfsdk:"default"`
}

type rateLimitConfigModel struct {
	PeriodType     types.String `tfsdk:"period_type"`
	PeriodSize     types.Int32  `tfsdk:"period_size"`
	AllowPerPeriod types.Int64  `tfsdk:"allow_per_period"`
}

func (r *apiKeySettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key_settings"
}

func (r *apiKeySettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Api Key Configurations. API Key authentication allows you to create API Keys for your end users as well as " +
			"your organizations in order to protect requests they make to your product. This resource is for configuring the API " +
			"key settings in your project.\n\nNote: API Keys are only available for use in non-test environments for some " +
			"pricing plans.",
		Attributes: map[string]schema.Attribute{
			"personal_api_keys_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Allow users to create personal API keys. The default setting is false.",
			},
			"org_api_keys_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Allow users to create API keys for their organization. The default setting is false.",
			},
			"invalidate_org_api_key_upon_user_removal": schema.BoolAttribute{
				Optional: true,
				Description: "If true, invalidate org API keys when the user that created them is removed from the organization. " +
					"The default setting is false.",
			},
			"api_key_config": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "API Key Configuration. This is for setting the options available to your users when creating an API key.",
				Attributes: map[string]schema.Attribute{
					"expiration_options": schema.SingleNestedAttribute{
						Required:    true,
						Description: "API Key Expiration Options. This is for setting the options available to your users when creating an API key.",
						Attributes: map[string]schema.Attribute{
							"options": schema.ListAttribute{
								Required: true,
								Description: "The options available for the expiration of an API key. Valid values are " +
									"`TwoWeeks`, `OneMonth`, `ThreeMonths`, `SixMonths`, `OneYear`, and `Never`.",
								ElementType: types.StringType,
							},
							"default": schema.StringAttribute{
								Required: true,
								Description: "The default expiration option for an API key. Valid values are " +
									"`TwoWeeks`, `OneMonth`, `ThreeMonths`, `SixMonths`, `OneYear`, and `Never`.",
							},
						},
					},
				},
			},
			"personal_api_key_rate_limit": schema.SingleNestedAttribute{
				Optional: true,
				Description: "Personal API Key Rate Limit. This sets the rate limit that will be applied to validations on your " +
					"end users' personal API keys. This is calculated and applied per user for all keys they own as opposed to per key." +
					"Note: Rate limits are only available on some pricing plans.",
				Attributes: map[string]schema.Attribute{
					"period_type": schema.StringAttribute{
						Required: true,
						Description: "The unit of time for time for calculating and applying your rate limit. Valid values are " +
							"`seconds`, `minutes`, `hours`, or `days`.",
						Validators: []validator.String{
							stringvalidator.OneOf("seconds", "minutes", "hours", "days"),
						},
					},
					"period_size": schema.Int32Attribute{
						Required:    true,
						Description: "The number of `period_type` units for calculating and applying your rate limit.",
						Validators: []validator.Int32{
							int32validator.AtLeast(1),
						},
					},
					"allow_per_period": schema.Int64Attribute{
						Required:    true,
						Description: "The number of requests allowed per period.",
						Validators: []validator.Int64{
							int64validator.AtLeast(1),
						},
					},
				},
			},
			"org_api_key_rate_limit": schema.SingleNestedAttribute{
				Optional: true,
				Description: "Organization API Key Rate Limit. This sets the rate limit that will be applied to validations on your " +
					"end users' organizations' API keys. This is calculated and applied per organization for all keys the organization " +
					"owns as opposed to per key. Note: Rate limits are only available on some pricing plans.",
				Attributes: map[string]schema.Attribute{
					"period_type": schema.StringAttribute{
						Required: true,
						Description: "The unit of time for time for calculating and applying your rate limit. Valid values are " +
							"`seconds`, `minutes`, `hours`, or `days`.",
						Validators: []validator.String{
							stringvalidator.OneOf("seconds", "minutes", "hours", "days"),
						},
					},
					"period_size": schema.Int32Attribute{
						Required:    true,
						Description: "The number of `period_type` units for calculating and applying your rate limit.",
						Validators: []validator.Int32{
							int32validator.AtLeast(1),
						},
					},
					"allow_per_period": schema.Int64Attribute{
						Required:    true,
						Description: "The number of requests allowed per period.",
						Validators: []validator.Int64{
							int64validator.AtLeast(1),
						},
					},
				},
			},
		},
	}
}

func (r *apiKeySettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *apiKeySettingsResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var plan apiKeySettingsResourceModel

	// Read Terraform plan data into the model
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate the plan data
	if !plan.PersonalApiKeysEnabled.ValueBool() && plan.PersonalApiKeyRateLimit != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("personal_api_key_rate_limit"),
			"Rate limit cannot be set when personal API keys are disabled",
			"Cannot set `personal_api_key_rate_limit` when `personal_api_keys_enabled` is false",
		)
		return
	}

	if !plan.OrgApiKeysEnabled.ValueBool() && plan.OrgApiKeyRateLimit != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("org_api_key_rate_limit"),
			"Rate limit cannot be set when org API keys are disabled",
			"Cannot set `org_api_key_rate_limit` when `org_api_keys_enabled` is false",
		)
		return
	}
}

func (r *apiKeySettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan apiKeySettingsResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the configuration in PropelAuth
	environmentConfigUpdate := propelauth.EnvironmentConfigUpdate{
		PersonalApiKeysEnabled:              plan.PersonalApiKeysEnabled.ValueBoolPointer(),
		OrgApiKeysEnabled:                   plan.OrgApiKeysEnabled.ValueBoolPointer(),
		InvalidateOrgApiKeysUponUserRemoval: plan.InvalidateOrgApiKeyUponUserRemoval.ValueBoolPointer(),
	}

	if plan.ApiKeyConfig != nil {
		expirationOptions := make([]string, len(plan.ApiKeyConfig.ExpirationOptions.Options))
		for i, option := range plan.ApiKeyConfig.ExpirationOptions.Options {
			expirationOptions[i] = option.ValueString()
		}
		environmentConfigUpdate.ApiKeyConfig = &propelauth.ApiKeyConfig{
			ExpirationOptions: propelauth.ApiKeyExpirationOptionSettings{
				Options: propelauth.CreateApiKeyExpirationOptions(expirationOptions),
				Default: plan.ApiKeyConfig.ExpirationOptions.Default.ValueString(),
			},
		}
	}

	if plan.PersonalApiKeyRateLimit != nil {
		environmentConfigUpdate.PersonalApiKeyRateLimit = &propelauth.RateLimitConfig{
			PeriodType:     plan.PersonalApiKeyRateLimit.PeriodType.ValueString(),
			PeriodSize:     plan.PersonalApiKeyRateLimit.PeriodSize.ValueInt32(),
			AllowPerPeriod: plan.PersonalApiKeyRateLimit.AllowPerPeriod.ValueInt64(),
		}
	}

	if plan.OrgApiKeyRateLimit != nil {
		environmentConfigUpdate.OrgApiKeyRateLimit = &propelauth.RateLimitConfig{
			PeriodType:     plan.OrgApiKeyRateLimit.PeriodType.ValueString(),
			PeriodSize:     plan.OrgApiKeyRateLimit.PeriodSize.ValueInt32(),
			AllowPerPeriod: plan.OrgApiKeyRateLimit.AllowPerPeriod.ValueInt64(),
		}
	}

	environmentConfigResponse, err := r.client.UpdateEnvironmentConfig(&environmentConfigUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting api key settings",
			"Could not set api key settings, unexpected error: "+err.Error(),
		)
		return
	}

	// Check that fields were updated to the new value if not empty
	if plan.PersonalApiKeysEnabled.ValueBoolPointer() != nil &&
		plan.PersonalApiKeysEnabled.ValueBool() != environmentConfigResponse.PersonalApiKeysEnabled {
		resp.Diagnostics.AddError(
			"Error updating api key settings",
			"PersonalApiKeysEnabled failed to update. The `personal_api_keys_enabled` is instead "+fmt.Sprintf("%t", environmentConfigResponse.PersonalApiKeysEnabled),
		)
		return
	}
	if plan.OrgApiKeysEnabled.ValueBoolPointer() != nil &&
		plan.OrgApiKeysEnabled.ValueBool() != environmentConfigResponse.OrgApiKeysEnabled {
		resp.Diagnostics.AddError(
			"Error updating api key settings",
			"OrgApiKeysEnabled failed to update. The `org_api_keys_enabled` is instead "+fmt.Sprintf("%t", environmentConfigResponse.OrgApiKeysEnabled),
		)
		return
	}
	if plan.InvalidateOrgApiKeyUponUserRemoval.ValueBoolPointer() != nil &&
		plan.InvalidateOrgApiKeyUponUserRemoval.ValueBool() != environmentConfigResponse.InvalidateOrgApiKeyUponUserRemoval {
		resp.Diagnostics.AddError(
			"Error updating api key settings",
			"InvalidateOrgApiKeysUponUserRemoval failed to update. The `invalidate_org_api_key_upon_user_removal` is instead "+fmt.Sprintf("%t", environmentConfigResponse.InvalidateOrgApiKeyUponUserRemoval),
		)
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a propelauth_api_key_settings resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *apiKeySettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state and read it into the model
	var state apiKeySettingsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// retrieve the environment config from PropelAuth
	environmentConfigResponse, err := r.client.GetEnvironmentConfig()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading PropelAuth api key settings",
			"Could not read PropelAuth api key settings: "+err.Error(),
		)
		return
	}

	// Save into the Terraform state only if the value is not null in Terraform.
	// Null, or unset values, in Terraform are left to be manually managed in the dashboard.
	if state.PersonalApiKeysEnabled.ValueBoolPointer() != nil {
		state.PersonalApiKeysEnabled = types.BoolValue(environmentConfigResponse.PersonalApiKeysEnabled)
	}
	if state.OrgApiKeysEnabled.ValueBoolPointer() != nil {
		state.OrgApiKeysEnabled = types.BoolValue(environmentConfigResponse.OrgApiKeysEnabled)
	}
	if state.InvalidateOrgApiKeyUponUserRemoval.ValueBoolPointer() != nil {
		state.InvalidateOrgApiKeyUponUserRemoval = types.BoolValue(environmentConfigResponse.InvalidateOrgApiKeyUponUserRemoval)
	}
	if state.ApiKeyConfig != nil {
		state.ApiKeyConfig.ExpirationOptions.Default = types.StringValue(environmentConfigResponse.ApiKeyConfig.ExpirationOptions.Default)
		remoteApiKeyExpirationOptions := environmentConfigResponse.ApiKeyConfig.ExpirationOptions.GetApiKeyExpirationOptions()
		if diffInOptions(state.ApiKeyConfig.ExpirationOptions.Options, remoteApiKeyExpirationOptions) {
			state.ApiKeyConfig.ExpirationOptions.Options = make([]types.String, len(remoteApiKeyExpirationOptions))
			for i, option := range remoteApiKeyExpirationOptions {
				state.ApiKeyConfig.ExpirationOptions.Options[i] = types.StringValue(option)
			}
		}
	}
	if state.PersonalApiKeyRateLimit != nil {
		state.PersonalApiKeyRateLimit.PeriodType = types.StringValue(environmentConfigResponse.PersonalApiKeyRateLimit.PeriodType)
		state.PersonalApiKeyRateLimit.PeriodSize = types.Int32Value(environmentConfigResponse.PersonalApiKeyRateLimit.PeriodSize)
		state.PersonalApiKeyRateLimit.AllowPerPeriod = types.Int64Value(environmentConfigResponse.PersonalApiKeyRateLimit.AllowPerPeriod)
	}
	if state.OrgApiKeyRateLimit != nil {
		state.OrgApiKeyRateLimit.PeriodType = types.StringValue(environmentConfigResponse.OrgApiKeyRateLimit.PeriodType)
		state.OrgApiKeyRateLimit.PeriodSize = types.Int32Value(environmentConfigResponse.OrgApiKeyRateLimit.PeriodSize)
		state.OrgApiKeyRateLimit.AllowPerPeriod = types.Int64Value(environmentConfigResponse.OrgApiKeyRateLimit.AllowPerPeriod)
	}

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *apiKeySettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan apiKeySettingsResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the configuration in PropelAuth
	environmentConfigUpdate := propelauth.EnvironmentConfigUpdate{
		PersonalApiKeysEnabled:              plan.PersonalApiKeysEnabled.ValueBoolPointer(),
		OrgApiKeysEnabled:                   plan.OrgApiKeysEnabled.ValueBoolPointer(),
		InvalidateOrgApiKeysUponUserRemoval: plan.InvalidateOrgApiKeyUponUserRemoval.ValueBoolPointer(),
	}

	if plan.ApiKeyConfig != nil {
		expirationOptions := make([]string, len(plan.ApiKeyConfig.ExpirationOptions.Options))
		for i, option := range plan.ApiKeyConfig.ExpirationOptions.Options {
			expirationOptions[i] = option.ValueString()
		}
		environmentConfigUpdate.ApiKeyConfig = &propelauth.ApiKeyConfig{
			ExpirationOptions: propelauth.ApiKeyExpirationOptionSettings{
				Options: propelauth.CreateApiKeyExpirationOptions(expirationOptions),
				Default: plan.ApiKeyConfig.ExpirationOptions.Default.ValueString(),
			},
		}
	}

	if plan.PersonalApiKeyRateLimit != nil {
		environmentConfigUpdate.PersonalApiKeyRateLimit = &propelauth.RateLimitConfig{
			PeriodType:     plan.PersonalApiKeyRateLimit.PeriodType.ValueString(),
			PeriodSize:     plan.PersonalApiKeyRateLimit.PeriodSize.ValueInt32(),
			AllowPerPeriod: plan.PersonalApiKeyRateLimit.AllowPerPeriod.ValueInt64(),
		}
	}

	if plan.OrgApiKeyRateLimit != nil {
		environmentConfigUpdate.OrgApiKeyRateLimit = &propelauth.RateLimitConfig{
			PeriodType:     plan.OrgApiKeyRateLimit.PeriodType.ValueString(),
			PeriodSize:     plan.OrgApiKeyRateLimit.PeriodSize.ValueInt32(),
			AllowPerPeriod: plan.OrgApiKeyRateLimit.AllowPerPeriod.ValueInt64(),
		}
	}

	environmentConfigResponse, err := r.client.UpdateEnvironmentConfig(&environmentConfigUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting api key settings",
			"Could not set api key settings, unexpected error: "+err.Error(),
		)
		return
	}

	// Check that fields were updated to the new value if not empty
	if plan.PersonalApiKeysEnabled.ValueBoolPointer() != nil &&
		plan.PersonalApiKeysEnabled.ValueBool() != environmentConfigResponse.PersonalApiKeysEnabled {
		resp.Diagnostics.AddError(
			"Error updating api key settings",
			"PersonalApiKeysEnabled failed to update. The `personal_api_keys_enabled` is instead "+fmt.Sprintf("%t", environmentConfigResponse.PersonalApiKeysEnabled),
		)
		return
	}
	if plan.OrgApiKeysEnabled.ValueBoolPointer() != nil &&
		plan.OrgApiKeysEnabled.ValueBool() != environmentConfigResponse.OrgApiKeysEnabled {
		resp.Diagnostics.AddError(
			"Error updating api key settings",
			"OrgApiKeysEnabled failed to update. The `org_api_keys_enabled` is instead "+fmt.Sprintf("%t", environmentConfigResponse.OrgApiKeysEnabled),
		)
		return
	}
	if plan.InvalidateOrgApiKeyUponUserRemoval.ValueBoolPointer() != nil &&
		plan.InvalidateOrgApiKeyUponUserRemoval.ValueBool() != environmentConfigResponse.InvalidateOrgApiKeyUponUserRemoval {
		resp.Diagnostics.AddError(
			"Error updating api key settings",
			"InvalidateOrgApiKeysUponUserRemoval failed to update. The `invalidate_org_api_key_upon_user_removal` is instead "+fmt.Sprintf("%t", environmentConfigResponse.InvalidateOrgApiKeyUponUserRemoval),
		)
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "updated a propelauth_api_key_settings resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *apiKeySettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "deleted a propelauth_api_key_settings resource")
}

func (r *apiKeySettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var state apiKeySettingsResourceModel

	// retrieve the environment config from PropelAuth
	environmentConfigResponse, err := r.client.GetEnvironmentConfig()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing PropelAuth api key settings",
			"Could not read PropelAuth api key settings: "+err.Error(),
		)
		return
	}

	// Save into the Terraform state only if the value is not null in Terraform.
	// Null, or unset values, in Terraform are left to be manually managed in the dashboard.
	state.PersonalApiKeysEnabled = types.BoolValue(environmentConfigResponse.PersonalApiKeysEnabled)
	state.OrgApiKeysEnabled = types.BoolValue(environmentConfigResponse.OrgApiKeysEnabled)
	state.InvalidateOrgApiKeyUponUserRemoval = types.BoolValue(environmentConfigResponse.InvalidateOrgApiKeyUponUserRemoval)

	apiKeyConfig := &apiKeyConfigResourceModel{}
	apiKeyConfig.ExpirationOptions.Default = types.StringValue(environmentConfigResponse.ApiKeyConfig.ExpirationOptions.Default)
	remoteApiKeyExpirationOptions := environmentConfigResponse.ApiKeyConfig.ExpirationOptions.GetApiKeyExpirationOptions()
	apiKeyConfig.ExpirationOptions.Options = make([]types.String, len(remoteApiKeyExpirationOptions))
	for i, option := range remoteApiKeyExpirationOptions {
		apiKeyConfig.ExpirationOptions.Options[i] = types.StringValue(option)
	}
	state.ApiKeyConfig = apiKeyConfig

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func diffInOptions(stateOptions []types.String, remoteOptions []string) bool {
	if len(stateOptions) != len(remoteOptions) {
		return true
	}

	for _, option := range stateOptions {
		if !propelauth.Contains(remoteOptions, option.ValueString()) {
			return true
		}
	}

	return false
}
