package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccImageResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccImageResourceConfig("logo"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"propelauth_image.test",
						"image_type",
						"logo",
					),
				),
			},
			// Update and Read testing
			{
				Config: testAccImageResourceConfig("favicon"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"propelauth_image.test",
						"image_type",
						"favicon",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImageResourceConfig(imageType string) string {
	return providerConfig + fmt.Sprintf(`
resource "propelauth_image" "test" {
  source = "${path.module}/../../examples/resources/propelauth_theme/git-merge.png"
  version = "0.0.0"
  image_type = %[1]q
}
`, imageType)
}
