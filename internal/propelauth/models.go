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
}
