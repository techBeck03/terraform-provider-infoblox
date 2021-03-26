package infoblox

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

func dataSourceSequentialAddressBlock() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSequentialAddressBlockRead,
		Schema: map[string]*schema.Schema{
			"cidr": {
				Type:             schema.TypeString,
				Description:      "Network for address block in IPv4 Address/CIDR format.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
			},
			"address_count": {
				Type:        schema.TypeInt,
				Description: "Number of IPs to allocate.",
				Required:    true,
			},
			"network_view": {
				Type:        schema.TypeString,
				Description: "The name of the network view in which the address block resides.",
				Optional:    true,
				Computed:    true,
			},
			"addresses": {
				Type:        schema.TypeList,
				Description: "List of sequential ip address objects.",
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
							Description: "IPv4 address.",
							Computed:    true,
						},
						"hostnames": {
							Type:        schema.TypeList,
							Description: "List of hostnames associated with IP address.",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"mac_address": {
							Type:        schema.TypeString,
							Description: "MAC address associated with IP address.",
							Computed:    true,
						},
						"network_view": {
							Type:        schema.TypeString,
							Description: "Network view associated with IP address.",
							Computed:    true,
						},
						"cidr": {
							Type:        schema.TypeString,
							Description: "CIDR associated with IP address.",
							Computed:    true,
						},
						"usage": {
							Type:        schema.TypeList,
							Description: "Usage associated with IP address.",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"types": {
							Type:        schema.TypeList,
							Description: "Types associated with IP address.",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"objects": {
							Type:        schema.TypeList,
							Description: "Objects associated with IP address.",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Status associated with IP address.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSequentialAddressBlockRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	cidr := d.Get("cidr").(string)
	count := d.Get("address_count").(int)

	addresses, err := client.GetSequentialAddressRange(infoblox.AddressQuery{
		CIDR:  cidr,
		Count: count,
	})

	if err != nil {
		return diag.FromErr(err)
	}
	prettyPrint(addresses)

	var ipAddressList []map[string]interface{}
	for _, address := range *addresses {
		var hostnames []string
		for _, hostname := range address.Hostnames {
			hostnames = append(hostnames, hostname)
		}
		var addressUses []string
		for _, addressUse := range address.Usage {
			addressUses = append(addressUses, addressUse)
		}
		var addressTypes []string
		for _, addressType := range address.Types {
			addressTypes = append(addressTypes, addressType)
		}
		var addressObjects []string
		for _, addressObject := range address.Objects {
			addressObjects = append(addressObjects, addressObject)
		}
		ipAddressList = append(ipAddressList, map[string]interface{}{
			"ref":          address.Ref,
			"ip_address":   address.IPAddress,
			"hostnames":    hostnames,
			"mac_address":  address.Mac,
			"network_view": address.NetworkView,
			"status":       address.Status,
			"usage":        addressUses,
			"types":        addressTypes,
			"objects":      addressObjects,
		})
	}

	d.Set("addresses", ipAddressList)
	d.SetId("test")

	return diags
}
