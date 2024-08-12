package provider

import (
	"context"
	"fmt"

	"terraform-provider-propelauth/internal/propelauth"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &userPropertySettingsResource{}
var _ resource.ResourceWithConfigure   = &userPropertySettingsResource{}

func NewUserPropertySettingsResource() resource.Resource {
	return &userPropertySettingsResource{}
}

// userPropertySettingsResource defines the resource implementation.
type userPropertySettingsResource struct {
	client *propelauth.PropelAuthClient
}

// userPropertySettingsResourceModel describes the resource data model.
type userPropertySettingsResourceModel struct {
	NameProperty *namePropertyModel `tfsdk:"name_property"`
	MetadataProperty *metadataPropertyModel `tfsdk:"metadata_property"`
	UsernameProperty *usernamePropertyModel `tfsdk:"username_property"`
	PictureUrlProperty *pictureUrlPropertyModel `tfsdk:"picture_url_property"`
	TosProperty *tosPropertyModel `tfsdk:"tos_property"`
	ReferralSourceProperty *referralSourcePropertyModel `tfsdk:"referral_source_property"`
	PhoneNumberProperty *phoneNumberPropertyModel `tfsdk:"phone_number_property"`
	CustomProperties []customPropertyModel `tfsdk:"custom_properties"`
}

type namePropertyModel struct {
	InJwt types.Bool `tfsdk:"in_jwt"`
}

type metadataPropertyModel struct {
	InJwt types.Bool `tfsdk:"in_jwt"`
	CollectViaSaml types.Bool `tfsdk:"collect_via_saml"`
}

type usernamePropertyModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	InJwt types.Bool `tfsdk:"in_jwt"`
}

type pictureUrlPropertyModel struct {
	InJwt types.Bool `tfsdk:"in_jwt"`
}

type tosPropertyModel struct {
	InJwt types.Bool `tfsdk:"in_jwt"`
	Required types.Bool `tfsdk:"required"`
	RequiredBy types.Int64 `tfsdk:"required_by"`
	TosLinks []tosLinkModel `tfsdk:"tos_links"`
}

type tosLinkModel struct {
	Url types.String `tfsdk:"url"`
	Name types.String `tfsdk:"name"`
}

type referralSourcePropertyModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	InJwt types.Bool `tfsdk:"in_jwt"`
	Required types.Bool `tfsdk:"required"`
	RequiredBy types.Int64 `tfsdk:"required_by"`
	UserWriteable types.String `tfsdk:"user_writable"`
	Options []types.String `tfsdk:"options"`
	ShowInAccount types.Bool `tfsdk:"show_in_account"`
	CollectViaSaml types.Bool `tfsdk:"collect_via_saml"`
}

type phoneNumberPropertyModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	ShowInAccount types.Bool `tfsdk:"show_in_account"`
	CollectViaSaml types.Bool `tfsdk:"collect_via_saml"`
	Required types.Bool `tfsdk:"required"`
	RequiredBy types.Int64 `tfsdk:"required_by"`
	UserWritable types.String `tfsdk:"user_writable"`
	InJwt types.Bool `tfsdk:"in_jwt"`
}

type customPropertyModel struct {
	Name types.String `tfsdk:"name"`
	DisplayName types.String `tfsdk:"display_name"`
	FieldType types.String `tfsdk:"field_type"`
	Required types.Bool `tfsdk:"required"`
	RequiredBy types.Int64 `tfsdk:"required_by"`
	InJwt types.Bool `tfsdk:"in_jwt"`
	IsUserFacing types.Bool `tfsdk:"is_user_facing"`
	CollectOnSignup types.Bool `tfsdk:"collect_on_signup"`
	CollectViaSaml types.Bool `tfsdk:"collect_via_saml"`
	ShowInAccount types.Bool `tfsdk:"show_in_account"`
	UserWritable types.String `tfsdk:"user_writable"`
	EnumValues []types.String `tfsdk:"enum_values"`
}

func (r *userPropertySettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_property_settings"
}

func (r *userPropertySettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "User Property Settings. User properties are fields that you can use to store information about your users. " +
			"You can use them to collect information about your users on sign up, like their name or how they heard about " +
			"your product. You can also use them to store information about your users as they use your product, like their " +
			"subscription status, external IDs, or just arbitrary JSON data.",
		Attributes: map[string]schema.Attribute{
			"name_property": schema.SingleNestedAttribute{
				Optional: true,
				Description: "Settings for the user's name property. If no block is provided, the name property will be disabled.",
				Attributes: map[string]schema.Attribute{
					"in_jwt": inJwtAttribute(true),
				},
			},
			"metadata_property": schema.SingleNestedAttribute{
				Optional: true,
				Description: "Settings for the user's metadata property. If no block is provided, the metadata property will be disabled.",
				Attributes: map[string]schema.Attribute{
					"in_jwt": inJwtAttribute(true),
					"collect_via_saml": collectViaSamlAttribute(false),
				},
			},
			"username_property": schema.SingleNestedAttribute{
				Optional: true,
				Description: "Settings for the user's username property. If no block is provided, the username property will be disabled.",
				Attributes: map[string]schema.Attribute{
					"display_name": displayNameAttribute("Username"),
					"in_jwt": inJwtAttribute(true),
				},
			},
			"picture_url_property": schema.SingleNestedAttribute{
				Optional: true,
				Description: "Settings for the user's picture URL property. If no block is provided, the picture URL property will be disabled.",
				Attributes: map[string]schema.Attribute{
					"in_jwt": inJwtAttribute(true),
				},
			},
			"tos_property": schema.SingleNestedAttribute{
				Optional: true,
				Description: "Settings for the user's Terms of Service property. If no block is provided, the terms of service property will be disabled.",
				Attributes: map[string]schema.Attribute{
					"in_jwt": inJwtAttribute(false),
					"required": requiredAttribute(true),
					"required_by": requiredByAttribute(0),
					"tos_links": schema.ListNestedAttribute{
						Optional: true,
						Description: "A list of Terms of Service links. Each link must have a URL and a name.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"url": schema.StringAttribute{
									Required: true,
									Description: "The URL of the Terms of Service link.",
								},
								"name": schema.StringAttribute{
									Required: true,
									Description: "The name of the Terms of Service link.",
								},
							},
						},
					},
				},
			},
			"referral_source_property": schema.SingleNestedAttribute{
				Optional: true,
				Description: "Settings for the user's referral source property. If no block is provided, the referral source property will be disabled.",
				Attributes: map[string]schema.Attribute{
					"display_name": displayNameAttribute("How did you hear about us?"),
					"in_jwt": inJwtAttribute(true),
					"required": requiredAttribute(true),
					"required_by": requiredByAttribute(0),
					"user_writable": userWriteableAttribute("WriteIfUnset"),
					"show_in_account": showInAccountAttribute(false),
					"collect_via_saml": collectViaSamlAttribute(false),
					"options": schema.ListAttribute{
						Optional: true,
						Computed: true,
						ElementType: types.StringType,
						Default: listdefault.StaticValue(types.ListValueMust(
								types.StringType,
								[]attr.Value{
									types.StringValue("Search engine"),
									types.StringValue("Recommendation"),
									types.StringValue("Social media"),
									types.StringValue("Blog post"),
									types.StringValue("Other"),
								},
						)),
						Description: "A list of options for the referral source property. If this is unset, the default options " +
							"will be used. These are `Search engine`, `Recommendation`, `Social media`, `Blog post`, `Other`.",
					},
				},
			},
			"phone_number_property": schema.SingleNestedAttribute{
				Optional: true,
				Description: "Settings for the user's phone number property. If no block is provided, the phone number property will be disabled.",
				Attributes: map[string]schema.Attribute{
					"display_name": displayNameAttribute("Phone number"),
					"show_in_account": showInAccountAttribute(false),
					"collect_via_saml": collectViaSamlAttribute(false),
					"required": requiredAttribute(true),
					"required_by": requiredByAttribute(0),
					"user_writable": userWriteableAttribute("WriteIfUnset"),
					"in_jwt": inJwtAttribute(false),
				},
			},
			"custom_properties": schema.ListNestedAttribute{
				Optional: true,
				Description: "Custom properties for the user. If no blocks are provided, no custom properties will be enabled.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
							Description: "The field name used to identify the property in the API and SDKs (e.g. external_id). " +
								"It cannot be changed after creation.",
						},
						"display_name": schema.StringAttribute{
							Required: true,
							Description: "The field name users see in the UI for the property.",
						},
						"field_type": schema.StringAttribute{
							Required: true,
							Description: "The type of the field. Accepted values are `Checkbox`, `Date`, `Enum`, " +
								"`Integer`, `Json`, `LongText`, `Text`, `Toggle`, and `Url`. Once set, this cannot be changed.",
							Validators: []validator.String{
								stringvalidator.OneOf("Checkbox", "Date", "Enum", "Integer", "Json", "LongText", "Text", "Toggle", "Url"),
							},
						},
						"required": requiredAttribute(true),
						"required_by": requiredByAttribute(0),
						"in_jwt": inJwtAttribute(true),
						"is_user_facing": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default: booldefault.StaticBool(true),
							Description: "Whether the property should be displayed in the user's account page hosted by PropelAuth. " +
								"The default value is `false`.",
						},
						"collect_on_signup": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default: booldefault.StaticBool(true),
							Description: "Whether the property should be collected from new users during the sign up flow. " +
								"The default value is `true`.",
						},
						"collect_via_saml": collectViaSamlAttribute(false),
						"show_in_account": showInAccountAttribute(true),
						"user_writable": userWriteableAttribute("Write"),
						"enum_values": schema.ListAttribute{
							Optional: true,
							Description: "A list of possible values for the property. This is only required for the `Enum` field type.",
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (r *userPropertySettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *userPropertySettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userPropertySettingsResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

	// Fetch the current user property settings from PropelAuth
	userPropertySettings, err := r.client.GetUserProperties()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading PropelAuth user properties settings",
			"Could not read PropelAuth user properties settings: " + err.Error(),
		)
		return
	}

	// Update the configuration in PropelAuth
	UpdateDefaultPropertiesFromPlan(&plan, userPropertySettings)
	UpdateCustomPropertiesFromPlan(&plan, userPropertySettings)

    _, err = r.client.UpdateUserProperties(userPropertySettings)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting user properties settings",
			"Could not set user properties settings, unexpected error: "+err.Error(),
		)
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a propelauth_user_properties_settings resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *userPropertySettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state and read it into the model
	var state userPropertySettingsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch the current user property settings from PropelAuth
	userPropertySettings, err := r.client.GetUserProperties()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading PropelAuth user properties settings",
			"Could not read PropelAuth user properties settings: " + err.Error(),
		)
		return
	}

	// Update the configuration in PropelAuth
	if userPropertySettings.NamePropertyEnabled() {
		namePropertySettings := userPropertySettings.GetNamePropertySettings()
		state.NameProperty = &namePropertyModel{
			InJwt: types.BoolValue(namePropertySettings.InJwt),
		}
	} else {
		state.NameProperty = nil
	}

	if userPropertySettings.MetadataPropertyEnabled() {
		metadataPropertySettings := userPropertySettings.GetMetadataPropertySettings()
		state.MetadataProperty = &metadataPropertyModel{
			InJwt: types.BoolValue(metadataPropertySettings.InJwt),
			CollectViaSaml: types.BoolValue(metadataPropertySettings.CollectViaSaml),
		}
	} else {
		state.MetadataProperty = nil
	}

	if userPropertySettings.UsernamePropertyEnabled() {
		usernamePropertySettings := userPropertySettings.GetUsernamePropertySettings()
		state.UsernameProperty = &usernamePropertyModel{
			DisplayName: types.StringValue(usernamePropertySettings.DisplayName),
			InJwt: types.BoolValue(usernamePropertySettings.InJwt),
		}
	} else {
		state.UsernameProperty = nil
	}

	if userPropertySettings.PictureUrlPropertyEnabled() {
		pictureUrlPropertySettings := userPropertySettings.GetPictureUrlPropertySettings()
		state.PictureUrlProperty = &pictureUrlPropertyModel{
			InJwt: types.BoolValue(pictureUrlPropertySettings.InJwt),
		}
	} else {
		state.PictureUrlProperty = nil
	}

	if userPropertySettings.PhoneNumberPropertyEnabled() {
		phoneNumberPropertySettings := userPropertySettings.GetPhoneNumberPropertySettings()
		state.PhoneNumberProperty = &phoneNumberPropertyModel{
			DisplayName: types.StringValue(phoneNumberPropertySettings.DisplayName),
			ShowInAccount: types.BoolValue(phoneNumberPropertySettings.ShowInAccount),
			CollectViaSaml: types.BoolValue(phoneNumberPropertySettings.CollectViaSaml),
			Required: types.BoolValue(phoneNumberPropertySettings.Required),
			RequiredBy: types.Int64Value(phoneNumberPropertySettings.RequiredBy),
			UserWritable: types.StringValue(phoneNumberPropertySettings.UserWritable),
			InJwt: types.BoolValue(phoneNumberPropertySettings.InJwt),
		}
	} else {
		state.PhoneNumberProperty = nil
	}

	if userPropertySettings.TosPropertyEnabled() {
		tosPropertySettings := userPropertySettings.GetTosPropertySettings()
		tosLinks := make([]tosLinkModel, len(tosPropertySettings.TosLinks))
		for i, tosLink := range tosPropertySettings.TosLinks {
			tosLinks[i] = tosLinkModel{
				Url: types.StringValue(tosLink.Url),
				Name: types.StringValue(tosLink.Name),
			}
		}
		state.TosProperty = &tosPropertyModel{
			InJwt: types.BoolValue(tosPropertySettings.InJwt),
			Required: types.BoolValue(tosPropertySettings.Required),
			RequiredBy: types.Int64Value(tosPropertySettings.RequiredBy),
			TosLinks: tosLinks,
		}
	} else {
		state.TosProperty = nil
	}

	if userPropertySettings.ReferralSourcePropertyEnabled() {
		referralSourcePropertySettings := userPropertySettings.GetReferralSourcePropertySettings()
		options := make([]types.String, len(referralSourcePropertySettings.Options))
		for i, option := range referralSourcePropertySettings.Options {
			options[i] = types.StringValue(option)
		}
		state.ReferralSourceProperty = &referralSourcePropertyModel{
			DisplayName: types.StringValue(referralSourcePropertySettings.DisplayName),
			InJwt: types.BoolValue(referralSourcePropertySettings.InJwt),
			Required: types.BoolValue(referralSourcePropertySettings.Required),
			RequiredBy: types.Int64Value(referralSourcePropertySettings.RequiredBy),
			UserWriteable: types.StringValue(referralSourcePropertySettings.UserWritable),
			Options: options,
			ShowInAccount: types.BoolValue(referralSourcePropertySettings.ShowInAccount),
			CollectViaSaml: types.BoolValue(referralSourcePropertySettings.CollectViaSaml),
		}
	} else {
		state.ReferralSourceProperty = nil
	}

	reconcileCustomProperties(&state, userPropertySettings)

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *userPropertySettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userPropertySettingsResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

	// Fetch the current user property settings from PropelAuth
	userPropertySettings, err := r.client.GetUserProperties()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading PropelAuth user properties settings",
			"Could not read PropelAuth user properties settings: " + err.Error(),
		)
		return
	}

	// Update the configuration in PropelAuth
	UpdateDefaultPropertiesFromPlan(&plan, userPropertySettings)
	UpdateCustomPropertiesFromPlan(&plan, userPropertySettings)

    _, err = r.client.UpdateUserProperties(userPropertySettings)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting user properties settings",
			"Could not set user properties settings, unexpected error: "+err.Error(),
		)
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "updated a propelauth_user_properties_settings resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *userPropertySettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "deleted a propelauth_user_properties_settings resource")
}

func UpdateDefaultPropertiesFromPlan(plan *userPropertySettingsResourceModel, userPropertySettings *propelauth.UserProperties) {
	if plan.NameProperty != nil {
		userPropertySettings.UpdateAndEnableNameProperty(propelauth.NamePropertySettings{
			InJwt: plan.NameProperty.InJwt.ValueBool(),
		})
	} else {
		userPropertySettings.DisableNameProperty()
	}

	if plan.MetadataProperty != nil {
		userPropertySettings.UpdateAndEnableMetadataProperty(propelauth.MetadataPropertySettings{
			InJwt: plan.MetadataProperty.InJwt.ValueBool(),
			CollectViaSaml: plan.MetadataProperty.CollectViaSaml.ValueBool(),
		})
	} else {
		userPropertySettings.DisableMetadataProperty()
	}

	if plan.UsernameProperty != nil {
		userPropertySettings.UpdateAndEnableUsernameProperty(propelauth.UsernamePropertySettings{
			DisplayName: plan.UsernameProperty.DisplayName.ValueString(),
			InJwt: plan.UsernameProperty.InJwt.ValueBool(),
		})
	} else {
		userPropertySettings.DisableUsernameProperty()
	}

	if plan.PictureUrlProperty != nil {
		userPropertySettings.UpdateAndEnablePictureUrlProperty(propelauth.PictureUrlPropertySettings{
			InJwt: plan.PictureUrlProperty.InJwt.ValueBool(),
		})
	} else {
		userPropertySettings.DisablePictureUrlProperty()
	}

	if plan.PhoneNumberProperty != nil {
		userPropertySettings.UpdateAndEnablePhoneNumberProperty(propelauth.PhoneNumberPropertySettings{
			DisplayName: plan.PhoneNumberProperty.DisplayName.ValueString(),
			ShowInAccount: plan.PhoneNumberProperty.ShowInAccount.ValueBool(),
			CollectViaSaml: plan.PhoneNumberProperty.CollectViaSaml.ValueBool(),
			Required: plan.PhoneNumberProperty.Required.ValueBool(),
			RequiredBy: plan.PhoneNumberProperty.RequiredBy.ValueInt64(),
			UserWritable: plan.PhoneNumberProperty.UserWritable.ValueString(),
			InJwt: plan.PhoneNumberProperty.InJwt.ValueBool(),
		})
	} else {
		userPropertySettings.DisablePhoneNumberProperty()
	}

	if plan.TosProperty != nil {
		tosPropertySettings := propelauth.TosPropertySettings{
			InJwt: plan.TosProperty.InJwt.ValueBool(),
			Required: plan.TosProperty.Required.ValueBool(),
			RequiredBy: plan.TosProperty.RequiredBy.ValueInt64(),
			TosLinks: make([]propelauth.TosLink, len(plan.TosProperty.TosLinks)),
		}
		for i, tosLink := range plan.TosProperty.TosLinks {
			tosPropertySettings.TosLinks[i] = propelauth.TosLink{
				Url: tosLink.Url.ValueString(),
				Name: tosLink.Name.ValueString(),
			}
		}
		userPropertySettings.UpdateAndEnableTosProperty(tosPropertySettings)
	} else {
		userPropertySettings.DisableTosProperty()
	}

	if plan.ReferralSourceProperty != nil {
		referralSourcePropertySettings := propelauth.ReferralSourcePropertySettings{
			DisplayName: plan.ReferralSourceProperty.DisplayName.ValueString(),
			InJwt: plan.ReferralSourceProperty.InJwt.ValueBool(),
			Required: plan.ReferralSourceProperty.Required.ValueBool(),
			RequiredBy: plan.ReferralSourceProperty.RequiredBy.ValueInt64(),
			UserWritable: plan.ReferralSourceProperty.UserWriteable.ValueString(),
			ShowInAccount: plan.ReferralSourceProperty.ShowInAccount.ValueBool(),
			CollectViaSaml: plan.ReferralSourceProperty.CollectViaSaml.ValueBool(),
			Options: make([]string, len(plan.ReferralSourceProperty.Options)),
		}
		for i, option := range plan.ReferralSourceProperty.Options {
			referralSourcePropertySettings.Options[i] = option.ValueString()
		}
		userPropertySettings.UpdateAndEnableReferralSourceProperty(referralSourcePropertySettings)
	} else {
		userPropertySettings.DisableReferralSourceProperty()
	}
}

func UpdateCustomPropertiesFromPlan(plan *userPropertySettingsResourceModel, userPropertySettings *propelauth.UserProperties) {
	customPropertyUpdates := make([]propelauth.CustomPropertySettings, len(plan.CustomProperties))
	for i, customProperty := range plan.CustomProperties {
		customPropertyUpdate := propelauth.CustomPropertySettings{
			Name: customProperty.Name.ValueString(),
			DisplayName: customProperty.DisplayName.ValueString(),
			FieldType: customProperty.FieldType.ValueString(),
			Required: customProperty.Required.ValueBool(),
			RequiredBy: customProperty.RequiredBy.ValueInt64(),
			InJwt: customProperty.InJwt.ValueBool(),
			IsUserFacing: customProperty.IsUserFacing.ValueBool(),
			CollectOnSignup: customProperty.CollectOnSignup.ValueBool(),
			CollectViaSaml: customProperty.CollectViaSaml.ValueBool(),
			ShowInAccount: customProperty.ShowInAccount.ValueBool(),
			UserWritable: customProperty.UserWritable.ValueString(),
			EnumValues: make([]string, len(customProperty.EnumValues)),
		}
		for j, enumValue := range customProperty.EnumValues {
			customPropertyUpdate.EnumValues[j] = enumValue.ValueString()
		}
		customPropertyUpdates[i] = customPropertyUpdate
	}
	
	for _, customPropertyUpdate := range customPropertyUpdates {
		userPropertySettings.UpsertCustomProperty(customPropertyUpdate)
	}
	userPropertySettings.DisableDroppedCustomProperties(customPropertyUpdates)
}

func reconcileCustomProperties(state *userPropertySettingsResourceModel, userPropertySettings *propelauth.UserProperties) {
	for i, customPropertyInState := range state.CustomProperties {
		activeCustomProperty, ok := userPropertySettings.GetEnabledCustomProperty(customPropertyInState.Name.ValueString())
		if !ok {
			customPropertyInState = customPropertyModel{}
		}
		convertedCustomPropertyInState := convertCustomPropertyFromModel(customPropertyInState)

		if !convertedCustomPropertyInState.IsEqual(activeCustomProperty) {
			state.CustomProperties[i] = convertCustomPropertyToModel(&activeCustomProperty)
		}
	}

	customPropertyNamesInState := make([]string, len(state.CustomProperties))
	for i, customProperty := range state.CustomProperties {
		customPropertyNamesInState[i] = customProperty.Name.ValueString()
	}

	hangingCustomProperties := userPropertySettings.GetHangingCustomProperties(customPropertyNamesInState)
	convertedCustomPropertiesFromHanging := make([]customPropertyModel, len(hangingCustomProperties))

	for _, hangingCustomProperty := range hangingCustomProperties {
		convertedCustomProperty := convertCustomPropertyToModel(&hangingCustomProperty)
		convertedCustomPropertiesFromHanging = append(convertedCustomPropertiesFromHanging, convertedCustomProperty)
	}

	state.CustomProperties = append(state.CustomProperties, convertedCustomPropertiesFromHanging...)
}

func convertCustomPropertyToModel(customProperty *propelauth.CustomPropertySettings) customPropertyModel {
	customPropertyModel := customPropertyModel{
		Name: types.StringValue(customProperty.Name),
		DisplayName: types.StringValue(customProperty.DisplayName),
		FieldType: types.StringValue(customProperty.FieldType),
		Required: types.BoolValue(customProperty.Required),
		RequiredBy: types.Int64Value(customProperty.RequiredBy),
		InJwt: types.BoolValue(customProperty.InJwt),
		IsUserFacing: types.BoolValue(customProperty.IsUserFacing),
		CollectOnSignup: types.BoolValue(customProperty.CollectOnSignup),
		CollectViaSaml: types.BoolValue(customProperty.CollectViaSaml),
		ShowInAccount: types.BoolValue(customProperty.ShowInAccount),
		UserWritable: types.StringValue(customProperty.UserWritable),
	}

	if customProperty.FieldType == "Enum" {
		enumValues := make([]types.String, len(customProperty.EnumValues))
		for i, enumValue := range customProperty.EnumValues {
			enumValues[i] = types.StringValue(enumValue)
		}
		customPropertyModel.EnumValues = enumValues
	}

	return customPropertyModel
}

func convertCustomPropertyFromModel(customPropertyModel customPropertyModel) propelauth.CustomPropertySettings {
	customProperty := propelauth.CustomPropertySettings{
		Name: customPropertyModel.Name.ValueString(),
		DisplayName: customPropertyModel.DisplayName.ValueString(),
		FieldType: customPropertyModel.FieldType.ValueString(),
		Required: customPropertyModel.Required.ValueBool(),
		RequiredBy: customPropertyModel.RequiredBy.ValueInt64(),
		InJwt: customPropertyModel.InJwt.ValueBool(),
		IsUserFacing: customPropertyModel.IsUserFacing.ValueBool(),
		CollectOnSignup: customPropertyModel.CollectOnSignup.ValueBool(),
		CollectViaSaml: customPropertyModel.CollectViaSaml.ValueBool(),
		ShowInAccount: customPropertyModel.ShowInAccount.ValueBool(),
		UserWritable: customPropertyModel.UserWritable.ValueString(),
	}

	if customPropertyModel.FieldType.ValueString() == "Enum" {
		enumValues := make([]string, len(customPropertyModel.EnumValues))
		for i, enumValue := range customPropertyModel.EnumValues {
			enumValues[i] = enumValue.ValueString()
		}
		customProperty.EnumValues = enumValues
	}

	return customProperty
}

func inJwtAttribute(defaultValue bool) schema.Attribute {
	return schema.BoolAttribute{
		Optional: true,
		Computed: true,
		Default: booldefault.StaticBool(defaultValue),
		Description: fmt.Sprintf("Whether the property should be included in the user token. " +
			"The default value is `%v`.", defaultValue),
	}
}

func collectViaSamlAttribute(defaultValue bool) schema.Attribute {
	return schema.BoolAttribute{
		Optional: true,
		Computed: true,
		Default: booldefault.StaticBool(defaultValue),
		Description: fmt.Sprintf("Whether the property should be collected for users during the enterprise SSO login flow. " +
			"The default value is `%v`.", defaultValue),
	}
}

func displayNameAttribute(defaultValue string) schema.Attribute {
	return schema.StringAttribute{
		Optional: true,
		Computed: true,
		Default: stringdefault.StaticString(defaultValue),
		Description: fmt.Sprintf("The field name users see in the UI for the property. " +
			"The default value is `%v`.", defaultValue),
	}
}

func requiredAttribute(defaultValue bool) schema.Attribute {
	return schema.BoolAttribute{
		Optional: true,
		Computed: true,
		Default: booldefault.StaticBool(defaultValue),
		Description: fmt.Sprintf("Whether the property is required for users. " +
			"The default value is `%v`.", defaultValue),
	}
}

func requiredByAttribute(defaultValue int64) schema.Attribute {
	return schema.Int64Attribute{
		Optional: true,
		Computed: true,
		Default: int64default.StaticInt64(defaultValue),
		Description: fmt.Sprintf("In epoch time. Only accounts created after this time are required to " +
			"provide this field. For example, a value of 0 means all accounts are required to provide " +
			"this field. The default value is `%v`.", defaultValue),
	}
}

func userWriteableAttribute(defaultValue string) schema.Attribute {
	return schema.StringAttribute{
		Optional: true,
		Computed: true,
		Default: stringdefault.StaticString(defaultValue),
		Validators: []validator.String{
			stringvalidator.OneOf("Write", "Read", "WriteIfUnset"),
		},
		Description: fmt.Sprintf("This setting determines whether the user can edit the value of the property " +
			"and how many times. Options are `Write`, `Read`, and `WriteIfUnset`. The default value is `%v`", 
			defaultValue),
	}
}

func showInAccountAttribute(defaultValue bool) schema.Attribute {
	return schema.BoolAttribute{
		Optional: true,
		Computed: true,
		Default: booldefault.StaticBool(defaultValue),
		Description: fmt.Sprintf("Whether the property should be displayed in the user's account page hosted by " +
			"PropelAuth. The default value is `%v`.", defaultValue),
	}
}
