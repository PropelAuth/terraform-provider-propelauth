package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserPropertySettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUserPropertySettingsResourceConfig(false, "https://example.com/tos", "Birthday"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"propelauth_user_property_settings.test",
						"phone_number_property.in_jwt",
						"false",
					),
					resource.TestCheckResourceAttr(
						"propelauth_user_property_settings.test",
						"tos_property.tos_links.0.url",
						"https://example.com/tos",
					),
					resource.TestCheckResourceAttr(
						"propelauth_user_property_settings.test",
						"custom_properties.0.display_name",
						"Birthday",
					),
				),
			},
			// Update and Read testing
			{
				Config: testAccUserPropertySettingsResourceConfig(true, "https://example.com/tos/v2", "Name Day"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"propelauth_user_property_settings.test",
						"phone_number_property.in_jwt",
						"true",
					),
					resource.TestCheckResourceAttr(
						"propelauth_user_property_settings.test",
						"tos_property.tos_links.0.url",
						"https://example.com/tos/v2",
					),
					resource.TestCheckResourceAttr(
						"propelauth_user_property_settings.test",
						"custom_properties.0.display_name",
						"Name Day",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccUserPropertySettingsResourceConfig(phoneInJwt bool, tosLink string, birthdayName string) string {
	return providerConfig + fmt.Sprintf(`
resource "propelauth_user_property_settings" "test" {
  name_property = {
    in_jwt = false
  }
  metadata_property = {
    in_jwt = true
  }
  username_property = {
    in_jwt = true
    display_name = "Account Name"
  }
  # picture_url_property = {} # leaving this unset unsures the property remains disabled
  phone_number_property = {
    in_jwt = %v
    required = true
    required_by = 0
    show_in_account = true
    collect_via_saml = true
    user_writable = "Write"
  }
  tos_property = {
    in_jwt = false
    required = true
    required_by = 0
    user_writable = "Write"
    tos_links = [
      {
        url = "%v",
        name = "Terms of Service"
      },
      {
        url = "https://example.com/privacy",
        name = "Privacy Policy"
      }
    ]
  }
  referral_source_property = {
    in_jwt = false
    display_name = "How did you find my awesome app?"
    required = true
    required_by = 0
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
      name = "birthday"
      display_name = "%v"
      field_type = "Date"
      in_jwt = true
      required = false
      user_writable = "Write"
    },
    {
      name = "favorite_ice_cream_flavor"
      display_name = "Favorite Ice Cream Flavor"
      field_type = "Enum"
      enum_values = [
        "Vanilla",
        "Chocolate",
        "Strawberry",
        "Mint Chocolate Chip",
        "Other"
      ]
      in_jwt = true
      required = true
      required_by = 0
      user_writable = "Write"
    },
    {
      name = "favorite_color"
      display_name = "Favorite Color"
      field_type = "Text"
      in_jwt = false
      required = false
      user_writable = "Write"
    },
    {
      name = "receive_newsletter"
      display_name = "I want to receive the newsletter"
      field_type = "Toggle"
      user_writable = "Write"
      show_in_account = false
    }
  ]
}
`, phoneInJwt, tosLink, birthdayName)
}
