package infoblox

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/techBeck03/go-ipmath"
)

var (
	fixedAddressIPAddressStatic  string
	fixedAddressIPAddressNetwork string
	fixedAddressIPAddressRange   string
)

func TestAccInfobloxFixedAddressBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: composeConfig(testAccProviderBaseConfig, testAccCheckInfobloxNetworkCreate(), testAccCheckInfobloxFixedAddressCreateStatic(), testAccCheckInfobloxFixedAddressCreateFromNetwork(), testAccCheckInfobloxFixedAddressCreateFromRange(), testAccCheckInfobloxRangeCreateStatic()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxFixedAddressExists("infoblox_network.new"),
					testAccCheckInfobloxFixedAddressExists("infoblox_fixed_address.static"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "ip_address", fixedAddressIPAddressStatic),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "hostname", "fixedAddress-test"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "comment", "test fixed address"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "match_client", "RESERVED"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "disable", "true"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.network", "ip_address", fixedAddressIPAddressNetwork),
					resource.TestCheckResourceAttr("infoblox_fixed_address.network", "hostname", "fixedAddress-test"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.network", "comment", "test fixed address"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.network", "match_client", "RESERVED"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.network", "disable", "true"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.network", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.network", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.network", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.range", "ip_address", fixedAddressIPAddressRange),
					resource.TestCheckResourceAttr("infoblox_fixed_address.range", "hostname", "fixedAddress-test"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.range", "comment", "test fixed address"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.range", "match_client", "RESERVED"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.range", "disable", "true"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.range", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.range", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.range", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
				),
			},
			{
				Config: composeConfig(testAccProviderBaseConfig, testAccCheckInfobloxNetworkCreate(), testAccCheckInfobloxFixedAddressUpdateStatic(), testAccCheckInfobloxFixedAddressUpdateFromNetwork(), testAccCheckInfobloxFixedAddressUpdateFromRange(), testAccCheckInfobloxRangeCreateStatic()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxFixedAddressExists("infoblox_network.new"),
					testAccCheckInfobloxFixedAddressExists("infoblox_fixed_address.static"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "ip_address", fixedAddressIPAddressStatic),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "hostname", "fixedAddress-test-update"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "comment", "test fixed address update"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "match_client", "RESERVED"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "disable", "false"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "extensible_attributes.Location", "{\"value\":\"CollegeStation2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.static", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.network", "ip_address", fixedAddressIPAddressNetwork),
					resource.TestCheckResourceAttr("infoblox_fixed_address.network", "hostname", "fixedAddress-test-update"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.network", "comment", "test fixed address update"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.network", "match_client", "RESERVED"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.network", "disable", "false"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.network", "extensible_attributes.Location", "{\"value\":\"CollegeStation2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.network", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.network", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.range", "ip_address", fixedAddressIPAddressRange),
					resource.TestCheckResourceAttr("infoblox_fixed_address.range", "hostname", "fixedAddress-test-update"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.range", "comment", "test fixed address update"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.range", "match_client", "RESERVED"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.range", "disable", "false"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.range", "extensible_attributes.Location", "{\"value\":\"CollegeStation2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.range", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.range", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
				),
			},
		},
	})
}

func testAccCheckInfobloxFixedAddressExists(resourceName string) resource.TestCheckFunc {
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

func testAccCheckInfobloxFixedAddressCreateStatic() string {
	networkIPAddress, _ := ipmath.NewIP(os.Getenv("INFOBLOX_TEST_NETWORK"))
	networkIPAddress.Add(10)
	fixedAddressIPAddressStatic = networkIPAddress.ToIPString()
	return fmt.Sprintf(`
resource "infoblox_fixed_address" "static" {
	ip_address        = "%s"
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
`, fixedAddressIPAddressStatic)
}

func testAccCheckInfobloxFixedAddressUpdateStatic() string {
	networkIPAddress, _ := ipmath.NewIP(os.Getenv("INFOBLOX_TEST_NETWORK"))
	networkIPAddress.Add(11)
	fixedAddressIPAddressStatic = networkIPAddress.ToIPString()
	return fmt.Sprintf(`
resource "infoblox_fixed_address" "static" {
	ip_address        = "%s"
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
`, fixedAddressIPAddressStatic)
}

func testAccCheckInfobloxFixedAddressCreateFromNetwork() string {
	networkIPAddress, _ := ipmath.NewIP(os.Getenv("INFOBLOX_TEST_NETWORK"))
	networkIPAddress.Add(1)
	fixedAddressIPAddressNetwork = networkIPAddress.ToIPString()
	return `
resource "infoblox_fixed_address" "network" {
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
}

func testAccCheckInfobloxFixedAddressUpdateFromNetwork() string {
	networkIPAddress, _ := ipmath.NewIP(os.Getenv("INFOBLOX_TEST_NETWORK"))
	networkIPAddress.Add(1)
	fixedAddressIPAddressNetwork = networkIPAddress.ToIPString()
	return `
resource "infoblox_fixed_address" "network" {
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
}

func testAccCheckInfobloxFixedAddressCreateFromRange() string {
	networkIPAddress, _ := ipmath.NewIP(os.Getenv("INFOBLOX_TEST_NETWORK"))
	networkIPAddress.Add(100)
	fixedAddressIPAddressRange = networkIPAddress.ToIPString()
	return `
resource "infoblox_fixed_address" "range" {
  hostname              = "fixedAddress-test"
  range_function_string = infoblox_range.static.range_function_string
	comment               = "test fixed address"
  disable               = true
	match_client          = "RESERVED"
	restart_if_needed     = true
	grid_ref              = data.infoblox_grid.grid.ref
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
}

func testAccCheckInfobloxFixedAddressUpdateFromRange() string {
	networkIPAddress, _ := ipmath.NewIP(os.Getenv("INFOBLOX_TEST_NETWORK"))
	networkIPAddress.Add(100)
	fixedAddressIPAddressRange = networkIPAddress.ToIPString()
	return `
resource "infoblox_fixed_address" "range" {
  hostname              = "fixedAddress-test-update"
  range_function_string = infoblox_range.static.range_function_string
	comment               = "test fixed address update"
  disable               = false
	match_client          = "RESERVED"
	restart_if_needed     = true
	grid_ref              = data.infoblox_grid.grid.ref
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
}
