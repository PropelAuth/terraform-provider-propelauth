# Configure how the front-end of your application integrates with PropelAuth.
resource "propelauth_fe_integration" "example" {
  environment          = "Test"
  application_url      = "http://localhost:3001"
  login_redirect_path  = "/home"
  logout_redirect_path = "/goodbye"
  additional_fe_locations = [
    {
      domain              = "http://localhost:3002"
      allow_any_subdomain = false
    }
  ]
}