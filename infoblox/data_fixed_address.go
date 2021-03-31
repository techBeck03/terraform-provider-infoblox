package infoblox

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
			"cidr": {
				Type:        schema.TypeString,
				Description: "The network to which this fixed address belongs, in IPv4 Address/CIDR format.",
				Computed:    true,
			},
			"comment": {
				Type:        schema.TypeString,
				Description: "Comment for the fixed address; maximum 256 characters.",
				Computed:    true,
			},
			"disable": {
				Type:        schema.TypeBool,
				Description: "Determines whether a fixed address is disabled or not. When this is set to False, the fixed address is enabled.",
				Computed:    true,
			},
			"extensible_attributes": {
				Type:        schema.TypeMap,
				Description: "Extensible attributes of fixed address (Values are JSON encoded).",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"hostname": {
				Type:          schema.TypeString,
				Description:   "This field contains the name of this fixed address.",
				Optional:      true,
				Computed:      true,
				AtLeastOneOf:  dataFixedAddressRequiredSearchFields,
				ConflictsWith: remove(dataFixedAddressRequiredSearchFields, "hostname", true),
			},
			"ip_address": {
				Type:             schema.TypeString,
				Description:      "The IPv4 Address of the fixed address.",
				Computed:         true,
				Optional:         true,
				AtLeastOneOf:     dataFixedAddressRequiredSearchFields,
				ConflictsWith:    remove(dataFixedAddressRequiredSearchFields, "ip_address", true),
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv4Address),
			},
			"mac": {
				Type:        schema.TypeString,
				Description: "The MAC address value for this fixed address.",
				Computed:    true,
			},
			"match_client": {
				Type:        schema.TypeString,
				Description: "The match_client value for this fixed address.",
				Computed:    true,
			},
			"network_view": {
				Type:        schema.TypeString,
				Description: "The name of the network view in which this fixed address resides.",
				Optional:    true,
				Computed:    true,
			},
			"option": {
				Type:        schema.TypeSet,
				Description: "An array of DHCP option structs that lists the DHCP options associated with the object.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the DHCP option.",
							Computed:    true,
						},
						"code": {
							Type:        schema.TypeInt,
							Description: "The code of the DHCP option.",
							Computed:    true,
						},
						"use_option": {
							Type:        schema.TypeBool,
							Description: "Only applies to special options that are displayed separately from other options and have a use flag.",
							Computed:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "Value of the DHCP option.",
							Computed:    true,
						},
						"vendor_class": {
							Type:        schema.TypeString,
							Description: "The name of the space this DHCP option is associated to.",
							Computed:    true,
						},
					},
				},
			},
			"query_params": {
				Type:        schema.TypeMap,
				Description: "Additional query parameters",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ref": {
				Type:          schema.TypeString,
				Description:   "Reference id of fixed address object.",
				Optional:      true,
				Computed:      true,
				AtLeastOneOf:  dataFixedAddressRequiredSearchFields,
				ConflictsWith: remove(dataFixedAddressRequiredSearchFields, "ref", true),
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
		if f == nil || len(f) == 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "No results found",
				Detail:   "The provided IP address did not match any fixed addresses",
			})
			return diags
		}
		if len(f) > 1 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Multiple data results found",
				Detail:   "The provided IP address matched multiple fixed addresses when one was expected",
			})
			return diags
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
		if f == nil || len(f) == 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "No results found",
				Detail:   "The provided hostname did not match any fixed addresses",
			})
			return diags
		}
		if len(f) > 1 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Multiple data results found",
				Detail:   "The provided hostname matched multiple fixed addresses when one was expected",
			})
			return diags
		}
		fixedAddress = f[0]
	}

	check := convertFixedAddressToResourceData(client, d, &fixedAddress)
	if check.HasError() {
		return check
	}

	d.SetId(fixedAddress.Ref)

	return diags
}
