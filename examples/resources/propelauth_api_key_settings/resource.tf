# Configure the API key settings for your PropelAuth project.
resource "propelauth_api_key_settings" "example" {
  personal_api_keys_enabled                = true
  org_api_keys_enabled                     = true
  invalidate_org_api_key_upon_user_removal = true
  api_key_config = {
    expiration_options = {
      options = [
        "TwoWeeks",
        "OneMonth",
        "ThreeMonths",
      ]
      default = "ThreeMonths"
    }
  }
}
