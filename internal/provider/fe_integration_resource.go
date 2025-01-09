package provider

import (
	"context"
	"fmt"
	"regexp"

	"terraform-provider-propelauth/internal/propelauth"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &feIntegrationResource{}
var _ resource.ResourceWithConfigure = &feIntegrationResource{}
var _ resource.ResourceWithValidateConfig = &feIntegrationResource{}
var _ resource.ResourceWithImportState = &feIntegrationResource{}

func NewFeIntegrationResource() resource.Resource {
	return &feIntegrationResource{}
}

// feIntegrationResource defines the resource implementation.
type feIntegrationResource struct {
	client *propelauth.PropelAuthClient
}

// feIntegrationResourceModel describes the resource data model.
type feIntegrationResourceModel struct {
	Environment           types.String                `tfsdk:"environment"`
	ApplicationUrl        types.String                `tfsdk:"application_url"`
	LoginRedirectPath     types.String                `tfsdk:"login_redirect_path"`
	LogoutRedirectPath    types.String                `tfsdk:"logout_redirect_path"`
	AdditionalFeLocations []additionalFeLocationModel `tfsdk:"additional_fe_locations"`
}

type additionalFeLocationModel struct {
	Domain            types.String `tfsdk:"domain"`
	AllowAnySubdomain types.Bool   `tfsdk:"allow_any_subdomain"`
}

func (r *feIntegrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fe_integration"
}

func (r *feIntegrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	urlPathRegex := regexp.MustCompile(`^\/[a-zA-Z0-9_\-\/]*$`)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Front-end Integration. This is for configuring the front-end integration for one of your project's environments.",
		Attributes: map[string]schema.Attribute{
			"environment": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("Test", "Staging", "Prod"),
				},
				Description: "The environment for which you are configuring the front-end integration. Accepted values are `Test`, `Staging`, and `Prod`.",
			},
			"application_url": schema.StringAttribute{
				Required: true,
				Description: "The URL of the application that will be integrated with PropelAuth. This is url is used in combination with " +
					"the `login_redirect_path` and `logout_redirect_path` to redirect users to your application after logging in or out. " +
					"`application_url` must be a valid URL and can only be set for `Test` environments. For `Staging` and `Prod` environments, " +
					"the `application_url` must be the verified domain for the environment or a subdomain of it. See " +
					" `propelauth_custom_domain_verification` resource for more information. Do not include any trailing path separator (`/`) " +
					"in your URL. For example, `https://any.subdomain.example.com` where `example.com` has been verified.",
			},
			"login_redirect_path": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(urlPathRegex, "Must be a valid URL path."),
				},
				Description: "The URL path to redirect users to after they log in. This path is appended to the `application_url` to form the " +
					"full URL. For example, if `application_url` is `https://example.com` and `login_redirect_path` is `/dashboard`, the " +
					"full URL will be `https://example.com/dashboard`.",
			},
			"logout_redirect_path": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(urlPathRegex, "Must be a valid URL path."),
				},
				Description: "The URL path to redirect users to after they log out. This path is appended to the `application_url` to form the " +
					"full URL. For example, if `application_url` is `https://example.com` and `logout_redirect_path` is `/goodbye`, the " +
					"full URL will be `https://example.com/goodbye`.",
			},
			"additional_fe_locations": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Additional front-end locations that are allowed to integrate with PropelAuth.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"domain": schema.StringAttribute{
							Required: true,
							Description: "A domain that will also be allowed to access user information. The domain must include a scheme. " +
								"For example, `https://example.com`.",
						},
						"allow_any_subdomain": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
							Description: "If true, any subdomain of the domain to integrate with PropelAuth is also allowed to access user info. " +
								"The default value is false.",
						},
					},
				},
			},
		},
	}
}

func (r *feIntegrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *feIntegrationResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var plan feIntegrationResourceModel

	// Read Terraform plan data into the model
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate the plan data
	isValidUrl, err := propelauth.IsValidUrlWithoutTrailingSlash(plan.ApplicationUrl.ValueString())
	if !isValidUrl {
		resp.Diagnostics.AddError(
			"Invalid application_url",
			"application_url must be a valid URL and cannot have a trailing slash. "+err.Error(),
		)
	}

	for i, location := range plan.AdditionalFeLocations {
		isValidUrl, err := propelauth.IsValidUrl(location.Domain.ValueString())
		if !isValidUrl {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Invalid domain for additional_fe_locations[%d]", i),
				"domain must be a valid URL. "+err.Error(),
			)
		}
	}
}

func (r *feIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan feIntegrationResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the front-end integration
	environment := plan.Environment.ValueString()
	update := convertPlanToUpdate(&plan)
	if environment == "Test" {
		_, err := r.client.UpdateTestFeIntegration(update)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error setting front-end intgeration info",
				"Could not set front-end integration info for test environment, unexpected error: "+err.Error(),
			)
			return
		}
	} else {
		_, err := r.client.UpdateLiveFeIntegration(environment, update)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error setting front-end intgeration info",
				"Could not set front-end integration info for live environment, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a propelauth_fe_integration resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *feIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state and read it into the model
	var state feIntegrationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// retrieve the front-end integration from PropelAuth
	if state.Environment.ValueString() == "Test" {
		fe_integration, err := r.client.GetTestFeIntegrationInfo()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading PropelAuth front-end integration",
				"Could not read PropelAuth front-end integration,: "+err.Error(),
			)
			return
		}
		updateStateForTestEnvironment(&state, fe_integration)
	} else {
		fe_integration, err := r.client.GetLiveFeIntegrationInfo(state.Environment.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading PropelAuth front-end integration",
				"Could not read PropelAuth front-end integration,: "+err.Error(),
			)
			return
		}
		updateStateForLiveEnvironment(&state, fe_integration)
	}

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *feIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan data into the model
	var plan feIntegrationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the front-end integration

	// Update the project info
	environment := plan.Environment.ValueString()
	update := convertPlanToUpdate(&plan)
	if environment == "Test" {
		_, err := r.client.UpdateTestFeIntegration(update)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error setting front-end intgeration info",
				"Could not set front-end integration info for test environment, unexpected error: "+err.Error(),
			)
			return
		}
	} else {
		_, err := r.client.UpdateLiveFeIntegration(environment, update)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error setting front-end intgeration info",
				"Could not set front-end integration info for live environment, unexpected error: "+err.Error(),
			)
			return
		}
	}

	tflog.Trace(ctx, "updated a propelauth_fe_integration resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *feIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "deleted a propelauth_fe_integration resource")
}

func (r *feIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("environment"), req, resp)
}

func convertPlanToUpdate(plan *feIntegrationResourceModel) propelauth.FeIntegrationUpdate {
	update := propelauth.FeIntegrationUpdate{
		ApplicationUrl:        plan.ApplicationUrl.ValueString(),
		LoginRedirectPath:     plan.LoginRedirectPath.ValueString(),
		LogoutRedirectPath:    plan.LogoutRedirectPath.ValueString(),
		AdditionalFeLocations: make([]propelauth.AdditionalFeLocation, len(plan.AdditionalFeLocations)),
	}

	for i, location := range plan.AdditionalFeLocations {
		update.AdditionalFeLocations[i] = propelauth.AdditionalFeLocation{
			Domain:            location.Domain.ValueString(),
			AllowAnySubdomain: location.AllowAnySubdomain.ValueBool(),
		}
	}

	return update
}

func updateStateForTestEnvironment(state *feIntegrationResourceModel, feIntegrationInfo *propelauth.TestFeIntegrationInfo) {
	state.ApplicationUrl = types.StringValue(feIntegrationInfo.TestEnvFeIntegrationApplicationUrl.ApplicationUrl)
	state.LoginRedirectPath = types.StringValue(feIntegrationInfo.LoginRedirectPath)
	state.LogoutRedirectPath = types.StringValue(feIntegrationInfo.LogoutRedirectPath)

	updateAdditionalLocationsInState(state, feIntegrationInfo.AdditionalFeLocations.AdditionalFeLocations)
}

func updateStateForLiveEnvironment(state *feIntegrationResourceModel, feIntegrationInfo *propelauth.FeIntegrationInfoForEnv) {
	state.ApplicationUrl = types.StringValue(feIntegrationInfo.ApplicationUrl)
	state.LoginRedirectPath = types.StringValue(feIntegrationInfo.LoginRedirectPath)
	state.LogoutRedirectPath = types.StringValue(feIntegrationInfo.LogoutRedirectPath)

	updateAdditionalLocationsInState(state, feIntegrationInfo.AdditionalFeLocations.AdditionalFeLocations)
}

func updateAdditionalLocationsInState(state *feIntegrationResourceModel, additionalLocations []propelauth.AdditionalFeLocation) {
	for i := range state.AdditionalFeLocations {
		if i >= len(additionalLocations) {
			state.AdditionalFeLocations = state.AdditionalFeLocations[:i]
			return
		}
		if state.AdditionalFeLocations[i].Domain.ValueString() != additionalLocations[i].Domain {
			state.AdditionalFeLocations[i].Domain = types.StringValue(additionalLocations[i].Domain)
		}
		if state.AdditionalFeLocations[i].AllowAnySubdomain.ValueBool() != additionalLocations[i].AllowAnySubdomain {
			state.AdditionalFeLocations[i].AllowAnySubdomain = types.BoolValue(additionalLocations[i].AllowAnySubdomain)
		}
	}

	if len(state.AdditionalFeLocations) < len(additionalLocations) {
		for i := len(state.AdditionalFeLocations); i < len(additionalLocations); i++ {
			state.AdditionalFeLocations = append(state.AdditionalFeLocations, additionalFeLocationModel{
				Domain:            types.StringValue(additionalLocations[i].Domain),
				AllowAnySubdomain: types.BoolValue(additionalLocations[i].AllowAnySubdomain),
			})
		}
	}
}
