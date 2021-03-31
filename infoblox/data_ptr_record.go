package infoblox

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

var (
	dataPtrRecordRequiredSearchFields = []string{
		"ref",
		"name",
		"pointer_domain_name",
		"ip_v4_address",
		"ip_v6_address",
	}
)

func dataSourcePtrRecord() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePtrRecordRead,
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
				Description: "The name for a DNS PTR record in punycode format.",
				Computed:    true,
			},
			"dns_pointer_domain_name": {
				Type:        schema.TypeString,
				Description: "The domain name of the DNS PTR record in punycode format.",
				Computed:    true,
			},
			"extensible_attributes": {
				Type:        schema.TypeMap,
				Description: "Extensible attributes of ptr record (Values are JSON encoded).",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ip_v4_address": {
				Type:             schema.TypeString,
				Description:      "The IPv4 Address of the record.",
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv4Address),
				AtLeastOneOf:     dataPtrRecordRequiredSearchFields,
				ConflictsWith:    []string{"ref", "ip_v6_address"},
			},
			"ip_v6_address": {
				Type:             schema.TypeString,
				Description:      "The IPv6 Address of the record.",
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv6Address),
				AtLeastOneOf:     dataPtrRecordRequiredSearchFields,
				ConflictsWith:    []string{"ref", "ip_v4_address"},
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the DNS PTR record in FQDN format.",
				Optional:      true,
				Computed:      true,
				AtLeastOneOf:  dataPtrRecordRequiredSearchFields,
				ConflictsWith: []string{"ref", "pointer_domain_name"},
			},
			"pointer_domain_name": {
				Type:          schema.TypeString,
				Description:   "The domain name of the DNS PTR record in FQDN format.",
				Optional:      true,
				Computed:      true,
				AtLeastOneOf:  dataPtrRecordRequiredSearchFields,
				ConflictsWith: []string{"ref", "name"},
			},
			"query_params": {
				Type:        schema.TypeMap,
				Description: "Additional query parameters.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ref": {
				Type:          schema.TypeString,
				Description:   "Reference id of ptr record object.",
				Optional:      true,
				Computed:      true,
				AtLeastOneOf:  dataPtrRecordRequiredSearchFields,
				ConflictsWith: remove(dataPtrRecordRequiredSearchFields, "ref", true),
			},
			"view": {
				Type:        schema.TypeString,
				Description: "Name of the DNS View in which the record resides.",
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

func dataSourcePtrRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var record infoblox.PtrRecord

	if ref, ok := d.GetOk("ref"); ok {
		r, err := client.GetPtrRecordByRef(ref.(string), nil)
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
		if view, ok := d.GetOk("view"); ok {
			resolvedQueryParams["view"] = view.(string)
		}
		if ip_v4_address, ok := d.GetOk("ip_v4_address"); ok {
			resolvedQueryParams["ipv4addr"] = ip_v4_address.(string)
		}
		if ip_v6_address, ok := d.GetOk("ip_v6_address"); ok {
			resolvedQueryParams["ipv6addr"] = ip_v6_address.(string)
		}
		if name, ok := d.GetOk("name"); ok {
			resolvedQueryParams["name"] = name.(string)
			r, err := client.GetPtrRecordByQuery(resolvedQueryParams)
			if err != nil {
				return diag.FromErr(err)
			}
			if r == nil || len(r) == 0 {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "No results found",
					Detail:   "The provided hostname did not match any PTR records",
				})
				return diags
			}
			if len(r) > 1 {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Multiple data results found",
					Detail:   "The provided hostname matched multiple PTR records when one was expected",
				})
				return diags
			}
			record = r[0]
		} else if pointer_domain_name, ok := d.GetOk("pointer_domain_name"); ok {
			resolvedQueryParams["ptrdname"] = pointer_domain_name.(string)
			r, err := client.GetPtrRecordByQuery(resolvedQueryParams)
			if err != nil {
				return diag.FromErr(err)
			}
			if r == nil || len(r) == 0 {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "No results found",
					Detail:   "The provided ptrdname did not match any PTR records",
				})
				return diags
			}
			if len(r) > 1 {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Multiple data results found",
					Detail:   "The provided ptrdname matched multiple PTR records when one was expected",
				})
				return diags
			}
			record = r[0]
		} else if dns_name, ok := d.GetOk("dns_name"); ok {
			resolvedQueryParams["dns_name"] = dns_name.(string)
			r, err := client.GetPtrRecordByQuery(resolvedQueryParams)
			if err != nil {
				return diag.FromErr(err)
			}
			if r == nil || len(r) == 0 {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "No results found",
					Detail:   "The provided DNS name did not match any PTR records",
				})
				return diags
			}
			if len(r) > 1 {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Multiple data results found",
					Detail:   "The provided DNS name matched multiple PTR records when one was expected",
				})
				return diags
			}
			record = r[0]
		} else if dns_pointer_domain_name, ok := d.GetOk("dns_pointer_domain_name"); ok {
			resolvedQueryParams["dns_ptrdname"] = dns_pointer_domain_name.(string)
			r, err := client.GetPtrRecordByQuery(resolvedQueryParams)
			if err != nil {
				return diag.FromErr(err)
			}
			if r == nil || len(r) == 0 {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "No results found",
					Detail:   "The provided DNS ptrdname did not match any PTR records",
				})
				return diags
			}
			if len(r) > 1 {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Multiple data results found",
					Detail:   "The provided DNS ptrdname matched multiple PTR records when one was expected",
				})
				return diags
			}
			record = r[0]
		} else {
			r, err := client.GetPtrRecordByQuery(resolvedQueryParams)
			if err != nil {
				return diag.FromErr(err)
			}
			if r == nil || len(r) == 0 {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "No results found",
					Detail:   "The provided IP address did not match any PTR records",
				})
				return diags
			}
			if len(r) > 1 {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Multiple data results found",
					Detail:   "The provided IP address matched multiple PTR records when one was expected",
				})
				return diags
			}
			record = r[0]
		}
	}

	check := convertPtrRecordToResourceData(client, d, &record)
	if check.HasError() {
		return check
	}
	d.SetId(record.Ref)

	return diags
}
