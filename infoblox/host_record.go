package infoblox

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

func resourceHostRecord() *schema.Resource {
	return &schema.Resource{
		// CreateContext: resourceHostRecordCreate,
		ReadContext: resourceHostRecordRead,
		Schema: map[string]*schema.Schema{
			"ref": {
				Type:        schema.TypeString,
				Description: "Reference id of host record object",
				Optional:    true,
				Computed:    true,
			},
			"hostname": {
				Type:        schema.TypeString,
				Description: "Hostname of host record",
				Optional:    true,
				Computed:    true,
			},
			"comment": {
				Type:        schema.TypeString,
				Description: "Comment string",
				Computed:    true,
			},
			"enable_dns": {
				Type:        schema.TypeString,
				Description: "Enable for DNS",
				Computed:    true,
			},
			"network_view": {
				Type:        schema.TypeString,
				Description: "Network view",
				Optional:    true,
				Computed:    true,
			},
			"view": {
				Type:        schema.TypeString,
				Description: "DNS view",
				Computed:    true,
			},
			"zone": {
				Type:        schema.TypeString,
				Description: "DNS zone",
				Optional:    true,
				Computed:    true,
			},
			"extensible_attributes": {
				Type:        schema.TypeMap,
				Description: "Extensible attributes of host record",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func convertHostRecordToResourceData(client *infoblox.Client, d *schema.ResourceData, record *infoblox.HostRecord) diag.Diagnostics {
	var diags diag.Diagnostics

	d.Set("ref", record.Ref)
	d.Set("hostname", record.Hostname)
	d.Set("comment", record.Comment)
	d.Set("enable_dns", record.EnableDNS)
	d.Set("network_view", record.NetworkView)
	d.Set("view", record.View)
	d.Set("zone", record.Zone)
	eas, err := client.ConvertEAsToJSONString(*record.ExtensibleAttributes)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	} else {
		d.Set("extensible_attributes", eas)
	}

	return diags
}

func resourceHostRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics
	var record infoblox.HostRecord

	ref := d.Get("ref").(string)
	hostname := d.Get("hostname").(string)
	networkView := d.Get("network_view").(string)

	if ref == "" && (hostname == "" || networkView == "") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error retrieving host record",
			Detail:   "Either `ref` or `hostname` + `network_view` must be supplied",
		})

		return diags
	}

	queryParams := d.Get("query_params").(map[string]string)

	if ref != "" {
		r, err := client.GetHostRecordByRef(ref, nil)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
		record = r
	} else {
		queryParams["name"] = hostname
		queryParams["network_view"] = networkView
		r, err := client.GetHostRecordByQuery(queryParams)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
		record = r
	}

	check := convertHostRecordToResourceData(client, d, &record)
	if check.HasError() {
		return check
	}

	d.SetId(record.Ref)

	return diags
}

// func resourceHostRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	client := m.(*infoblox.Client)

// 	// Warning or errors can be collected in a slice type
// 	var diags diag.Diagnostics

// 	return diags
// }
