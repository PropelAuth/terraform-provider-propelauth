---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "propelauth_organization_configuration Resource - propelauth"
subcategory: ""
description: |-
  Organization Configuration. This is for configuring your global organization settings in PropelAuth. Settings on specific organizations can be managed in the dashboard.
---

# propelauth_organization_configuration (Resource)

Organization Configuration. This is for configuring your global organization settings in PropelAuth. Settings on specific organizations can be managed in the dashboard.

## Example Usage

```terraform
# Configure how your global organization settings in PropelAuth.
resource "propelauth_organization_configuration" "example" {
  has_orgs                        = true
  orgs_metaname                   = "Company"
  max_num_orgs_users_can_be_in    = 1
  users_can_delete_their_own_orgs = true
  orgs_can_require_2fa            = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `customer_org_audit_log_settings` (Attributes) Settings for enabling whether and configuring how your customer organizations will have access to their own audit log.

Note: This feature is only available for use by your customer organizations in non-test environments for some pricing plans. (see [below for nested schema](#nestedatt--customer_org_audit_log_settings))
- `default_to_saml_login` (Boolean) This is an advanced setting that only applies if SAML is enabled. If true, affected users will be directed to SAML by default in the hosted pages.The default setting is false.
- `has_orgs` (Boolean) This is the top level setting for whether organizations are in your PropelAuth integration.If false, all other organization settings are ignored. The default setting is true.
- `max_num_orgs_users_can_be_in` (Number) This is the maximum number of organizations a user can be a member of. If a user tries to exceed this number, they will be asked to leave an existing organization. The default setting is 10.
- `orgs_can_require_2fa` (Boolean) If true, organizations can require their users to use 2FA.The default setting is false. Warning: This is only applied in prod for some billing plans
- `orgs_can_setup_saml` (Boolean) If true, your users can setup a SAML connection for their organization. This allows them to log into your product using their existing work account managed by an Identity Provider like Okta, Azure/Entra, Google, and more. The default setting is false. Warning: This is only applied in prod for some billing plans
- `orgs_metaname` (String) What name do you use for organizations? This will update the copy across your hosted pages.The default setting is 'Organization'.
- `skip_saml_role_mapping_step` (Boolean) This is an advanced setting that only applies if SAML is enabled. If true, end users setting up SAML for their organization will not see the role-mapping step. The default setting is false.
- `use_org_name_for_saml` (Boolean) This is an advanced setting that only applies if SAML is enabled. If true, users can look up and be redirected to their SSO provider using their organization's name.The default setting is false which means the SAML provider is instead inferred from their email address.
- `users_can_create_orgs` (Boolean) If true, users have access to the 'Create Org' UI, allowing them to create their own organizations.The default setting is true.
- `users_can_delete_their_own_orgs` (Boolean) If true, users with the requisite permission will be able to delete their organizations. The default setting is false.
- `users_must_be_in_an_organization` (Boolean) If true, users will be required to create or join an organization as part of the signup process. The default setting is false.

<a id="nestedatt--customer_org_audit_log_settings"></a>
### Nested Schema for `customer_org_audit_log_settings`

Required:

- `all_orgs_can_view_their_audit_log` (Boolean) If true, all of your customer organization will automatically have access to this feature. Otherwise, you will need to enable it for each organization individually.
- `enabled` (Boolean) If enabled, your customer organizations will have access to their own audit log.
- `include_api_key_actions` (Boolean) If true, the audit log will include actions that were triggered by your BE service utilizing PropelAuth APIs.
- `include_employee_actions` (Boolean) If true, the audit log will include actions that were triggered by a member of your team using the PropelAuth dashboard. The person who triggered the action will be anonymous to your customer.
- `include_impersonation` (Boolean) If true, the audit log will include actions that were triggered by a member of your team impersonating one of their organization members. The impersonator will be anonymous to your customer.

## Import

Import is supported using the following syntax:

```shell
# As there is only one organization_configuration per project there's no need to specify the id,
# but terraform import requires an id to be specified, so we can use an arbitrary string here.
terraform import propelauth_organization_configuration.example arbitrary_string_here
```
