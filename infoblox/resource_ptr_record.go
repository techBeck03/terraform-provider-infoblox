package infoblox

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

func resourcePtrRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePtrRecordCreate,
		ReadContext:   resourcePtrRecordRead,
		UpdateContext: resourcePtrRecordUpdate,
		DeleteContext: resourcePtrRecordDelete,
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
				Description: "The name for a DNS PTR record in punycode format.",
				Computed:    true,
			},
			"dns_pointer_domain_name": {
				Type:        schema.TypeString,
				Description: "The domain name of the DNS PTR record in punycode format.",
				Computed:    true,
			},
			"extensible_attributes": {
				Type:             schema.TypeMap,
				Description:      "Extensible attributes of ptr record (Values are JSON encoded).",
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validateEa,
				DiffSuppressFunc: eaSuppressDiff,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ip_v4_address": {
				Type:             schema.TypeString,
				Description:      "The IPv4 Address of the record.",
				Optional:         true,
				AtLeastOneOf:     []string{"ip_v4_address", "ip_v6_address"},
				ConflictsWith:    []string{"ip_v6_address"},
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv4Address),
			},
			"ip_v6_address": {
				Type:             schema.TypeString,
				Description:      "The IPv6 Address of the record.",
				Optional:         true,
				AtLeastOneOf:     []string{"ip_v4_address", "ip_v6_address"},
				ConflictsWith:    []string{"ip_v4_address"},
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv6Address),
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the DNS PTR record in FQDN format.",
				Optional:      true,
				Computed:      true,
				AtLeastOneOf:  []string{"name", "pointer_domain_name"},
				ConflictsWith: []string{"pointer_domain_name"},
			},
			"pointer_domain_name": {
				Type:          schema.TypeString,
				Description:   "The domain name of the DNS PTR record in FQDN format.",
				Optional:      true,
				Computed:      true,
				AtLeastOneOf:  []string{"name", "pointer_domain_name"},
				ConflictsWith: []string{"name"},
			},
			"ref": {
				Type:        schema.TypeString,
				Description: "Reference id of ptr record object.",
				Computed:    true,
			},
			"view": {
				Type:        schema.TypeString,
				Description: "Name of the DNS View in which the record resides.",
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

func convertPtrRecordToResourceData(client *infoblox.Client, d *schema.ResourceData, record *infoblox.PtrRecord) diag.Diagnostics {
	var diags diag.Diagnostics

	d.Set("ref", record.Ref)
	d.Set("name", record.Name)
	d.Set("pointer_domain_name", record.PointerDomainName)
	d.Set("ip_v4_address", record.IPv4Address)
	d.Set("ip_v6_address", record.IPv6Address)
	d.Set("dns_name", record.DNSName)
	d.Set("dns_pointer_domain_name", record.DNSPointerDomainName)
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

func convertResourceDataToPtrRecord(client *infoblox.Client, d *schema.ResourceData) (*infoblox.PtrRecord, error) {
	var record infoblox.PtrRecord

	record.Name = d.Get("name").(string)
	record.PointerDomainName = d.Get("pointer_domain_name").(string)
	record.IPv4Address = d.Get("ip_v4_address").(string)
	record.IPv6Address = d.Get("ip_v6_address").(string)
	record.DNSName = d.Get("dns_name").(string)
	record.DNSPointerDomainName = d.Get("dns_pointer_domain_name").(string)
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

func resourcePtrRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics
	ref := d.Id()

	record, err := client.GetPtrRecordByRef(ref, nil)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	check := convertPtrRecordToResourceData(client, d, &record)
	if check.HasError() {
		return check
	}

	d.SetId(record.Ref)

	return diags
}

func resourcePtrRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics

	record, err := convertResourceDataToPtrRecord(client, d)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	err = client.CreatePtrRecord(record)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if diags.HasError() {
		return diags
	}

	d.SetId(record.Ref)
	return resourcePtrRecordRead(ctx, d, m)
}

func resourcePtrRecordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	client := m.(*infoblox.Client)

	var record infoblox.PtrRecord

	if d.HasChange("name") {
		record.Name = d.Get("name").(string)
	}
	if d.HasChange("pointer_domain_name") {
		record.PointerDomainName = d.Get("pointer_domain_name").(string)
	}
	if d.HasChange("ip_v4_address") {
		record.IPv4Address = d.Get("ip_v4_address").(string)
	}
	if d.HasChange("ip_v6_address") {
		record.IPv6Address = d.Get("ip_v6_address").(string)
	}
	if d.HasChange("dns_name") {
		record.DNSName = d.Get("dns_name").(string)
	}
	if d.HasChange("dns_pointer_domain_name") {
		record.DNSPointerDomainName = d.Get("dns_pointer_domain_name").(string)
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
	changedRecord, err := client.UpdatePtrRecord(d.Id(), record)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	d.SetId(changedRecord.Ref)
	return resourcePtrRecordRead(ctx, d, m)
}

func resourcePtrRecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics
	ref := d.Id()

	err := client.DeletePtrRecord(ref)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
