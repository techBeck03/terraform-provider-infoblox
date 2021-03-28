package infoblox

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var cNameDomainName = os.Getenv("INFOBLOX_DOMAIN")

func TestAccInfobloxCnameRecordBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: composeConfig(testProviderARecordCreate, testProviderCnameRecordCreate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxCnameRecordExists("infoblox_a_record.new"),
					testAccCheckInfobloxCnameRecordExists("infoblox_cname_record.new"),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "alias", fmt.Sprintf("alias-infoblox-test.%s", cNameDomainName)),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "canonical", fmt.Sprintf("infoblox-test.%s", cNameDomainName)),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "comment", "test cname record"),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "disable", "true"),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\"}"),
				),
			},
			{
				Config: composeConfig(testProviderARecordCreate, testProviderCnameRecordUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxCnameRecordExists("infoblox_a_record.new"),
					testAccCheckInfobloxCnameRecordExists("infoblox_cname_record.new"),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "alias", fmt.Sprintf("alias-infoblox-test2.%s", cNameDomainName)),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "canonical", fmt.Sprintf("infoblox-test.%s", cNameDomainName)),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "comment", "test cname record update"),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "disable", "false"),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "extensible_attributes.Location", "{\"value\":\"CollegeStation2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins2\",\"type\":\"STRING\"}"),
				),
			},
		},
	})
}

func testAccCheckInfobloxCnameRecordExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No connection group set")
		}

		return nil
	}
}

var testProviderCnameRecordCreate = fmt.Sprintf(`
resource "infoblox_cname_record" "new" {
  alias     = "alias-infoblox-test.%s"
  comment   = "test cname record"
  canonical = infoblox_a_record.new.hostname
  disable   = true
  extensible_attributes = {
    Owner = jsonencode({
      value = "leroyjenkins",
      type  = "STRING"
    })
    Location = jsonencode({
      value = "CollegeStation",
      type  = "STRING"
    })
  }
}
`, cNameDomainName)

var testProviderCnameRecordUpdate = fmt.Sprintf(`
resource "infoblox_cname_record" "new" {
  alias     = "alias-infoblox-test2.%s"
  comment   = "test cname record update"
  canonical = infoblox_a_record.new.hostname
  disable   = false
  extensible_attributes = {
    Owner = jsonencode({
      value = "leroyjenkins2",
      type  = "STRING"
    })
    Location = jsonencode({
      value = "CollegeStation2",
      type  = "STRING"
    })
  }
}
`, cNameDomainName)
