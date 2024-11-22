package provider

import (
	"context"
	"fmt"
	"terraform-provider-propelauth/internal/propelauth"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &SocialLoginRedirectDataSource{}
)

// NewSocialLoginRedirectDataSource is a helper function to simplify the provider implementation.
func NewSocialLoginRedirectDataSource() datasource.DataSource {
	return &SocialLoginRedirectDataSource{}
}

type SocialLoginRedirectDataSource struct {
	client *propelauth.PropelAuthClient
}

type SocialLoginRedirectDataSourceModel struct {
	Environment    types.String `tfsdk:"environment"`
	SocialProvider types.String `tfsdk:"social_provider"`
	RedirectUrl    types.String `tfsdk:"redirect_url"`
}

// Metadata returns the data source type name.
func (d *SocialLoginRedirectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_social_login_redirect"
}

// Schema defines the schema for the data source.
func (d *SocialLoginRedirectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the redirect url needed for configuring an OIDC provider for a Social Login in PropelAuth.",
		Attributes: map[string]schema.Attribute{
			"environment": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("Test", "Staging", "Prod"),
				},
				Description: "The environment for which you are configuring the social login. Accepted values are `Test`, `Staging`, and `Prod`.",
			},
			"social_provider": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("Google", "Microsoft", "GitHub", "Slack", "LinkedIn", "Atlassian", "Apple", "Salesforce", "QuickBooks", "Xero", "Salesloft", "Outreach"),
				},
				Description: "The social login provider for which you are configuring and need the redirect URL. Accepted values are " +
					"`Google`, `Microsoft`, `GitHub`, `Slack`, `LinkedIn`, `Atlassian`, `Apple`, `Salesforce`, `QuickBooks`, `Xero`, `Salesloft`, and `Outreach`.",
			},
			"redirect_url": schema.StringAttribute{
				Computed:    true,
				Description: "The redirect URL to be white-listed in the OIDC configuration of the social login provider.",
			},
		},
	}
}

func (d *SocialLoginRedirectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*propelauth.PropelAuthClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data source Configure Type",
			fmt.Sprintf("Expected *propelauth.PropelAuthClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *SocialLoginRedirectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state SocialLoginRedirectDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch the data from the PropelAuth API
	socialLoginRedirectUrl, err := d.client.GetSocialLoginRedirectUrl(state.Environment.ValueString(), state.SocialProvider.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to fetch social login redirect url from PropelAuth API", err.Error())
		return
	}
	state.RedirectUrl = types.StringValue(*socialLoginRedirectUrl)

	// Write the data to the Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
