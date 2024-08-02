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
}
