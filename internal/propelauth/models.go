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