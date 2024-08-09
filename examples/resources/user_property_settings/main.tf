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

resource "propelauth_user_property_settings" "my_user_property_settings" {
  name_property = {
    in_jwt = false
  }
  metadata_property = {
    in_jwt = true
  }
  username_property = {
    in_jwt = true
    display_name = "Account Name"
  }
  # picture_url_property = {} # leaving this unset unsures the property remains disabled
  phone_number_property = {
    in_jwt = false
    required = true
    required_by = 0
    show_in_account = true
    collect_via_saml = true
    user_writable = "Write"
  }
  tos_property = {
    in_jwt = false
    required = true
    required_by = 0
    user_writable = "Write"
    tos_links = [
      {
        url = "https://example.com/tos",
        name = "Terms of Service"
      },
      {
        url = "https://example.com/privacy",
        name = "Privacy Policy"
      }
    ]
  }
  referral_source_property = {
    in_jwt = false
    display_name = "How did you find my awesome app?"
    required = true
    required_by = 0
    user_writable = "WriteIfUnset"
    options = [
        "Google",
        "Facebook",
        "Twitter",
        "LinkedIn",
        "Other"
    ]
  }
}