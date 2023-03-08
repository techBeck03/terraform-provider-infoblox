package infoblox

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	cNameDomainName           = os.Getenv("INFOBLOX_DOMAIN")
	cNameRecordHostnameCreate = fmt.Sprintf("cname-infoblox-test.%s", cNameDomainName)
	cNameRecordHostnameUpdate = fmt.Sprintf("cname-infoblox-test-update.%s", cNameDomainName)
)

func TestAccInfobloxCnameRecordBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: composeConfig(testAccProviderBaseConfig, testAccCheckInfobloxARecordCreate(), testAccCheckInfobloxCnameRecordCreate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxCnameRecordExists("infoblox_a_record.new"),
					testAccCheckInfobloxCnameRecordExists("infoblox_cname_record.new"),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "alias", cNameRecordHostnameCreate),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "canonical", aRecordHostnameCreate),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "comment", "test cname record"),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "disable", "true"),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
					testAccCheckInfobloxCnameRecordExists("data.infoblox_cname_record.alias"),
					testAccCheckInfobloxCnameRecordExists("data.infoblox_cname_record.ref"),
					testAccCheckInfobloxCnameRecordExists("data.infoblox_cname_record.canonical"),
					resource.TestCheckResourceAttr("data.infoblox_cname_record.alias", "alias", cNameRecordHostnameCreate),
					resource.TestCheckResourceAttr("data.infoblox_cname_record.alias", "canonical", aRecordHostnameCreate),
					resource.TestCheckResourceAttr("data.infoblox_cname_record.alias", "comment", "test cname record"),
					resource.TestCheckResourceAttr("data.infoblox_cname_record.alias", "disable", "true"),
					resource.TestCheckResourceAttr("data.infoblox_cname_record.alias", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("data.infoblox_cname_record.alias", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("data.infoblox_cname_record.alias", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
				),
			},
			{
				Config: composeConfig(testAccProviderBaseConfig, testAccCheckInfobloxARecordCreate(), testAccCheckInfobloxCnameRecordUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxCnameRecordExists("infoblox_a_record.new"),
					testAccCheckInfobloxCnameRecordExists("infoblox_cname_record.new"),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "alias", cNameRecordHostnameUpdate),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "canonical", aRecordHostnameCreate),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "comment", "test cname record update"),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "disable", "false"),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "extensible_attributes.Location", "{\"value\":\"CollegeStation2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_cname_record.new", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
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
			return fmt.Errorf("Resource: %s not set", resourceName)
		}

		return nil
	}
}

var testAccCheckInfobloxCnameRecordCreate = fmt.Sprintf(`
resource "infoblox_cname_record" "new" {
  alias     = "%s"
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
data "infoblox_cname_record" "alias" {
  alias = infoblox_cname_record.new.alias
}
data "infoblox_cname_record" "ref" {
  ref = infoblox_cname_record.new.ref
}
data "infoblox_cname_record" "canonical" {
  canonical = infoblox_cname_record.new.canonical
}
`, cNameRecordHostnameCreate)

var testAccCheckInfobloxCnameRecordUpdate = fmt.Sprintf(`
resource "infoblox_cname_record" "new" {
  alias     = "%s"
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
`, cNameRecordHostnameUpdate)
