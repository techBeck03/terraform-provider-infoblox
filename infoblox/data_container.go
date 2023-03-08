package infoblox

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

func dataSourceContainer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerRead,
		Schema: map[string]*schema.Schema{
			"cidr": {
				Type:             schema.TypeString,
				Description:      "The container network address in IPv4 Address/CIDR format.",
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
				ConflictsWith:    []string{"ref"},
				AtLeastOneOf:     []string{"cidr", "ref"},
			},
			"comment": {
				Type:        schema.TypeString,
				Description: "Comment for the container; maximum 256 characters.",
				Computed:    true,
			},
			"extensible_attributes": {
				Type:        schema.TypeMap,
				Description: "Extensible attributes of network (Values are JSON encoded).",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"network_view": {
				Type:        schema.TypeString,
				Description: "The name of the network view in which this network resides.",
				Computed:    true,
			},
			"ref": {
				Type:          schema.TypeString,
				Description:   "Reference id of A container object.",
				Optional:      true,
				ConflictsWith: []string{"cidr"},
				AtLeastOneOf:  []string{"cidr", "ref"},
			},
		},
	}
}

func dataSourceContainerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics
	var container infoblox.NetworkContainer

	ref := d.Get("ref").(string)
	if ref != "" {
		c, err := client.GetContainerByRef(ref, nil)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
		container = c
	} else {
		cidr := d.Get("cidr").(string)
		query_params := make(map[string]string)
		query_params["network"] = cidr
		c, err := client.GetContainerByQuery(query_params)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
		if len(c) > 1 {
			diags = append(diags, diag.Errorf("More than one container found for supplied cidr")...)
			return diags
		}
		container = c[0]
	}

	check := convertContainerToResourceData(client, d, &container)
	if check.HasError() {
		return check
	}
	d.SetId(container.Ref)

	return diags
}
