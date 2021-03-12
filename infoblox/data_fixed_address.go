package infoblox

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

var (
	dataFixedAddressRequiredSearchFields = []string{
		"hostname",
		"ref",
		"ip_address",
	}
)

func dataSourceFixedAddress() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFixedAddressRead,
		Schema: map[string]*schema.Schema{
			"ref": {
				Type:          schema.TypeString,
				Description:   "Reference id of host fixed address object",
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"hostname", "ip_address"},
				AtLeastOneOf:  dataFixedAddressRequiredSearchFields,
			},
			"hostname": {
				Type:          schema.TypeString,
				Description:   "Hostname of host fixed address",
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"ref", "ip_address"},
				AtLeastOneOf:  dataFixedAddressRequiredSearchFields,
			},
			"cidr": {
				Type:        schema.TypeString,
				Description: "Hostname of host fixed address",
				Computed:    true,
			},
			"query_params": {
				Type:        schema.TypeMap,
				Description: "Additional query parameters",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"comment": {
				Type:        schema.TypeString,
				Description: "Comment string",
				Computed:    true,
			},
			"ip_address": {
				Type:          schema.TypeBool,
				Description:   "IPv4 address",
				Computed:      true,
				Optional:      true,
				ConflictsWith: []string{"ref", "hostname"},
				AtLeastOneOf:  dataFixedAddressRequiredSearchFields,
			},
			"network_view": {
				Type:        schema.TypeString,
				Description: "Network view",
				Optional:    true,
				Computed:    true,
			},
			"mac": {
				Type:        schema.TypeString,
				Description: "MAC address",
				Computed:    true,
			},
			"disable": {
				Type:        schema.TypeBool,
				Description: "Disabled",
				Computed:    true,
			},
			"option": {
				Type:        schema.TypeSet,
				Description: "DHCP options associated with network",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Name of DHCP option",
							Computed:    true,
						},
						"code": {
							Type:        schema.TypeInt,
							Description: "Option numberic id",
							Computed:    true,
						},
						"use_option": {
							Type:        schema.TypeBool,
							Description: "Use this dhcp option",
							Computed:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "Value of option",
							Computed:    true,
						},
						"vendor_class": {
							Type:        schema.TypeString,
							Description: "Value of option",
							Computed:    true,
						},
					},
				},
			},
			"extensible_attributes": {
				Type:        schema.TypeMap,
				Description: "Extensible attributes of host fixed address",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceFixedAddressRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var fixedAddress infoblox.FixedAddress

	if ref, ok := d.GetOk("ref"); ok {
		f, err := client.GetFixedAddressByRef(ref.(string), nil)
		if err != nil {
			return diag.FromErr(err)
		}
		fixedAddress = f
	} else if ipAddress, ok := d.GetOk("ip_address"); ok {
		queryParams := d.Get("query_params").(map[string]interface{})
		resolvedQueryParams := make(map[string]string)

		for k, v := range queryParams {
			resolvedQueryParams[k] = v.(string)
		}
		resolvedQueryParams["ipv4addr"] = ipAddress.(string)
		f, err := client.GetFixedAddressByQuery(resolvedQueryParams)
		if err != nil {
			return diag.FromErr(err)
		}
		fixedAddress = f[0]
	} else if ipAddress, ok := d.GetOk("hostname"); ok {
		queryParams := d.Get("query_params").(map[string]interface{})
		resolvedQueryParams := make(map[string]string)

		for k, v := range queryParams {
			resolvedQueryParams[k] = v.(string)
		}
		resolvedQueryParams["hostname"] = ipAddress.(string)
		f, err := client.GetFixedAddressByQuery(resolvedQueryParams)
		if err != nil {
			return diag.FromErr(err)
		}
		fixedAddress =
			f[0]
	}

	d.SetId(fixedAddress.Ref)

	return diags
}
