package provider

import (
	"context"
	"fmt"
	"time"

	"terraform-provider-propelauth/internal/propelauth"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
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
var _ resource.Resource = &customDomainVerificationResource{}
var _ resource.ResourceWithConfigure = &customDomainVerificationResource{}
var _ resource.ResourceWithImportState = &customDomainVerificationResource{}

func NewCustomDomainVerificationResource() resource.Resource {
	return &customDomainVerificationResource{}
}

// customDomainVerificationResource defines the resource implementation.
type customDomainVerificationResource struct {
	client *propelauth.PropelAuthClient
}

// projectInfoResourceModel describes the resource data model.
type customDomainVerificationResourceModel struct {
	Environment types.String   `tfsdk:"environment"`
	Domain      types.String   `tfsdk:"domain"`
	Timeouts    timeouts.Value `tfsdk:"timeouts"`
}

func (r *customDomainVerificationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_domain_verification"
}

func (r *customDomainVerificationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Custom Domain Verification resource. This is for verifying a custom domain for Production or Staging.",
		Attributes: map[string]schema.Attribute{
			"environment": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("Staging", "Prod"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The environment for which you are configuring the custom domain. Accepted values are `Staging`, `Prod`.",
			},
			"domain": schema.StringAttribute{
				Required:    true,
				Description: "The domain to verify.",
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
			}),
		},
	}
}

func (r *customDomainVerificationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *customDomainVerificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan customDomainVerificationResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, err := plan.Timeouts.Create(ctx, 5*time.Minute)
	if err != nil {
		resp.Diagnostics.AddError("Error creating a timeout", "Could not create a timeout for the custom domain verification.")
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	// Retry interval
	retryInterval := 30 * time.Second

	// Verify the custom domain with retries
	environment := plan.Environment.ValueString()
	for {
		select {
		case <-ctx.Done():
			resp.Diagnostics.AddError("Timeout exceeded", "Could not verify custom domain within the timeout. It can take a few minutes for the DNS records to propagate, please verify the records are set and try again.")
			return
		default:
			verificationErr := r.client.VerifyCustomDomainInfo(environment, false)
			if verificationErr == nil {
				// Verification successful
				// Set the data from the state into the response
				resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
				return
			}

			// Log the retry attempt
			tflog.Warn(ctx, "Unable to verify the custom domain. It can take a few minutes for the DNS records to propagate. Retrying in 30 seconds...")

			// Wait for the retry interval before the next attempt
			time.Sleep(retryInterval)
		}
	}

}

func (r *customDomainVerificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read the data from the state
	var state customDomainVerificationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a custom_domain resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *customDomainVerificationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan customDomainVerificationResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state customDomainVerificationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateTimeout, err := plan.Timeouts.Update(ctx, 5*time.Minute)
	if err != nil {
		resp.Diagnostics.AddError("Error creating a timeout", "Could not create a timeout for the custom domain verification.")
		return
	}

	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	// Retry interval
	retryInterval := 30 * time.Second

	// Verify the custom domain with retries
	environment := plan.Environment.ValueString()
	isSwitching := plan.Domain.ValueString() != state.Domain.ValueString()
	for {
		select {
		case <-ctx.Done():
			resp.Diagnostics.AddError("Timeout exceeded", "Could not verify custom domain within the timeout. It can take a few minutes for the DNS records to propagate, please verify the records are set and try again.")
			return
		default:
			verificationErr := r.client.VerifyCustomDomainInfo(environment, isSwitching)
			if verificationErr == nil {
				// Verification successful
				// Set the data from the state into the response
				resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
				return
			}

			// Log the retry attempt
			tflog.Warn(ctx, "Unable to verify the custom domain. It can take a few minutes for the DNS records to propagate. Retrying in 30 seconds...")

			// Wait for the retry interval before the next attempt
			time.Sleep(retryInterval)
		}
	}
}

func (r *customDomainVerificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "deleted a custom_domain resource")
}

func (r *customDomainVerificationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var state customDomainVerificationResourceModel

	environment := req.ID
	if environment != "Staging" && environment != "Prod" {
		resp.Diagnostics.AddError("Invalid import ID", "The import ID must be either `Staging` or `Prod`.")
		return
	}

	customDomainInfo, err := r.client.GetCustomDomainInfo(environment, false)
	if err != nil {
		resp.Diagnostics.AddError("Error fetching custom domain info", "Could not fetch custom domain info for the environment.")
		return
	}

	state.Environment = types.StringValue(environment)
	state.Domain = types.StringValue(customDomainInfo.Domain)
	// need to manually pull timeouts from hcl to pass validation
	var timeouts timeouts.Value
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Timeouts = timeouts

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
