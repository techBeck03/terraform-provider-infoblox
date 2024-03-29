package infoblox

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

var (
	dataAliasRecordRequiredSearchFields = []string{
		"name",
		"ref",
	}
)

func dataSourceAliasRecord() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAliasRecordRead,
		Schema: map[string]*schema.Schema{
			"comment": {
				Type:        schema.TypeString,
				Description: "Comment for the record; maximum 256 characters.",
				Computed:    true,
			},
			"disable": {
				Type:        schema.TypeBool,
				Description: "Determines if the record is disabled or not. False means that the record is enabled.",
				Computed:    true,
			},
			"dns_name": {
				Type:        schema.TypeString,
				Description: "The name for an Alias record in punycode format.",
				Computed:    true,
			},
			"dns_target_name": {
				Type:        schema.TypeString,
				Description: "Target name in punycode format.",
				Computed:    true,
			},
			"extensible_attributes": {
				Type:        schema.TypeMap,
				Description: "Extensible attributes of alias record (Values are JSON encoded).",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name for an Alias record in FQDN format.",
				Optional:      true,
				Computed:      true,
				AtLeastOneOf:  dataAliasRecordRequiredSearchFields,
				ConflictsWith: []string{"ref"},
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
				Description:   "Reference id of alias record object.",
				Optional:      true,
				Computed:      true,
				AtLeastOneOf:  dataAliasRecordRequiredSearchFields,
				ConflictsWith: []string{"name"},
			},
			"target_name": {
				Type:        schema.TypeString,
				Description: "Target name in FQDN format.",
				Computed:    true,
			},
			"target_type": {
				Type:        schema.TypeString,
				Description: "Target type.",
				Computed:    true,
			},
			"view": {
				Type:        schema.TypeString,
				Description: "The name of the DNS View in which the record resides.",
				Optional:    true,
				Computed:    true,
			},
			"zone": {
				Type:        schema.TypeString,
				Description: "The name of the zone in which the record resides.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceAliasRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var record infoblox.AliasRecord

	if ref, ok := d.GetOk("ref"); ok {
		r, err := client.GetAliasRecordByRef(ref.(string), nil)
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
		if view, ok := d.GetOk("view"); ok {
			resolvedQueryParams["view"] = view.(string)
		}
		if name, ok := d.GetOk("name"); ok {
			resolvedQueryParams["name"] = name.(string)
			r, err := client.GetAliasRecordByQuery(resolvedQueryParams)
			if err != nil {
				return diag.FromErr(err)
			}
			if r == nil || len(r) == 0 {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "No results found",
					Detail:   "The provided hostname did not match any alias records",
				})
				return diags
			}
			if len(r) > 1 {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Multiple data results found",
					Detail:   "The provided hostname matched multiple alias records when one was expected",
				})
				return diags
			}
			record = r[0]
		} else if dns_name, ok := d.GetOk("dns_name"); ok {
			resolvedQueryParams["dns_name"] = dns_name.(string)
			r, err := client.GetAliasRecordByQuery(resolvedQueryParams)
			if err != nil {
				return diag.FromErr(err)
			}
			if r == nil || len(r) == 0 {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "No results found",
					Detail:   "The provided DNS name did not match any alias records",
				})
				return diags
			}
			if len(r) > 1 {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Multiple data results found",
					Detail:   "The provided DNS name matched multiple alias records when one was expected",
				})
				return diags
			}
			record = r[0]
		}
	}

	check := convertAliasRecordToResourceData(client, d, &record)
	if check.HasError() {
		return check
	}
	d.SetId(record.Ref)

	return diags
}
