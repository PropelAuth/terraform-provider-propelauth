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
var _ resource.Resource = &customDomainResource{}
var _ resource.ResourceWithConfigure   = &customDomainResource{}

func NewCustomDomainResource() resource.Resource {
	return &customDomainResource{}
}

// customDomainResource defines the resource implementation.
type customDomainResource struct {
	client *propelauth.PropelAuthClient
}

// projectInfoResourceModel describes the resource data model.
type customDomainResourceModel struct {
	Environment types.String `tfsdk:"environment"`
	Domain types.String `tfsdk:"domain"`
	Subdomain types.String `tfsdk:"subdomain"`
	TxtRecordKey types.String `tfsdk:"txt_record_key"`
	TxtRecordValue types.String `tfsdk:"txt_record_value"`
	CnameRecordKey types.String `tfsdk:"cname_record_key"`
	CnameRecordValue types.String `tfsdk:"cname_record_value"`
}

func (r *customDomainResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_domain"
}

func (r *customDomainResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Custom Domain resource. This is for configuring a custom domain for Production or Staging.",
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
			"domain": schema.StringAttribute{
				Required: true,
				Description: "The domain name for the custom domain.",
			},
			"subdomain": schema.StringAttribute{
				Optional: true,
				Description: "The subdomain for the custom domain. This is optional.",
			},
			"txt_record_key": schema.StringAttribute{
				Computed: true,
				Description: "The TXT record key for the custom domain.",
			},
			"txt_record_value": schema.StringAttribute{
				Computed: true,
				Description: "The TXT record value for the custom domain.",
			},
			"cname_record_key": schema.StringAttribute{
				Computed: true,
				Description: "The CNAME record key for the custom domain.",
			},
			"cname_record_value": schema.StringAttribute{
				Computed: true,
				Description: "The CNAME record value for the custom domain.",
			},
		},
	}
}

func (r *customDomainResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *customDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan customDomainResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

	// Update the custom domain info
	environment := plan.Environment.ValueString()
	domain := plan.Domain.ValueString()
	subdomain := plan.Subdomain.ValueStringPointer()
	customDomainInfo, err := r.client.UpdateCustomDomainInfo(environment, domain, subdomain)
    if err != nil {
        resp.Diagnostics.AddError(
			"Error setting custom domain info",
			"Could not set custom domain info, unexpected error: " + err.Error(),
		)
        return
    }


	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a propelauth_project_info resource")

	// Save data into Terraform state
	plan.CnameRecordKey = types.StringValue(customDomainInfo.CnameRecordKey)
	plan.CnameRecordValue = types.StringValue(customDomainInfo.CnameRecordValue)
	plan.TxtRecordKey = types.StringValue(customDomainInfo.TxtRecordKey)
	plan.TxtRecordValue = types.StringValue(customDomainInfo.TxtRecordValue)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *customDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read the data from the state
	var state customDomainResourceModel
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
	state.Domain = types.StringValue(customDomainInfo.Domain)
	state.Subdomain = types.StringPointerValue(customDomainInfo.Subdomain)
	state.TxtRecordKey = types.StringValue(customDomainInfo.TxtRecordKey)
	state.TxtRecordValue = types.StringValue(customDomainInfo.TxtRecordValue)
	state.CnameRecordKey = types.StringValue(customDomainInfo.CnameRecordKey)
	state.CnameRecordValue = types.StringValue(customDomainInfo.CnameRecordValue)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *customDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan customDomainResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the custom domain info
	environment := plan.Environment.ValueString()
	domain := plan.Domain.ValueString()
	subdomain := plan.Subdomain.ValueStringPointer()
	customDomainInfo, err := r.client.UpdateCustomDomainInfo(environment, domain, subdomain)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting custom domain info",
			"Could not set custom domain info, unexpected error: " + err.Error(),
		)
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "updated a custom_domain resource")

	// Save data into Terraform state
	plan.CnameRecordKey = types.StringValue(customDomainInfo.CnameRecordKey)
	plan.CnameRecordValue = types.StringValue(customDomainInfo.CnameRecordValue)
	plan.TxtRecordKey = types.StringValue(customDomainInfo.TxtRecordKey)
	plan.TxtRecordValue = types.StringValue(customDomainInfo.TxtRecordValue)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *customDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "deleted a custom_domain resource")
}
