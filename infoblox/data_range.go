package infoblox

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

var (
	dataRangeRequiredFields = []string{
		"ref",
		"cidr",
	}
)

func dataSourceRange() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRangeRead,
		Schema: map[string]*schema.Schema{
			"address_list": {
				Type:        schema.TypeList,
				Description: "The list of IP Addresses associated with this range",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cidr": {
				Type:             schema.TypeString,
				Description:      "The network to which this range belongs, in IPv4 Address/CIDR format.",
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
				AtLeastOneOf:     dataRangeRequiredFields,
				ConflictsWith:    []string{"ref"},
			},
			"comment": {
				Type:        schema.TypeString,
				Description: "Comment for the range; maximum 256 characters.",
				Computed:    true,
			},
			"disable_dhcp": {
				Type:        schema.TypeBool,
				Description: "Disable for DHCP.",
				Computed:    true,
			},
			"end_address": {
				Type:             schema.TypeString,
				Description:      "The IPv4 Address end address of the range.",
				Optional:         true,
				Computed:         true,
				RequiredWith:     []string{"cidr"},
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv4Address),
			},
			"extensible_attributes": {
				Type:        schema.TypeMap,
				Description: "Extensible attributes of range object (Values are JSON encoded).",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"member": {
				Type:        schema.TypeSet,
				Description: "Grid members associated with range.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"struct": {
							Type:        schema.TypeString,
							Description: "Struct type of member.",
							Computed:    true,
						},
						"ip_v4_address": {
							Type:        schema.TypeString,
							Description: "IPv4 address.",
							Computed:    true,
						},
						"ip_v6_address": {
							Type:        schema.TypeString,
							Description: "IPv6 address.",
							Computed:    true,
						},
						"hostname": {
							Type:        schema.TypeString,
							Description: "Hostname of member.",
							Computed:    true,
						},
					},
				},
			},
			"network_view": {
				Type:        schema.TypeString,
				Description: "The name of the network view in which this range resides.",
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
				Description: "Additional query parameters.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"range_function_string": {
				Type:        schema.TypeString,
				Description: "String representation of start and end addresses to be used with function calls.",
				Computed:    true,
			},
			"ref": {
				Type:          schema.TypeString,
				Description:   "Reference id of range object.",
				Optional:      true,
				Computed:      true,
				AtLeastOneOf:  dataRangeRequiredFields,
				ConflictsWith: []string{"cidr"},
			},
			"start_address": {
				Type:             schema.TypeString,
				Description:      "The IPv4 Address starting address of the range.",
				Optional:         true,
				Computed:         true,
				RequiredWith:     []string{"cidr"},
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv4Address),
			},
		},
	}
}

func dataSourceRangeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var addressRange infoblox.Range

	if ref, ok := d.GetOk("ref"); ok {
		r, err := client.GetRangeByRef(ref.(string), nil)
		if err != nil {
			return diag.FromErr(err)
		}
		addressRange = r
	} else {
		queryParams := d.Get("query_params").(map[string]interface{})
		resolvedQueryParams := make(map[string]string)

		for k, v := range queryParams {
			resolvedQueryParams[k] = v.(string)
		}
		resolvedQueryParams["network"] = d.Get("cidr").(string)
		resolvedQueryParams["start_addr"] = d.Get("start_address").(string)
		resolvedQueryParams["end_addr"] = d.Get("end_address").(string)
		r, err := client.GetRangeByQuery(resolvedQueryParams)
		if err != nil {
			return diag.FromErr(err)
		}
		if r == nil || len(r) == 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "No results found",
				Detail:   "The provided cidr, start_addr, end_addr did not match any range",
			})
			return diags
		}
		if len(r) > 1 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Multiple data results found",
				Detail:   "The provided cidr, start_addr, end_addr matched multiple ranges",
			})
			return diags
		}
		addressRange = r[0]
	}

	check := convertRangeToResourceData(client, d, &addressRange)
	if check.HasError() {
		return check
	}
	d.SetId(addressRange.Ref)

	return diags
}
