package propelauth

import (
	"github.com/google/uuid"
)

type ProjectInfoResponse struct {
	Name string `json:"name"`
	ProjectId uuid.UUID `json:"project_id"`
	TestRealmId uuid.UUID `json:"test_realm_id"`
	StageRealmId uuid.UUID `json:"stage_realm_id"`
	ProdRealmId uuid.UUID `json:"prod_realm_id"`
}

type ProjectInfoUpdateRequest struct {
	Name string `json:"name,omitempty"`
}

type EnvironmentConfigUpdate struct {
	AllowUsersToSignupWithPersonalEmail *bool `json:"allow_users_to_signup_with_personal_email,omitempty"`
	HasPasswordLogin *bool `json:"has_password_login,omitempty"`
	HasPasswordlessLogin *bool `json:"has_passwordless_login,omitempty"`
	WaitlistUsersEnabled *bool `json:"waitlist_users_enabled,omitempty"`
	UserAutologoutSeconds *int64 `json:"user_autologout_seconds,omitempty"`
	UserAutologoutType string `json:"user_autologout_type,omitempty"`
	UsersCanDeleteOwnAccount *bool `json:"users_can_delete_own_account,omitempty"`
	UsersCanChangeEmail *bool `json:"users_can_change_email,omitempty"`
	IncludeLoginMethod *bool `json:"include_login_method,omitempty"`
	HasOrgs *bool `json:"has_orgs,omitempty"`
	MaxNumOrgsUsersCanBeIn *int32 `json:"max_num_orgs_users_can_be_in,omitempty"`
	OrgsMetaname string `json:"orgs_metaname,omitempty"`
	UsersCanCreateOrgs *bool `json:"users_can_create_orgs,omitempty"`
	UsersCanDeleteTheirOwnOrgs *bool `json:"users_can_delete_their_own_orgs,omitempty"`
	UsersMustBeInAnOrganization *bool `json:"users_must_be_in_an_organization,omitempty"`
	OrgsCanSetupSaml *bool `json:"orgs_can_setup_saml,omitempty"`
	UseOrgNameForSaml *bool `json:"use_org_name_for_saml,omitempty"`
	DefaultToSamlLogin *bool `json:"default_to_saml_login,omitempty"`
	OrgsCanRequire2fa *bool `json:"orgs_can_require_2fa,omitempty"`
	Theme *Theme `json:"theme,omitempty"`
	LogoImageId string `json:"logo_image_id,omitempty"`
	FaviconImageId string `json:"favicon_image_id,omitempty"`
	BackgroundImageId string `json:"background_image_id,omitempty"`
}

type EnvironmentConfigResponse struct {
	AllowUsersToSignupWithPersonalEmail bool `json:"allow_users_to_signup_with_personal_email"`
	HasPasswordLogin bool `json:"has_password_login"`
	HasPasswordlessLogin bool `json:"has_passwordless_login"`
	WaitlistUsersEnabled bool `json:"waitlist_users_enabled"`
	UserAutologoutSeconds int64 `json:"user_autologout_seconds"`
	UserAutologoutType string `json:"user_autologout_type"`
	UsersCanDeleteOwnAccount bool `json:"users_can_delete_own_account"`
	UsersCanChangeEmail bool `json:"users_can_change_email"`
	IncludeLoginMethod bool `json:"include_login_method"`
	HasOrgs bool `json:"has_orgs"`
	MaxNumOrgsUsersCanBeIn int32 `json:"max_num_orgs_users_can_be_in"`
	OrgsMetaname string `json:"orgs_metaname"`
	UsersCanCreateOrgs bool `json:"users_can_create_orgs"`
	UsersCanDeleteTheirOwnOrgs bool `json:"users_can_delete_their_own_orgs"`
	UsersMustBeInAnOrganization bool `json:"users_must_be_in_an_organization"`
	OrgsCanSetupSaml bool `json:"orgs_can_setup_saml"`
	UseOrgNameForSaml bool `json:"use_org_name_for_saml"`
	DefaultToSamlLogin bool `json:"default_to_saml_login"`
	OrgsCanRequire2fa bool `json:"orgs_can_require_2fa"`
	Theme Theme `json:"theme"`
	LogoUrl string `json:"logo_url"`
	FaviconUrl string `json:"favicon_url"`
	BackgroundUrl string `json:"background_url"`
}

type Theme struct {
	ThemeType string `json:"theme_type"` // always "CustomV2" for now
	BodyFont string `json:"font_family"`
	HeaderFont string `json:"secondary_font_family"`
	DisplayProjectName bool `json:"display_project_name"`
	LoginLayout string `json:"login_ui_theme"`
	BackgroundType string `json:"background_type"`
	BackgroundColor RgbColor `json:"background_color"`
	SecondaryBackgroundColor RgbColor `json:"secondary_background_color"` // secondary background color in gradient
	GradientAngle int32 `json:"gradient_angle"` // angle of gradient in degrees
	BackgroundTextColor RgbColor `json:"background_text_color"` // header text color in split layout
	SecondaryBackgroundTextColor RgbColor `json:"secondary_background_text_color"` // subheader text color in split layout
	BorderColor RgbColor `json:"border_color"`
	FrameBackgroundColor RgbColor `json:"foreground_color"` // input/frame background color
	FrameTextColor RgbColor `json:"foreground_text_color"` // frame/input text color
	FrameSecondaryTextColor RgbColor `json:"foreground_secondary_text_color"` // currently always same as ForegroundTextColor
	PrimaryColor RgbColor `json:"success_button_color"`
	PrimaryTextColor RgbColor `json:"success_button_text_color"`
	ErrorButtonColor RgbColor `json:"error_button_color"`
	ErrorButtonTextColor RgbColor `json:"error_button_text_color"`
	Splitscreen *SplitscreenParams `json:"splitscreen,omitempty"`
	ManagementPagesTheme ManagementPagesTheme `json:"management_pages_theme"`
}

type RgbColor struct {
	Red uint8 `json:"red"`
	Green uint8 `json:"green"`
	Blue uint8 `json:"blue"`
}

type SplitscreenParams struct {
	Direction string `json:"direction"`
	ContentType string `json:"content_type"`
	Header string `json:"header"`
	Subheader string `json:"subheader"`
}

type ManagementPagesTheme struct {
	DisplayNavbar bool `json:"display_sidebar"`
	MainBackgroundColor RgbColor `json:"main_background_color"`
	MainTextColor RgbColor `json:"main_text_color"`
	NavbarBackgroundColor RgbColor `json:"sidebar_background_color"`
	NavbarTextColor RgbColor `json:"sidebar_text_color"`
	BorderColor RgbColor `json:"border_color"`
	ActionButtonColor RgbColor `json:"action_button_color"`
	ActionButtonTextColor RgbColor `json:"action_button_text_color"`
}

type ImageUploadResponse struct {
	ImageId string `json:"image_id"`
}