package infoblox

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

var (
	dataGridMemberRequiredSearchFields = []string{
		"hostname",
		"ref",
	}
)

func dataSourceGridMember() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGridMemberRead,
		Schema: map[string]*schema.Schema{
			"ref": {
				Type:          schema.TypeString,
				Description:   "Reference id of member.",
				Computed:      true,
				Optional:      true,
				AtLeastOneOf:  dataGridMemberRequiredSearchFields,
				ConflictsWith: remove(dataGridMemberRequiredSearchFields, "ref", true),
			},
			"hostname": {
				Type:          schema.TypeString,
				Description:   "Hostname of member in FQDN format.",
				Computed:      true,
				Optional:      true,
				AtLeastOneOf:  dataGridMemberRequiredSearchFields,
				ConflictsWith: remove(dataGridMemberRequiredSearchFields, "hostname", true),
			},
			"config_address_type": {
				Type:        schema.TypeString,
				Description: "Configured IP address type.",
				Computed:    true,
			},
			"service_type_configuration": {
				Type:        schema.TypeString,
				Description: "Service type configuration.",
				Computed:    true,
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

func dataSourceGridMemberRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var member infoblox.GridMember

	queryParams := d.Get("query_params").(map[string]interface{})
	resolvedQueryParams := make(map[string]string)

	for k, v := range queryParams {
		resolvedQueryParams[k] = v.(string)
	}

	if ref, ok := d.GetOk("ref"); ok {
		m, err := client.GetGridMembersByRef(ref.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		member = m
	} else if hostname, ok := d.GetOk("hostname"); ok {
		resolvedQueryParams["host_name"] = hostname.(string)
		members, err := client.GetGridMembersByQuery(resolvedQueryParams)

		if err != nil {
			return diag.FromErr(err)
		}

		if members == nil || len(members) == 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "No results found",
				Detail:   "The provided hostname did not match any grid members",
			})
			return diags
		}
		if len(members) > 1 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Multiple values returned",
				Detail:   "Grid member query returned multiple values when one was expected",
			})
			return diags
		}
		member = members[0]
	}

	d.Set("ref", member.Ref)
	d.Set("hostname", member.Hostname)
	d.Set("config_address_type", member.ConfigAddressType)
	d.Set("service_type_configuration", member.ServiceTypeConfiguration)

	d.SetId(member.Ref)

	return diags
}
