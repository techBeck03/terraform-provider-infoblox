package infoblox

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	networkDomainName  = os.Getenv("INFOBLOX_DOMAIN")
	gridMemberHostname = fmt.Sprintf("infoblox.%s", networkDomainName)
)

func TestAccInfobloxNetworkBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testProviderNetworkCreate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxNetworkExists("infoblox_network.new"),
					resource.TestCheckResourceAttr("infoblox_network.new", "cidr", "172.19.4.0/24"),
					resource.TestCheckResourceAttr("infoblox_network.new", "comment", "example network"),
					resource.TestCheckResourceAttr("infoblox_network.new", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_network.new", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_network.new", "extensible_attributes.Gateway", "{\"value\":\"172.19.4.1\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_network.new", "option.0.name", "routers"),
					resource.TestCheckResourceAttr("infoblox_network.new", "option.0.code", "3"),
					resource.TestCheckResourceAttr("infoblox_network.new", "option.0.value", "172.19.4.1"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.gw", "ip_address", "172.19.4.1"),
					resource.TestCheckResourceAttr("infoblox_fixed_address.gw", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
				),
			},
			{
				Config: testProviderNetworkUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxNetworkExists("infoblox_network.new"),
					resource.TestCheckResourceAttr("infoblox_network.new", "comment", "example network update"),
					resource.TestCheckResourceAttr("infoblox_network.new", "extensible_attributes.Location", "{\"value\":\"CollegeStation2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_network.new", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_network.new", "extensible_attributes.Gateway", "{\"value\":\"172.19.4.1\",\"type\":\"STRING\"}"),
				),
			},
		},
	})
}

func testAccCheckInfobloxNetworkExists(resourceName string) resource.TestCheckFunc {
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

var testProviderNetworkCreate = fmt.Sprintf(`
data "infoblox_grid" "grid" {
	name = "Infoblox"
}

data "infoblox_grid_member" "member" {
	hostname = "%s"
}
resource "infoblox_network" "new" {
	cidr       = "172.19.4.0/24"
	comment    = "example network"
	network_view      = "default"
	restart_if_needed = true
	grid_ref          = data.infoblox_grid.grid.ref
	member {
	  hostname = data.infoblox_grid_member.member.hostname
	}
	option {
	  code  = 3
	  name  = "routers"
	  value = "172.19.4.1"
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
	  Gateway = jsonencode({
		value = "172.19.4.1",
		type  = "STRING"
	  })
	}
}
resource "infoblox_fixed_address" "gw" {
	ip_address = "172.19.4.1"
	cidr              = infoblox_network.new.cidr
	comment           = "Default Gateway"
	match_client      = "RESERVED"
	restart_if_needed = true
	grid_ref          = data.infoblox_grid.grid.ref
	member {
	  hostname = data.infoblox_grid_member.member.hostname
	}
  }
`, gridMemberHostname)

var testProviderNetworkUpdate = fmt.Sprintf(`
data "infoblox_grid" "grid" {
	name = "Infoblox"
}

data "infoblox_grid_member" "member" {
	hostname = "%s"
}
resource "infoblox_network" "new" {
	cidr       = "172.19.4.0/24"
	comment    = "example network update"
	network_view      = "default"
	restart_if_needed = true
	grid_ref          = data.infoblox_grid.grid.ref
	member {
	  hostname = data.infoblox_grid_member.member.hostname
	}
	option {
	  code  = 3
	  name  = "routers"
	  value = "172.19.4.1"
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
	  Gateway = jsonencode({
		value = "172.19.4.1",
		type  = "STRING"
	  })
	}
}
resource "infoblox_fixed_address" "gw" {
	ip_address = "172.19.4.1"
	cidr              = infoblox_network.new.cidr
	comment           = "Default Gateway"
	match_client      = "RESERVED"
	restart_if_needed = true
	grid_ref          = data.infoblox_grid.grid.ref
	member {
	  hostname = data.infoblox_grid_member.member.hostname
	}
}
`, gridMemberHostname)
