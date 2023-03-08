package infoblox

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	containerIPAddress = os.Getenv("INFOBLOX_CONTAINER_NETWORK")
)

func TestAccInfobloxContainerBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: composeConfig(testAccProviderBaseConfig, testAccCheckInfobloxContainerCreate()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxContainerExists("infoblox_container.new"),
					resource.TestCheckResourceAttr("infoblox_container.new", "cidr", containerIPAddress),
					resource.TestCheckResourceAttr("infoblox_container.new", "comment", "test container"),
					resource.TestCheckResourceAttr("infoblox_container.new", "network_view", "default"),
					resource.TestCheckResourceAttr("infoblox_container.new", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_container.new", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_container.new", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
					testAccCheckInfobloxContainerExists("data.infoblox_container.ref"),
					resource.TestCheckResourceAttr("data.infoblox_container.ref", "cidr", containerIPAddress),
					resource.TestCheckResourceAttr("data.infoblox_container.ref", "comment", "test container"),
					resource.TestCheckResourceAttr("data.infoblox_container.ref", "network_view", "default"),
					resource.TestCheckResourceAttr("data.infoblox_container.ref", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("data.infoblox_container.ref", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("data.infoblox_container.ref", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
				),
			},
			{
				Config: composeConfig(testAccProviderBaseConfig, testAccCheckInfobloxContainerUpdate()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxContainerExists("infoblox_container.new"),
					resource.TestCheckResourceAttr("infoblox_container.new", "cidr", containerIPAddress),
					resource.TestCheckResourceAttr("infoblox_container.new", "comment", "test container update"),
					resource.TestCheckResourceAttr("infoblox_container.new", "network_view", "default"),
					resource.TestCheckResourceAttr("infoblox_container.new", "extensible_attributes.Location", "{\"value\":\"CollegeStation2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_container.new", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_container.new", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
				),
			},
		},
	})
}

func testAccCheckInfobloxContainerExists(resourceName string) resource.TestCheckFunc {
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

func testAccCheckInfobloxContainerCreate() string {
	return fmt.Sprintf(`
  resource "infoblox_container" "new"{
    cidr = "%s"
    comment    = "test container"
	network_view = "default"
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
  data "infoblox_container" "ref" {
    ref = infoblox_container.new.ref
  }
`, containerIPAddress)
}

func testAccCheckInfobloxContainerUpdate() string {

	return fmt.Sprintf(`
resource "infoblox_container" "new"{
	cidr = "%s"
	comment    = "test container update"
	network_view = "default"
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
`, containerIPAddress)
}
