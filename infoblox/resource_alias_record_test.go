package infoblox

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	aliasRecordDomainName     = os.Getenv("INFOBLOX_DOMAIN")
	aliasRecordHostnameCreate = fmt.Sprintf("alias-infoblox-test.%s", aliasRecordDomainName)
	aliasRecordHostnameUpdate = fmt.Sprintf("alias-infoblox-test-update.%s", aliasRecordDomainName)
)

func TestAccInfobloxAliasRecordBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: composeConfig(testAccProviderBaseConfig, testAccCheckInfobloxARecordCreate(), testAccCheckInfobloxAliasRecordCreate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxAliasRecordExists("infoblox_a_record.new"),
					testAccCheckInfobloxAliasRecordExists("infoblox_alias_record.new"),
					resource.TestCheckResourceAttr("infoblox_alias_record.new", "name", aliasRecordHostnameCreate),
					resource.TestCheckResourceAttr("infoblox_alias_record.new", "target_name", aRecordHostnameCreate),
					resource.TestCheckResourceAttr("infoblox_alias_record.new", "target_type", "A"),
					resource.TestCheckResourceAttr("infoblox_alias_record.new", "comment", "test alias record"),
					resource.TestCheckResourceAttr("infoblox_alias_record.new", "disable", "true"),
					resource.TestCheckResourceAttr("infoblox_alias_record.new", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_alias_record.new", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_alias_record.new", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
					testAccCheckInfobloxAliasRecordExists("data.infoblox_alias_record.hostname"),
					testAccCheckInfobloxAliasRecordExists("data.infoblox_alias_record.ref"),
					resource.TestCheckResourceAttr("data.infoblox_alias_record.hostname", "name", aliasRecordHostnameCreate),
					resource.TestCheckResourceAttr("data.infoblox_alias_record.hostname", "target_name", aRecordHostnameCreate),
					resource.TestCheckResourceAttr("data.infoblox_alias_record.hostname", "target_type", "A"),
					resource.TestCheckResourceAttr("data.infoblox_alias_record.hostname", "comment", "test alias record"),
					resource.TestCheckResourceAttr("data.infoblox_alias_record.hostname", "disable", "true"),
					resource.TestCheckResourceAttr("data.infoblox_alias_record.hostname", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("data.infoblox_alias_record.hostname", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("data.infoblox_alias_record.hostname", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
				),
			},
			{
				Config: composeConfig(testAccProviderBaseConfig, testAccCheckInfobloxARecordCreate(), testAccCheckInfobloxAliasRecordUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxAliasRecordExists("infoblox_a_record.new"),
					testAccCheckInfobloxAliasRecordExists("infoblox_alias_record.new"),
					resource.TestCheckResourceAttr("infoblox_alias_record.new", "name", aliasRecordHostnameUpdate),
					resource.TestCheckResourceAttr("infoblox_alias_record.new", "target_name", aRecordHostnameCreate),
					resource.TestCheckResourceAttr("infoblox_alias_record.new", "target_type", "A"),
					resource.TestCheckResourceAttr("infoblox_alias_record.new", "comment", "test alias record update"),
					resource.TestCheckResourceAttr("infoblox_alias_record.new", "disable", "false"),
					resource.TestCheckResourceAttr("infoblox_alias_record.new", "extensible_attributes.Location", "{\"value\":\"CollegeStation2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_alias_record.new", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_alias_record.new", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
				),
			},
		},
	})
}

func testAccCheckInfobloxAliasRecordExists(resourceName string) resource.TestCheckFunc {
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

var testAccCheckInfobloxAliasRecordCreate = fmt.Sprintf(`
resource "infoblox_alias_record" "new"{
	name        = "%s"
  target_name = infoblox_a_record.new.hostname
  target_type = "A"
  disable = true
  comment    = "test alias record"
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
data "infoblox_alias_record" "hostname" {
  name = infoblox_alias_record.new.name
}
data "infoblox_alias_record" "ref" {
  ref = infoblox_alias_record.new.ref
}
`, aliasRecordHostnameCreate)

var testAccCheckInfobloxAliasRecordUpdate = fmt.Sprintf(`
resource "infoblox_alias_record" "new"{
	name        = "%s"
  target_name = infoblox_a_record.new.hostname
  target_type = "A"
  disable = false
  comment    = "test alias record update"
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
`, aliasRecordHostnameUpdate)
