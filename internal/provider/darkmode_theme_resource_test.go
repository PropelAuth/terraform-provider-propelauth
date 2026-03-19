package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDarkmodeThemeResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDarkmodeThemeResourceConfig("#f70000"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"propelauth_darkmode_theme.test",
						"login_page_theme.gradient_background_parameters.background_gradient_end_color",
						"#f70000",
					),
				),
			},
			// Update and Read testing
			{
				Config: testAccDarkmodeThemeResourceConfig("#000000"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"propelauth_darkmode_theme.test",
						"login_page_theme.gradient_background_parameters.background_gradient_end_color",
						"#000000",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDarkmodeThemeResourceConfig(backgroundGradientEndColor string) string {
	return providerConfig + fmt.Sprintf(`
resource "propelauth_darkmode_theme" "test" {
  header_font = "Fraunces"
  body_font = "PlusJakartaSans"
  login_page_theme = {
    layout = "Frameless"
    background_type = "Gradient"
    gradient_background_parameters = {
      background_gradient_start_color = "#0c1cf7"
      background_gradient_end_color = %[1]q
      background_gradient_angle = 45
      background_text_color = "#000000"
    }
    frame_background_color = "#000000"
    frame_text_color = "#700278"
    primary_color = "#02927d"
    primary_text_color = "#ffffff"
    error_color = "#cf222e"
    error_button_text_color = "#ffffff"
    border_color = "#ffffff"
  }
  management_pages_theme = {
    main_background_color = "#000000"
    main_text_color = "#495a51"
    navbar_background_color = "#e4c7f2"
    navbar_text_color = "#495a51"
    action_button_color = "#629c75"
    action_button_text_color = "#f7f7f7"
    border_color = "#fcfcfc"
  }
}
`, backgroundGradientEndColor)
}
