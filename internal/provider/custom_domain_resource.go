package provider

import (
	"context"
	"fmt"
	"strings"
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
var _ resource.ResourceWithConfigure = &customDomainResource{}
var _ resource.ResourceWithImportState = &customDomainResource{}

func NewCustomDomainResource() resource.Resource {
	return &customDomainResource{}
}

// customDomainResource defines the resource implementation.
type customDomainResource struct {
	client *propelauth.PropelAuthClient
}

// projectInfoResourceModel describes the resource data model.
type customDomainResourceModel struct {
	Environment                 types.String `tfsdk:"environment"`
	Domain                      types.String `tfsdk:"domain"`
	Subdomain                   types.String `tfsdk:"subdomain"`
	TxtRecordKey                types.String `tfsdk:"txt_record_key"`
	TxtRecordKeyWithoutDomain   types.String `tfsdk:"txt_record_key_without_domain"`
	TxtRecordValue              types.String `tfsdk:"txt_record_value"`
	CnameRecordKey              types.String `tfsdk:"cname_record_key"`
	CnameRecordKeyWithoutDomain types.String `tfsdk:"cname_record_key_without_domain"`
	CnameRecordValue            types.String `tfsdk:"cname_record_value"`
}

func (r *customDomainResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_domain"
}

func (r *customDomainResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "This resource just sets up the process of verifying the domain. " +
			"It will return the TXT and CNAME records that you need to add to your DNS settings. " +
			"You will need to add these records to your DNS settings manually or using Terraform. " +
			"Then, the `propelauth_custom_domain_verification` resource will verify the domain.",
		Attributes: map[string]schema.Attribute{
			"environment": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("Staging", "Prod"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The environment for which you are configuring the custom domain. Accepted values are `Staging` and `Prod`.",
			},
			"domain": schema.StringAttribute{
				Required: true,
				Description: "The domain name for the custom domain. Your resulting auth domain will be `auth.<domain>`. " +
					"You can also specify a subdomain like prod.example.com which will result in auth.prod.example.com",
			},
			"subdomain": schema.StringAttribute{
				Optional: true,
				Description: "The subdomain for the custom domain. This is optional, but recommended, as it will " +
					"allow PropelAuth to automatically redirect users to your application after they login.",
			},
			"txt_record_key": schema.StringAttribute{
				Computed:    true,
				Description: "The TXT record key for the custom domain.",
			},
			"txt_record_key_without_domain": schema.StringAttribute{
				Computed:    true,
				Description: "The TXT record key for the custom domain without the domain (e.g. just auth instead of auth.example.com) .",
			},
			"txt_record_value": schema.StringAttribute{
				Computed:    true,
				Description: "The TXT record value for the custom domain.",
			},
			"cname_record_key": schema.StringAttribute{
				Computed:    true,
				Description: "The CNAME record key for the custom domain.",
			},
			"cname_record_key_without_domain": schema.StringAttribute{
				Computed:    true,
				Description: "The CNAME record key for the custom domain without the domain (e.g. just auth instead of auth.example.com) .",
			},
			"cname_record_value": schema.StringAttribute{
				Computed:    true,
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
	customDomainInfo, err := r.client.UpdateCustomDomainInfo(environment, domain, subdomain, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting custom domain info",
			"Could not set custom domain info, unexpected error: "+err.Error(),
		)
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a propelauth_custom_domain resource")

	// Derive fields so they can be used either with or without the full domain
	cnameRecordKeyParts := strings.Split(*customDomainInfo.CnameRecordKey, ".")
	cnameRecordKeyWithoutDomain := strings.Join(cnameRecordKeyParts[:len(cnameRecordKeyParts)-2], ".")

	txtRecordKeyParts := strings.Split(*customDomainInfo.TxtRecordKey, ".")
	txtRecordKeyWithoutDomain := strings.Join(txtRecordKeyParts[:len(txtRecordKeyParts)-2], ".")

	// Save data into Terraform state
	plan.CnameRecordKey = types.StringPointerValue(customDomainInfo.CnameRecordKey)
	plan.CnameRecordValue = types.StringPointerValue(customDomainInfo.CnameRecordValue)
	plan.CnameRecordKeyWithoutDomain = types.StringValue(cnameRecordKeyWithoutDomain)
	plan.TxtRecordKey = types.StringPointerValue(customDomainInfo.TxtRecordKey)
	plan.TxtRecordKeyWithoutDomain = types.StringValue(txtRecordKeyWithoutDomain)
	plan.TxtRecordValue = types.StringPointerValue(customDomainInfo.TxtRecordValue)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *customDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read the data from the state
	var state customDomainResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the main env's custom domain info
	environment := state.Environment.ValueString()

	// isSwitching := state.IsPending.ValueBool()
	customDomainInfo, err := r.client.GetCustomDomainInfo(environment, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting custom domain info",
			"Could not get custom domain info, unexpected error: "+err.Error(),
		)
		return
	}

	isDomainOrSubdomainChanged := state.Domain.ValueString() != customDomainInfo.Domain || state.Subdomain.ValueStringPointer() != customDomainInfo.Subdomain

	isSwitching := isDomainOrSubdomainChanged && customDomainInfo.IsVerified
	if isSwitching {
		// If the domain is switching, fetch the pending state instead.
		customDomainInfo, err = r.client.GetCustomDomainInfo(environment, true)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting custom domain info",
				"Could not get custom domain info, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a custom_domain resource")

	// Set the data from the state into the response
	state.Domain = types.StringValue(customDomainInfo.Domain)
	state.Subdomain = types.StringPointerValue(customDomainInfo.Subdomain)

	// So as not to need an update after verification of a domain,
	// these fields are only updated if the domain is not verified, which is when they are
	// returned.
	if customDomainInfo.TxtRecordKey != nil {
		state.TxtRecordKey = types.StringPointerValue(customDomainInfo.TxtRecordKey)

		txtRecordKeyParts := strings.Split(*customDomainInfo.TxtRecordKey, ".")
		txtRecordKeyWithoutDomain := strings.Join(txtRecordKeyParts[:len(txtRecordKeyParts)-2], ".")
		state.TxtRecordKeyWithoutDomain = types.StringValue(txtRecordKeyWithoutDomain)
	}
	if customDomainInfo.TxtRecordValue != nil {
		state.TxtRecordValue = types.StringPointerValue(customDomainInfo.TxtRecordValue)
	}
	if customDomainInfo.CnameRecordKey != nil {
		state.CnameRecordKey = types.StringPointerValue(customDomainInfo.CnameRecordKey)

		cnameRecordKeyParts := strings.Split(*customDomainInfo.CnameRecordKey, ".")
		cnameRecordKeyWithoutDomain := strings.Join(cnameRecordKeyParts[:len(cnameRecordKeyParts)-2], ".")
		state.CnameRecordKeyWithoutDomain = types.StringValue(cnameRecordKeyWithoutDomain)
	}
	if customDomainInfo.CnameRecordValue != nil {
		state.CnameRecordValue = types.StringPointerValue(customDomainInfo.CnameRecordValue)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *customDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan data into the model
	var plan customDomainResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the current state data
	var state customDomainResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Re-fetch the main env's custom domain info to check
	// if its verification status has changed.
	customDomainInfo, err := r.client.GetCustomDomainInfo(state.Environment.ValueString(), false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting custom domain info",
			"Could not get custom domain info, unexpected error: "+err.Error(),
		)
		return
	}

	isVerified := customDomainInfo.IsVerified
	isPending := customDomainInfo.IsPending

	// Update the custom domain info
	environment := plan.Environment.ValueString()
	domain := plan.Domain.ValueString()
	subdomain := plan.Subdomain.ValueStringPointer()
	isSwitching := isPending || (!isPending && isVerified)
	customDomainInfo, err = r.client.UpdateCustomDomainInfo(environment, domain, subdomain, isSwitching)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting custom domain info",
			"Could not set custom domain info, unexpected error: "+err.Error(),
		)
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "updated a propelauth_custom_domain resource")

	// Derive fields so they can be used either with or without the full domain
	cnameRecordKeyParts := strings.Split(*customDomainInfo.CnameRecordKey, ".")
	cnameRecordKeyWithoutDomain := strings.Join(cnameRecordKeyParts[:len(cnameRecordKeyParts)-2], ".")

	txtRecordKeyParts := strings.Split(*customDomainInfo.TxtRecordKey, ".")
	txtRecordKeyWithoutDomain := strings.Join(txtRecordKeyParts[:len(txtRecordKeyParts)-2], ".")

	// Save data into Terraform state
	plan.CnameRecordKey = types.StringPointerValue(customDomainInfo.CnameRecordKey)
	plan.CnameRecordKeyWithoutDomain = types.StringValue(cnameRecordKeyWithoutDomain)
	plan.CnameRecordValue = types.StringPointerValue(customDomainInfo.CnameRecordValue)
	plan.TxtRecordKey = types.StringPointerValue(customDomainInfo.TxtRecordKey)
	plan.TxtRecordKeyWithoutDomain = types.StringValue(txtRecordKeyWithoutDomain)
	plan.TxtRecordValue = types.StringPointerValue(customDomainInfo.TxtRecordValue)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *customDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "deleted a propelauth_custom_domain resource")
}

func (r *customDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var state customDomainResourceModel

	// Get the main env's custom domain info
	environment := req.ID
	if environment != "Staging" && environment != "Prod" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID must be either `Staging` or `Prod`.",
		)
		return
	}

	customDomainInfo, err := r.client.GetCustomDomainInfo(environment, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing custom domain info",
			"Could not get custom domain info, unexpected error: "+err.Error(),
		)
		return
	}

	// Set the state data from the response
	state.Environment = types.StringValue(environment)
	state.Domain = types.StringValue(customDomainInfo.Domain)
	state.Subdomain = types.StringPointerValue(customDomainInfo.Subdomain)

	// So as not to need an update after verification of a domain,
	// these fields are only updated if the domain is not verified, which is when they are
	// returned.
	if customDomainInfo.TxtRecordKey != nil {
		state.TxtRecordKey = types.StringPointerValue(customDomainInfo.TxtRecordKey)

		txtRecordKeyParts := strings.Split(*customDomainInfo.TxtRecordKey, ".")
		txtRecordKeyWithoutDomain := strings.Join(txtRecordKeyParts[:len(txtRecordKeyParts)-2], ".")
		state.TxtRecordKeyWithoutDomain = types.StringValue(txtRecordKeyWithoutDomain)
	}
	if customDomainInfo.TxtRecordValue != nil {
		state.TxtRecordValue = types.StringPointerValue(customDomainInfo.TxtRecordValue)
	}
	if customDomainInfo.CnameRecordKey != nil {
		state.CnameRecordKey = types.StringPointerValue(customDomainInfo.CnameRecordKey)

		cnameRecordKeyParts := strings.Split(*customDomainInfo.CnameRecordKey, ".")
		cnameRecordKeyWithoutDomain := strings.Join(cnameRecordKeyParts[:len(cnameRecordKeyParts)-2], ".")
		state.CnameRecordKeyWithoutDomain = types.StringValue(cnameRecordKeyWithoutDomain)
	}
	if customDomainInfo.CnameRecordValue != nil {
		state.CnameRecordValue = types.StringPointerValue(customDomainInfo.CnameRecordValue)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
