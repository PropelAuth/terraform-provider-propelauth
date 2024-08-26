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
var _ resource.Resource = &customDomainVerificationResource{}
var _ resource.ResourceWithConfigure   = &customDomainVerificationResource{}

func NewCustomDomainVerificationResource() resource.Resource {
	return &customDomainVerificationResource{}
}

// customDomainVerificationResource defines the resource implementation.
type customDomainVerificationResource struct {
	client *propelauth.PropelAuthClient
}

// projectInfoResourceModel describes the resource data model.
type customDomainVerificationResourceModel struct {
	Environment types.String `tfsdk:"environment"`
	IsVerified types.Bool `tfsdk:"is_verified"`
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
					stringvalidator.OneOf("Staging", "Prod", "PendingStaging", "PendingProd"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The environment for which you are configuring the custom domain. Accepted values are `Staging`, `Prod`, `PendingStaging`, and `PendingProd`. Use the `Pending` environments for switching.",
			},
			"is_verified": schema.BoolAttribute{
				Computed: true,
				Description: "Whether the custom domain has been verified.",
			},
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

	// Verify the custom domain
	environment := plan.Environment.ValueString()
	err := r.client.VerifyCustomDomainInfo(environment)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error verifying custom domain",
			"Could not verify custom domain, unexpected error: " + err.Error(),
		)
		return
	}

	// Call the read method to get the updated state
	customDomainInfo, err := r.client.GetCustomDomainInfo(environment)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting custom domain info",
			"Could not get custom domain info, unexpected error: " + err.Error(),
		)
		return
	}

	// Set the data from the state into the response
	plan.IsVerified = types.BoolValue(customDomainInfo.IsVerified)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *customDomainVerificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read the data from the state
	var state customDomainVerificationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the custom domain info
	environment := state.Environment.ValueString()
	customDomainInfo, err := r.client.GetCustomDomainInfo(environment)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting custom domain info",
			"Could not get custom domain info, unexpected error: " + err.Error(),
		)
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a custom_domain resource")

	// Set the data from the state into the response
	state.IsVerified = types.BoolValue(customDomainInfo.IsVerified)
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

	// Verify the custom domain
	environment := plan.Environment.ValueString()
	err := r.client.VerifyCustomDomainInfo(environment)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error verifying custom domain",
			"Could not verify custom domain, unexpected error: " + err.Error(),
		)
		return
	}

	// Call the read method to get the updated state
	customDomainInfo, err := r.client.GetCustomDomainInfo(environment)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting custom domain info",
			"Could not get custom domain info, unexpected error: " + err.Error(),
		)
		return
	}

	// Set the data from the state into the response
	plan.IsVerified = types.BoolValue(customDomainInfo.IsVerified)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *customDomainVerificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "deleted a custom_domain resource")
}
