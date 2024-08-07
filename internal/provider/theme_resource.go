package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/reiver/go-hexcolor"

	"terraform-provider-propelauth/internal/propelauth"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &themeResource{}
var _ resource.ResourceWithConfigure = &themeResource{}
var _ resource.ResourceWithValidateConfig = &themeResource{}

func NewThemeResource() resource.Resource {
	return &themeResource{}
}

// themeResource defines the resource implementation.
type themeResource struct {
	client *propelauth.PropelAuthClient
}

// themeResourceModel describes the resource data model.
type themeResourceModel struct {
	// LogoImage image `tfsdk:"logo_image"`
	// FaviconImage image `tfsdk:"favicon_image"`
	HeaderFont types.String `tfsdk:"header_font"`
	BodyFont types.String `tfsdk:"body_font"`
	LoginPageTheme loginPageTheme  `tfsdk:"login_page_theme"`
	ManagementPagesTheme managementPagesTheme  `tfsdk:"management_pages_theme"`
	DisplayProjectName types.Bool `tfsdk:"display_project_name"`
}

type loginPageTheme struct {
	Layout types.String `tfsdk:"layout"`
	BackgroundType types.String `tfsdk:"background_type"`
	SolidBackgroundParameters *solidBackgroundParameters `tfsdk:"solid_background_parameters"`
	GradientBackgroundParameters *gradientBackgroundParameters `tfsdk:"gradient_background_parameters"`
	ImageBackgroundParameters *imageBackgroundParameters `tfsdk:"image_background_parameters"`
	FrameBackgroundColor types.String `tfsdk:"frame_background_color"`
	FrameTextColor types.String `tfsdk:"frame_text_color"`
	PrimaryColor types.String `tfsdk:"primary_color"`
	PrimaryTextColor types.String `tfsdk:"primary_text_color"`
	ErrorColor types.String `tfsdk:"error_color"`
	ErrorButtonTextColor types.String `tfsdk:"error_button_text_color"`
	BorderColor types.String `tfsdk:"border_color"`
	SplitLoginPageParameters *splitLoginPageParameters `tfsdk:"split_login_page_parameters"`
}

type splitLoginPageParameters struct {
	Direction types.String `tfsdk:"direction"`
	ContentType types.String `tfsdk:"content_type"`
	Header types.String `tfsdk:"header"`
	Subheader types.String `tfsdk:"subheader"`
	SecondaryBackgroundTextColor types.String `tfsdk:"secondary_background_text_color"`
}

type solidBackgroundParameters struct {
	BackgroundColor types.String `tfsdk:"background_color"`
	BackgroundTextColor types.String `tfsdk:"background_text_color"`
}

type gradientBackgroundParameters struct {
	BackgroundGradientStartColor types.String `tfsdk:"background_gradient_start_color"`
	BackgroundGradientEndColor types.String `tfsdk:"background_gradient_end_color"`
	BackgroundGradientAngle types.Int32 `tfsdk:"background_gradient_angle"`
	BackgroundTextColor types.String `tfsdk:"background_text_color"`
}

type imageBackgroundParameters struct {
	BackgroundImage image `tfsdk:"background_image"`
	DefaultBackgroundColor types.String `tfsdk:"default_background_color"`
	BackgroundTextColor types.String `tfsdk:"background_text_color"`
}

type managementPagesTheme struct {
	MainBackgroundColor types.String `tfsdk:"main_background_color"`
	MainTextColor types.String `tfsdk:"main_text_color"`
	NavbarBackgroundColor types.String `tfsdk:"navbar_background_color"`
	NavbarTextColor types.String `tfsdk:"navbar_text_color"`
	ActionButtonColor types.String `tfsdk:"action_button_color"`
	ActionButtonTextColor types.String `tfsdk:"action_button_text_color"`
	BorderColor types.String `tfsdk:"border_color"`
	DisplayNavbar types.Bool `tfsdk:"display_navbar"`
}

type image struct {
	Content types.String `tfsdk:"content"`
	ContentSha1 types.String `tfsdk:"content_sha1"`
}

func (r *themeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_theme"
}

func (r *themeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	hexcode_regex := regexp.MustCompile(`^#(?:[0-9a-f]{3}){1,2}$`)

	resp.Schema = schema.Schema{
		Description: "Hosted Pages Look & Feel. This is for configuring the look and feel of your PropelAuth hosted pages.",
		Attributes: map[string]schema.Attribute{
			"header_font": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default: 		   stringdefault.StaticString("Inter"),
				Validators: []validator.String{
					stringvalidator.OneOf(
						"Roboto", "Inter", "OpenSans", "Montserrat", "Lato", "Poppins", "Raleway", "Jost",
						"Fraunces", "Caveat", "PlusJakartaSans",
					),
				},
				Description: "The font used for all headings in your hosted pages written in PascalCase. This includes both login and management pages. " +
					"Options include `Roboto`, `Inter`, `OpenSans`, `Montserrat`, `Lato`, `Poppins`, `Raleway`, `Jost`, " +
					"`Fraunces`, `Caveat`, `PlusJakartaSans`, etc." +
					"The default value is `Inter`.",
			},
			"body_font": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default: 		   stringdefault.StaticString("Inter"),
				Description: "The font used for all body text in your hosted pages. This includes both login and management pages. " +
					"The available options are the same as for `header_font`. The default value is `Inter`.",
			},
			"display_project_name": schema.BoolAttribute{
				Optional: 		  true,
				Computed: true,
				Default: 		   booldefault.StaticBool(true),
				Description: "If true, the project name is displayed in the header of the login page. " +
					"The default value is `true`.",
			},
			"login_page_theme": schema.SingleNestedAttribute{
				Description: "The theme for the login page.",
				Required: true,
				Attributes: map[string]schema.Attribute{
					"layout": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default: 		   stringdefault.StaticString("Frame"),
						Validators: []validator.String{
							stringvalidator.OneOf("Frame", "Frameless", "Split"),
						},
						Description: "The layout of the login page. Options include `Frame`, `Frameless`, and `Split`. " +
							"The default value is `Frame`.",
					},
					"background_type": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default: 		   stringdefault.StaticString("Solid"),
						Validators: []validator.String{
							stringvalidator.OneOf("Solid", "Gradient", "Image"),
						},
						Description: "The type of background for the login page. Options include `Solid`, `Gradient`, and `Image`. " +
							"The default value is `Solid`.",
					},
					"solid_background_parameters": schema.SingleNestedAttribute{
						Optional: true,
						Description: "The parameters required for a solid background in the login page.",
						Attributes: map[string]schema.Attribute{
							"background_color": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default: 		   stringdefault.StaticString("#f7f7f7"),
								Validators: []validator.String{
									stringvalidator.RegexMatches(hexcode_regex, "background_color must be a valid hex color code with lowercase characters."),
								},
								Description: "The color of a solid background in the login page. The default value is `#f7f7f7`.",
							},
							"background_text_color": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default: 		   stringdefault.StaticString("#363636"),
								Validators: []validator.String{
									stringvalidator.RegexMatches(hexcode_regex, "background_text_color must be a valid hex color code with lowercase characters."),
								},
								Description: "The color of the text on a solid background in the login page. The default value is `#363636`.",
							},
						},
					},
					"gradient_background_parameters": schema.SingleNestedAttribute{
						Optional: true,
						Description: "The parameters required for a gradient background in the login page.",
						Attributes: map[string]schema.Attribute{
							"background_gradient_start_color": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default: 		   stringdefault.StaticString("#f7f7f7"),
								Validators: []validator.String{
									stringvalidator.RegexMatches(hexcode_regex, "background_gradient_start_color must be a valid hex color code with lowercase characters."),
								},
								Description: "The start color of a gradient background in the login page. The default value is `#f7f7f7`.",
							},
							"background_gradient_end_color": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default: 		   stringdefault.StaticString("#f7f7f7"),
								Validators: []validator.String{
									stringvalidator.RegexMatches(hexcode_regex, "background_gradient_end_color must be a valid hex color code with lowercase characters."),
								},
								Description: "The end color of a gradient background in the login page. The default value is `#f7f7f7`.",
							},
							"background_gradient_angle": schema.Int32Attribute{
								Optional: true,
								Computed: true,
								Default: 		   int32default.StaticInt32(135),
								Description: "The angle of the gradient background in the login page. The default value is `135`.",
							},
							"background_text_color": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default: 		   stringdefault.StaticString("#363636"),
								Validators: []validator.String{
									stringvalidator.RegexMatches(hexcode_regex, "background_text_color must be a valid hex color code with lowercase characters."),
								},
								Description: "The color of the text on a gradient background in the login page. The default value is `#363636`.",
							},
						},
					},
					"image_background_parameters": schema.SingleNestedAttribute{
						Optional: true,
						Description: "The parameters required for an image background in the login page.",
						Attributes: map[string]schema.Attribute{
							"background_image": schema.SingleNestedAttribute{
								Optional: true,
								Computed: true,
								Description: "The image used as the background in the login page.",
								Attributes: map[string]schema.Attribute{
									"content": schema.StringAttribute{
										Optional: true,
										Computed: true,
										Description: "The content of the image.",
									},
									"content_sha1": schema.StringAttribute{
										Optional: true,
										Computed: true,
										Description: "The SHA1 hash of the image content.",
									},
								},
							},
							"default_background_color": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default: 		   stringdefault.StaticString("#f7f7f7"),
								Validators: []validator.String{
									stringvalidator.RegexMatches(hexcode_regex, "default_background_color must be a valid hex color code with lowercase characters."),
								},
								Description: "The default color behind the background image in the login page. The default value is `#f7f7f7`.",
							},
							"background_text_color": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default: 		   stringdefault.StaticString("#363636"),
								Validators: []validator.String{
									stringvalidator.RegexMatches(hexcode_regex, "background_text_color must be a valid hex color code with lowercase characters."),
								},
								Description: "The color of the text on an image background in the login page. The default value is `#363636`.",
							},
						},
					},
					"frame_background_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default: 		   stringdefault.StaticString("#ffffff"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "frame_background_color must be a valid hex color code with lowercase characters."),
						},
						Description: "The background color within the frame in the login page. If the the `layout` is `Frameless`, " +
							"this color is applied to the background of the input components on the page. The default value is `#ffffff`.",
					},
					"frame_text_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default: 		   stringdefault.StaticString("#0f0f0f"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "frame_text_color must be a valid hex color code with lowercase characters."),
						},
						Description: "The color of the text within the frame in the login page.  If the the `layout` is `Frameless`, " +
							"this color is applied to text within input components on the page. The default value is `#0f0f0f`.",
					},
					"primary_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default: 		   stringdefault.StaticString("#50c878"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "primary_color must be a valid hex color code with lowercase characters."),
						},
						Description: "The primary color of action buttons and links in the login page. " +
							"The default value is `#50c878`.",
					},
					"primary_text_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default: 		   stringdefault.StaticString("#f7f7f7"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "primary_text_color must be a valid hex color code with lowercase characters."),
						},
						Description: "The color of the text on action buttons in the login page. The default value is `#f7f7f7`.",
					},
					"error_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default: 		   stringdefault.StaticString("#cf222e"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "error_color must be a valid hex color code with lowercase characters."),
						},
						Description: "The color for error messages and cancel button in the login page. " +
							"The default value is `#cf222e`.",
					},
					"error_button_text_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default: 		   stringdefault.StaticString("#ffffff"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "error_button_text_color must be a valid hex color code with lowercase characters."),
						},
						Description: "The color of the text on error messages and cancel button in the login page. " +
							"The default value is `#ffffff`.",
					},
					"border_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default: 		   stringdefault.StaticString("#e4e4e4"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "border_color must be a valid hex color code with lowercase characters."),
						},
						Description: "The color of the borders in the login page. The default value is `#e4e4e4`.",
					},
					"split_login_page_parameters": schema.SingleNestedAttribute{
						Optional: true,
						Description: "The extra parameters required to configure a split login page.",
						Attributes: map[string]schema.Attribute{
							"direction": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default: 		   stringdefault.StaticString("Left"),
								Validators: []validator.String{
									stringvalidator.OneOf("Left", "Right"),
								},
								Description: "The side of the screen where all the login components are placed. " +
									"Options include `Left` and `Right`. The default value is `Left`.",
							},
							"content_type": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default: 		   stringdefault.StaticString("None"),
								Validators: []validator.String{
									stringvalidator.OneOf("None", "Text"),
								},
								Description: "The type of content displayed on the side of the screen opposite the login components. " +
									"Currently, options include `None` and `Text`. The default value is `None`.",
							},
							"header": schema.StringAttribute{
								Optional: true,
								Description: "The header text displayed on the side of the screen opposite the login components. " +
									"This is only displayed if `content_type` is `Text`.",
							},
							"subheader": schema.StringAttribute{
								Optional: true,
								Description: "The subheader text displayed on the side of the screen opposite the login components. " +
									"This is only displayed if `content_type` is `Text`.",
							},
							"secondary_background_text_color": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default: 		   stringdefault.StaticString("#363636"),
								Validators: []validator.String{
									stringvalidator.RegexMatches(hexcode_regex, "secondary_background_text_color must be a valid hex color code with lowercase characters."),
								},
								Description: "The color of the subheader on the side of the screen opposite the login components. " +
									"The header text in the same area uses the `background_text_color`. " +
									"The default value is `#363636`.",
							},
						},
					},
				},
			},
			"management_pages_theme": schema.SingleNestedAttribute{
				Description: "The theme for the account and organization management pages.",
				Required: true,
				Attributes: map[string]schema.Attribute{
					"main_background_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default: 		   stringdefault.StaticString("#f7f7f7"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "main_background_color must be a valid hex color code with lowercase characters."),
						},
						Description: "The background color of the main content area in the management pages. The default value is `#f7f7f7`.",
					},
					"main_text_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default: 		   stringdefault.StaticString("#363636"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "main_text_color must be a valid hex color code with lowercase characters."),
						},
						Description: "The color of the text in the main content area of the management pages. The default value is `#363636`.",
					},
					"navbar_background_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default: 		   stringdefault.StaticString("#ffffff"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "navbar_background_color must be a valid hex color code with lowercase characters."),
						},
						Description: "The background color of the navigation bar in the management pages. The default value is `#ffffff`.",
					},
					"navbar_text_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default: 		   stringdefault.StaticString("#0f0f0f"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "navbar_text_color must be a valid hex color code with lowercase characters."),
						},
						Description: "The color of the text in the navigation bar in the management pages. The default value is `#0f0f0f`.",
					},
					"action_button_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default: 		   stringdefault.StaticString("#50c878"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "action_button_color must be a valid hex color code with lowercase characters."),
						},
						Description: "The color of action buttons in the management pages. The default value is `#50c878`.",
					},
					"action_button_text_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default: 		   stringdefault.StaticString("#f7f7f7"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "action_button_text_color must be a valid hex color code with lowercase characters."),
						},
						Description: "The color of the text on action buttons in the management pages. The default value is `#f7f7f7`.",
					},
					"border_color": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default: 		   stringdefault.StaticString("#e4e4e4"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(hexcode_regex, "border_color must be a valid hex color code with lowercase characters."),
						},
						Description: "The color of the border between the navbar and the main content area in the management pages. " +
							"The default value is `#e4e4e4`.",
					},
					"display_navbar": schema.BoolAttribute{
						Optional: 		  true,
						Computed: true,
						Default: 		   booldefault.StaticBool(true),
						Description: "If true, the sidebar is displayed in the management pages. The default value is `true`.",
					},
				},
			},
		},
	}
}

func (r *themeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *themeResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
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
				"`Solid` `background_type` requires `solid_background_parameters` to be set.",
			)
			return
		}
		if plan.LoginPageTheme.GradientBackgroundParameters != nil || plan.LoginPageTheme.ImageBackgroundParameters != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("login_page_theme"),
				"Invalid background parameters",
				"`Solid` `background_type` should not have gradient or image background parameters set.",
			)
			return
		}
	}

	if plan.LoginPageTheme.BackgroundType.ValueString() == "Gradient" {
		if plan.LoginPageTheme.GradientBackgroundParameters == nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("login_page_theme"),
				"Missing gradient_background_parameters",
				"`Gradient` `background_type` requires `gradient_background_parameters` to be set.",
			)
			return
		}

		if plan.LoginPageTheme.SolidBackgroundParameters != nil || plan.LoginPageTheme.ImageBackgroundParameters != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("login_page_theme"),
				"Invalid background parameters",
				"`Gradient` `background_type` should not have solid or image background parameters set.",
			)
			return
		}
	}

	if plan.LoginPageTheme.BackgroundType.ValueString() == "Image" {
		if plan.LoginPageTheme.ImageBackgroundParameters == nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("login_page_theme"),
				"Missing `image_background_parameters`",
				"`Image` `background_type` requires `image_background_parameters` to be set.",
			)
			return
		}

		if plan.LoginPageTheme.SolidBackgroundParameters != nil || plan.LoginPageTheme.GradientBackgroundParameters != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("login_page_theme"),
				"Invalid background parameters",
				"`Image` `background_type` should not have solid or gradient background parameters set.",
			)
			return
		}
	}

	if plan.LoginPageTheme.Layout.ValueString() == "Split" && plan.LoginPageTheme.SplitLoginPageParameters == nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("login_page_theme"),
			"Missing `split_login_page_parameters`",
			"`Split` login page layout requires `split_login_page_parameters` to be set.",
		)
		return
	}
}

func (r *themeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan themeResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

	// Update the configuration in PropelAuth
	environmentConfigUpdate := propelauth.EnvironmentConfigUpdate{
		Theme: convertPlanToTheme(&plan),
	}

    _, err := r.client.UpdateEnvironmentConfig(&environmentConfigUpdate)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error setting propelauth theme",
            "Could not set propelauth theme, unexpected error: "+err.Error(),
        )
        return
    }

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a propelauth_theme resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *themeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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
            "Error Reading PropelAuth propelauth theme",
            "Could not read PropelAuth propelauth theme: " + err.Error(),
        )
        return
    }

	// overwrite the state with the retrieved data
	updateStateFromTheme(environmentConfig.Theme, &state)

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *themeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan themeResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

	// Update the configuration in PropelAuth
	environmentConfigUpdate := propelauth.EnvironmentConfigUpdate{
		Theme: convertPlanToTheme(&plan),
	}

    _, err := r.client.UpdateEnvironmentConfig(&environmentConfigUpdate)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error setting propelauth theme",
            "Could not set propelauth theme, unexpected error: "+err.Error(),
        )
        return
    }

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a propelauth_theme resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *themeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "deleted a propelauth_theme resource")
}

func convertPlanToTheme(plan *themeResourceModel) *propelauth.Theme {
	theme := propelauth.Theme{
		ThemeType: "CustomV2",
		HeaderFont: plan.HeaderFont.ValueString(),
		BodyFont: plan.BodyFont.ValueString(),
		DisplayProjectName: plan.DisplayProjectName.ValueBool(),
		LoginLayout: plan.LoginPageTheme.Layout.ValueString(),
		BackgroundType: plan.LoginPageTheme.BackgroundType.ValueString(),
		FrameBackgroundColor: convertHexColorToRgb(plan.LoginPageTheme.FrameBackgroundColor.ValueString()),
		FrameTextColor: convertHexColorToRgb(plan.LoginPageTheme.FrameTextColor.ValueString()),
		FrameSecondaryTextColor: convertHexColorToRgb(plan.LoginPageTheme.FrameTextColor.ValueString()),
		PrimaryColor: convertHexColorToRgb(plan.LoginPageTheme.PrimaryColor.ValueString()),
		PrimaryTextColor: convertHexColorToRgb(plan.LoginPageTheme.PrimaryTextColor.ValueString()),
		ErrorButtonColor: convertHexColorToRgb(plan.LoginPageTheme.ErrorColor.ValueString()),
		ErrorButtonTextColor: convertHexColorToRgb(plan.LoginPageTheme.ErrorButtonTextColor.ValueString()),
		BorderColor: convertHexColorToRgb(plan.LoginPageTheme.BorderColor.ValueString()),
		ManagementPagesTheme: propelauth.ManagementPagesTheme{
			MainBackgroundColor: convertHexColorToRgb(plan.ManagementPagesTheme.MainBackgroundColor.ValueString()),
			MainTextColor: convertHexColorToRgb(plan.ManagementPagesTheme.MainTextColor.ValueString()),
			NavbarBackgroundColor: convertHexColorToRgb(plan.ManagementPagesTheme.NavbarBackgroundColor.ValueString()),
			NavbarTextColor: convertHexColorToRgb(plan.ManagementPagesTheme.NavbarTextColor.ValueString()),
			ActionButtonColor: convertHexColorToRgb(plan.ManagementPagesTheme.ActionButtonColor.ValueString()),
			ActionButtonTextColor: convertHexColorToRgb(plan.ManagementPagesTheme.ActionButtonTextColor.ValueString()),
			BorderColor: convertHexColorToRgb(plan.ManagementPagesTheme.BorderColor.ValueString()),
			DisplayNavbar: plan.ManagementPagesTheme.DisplayNavbar.ValueBool(),
		},
	}

	if plan.LoginPageTheme.BackgroundType.ValueString() == "Solid" {
		theme.BackgroundColor = convertHexColorToRgb(plan.LoginPageTheme.SolidBackgroundParameters.BackgroundColor.ValueString())
		theme.BackgroundTextColor = convertHexColorToRgb(plan.LoginPageTheme.SolidBackgroundParameters.BackgroundTextColor.ValueString())
	}

	if plan.LoginPageTheme.BackgroundType.ValueString() == "Gradient" {
		theme.BackgroundColor = convertHexColorToRgb(plan.LoginPageTheme.GradientBackgroundParameters.BackgroundGradientStartColor.ValueString())
		theme.SecondaryBackgroundColor = convertHexColorToRgb(plan.LoginPageTheme.GradientBackgroundParameters.BackgroundGradientEndColor.ValueString())
		theme.BackgroundTextColor = convertHexColorToRgb(plan.LoginPageTheme.GradientBackgroundParameters.BackgroundTextColor.ValueString())
		theme.GradientAngle = plan.LoginPageTheme.GradientBackgroundParameters.BackgroundGradientAngle.ValueInt32()
	}

	if plan.LoginPageTheme.BackgroundType.ValueString() == "Image" {
		theme.BackgroundColor = convertHexColorToRgb(plan.LoginPageTheme.ImageBackgroundParameters.DefaultBackgroundColor.ValueString())
		theme.BackgroundTextColor = convertHexColorToRgb(plan.LoginPageTheme.ImageBackgroundParameters.BackgroundTextColor.ValueString())
	}

	if plan.LoginPageTheme.Layout.ValueString() == "Split" {
		var splitscreenParams propelauth.SplitscreenParams
		splitscreenParams.Direction = plan.LoginPageTheme.SplitLoginPageParameters.Direction.ValueString()
		splitscreenParams.ContentType = plan.LoginPageTheme.SplitLoginPageParameters.ContentType.ValueString()
		splitscreenParams.Header = plan.LoginPageTheme.SplitLoginPageParameters.Header.ValueString()
		splitscreenParams.Subheader = plan.LoginPageTheme.SplitLoginPageParameters.Subheader.ValueString()
		theme.Splitscreen = &splitscreenParams
		theme.SecondaryBackgroundTextColor = convertHexColorToRgb(plan.LoginPageTheme.SplitLoginPageParameters.SecondaryBackgroundTextColor.ValueString())
	}

	return &theme
}

func updateStateFromTheme(theme propelauth.Theme, state *themeResourceModel) {
	state.HeaderFont = types.StringValue(theme.HeaderFont)
	state.BodyFont = types.StringValue(theme.BodyFont)
	state.DisplayProjectName = types.BoolValue(theme.DisplayProjectName)
	state.LoginPageTheme = loginPageTheme{
		Layout: types.StringValue(theme.LoginLayout),
		BackgroundType: types.StringValue(theme.BackgroundType),
		FrameBackgroundColor: types.StringValue(convertRgbToHexColor(theme.FrameBackgroundColor)),
		FrameTextColor: types.StringValue(convertRgbToHexColor(theme.FrameTextColor)),
		PrimaryColor: types.StringValue(convertRgbToHexColor(theme.PrimaryColor)),
		PrimaryTextColor: types.StringValue(convertRgbToHexColor(theme.PrimaryTextColor)),
		ErrorColor: types.StringValue(convertRgbToHexColor(theme.ErrorButtonColor)),
		ErrorButtonTextColor: types.StringValue(convertRgbToHexColor(theme.ErrorButtonTextColor)),
		BorderColor: types.StringValue(convertRgbToHexColor(theme.BorderColor)),
	}
	state.ManagementPagesTheme = managementPagesTheme{
		MainBackgroundColor: types.StringValue(convertRgbToHexColor(theme.ManagementPagesTheme.MainBackgroundColor)),
		MainTextColor: types.StringValue(convertRgbToHexColor(theme.ManagementPagesTheme.MainTextColor)),
		NavbarBackgroundColor: types.StringValue(convertRgbToHexColor(theme.ManagementPagesTheme.NavbarBackgroundColor)),
		NavbarTextColor: types.StringValue(convertRgbToHexColor(theme.ManagementPagesTheme.NavbarTextColor)),
		ActionButtonColor: types.StringValue(convertRgbToHexColor(theme.ManagementPagesTheme.ActionButtonColor)),
		ActionButtonTextColor: types.StringValue(convertRgbToHexColor(theme.ManagementPagesTheme.ActionButtonTextColor)),
		BorderColor: types.StringValue(convertRgbToHexColor(theme.ManagementPagesTheme.BorderColor)),
		DisplayNavbar: types.BoolValue(theme.ManagementPagesTheme.DisplayNavbar),
	}

	if theme.BackgroundType == "Solid" {
		state.LoginPageTheme.SolidBackgroundParameters = &solidBackgroundParameters{
			BackgroundColor: types.StringValue(convertRgbToHexColor(theme.BackgroundColor)),
			BackgroundTextColor: types.StringValue(convertRgbToHexColor(theme.BackgroundTextColor)),
		}
	}

	if theme.BackgroundType == "Gradient" {
		state.LoginPageTheme.GradientBackgroundParameters = &gradientBackgroundParameters{
			BackgroundGradientStartColor: types.StringValue(convertRgbToHexColor(theme.BackgroundColor)),
			BackgroundGradientEndColor: types.StringValue(convertRgbToHexColor(theme.SecondaryBackgroundColor)),
			BackgroundGradientAngle: types.Int32Value(theme.GradientAngle),
			BackgroundTextColor: types.StringValue(convertRgbToHexColor(theme.BackgroundTextColor)),
		}
	}

	if theme.BackgroundType == "Image" {
		state.LoginPageTheme.ImageBackgroundParameters = &imageBackgroundParameters{
			// TODO: once we support image content, impl populating BackgrondImage
			DefaultBackgroundColor: types.StringValue(convertRgbToHexColor(theme.BackgroundColor)),
			BackgroundTextColor: types.StringValue(convertRgbToHexColor(theme.BackgroundTextColor)),
		}
	}

	if theme.LoginLayout == "Split" {
		state.LoginPageTheme.SplitLoginPageParameters = &splitLoginPageParameters{
			Direction: types.StringValue(theme.Splitscreen.Direction),
			ContentType: types.StringValue(theme.Splitscreen.ContentType),
			Header: types.StringValue(theme.Splitscreen.Header),
			Subheader: types.StringValue(theme.Splitscreen.Subheader),
			SecondaryBackgroundTextColor: types.StringValue(convertRgbToHexColor(theme.SecondaryBackgroundTextColor)),
		}
	}
}

func convertHexColorToRgb(hexColor string) propelauth.RgbColor {
	red, green, blue, err := hexcolor.Parse(hexColor)
	if err != nil {
		return propelauth.RgbColor{} // shouldn't happen since we validate all hex color inputs
	} else {
		return propelauth.RgbColor{
			Red: red,
			Green: green,
			Blue: blue,
		}
	}
}

func convertRgbToHexColor(rgb propelauth.RgbColor) string {
	return strings.ToLower(hexcolor.Format(rgb.Red, rgb.Green, rgb.Blue))
}