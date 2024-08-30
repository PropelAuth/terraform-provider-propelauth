package propelauth

import (
	"encoding/json"
)

// ValidateRolesAndPermissions - Validates an update to roles and permissions without applying it.
func (c *PropelAuthClient) ValidateRolesAndPermissions(candidateUpdate rolesAndPermissionsUpdate) (bool, error) {
	updateJson, err := json.Marshal(candidateUpdate)
	if err != nil {
		return false, err
	}

	_, err = c.post("roles_and_permissions/validate", updateJson)
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetRolesAndPermissions - Returns the roles and permissions.
func (c *PropelAuthClient) GetRolesAndPermissions() (*RolesAndPermissions, error) {
	res, err := c.get("roles_and_permissions", nil)
	if err != nil {
		return nil, err
	}

	rolesAndPermissions := RolesAndPermissions{}
	err = json.Unmarshal(res.BodyBytes, &rolesAndPermissions)
	if err != nil {
		return nil, err
	}

	return &rolesAndPermissions, nil
}

// UpdateRolesAndPermissions - Updates the roles and permissions.
func (c *PropelAuthClient) UpdateRolesAndPermissions(update rolesAndPermissionsUpdate) (*RolesAndPermissions, error) {
	updateJson, err := json.Marshal(update)
	if err != nil {
		return nil, err
	}

	_, err = c.post("roles_and_permissions", updateJson)
	if err != nil {
		return nil, err
	}

	return c.GetRolesAndPermissions()
}

type RolesAndPermissionsUpdateBuilder struct {
	multipleRolesPerUser bool
	defaultRole string
	defaultOwnerRole string
	roleHierarchy []string
	roles map[string]RoleDefinition
	permissions []Permission
	oldToNewRoleMapping map[string]*string
	oldRoleNames []string
}

func NewRolesAndPermissionsUpdateBuilder() *RolesAndPermissionsUpdateBuilder {
	return &RolesAndPermissionsUpdateBuilder{
		roles: make(map[string]RoleDefinition),
		permissions: []Permission{},
		oldToNewRoleMapping: make(map[string]*string),
		oldRoleNames: make([]string, 0),
	}
}

func (b *RolesAndPermissionsUpdateBuilder) SetMultipleRolesPerUser(multipleRolesPerUser bool) *RolesAndPermissionsUpdateBuilder {
	b.multipleRolesPerUser = multipleRolesPerUser
	return b
}

func (b *RolesAndPermissionsUpdateBuilder) SetDefaultRole(defaultRole string) *RolesAndPermissionsUpdateBuilder {
	b.defaultRole = defaultRole
	return b
}

func (b *RolesAndPermissionsUpdateBuilder) SetDefaultOwnerRole(defaultOwnerRole string) *RolesAndPermissionsUpdateBuilder {
	b.defaultOwnerRole = defaultOwnerRole
	return b
}

func (b *RolesAndPermissionsUpdateBuilder) SetRoleHierarchy(roleHierarchy []string) *RolesAndPermissionsUpdateBuilder {
	b.roleHierarchy = roleHierarchy
	return b
}

func (b *RolesAndPermissionsUpdateBuilder) InsertRole(roleName string, roleDefinition RoleDefinition) *RolesAndPermissionsUpdateBuilder {
	b.roles[roleName] = roleDefinition
	return b
}

func (b *RolesAndPermissionsUpdateBuilder) InsertPermission(permission Permission) *RolesAndPermissionsUpdateBuilder {
	b.permissions = append(b.permissions, permission)
	return b
}

func (b *RolesAndPermissionsUpdateBuilder) InsertOldToNewRoleMapping(oldRoleName string, newRoleName string) *RolesAndPermissionsUpdateBuilder {
	b.oldToNewRoleMapping[oldRoleName] = &newRoleName
	return b
}

func (b *RolesAndPermissionsUpdateBuilder) InsertOldRoleName(oldRoleName string) *RolesAndPermissionsUpdateBuilder {
	b.oldRoleNames = append(b.oldRoleNames, oldRoleName)
	return b
}

func (b *RolesAndPermissionsUpdateBuilder) Build() rolesAndPermissionsUpdate {
	updateRequest := rolesAndPermissionsUpdate{
		RolesAndPermissions: RolesAndPermissions{},
		RoleMigrationMap: RoleMigrationMap{
			OldToNewRoleMapping: make(map[string]*string),
		},
	}

	if b.multipleRolesPerUser {
		updateRequest.RolesAndPermissions.OrgRoleStructure = "multi_role"
	} else {
		updateRequest.RolesAndPermissions.OrgRoleStructure = "single_role_in_hierarchy"
	}

	updateRequest.RolesAndPermissions.DefaultRole = b.defaultRole
	updateRequest.RolesAndPermissions.DefaultOwnerRole = b.defaultOwnerRole
	updateRequest.RolesAndPermissions.Permissions = b.permissions
	
	if b.multipleRolesPerUser {
		for _, role := range b.roles {
			updateRequest.RolesAndPermissions.Roles = append(updateRequest.RolesAndPermissions.Roles, role)
		}
	} else {
		for _, role := range b.roleHierarchy {
			updateRequest.RolesAndPermissions.Roles = append(updateRequest.RolesAndPermissions.Roles, b.roles[role])
		}
	}

	// build role-to-role mapping for PropelAuth to migrate roles that have changed names or been deleted
	// 1. roles that have name changes
	oldToNewRoleMapping := b.oldToNewRoleMapping
	// 2. roles that are unchanged or have been added
	for roleName := range b.roles {
		_, ok := oldToNewRoleMapping[roleName]
		if ok {
			continue
		} else {
			forcedCopyOfRoleName := roleName // need a hard copy since go reuses the same variable when ranging over a map
			oldToNewRoleMapping[roleName] = &forcedCopyOfRoleName
		}
	}
	// 3. roles that have been removed
	for _, roleName := range b.oldRoleNames {
		_, ok := oldToNewRoleMapping[roleName]
		if ok {
			continue
		} else {
			oldToNewRoleMapping[roleName] = nil
		}
	}
	updateRequest.RoleMigrationMap.OldToNewRoleMapping = oldToNewRoleMapping

	return updateRequest
}

func (r *RolesAndPermissions) GetPermission(permissionName string) (*Permission, bool) {
	for _, permission := range r.Permissions {
		if permission.Name == permissionName {
			return &permission, true
		}
	}
	return nil, false
}

func (r *RolesAndPermissions) GetHeirarchy() []string {
	if r.OrgRoleStructure == "multi_role" {
		return nil
	} else {
		hierarchy := make([]string, len(r.Roles))
		for i, role := range r.Roles {
			hierarchy[i] = role.Name
		}
		return hierarchy
	}
}