package propelauth

import (
	"github.com/google/uuid"
)

type ProjectInfoResponse struct {
	Name         string    `json:"name"`
	ProjectId    uuid.UUID `json:"project_id"`
	TestRealmId  uuid.UUID `json:"test_realm_id"`
	StageRealmId uuid.UUID `json:"stage_realm_id"`
	ProdRealmId  uuid.UUID `json:"prod_realm_id"`
}

type ProjectInfoUpdateRequest struct {
	Name string `json:"name,omitempty"`
}

type EnvironmentConfigUpdate struct {
	AllowUsersToSignupWithPersonalEmail *bool            `json:"allow_users_to_signup_with_personal_email,omitempty"`
	SignupDomainAllowlistEnabled        *bool            `json:"signup_domain_allowlist_enabled,omitempty"`
	SignupDomainAllowlist               []string         `json:"signup_domain_allowlist,omitempty"`
	SignupDomainBlocklistEnabled        *bool            `json:"signup_domain_blocklist_enabled,omitempty"`
	SignupDomainBlocklist               []string         `json:"signup_domain_blocklist,omitempty"`
	HasPasswordLogin                    *bool            `json:"has_password_login,omitempty"`
	HasPasswordlessLogin                *bool            `json:"has_passwordless_login,omitempty"`
	WaitlistUsersEnabled                *bool            `json:"waitlist_users_enabled,omitempty"`
	UserAutologoutSeconds               *int64           `json:"user_autologout_seconds,omitempty"`
	UserAutologoutType                  string           `json:"user_autologout_type,omitempty"`
	UsersCanDeleteOwnAccount            *bool            `json:"users_can_delete_own_account,omitempty"`
	UsersCanChangeEmail                 *bool            `json:"users_can_change_email,omitempty"`
	IncludeLoginMethod                  *bool            `json:"include_login_method,omitempty"`
	HasOrgs                             *bool            `json:"has_orgs,omitempty"`
	MaxNumOrgsUsersCanBeIn              *int32           `json:"max_num_orgs_users_can_be_in,omitempty"`
	OrgsMetaname                        string           `json:"orgs_metaname,omitempty"`
	UsersCanCreateOrgs                  *bool            `json:"users_can_create_orgs,omitempty"`
	UsersCanDeleteTheirOwnOrgs          *bool            `json:"users_can_delete_their_own_orgs,omitempty"`
	UsersMustBeInAnOrganization         *bool            `json:"users_must_be_in_an_organization,omitempty"`
	OrgsCanSetupSaml                    *bool            `json:"orgs_can_setup_saml,omitempty"`
	UseOrgNameForSaml                   *bool            `json:"use_org_name_for_saml,omitempty"`
	DefaultToSamlLogin                  *bool            `json:"default_to_saml_login,omitempty"`
	SkipSamlRoleMappingStep             *bool            `json:"skip_saml_role_mapping_step,omitempty"`
	OrgsCanRequire2fa                   *bool            `json:"orgs_can_require_2fa,omitempty"`
	Theme                               *Theme           `json:"theme,omitempty"`
	LogoImageId                         string           `json:"logo_image_id,omitempty"`
	FaviconImageId                      string           `json:"favicon_image_id,omitempty"`
	BackgroundImageId                   string           `json:"background_image_id,omitempty"`
	PersonalApiKeysEnabled              *bool            `json:"personal_api_keys_enabled,omitempty"`
	PersonalApiKeyRateLimit             *RateLimitConfig `json:"personal_api_key_rate_limit,omitempty"`
	OrgApiKeysEnabled                   *bool            `json:"org_api_keys_enabled,omitempty"`
	OrgApiKeyRateLimit                  *RateLimitConfig `json:"org_api_key_rate_limit,omitempty"`
	InvalidateOrgApiKeysUponUserRemoval *bool            `json:"invalidate_org_api_key_upon_user_removal,omitempty"`
	ApiKeyConfig                        *ApiKeyConfig    `json:"api_key_config,omitempty"`
	OrgsCanViewOrgAuditLog              *bool            `json:"orgs_can_view_org_audit_log,omitempty"`
	AllOrgsCanViewOrgAuditLog           *bool            `json:"all_orgs_can_view_org_audit_log,omitempty"`
	OrgAuditLogIncludesImpersonation    *bool            `json:"org_audit_log_includes_impersonation,omitempty"`
	OrgAuditLogIncludesApiKeys          *bool            `json:"org_audit_log_includes_api_keys,omitempty"`
	OrgAuditLogIncludesEmployees        *bool            `json:"org_audit_log_includes_employees,omitempty"`
}

type EnvironmentConfigResponse struct {
	AllowUsersToSignupWithPersonalEmail bool            `json:"allow_users_to_signup_with_personal_email"`
	SignupDomainAllowlistEnabled        bool            `json:"signup_domain_allowlist_enabled"`
	SignupDomainAllowlist               []string        `json:"signup_domain_allowlist"`
	SignupDomainBlocklistEnabled        bool            `json:"signup_domain_blocklist_enabled"`
	SignupDomainBlocklist               []string        `json:"signup_domain_blocklist"`
	HasPasswordLogin                    bool            `json:"has_password_login"`
	HasPasswordlessLogin                bool            `json:"has_passwordless_login"`
	WaitlistUsersEnabled                bool            `json:"waitlist_users_enabled"`
	UserAutologoutSeconds               int64           `json:"user_autologout_seconds"`
	UserAutologoutType                  string          `json:"user_autologout_type"`
	UsersCanDeleteOwnAccount            bool            `json:"users_can_delete_own_account"`
	UsersCanChangeEmail                 bool            `json:"users_can_change_email"`
	IncludeLoginMethod                  bool            `json:"include_login_method"`
	HasOrgs                             bool            `json:"has_orgs"`
	MaxNumOrgsUsersCanBeIn              int32           `json:"max_num_orgs_users_can_be_in"`
	OrgsMetaname                        string          `json:"orgs_metaname"`
	UsersCanCreateOrgs                  bool            `json:"users_can_create_orgs"`
	UsersCanDeleteTheirOwnOrgs          bool            `json:"users_can_delete_their_own_orgs"`
	UsersMustBeInAnOrganization         bool            `json:"users_must_be_in_an_organization"`
	OrgsCanSetupSaml                    bool            `json:"orgs_can_setup_saml"`
	UseOrgNameForSaml                   bool            `json:"use_org_name_for_saml"`
	DefaultToSamlLogin                  bool            `json:"default_to_saml_login"`
	SkipSamlRoleMappingStep             bool            `json:"skip_saml_role_mapping_step"`
	OrgsCanRequire2fa                   bool            `json:"orgs_can_require_2fa"`
	Theme                               Theme           `json:"theme"`
	LogoUrl                             string          `json:"logo_url"`
	FaviconUrl                          string          `json:"favicon_url"`
	BackgroundUrl                       string          `json:"background_url"`
	PersonalApiKeysEnabled              bool            `json:"personal_api_keys_enabled"`
	PersonalApiKeyRateLimit             RateLimitConfig `json:"personal_api_key_rate_limit"`
	OrgApiKeysEnabled                   bool            `json:"org_api_keys_enabled"`
	OrgApiKeyRateLimit                  RateLimitConfig `json:"org_api_key_rate_limit"`
	InvalidateOrgApiKeyUponUserRemoval  bool            `json:"invalidate_org_api_key_upon_user_removal"`
	ApiKeyConfig                        ApiKeyConfig    `json:"api_key_config"`
	OrgsCanViewOrgAuditLog              bool            `json:"orgs_can_view_org_audit_log"`
	AllOrgsCanViewOrgAuditLog           bool            `json:"all_orgs_can_view_org_audit_log"`
	OrgAuditLogIncludesImpersonation    bool            `json:"org_audit_log_includes_impersonation"`
	OrgAuditLogIncludesApiKeys          bool            `json:"org_audit_log_includes_api_keys"`
	OrgAuditLogIncludesEmployees        bool            `json:"org_audit_log_includes_employees"`
}

type Theme struct {
	ThemeType                    string               `json:"theme_type"` // always "CustomV2" for now
	BodyFont                     string               `json:"font_family"`
	HeaderFont                   string               `json:"secondary_font_family"`
	DisplayProjectName           bool                 `json:"display_project_name"`
	LoginLayout                  string               `json:"login_ui_theme"`
	BackgroundType               string               `json:"background_type"`
	BackgroundColor              RgbColor             `json:"background_color"`
	SecondaryBackgroundColor     RgbColor             `json:"secondary_background_color"`      // secondary background color in gradient
	GradientAngle                int32                `json:"gradient_angle"`                  // angle of gradient in degrees
	BackgroundTextColor          RgbColor             `json:"background_text_color"`           // header text color in split layout
	SecondaryBackgroundTextColor RgbColor             `json:"secondary_background_text_color"` // subheader text color in split layout
	BorderColor                  RgbColor             `json:"border_color"`
	FrameBackgroundColor         RgbColor             `json:"foreground_color"`                // input/frame background color
	FrameTextColor               RgbColor             `json:"foreground_text_color"`           // frame/input text color
	FrameSecondaryTextColor      RgbColor             `json:"foreground_secondary_text_color"` // currently always same as ForegroundTextColor
	PrimaryColor                 RgbColor             `json:"success_button_color"`
	PrimaryTextColor             RgbColor             `json:"success_button_text_color"`
	ErrorButtonColor             RgbColor             `json:"error_button_color"`
	ErrorButtonTextColor         RgbColor             `json:"error_button_text_color"`
	Splitscreen                  *SplitscreenParams   `json:"splitscreen,omitempty"`
	ManagementPagesTheme         ManagementPagesTheme `json:"management_pages_theme"`
}

type RgbColor struct {
	Red   uint8 `json:"red"`
	Green uint8 `json:"green"`
	Blue  uint8 `json:"blue"`
}

type SplitscreenParams struct {
	Direction   string `json:"direction"`
	ContentType string `json:"content_type"`
	Header      string `json:"header"`
	Subheader   string `json:"subheader"`
}

type ManagementPagesTheme struct {
	DisplayNavbar         bool     `json:"display_sidebar"`
	MainBackgroundColor   RgbColor `json:"main_background_color"`
	MainTextColor         RgbColor `json:"main_text_color"`
	NavbarBackgroundColor RgbColor `json:"sidebar_background_color"`
	NavbarTextColor       RgbColor `json:"sidebar_text_color"`
	BorderColor           RgbColor `json:"border_color"`
	ActionButtonColor     RgbColor `json:"action_button_color"`
	ActionButtonTextColor RgbColor `json:"action_button_text_color"`
}

type ImageUploadResponse struct {
	ImageId string `json:"image_id"`
}

type ApiKeyConfig struct {
	ExpirationOptions ApiKeyExpirationOptionSettings `json:"expiration_options"`
}

type ApiKeyExpirationOptionSettings struct {
	Options ApiKeyExpirationOptions `json:"options"`
	Default string                  `json:"default"`
}

type ApiKeyExpirationOptions struct {
	TwoWeeks    bool `json:"TwoWeeks"`
	OneMonth    bool `json:"OneMonth"`
	ThreeMonths bool `json:"ThreeMonths"`
	SixMonths   bool `json:"SixMonths"`
	OneYear     bool `json:"OneYear"`
	Never       bool `json:"Never"`
}

type RateLimitConfig struct {
	PeriodType     string `json:"period_type"`
	PeriodSize     int32  `json:"period_size"`
	AllowPerPeriod int64  `json:"allow_per_period"`
}

type RealmConfigUpdate struct {
	AutoConfirmEmails                     *bool `json:"auto_confirm_emails,omitempty"`
	AllowPublicSignups                    *bool `json:"allow_public_signups,omitempty"`
	WaitlistUsersRequireEmailConfirmation *bool `json:"waitlist_users_require_email_confirmation,omitempty"`
}

type RealmConfigsResponse struct {
	Test    RealmConfigResponse  `json:"test"`
	Staging *RealmConfigResponse `json:"staging"`
	Prod    *RealmConfigResponse `json:"prod"`
}

type RealmConfigResponse struct {
	AutoConfirmEmails                     bool   `json:"auto_confirm_emails"`
	AllowPublicSignups                    bool   `json:"allow_public_signups"`
	WaitlistUsersRequireEmailConfirmation bool   `json:"waitlist_users_require_email_confirmation"`
	AuthHostname                          string `json:"auth_hostname"`
}

type UserProperties struct {
	Fields []UserProperty `json:"fields"`
}

type UserProperty struct {
	Name            string               `json:"name"`
	DisplayName     string               `json:"display_name"`
	FieldType       string               `json:"field_type"`
	Required        bool                 `json:"required"`
	RequiredBy      int64                `json:"required_by"`
	InJwt           bool                 `json:"in_jwt"`
	IsEnabled       bool                 `json:"is_enabled"`
	IsUserFacing    bool                 `json:"is_user_facing"`
	CollectOnSignup bool                 `json:"collect_on_signup"`
	CollectViaSaml  bool                 `json:"collect_via_saml"`
	ShowInAccount   bool                 `json:"show_in_account"`
	UserWritable    string               `json:"user_writable"`
	Metadata        userPropertyMetadata `json:"metadata"`
}

type userPropertyMetadata struct {
	TosLinks   []TosLink `json:"tos_links,omitempty"`
	EnumValues []string  `json:"enum_values,omitempty"`
}

type TosLink struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

type FeIntegrationInfoResponse struct {
	Test    TestFeIntegrationInfo   `json:"test"`
	Staging FeIntegrationInfoForEnv `json:"stage"`
	Prod    FeIntegrationInfoForEnv `json:"prod"`
}

type TestFeIntegrationInfo struct {
	AuthUrl                            string                             `json:"auth_url_origin"`
	LoginRedirectPath                  string                             `json:"login_redirect_path"`
	LogoutRedirectPath                 string                             `json:"logout_redirect_path"`
	AdditionalFeLocations              AdditionalFeLocations              `json:"allowed_urls"`
	TestEnvFeIntegrationApplicationUrl testEnvFeIntegrationApplicationUrl `json:"test_env"`
}

type FeIntegrationInfoForEnv struct {
	AuthUrl               string                `json:"auth_url_origin"`
	ApplicationUrl        string                `json:"application_hostname_with_scheme"`
	LoginRedirectPath     string                `json:"login_redirect_path"`
	LogoutRedirectPath    string                `json:"logout_redirect_path"`
	AdditionalFeLocations AdditionalFeLocations `json:"allowed_urls"`
	VerifiedDomain        string                `json:"verified_domain"`
}

type AdditionalFeLocations struct {
	AdditionalFeLocations []AdditionalFeLocation `json:"allowed_urls"`
}

type AdditionalFeLocation struct {
	Domain            string `json:"base_domain"`
	AllowAnySubdomain bool   `json:"allow_any_subdomain_match"`
}

type feIntegrationUpdateRequest struct {
	AdditionalFeLocations              AdditionalFeLocations               `json:"allowed_urls"`
	LogoutRedirectPath                 string                              `json:"logout_redirect_path"`
	LoginRedirectPath                  string                              `json:"login_redirect_path"`
	ApplicationHostnameWithScheme      string                              `json:"application_hostname_with_scheme,omitempty"`
	TestEnvFeIntegrationApplicationUrl *testEnvFeIntegrationApplicationUrl `json:"test_env,omitempty"`
}

type testEnvFeIntegrationApplicationUrl struct {
	ApplicationUrl string `json:"scheme_and_domain"`
	Port           int    `json:"port,omitempty"`
	Type           string `json:"type"`
}

type BeIntegrationInfoResponse struct {
	Test    BeIntegrationInfo `json:"test"`
	Staging BeIntegrationInfo `json:"stage"`
	Prod    BeIntegrationInfo `json:"prod"`
}

type BeIntegrationInfo struct {
	AuthUrl     string `json:"auth_url_origin"`
	VerifierKey string `json:"verifier_key"`
	Issuer      string `json:"issuer"`
}
type BeApiKeyCreateRequest struct {
	Name       string `json:"name"`
	IsReadOnly bool   `json:"readonly"`
}

type BeApiKeyUpdateRequest struct {
	ApiKeyId string `json:"api_key_id"`
	Name     string `json:"name"`
}

type BeApiKey struct {
	ApiKey     string `json:"api_key"`
	ApiKeyId   string `json:"api_key_id"`
	Name       string `json:"name"`
	IsReadOnly bool   `json:"readonly"`
}

type CustomDomainInfoResponse struct {
	Domain           string  `json:"domain"`
	Subdomain        *string `json:"subdomain"`
	IsVerified       bool    `json:"is_verified"`
	IsPending        bool    `json:"is_pending"`
	TxtRecordKey     *string `json:"txt_record_key"`
	TxtRecordValue   *string `json:"txt_record_value"`
	CnameRecordKey   *string `json:"cname_record_key"`
	CnameRecordValue *string `json:"cname_record_value"`
}

type customDomainUpdateRequest struct {
	Domain      string  `json:"domain"`
	Subdomain   *string `json:"subdomain,omitempty"`
	Environment string  `json:"environment"`
	IsSwitching bool    `json:"is_switching"`
}

type customDomainVerifyRequest struct {
	Environment string `json:"environment"`
	IsSwitching bool   `json:"is_switching"`
}

type RolesAndPermissions struct {
	Roles            []RoleDefinition `json:"roles"`
	Permissions      []Permission     `json:"available_external_permissions"`
	DefaultRole      string           `json:"default_role"`
	DefaultOwnerRole string           `json:"default_owner_role"`
	OrgRoleStructure string           `json:"org_role_structure"`
}

type RoleDefinition struct {
	Name                 string   `json:"name"`
	CanInvite            bool     `json:"can_invite"`
	CanChangeRoles       bool     `json:"can_change_roles"`
	CanManageApiKeys     bool     `json:"can_manage_api_keys"`
	CanRemoveUsers       bool     `json:"can_remove_users"`
	CanSetupSaml         bool     `json:"can_setup_saml"`
	CanViewOtherMembers  bool     `json:"can_view_other_members"`
	CanDeleteOrg         bool     `json:"can_delete_org"`
	CanEditOrgAccess     bool     `json:"can_edit_org_access"`
	CanUpdateOrgMetadata bool     `json:"can_update_org_metadata"`
	RolesCanManage       []string `json:"roles_can_manage,omitempty"`
	Disabled             bool     `json:"disabled"`
	IsVisibleToEndUser   bool     `json:"is_visible_to_end_user"`
	ExternalPermissions  []string `json:"external_permissions,omitempty"`
	Description          *string  `json:"description,omitempty"`
}

type Permission struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	DisplayName *string `json:"display_name,omitempty"`
}

type rolesAndPermissionsUpdate struct {
	RolesAndPermissions RolesAndPermissions `json:"org_definition"`
	RoleMigrationMap    RoleMigrationMap    `json:"role_to_role"`
}

type RoleMigrationMap struct {
	OldToNewRoleMapping map[string]*string `json:"role_map"`
}

type AllSocialLoginInfoResponse struct {
	Google     SocialLoginInfo `json:"google"`
	GitHub     SocialLoginInfo `json:"github"`
	Slack      SocialLoginInfo `json:"slack"`
	Microsoft  SocialLoginInfo `json:"microsoft"`
	LinkedIn   SocialLoginInfo `json:"linkedin"`
	Salesforce SocialLoginInfo `json:"salesforce"`
	Outreach   SocialLoginInfo `json:"outreach"`
	QuickBooks SocialLoginInfo `json:"quickbooks"`
	Xero       SocialLoginInfo `json:"xero"`
	Salesloft  SocialLoginInfo `json:"salesloft"`
	Atlassian  SocialLoginInfo `json:"atlassian"`
	Apple      SocialLoginInfo `json:"apple"`
}

type SocialLoginInfo struct {
	ClientId           string `json:"client_id"`
	TestRedirectUrl    string `json:"test_redirect_url"`
	StagingRedirectUrl string `json:"stage_redirect_url"`
	ProdRedirectUrl    string `json:"prod_redirect_url"`
}

type SocialLoginUpdateRequest struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Enabled      bool   `json:"enabled"`
}

type OauthClientRequest struct {
	RedirectUris []string `json:"redirect_uris"`
}

type OauthClientCreationResponse struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type OauthClientInfo struct {
	ClientId     string   `json:"client_id"`
	RedirectUris []string `json:"redirect_uris"`
}

type ApiKeyAlert struct {
	Enabled           bool  `json:"enabled"`
	AdvanceNoticeDays int32 `json:"advance_notice_days"`
}
