package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-propelauth/internal/propelauth"
)

// Ensure Provider satisfies various provider interfaces.
var _ provider.Provider = &propelauthProvider{}
var _ provider.ProviderWithFunctions = &propelauthProvider{}

// propelauthProvider defines the provider implementation.
type propelauthProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// propelauthProviderModel describes the provider data model.
type propelauthProviderModel struct {
	TenantId types.String `tfsdk:"tenant_id"`
	ProjectId types.String `tfsdk:"project_id"`
	ApiKey types.String `tfsdk:"api_key"`
}

func (p *propelauthProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "propelauth"
	resp.Version = p.version
}

func (p *propelauthProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
        Description: "Manage your PropelAuth integration for authentication, B2B authorization, and user management.",
		Attributes: map[string]schema.Attribute{
            "tenant_id": schema.StringAttribute{
                Optional: true,
				Description: "Your PropelAuth Tenant ID. This can be retrieved from Infrastructure Integration page of the PropelAuth Dashboard.",
            },
            "project_id": schema.StringAttribute{
				Optional: true,
				Description: "Your PropelAuth Project ID. This can be retrieved from Infrastructure Integration page of the PropelAuth Dashboard.",
            },
            "api_key": schema.StringAttribute{
				Optional: true,
                Sensitive: true,
				Description: "A PropelAuth Infrastructure Integration Key for your project. " +
					"You can generate one on Infrastructure Integration page of the PropelAuth Dashboard.",
            },
        },
	}
}

func (p *propelauthProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring PropelAuth client")

	// Retrieve provider data from configuration
    var config propelauthProviderModel
    diags := req.Config.Get(ctx, &config)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // If practitioner provided a configuration value for any of the
    // attributes, it must be a known value.

    if config.TenantId.IsUnknown() {
        resp.Diagnostics.AddAttributeError(
            path.Root("tenant_id"),
            "Unknown PropelAuth API TenantId",
            "The provider cannot create the PropelAuth API client as there is an unknown configuration value for the PropelAuth API tenant_id. "+
                "Either target apply the source of the value first, set the value statically in the configuration, or use the PROPELAUTH_TENANT_ID environment variable.",
        )
    }

    if config.ProjectId.IsUnknown() {
        resp.Diagnostics.AddAttributeError(
            path.Root("project_id"),
            "Unknown PropelAuth API ProjectId",
            "The provider cannot create the PropelAuth API client as there is an unknown configuration value for the PropelAuth API project_id. "+
                "Either target apply the source of the value first, set the value statically in the configuration, or use the PROPELAUTH_PROJECT_ID environment variable.",
        )
    }

    if config.ApiKey.IsUnknown() {
        resp.Diagnostics.AddAttributeError(
            path.Root("api_key"),
            "Unknown PropelAuth API ApiKey",
            "The provider cannot create the PropelAuth API client as there is an unknown configuration value for the PropelAuth API api_key. "+
                "Either target apply the source of the value first, set the value statically in the configuration, or use the PROPELAUTH_API_KEY environment variable.",
        )
    }

    if resp.Diagnostics.HasError() {
        return
    }

    // Default values to environment variables, but override
    // with Terraform configuration value if set.

    tenantId := os.Getenv("PROPELAUTH_TENANT_ID")
    projectId := os.Getenv("PROPELAUTH_PROJECT_ID")
    apiKey := os.Getenv("PROPELAUTH_API_KEY")

    if !config.TenantId.IsNull() {
        tenantId = config.TenantId.ValueString()
    }

    if !config.ProjectId.IsNull() {
        projectId = config.ProjectId.ValueString()
    }

    if !config.ApiKey.IsNull() {
        apiKey = config.ApiKey.ValueString()
    }

    // If any of the expected configurations are missing, return
    // errors with provider-specific guidance.

    if tenantId == "" {
        resp.Diagnostics.AddAttributeError(
            path.Root("tenant_id"),
            "Missing PropelAuth API TenantId",
            "The provider cannot create the PropelAuth API client as there is a missing or empty value for the PropelAuth API tenant_id. "+
                "Set the tenant_id value in the configuration or use the PROPELAUTH_TENANT_ID environment variable. "+
                "If either is already set, ensure the value is not empty.",
        )
    }

    if projectId == "" {
        resp.Diagnostics.AddAttributeError(
            path.Root("project_id"),
            "Missing PropelAuth API ProjectId",
            "The provider cannot create the PropelAuth API client as there is a missing or empty value for the PropelAuth API project_id. "+
                "Set the project_id value in the configuration or use the PROPELAUTH_PROJECT_ID environment variable. "+
                "If either is already set, ensure the value is not empty.",
        )
    }

    if apiKey == "" {
        resp.Diagnostics.AddAttributeError(
            path.Root("api_key"),
            "Missing PropelAuth API ApiKey",
            "The provider cannot create the PropelAuth API client as there is a missing or empty value for the PropelAuth API api_key. "+
                "Set the api_key value in the configuration or use the PROPELAUTH_API_KEY environment variable. "+
                "If either is already set, ensure the value is not empty.",
        )
    }

    if resp.Diagnostics.HasError() {
        return
    }

	ctx = tflog.SetField(ctx, "propelauth_tenant_id", tenantId)
    ctx = tflog.SetField(ctx, "propelauth_project_id", projectId)
    ctx = tflog.SetField(ctx, "propelauth_api_key", apiKey)
    ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "propelauth_api_key")

    tflog.Debug(ctx, "Creating PropelAuth API client")

    // Create a new PropelAuth client using the configuration values
    client, err := propelauth.NewClient(&tenantId, &projectId, &apiKey)
    if err != nil {
        resp.Diagnostics.AddError(
            "Unable to Create PropelAuth API Client",
            "An unexpected error occurred when creating the PropelAuth API client. "+
                "If the error is not clear, please contact the provider developers.\n\n"+
                "PropelAuth Client Error: " + err.Error(),
        )
        return
    }

    // Make the PropelAuth client available during DataSource and Resource
    // type Configure methods.
    resp.DataSourceData = client
    resp.ResourceData = client
}

func (p *propelauthProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectInfoResource,
        NewBasicAuthConfigurationResource,
        NewOrganizationConfigurationResource,
        NewThemeResource,
	}
}

func (p *propelauthProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *propelauthProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &propelauthProvider{
			version: version,
		}
	}
}
