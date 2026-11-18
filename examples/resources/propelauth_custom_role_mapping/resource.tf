# Configure a custom role mapping. Every role defined in the project must be
# present in the `roles` map, and each role's `is_internal` must match the
# project's role definition.
resource "propelauth_custom_role_mapping" "example" {
  name = "Example Mapping"
  roles = {
    "Owner" = {
      can_view_other_members  = true
      can_invite              = true
      can_change_roles        = true
      can_manage_api_keys     = true
      can_remove_users        = true
      can_setup_saml          = true
      can_delete_org          = true
      can_edit_org_access     = true
      can_update_org_metadata = true
      roles_can_manage        = ["Owner", "Admin", "Member", "Support"]
      permissions = [
        "doc::read",
        "doc::write",
        "ticket::read",
        "ticket::write"
      ]
    }
    "Admin" = {
      can_view_other_members = true
      can_invite             = true
      can_change_roles       = true
      can_remove_users       = true
      roles_can_manage       = ["Admin", "Member", "Support"]
      permissions = [
        "doc::read",
        "doc::write",
        "ticket::read"
      ]
    }
    "Member" = {
      permissions = [
        "doc::read",
        "ticket::read"
      ]
    }
    "Support" = {
      is_internal             = true
      can_view_other_members  = true
      can_invite              = false
      can_change_roles        = false
      can_manage_api_keys     = false
      can_remove_users        = false
      can_setup_saml          = false
      can_delete_org          = false
      can_edit_org_access     = false
      can_update_org_metadata = false
      permissions = [
        "doc::read",
        "ticket::read",
        "ai::deploy"
      ]
    }
  }
}
