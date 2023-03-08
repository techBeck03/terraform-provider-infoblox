package infoblox

import (
	"net"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]func() (*schema.Provider, error)
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]func() (*schema.Provider, error){
		"infoblox": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if err := os.Getenv("INFOBLOX_HOSTNAME"); err == "" {
		t.Fatal("INFOBLOX_HOSTNAME must be set for acceptance tests")
	}
	if err := os.Getenv("INFOBLOX_PORT"); err == "" {
		t.Fatal("INFOBLOX_PORT must be set for acceptance tests")
	}
	if err := os.Getenv("INFOBLOX_VERSION"); err == "" {
		t.Fatal("INFOBLOX_VERSION must be set for acceptance tests")
	}
	if err := os.Getenv("INFOBLOX_USERNAME"); err == "" {
		t.Fatal("INFOBLOX_USERNAME must be set for acceptance tests")
	}
	if err := os.Getenv("INFOBLOX_PASSWORD"); err == "" {
		t.Fatal("INFOBLOX_PASSWORD must be set for acceptance tests")
	}
	if err := os.Getenv("INFOBLOX_DISABLE_TLS"); err == "" {
		t.Fatal("INFOBLOX_DISABLE_TLS must be set for acceptance tests")
	}
	testNetwork := os.Getenv("INFOBLOX_TEST_NETWORK")
	if testNetwork == "" {
		t.Fatal("INFOBLOX_TEST_NETWORK must be set for acceptance tests")
	} else {
		_, net, err := net.ParseCIDR(testNetwork)
		if err != nil {
			t.Fatal(err)
		}
		prefix, _ := net.Mask.Size()
		if prefix != 24 {
			t.Fatal("INFOBLOX_TEST_NETWORK must be a properly formatted /24 CIDR network string")
		} else if testNetwork != net.String() {
			t.Fatalf("INFOBLOX_TEST_NETWORK expected to be: %s but found: %s", net.String(), testNetwork)
		}
	}
}

var testAccProviderBaseConfig = `
  provider infoblox {
    orchestrator_extensible_attributes = {
      Orchestrator = jsonencode({
        value = "Terraform",
        type  = "ENUM"
      })
    }
  }
`
