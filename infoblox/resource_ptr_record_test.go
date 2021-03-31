package infoblox

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccInfobloxPtrRecordBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: composeConfig(testAccProviderBaseConfig, testAccCheckInfobloxARecordCreate(), testAccCheckInfobloxPtrRecordCreate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxPtrRecordExists("infoblox_a_record.new"),
					testAccCheckInfobloxPtrRecordExists("infoblox_ptr_record.new"),
					resource.TestCheckResourceAttr("infoblox_ptr_record.new", "pointer_domain_name", aRecordHostnameCreate),
					resource.TestCheckResourceAttr("infoblox_ptr_record.new", "ip_v4_address", aRecordIPAddress),
					resource.TestCheckResourceAttr("infoblox_ptr_record.new", "comment", "test ptr record"),
					resource.TestCheckResourceAttr("infoblox_ptr_record.new", "disable", "true"),
					resource.TestCheckResourceAttr("infoblox_ptr_record.new", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_ptr_record.new", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_ptr_record.new", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
				),
			},
			{
				Config: composeConfig(testAccProviderBaseConfig, testAccCheckInfobloxARecordUpdate(), testAccCheckInfobloxPtrRecordUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxPtrRecordExists("infoblox_a_record.new"),
					testAccCheckInfobloxPtrRecordExists("infoblox_ptr_record.new"),
					resource.TestCheckResourceAttr("infoblox_ptr_record.new", "pointer_domain_name", aRecordHostnameUpdate),
					resource.TestCheckResourceAttr("infoblox_ptr_record.new", "ip_v4_address", aRecordIPAddress),
					resource.TestCheckResourceAttr("infoblox_ptr_record.new", "comment", "test ptr record update"),
					resource.TestCheckResourceAttr("infoblox_ptr_record.new", "disable", "false"),
					resource.TestCheckResourceAttr("infoblox_ptr_record.new", "extensible_attributes.Location", "{\"value\":\"CollegeStation2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_ptr_record.new", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_ptr_record.new", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
				),
			},
		},
	})
}

func testAccCheckInfobloxPtrRecordExists(resourceName string) resource.TestCheckFunc {
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

var testAccCheckInfobloxPtrRecordCreate = ` 
resource "infoblox_ptr_record" "new" {
  pointer_domain_name = infoblox_a_record.new.hostname
  ip_v4_address       = infoblox_a_record.new.ip_address
  comment    = "test ptr record"
	disable    = true
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
`

var testAccCheckInfobloxPtrRecordUpdate = `
resource "infoblox_ptr_record" "new" {
  pointer_domain_name = infoblox_a_record.new.hostname
  ip_v4_address       = infoblox_a_record.new.ip_address
  comment    = "test ptr record update"
	disable    = false
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
`
