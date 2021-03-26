package infoblox

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

var (
	dataNetworkRequiredSearchFields = []string{
		"cidr",
		"ref",
	}
)

func dataSourceNetwork() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkRead,
		Schema: map[string]*schema.Schema{
			"ref": {
				Type:          schema.TypeString,
				Description:   "Reference id of network object.",
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"cidr"},
				AtLeastOneOf:  dataNetworkRequiredSearchFields,
			},
			"cidr": {
				Type:          schema.TypeString,
				Description:   "The network address in IPv4 Address/CIDR format.",
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"ref"},
				AtLeastOneOf:  dataNetworkRequiredSearchFields,
			},
			"comment": {
				Type:        schema.TypeString,
				Description: "Comment for the record; maximum 256 characters.",
				Computed:    true,
			},
			"disable_dhcp": {
				Type:        schema.TypeBool,
				Description: "Disable for DHCP.",
				Computed:    true,
			},
			"network_view": {
				Type:        schema.TypeString,
				Description: "Network view",
				Optional:    true,
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
			"member": {
				Type:        schema.TypeSet,
				Description: "Grid members associated with network.",
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
			"extensible_attributes": {
				Type:        schema.TypeMap,
				Description: "Extensible attributes of network (Values are JSON encoded).",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var record infoblox.Network

	ref := d.Get("ref").(string)
	cidr := d.Get("cidr").(string)
	networkView := d.Get("network_view").(string)

	queryParams := d.Get("query_params").(map[string]interface{})
	resolvedQueryParams := make(map[string]string)

	for k, v := range queryParams {
		resolvedQueryParams[k] = v.(string)
	}

	if ref != "" {
		r, err := client.GetNetworkByRef(ref, nil)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
		record = r
	} else {
		resolvedQueryParams["network"] = cidr
		if networkView != "" {
			resolvedQueryParams["network_view"] = networkView
		}
		r, err := client.GetNetworkByQuery(resolvedQueryParams)
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

	check := convertNetworkToResourceData(client, d, &record)
	if check.HasError() {
		return check
	}
	d.SetId(record.Ref)

	return diags
}
