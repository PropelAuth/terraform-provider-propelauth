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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &darkmodeThemeResource{}
var _ resource.ResourceWithConfigure = &darkmodeThemeResource{}
var _ resource.ResourceWithValidateConfig = &darkmodeThemeResource{}
var _ resource.ResourceWithImportState = &darkmodeThemeResource{}

func NewDarkmodeThemeResource() resource.Resource {
	return &darkmodeThemeResource{}
}

// darkmodeThemeResource defines the resource implementation.
type darkmodeThemeResource struct {
	client *propelauth.PropelAuthClient
}

func (r *darkmodeThemeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_darkmode_theme"
}

func (r *darkmodeThemeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	hexcode_regex := regexp.MustCompile(`^#(?:[0-9a-f]{3}){1,2}$`)

	resp.Schema = schema.Schema{
		Description: "Darkmode Pages Look & Feel. This is for creating a darkmode theme for your PropelAuth hosted pages." +
			"The parameters and behavior are identical to the `propelauth_theme` resource, except this enables an optional darkmode " +
			"version for your users to toggle to. Altering these settings does not affect the primary theme.",
		Attributes: map[string]schema.Attribute{
			"header_font": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("Inter"),
				Validators: []validator.String{
					stringvalidator.OneOf(
						"Roboto", "Inter", "OpenSans", "Montserrat", "Lato", "Poppins", "Raleway", "Jost",
						"Fraunces", "Caveat", "PlusJakartaSans",
					),
				},
				Description: "The font used for all headings in your hosted pages written in PascalCase. This includes both login and management pages. " +
					"Options include `Roboto`, `Inter`, `OpenSans`, `Montserrat`, `Lato`, `Poppins`, `Raleway`, `Jost`, " +
					"`Fraunces`, `Caveat`, `PlusJakartaSans`, etc" +
					"The default value is `Inter`",
			},
			"body_font": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("Inter"),
				Description: "The font used for all body text in your hosted pages. This includes both login and management pages. " +
					"The available options are the same as for `header_font`. The default value is `Inter`",
			},
			"display_project_name": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
				Description: "If true, the project name is displayed in the header of the login page. " +
					"The default value is `true`",
			},
			"login_page_theme": schema.SingleNestedAttribute{
				Description: "The theme for the login page",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"layout": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("Frame"),
						Validators: []validator.String{
							stringvalidator.OneOf("Frame", "Frameless", "SplitScreen"),
						},
						Description: "The layout of the login page. Options include `Frame`, `Frameless`, and `SplitScreen`. " +
							"The default value is `Frame`",
					},
					"background_type": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("Solid"),
						Validators: []validator.String{
							stringvalidator.OneOf("Solid", "Gradient", "Image"),
						},
						Description: "The type of background for the login page. Options include `Solid`, `Gradient`, and `Image`. " +
							"The default value is `Solid`",
					},
					"solid_background_parameters": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "The parameters required for a solid background in the login page",
						Attributes: map[string]schema.Attribute{
							"background_color": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString("#f7f7f7"),
								Validators: []validator.String{
									stringvalidator.RegexMatches(hexcode_regex, "background_color must be a valid hex color code with lowercase characters"),
								},
								Description: "The color of a solid background in the login page. The default value is `#f7f7f7`",
							},
							"background_text_color": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString("#363636"),
								Validators: []validator.String{
									stringvalidator.RegexMatches(hexcode_regex, "background_text_color must be a valid hex color code with lowercase characters"),
								},
								Description: "The color of the text on a solid background in the login page. The default value is `#363636`",
							},
						},
					},
					"gradient_background_parameters": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "The parameters required for a gradient background in the login page",
						Attributes: map[string]schema.Attribute{
							"background_gradient_start_color": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString("#f7f7f7"),
								Validators: []validator.String{
									stringvalidator.RegexMatches(hexcode_regex, "background_gradient_start_color must be a valid hex color code with lowercase characters"),
								},
								Description: "The start color of a gradient background in the login page. The default value is `#f7f7f7`",
							},
							"background_gradient_end_color": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString("#f7f7f7"),
								Validators: []validator.String{
									stringvalidator.RegexMatches(hexcode_regex, "background_gradient_end_color must be a valid hex color code with lowercase characters"),
								},
								Description: "The end color of a gradient background in the login page. The default value is `#f7f7f7`",
							},
							"background_gradient_angle": schema.Int32Attribute{
								Optional:    true,
								Computed:    true,
								Default:     int32default.StaticInt32(135),
								Description: "The angle of the gradient background in the login page. The default value is `135`",
							},
							"background_text_color": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString("#363636"),
								Validators: []validator.String{
									stringvalidator.RegexMatches(hexcode_regex, "background_text_color must be a valid hex color code with lowercase characters"),
								},
								Description: "The color of the text on a gradient background in the login page. The default value is `#363636`",
							},
						},
					},
					"image_background_parameters": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "The parameters required for an image background in the login page",
						Attributes: map[string]schema.Attribute{
							"default_background_color": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString("#f7f7f7"),
								Validators: []validator.String{
									stringvalidator.RegexMatches(hexcode_regex, "default_background_color must be a valid hex color code with lowercase characters"),
								},
								Description: "The default color behind the background image in the login page. The default value is `#f7f7f7`",
							},
							"background_text_color": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString("#363636"),
								Validators: []validator.String{
									stringvalidator.RegexMatches(hexcode_regex, "background_text_color must be a valid hex color code with lowercase characters"),
								},
								Description: "The color of the text on an image background in the login page. The default value is `#363636`",
							},
						},
					},
					"frame_background_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("#ffffff"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "frame_background_color must be a valid hex color code with lowercase characters"),
						},
						Description: "The background color within the frame in the login page. If the the `layout` is `Frameless`, " +
							"this color is applied to the background of the input components on the page. The default value is `#ffffff`",
					},
					"frame_text_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("#0f0f0f"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "frame_text_color must be a valid hex color code with lowercase characters"),
						},
						Description: "The color of the text within the frame in the login page.  If the the `layout` is `Frameless`, " +
							"this color is applied to text within input components on the page. The default value is `#0f0f0f`",
					},
					"primary_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("#50c878"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "primary_color must be a valid hex color code with lowercase characters"),
						},
						Description: "The primary color of action buttons and links in the login page. " +
							"The default value is `#50c878`",
					},
					"primary_text_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("#f7f7f7"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "primary_text_color must be a valid hex color code with lowercase characters"),
						},
						Description: "The color of the text on action buttons in the login page. The default value is `#f7f7f7`",
					},
					"error_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("#cf222e"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "error_color must be a valid hex color code with lowercase characters"),
						},
						Description: "The color for error messages and cancel button in the login page. " +
							"The default value is `#cf222e`",
					},
					"error_button_text_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("#ffffff"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "error_button_text_color must be a valid hex color code with lowercase characters"),
						},
						Description: "The color of the text on error messages and cancel button in the login page. " +
							"The default value is `#ffffff`",
					},
					"border_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("#e4e4e4"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "border_color must be a valid hex color code with lowercase characters"),
						},
						Description: "The color of the borders in the login page. The default value is `#e4e4e4`",
					},
					"split_login_page_parameters": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "The extra parameters required to configure a split login page",
						Attributes: map[string]schema.Attribute{
							"direction": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString("Left"),
								Validators: []validator.String{
									stringvalidator.OneOf("Left", "Right"),
								},
								Description: "The side of the screen where all the login components are placed. " +
									"Options include `Left` and `Right`. The default value is `Left`",
							},
							"content_type": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString("None"),
								Validators: []validator.String{
									stringvalidator.OneOf("None", "Text"),
								},
								Description: "The type of content displayed on the side of the screen opposite the login components. " +
									"Currently, options include `None` and `Text`. The default value is `None`",
							},
							"header": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString(""),
								Description: "The header text displayed on the side of the screen opposite the login components. " +
									"This is only displayed if `content_type` is `Text`",
							},
							"subheader": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString(""),
								Description: "The subheader text displayed on the side of the screen opposite the login components. " +
									"This is only displayed if `content_type` is `Text`",
							},
							"secondary_background_text_color": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString("#363636"),
								Validators: []validator.String{
									stringvalidator.RegexMatches(hexcode_regex, "secondary_background_text_color must be a valid hex color code with lowercase characters"),
								},
								Description: "The color of the subheader on the side of the screen opposite the login components. " +
									"The header text in the same area uses the `background_text_color`. " +
									"The default value is `#363636`",
							},
						},
					},
				},
			},
			"management_pages_theme": schema.SingleNestedAttribute{
				Description: "The theme for the account and organization management pages",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"main_background_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("#f7f7f7"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "main_background_color must be a valid hex color code with lowercase characters"),
						},
						Description: "The background color of the main content area in the management pages. The default value is `#f7f7f7`",
					},
					"main_text_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("#363636"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "main_text_color must be a valid hex color code with lowercase characters"),
						},
						Description: "The color of the text in the main content area of the management pages. The default value is `#363636`",
					},
					"navbar_background_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("#ffffff"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "navbar_background_color must be a valid hex color code with lowercase characters"),
						},
						Description: "The background color of the navigation bar in the management pages. The default value is `#ffffff`",
					},
					"navbar_text_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("#0f0f0f"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "navbar_text_color must be a valid hex color code with lowercase characters"),
						},
						Description: "The color of the text in the navigation bar in the management pages. The default value is `#0f0f0f`",
					},
					"action_button_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("#50c878"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "action_button_color must be a valid hex color code with lowercase characters"),
						},
						Description: "The color of action buttons in the management pages. The default value is `#50c878`",
					},
					"action_button_text_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("#f7f7f7"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "action_button_text_color must be a valid hex color code with lowercase characters"),
						},
						Description: "The color of the text on action buttons in the management pages. The default value is `#f7f7f7`",
					},
					"border_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("#e4e4e4"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "border_color must be a valid hex color code with lowercase characters"),
						},
						Description: "The color of the border between the navbar and the main content area in the management pages. " +
							"The default value is `#e4e4e4`",
					},
					"display_navbar": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
						Description: "If true, the sidebar is displayed in the management pages. The default value is `true`",
					},
				},
			},
		},
	}
}

func (r *darkmodeThemeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*propelauth.PropelAuthClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *propelauth.PropelAuthClient, got: %T. Please report this issue to the provider developers", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *darkmodeThemeResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var plan themeResourceModel

	// Read Terraform plan data into the model
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate the plan data

	if plan.LoginPageTheme.BackgroundType.ValueString() == "Solid" {
		if plan.LoginPageTheme.SolidBackgroundParameters == nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("login_page_theme"),
				"Missing solid_background_parameters",
				"`Solid` `background_type` requires `solid_background_parameters` to be set",
			)
			return
		}
		if plan.LoginPageTheme.GradientBackgroundParameters != nil || plan.LoginPageTheme.ImageBackgroundParameters != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("login_page_theme"),
				"Invalid background parameters",
				"`Solid` `background_type` should not have gradient or image background parameters set",
			)
			return
		}
	}

	if plan.LoginPageTheme.BackgroundType.ValueString() == "Gradient" {
		if plan.LoginPageTheme.GradientBackgroundParameters == nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("login_page_theme"),
				"Missing gradient_background_parameters",
				"`Gradient` `background_type` requires `gradient_background_parameters` to be set",
			)
			return
		}

		if plan.LoginPageTheme.SolidBackgroundParameters != nil || plan.LoginPageTheme.ImageBackgroundParameters != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("login_page_theme"),
				"Invalid background parameters",
				"`Gradient` `background_type` should not have solid or image background parameters set",
			)
			return
		}
	}

	if plan.LoginPageTheme.BackgroundType.ValueString() == "Image" {
		if plan.LoginPageTheme.ImageBackgroundParameters == nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("login_page_theme"),
				"Missing `image_background_parameters`",
				"`Image` `background_type` requires `image_background_parameters` to be set",
			)
			return
		}

		if plan.LoginPageTheme.SolidBackgroundParameters != nil || plan.LoginPageTheme.GradientBackgroundParameters != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("login_page_theme"),
				"Invalid background parameters",
				"`Image` `background_type` should not have solid or gradient background parameters set",
			)
			return
		}
	}

	if plan.LoginPageTheme.Layout.ValueString() == "SplitScreen" && plan.LoginPageTheme.SplitLoginPageParameters == nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("login_page_theme"),
			"Missing `split_login_page_parameters`",
			"`SplitScreen` login page layout requires `split_login_page_parameters` to be set",
		)
		return
	}

	if plan.LoginPageTheme.SolidBackgroundParameters == nil && plan.LoginPageTheme.GradientBackgroundParameters == nil && plan.LoginPageTheme.ImageBackgroundParameters == nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("login_page_theme"),
			"Missing background parameters",
			"Either `solid_background_parameters`, `gradient_background_parameters`, or `image_background_parameters` must be set if even to an empty object",
		)
		return
	}
}

func (r *darkmodeThemeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan themeResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the configuration in PropelAuth
	enableDarkmodeTheme := true
	environmentConfigUpdate := propelauth.EnvironmentConfigUpdate{
		DarkmodeTheme:       convertPlanToTheme(&plan),
		EnableDarkmodeTheme: &enableDarkmodeTheme,
	}
	environmentConfig, err := r.client.UpdateEnvironmentConfig(&environmentConfigUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting propelauth darkmode theme",
			"Could not set propelauth darkmode theme, unexpected error: "+err.Error(),
		)
		return
	}

	// overwrite the computed state with the retrieved data
	updateStateFromTheme(environmentConfig.DarkmodeTheme, &plan)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a propelauth_darkmode_theme resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *darkmodeThemeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state and read it into the model
	var state themeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// retrieve the environment config from PropelAuth
	environmentConfig, err := r.client.GetEnvironmentConfig()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading PropelAuth propelauth darkmode theme",
			"Could not read PropelAuth propelauth darkmode theme: "+err.Error(),
		)
		return
	}

	// overwrite the state with the retrieved data
	updateStateFromTheme(environmentConfig.DarkmodeTheme, &state)

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *darkmodeThemeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan themeResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the configuration in PropelAuth
	enableDarkmodeTheme := true
	environmentConfigUpdate := propelauth.EnvironmentConfigUpdate{
		DarkmodeTheme:       convertPlanToTheme(&plan),
		EnableDarkmodeTheme: &enableDarkmodeTheme,
	}
	environmentConfig, err := r.client.UpdateEnvironmentConfig(&environmentConfigUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting propelauth darkmode theme",
			"Could not set propelauth darkmode theme, unexpected error: "+err.Error(),
		)
		return
	}

	// overwrite the computed state with the retrieved data
	updateStateFromTheme(environmentConfig.DarkmodeTheme, &plan)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a propelauth_darkmode_theme resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *darkmodeThemeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	enableDarkmodeTheme := false
	environmentConfigUpdate := propelauth.EnvironmentConfigUpdate{
		EnableDarkmodeTheme: &enableDarkmodeTheme,
	}
	environmentConfig, err := r.client.UpdateEnvironmentConfig(&environmentConfigUpdate)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting propelauth darkmode theme",
			"Could not set propelauth darkmode theme, unexpected error: "+err.Error(),
		)
		return
	}
	if environmentConfig.EnableDarkmodeTheme {
		resp.Diagnostics.AddError(
			"Error deleting propelauth darkmode theme",
			"Failed to disable propelauth darkmode theme",
		)
		return
	}

	tflog.Trace(ctx, "deleted a propelauth_darkmode_theme resource")
}

func (r *darkmodeThemeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// retrieve the environment config from PropelAuth
	environmentConfig, err := r.client.GetEnvironmentConfig()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading PropelAuth propelauth darkmode theme",
			"Could not read PropelAuth propelauth darkmode theme: "+err.Error(),
		)
		return
	}

	var state themeResourceModel
	updateStateFromTheme(environmentConfig.DarkmodeTheme, &state)

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
