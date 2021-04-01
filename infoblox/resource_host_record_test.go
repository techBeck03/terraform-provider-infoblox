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
	hostRecordDomainName            = os.Getenv("INFOBLOX_DOMAIN")
	hostRecordHostnameCreateStatic  = fmt.Sprintf("infoblox-test-host-static.%s", hostRecordDomainName)
	hostRecordHostnameUpdateStatic  = fmt.Sprintf("infoblox-test-host-static-update.%s", hostRecordDomainName)
	hostRecordHostnameCreateNetwork = fmt.Sprintf("infoblox-test-host-network.%s", hostRecordDomainName)
	hostRecordHostnameUpdateNetwork = fmt.Sprintf("infoblox-test-host-network-update.%s", hostRecordDomainName)
	hostRecordHostnameCreateRange   = fmt.Sprintf("infoblox-test-host-range.%s", hostRecordDomainName)
	hostRecordHostnameUpdateRange   = fmt.Sprintf("infoblox-test-host-range-update.%s", hostRecordDomainName)
	hostRecordIPAddressStatic       string
	hostRecordIPAddressNetwork      string
	hostRecordIPAddressRange        string
)

func TestAccInfobloxHostRecordBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: composeConfig(testAccProviderBaseConfig, testAccCheckInfobloxNetworkCreate(), testAccCheckInfobloxHostRecordCreateStatic(), testAccCheckInfobloxHostRecordCreateFromNetwork(), testAccCheckInfobloxHostRecordCreateFromRange(), testAccCheckInfobloxRangeCreateStatic()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxHostRecordExists("infoblox_network.new"),
					testAccCheckInfobloxHostRecordExists("infoblox_host_record.static"),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "hostname", hostRecordHostnameCreateStatic),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "ip_v4_address.0.ip_address", hostRecordIPAddressStatic),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "comment", "test host record"),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "enable_dns", "true"),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
					testAccCheckInfobloxHostRecordExists("infoblox_host_record.network"),
					resource.TestCheckResourceAttr("infoblox_host_record.network", "hostname", hostRecordHostnameCreateNetwork),
					resource.TestCheckResourceAttr("infoblox_host_record.network", "ip_v4_address.0.ip_address", hostRecordIPAddressNetwork),
					resource.TestCheckResourceAttr("infoblox_host_record.network", "comment", "test host record"),
					resource.TestCheckResourceAttr("infoblox_host_record.network", "enable_dns", "true"),
					resource.TestCheckResourceAttr("infoblox_host_record.network", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
					resource.TestCheckResourceAttr("infoblox_host_record.network", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
					resource.TestCheckResourceAttr("infoblox_host_record.network", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
					testAccCheckInfobloxHostRecordExists("infoblox_range.static"),
					testAccCheckInfobloxHostRecordExists("infoblox_host_record.range"),
					resource.TestCheckResourceAttr("infoblox_host_record.range", "hostname", hostRecordHostnameCreateRange),
					resource.TestCheckResourceAttr("infoblox_host_record.range", "ip_v4_address.0.ip_address", hostRecordIPAddressRange),
					resource.TestCheckResourceAttr("infoblox_host_record.range", "comment", "test host record"),
					resource.TestCheckResourceAttr("infoblox_host_record.range", "enable_dns", "true"),
					resource.TestCheckResourceAttr("infoblox_host_record.range", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
					resource.TestCheckResourceAttr("infoblox_host_record.range", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
					resource.TestCheckResourceAttr("infoblox_host_record.range", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
					testAccCheckInfobloxHostRecordExists("data.infoblox_host_record.ref"),
					testAccCheckInfobloxHostRecordExists("data.infoblox_host_record.hostname"),
					resource.TestCheckResourceAttr("data.infoblox_host_record.hostname", "hostname", hostRecordHostnameCreateStatic),
					resource.TestCheckResourceAttr("data.infoblox_host_record.hostname", "ip_v4_address.0.ip_address", hostRecordIPAddressStatic),
					resource.TestCheckResourceAttr("data.infoblox_host_record.hostname", "comment", "test host record"),
					resource.TestCheckResourceAttr("data.infoblox_host_record.hostname", "enable_dns", "true"),
					resource.TestCheckResourceAttr("data.infoblox_host_record.hostname", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
					resource.TestCheckResourceAttr("data.infoblox_host_record.hostname", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\",\"inheritance_source\":{\"_ref\":\"network/ZG5zLm5ldHdvcmskMTcyLjE5LjQuMC8yNC8w:172.19.4.0/24/default\"}}"),
					resource.TestCheckResourceAttr("data.infoblox_host_record.hostname", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
				),
			},
			{
				Config: composeConfig(testAccProviderBaseConfig, testAccCheckInfobloxNetworkCreate(), testAccCheckInfobloxHostRecordUpdateStatic(), testAccCheckInfobloxHostRecordUpdateFromNetwork(), testAccCheckInfobloxHostRecordUpdateFromRange(), testAccCheckInfobloxRangeCreateStatic()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxHostRecordExists("infoblox_network.new"),
					testAccCheckInfobloxHostRecordExists("infoblox_host_record.static"),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "hostname", hostRecordHostnameUpdateStatic),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "ip_v4_address.0.ip_address", hostRecordIPAddressStatic),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "comment", "test host record update"),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "enable_dns", "true"),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "extensible_attributes.Location", "{\"value\":\"CollegeStation2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_host_record.static", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
					testAccCheckInfobloxHostRecordExists("infoblox_host_record.network"),
					resource.TestCheckResourceAttr("infoblox_host_record.network", "hostname", hostRecordHostnameUpdateNetwork),
					resource.TestCheckResourceAttr("infoblox_host_record.network", "ip_v4_address.0.ip_address", hostRecordIPAddressNetwork),
					resource.TestCheckResourceAttr("infoblox_host_record.network", "comment", "test host record update"),
					resource.TestCheckResourceAttr("infoblox_host_record.network", "enable_dns", "true"),
					resource.TestCheckResourceAttr("infoblox_host_record.network", "extensible_attributes.Location", "{\"value\":\"CollegeStation2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_host_record.network", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_host_record.network", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
					testAccCheckInfobloxHostRecordExists("infoblox_host_record.range"),
					resource.TestCheckResourceAttr("infoblox_host_record.range", "hostname", hostRecordHostnameUpdateRange),
					resource.TestCheckResourceAttr("infoblox_host_record.range", "ip_v4_address.0.ip_address", hostRecordIPAddressRange),
					resource.TestCheckResourceAttr("infoblox_host_record.range", "comment", "test host record update"),
					resource.TestCheckResourceAttr("infoblox_host_record.range", "enable_dns", "true"),
					resource.TestCheckResourceAttr("infoblox_host_record.range", "extensible_attributes.Location", "{\"value\":\"CollegeStation2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_host_record.range", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_host_record.range", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
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
        type  = "STRING",
      })
      Location = jsonencode({
        value = "CollegeStation",
        type  = "STRING"
      })
    }
  }
  data "infoblox_host_record" "hostname" {
    hostname = infoblox_host_record.static.hostname
  }
  data "infoblox_host_record" "ref" {
    ref = infoblox_host_record.static.ref
  }
  `, hostRecordHostnameCreateStatic, hostRecordIPAddressStatic)
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
  `, hostRecordHostnameUpdateStatic, hostRecordIPAddressStatic)
}

func testAccCheckInfobloxHostRecordCreateFromNetwork() string {
	networkIPAddress, _ := ipmath.NewIP(os.Getenv("INFOBLOX_TEST_NETWORK"))
	networkIPAddress.Add(1)
	hostRecordIPAddressNetwork = networkIPAddress.ToIPString()
	return fmt.Sprintf(`
  resource "infoblox_host_record" "network" {
    depends_on = [ infoblox_network.new ]
    hostname   = "%s"
    comment    = "test host record"
    enable_dns = true
    ip_v4_address {
      network    = infoblox_network.new.cidr
      use_for_ea_inheritance = true
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
    }
  }
`, hostRecordHostnameCreateNetwork)
}

func testAccCheckInfobloxHostRecordUpdateFromNetwork() string {
	networkIPAddress, _ := ipmath.NewIP(os.Getenv("INFOBLOX_TEST_NETWORK"))
	networkIPAddress.Add(1)
	hostRecordIPAddressNetwork = networkIPAddress.ToIPString()
	return fmt.Sprintf(`
  resource "infoblox_host_record" "network" {
    depends_on = [ infoblox_network.new ]
    hostname   = "%s"
    comment    = "test host record update"
    enable_dns = true
    ip_v4_address {
      network    = infoblox_network.new.cidr
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
`, hostRecordHostnameUpdateNetwork)
}

func testAccCheckInfobloxHostRecordCreateFromRange() string {
	networkIPAddress, _ := ipmath.NewIP(os.Getenv("INFOBLOX_TEST_NETWORK"))
	networkIPAddress.Add(100)
	hostRecordIPAddressRange = networkIPAddress.ToIPString()
	return fmt.Sprintf(`
	  resource "infoblox_host_record" "range" {
	    depends_on = [ infoblox_network.new ]
	    hostname   = "%s"
	    comment    = "test host record"
	    enable_dns = true
	    ip_v4_address {
	      range_function_string  = infoblox_range.static.range_function_string
	      use_for_ea_inheritance = true
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
	    }
	  }
	`, hostRecordHostnameCreateRange)
}

func testAccCheckInfobloxHostRecordUpdateFromRange() string {
	networkIPAddress, _ := ipmath.NewIP(os.Getenv("INFOBLOX_TEST_NETWORK"))
	networkIPAddress.Add(100)
	hostRecordIPAddressRange = networkIPAddress.ToIPString()
	return fmt.Sprintf(`
	  resource "infoblox_host_record" "range" {
	    depends_on = [ infoblox_network.new ]
	    hostname   = "%s"
	    comment    = "test host record update"
	    enable_dns = true
	    ip_v4_address {
	      range_function_string  = infoblox_range.static.range_function_string
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
	`, hostRecordHostnameUpdateRange)
}
