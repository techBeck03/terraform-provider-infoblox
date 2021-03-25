package infoblox

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

func resourceAliasRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAliasRecordCreate,
		ReadContext:   resourceAliasRecordRead,
		UpdateContext: resourceAliasRecordUpdate,
		DeleteContext: resourceAliasRecordDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
		CustomizeDiff: customdiff.Sequence(
			makeEACustomDiff("extensible_attributes"),
		),
		Schema: map[string]*schema.Schema{
			"ref": {
				Type:        schema.TypeString,
				Description: "Reference id of alias record object.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name for an Alias record in FQDN format.",
				Required:    true,
			},
			"target_name": {
				Type:        schema.TypeString,
				Description: "Target name in FQDN format.",
				Required:    true,
			},
			"target_type": {
				Type:             schema.TypeString,
				Description:      "Target type.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"A", "AAA", "MX", "NAPTR", "PTR", "SPF", "SRV", "TXT"}, true)),
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
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
			"view": {
				Type:        schema.TypeString,
				Description: "The name of the DNS View in which the record resides.",
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"zone": {
				Type:        schema.TypeString,
				Description: "The name of the zone in which the record resides.",
				Computed:    true,
			},
			"extensible_attributes": {
				Type:             schema.TypeMap,
				Description:      "Extensible attributes of alias record (Values are JSON encoded).",
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validateEa,
				DiffSuppressFunc: eaSuppressDiff,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func convertAliasRecordToResourceData(client *infoblox.Client, d *schema.ResourceData, record *infoblox.AliasRecord) diag.Diagnostics {
	var diags diag.Diagnostics

	d.Set("ref", record.Ref)
	d.Set("name", record.Name)
	d.Set("target_name", record.Target)
	d.Set("target_type", record.TargetType)
	d.Set("dns_name", record.DNSName)
	d.Set("dns_target_name", record.DNSTargetName)
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

func convertResourceDataToAliasRecord(client *infoblox.Client, d *schema.ResourceData) (*infoblox.AliasRecord, error) {
	var record infoblox.AliasRecord

	record.Name = d.Get("name").(string)
	record.Target = d.Get("target_name").(string)
	record.TargetType = d.Get("target_type").(string)
	record.DNSName = d.Get("dns_name").(string)
	record.DNSTargetName = d.Get("dns_target_name").(string)
	record.Comment = d.Get("comment").(string)
	record.Disable = newBool(d.Get("disable").(bool))
	record.View = d.Get("view").(string)
	record.Zone = d.Get("zone").(string)

	eaMap := d.Get("extensible_attributes").(map[string]interface{})
	if len(eaMap) > 0 {
		eas, err := createExtensibleAttributesFromJSON(client, eaMap)
		if err != nil {
			return &record, err
		}
		record.ExtensibleAttributes = &eas
	}

	if client.OrchestratorEAs != nil && len(*client.OrchestratorEAs) > 0 {
		for k, v := range *client.OrchestratorEAs {
			(*record.ExtensibleAttributes)[k] = v
		}
	}

	return &record, nil
}

func resourceAliasRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics
	ref := d.Id()

	record, err := client.GetAliasRecordByRef(ref, nil)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	check := convertAliasRecordToResourceData(client, d, &record)
	if check.HasError() {
		return check
	}

	d.SetId(record.Ref)

	return diags
}

func resourceAliasRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics

	record, err := convertResourceDataToAliasRecord(client, d)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	err = client.CreateAliasRecord(record)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if diags.HasError() {
		return diags
	}

	d.SetId(record.Ref)
	return resourceAliasRecordRead(ctx, d, m)
}

func resourceAliasRecordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	client := m.(*infoblox.Client)

	var record infoblox.AliasRecord

	if d.HasChange("name") {
		record.Name = d.Get("name").(string)
	}
	if d.HasChange("target_name") {
		record.Target = d.Get("target_name").(string)
	}
	if d.HasChange("target_type") {
		record.TargetType = d.Get("target_type").(string)
	}
	if d.HasChange("dns_name") {
		record.DNSName = d.Get("dns_name").(string)
	}
	if d.HasChange("dns_target_name") {
		record.DNSTargetName = d.Get("dns_target_name").(string)
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
		eaMap := d.Get("extensible_attributes").(map[string]interface{})
		if len(eaMap) > 0 {
			eas, err := createExtensibleAttributesFromJSON(client, eaMap)
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
				return diags
			}
			record.ExtensibleAttributes = &eas
		}
		if client.OrchestratorEAs != nil && len(*client.OrchestratorEAs) > 0 {
			for k, v := range *client.OrchestratorEAs {
				(*record.ExtensibleAttributes)[k] = v
			}
		}
	}
	changedRecord, err := client.UpdateAliasRecord(d.Id(), record)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	d.SetId(changedRecord.Ref)
	return resourceAliasRecordRead(ctx, d, m)
}

func resourceAliasRecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics
	ref := d.Id()

	err := client.DeleteAliasRecord(ref)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
