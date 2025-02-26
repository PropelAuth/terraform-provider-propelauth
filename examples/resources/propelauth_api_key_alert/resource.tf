# Create a PropelAuth API key alert.
resource "propelauth_api_key_alert" "example" {
  # Cannot setup api key alerting for end users without a production environment.
  depends_on          = [propelauth_custom_domain_verification.my_prod_domain_verification]
  advance_notice_days = 30
}
