package infoblox

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

var (
	cNameRecordRequiredIPFields = []string{
		"network",
		"ip_v4_address",
		"range_function_string",
	}
)

func resourceCNameRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCNameRecordCreate,
		ReadContext:   resourceCNameRecordRead,
		UpdateContext: resourceCNameRecordUpdate,
		DeleteContext: resourceCNameRecordDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
		CustomizeDiff: customdiff.Sequence(
			makeEACustomDiff("extensible_attributes"),
		),
		Schema: map[string]*schema.Schema{
			"alias": {
				Type:        schema.TypeString,
				Description: "The name for a CNAME record in FQDN format.",
				Required:    true,
			},
			"canonical": {
				Type:        schema.TypeString,
				Description: "Canonical name in FQDN format.",
				Required:    true,
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
			"dns_canonical": {
				Type:        schema.TypeString,
				Description: "Canonical name in punycode format.",
				Computed:    true,
			},
			"dns_name": {
				Type:        schema.TypeString,
				Description: "The name for the CNAME record in punycode format.",
				Computed:    true,
			},
			"extensible_attributes": {
				Type:             schema.TypeMap,
				Description:      "Extensible attributes of cname record (Values are JSON encoded).",
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validateEa,
				DiffSuppressFunc: eaSuppressDiff,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ref": {
				Type:        schema.TypeString,
				Description: "Reference id of cname record object.",
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

func convertCNameRecordToResourceData(client *infoblox.Client, d *schema.ResourceData, record *infoblox.CNameRecord) diag.Diagnostics {
	var diags diag.Diagnostics

	d.Set("ref", record.Ref)
	d.Set("alias", record.Alias)
	d.Set("canonical", record.Canonical)
	d.Set("dns_name", record.DNSName)
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

func convertResourceDataToCNameRecord(client *infoblox.Client, d *schema.ResourceData) (*infoblox.CNameRecord, error) {
	var record infoblox.CNameRecord

	record.Alias = d.Get("alias").(string)
	record.Canonical = d.Get("canonical").(string)
	record.DNSName = d.Get("dns_name").(string)
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

func resourceCNameRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics
	ref := d.Id()

	record, err := client.GetCNameRecordByRef(ref, nil)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	check := convertCNameRecordToResourceData(client, d, &record)
	if check.HasError() {
		return check
	}

	d.SetId(record.Ref)

	return diags
}

func resourceCNameRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics

	record, err := convertResourceDataToCNameRecord(client, d)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	err = client.CreateCNameRecord(record)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if diags.HasError() {
		return diags
	}

	d.SetId(record.Ref)
	return resourceCNameRecordRead(ctx, d, m)
}

func resourceCNameRecordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	client := m.(*infoblox.Client)

	var record infoblox.CNameRecord

	if d.HasChange("alias") {
		record.Alias = d.Get("alias").(string)
	}
	if d.HasChange("canonical") {
		record.Canonical = d.Get("canonical").(string)
	}
	if d.HasChange("dns_name") {
		record.DNSName = d.Get("dns_name").(string)
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
	if extensibleAttributes, ok := d.GetOk("extensible_attributes"); ok {
		eaMap := extensibleAttributes.(map[string]interface{})
		if len(eaMap) > 0 {
			eas, err := createExtensibleAttributesFromJSON(eaMap)
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
				return diags
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
	}
	changedRecord, err := client.UpdateCNameRecord(d.Id(), record)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	d.SetId(changedRecord.Ref)
	return resourceCNameRecordRead(ctx, d, m)
}

func resourceCNameRecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics
	ref := d.Id()

	err := client.DeleteCNameRecord(ref)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
