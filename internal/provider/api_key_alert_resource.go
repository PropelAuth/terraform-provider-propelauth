package provider

import (
	"context"
	"fmt"

	"terraform-provider-propelauth/internal/propelauth"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &apiKeyAlertResource{}
var _ resource.ResourceWithConfigure = &apiKeyAlertResource{}
var _ resource.ResourceWithImportState = &apiKeyAlertResource{}

func NewApiKeyAlertResource() resource.Resource {
	return &apiKeyAlertResource{}
}

type apiKeyAlertResource struct {
	client *propelauth.PropelAuthClient
}

type apiKeyAlertResourceModel struct {
	AdvanceNoticeDays types.Int32 `tfsdk:"advance_notice_days"`
}

func (r *apiKeyAlertResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key_alert"
}

func (r *apiKeyAlertResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "API Key Alerting. This sets up alerting for end users emailing them when their API key is about to expire. " +
			"Note: API key alerts are only available on some pricing plans. These alerts are only sent for users in production " +
			"environments and can only be set up if you have a production environment.",
		Attributes: map[string]schema.Attribute{
			"advance_notice_days": schema.Int32Attribute{
				Required: true,
				Validators: []validator.Int32{
					int32validator.Between(1, 90),
				},
				Description: "The number of days before an API key expires by which time end users will receive an email alert.",
			},
		},
	}
}

func (r *apiKeyAlertResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *apiKeyAlertResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan apiKeyAlertResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the api key alert
	err := r.client.UpdateApiKeyAlert(plan.AdvanceNoticeDays.ValueInt32())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting propelauth_api_key_alert",
			"Could not set api key alert info for test environment, unexpected error: "+err.Error(),
		)
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a propelauth_api_key_alert resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *apiKeyAlertResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state apiKeyAlertResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// retrieve the api key alert from PropelAuth
	alertSettings, err := r.client.GetApiKeyAlert()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading PropelAuth  API Key Alert",
			"Could not read PropelAuth api key alert,: "+err.Error(),
		)
		return
	}
	if alertSettings.Enabled {
		state.AdvanceNoticeDays = types.Int32Value(alertSettings.AdvanceNoticeDays)
	} else {
		state.AdvanceNoticeDays = types.Int32Null()
	}

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *apiKeyAlertResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan apiKeyAlertResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the api key alert
	err := r.client.UpdateApiKeyAlert(plan.AdvanceNoticeDays.ValueInt32())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting propelauth_api_key_alert",
			"Could not set api key alert info for test environment, unexpected error: "+err.Error(),
		)
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a propelauth_api_key_alert resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *apiKeyAlertResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	err := r.client.DeleteApiKeyAlert()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting PropelAuth API Key Alert",
			"Could not delete be api key, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "deleted a propelauth_api_key_alert resource")
}

func (r *apiKeyAlertResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var state apiKeyAlertResourceModel

	// retrieve the environment config from PropelAuth
	apiKeyAlert, err := r.client.GetApiKeyAlert()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing PropelAuth API Key Alert",
			"Could not read PropelAuth API key alert: "+err.Error(),
		)
		return
	}

	// Save into the Terraform state all values from the dashboard.
	if apiKeyAlert.Enabled {
		state.AdvanceNoticeDays = types.Int32Value(apiKeyAlert.AdvanceNoticeDays)
	} else {
		resp.Diagnostics.AddError(
			"Error Importing PropelAuth API Key Alert",
			"Could not find any API key alert in your PropelAuth project to import.",
		)
		return
	}

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
