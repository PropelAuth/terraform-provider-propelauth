---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "propelauth_social_login Resource - propelauth"
subcategory: ""
description: |-
  Backend API Key resource. This is for configuring the basic BE API key information in PropelAuth.
---

# propelauth_social_login (Resource)

Backend API Key resource. This is for configuring the basic BE API key information in PropelAuth.

## Example Usage

```terraform
variable "github_client_secret" {
  type      = string
  sensitive = true
}

resource "propelauth_social_login" "github_sso" {
  social_provider = "GitHub"
  client_id       = "my-client-id"
  client_secret   = var.github_client_secret
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `client_id` (String) The client ID. This is a unique identifier for the oauth client that can be retrieved from the OIDC provider.
- `client_secret` (String, Sensitive) The client secret for the oauth client that can be retrieved from the OIDC provider.
- `social_provider` (String) The OIDC provider for the Social Login you're configuring. This is only for internal dislay purposes.Accepted values are `Google`, `Microsoft`, `GitHub`, `Slack`, `LinkedIn`, `Atlassian`, `Apple`, `Salesforce`, `QuickBooks`, `Xero`, `Salesloft`, and `Outreach`.

## Import

Import is supported using the following syntax:

```shell
# Import an existing social login integration by the social_provider name
terraform import propelauth_social_login.github_sso GitHub
# or
terraform import propelauth_social_login.mircosoft_sso Microsoft
```