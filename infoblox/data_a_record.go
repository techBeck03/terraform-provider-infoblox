package infoblox

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

var (
	dataARecordRequiredSearchFields = []string{
		"hostname",
		"ref",
		"ip_address",
	}
)

func dataSourceARecord() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceARecordRead,
		Schema: map[string]*schema.Schema{
			"ref": {
				Type:          schema.TypeString,
				Description:   "Reference id of A record object.",
				Optional:      true,
				Computed:      true,
				AtLeastOneOf:  dataARecordRequiredSearchFields,
				ConflictsWith: remove(dataARecordRequiredSearchFields, "ref", true),
			},
			"hostname": {
				Type:          schema.TypeString,
				Description:   "Hostname of A record.",
				Optional:      true,
				Computed:      true,
				AtLeastOneOf:  dataARecordRequiredSearchFields,
				ConflictsWith: remove(dataARecordRequiredSearchFields, "hostname", true),
			},
			"dns_name": {
				Type:        schema.TypeString,
				Description: "The name for an A record in punycode format.",
				Computed:    true,
			},
			"ip_address": {
				Type:             schema.TypeString,
				Description:      "The IPv4 Address of the record.",
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv4Address),
				AtLeastOneOf:     dataARecordRequiredSearchFields,
				ConflictsWith:    remove(dataARecordRequiredSearchFields, "ip_address", true),
			},
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
			"view": {
				Type:        schema.TypeString,
				Description: "The name of the DNS view in which the record resides.",
				Optional:    true,
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
				Description: "Additional query parameters",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"extensible_attributes": {
				Type:        schema.TypeMap,
				Description: "Extensible attributes of A record (Values are JSON encoded).",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceARecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var record infoblox.ARecord

	if ref, ok := d.GetOk("ref"); ok {
		r, err := client.GetARecordByRef(ref.(string), nil)
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
		if hostname, ok := d.GetOk("hostname"); ok {
			resolvedQueryParams["name"] = hostname.(string)
			r, err := client.GetARecordByQuery(resolvedQueryParams)
			if err != nil {
				return diag.FromErr(err)
			}
			record = r[0]
		} else if dns_name, ok := d.GetOk("dns_name"); ok {
			resolvedQueryParams["dns_name"] = dns_name.(string)
			r, err := client.GetARecordByQuery(resolvedQueryParams)
			if err != nil {
				return diag.FromErr(err)
			}
			record = r[0]
		} else if ip_address, ok := d.GetOk("ip_address"); ok {
			resolvedQueryParams["ipv4addr"] = ip_address.(string)
			r, err := client.GetARecordByQuery(resolvedQueryParams)
			if err != nil {
				return diag.FromErr(err)
			}
			record = r[0]
		}
	}

	check := convertARecordToResourceData(client, d, &record)
	if check.HasError() {
		return check
	}
	d.SetId(record.Ref)

	return diags
}
