terraform {
  required_providers {
    propelauth = {
      source = "registry.terraform.io/propelauth/propelauth"
    }
  }
}

provider "propelauth" {
#   tenant_id  = "<PROPELAUTH_TENANT_ID>"  # or PROPELAUTH_TENANT_ID environment variable
#   project_id = "<PROPELAUTH_PROJECT_ID>" # or PROPELAUTH_PROJECT_ID environment variable
#   api_key    = "<PROPELAUTH_API_KEY>"    # or PROPELAUTH_API_KEY environment variable
}

resource "propelauth_roles_and_permissions" "my_roles_and_permissions" {
  permissions = [
	  {
	    name = "doc::read"
	    display_name = "Can read documents." # optional
	    description = "A description here." # optional
	  },
	  {
	    name = "doc::write"
	  },
	  {
	    name = "ticket::read"
	  },
	  {
	    name = "ticket::write"
	  },
	  {
	    name = "ai::deploy"
	  }
	]
	roles = {
    "Owner" = {
      can_view_other_members = true
      can_invite = true
      can_change_roles = true
      can_manage_api_keys = false
      can_remove_users = true
      can_setup_saml = true
      can_delete_org = true
      can_edit_org_access = false
      can_update_org_metadata = false
      permissions = [
        "doc::read",
        "doc::write"
      ]
    }
    "Admin" = {
      can_view_other_members = true
      can_invite = true
      can_change_roles = true
      can_manage_api_keys = false
      can_remove_users = true
      can_setup_saml = false
      can_delete_org = false
      can_edit_org_access = false
      can_update_org_metadata = false
      permissions = [
        "doc::read",
        "doc::write"
      ]
    }
    "Member" = {
      permissions = [
        "doc::read"
      ]
    }
	}
	role_hierarchy = [ # only for multiple_roles_per_user = false
	  "Owner",
	  "Admin",
	  "Member"
	]
	default_role = "Member"
	default_owner_role = "Owner"
	# default_mapping_name = "Free Tier"
}

resource "propelauth_role_permissions_override" "premium_permissions" {
  depends_on = [ propelauth_roles_and_permissions.my_roles_and_permissions ]
  name = "Premium"
  roles = {
    "Owner" = {
      can_manage_api_keys = true
      can_edit_org_access = true
      can_update_org_metadata = true
      permissions = [
        "doc::read",
        "doc::write",
        "ticket::read",
        "ticket::write",
        "ai::deploy"
      ]
    }
    "Admin" = {
      can_manage_api_keys = true
      can_edit_org_access = true
      can_update_org_metadata = true
      permissions = [
        "doc::read",
        "doc::write",
        "ticket::read",
        "ai::deploy"
      ]
    }
    "Member" = {
      permissions = [
        "doc::read",
        "ticket::read"
      ]
    }
  }
}

output "test_role_permissions_override" {
  value = propelauth_role_permissions_override.premium_permissions
}
