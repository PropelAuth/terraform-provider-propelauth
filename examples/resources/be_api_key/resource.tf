# Create a PropelAuth API key for use in the Back End of your application.
resource "propelauth_be_api_key" "example" {
  environment = "Prod"
  name        = "Test API Key"
  read_only   = false
}

output "be_api_key_result" {
  value     = propelauth_be_api_key.example.api_key
  sensitive = true
}