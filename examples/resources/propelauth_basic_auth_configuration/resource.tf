# Configure basic authentication settings for your PropelAuth project.
resource "propelauth_basic_auth_configuration" "example" {
  has_password_login     = true
  has_passwordless_login = true
}