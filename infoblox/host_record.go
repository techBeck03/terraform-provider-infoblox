package infoblox

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

func resourceHostRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHostRecordCreate,
		ReadContext:   resourceHostRecordRead,
		UpdateContext: resourceHostRecordUpdate,
		DeleteContext: resourceHostRecordDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
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
				Optional:    true,
				Computed:    true,
			},
			"network": {
				Type:             schema.TypeString,
				Description:      "Network for host record in CIDR notation",
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
			},
			"enable_dns": {
				Type:        schema.TypeBool,
				Description: "Enable for DNS",
				Optional:    true,
				Computed:    true,
			},
			"network_view": {
				Type:        schema.TypeString,
				Description: "Network view",
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
			},
			"view": {
				Type:        schema.TypeString,
				Description: "DNS view",
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"zone": {
				Type:        schema.TypeString,
				Description: "DNS zone",
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"ip_v4_addresses": {
				Type:        schema.TypeList,
				Description: "IPv4 addresses associated with host record",
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ref": {
							Type:        schema.TypeString,
							Description: "Reference id of address object",
							Computed:    true,
						},
						"ip_address": {
							Type:        schema.TypeString,
							Description: "IP address",
							Required:    true,
						},
						"hostname": {
							Type:        schema.TypeString,
							Description: "Hostname associated with IP address",
							Optional:    true,
							Computed:    true,
						},
						"network": {
							Type:        schema.TypeString,
							Description: "Network associated with IP address",
							Optional:    true,
							Computed:    true,
						},
						"mac_address": {
							Type:        schema.TypeString,
							Description: "MAC address associated with IP address",
							Optional:    true,
							Computed:    true,
						},
						"configure_for_dhcp": {
							Type:        schema.TypeBool,
							Description: "Configure IP for DHCP",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"extensible_attributes": {
				Type:             schema.TypeMap,
				Description:      "Extensible attributes of host record",
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

func convertHostRecordToResourceData(client *infoblox.Client, d *schema.ResourceData, record *infoblox.HostRecord) diag.Diagnostics {
	var diags diag.Diagnostics

	d.Set("ref", record.Ref)
	d.Set("hostname", record.Hostname)
	d.Set("comment", record.Comment)
	d.Set("enable_dns", record.EnableDNS)
	d.Set("network_view", record.NetworkView)
	d.Set("view", record.View)
	d.Set("zone", record.Zone)

	var ipAddressList []map[string]interface{}
	for _, address := range record.IPv4Addrs {
		ipAddressList = append(ipAddressList, map[string]interface{}{
			"ref":                address.Ref,
			"ip_address":         address.IPAddress,
			"hostname":           address.Host,
			"network":            address.CIDR,
			"mac_address":        address.Mac,
			"configure_for_dhcp": address.ConfigureForDHCP,
		})
	}

	d.Set("ip_v4_addresses", ipAddressList)

	eas, err := client.ConvertEAsToJSONString(*record.ExtensibleAttributes)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	} else {
		d.Set("extensible_attributes", eas)
	}

	return diags
}

func convertResourceDataToHostRecord(client *infoblox.Client, d *schema.ResourceData) (*infoblox.HostRecord, error) {
	var record infoblox.HostRecord

	record.Hostname = d.Get("hostname").(string)
	record.Comment = d.Get("comment").(string)
	network := d.Get("network").(string)
	record.EnableDNS = newBool(d.Get("enable_dns").(bool))
	record.NetworkView = d.Get("network_view").(string)
	record.View = d.Get("view").(string)
	record.Zone = d.Get("zone").(string)

	ipAddressList := d.Get("ip_v4_addresses").([]interface{})
	record.IPv4Addrs = []infoblox.IPv4Addr{}
	if len(ipAddressList) == 0 {
		if network == "" {
			return &record, fmt.Errorf("`network` is required when not specifying an IP address")
		}
		record.IPv4Addrs = append(record.IPv4Addrs, infoblox.IPv4Addr{
			IPAddress: fmt.Sprintf("func:nextavailableip:%s", network),
		})
	} else {
		for _, address := range ipAddressList {
			var ipv4Addr infoblox.IPv4Addr
			ipv4Addr.IPAddress = address.(map[string]interface{})["ip_address"].(string)
			if address.(map[string]interface{})["host"].(string) != "" {
				ipv4Addr.Host = address.(map[string]interface{})["host"].(string)
			}
			ipv4Addr.ConfigureForDHCP = newBool(address.(map[string]interface{})["configure_for_dhcp"].(bool))
			if address.(map[string]interface{})["mac_address"].(string) != "" {
				ipv4Addr.Mac = address.(map[string]interface{})["mac_address"].(string)
			}
			record.IPv4Addrs = append(record.IPv4Addrs, ipv4Addr)
		}
	}

	eaMap := d.Get("extensible_attributes").(map[string]interface{})
	if len(eaMap) > 0 {
		eas, err := createExtensibleAttributesFromJSON(client, eaMap)
		if err != nil {
			return &record, err
		}
		record.ExtensibleAttributes = &eas
	}

	return &record, nil
}

func resourceHostRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics
	ref := d.Id()

	record, err := client.GetHostRecordByRef(ref, nil)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	check := convertHostRecordToResourceData(client, d, &record)
	if check.HasError() {
		return check
	}

	d.SetId(record.Ref)

	return diags
}

func resourceHostRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics

	record, err := convertResourceDataToHostRecord(client, d)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	prettyPrint(record)
	err = client.CreateHostRecord(record)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if diags.HasError() {
		return diags
	}

	d.SetId(record.Ref)
	return resourceHostRecordRead(ctx, d, m)
}

func resourceHostRecordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	client := m.(*infoblox.Client)

	var record infoblox.HostRecord

	if d.HasChange("hostname") {
		record.Hostname = d.Get("hostname").(string)
	}
	if d.HasChange("comment") {
		record.Comment = d.Get("comment").(string)
	}
	if d.HasChange("enable_dns") {
		record.EnableDNS = newBool(d.Get("enable_dns").(bool))
	}
	if d.HasChange("network_view") {
		record.NetworkView = d.Get("network_view").(string)
	}
	if d.HasChange("view") {
		record.NetworkView = d.Get("view").(string)
	}
	if d.HasChange("zone") {
		record.NetworkView = d.Get("zone").(string)
	}
	if d.HasChange("ip_v4_addresses") {
		var ipAddressList []map[string]interface{}
		for _, address := range record.IPv4Addrs {
			ipAddressList = append(ipAddressList, map[string]interface{}{
				"ref":                address.Ref,
				"ip_address":         address.IPAddress,
				"hostname":           address.Host,
				"network":            address.CIDR,
				"mac_address":        address.Mac,
				"configure_for_dhcp": address.ConfigureForDHCP,
			})
		}

		d.Set("ip_v4_addresses", ipAddressList)
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
	}
	changedRecord, err := client.UpdateHostRecord(d.Id(), record)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	d.SetId(changedRecord.Ref)
	return resourceHostRecordRead(ctx, d, m)
}

func resourceHostRecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics
	ref := d.Id()

	err := client.DeleteHostRecord(ref)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
