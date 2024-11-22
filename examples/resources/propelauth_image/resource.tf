# Set the logo for your PropelAuth project.
resource "propelauth_image" "logo_example" {
  source     = "${path.module}/example-logo-image.png"
  version    = "0.1.0"
  image_type = "logo"
}

# Set the favicon for your PropelAuth project.
resource "propelauth_image" "favicon_example" {
  source     = "${path.module}/example-favicon-image.png"
  version    = "0.1.0"
  image_type = "favicon"
}

# And in the case where you've updated the image file at the same source path,
# you can increment the version to force the update in PropelAuth.
resource "propelauth_image" "background_example" {
  source     = "${path.module}/example-bg-image.png"
  version    = "0.1.1"
  image_type = "background"
}