---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "propelauth_social_login_redirect Data Source - propelauth"
subcategory: ""
description: |-
  Retrieves the redirect url needed for configuring an OIDC provider for a Social Login in PropelAuth.
---

# propelauth_social_login_redirect (Data Source)

Retrieves the redirect url needed for configuring an OIDC provider for a Social Login in PropelAuth.

## Example Usage

```terraform
# Retrieve the redirect url needed for configuring a google social login 
# in your test environment.
data "propelauth_social_login_redirect" "google_test_redirect" {
  environment = "Test"
  provider    = "Google"
}

output "google_test_redirect_result" {
  value = propelauth_social_login_redirect.google_test_redirect.redirect_url
}

# To do the same for a production or staging environment,
# you must first verify a domain for the target environment,
# or no social login redirect URLs will be available.
data "propelauth_social_login_redirect" "google_prod_redirect" {
  environment = "Prod"
  provider    = "Google"
  depends_on  = [data.propelauth_custom_domain_verification.my_prod_custom_domain_verification]
}

output "google_prod_redirect_result" {
  value = propelauth_social_login_redirect.google_prod_redirect.redirect_url
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment` (String) The environment for which you are configuring the social login. Accepted values are `Test`, `Staging`, and `Prod`.
- `social_provider` (String) The social login provider for which you are configuring and need the redirect URL. Accepted values are `Google`, `Microsoft`, `GitHub`, `Slack`, `LinkedIn`, `Atlassian`, `Apple`, `Salesforce`, `QuickBooks`, `Xero`, `Salesloft`, and `Outreach`.

### Read-Only

- `redirect_url` (String) The redirect URL to be white-listed in the OIDC configuration of the social login provider.