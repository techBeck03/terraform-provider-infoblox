package infoblox

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

var (
	dataCNameRecordRequiredSearchFields = []string{
		"alias",
		"ref",
		"canonical",
		"dns_name",
	}
)

func dataSourceCNameRecord() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCNameRecordRead,
		Schema: map[string]*schema.Schema{
			"ref": {
				Type:         schema.TypeString,
				Description:  "Reference id of A record object",
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: dataCNameRecordRequiredSearchFields,
			},
			"alias": {
				Type:          schema.TypeString,
				Description:   "Alias of A record",
				Optional:      true,
				Computed:      true,
				AtLeastOneOf:  dataCNameRecordRequiredSearchFields,
				ConflictsWith: []string{"ref"},
			},
			"canonical": {
				Type:          schema.TypeString,
				Description:   "Canonical name",
				Optional:      true,
				Computed:      true,
				AtLeastOneOf:  dataCNameRecordRequiredSearchFields,
				ConflictsWith: []string{"ref"},
			},
			"dns_name": {
				Type:          schema.TypeString,
				Description:   "DNS name of A record",
				Optional:      true,
				Computed:      true,
				AtLeastOneOf:  dataCNameRecordRequiredSearchFields,
				ConflictsWith: []string{"ref"},
			},
			"comment": {
				Type:        schema.TypeString,
				Description: "Comment string",
				Computed:    true,
			},
			"disable": {
				Type:        schema.TypeBool,
				Description: "Disable",
				Computed:    true,
			},
			"view": {
				Type:        schema.TypeString,
				Description: "DNS view",
				Optional:    true,
				Computed:    true,
			},
			"zone": {
				Type:        schema.TypeString,
				Description: "DNS zone",
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
			"extensible_attributes": {
				Type:        schema.TypeMap,
				Description: "Extensible attributes of A record",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceCNameRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var record infoblox.CNameRecord

	if ref, ok := d.GetOk("ref"); ok {
		r, err := client.GetCNameRecordByRef(ref.(string), nil)
		if err != nil {
			return diag.FromErr(err)
		}
		record = r
	} else {
		queryParams := d.Get("query_params").(map[string]interface{})
		resolvedQueryParams := make(map[string]string)

		for k, v := range queryParams {
			resolvedQueryParams[k] = v.(string)
		}
		if zone, ok := d.GetOk("zone"); ok {
			resolvedQueryParams["zone"] = zone.(string)
		}
		if view, ok := d.GetOk("view"); ok {
			resolvedQueryParams["view"] = view.(string)
		}
		if zone, ok := d.GetOk("zone"); ok {
			resolvedQueryParams["zone"] = zone.(string)
		}
		if alias, ok := d.GetOk("alias"); ok {
			resolvedQueryParams["name"] = alias.(string)
			r, err := client.GetCNameRecordByQuery(resolvedQueryParams)
			if err != nil {
				return diag.FromErr(err)
			}
			record = r[0]
		} else if canonical, ok := d.GetOk("canonical"); ok {
			resolvedQueryParams["canonical"] = canonical.(string)
			r, err := client.GetCNameRecordByQuery(resolvedQueryParams)
			if err != nil {
				return diag.FromErr(err)
			}
			record = r[0]
		} else if dns_name, ok := d.GetOk("dns_name"); ok {
			resolvedQueryParams["dns_name"] = dns_name.(string)
			r, err := client.GetCNameRecordByQuery(resolvedQueryParams)
			if err != nil {
				return diag.FromErr(err)
			}
			record = r[0]
		}
	}

	check := convertCNameRecordToResourceData(client, d, &record)
	if check.HasError() {
		return check
	}
	d.SetId(record.Ref)

	return diags
}
