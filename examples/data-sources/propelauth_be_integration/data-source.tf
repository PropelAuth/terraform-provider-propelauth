# Retrieve the details of a Back-end Integration to PropelAuth by environment.
data "propelauth_be_integration" "example" {
  environment = "Test"
}