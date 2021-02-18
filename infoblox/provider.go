package infoblox

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("INFOBLOX_HOSTNAME", nil),
			},
			"port": {
				Type:        schema.TypeString,
				Default:     "443",
				DefaultFunc: schema.EnvDefaultFunc("INFOBLOX_PORT", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("INFOBLOX_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("INFOBLOX_PASSWORD", nil),
			},
			"version": {
				Type:        schema.TypeString,
				Default:     "2.11",
				DefaultFunc: schema.EnvDefaultFunc("INFOBLOX_VERSION", nil),
			},
			"disable_tls_verification": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("INFOBLOX_DISABLE_TLS", false),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"host_record": resourceHostRecord(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"host_record": dataSourceHostRecord(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	hostname := d.Get("hostname").(string)
	port := d.Get("port").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	version := d.Get("version").(string)
	disableTLS := d.Get("disable_tls_verification").(bool)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	config := infoblox.Config{
		Host:                   hostname,
		Port:                   port,
		Username:               username,
		Password:               password,
		Version:                version,
		DisableTLSVerification: disableTLS,
	}

	// Check for required provider parameters
	check := validate(config)

	if check.HasError() {
		return nil, check
	}

	client := infoblox.New(config)

	return &client, diags
}

// validate validates the config needed to initialize a infoblox client,
// returning a single error with all validation errors, or nil if no error.
func validate(config infoblox.Config) diag.Diagnostics {
	var diags diag.Diagnostics

	if config.Host == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Missing provider parameter",
			Detail:   "Hostname must be configured for the infoblox provider",
		})
	}
	if config.Username == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Missing provider parameter",
			Detail:   "Username must be configured for the infoblox provider",
		})
	}
	if config.Password == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Missing provider parameter",
			Detail:   "Password must be configured for the infoblox provider",
		})
	}
	return diags
}
