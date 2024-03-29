package infoblox

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

var (
	dataGridRequiredSearchFields = []string{
		"name",
		"ref",
	}
)

func dataSourceGrid() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGridRead,
		Schema: map[string]*schema.Schema{
			"ref": {
				Type:          schema.TypeString,
				Description:   "Reference id of grid object.",
				Computed:      true,
				Optional:      true,
				AtLeastOneOf:  dataGridRequiredSearchFields,
				ConflictsWith: remove(dataGridRequiredSearchFields, "ref", true),
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "Name of grid.",
				Computed:      true,
				Optional:      true,
				AtLeastOneOf:  dataGridRequiredSearchFields,
				ConflictsWith: remove(dataGridRequiredSearchFields, "name", true),
			},
			"service_status": {
				Type:        schema.TypeString,
				Description: "Service status of grid.",
				Computed:    true,
			},
			"dns_resolvers": {
				Type:        schema.TypeList,
				Description: "List of DNS resolvers.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"dns_search_domains": {
				Type:        schema.TypeList,
				Description: "List of DNS search domains.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"query_params": {
				Type:          schema.TypeMap,
				Description:   "Additional query parameters.",
				Optional:      true,
				ConflictsWith: []string{"ref"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceGridRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var grid infoblox.Grid

	queryParams := d.Get("query_params").(map[string]interface{})
	resolvedQueryParams := make(map[string]string)

	for k, v := range queryParams {
		resolvedQueryParams[k] = v.(string)
	}

	if ref, ok := d.GetOk("ref"); ok {
		g, err := client.GetGridByRef(ref.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		grid = g
	} else if name, ok := d.GetOk("name"); ok {
		resolvedQueryParams["name"] = name.(string)
		grids, err := client.GetGridsByQuery(resolvedQueryParams)

		if err != nil {
			return diag.FromErr(err)
		}
		if grids == nil || len(grids) == 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "No results found",
				Detail:   "The provided name did not match any grids",
			})
			return diags
		}
		if len(grids) > 1 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Multiple values returned",
				Detail:   "Grid query returned multiple values when one was expected",
			})
			return diags
		}

		grid = grids[0]
	}

	d.Set("ref", grid.Ref)
	d.Set("name", grid.Name)
	d.Set("service_status", grid.ServiceStatus)

	var dnsResolversList []string
	for _, resolver := range grid.DNSResolverSetting.Resolvers {
		dnsResolversList = append(dnsResolversList, resolver)
	}

	d.Set("dns_resolvers", dnsResolversList)

	var dnsSearchDomainsList []string
	for _, domain := range grid.DNSResolverSetting.SearchDomains {
		dnsSearchDomainsList = append(dnsSearchDomainsList, domain)
	}

	d.Set("dns_search_domains", dnsSearchDomainsList)

	d.SetId(grid.Ref)

	return diags
}
