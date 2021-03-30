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
	hostRecordDomainName       = os.Getenv("INFOBLOX_DOMAIN")
	hostRecordHostnameCreate   = fmt.Sprintf("infoblox-test-host.%s", hostRecordDomainName)
	hostRecordHostnameUpdate   = fmt.Sprintf("infoblox-test-host-update.%s", hostRecordDomainName)
	hostRecordIPAddressStatic  string
	hostRecordIPAddressNetwork string
	hostRecordIPAddressRange   string
)

func TestAccInfobloxHostRecordBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Config: composeConfig(testAccProviderBaseConfig, testAccCheckInfobloxNetworkCreate(), testAccCheckInfobloxHostRecordCreateStatic(), testAccCheckInfobloxHostRecordCreateFromNetwork(), testAccCheckInfobloxHostRecordCreateFromRange(), testAccCheckInfobloxRangeCreateStatic()),
				Config: composeConfig(testAccProviderBaseConfig, testAccCheckInfobloxNetworkCreate(), testAccCheckInfobloxHostRecordCreateStatic()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxHostRecordExists("infoblox_network.new"),
					testAccCheckInfobloxHostRecordExists("infoblox_host_record.static"),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "hostname", hostRecordHostnameCreate),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "ip_v4_address.0.ip_address", hostRecordIPAddressStatic),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "comment", "test host record"),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "enable_dns", "true"),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
					// resource.TestCheckResourceAttr("infoblox_host_record.network", "ip_address", hostRecordIPAddressNetwork),
					// resource.TestCheckResourceAttr("infoblox_host_record.network", "hostname", "hostRecord-test"),
					// resource.TestCheckResourceAttr("infoblox_host_record.network", "comment", "test fixed address"),
					// resource.TestCheckResourceAttr("infoblox_host_record.network", "match_client", "RESERVED"),
					// resource.TestCheckResourceAttr("infoblox_host_record.network", "disable", "true"),
					// resource.TestCheckResourceAttr("infoblox_host_record.network", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
					// resource.TestCheckResourceAttr("infoblox_host_record.network", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
					// resource.TestCheckResourceAttr("infoblox_host_record.network", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
					// resource.TestCheckResourceAttr("infoblox_host_record.range", "ip_address", hostRecordIPAddressRange),
					// resource.TestCheckResourceAttr("infoblox_host_record.range", "hostname", "hostRecord-test"),
					// resource.TestCheckResourceAttr("infoblox_host_record.range", "comment", "test fixed address"),
					// resource.TestCheckResourceAttr("infoblox_host_record.range", "match_client", "RESERVED"),
					// resource.TestCheckResourceAttr("infoblox_host_record.range", "disable", "true"),
					// resource.TestCheckResourceAttr("infoblox_host_record.range", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
					// resource.TestCheckResourceAttr("infoblox_host_record.range", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
					// resource.TestCheckResourceAttr("infoblox_host_record.range", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
				),
			},
			{
				// Config: composeConfig(testAccProviderBaseConfig, testAccCheckInfobloxNetworkCreate(), testAccCheckInfobloxHostRecordUpdateStatic(), testAccCheckInfobloxHostRecordUpdateFromNetwork(), testAccCheckInfobloxHostRecordUpdateFromRange(), testAccCheckInfobloxRangeCreateStatic()),
				Config: composeConfig(testAccProviderBaseConfig, testAccCheckInfobloxNetworkCreate(), testAccCheckInfobloxHostRecordUpdateStatic()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxHostRecordExists("infoblox_network.new"),
					testAccCheckInfobloxHostRecordExists("infoblox_host_record.static"),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "hostname", hostRecordHostnameUpdate),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "ip_v4_address.0.ip_address", hostRecordIPAddressStatic),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "comment", "test host record update"),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "enable_dns", "true"),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "extensible_attributes.Location", "{\"value\":\"CollegeStation2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
					// resource.TestCheckResourceAttr("infoblox_host_record.network", "ip_address", hostRecordIPAddressNetwork),
					// resource.TestCheckResourceAttr("infoblox_host_record.network", "hostname", "hostRecord-test-update"),
					// resource.TestCheckResourceAttr("infoblox_host_record.network", "comment", "test fixed address update"),
					// resource.TestCheckResourceAttr("infoblox_host_record.network", "match_client", "RESERVED"),
					// resource.TestCheckResourceAttr("infoblox_host_record.network", "disable", "false"),
					// resource.TestCheckResourceAttr("infoblox_host_record.network", "extensible_attributes.Location", "{\"value\":\"CollegeStation2\",\"type\":\"STRING\"}"),
					// resource.TestCheckResourceAttr("infoblox_host_record.network", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins2\",\"type\":\"STRING\"}"),
					// resource.TestCheckResourceAttr("infoblox_host_record.network", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
					// resource.TestCheckResourceAttr("infoblox_host_record.range", "ip_address", hostRecordIPAddressRange),
					// resource.TestCheckResourceAttr("infoblox_host_record.range", "hostname", "hostRecord-test-update"),
					// resource.TestCheckResourceAttr("infoblox_host_record.range", "comment", "test fixed address update"),
					// resource.TestCheckResourceAttr("infoblox_host_record.range", "match_client", "RESERVED"),
					// resource.TestCheckResourceAttr("infoblox_host_record.range", "disable", "false"),
					// resource.TestCheckResourceAttr("infoblox_host_record.range", "extensible_attributes.Location", "{\"value\":\"CollegeStation2\",\"type\":\"STRING\"}"),
					// resource.TestCheckResourceAttr("infoblox_host_record.range", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins2\",\"type\":\"STRING\"}"),
					// resource.TestCheckResourceAttr("infoblox_host_record.range", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
				),
			},
		},
	})
}

func testAccCheckInfobloxHostRecordExists(resourceName string) resource.TestCheckFunc {
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

func testAccCheckInfobloxHostRecordCreateStatic() string {
	networkIPAddress, _ := ipmath.NewIP(os.Getenv("INFOBLOX_TEST_NETWORK"))
	networkIPAddress.Add(10)
	hostRecordIPAddressStatic = networkIPAddress.ToIPString()
	return fmt.Sprintf(`
  resource "infoblox_host_record" "static" {
    depends_on = [ infoblox_network.new ]
    hostname   = "%s"
    comment    = "test host record"
    enable_dns = true
    ip_v4_address {
      ip_address = "%s"
      use_for_ea_inheritance = true
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
  `, hostRecordHostnameCreate, hostRecordIPAddressStatic)
}

func testAccCheckInfobloxHostRecordUpdateStatic() string {
	networkIPAddress, _ := ipmath.NewIP(os.Getenv("INFOBLOX_TEST_NETWORK"))
	networkIPAddress.Add(11)
	hostRecordIPAddressStatic = networkIPAddress.ToIPString()
	return fmt.Sprintf(`
  resource "infoblox_host_record" "static" {
    depends_on = [ infoblox_network.new ]
    hostname   = "%s"
    comment    = "test host record update"
    enable_dns = true
    ip_v4_address {
      ip_address = "%s"
      use_for_ea_inheritance = true
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
  `, hostRecordHostnameUpdate, hostRecordIPAddressStatic)
}

func testAccCheckInfobloxHostRecordCreateFromNetwork() string {
	networkIPAddress, _ := ipmath.NewIP(os.Getenv("INFOBLOX_TEST_NETWORK"))
	networkIPAddress.Add(1)
	hostRecordIPAddressNetwork = networkIPAddress.ToIPString()
	return `
resource "infoblox_host_record" "network" {
  hostname          = "hostRecord-test"
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

func testAccCheckInfobloxHostRecordUpdateFromNetwork() string {
	networkIPAddress, _ := ipmath.NewIP(os.Getenv("INFOBLOX_TEST_NETWORK"))
	networkIPAddress.Add(1)
	hostRecordIPAddressNetwork = networkIPAddress.ToIPString()
	return `
resource "infoblox_host_record" "network" {
  hostname          = "hostRecord-test-update"
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

func testAccCheckInfobloxHostRecordCreateFromRange() string {
	networkIPAddress, _ := ipmath.NewIP(os.Getenv("INFOBLOX_TEST_NETWORK"))
	networkIPAddress.Add(100)
	hostRecordIPAddressRange = networkIPAddress.ToIPString()
	return `
resource "infoblox_host_record" "range" {
  hostname              = "hostRecord-test"
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

func testAccCheckInfobloxHostRecordUpdateFromRange() string {
	networkIPAddress, _ := ipmath.NewIP(os.Getenv("INFOBLOX_TEST_NETWORK"))
	networkIPAddress.Add(100)
	hostRecordIPAddressRange = networkIPAddress.ToIPString()
	return `
resource "infoblox_host_record" "range" {
  hostname              = "hostRecord-test-update"
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
