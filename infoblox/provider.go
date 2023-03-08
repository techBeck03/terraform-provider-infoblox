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
				Description: "Infoblox server hostname",
			},
			"port": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("INFOBLOX_PORT", 443),
				Description: "Infoblox server port",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("INFOBLOX_USERNAME", nil),
				Description: "Infoblox server username",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("INFOBLOX_PASSWORD", nil),
				Description: "Infoblox server password",
			},
			"wapi_version": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("INFOBLOX_VERSION", "2.11"),
				Description: "Infoblox server wapi version",
			},
			"disable_tls_verification": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("INFOBLOX_DISABLE_TLS", false),
				Description: "Disable tls verification",
			},
			"orchestrator_extensible_attributes": {
				Type:             schema.TypeMap,
				Description:      "Extensible attributes applied to all objects configured by provider",
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validateEa,
				DiffSuppressFunc: eaSuppressDiff,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"infoblox_container":     resourceContainer(),
			"infoblox_host_record":   resourceHostRecord(),
			"infoblox_network":       resourceNetwork(),
			"infoblox_range":         resourceRange(),
			"infoblox_fixed_address": resourceFixedAddress(),
			"infoblox_a_record":      resourceARecord(),
			"infoblox_cname_record":  resourceCNameRecord(),
			"infoblox_alias_record":  resourceAliasRecord(),
			"infoblox_ptr_record":    resourcePtrRecord(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"infoblox_container":                dataSourceContainer(),
			"infoblox_host_record":              dataSourceHostRecord(),
			"infoblox_network":                  dataSourceNetwork(),
			"infoblox_grid":                     dataSourceGrid(),
			"infoblox_grid_member":              dataSourceGridMember(),
			"infoblox_sequential_address_block": dataSourceSequentialAddressBlock(),
			"infoblox_range":                    dataSourceRange(),
			"infoblox_a_record":                 dataSourceARecord(),
			"infoblox_cname_record":             dataSourceCNameRecord(),
			"infoblox_alias_record":             dataSourceAliasRecord(),
			"infoblox_ptr_record":               dataSourcePtrRecord(),
			"infoblox_fixed_address":            dataSourceFixedAddress(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	hostname := d.Get("hostname").(string)
	port := d.Get("port").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	wapiVersion := d.Get("wapi_version").(string)
	disableTLS := d.Get("disable_tls_verification").(bool)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	config := infoblox.Config{
		Host:                   hostname,
		Port:                   port,
		Username:               username,
		Password:               password,
		Version:                wapiVersion,
		DisableTLSVerification: disableTLS,
	}

	// Check for required provider parameters
	check := validate(config)

	if check.HasError() {
		return nil, check
	}

	client := infoblox.New(config)

	eaMap := d.Get("orchestrator_extensible_attributes").(map[string]interface{})
	if len(eaMap) > 0 {
		eas, err := createExtensibleAttributesFromJSON(eaMap)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		client.OrchestratorEAs = &eas
	}

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
