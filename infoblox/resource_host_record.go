package infoblox

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

var (
	hostRecordRequiredIPFields = []string{
		"ip_v4_address.*.network",
		"ip_v4_address.*.ip_address",
		"ip_v4_address.*.range_function_string",
	}
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
		CustomizeDiff: customdiff.Sequence(
			makeEACustomDiff("extensible_attributes"),
			hostRecordAddressDiff,
		),
		Schema: map[string]*schema.Schema{
			"comment": {
				Type:             schema.TypeString,
				Description:      "Comment for the record; maximum 256 characters.",
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 256)),
			},
			"enable_dns": {
				Type:        schema.TypeBool,
				Description: "When false, the host does not have parent zone information.",
				Optional:    true,
				Computed:    true,
			},
			"extensible_attributes": {
				Type:             schema.TypeMap,
				Description:      "Extensible attributes of host record (Values are JSON encoded).",
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
				Description: "The host name in FQDN format.",
				Required:    true,
			},
			"ip_v4_address": {
				Type:        schema.TypeSet,
				Description: "IPv4 addresses associated with host record.",
				Optional:    true,
				Computed:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"configure_for_dhcp": {
							Type:        schema.TypeBool,
							Description: "Set this to True to enable the DHCP configuration for this host address.",
							Optional:    true,
							Computed:    true,
						},
						"hostname": {
							Type:        schema.TypeString,
							Description: "Hostname associated with IP address.",
							Computed:    true,
						},
						"ip_address": {
							Type:             schema.TypeString,
							Description:      "IP address.",
							Optional:         true,
							ForceNew:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv4Address),
						},
						"mac_address": {
							Type:        schema.TypeString,
							Description: "MAC address associated with IP address.",
							Optional:    true,
							Computed:    true,
						},
						"network": {
							Type:             schema.TypeString,
							Description:      "Network for host record in CIDR notation (next_available_ip will be retrieved from this network).",
							Optional:         true,
							Computed:         true,
							ForceNew:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
						},
						"range_function_string": {
							Type:        schema.TypeString,
							Description: "Range start and end string for next_available_ip function calls.",
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
						},
						"ref": {
							Type:        schema.TypeString,
							Description: "Reference id of address object.",
							Computed:    true,
						},
						"use_for_ea_inheritance": {
							Type:        schema.TypeBool,
							Description: "Set this to True when using this host address for EA inheritance.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"network_view": {
				Type:        schema.TypeString,
				Description: "The name of the network view in which the host record resides.",
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
			},
			"ref": {
				Type:        schema.TypeString,
				Description: "Reference id of host record object.",
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

func convertHostRecordToResourceData(client *infoblox.Client, d *schema.ResourceData, record *infoblox.HostRecord) diag.Diagnostics {
	var diags diag.Diagnostics

	d.Set("ref", record.Ref)
	d.Set("hostname", record.Hostname)
	d.Set("comment", record.Comment)
	d.Set("enable_dns", record.EnableDNS)
	d.Set("network_view", record.NetworkView)
	d.Set("view", record.View)
	d.Set("zone", record.Zone)

	configuredAddressList := d.Get("ip_v4_address").(*schema.Set).List()
	var ipAddressList []map[string]interface{}
	for _, address := range record.IPv4Addrs {
		var configuredIndex int
		for i, v := range configuredAddressList {
			addr := v.(map[string]interface{})["ip_address"].(string)
			if addr == address.IPAddress {
				configuredIndex = i
			}
		}
		ipAddressList = append(ipAddressList, map[string]interface{}{
			"ref":                    address.Ref,
			"ip_address":             address.IPAddress,
			"hostname":               address.Host,
			"network":                configuredAddressList[configuredIndex].(map[string]interface{})["network"].(string),
			"range_function_string":  configuredAddressList[configuredIndex].(map[string]interface{})["range_function_string"].(string),
			"mac_address":            address.Mac,
			"configure_for_dhcp":     address.ConfigureForDHCP,
			"use_for_ea_inheritance": address.UseForEAInheritance,
		})
	}

	d.Set("ip_v4_address", ipAddressList)

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
	record.EnableDNS = newBool(d.Get("enable_dns").(bool))
	record.NetworkView = d.Get("network_view").(string)
	record.View = d.Get("view").(string)
	record.Zone = d.Get("zone").(string)

	ipAddressList := d.Get("ip_v4_address").(*schema.Set).List()
	record.IPv4Addrs = []infoblox.IPv4Addr{}
	for _, address := range ipAddressList {
		var ipv4Addr infoblox.IPv4Addr
		if address.(map[string]interface{})["ip_address"].(string) != "" {
			ipv4Addr.IPAddress = address.(map[string]interface{})["ip_address"].(string)
		} else if address.(map[string]interface{})["network"].(string) != "" {
			ipv4Addr.IPAddress = fmt.Sprintf("func:nextavailableip:%s", address.(map[string]interface{})["network"].(string))
		} else if address.(map[string]interface{})["range_function_string"].(string) != "" {
			ipv4Addr.IPAddress = fmt.Sprintf("func:nextavailableip:%s", address.(map[string]interface{})["range_function_string"].(string))
		}
		if address.(map[string]interface{})["hostname"] != "" {
			ipv4Addr.Host = address.(map[string]interface{})["hostname"].(string)
		}
		ipv4Addr.ConfigureForDHCP = newBool(address.(map[string]interface{})["configure_for_dhcp"].(bool))
		ipv4Addr.UseForEAInheritance = newBool(address.(map[string]interface{})["use_for_ea_inheritance"].(bool))
		if address.(map[string]interface{})["mac_address"].(string) != "" {
			ipv4Addr.Mac = address.(map[string]interface{})["mac_address"].(string)
		}
		record.IPv4Addrs = append(record.IPv4Addrs, ipv4Addr)
	}

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
	if d.HasChange("network") {
		if network, ok := d.GetOk("network"); ok && network.(string) != "" {
			record.IPv4Addrs = append(record.IPv4Addrs, infoblox.IPv4Addr{
				IPAddress: fmt.Sprintf("func:nextavailableip:%s", network.(string)),
			})
		}
	}
	if d.HasChange("range_function_string") {
		if rangeFunctionString, ok := d.GetOk("range_function_string"); ok && rangeFunctionString.(string) != "" {
			record.IPv4Addrs = append(record.IPv4Addrs, infoblox.IPv4Addr{
				IPAddress: fmt.Sprintf("func:nextavailableip:%s", rangeFunctionString.(string)),
			})
		}
	}
	if d.HasChange("network_view") {
		record.NetworkView = d.Get("network_view").(string)
	}
	if d.HasChange("view") {
		record.View = d.Get("view").(string)
	}
	if d.HasChange("zone") {
		record.Zone = d.Get("zone").(string)
	}
	if d.HasChange("ip_v4_address") {
		if ipAddress, ok := d.GetOk("ip_v4_address"); ok && len(ipAddress.(*schema.Set).List()) > 0 {
			record.IPv4Addrs = []infoblox.IPv4Addr{}
			for _, address := range ipAddress.(*schema.Set).List() {
				var ipv4Addr infoblox.IPv4Addr

				ipv4Addr.IPAddress = address.(map[string]interface{})["ip_address"].(string)
				if address.(map[string]interface{})["hostname"] != "" {
					ipv4Addr.Host = address.(map[string]interface{})["hostname"].(string)
				}

				ipv4Addr.ConfigureForDHCP = newBool(address.(map[string]interface{})["configure_for_dhcp"].(bool))
				ipv4Addr.UseForEAInheritance = newBool(address.(map[string]interface{})["use_for_ea_inheritance"].(bool))
				if address.(map[string]interface{})["mac_address"].(string) != "" {
					ipv4Addr.Mac = address.(map[string]interface{})["mac_address"].(string)
				}
				record.IPv4Addrs = append(record.IPv4Addrs, ipv4Addr)
			}
		}
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
