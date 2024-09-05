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
var _ resource.Resource = &apiKeySettingsResource{}
var _ resource.ResourceWithConfigure = &apiKeySettingsResource{}

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
}

type apiKeyConfigResourceModel struct {
	ExpirationOptions apiKeyExpirationOptionsResourceModel `tfsdk:"expiration_options"`
}

type apiKeyExpirationOptionsResourceModel struct {
	Options []types.String `tfsdk:"options"`
	Default types.String   `tfsdk:"default"`
}

func (r *apiKeySettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key_settings"
}

func (r *apiKeySettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Api Key Configurations. This is for configuring the API global settings for t.",
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
