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
	networkDomainName     = os.Getenv("INFOBLOX_DOMAIN")
	gridMemberHostname    = fmt.Sprintf("infoblox.%s", networkDomainName)
	networkNetworkAddress string
	networkGatewayAddress string
)

func TestAccInfobloxNetworkBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: composeConfig(testAccProviderBaseConfig, testAccCheckInfobloxNetworkCreate()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxNetworkExists("infoblox_network.new"),
					resource.TestCheckResourceAttr("infoblox_network.new", "cidr", networkNetworkAddress),
					resource.TestCheckResourceAttr("infoblox_network.new", "comment", "test network"),
					resource.TestCheckResourceAttr("infoblox_network.new", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_network.new", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_network.new", "extensible_attributes.Gateway", fmt.Sprintf("{\"value\":\"%s\",\"type\":\"STRING\"}", networkGatewayAddress)),
					resource.TestCheckResourceAttr("infoblox_network.new", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
					resource.TestCheckResourceAttr("infoblox_network.new", "option.0.name", "routers"),
					resource.TestCheckResourceAttr("infoblox_network.new", "option.0.code", "3"),
					resource.TestCheckResourceAttr("infoblox_network.new", "option.0.value", networkGatewayAddress),
				),
			},
			{
				Config: composeConfig(testAccProviderBaseConfig, testAccCheckInfobloxNetworkUpdate()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxNetworkExists("infoblox_network.new"),
					resource.TestCheckResourceAttr("infoblox_network.new", "cidr", networkNetworkAddress),
					resource.TestCheckResourceAttr("infoblox_network.new", "comment", "test network update"),
					resource.TestCheckResourceAttr("infoblox_network.new", "extensible_attributes.Location", "{\"value\":\"CollegeStation2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_network.new", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_network.new", "extensible_attributes.Gateway", fmt.Sprintf("{\"value\":\"%s\",\"type\":\"STRING\"}", networkGatewayAddress)),
					resource.TestCheckResourceAttr("infoblox_network.new", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
					resource.TestCheckResourceAttr("infoblox_network.new", "option.0.name", "routers"),
					resource.TestCheckResourceAttr("infoblox_network.new", "option.0.code", "3"),
					resource.TestCheckResourceAttr("infoblox_network.new", "option.0.value", networkGatewayAddress),
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

func testAccCheckInfobloxNetworkCreate() string {
	networkIPAddress, _ := ipmath.NewIP(os.Getenv("INFOBLOX_TEST_NETWORK"))
	networkNetworkAddress, _ = networkIPAddress.ToCIDRString()
	networkIPAddress.Add(1)
	networkGatewayAddress = networkIPAddress.ToIPString()
	return fmt.Sprintf(`
data "infoblox_grid" "grid" {
	name = "Infoblox"
}

data "infoblox_grid_member" "member" {
	hostname = "%s"
}
resource "infoblox_network" "new" {
	cidr       = "%s"
	comment    = "test network"
	network_view      = "default"
	restart_if_needed = true
	grid_ref          = data.infoblox_grid.grid.ref
	member {
	  hostname = data.infoblox_grid_member.member.hostname
	}
	option {
	  code  = 3
	  name  = "routers"
	  value = "%s"
	}
	extensible_attributes = {
	  Owner = jsonencode({
		value = "leroyjenkins",
		type  = "STRING",
	  })
	  Location = jsonencode({
		value = "CollegeStation",
		type  = "STRING"
	  })
	  Gateway = jsonencode({
		value = "%s",
		type  = "STRING"
	  })
	}
}
`, gridMemberHostname, networkNetworkAddress, networkGatewayAddress, networkGatewayAddress)
}

func testAccCheckInfobloxNetworkUpdate() string {
	networkIPAddress, _ := ipmath.NewIP(os.Getenv("INFOBLOX_TEST_NETWORK"))
	networkNetworkAddress, _ = networkIPAddress.ToCIDRString()
	networkIPAddress.Add(2)
	networkGatewayAddress = networkIPAddress.ToIPString()
	return fmt.Sprintf(`
data "infoblox_grid" "grid" {
	name = "Infoblox"
}

data "infoblox_grid_member" "member" {
	hostname = "%s"
}
resource "infoblox_network" "new" {
	cidr       = "%s"
	comment    = "test network update"
	network_view      = "default"
	restart_if_needed = true
	grid_ref          = data.infoblox_grid.grid.ref
	member {
	  hostname = data.infoblox_grid_member.member.hostname
	}
	option {
	  code  = 3
	  name  = "routers"
	  value = "%s"
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
		value = "%s",
		type  = "STRING"
	  })
	}
}
`, gridMemberHostname, networkNetworkAddress, networkGatewayAddress, networkGatewayAddress)
}
