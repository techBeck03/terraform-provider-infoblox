package infoblox

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccInfobloxfixedAddressRecordBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: composeConfig(testProviderNetworkCreate, testProviderfixedAddressRecordCreateStatic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxfixedAddressRecordExists("infoblox_network.new"),
					testAccCheckInfobloxfixedAddressRecordExists("infoblox_fixed_address.static"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "ip_address", "172.19.4.251"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "hostname", "fixedAddress-test"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "comment", "test fixed address"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "match_client", "RESERVED"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "disable", "true"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
				),
			},
			{
				Config: composeConfig(testProviderNetworkCreate, testProviderfixedAddressRecordUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxfixedAddressRecordExists("infoblox_network.new"),
					testAccCheckInfobloxfixedAddressRecordExists("infoblox_fixed_address.static"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "ip_address", "172.19.4.251"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "hostname", "fixedAddress-test-update"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "comment", "test fixed address update"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "match_client", "RESERVED"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "disable", "false"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "extensible_attributes.Location", "{\"value\":\"CollegeStation2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins2\",\"type\":\"STRING\"}"),
				),
			},
		},
	})
}

func testAccCheckInfobloxfixedAddressRecordExists(resourceName string) resource.TestCheckFunc {
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

var testProviderfixedAddressRecordCreateStatic = `
resource "infoblox_fixed_address" "static" {
	ip_address        = "172.19.4.251"
  hostname          = "fixedAddress-test"
  cidr              = infoblox_network.new.cidr
	comment           = "test fixed address"
  disable           = true
	match_client      = "RESERVED"
	restart_if_needed = true
	grid_ref          = data.infoblox_grid.grid.ref
	member {
	  hostname = data.infoblox_grid_member.member.hostname
	}
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

var testProviderfixedAddressRecordUpdate = `
resource "infoblox_fixed_address" "static" {
	ip_address        = "172.19.4.251"
  hostname          = "fixedAddress-test-update"
  cidr              = infoblox_network.new.cidr
	comment           = "test fixed address update"
  disable           = false
	match_client      = "RESERVED"
	restart_if_needed = true
	grid_ref          = data.infoblox_grid.grid.ref
	member {
	  hostname = data.infoblox_grid_member.member.hostname
	}
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
