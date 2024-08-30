package propelauth

import (
	"encoding/json"
	"fmt"
)

// CreateRolePermissionsOverride - Creates a role mapping.
func (c *PropelAuthClient) CreateRolePermissionsOverride(rolesOverride RolePermissionsOverride) (*string, error) {
	rolesOverrideJson, err := json.Marshal(rolesOverride)
	if err != nil {
		return nil, err
	}

	res, err := c.post("custom_role_mappings", rolesOverrideJson)
	if err != nil {
		return nil, err
	}
	createdRoleMapping := RolePermissionsOverrideCreationResponse{}
	err = json.Unmarshal(res.BodyBytes, &createdRoleMapping)
	if err != nil {
		return nil, err
	}

	return &createdRoleMapping.MappingId, nil
}

// GetRolePermissionsOverride - Returns the role mapping for the given id.
func (c *PropelAuthClient) GetRolePermissionsOverride(mappingId string) (*RolePermissionsOverride, error) {
	res, err := c.get(fmt.Sprintf("custom_role_mappings/%s", mappingId), nil)
	if err != nil {
		return nil, err
	}

	rolesOverride := RolePermissionsOverride{}
	err = json.Unmarshal(res.BodyBytes, &rolesOverride)
	if err != nil {
		return nil, err
	}

	return &rolesOverride, nil
}

// UpdateRolePermissionsOverride - Updates the role mapping for the given id.
func (c *PropelAuthClient) UpdateRolePermissionsOverride(mappingId string, rolesOverride RolePermissionsOverride) (*RolePermissionsOverride, error) {
	rolesOverrideJson, err := json.Marshal(rolesOverride)
	if err != nil {
		return nil, err
	}

	_, err = c.put(fmt.Sprintf("custom_role_mappings/%s", mappingId), rolesOverrideJson)
	if err != nil {
		return nil, err
	}

	return c.GetRolePermissionsOverride(mappingId)
}

// DeleteRolePermissionsOverride - Deletes the role mapping for the given id.
func (c *PropelAuthClient) DeleteRolePermissionsOverride(mappingId string) error {
	_, err := c.delete(fmt.Sprintf("custom_role_mappings/%s", mappingId), nil)
	if err != nil {
		return err
	}

	return nil
}
