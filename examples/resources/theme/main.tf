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

resource "propelauth_image" "logo" {
  source     = "${path.module}/git-merge.png"
  version    = "0.0.0"
  image_type = "logo"
}

resource "propelauth_image" "favicon" {
  source     = "${path.module}/git-merge.png"
  version    = "0.0.0"
  image_type = "favicon"
}

resource "propelauth_theme" "my_theme" {
  header_font = "Fraunces"
  body_font   = "PlusJakartaSans"
  login_page_theme = {
    layout          = "Frameless"
    background_type = "Gradient"
    gradient_background_parameters = {
      background_gradient_start_color = "#0c1cf7"
      background_gradient_end_color   = "#000000"
      background_gradient_angle       = 45
      background_text_color           = "#ffffff"
    }
    frame_background_color  = "#ffffff"
    frame_text_color        = "#700278"
    primary_color           = "#02927d"
    primary_text_color      = "#ffffff"
    error_color             = "#cf222e"
    error_button_text_color = "#ffffff"
    border_color            = "#000000"
  }
  management_pages_theme = {
    main_background_color    = "#c1a0eb"
    main_text_color          = "#2d4036"
    navbar_background_color  = "#e4c7f2"
    navbar_text_color        = "#2d4036"
    action_button_color      = "#629c75"
    action_button_text_color = "#f7f7f7"
    border_color             = "#fcfcfc"
  }
}

output "theme_result" {
  value = propelauth_theme.my_theme
}