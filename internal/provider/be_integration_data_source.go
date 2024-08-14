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
  _ datasource.DataSource = &beIntegrationDataSource{}
)

// NewBeIntegrationDataSource is a helper function to simplify the provider implementation.
func NewBeIntegrationDataSource() datasource.DataSource {
  return &beIntegrationDataSource{}
}

type beIntegrationDataSource struct {
	client *propelauth.PropelAuthClient
}

type beIntegrationDataSourceModel struct {
	Environment types.String `tfsdk:"environment"`
	AuthUrl types.String `tfsdk:"auth_url"`
	PublicKey types.String `tfsdk:"public_key"`
	Issuer types.String `tfsdk:"issuer"`
}

// Metadata returns the data source type name.
func (d *beIntegrationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
  resp.TypeName = req.ProviderTypeName + "_be_integration"
}

// Schema defines the schema for the data source.
func (d *beIntegrationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve the parameters for a backend integration with one of your PropelAuth environments.",
		Attributes: map[string]schema.Attribute{
			"environment": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("Test", "Staging", "Prod"),
				},
				Description: "The environment for which you are configuring the backend integration. Accepted values are `Test`, `Staging`, and `Prod`.",
			},
			"auth_url": schema.StringAttribute{
				Computed: true,
				Description: "The URL to the authentication endpoint for the environment. This is needed in PropelAuth backend libraries.",
			},
			"public_key": schema.StringAttribute{
				Computed: true,
				Description: "Your public key that can be used to verify access tokens. This is optional in our backend libraries, and " +
					"if unspecified, the libraries will fetch it for you.",
			},
			"issuer": schema.StringAttribute{
				Computed: true,
				Description: "A value that we verify in the access token. This is optional in our backend libraries, and " +
					"if unspecified, the libraries will fetch it for you.",
			},
		},
	}
}

func (d *beIntegrationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *beIntegrationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state beIntegrationDataSourceModel

	// Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch the data from the PropelAuth API
	beIntegrationInfo, err := d.client.GetBeIntegrationInfo(state.Environment.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to fetch data from PropelAuth API", err.Error())
		return
	}
	state.AuthUrl = types.StringValue(beIntegrationInfo.AuthUrl)
	state.PublicKey = types.StringValue(beIntegrationInfo.VerifierKey)
	state.Issuer = types.StringValue(beIntegrationInfo.Issuer)

	// Write the data to the Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
