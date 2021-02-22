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

var (
	networkRequiredIPFields = []string{
		"network",
		"ip_v4_address",
	}
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		// CreateContext: resourceNetworkCreate,
		ReadContext: resourceNetworkRead,
		// UpdateContext: resourceNetworkUpdate,
		// DeleteContext: resourceNetworkDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
		CustomizeDiff: customdiff.Sequence(
			eaCustomDiff,
		),
		Schema: map[string]*schema.Schema{
			"ref": {
				Type:          schema.TypeString,
				Description:   "Reference id of network object",
				Computed:      true,
				ConflictsWith: []string{"hostname"},
				AtLeastOneOf:  dataNetworkRequiredSearchFields,
			},
			"network": {
				Type:          schema.TypeString,
				Description:   "CIDR of network",
				Required:      true,
				ConflictsWith: []string{"ref"},
				AtLeastOneOf:  dataNetworkRequiredSearchFields,
			},
			"comment": {
				Type:        schema.TypeString,
				Description: "Comment string",
				Optional:    true,
				Computed:    true,
			},
			"disable": {
				Type:        schema.TypeBool,
				Description: "Disable for DHCP",
				Optional:    true,
				Computed:    true,
			},
			"network_view": {
				Type:        schema.TypeString,
				Description: "Network view",
				Optional:    true,
				Computed:    true,
			},
			"member": {
				Type:        schema.TypeSet,
				Description: "Grid members associated with network",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"struct": {
							Type:             schema.TypeString,
							Description:      "Struct type of member",
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"DHCPMEMBER", "MSDHCPSERVER"}, true)),
							StateFunc: func(val interface{}) string {
								return strings.ToUpper(val.(string))
							},
						},
						"ip_v4_address": {
							Type:             schema.TypeString,
							Description:      "IPv4 address",
							Optional:         true,
							Computed:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv4Address),
						},
						"ip_v6_address": {
							Type:             schema.TypeString,
							Description:      "IPv6 address",
							Optional:         true,
							Computed:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv6Address),
						},
						"hostname": {
							Type:        schema.TypeString,
							Description: "Hostname of member",
							Required:    true,
						},
					},
				},
			},
			"option": {
				Type:        schema.TypeSet,
				Description: "DHCP options associated with network",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Name of DHCP option",
							Required:    true,
						},
						"code": {
							Type:        schema.TypeInt,
							Description: "Code of the DHCP option",
							Computed:    true,
						},
						"use_option": {
							Type:        schema.TypeBool,
							Description: "Use this dhcp option",
							Required:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "Value of option",
							Required:    true,
						},
						"vendor_class": {
							Type:        schema.TypeString,
							Description: "Value of option",
							Default:     "DHCP",
						},
					},
				},
			},
			"extensible_attributes": {
				Type:        schema.TypeMap,
				Description: "Extensible attributes of network",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func convertNetworkToResourceData(client *infoblox.Client, d *schema.ResourceData, network *infoblox.Network) diag.Diagnostics {
	var diags diag.Diagnostics

	d.Set("ref", network.Ref)
	d.Set("network", network.CIDR)
	d.Set("comment", network.Comment)
	d.Set("disable", network.DisableDHCP)
	d.Set("network_view", network.NetworkView)

	var memberList []map[string]interface{}
	for _, member := range network.Members {
		memberList = append(memberList, map[string]interface{}{
			"struct":        member.StructType,
			"hostname":      member.Hostname,
			"ip_v4_address": member.IPV4Address,
			"ip_v6_address": member.IPV6Address,
		})
	}

	d.Set("member", memberList)

	var optionList []map[string]interface{}
	for _, option := range network.Options {
		optionList = append(optionList, map[string]interface{}{
			"name":         option.Name,
			"code":         option.Code,
			"use_option":   option.UseOption,
			"value":        option.Value,
			"vendor_class": option.VendorClass,
		})
	}

	d.Set("options", optionList)

	eas, err := client.ConvertEAsToJSONString(*network.ExtensibleAttributes)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	} else {
		d.Set("extensible_attributes", eas)
	}

	return diags
}

// func convertResourceDataToNetwork(client *infoblox.Client, d *schema.ResourceData) (*infoblox.Network, error) {
// 	var network infoblox.Network

// 	network.Hostname = d.Get("hostname").(string)
// 	network.Comment = d.Get("comment").(string)
// 	network := d.Get("network").(string)
// 	network.EnableDNS = newBool(d.Get("enable_dns").(bool))
// 	network.NetworkView = d.Get("network_view").(string)
// 	network.View = d.Get("view").(string)
// 	network.Zone = d.Get("zone").(string)

// 	ipAddressList := d.Get("ip_v4_address").(*schema.Set).List()
// 	network.IPv4Addrs = []infoblox.IPv4Addr{}
// 	if len(ipAddressList) == 0 {
// 		network.IPv4Addrs = append(network.IPv4Addrs, infoblox.IPv4Addr{
// 			IPAddress: fmt.Sprintf("func:nextavailableip:%s", network),
// 		})
// 	} else {
// 		for _, address := range ipAddressList {
// 			var ipv4Addr infoblox.IPv4Addr
// 			ipv4Addr.IPAddress = address.(map[string]interface{})["ip_address"].(string)
// 			if address.(map[string]interface{})["hostname"] != "" {
// 				ipv4Addr.Host = address.(map[string]interface{})["hostname"].(string)
// 			}
// 			ipv4Addr.ConfigureForDHCP = newBool(address.(map[string]interface{})["configure_for_dhcp"].(bool))
// 			if address.(map[string]interface{})["mac_address"].(string) != "" {
// 				ipv4Addr.Mac = address.(map[string]interface{})["mac_address"].(string)
// 			}
// 			network.IPv4Addrs = append(network.IPv4Addrs, ipv4Addr)
// 		}
// 	}

// 	eaMap := d.Get("extensible_attributes").(map[string]interface{})
// 	if len(eaMap) > 0 {
// 		eas, err := createExtensibleAttributesFromJSON(client, eaMap)
// 		if err != nil {
// 			return &network, err
// 		}
// 		network.ExtensibleAttributes = &eas
// 	}

// 	if client.OrchestratorEAs != nil && len(*client.OrchestratorEAs) > 0 {
// 		for k, v := range *client.OrchestratorEAs {
// 			(*network.ExtensibleAttributes)[k] = v
// 		}
// 	}

// 	return &network, nil
// }

func resourceNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics
	ref := d.Id()

	network, err := client.GetNetworkByRef(ref, nil)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	check := convertNetworkToResourceData(client, d, &network)
	if check.HasError() {
		return check
	}

	d.SetId(network.Ref)

	return diags
}

// func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	client := m.(*infoblox.Client)

// 	var diags diag.Diagnostics

// 	network, err := convertResourceDataToNetwork(client, d)
// 	if err != nil {
// 		diags = append(diags, diag.FromErr(err)...)
// 		return diags
// 	}

// 	err = client.CreateNetwork(network)
// 	if err != nil {
// 		diags = append(diags, diag.FromErr(err)...)
// 		return diags
// 	}

// 	if diags.HasError() {
// 		return diags
// 	}

// 	d.SetId(network.Ref)
// 	return resourceNetworkRead(ctx, d, m)
// }

// func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
// 	client := m.(*infoblox.Client)

// 	var network infoblox.Network

// 	if d.HasChange("hostname") {
// 		network.Hostname = d.Get("hostname").(string)
// 	}
// 	if d.HasChange("comment") {
// 		network.Comment = d.Get("comment").(string)
// 	}
// 	if d.HasChange("enable_dns") {
// 		network.EnableDNS = newBool(d.Get("enable_dns").(bool))
// 	}
// 	if d.HasChange("network") {
// 		network := d.Get("network").(string)
// 		network.IPv4Addrs = append(network.IPv4Addrs, infoblox.IPv4Addr{
// 			IPAddress: fmt.Sprintf("func:nextavailableip:%s", network),
// 		})
// 	}
// 	if d.HasChange("network_view") {
// 		network.NetworkView = d.Get("network_view").(string)
// 	}
// 	if d.HasChange("view") {
// 		network.NetworkView = d.Get("view").(string)
// 	}
// 	if d.HasChange("zone") {
// 		network.NetworkView = d.Get("zone").(string)
// 	}
// 	if d.HasChange("ip_v4_address") {
// 		ipAddressList := d.Get("ip_v4_address").(*schema.Set).List()
// 		network.IPv4Addrs = []infoblox.IPv4Addr{}
// 		for _, address := range ipAddressList {
// 			var ipv4Addr infoblox.IPv4Addr

// 			ipv4Addr.IPAddress = address.(map[string]interface{})["ip_address"].(string)
// 			if address.(map[string]interface{})["hostname"] != "" {
// 				ipv4Addr.Host = address.(map[string]interface{})["hostname"].(string)
// 			}

// 			ipv4Addr.ConfigureForDHCP = newBool(address.(map[string]interface{})["configure_for_dhcp"].(bool))
// 			if address.(map[string]interface{})["mac_address"].(string) != "" {
// 				ipv4Addr.Mac = address.(map[string]interface{})["mac_address"].(string)
// 			}
// 			network.IPv4Addrs = append(network.IPv4Addrs, ipv4Addr)
// 		}
// 	}
// 	if d.HasChange("extensible_attributes") {
// 		eaMap := d.Get("extensible_attributes").(map[string]interface{})
// 		if len(eaMap) > 0 {
// 			eas, err := createExtensibleAttributesFromJSON(client, eaMap)
// 			if err != nil {
// 				diags = append(diags, diag.FromErr(err)...)
// 				return diags
// 			}
// 			network.ExtensibleAttributes = &eas
// 		}
// 		if client.OrchestratorEAs != nil && len(*client.OrchestratorEAs) > 0 {
// 			for k, v := range *client.OrchestratorEAs {
// 				(*network.ExtensibleAttributes)[k] = v
// 			}
// 		}
// 	}
// 	changedRecord, err := client.UpdateNetwork(d.Id(), network)
// 	if err != nil {
// 		diags = append(diags, diag.FromErr(err)...)
// 		return diags
// 	}

// 	d.SetId(changedRecord.Ref)
// 	return resourceNetworkRead(ctx, d, m)
// }

// func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	client := m.(*infoblox.Client)

// 	var diags diag.Diagnostics
// 	ref := d.Id()

// 	err := client.DeleteNetwork(ref)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}

// 	return diags
// }
