package infoblox

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var aRecordDomainName = os.Getenv("INFOBLOX_DOMAIN")

func TestAccInfobloxARecordBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testProviderARecordCreate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxARecordExists("infoblox_a_record.new"),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "ip_address", "172.19.4.6"),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "comment", "test a record"),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "disable", "true"),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "hostname", fmt.Sprintf("infoblox-test.%s", aRecordDomainName)),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\"}"),
				),
			},
			{
				Config: testProviderARecordUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxARecordExists("infoblox_a_record.new"),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "ip_address", "172.19.4.7"),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "comment", "test a record update"),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "disable", "false"),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "hostname", fmt.Sprintf("infoblox-test2.%s", aRecordDomainName)),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "extensible_attributes.Location", "{\"value\":\"CollegeStation2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins2\",\"type\":\"STRING\"}"),
				),
			},
		},
	})
}

func testAccCheckInfobloxARecordExists(resourceName string) resource.TestCheckFunc {
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

var testProviderARecordCreate = fmt.Sprintf(`
resource "infoblox_a_record" "new"{
	ip_address = "172.19.4.6"
	comment    = "test a record"
	hostname   = "infoblox-test.%s"
	disable    = true
	extensible_attributes = {
	  Location = jsonencode({
		value = "CollegeStation",
		type  = "STRING"
	  })
	  Owner = jsonencode({
		value = "leroyjenkins",
		type  = "STRING"
	  })
	}
}
`, aRecordDomainName)

var testProviderARecordUpdate = fmt.Sprintf(`
resource "infoblox_a_record" "new"{
	ip_address = "172.19.4.7"
	comment    = "test a record update"
	hostname   = "infoblox-test2.%s"
	disable    = false
	extensible_attributes = {
	  Location = jsonencode({
		value = "CollegeStation2",
		type  = "STRING"
	  })
	  Owner = jsonencode({
		value = "leroyjenkins2",
		type  = "STRING"
	  })
	}
}
`, aRecordDomainName)
