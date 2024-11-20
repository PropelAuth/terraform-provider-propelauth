---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "propelauth_environment_level_auth_configuration Resource - propelauth"
subcategory: ""
description: |-
  Environment-level Auth Configuration. This is for configuring elements of the signup and login experience in PropelAuth that you may want to differ between test and production environments.
---

# propelauth_environment_level_auth_configuration (Resource)

Environment-level Auth Configuration. This is for configuring elements of the signup and login experience in PropelAuth that you may want to differ between test and production environments.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment` (String) The environment for which you are configuring the login and signup experience. Accepted values are `Test`, `Staging`, and `Prod`.

### Optional

- `allow_public_signups` (Boolean) If true, new users will be able to sign up for your product directly in the PropelAuth hosted pages.The default setting is true for all environments.
- `require_email_confirmation` (Boolean) If true, all users are required to have confirmed email addresses. Whenever PropelAuth doesn't know for certain whether a user's email adderss is in fact owned by them, PropelAuth will trigger an email confirmation flow. The default setting is true for `Prod` and `Staging` environments but is false for `Test` for ease of development.
- `waitlist_users_require_email_confirmation` (Boolean) If true, all waitlisted users are required to have confirmed email addresses. Whenever PropelAuth doesn't know for certain whether a waitlisted user's email adderss is in fact owned by them, PropelAuth will trigger an email confirmation flow. The default setting is false for all environments.