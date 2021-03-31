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
	aRecordDomainName     = os.Getenv("INFOBLOX_DOMAIN")
	aRecordIPAddress      string
	aRecordHostnameCreate = fmt.Sprintf("infoblox-test-ptr.%s", aRecordDomainName)
	aRecordHostnameUpdate = fmt.Sprintf("infoblox-test-ptr-update.%s", aRecordDomainName)
)

func TestAccInfobloxARecordBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: composeConfig(testAccProviderBaseConfig, testAccCheckInfobloxARecordCreate()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxARecordExists("infoblox_a_record.new"),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "ip_address", aRecordIPAddress),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "comment", "test a record"),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "disable", "true"),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "hostname", aRecordHostnameCreate),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
					testAccCheckInfobloxARecordExists("data.infoblox_a_record.hostname"),
					testAccCheckInfobloxARecordExists("data.infoblox_a_record.ref"),
					resource.TestCheckResourceAttr("data.infoblox_a_record.hostname", "ip_address", aRecordIPAddress),
					resource.TestCheckResourceAttr("data.infoblox_a_record.hostname", "comment", "test a record"),
					resource.TestCheckResourceAttr("data.infoblox_a_record.hostname", "disable", "true"),
					resource.TestCheckResourceAttr("data.infoblox_a_record.hostname", "hostname", aRecordHostnameCreate),
					resource.TestCheckResourceAttr("data.infoblox_a_record.hostname", "extensible_attributes.Location", "{\"value\":\"CollegeStation\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("data.infoblox_a_record.hostname", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("data.infoblox_a_record.hostname", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
				),
			},
			{
				Config: composeConfig(testAccProviderBaseConfig, testAccCheckInfobloxARecordUpdate()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInfobloxARecordExists("infoblox_a_record.new"),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "ip_address", aRecordIPAddress),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "comment", "test a record update"),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "disable", "false"),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "hostname", aRecordHostnameUpdate),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "extensible_attributes.Location", "{\"value\":\"CollegeStation2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "extensible_attributes.Owner", "{\"value\":\"leroyjenkins2\",\"type\":\"STRING\"}"),
					resource.TestCheckResourceAttr("infoblox_a_record.new", "extensible_attributes.Orchestrator", "{\"value\":\"Terraform\",\"type\":\"ENUM\"}"),
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
			return fmt.Errorf("Resource: %s not set", resourceName)
		}

		return nil
	}
}

func testAccCheckInfobloxARecordCreate() string {
	aRecordNetworkAddress, _ := ipmath.NewIP(os.Getenv("INFOBLOX_TEST_NETWORK"))
	aRecordNetworkAddress.Add(2)
	aRecordIPAddress = aRecordNetworkAddress.ToIPString()
	return fmt.Sprintf(`
  resource "infoblox_a_record" "new"{
    ip_address = "%s"
    comment    = "test a record"
    hostname   = "%s"
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
  data "infoblox_a_record" "hostname" {
    hostname = infoblox_a_record.new.hostname
  }
  data "infoblox_a_record" "ref" {
    ref = infoblox_a_record.new.ref
  }
`, aRecordIPAddress, aRecordHostnameCreate)
}

func testAccCheckInfobloxARecordUpdate() string {
	aRecordNetworkAddress, _ := ipmath.NewIP(os.Getenv("INFOBLOX_TEST_NETWORK"))
	aRecordNetworkAddress.Add(3)
	aRecordIPAddress = aRecordNetworkAddress.ToIPString()
	return fmt.Sprintf(`
resource "infoblox_a_record" "new"{
	ip_address = "%s"
	comment    = "test a record update"
	hostname   = "%s"
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
`, aRecordIPAddress, aRecordHostnameUpdate)
}
