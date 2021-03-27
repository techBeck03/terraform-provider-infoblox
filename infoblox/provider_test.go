package infoblox

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	var rc terraform.ResourceConfig
	rc.Config = map[string]interface{}{
		"orchestrator_extensible_attributes": map[string]string{
			"Orchestrator": "{\"value\":\"Terraform\",\"type\":\"ENUM\"}",
		},
	}
	testAccProvider.Configure(nil, &rc)
	testAccProviders = map[string]*schema.Provider{
		"infoblox": testAccProvider,
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
}
