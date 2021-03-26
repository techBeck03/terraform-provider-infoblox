package infoblox

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

var (
	dataHostRecordRequiredSearchFields = []string{
		"hostname",
		"ref",
	}
)

func dataSourceHostRecord() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHostRecordRead,
		Schema: map[string]*schema.Schema{
			"ref": {
				Type:          schema.TypeString,
				Description:   "Reference id of host record object.",
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"hostname"},
				AtLeastOneOf:  dataHostRecordRequiredSearchFields,
			},
			"hostname": {
				Type:          schema.TypeString,
				Description:   "The host name in FQDN format.",
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"ref"},
				AtLeastOneOf:  dataHostRecordRequiredSearchFields,
			},
			"comment": {
				Type:        schema.TypeString,
				Description: "Comment for the record; maximum 256 characters.",
				Computed:    true,
			},
			"enable_dns": {
				Type:        schema.TypeBool,
				Description: "When false, the host does not have parent zone information.",
				Computed:    true,
			},
			"network_view": {
				Type:        schema.TypeString,
				Description: "The name of the network view in which the host record resides.",
				Optional:    true,
				Computed:    true,
			},
			"view": {
				Type:        schema.TypeString,
				Description: "The name of the DNS view in which the record resides.",
				Computed:    true,
			},
			"zone": {
				Type:        schema.TypeString,
				Description: "The name of the zone in which the record resides.",
				Optional:    true,
				Computed:    true,
			},
			"query_params": {
				Type:        schema.TypeMap,
				Description: "Additional query parameters.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ip_v4_address": {
				Type:        schema.TypeSet,
				Description: "IPv4 addresses associated with host record.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ref": {
							Type:        schema.TypeString,
							Description: "Reference id of address object.",
							Computed:    true,
						},
						"ip_address": {
							Type:        schema.TypeString,
							Description: "IP address.",
							Computed:    true,
						},
						"hostname": {
							Type:        schema.TypeString,
							Description: "Hostname associated with IP address.",
							Computed:    true,
						},
						"network": {
							Type:        schema.TypeString,
							Description: "Network associated with IP address.",
							Computed:    true,
						},
						"mac_address": {
							Type:        schema.TypeString,
							Description: "MAC address associated with IP address.",
							Computed:    true,
						},
						"configure_for_dhcp": {
							Type:        schema.TypeBool,
							Description: "Set this to True to enable the DHCP configuration for this host address.",
							Computed:    true,
						},
					},
				},
			},
			"extensible_attributes": {
				Type:        schema.TypeMap,
				Description: "Extensible attributes of host record (Values are JSON encoded).",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceHostRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var record infoblox.HostRecord

	ref := d.Get("ref").(string)
	hostname := d.Get("hostname").(string)
	networkView := d.Get("network_view").(string)

	queryParams := d.Get("query_params").(map[string]interface{})
	resolvedQueryParams := make(map[string]string)

	for k, v := range queryParams {
		resolvedQueryParams[k] = v.(string)
	}

	if ref != "" {
		r, err := client.GetHostRecordByRef(ref, nil)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
		record = r
	} else {
		resolvedQueryParams["name"] = hostname
		if networkView != "" {
			resolvedQueryParams["network_view"] = networkView
		}
		r, err := client.GetHostRecordByQuery(resolvedQueryParams)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
		if len(r) > 1 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Multiple data results found",
				Detail:   "The provided hostname matched multiple record hosts",
			})
			return diags
		}
		record = r[0]
	}

	check := convertHostRecordToResourceData(client, d, &record)
	if check.HasError() {
		return check
	}
	d.SetId(record.Ref)

	return diags
}
