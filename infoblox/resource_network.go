package infoblox

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkCreate,
		ReadContext:   resourceNetworkRead,
		UpdateContext: resourceNetworkUpdate,
		DeleteContext: resourceNetworkDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
		CustomizeDiff: customdiff.Sequence(
			makeEACustomDiff("extensible_attributes"),
			makeEACustomDiff("gateway_extensible_attributes"),
			optionCustomDiff,
		),
		Schema: map[string]*schema.Schema{
			"ref": {
				Type:        schema.TypeString,
				Description: "Reference id of network object.",
				Computed:    true,
			},
			"cidr": {
				Type:             schema.TypeString,
				Description:      "The network address in IPv4 Address/CIDR format.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
			},
			"gateway_ip": {
				Type:             schema.TypeString,
				Description:      "Default gateway IPv4 address.",
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPAddress),
			},
			"gateway_comment": {
				Type:             schema.TypeString,
				Description:      "Comment for gateway reservation.",
				Optional:         true,
				Default:          "Gateway",
				AtLeastOneOf:     []string{"gateway_ip"},
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPAddress),
			},
			"gateway_extensible_attributes": {
				Type:             schema.TypeMap,
				Description:      "Extensible attributes for gateway fixed reservation.",
				Optional:         true,
				Computed:         true,
				AtLeastOneOf:     []string{"gateway_ip"},
				ValidateDiagFunc: validateEa,
				DiffSuppressFunc: eaSuppressDiff,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"comment": {
				Type:             schema.TypeString,
				Description:      "Comment for the record; maximum 256 characters.",
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 256)),
			},
			"disable_dhcp": {
				Type:        schema.TypeBool,
				Description: "Disable for DHCP.",
				Optional:    true,
				Computed:    true,
			},
			"grid_ref": {
				Type:         schema.TypeString,
				Description:  "Ref for grid needed for restarting services.",
				Optional:     true,
				RequiredWith: []string{"restart_if_needed"},
			},
			"restart_if_needed": {
				Type:        schema.TypeBool,
				Description: "Restart dhcp services if needed.",
				Optional:    true,
			},
			"network_view": {
				Type:        schema.TypeString,
				Description: "The name of the network view in which this network resides.",
				Optional:    true,
				Computed:    true,
			},
			"member": {
				Type:        schema.TypeList,
				Description: "Grid members associated with network.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"struct": {
							Type:             schema.TypeString,
							Description:      "Struct type of member.",
							Optional:         true,
							Default:          "dhcpmember",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"dhcpmember", "msdhcpserver"}, true)),
							StateFunc: func(val interface{}) string {
								return strings.ToLower(val.(string))
							},
						},
						"ip_v4_address": {
							Type:             schema.TypeString,
							Description:      "IPv4 address.",
							Optional:         true,
							Computed:         true,
							ConflictsWith:    []string{"member.0.ip_v6_address"},
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv4Address),
						},
						"ip_v6_address": {
							Type:             schema.TypeString,
							Description:      "IPv6 address.",
							Optional:         true,
							Computed:         true,
							ConflictsWith:    []string{"member.0.ip_v4_address"},
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv6Address),
						},
						"hostname": {
							Type:        schema.TypeString,
							Description: "Hostname of member.",
							Required:    true,
						},
					},
				},
			},
			"option": {
				Type:        schema.TypeSet,
				Description: "An array of DHCP option structs that lists the DHCP options associated with the object.",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the DHCP option.",
							Required:    true,
						},
						"code": {
							Type:        schema.TypeInt,
							Description: "The code of the DHCP option.",
							Optional:    true,
							Computed:    true,
						},
						"use_option": {
							Type:        schema.TypeBool,
							Description: "Only applies to special options that are displayed separately from other options and have a use flag.",
							Optional:    true,
							Default:     true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "Value of the DHCP option.",
							Required:    true,
						},
						"vendor_class": {
							Type:        schema.TypeString,
							Description: "The name of the space this DHCP option is associated to.",
							Optional:    true,
							Default:     "DHCP",
						},
					},
				},
			},
			"extensible_attributes": {
				Type:             schema.TypeMap,
				Description:      "Extensible attributes of network (Values are JSON encoded).",
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

func convertNetworkToResourceData(client *infoblox.Client, d *schema.ResourceData, network *infoblox.Network) diag.Diagnostics {
	var diags diag.Diagnostics

	d.Set("ref", network.Ref)
	d.Set("cidr", network.CIDR)
	d.Set("comment", network.Comment)
	d.Set("disable_dhcp", network.DisableDHCP)
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

	d.Set("option", optionList)

	_, ok := d.GetOk("gateway_extensible_attributes")
	if ok {
		existingFixedAddress, err := client.GetFixedAddressByQuery(map[string]string{
			"network":  d.Get("cidr").(string),
			"ipv4addr": d.Get("gateway_ip").(string),
		})
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
		} else {
			if len(existingFixedAddress) != 1 {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Error retrieving gateway fixed address",
					Detail:   fmt.Sprintf("Multiple fixed addresses found for gateway_ip: %s", d.Get("gateway_ip").(string)),
				})
			} else {
				gwEAs, err := client.ConvertEAsToJSONString(*existingFixedAddress[0].ExtensibleAttributes)
				if err != nil {
					diags = append(diags, diag.FromErr(err)...)
				} else {
					d.Set("gateway_extensible_attributes", gwEAs)
				}
			}
		}

	}

	eas, err := client.ConvertEAsToJSONString(*network.ExtensibleAttributes)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	} else {
		d.Set("extensible_attributes", eas)
	}

	return diags
}

func convertResourceDataToNetwork(client *infoblox.Client, d *schema.ResourceData) (*infoblox.Network, error) {
	var network infoblox.Network

	network.CIDR = d.Get("cidr").(string)
	network.Comment = d.Get("comment").(string)
	network.DisableDHCP = newBool(d.Get("disable_dhcp").(bool))
	network.NetworkView = d.Get("network_view").(string)

	memberList := d.Get("member").([]interface{})
	network.Members = []infoblox.Member{}
	if len(memberList) > 0 {
		for _, member := range memberList {
			network.Members = append(network.Members, infoblox.Member{
				StructType:  member.(map[string]interface{})["struct"].(string),
				Hostname:    member.(map[string]interface{})["hostname"].(string),
				IPV4Address: member.(map[string]interface{})["ip_v4_address"].(string),
				IPV6Address: member.(map[string]interface{})["ip_v6_address"].(string),
			})
		}
	}

	optionList := d.Get("option").(*schema.Set).List()
	network.Options = []infoblox.Option{}
	if len(optionList) > 0 {
		for _, option := range optionList {
			network.Options = append(network.Options, infoblox.Option{
				Name:        option.(map[string]interface{})["name"].(string),
				UseOption:   newBool(option.(map[string]interface{})["use_option"].(bool)),
				Value:       option.(map[string]interface{})["value"].(string),
				VendorClass: option.(map[string]interface{})["vendor_class"].(string),
			})
		}
	}

	eaMap := d.Get("extensible_attributes").(map[string]interface{})
	if len(eaMap) > 0 {
		eas, err := createExtensibleAttributesFromJSON(client, eaMap)
		if err != nil {
			return &network, err
		}
		network.ExtensibleAttributes = &eas
	}

	if client.OrchestratorEAs != nil && len(*client.OrchestratorEAs) > 0 {
		for k, v := range *client.OrchestratorEAs {
			(*network.ExtensibleAttributes)[k] = v
		}
	}

	return &network, nil
}

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

func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics

	gw, ok := d.GetOk("gateway_ip")
	if ok {
		_, parsedNetwork, _ := net.ParseCIDR(d.Get("cidr").(string))
		if parsedNetwork.Contains(net.ParseIP(gw.(string))) != true {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid gateway IP address",
				Detail:   fmt.Sprintf("Gateway address: %s is not within network CIDR: %s", gw.(string), d.Get("cidr").(string)),
			})
			return diags
		}
	}

	network, err := convertResourceDataToNetwork(client, d)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	err = client.CreateNetwork(network)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if diags.HasError() {
		return diags
	}

	if ok {
		gwFixedAddress := infoblox.FixedAddress{
			IPAddress:   gw.(string),
			Comment:     d.Get("gateway_comment").(string),
			MatchClient: "RESERVED",
		}
		eaMap := d.Get("gateway_extensible_attributes").(map[string]interface{})
		if len(eaMap) > 0 {
			eas, err := createExtensibleAttributesFromJSON(client, eaMap)
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
			gwFixedAddress.ExtensibleAttributes = &eas
		}

		if client.OrchestratorEAs != nil && len(*client.OrchestratorEAs) > 0 {
			for k, v := range *client.OrchestratorEAs {
				(*gwFixedAddress.ExtensibleAttributes)[k] = v
			}
		}
		err = client.CreateFixedAddress(&gwFixedAddress)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	if d.Get("restart_if_needed").(bool) && len(network.Members) == 1 {
		err = client.RestartServices(d.Get("grid_ref").(string), infoblox.GridServiceRestartRequest{
			RestartOption: "RESTART_IF_NEEDED",
			Services:      []string{"DHCP"},
			Members:       []string{network.Members[0].Hostname},
		})
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
		// time.Sleep(2 * time.Second)
	}

	if diags.HasError() {
		return diags
	}

	d.SetId(network.Ref)
	return resourceNetworkRead(ctx, d, m)
}

func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	client := m.(*infoblox.Client)

	var network infoblox.Network

	if d.HasChange("cidr") {
		network.CIDR = d.Get("cidr").(string)
	}
	if d.HasChange("gateway_ip") {
		gw, ok := d.GetOk("gateway_ip")
		if ok {
			old, _ := d.GetChange("gateway_ip")
			_, parsedNetwork, _ := net.ParseCIDR(d.Get("cidr").(string))
			if parsedNetwork.Contains(net.ParseIP(gw.(string))) != true {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Invalid gateway IP address",
					Detail:   fmt.Sprintf("Gateway address: %s is not within network CIDR: %s", gw.(string), d.Get("cidr").(string)),
				})
				d.Set("gateway_ip", old.(string))
				return diags
			}
		}
	}
	if d.HasChange("comment") {
		network.Comment = d.Get("comment").(string)
	}
	if d.HasChange("disable_dhcp") {
		network.DisableDHCP = newBool(d.Get("disable_dhcp").(bool))
	}
	if d.HasChange("grid_ref") {
		network.Comment = d.Get("grid_ref").(string)
	}
	if d.HasChange("restart_if_needed") {
		network.DisableDHCP = newBool(d.Get("restart_if_needed").(bool))
	}
	if d.HasChange("network_view") {
		network.NetworkView = d.Get("network_view").(string)
	}
	if d.HasChange("member") {
		memberList := d.Get("member").([]interface{})
		network.Members = []infoblox.Member{}
		if len(memberList) > 0 {
			for _, member := range memberList {
				network.Members = append(network.Members, infoblox.Member{
					StructType:  member.(map[string]interface{})["struct"].(string),
					Hostname:    member.(map[string]interface{})["hostname"].(string),
					IPV4Address: member.(map[string]interface{})["ip_v4_address"].(string),
					IPV6Address: member.(map[string]interface{})["ip_v6_address"].(string),
				})
			}
		}
	}
	if d.HasChange("option") {
		optionList := d.Get("option").(*schema.Set).List()
		network.Options = []infoblox.Option{}
		if len(optionList) > 0 {
			for _, option := range optionList {
				network.Options = append(network.Options, infoblox.Option{
					Name:        option.(map[string]interface{})["name"].(string),
					UseOption:   newBool(option.(map[string]interface{})["use_option"].(bool)),
					Value:       option.(map[string]interface{})["value"].(string),
					VendorClass: option.(map[string]interface{})["vendor_class"].(string),
				})
			}
		}
	}
	if d.HasChange("extensible_attributes") {
		eaMap := d.Get("extensible_attributes").(map[string]interface{})
		if len(eaMap) > 0 {
			eas, err := createExtensibleAttributesFromJSON(client, eaMap)
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
				return diags
			}
			network.ExtensibleAttributes = &eas
		}
		if client.OrchestratorEAs != nil && len(*client.OrchestratorEAs) > 0 {
			for k, v := range *client.OrchestratorEAs {
				(*network.ExtensibleAttributes)[k] = v
			}
		}
	}
	changedNetwork, err := client.UpdateNetwork(d.Id(), network)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	if d.HasChanges("gateway_ip", "gateway_comment", "gateway_extensible_attributes") {
		_, ok := d.GetOk("gateway_ip")
		old, new := d.GetChange("gateway_ip")
		if new.(string) != "" || (ok && d.HasChanges("gateway_comment", "gateway_etensible_attributes")) {
			gwFixedAddress := infoblox.FixedAddress{
				IPAddress:   new.(string),
				Comment:     d.Get("gateway_comment").(string),
				MatchClient: "RESERVED",
			}
			eaMap := d.Get("gateway_extensible_attributes").(map[string]interface{})
			if len(eaMap) > 0 {
				eas, err := createExtensibleAttributesFromJSON(client, eaMap)
				if err != nil {
					diags = append(diags, diag.FromErr(err)...)
				}
				gwFixedAddress.ExtensibleAttributes = &eas
			}

			if client.OrchestratorEAs != nil && len(*client.OrchestratorEAs) > 0 {
				for k, v := range *client.OrchestratorEAs {
					(*gwFixedAddress.ExtensibleAttributes)[k] = v
				}
			}
			if old != "" {
				existingFixedAddress, err := client.GetFixedAddressByQuery(map[string]string{
					"network":  d.Get("cidr").(string),
					"ipv4addr": old.(string),
				})
				if err != nil {
					diags = append(diags, diag.FromErr(err)...)
				} else {
					_, err = client.UpdateFixedAddress(existingFixedAddress[0].Ref, gwFixedAddress)
					if err != nil {
						diags = append(diags, diag.FromErr(err)...)
					}
				}
			} else {
				err = client.CreateFixedAddress(&gwFixedAddress)
				if err != nil {
					diags = append(diags, diag.FromErr(err)...)
				}
			}
		} else {
			existingFixedAddress, err := client.GetFixedAddressByQuery(map[string]string{
				"network":  d.Get("cidr").(string),
				"ipv4addr": old.(string),
			})
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
			} else {
				err = client.DeleteFixedAddress(existingFixedAddress[0].Ref)
				if err != nil {
					diags = append(diags, diag.FromErr(err)...)
				}
			}
		}
	}

	d.SetId(changedNetwork.Ref)
	if d.Get("restart_if_needed").(bool) && len(changedNetwork.Members) == 1 {
		err := client.RestartServices(d.Get("grid_ref").(string), infoblox.GridServiceRestartRequest{
			RestartOption: "RESTART_IF_NEEDED",
			Services:      []string{"DHCP"},
			Members:       []string{changedNetwork.Members[0].Hostname},
		})
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
		// time.Sleep(2 * time.Second)
	}

	if diags.HasError() {
		return diags
	}

	return resourceNetworkRead(ctx, d, m)
}

func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics
	ref := d.Id()

	network, err := convertResourceDataToNetwork(client, d)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	err = client.DeleteNetwork(ref)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.Get("restart_if_needed").(bool) && len(network.Members) == 1 {
		err = client.RestartServices(d.Get("grid_ref").(string), infoblox.GridServiceRestartRequest{
			RestartOption: "RESTART_IF_NEEDED",
			Services:      []string{"DHCP"},
			Members:       []string{network.Members[0].Hostname},
		})
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
	}

	return diags
}
