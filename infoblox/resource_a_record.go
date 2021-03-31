package infoblox

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

func resourceARecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceARecordCreate,
		ReadContext:   resourceARecordRead,
		UpdateContext: resourceARecordUpdate,
		DeleteContext: resourceARecordDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
		CustomizeDiff: customdiff.Sequence(
			makeEACustomDiff("extensible_attributes"),
		),
		Schema: map[string]*schema.Schema{
			"comment": {
				Type:             schema.TypeString,
				Description:      "Comment for the record; maximum 256 characters.",
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 256)),
			},
			"disable": {
				Type:        schema.TypeBool,
				Description: "Determines if the record is disabled or not. False means that the record is enabled.",
				Optional:    true,
				Computed:    true,
			},
			"dns_name": {
				Type:        schema.TypeString,
				Description: "The name for an A record in punycode format.",
				Computed:    true,
			},
			"extensible_attributes": {
				Type:             schema.TypeMap,
				Description:      "Extensible attributes of A record (Values are JSON encoded).",
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validateEa,
				DiffSuppressFunc: eaSuppressDiff,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"hostname": {
				Type:        schema.TypeString,
				Description: "Hostname of A record.",
				Required:    true,
			},
			"ip_address": {
				Type:             schema.TypeString,
				Description:      "The IPv4 Address of the record.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv4Address),
			},
			"ref": {
				Type:        schema.TypeString,
				Description: "Reference id of A record object.",
				Computed:    true,
			},
			"view": {
				Type:        schema.TypeString,
				Description: "The name of the DNS view in which the record resides.",
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"zone": {
				Type:        schema.TypeString,
				Description: "The name of the zone in which the record resides.",
				Computed:    true,
			},
		},
	}
}

func convertARecordToResourceData(client *infoblox.Client, d *schema.ResourceData, record *infoblox.ARecord) diag.Diagnostics {
	var diags diag.Diagnostics

	d.Set("ref", record.Ref)
	d.Set("hostname", record.Hostname)
	d.Set("dns_name", record.DNSName)
	d.Set("ip_address", record.IPAddress)
	d.Set("comment", record.Comment)
	d.Set("disable", record.Disable)
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

func convertResourceDataToARecord(client *infoblox.Client, d *schema.ResourceData) (*infoblox.ARecord, error) {
	var record infoblox.ARecord

	record.Hostname = d.Get("hostname").(string)
	record.DNSName = d.Get("dns_name").(string)
	record.IPAddress = d.Get("ip_address").(string)
	record.Comment = d.Get("comment").(string)
	record.Disable = newBool(d.Get("disable").(bool))
	record.View = d.Get("view").(string)
	record.Zone = d.Get("zone").(string)

	eaMap := d.Get("extensible_attributes").(map[string]interface{})
	if len(eaMap) > 0 {
		eas, err := createExtensibleAttributesFromJSON(eaMap)
		if err != nil {
			return &record, err
		}
		record.ExtensibleAttributes = &eas
	}

	if client.OrchestratorEAs != nil && len(*client.OrchestratorEAs) > 0 {
		if record.ExtensibleAttributes == nil {
			record.ExtensibleAttributes = &infoblox.ExtensibleAttribute{}
		}
		for k, v := range *client.OrchestratorEAs {
			(*record.ExtensibleAttributes)[k] = v
		}
	}

	return &record, nil
}

func resourceARecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics
	ref := d.Id()

	record, err := client.GetARecordByRef(ref, nil)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	check := convertARecordToResourceData(client, d, &record)
	if check.HasError() {
		return check
	}

	d.SetId(record.Ref)

	return diags
}

func resourceARecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics

	record, err := convertResourceDataToARecord(client, d)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	err = client.CreateARecord(record)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if diags.HasError() {
		return diags
	}

	d.SetId(record.Ref)
	return resourceARecordRead(ctx, d, m)
}

func resourceARecordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	client := m.(*infoblox.Client)

	var record infoblox.ARecord

	if d.HasChange("hostname") {
		record.Hostname = d.Get("hostname").(string)
	}
	if d.HasChange("dns_name") {
		record.DNSName = d.Get("dns_name").(string)
	}
	if d.HasChange("ip_address") {
		record.IPAddress = d.Get("ip_address").(string)
	}
	if d.HasChange("comment") {
		record.Comment = d.Get("comment").(string)
	}
	if d.HasChange("disable") {
		record.Disable = newBool(d.Get("disable").(bool))
	}
	if d.HasChange("view") {
		record.View = d.Get("view").(string)
	}
	if d.HasChange("zone") {
		record.Zone = d.Get("zone").(string)
	}
	if d.HasChange("extensible_attributes") {
		old, new := d.GetChange("extensible_attributes")
		oldKeys := Keys(old.(map[string]interface{}))
		oldEAs, err := createExtensibleAttributesFromJSON(old.(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		newKeys := Keys(new.(map[string]interface{}))
		newEAs, err := createExtensibleAttributesFromJSON(new.(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		removeEAs := sliceDiff(oldKeys, newKeys, false)
		if len(removeEAs) > 0 {
			record.ExtensibleAttributesRemove = &infoblox.ExtensibleAttribute{}
			for _, v := range removeEAs {
				(*record.ExtensibleAttributesRemove)[v] = oldEAs[v]
			}
		}
		for k, v := range newEAs {
			if !Contains(oldKeys, k) || (Contains(oldKeys, k) && v.Value != oldEAs[k].Value) {
				if record.ExtensibleAttributesAdd == nil {
					record.ExtensibleAttributesAdd = &infoblox.ExtensibleAttribute{}
				}
				(*record.ExtensibleAttributesAdd)[k] = v
			}
		}
		if client.OrchestratorEAs != nil && len(*client.OrchestratorEAs) > 0 {
			if record.ExtensibleAttributesAdd == nil {
				record.ExtensibleAttributesAdd = &infoblox.ExtensibleAttribute{}
			}
			for k, v := range *client.OrchestratorEAs {
				(*record.ExtensibleAttributesAdd)[k] = v
			}
		}
	}
	changedRecord, err := client.UpdateARecord(d.Id(), record)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	d.SetId(changedRecord.Ref)
	return resourceARecordRead(ctx, d, m)
}

func resourceARecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics
	ref := d.Id()

	err := client.DeleteARecord(ref)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
