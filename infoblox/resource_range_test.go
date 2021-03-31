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
	rangeNetworkAddress         string
	rangeStartAddressSequential string
	rangeEndAddressSequential   string
	rangeStartAddressStatic     string
	rangeEndAddressStatic       string
)

func TestAccInfobloxRangeBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: composeConfig(testAccProviderBaseConfig, testAccCheckInfobloxNetworkCreate(), testAccCheckInfobloxRangeCreateSequential(), testAccCheckInfobloxRangeCreateStatic()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxNetworkExists("infoblox_range.sequential"),
					resource.TestCheckResourceAttr("infoblox_range.sequential", "cidr", rangeNetworkAddress),
					resource.TestCheckResourceAttr("infoblox_range.sequential", "comment", "test range"),
					resource.TestCheckResourceAttr("infoblox_range.sequential", "sequential_count", "10"),
					resource.TestCheckResourceAttr("infoblox_range.sequential", "start_address", rangeStartAddressSequential),
					resource.TestCheckResourceAttr("infoblox_range.sequential", "end_address", rangeEndAddressSequential),
					resource.TestCheckResourceAttr("infoblox_range.sequential", "range_function_string", fmt.Sprintf("%s-%s", rangeStartAddressSequential, rangeEndAddressSequential)),
					resource.TestCheckResourceAttr("infoblox_range.sequential", "disable_dhcp", "true"),
					resource.TestCheckResourceAttr("infoblox_range.sequential", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
					resource.TestCheckResourceAttr("infoblox_range.sequential", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
					resource.TestCheckResourceAttr("infoblox_range.sequential", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
					testAccCheckInfobloxNetworkExists("infoblox_range.static"),
					resource.TestCheckResourceAttr("infoblox_range.static", "cidr", rangeNetworkAddress),
					resource.TestCheckResourceAttr("infoblox_range.static", "comment", "test range"),
					resource.TestCheckResourceAttr("infoblox_range.static", "start_address", rangeStartAddressStatic),
					resource.TestCheckResourceAttr("infoblox_range.static", "end_address", rangeEndAddressStatic),
					resource.TestCheckResourceAttr("infoblox_range.static", "range_function_string", fmt.Sprintf("%s-%s", rangeStartAddressStatic, rangeEndAddressStatic)),
					resource.TestCheckResourceAttr("infoblox_range.static", "disable_dhcp", "true"),
					resource.TestCheckResourceAttr("infoblox_range.static", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
					resource.TestCheckResourceAttr("infoblox_range.static", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
					resource.TestCheckResourceAttr("infoblox_range.static", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
				),
			},
			{
				Config: composeConfig(testAccProviderBaseConfig, testAccCheckInfobloxNetworkCreate(), testAccCheckInfobloxRangeUpdateSequential(), testAccCheckInfobloxRangeUpdateStatic()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxNetworkExists("infoblox_range.sequential"),
					resource.TestCheckResourceAttr("infoblox_range.sequential", "cidr", rangeNetworkAddress),
					resource.TestCheckResourceAttr("infoblox_range.sequential", "comment", "test range update"),
					resource.TestCheckResourceAttr("infoblox_range.sequential", "sequential_count", "12"),
					resource.TestCheckResourceAttr("infoblox_range.sequential", "start_address", rangeStartAddressSequential),
					resource.TestCheckResourceAttr("infoblox_range.sequential", "end_address", rangeEndAddressSequential),
					resource.TestCheckResourceAttr("infoblox_range.sequential", "range_function_string", fmt.Sprintf("%s-%s", rangeStartAddressSequential, rangeEndAddressSequential)),
					resource.TestCheckResourceAttr("infoblox_range.sequential", "disable_dhcp", "false"),
					resource.TestCheckResourceAttr("infoblox_range.sequential", "extensible_attributes.Location", "{\"value\":\"CollegeStation2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_range.sequential", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_range.sequential", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
					testAccCheckInfobloxNetworkExists("infoblox_range.static"),
					resource.TestCheckResourceAttr("infoblox_range.static", "cidr", rangeNetworkAddress),
					resource.TestCheckResourceAttr("infoblox_range.static", "comment", "test range update"),
					resource.TestCheckResourceAttr("infoblox_range.static", "start_address", rangeStartAddressStatic),
					resource.TestCheckResourceAttr("infoblox_range.static", "end_address", rangeEndAddressStatic),
					resource.TestCheckResourceAttr("infoblox_range.static", "range_function_string", fmt.Sprintf("%s-%s", rangeStartAddressStatic, rangeEndAddressStatic)),
					resource.TestCheckResourceAttr("infoblox_range.static", "disable_dhcp", "false"),
					resource.TestCheckResourceAttr("infoblox_range.static", "extensible_attributes.Location", "{\"value\":\"CollegeStation2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_range.static", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_range.static", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
				),
			},
		},
	})
}

func testAccCheckInfobloxRangeExists(resourceName string) resource.TestCheckFunc {
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

func testAccCheckInfobloxRangeCreateSequential() string {
	networkIPAddress, _ := ipmath.NewIP(os.Getenv("INFOBLOX_TEST_NETWORK"))
	rangeNetworkAddress, _ = networkIPAddress.ToCIDRString()
	networkIPAddress.Inc()
	rangeStartAddressSequential = networkIPAddress.ToIPString()
	networkIPAddress.Add(9)
	rangeEndAddressSequential = networkIPAddress.ToIPString()
	return `
resource "infoblox_range" "sequential" {
	cidr             = infoblox_network.new.cidr
	comment          = "test range"
	sequential_count = 10
	disable_dhcp     = true
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

func testAccCheckInfobloxRangeUpdateSequential() string {
	networkIPAddress, _ := ipmath.NewIP(os.Getenv("INFOBLOX_TEST_NETWORK"))
	rangeNetworkAddress, _ = networkIPAddress.ToCIDRString()
	networkIPAddress.Inc()
	rangeStartAddressSequential = networkIPAddress.ToIPString()
	networkIPAddress.Add(11)
	rangeEndAddressSequential = networkIPAddress.ToIPString()
	return `
resource "infoblox_range" "sequential" {
	cidr             = infoblox_network.new.cidr
	comment          = "test range update"
	sequential_count = 12
	disable_dhcp     = false
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

func testAccCheckInfobloxRangeCreateStatic() string {
	networkIPAddress, _ := ipmath.NewIP(os.Getenv("INFOBLOX_TEST_NETWORK"))
	rangeNetworkAddress, _ = networkIPAddress.ToCIDRString()
	networkIPAddress.Add(100)
	rangeStartAddressStatic = networkIPAddress.ToIPString()
	networkIPAddress.Add(9)
	rangeEndAddressStatic = networkIPAddress.ToIPString()
	return fmt.Sprintf(`
resource "infoblox_range" "static" {
	cidr             = infoblox_network.new.cidr
	comment          = "test range"
	start_address    = "%s"
  end_address      = "%s"
	disable_dhcp     = true
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
`, rangeStartAddressStatic, rangeEndAddressStatic)
}

func testAccCheckInfobloxRangeUpdateStatic() string {
	networkIPAddress, _ := ipmath.NewIP(os.Getenv("INFOBLOX_TEST_NETWORK"))
	rangeNetworkAddress, _ = networkIPAddress.ToCIDRString()
	networkIPAddress.Add(98)
	rangeStartAddressStatic = networkIPAddress.ToIPString()
	networkIPAddress.Add(11)
	rangeEndAddressStatic = networkIPAddress.ToIPString()
	return fmt.Sprintf(`
resource "infoblox_range" "static" {
	cidr             = infoblox_network.new.cidr
	comment          = "test range update"
  start_address    = "%s"
  end_address      = "%s"
	disable_dhcp     = false
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
`, rangeStartAddressStatic, rangeEndAddressStatic)
}
