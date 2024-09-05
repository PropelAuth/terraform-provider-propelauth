# Configure how your global organization settings in PropelAuth.
resource "propelauth_organization_configuration" "example" {
  has_orgs                        = true
  orgs_metaname                   = "Company"
  max_num_orgs_users_can_be_in    = 1
  users_can_delete_their_own_orgs = true
  orgs_can_require_2fa            = true
}