# Configure all the properties that can be set on a user in your PropelAuth project.
resource "propelauth_user_property_settings" "example" {
  name_property = {
    in_jwt = false
  }
  metadata_property = {
    in_jwt = true
  }
  username_property = {
    in_jwt       = true
    display_name = "Account Name"
  }
  # picture_url_property = {} # leaving this unset unsures the property remains disabled
  phone_number_property = {
    in_jwt           = false
    required         = true
    required_by      = 0
    show_in_account  = true
    collect_via_saml = true
    user_writable    = "Write"
  }
  tos_property = {
    in_jwt        = false
    required      = true
    required_by   = 0
    user_writable = "Write"
    tos_links = [
      {
        url  = "https://example.com/tos",
        name = "Terms of Service"
      },
      {
        url  = "https://example.com/privacy",
        name = "Privacy Policy"
      }
    ]
  }
  referral_source_property = {
    in_jwt        = false
    display_name  = "How did you find my awesome app?"
    required      = true
    required_by   = 0
    user_writable = "WriteIfUnset"
    options = [
      "Google",
      "Facebook",
      "Twitter",
      "LinkedIn",
      "Other"
    ]
  }
  custom_properties = [
    {
      name          = "birthday"
      display_name  = "Birthday"
      field_type    = "Date"
      in_jwt        = true
      required      = false
      user_writable = "Write"
    },
    {
      name         = "favorite_ice_cream_flavor"
      display_name = "Favorite Ice Cream Flavor"
      field_type   = "Enum"
      enum_values = [
        "Vanilla",
        "Chocolate",
        "Strawberry",
        "Mint Chocolate Chip",
        "Other"
      ]
      in_jwt        = true
      required      = true
      required_by   = 0
      user_writable = "Write"
    },
    {
      name          = "favorite_color"
      display_name  = "Favorite Color"
      field_type    = "Text"
      in_jwt        = false
      required      = false
      user_writable = "Write"
    },
    {
      name            = "receive_newsletter"
      display_name    = "I want to receive the newsletter"
      field_type      = "Toggle"
      user_writable   = "Write"
      show_in_account = false
    }
  ]
}