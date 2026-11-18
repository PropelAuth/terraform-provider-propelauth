package propelauth

import (
	"encoding/json"
)

// Custom role mapping endpoints on the infrastructure API:
//   GET    custom_role_mappings/{id} - read a mapping and its roles
//   PATCH  custom_role_mappings/{id} - replace a mapping's roles (body mirrors the GET shape)
//   POST   custom_role_mappings      - create a mapping, returns {"mapping_id": "..."}
//   DELETE custom_role_mappings/{id} - delete a mapping

type CustomRoleMapping struct {
	MappingId string                  `json:"mapping_id,omitempty"`
	Name      string                  `json:"name"`
	Roles     []CustomRoleMappingRole `json:"roles"`
}

type CustomRoleMappingRole struct {
	Name                    string   `json:"name"`
	CanInvite               bool     `json:"can_invite"`
	CanChangeRoles          bool     `json:"can_change_roles"`
	CanManageApiKeys        bool     `json:"can_manage_api_keys"`
	CanRemoveUsers          bool     `json:"can_remove_users"`
	CanSetupSaml            bool     `json:"can_setup_saml"`
	CanViewOtherMembers     bool     `json:"can_view_other_members"`
	CanDeleteOrg            bool     `json:"can_delete_org"`
	CanEditOrgAccess        bool     `json:"can_edit_org_access"`
	CanUpdateOrgMetadata    bool     `json:"can_update_org_metadata"`
	CanViewOrgAuditLog      bool     `json:"can_view_org_audit_log"`
	CanManageMcpConnections bool     `json:"can_manage_mcp_connections"`
	RolesCanManage          []string `json:"roles_can_manage"`
	Disabled                bool     `json:"disabled"`
	Deprecated              bool     `json:"deprecated"`
	Description             *string  `json:"description,omitempty"`
	ShowDescriptionTooltip  *bool    `json:"show_description_tooltip,omitempty"`
	ExternalPermissions     []string `json:"external_permissions"`
	McpOrgScopes            []string `json:"mcp_org_scopes"`
	IsVisibleToEndUser      bool     `json:"is_visible_to_end_user"`
}

type createCustomRoleMappingResponse struct {
	MappingId string `json:"mapping_id"`
}

// GetCustomRoleMapping - Returns a custom role mapping and its roles.
func (c *PropelAuthClient) GetCustomRoleMapping(mappingId string) (*CustomRoleMapping, error) {
	res, err := c.get("custom_role_mappings/" + mappingId)
	if err != nil {
		return nil, err
	}

	mapping := CustomRoleMapping{}
	if err := json.Unmarshal(res.BodyBytes, &mapping); err != nil {
		return nil, err
	}

	return &mapping, nil
}

// CreateCustomRoleMapping - Creates a custom role mapping. The response carries
// only the new mapping's id, so the created mapping is fetched for the return.
func (c *PropelAuthClient) CreateCustomRoleMapping(mapping CustomRoleMapping) (*CustomRoleMapping, error) {
	body, err := json.Marshal(mapping)
	if err != nil {
		return nil, err
	}

	res, err := c.post("custom_role_mappings", body)
	if err != nil {
		return nil, err
	}

	created := createCustomRoleMappingResponse{}
	if err := json.Unmarshal(res.BodyBytes, &created); err != nil {
		return nil, err
	}

	return c.GetCustomRoleMapping(created.MappingId)
}

// UpdateCustomRoleMapping - Replaces a custom role mapping's roles.
func (c *PropelAuthClient) UpdateCustomRoleMapping(mappingId string, mapping CustomRoleMapping) (*CustomRoleMapping, error) {
	mapping.MappingId = mappingId
	body, err := json.Marshal(mapping)
	if err != nil {
		return nil, err
	}

	if _, err := c.patch("custom_role_mappings/"+mappingId, body); err != nil {
		return nil, err
	}

	return c.GetCustomRoleMapping(mappingId)
}

// DeleteCustomRoleMapping - Deletes a custom role mapping. A missing mapping is
// treated as already deleted.
func (c *PropelAuthClient) DeleteCustomRoleMapping(mappingId string) error {
	_, err := c.delete("custom_role_mappings/"+mappingId, nil)
	if err != nil && !IsPropelAuthNotFoundError(err) {
		return err
	}

	return nil
}
