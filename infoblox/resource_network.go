package infoblox

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/techBeck03/go-ipmath"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkCreate,
		ReadContext:   resourceNetworkRead,
		UpdateContext: resourceNetworkUpdate,
		DeleteContext: resourceNetworkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: customdiff.Sequence(
			makeEACustomDiffNetwork("extensible_attributes"),
			optionCustomDiff,
		),
		Schema: map[string]*schema.Schema{
			"cidr": {
				Type:             schema.TypeString,
				Description:      "The network address in IPv4 Address/CIDR format.",
				Optional:         true,
				Computed:         true,
				AtLeastOneOf:     []string{"cidr", "parent_cidr", "ea_search"},
				ConflictsWith:    []string{"ea_search", "parent_cidr"},
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
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
				Default:     false,
			},
			"ea_search": {
				Type:          schema.TypeMap,
				Description:   "Ea search criteria for next_available_network function",
				Optional:      true,
				Default:       "",
				ConflictsWith: []string{"cidr", "parent_cidr"},
				AtLeastOneOf:  []string{"cidr", "parent_cidr", "ea_search"},
				RequiredWith:  []string{"prefix_length"},
				ForceNew:      true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
			"gateway_ea": {
				Type:         schema.TypeString,
				Description:  "Name of extensible attribute for gateway address",
				Optional:     true,
				RequiredWith: []string{"gateway_offset"},
			},
			"gateway_ip": {
				Type:         schema.TypeString,
				Description:  "IP address of default gateway if auto-created",
				Computed:     true,
				RequiredWith: []string{"gateway_offset"},
			},
			"gateway_label": {
				Type:        schema.TypeString,
				Description: "Name to apply to gateway reservation",
				Optional:    true,
				Default:     "Gateway",
			},
			"gateway_offset": {
				Type:             schema.TypeInt,
				Description:      "Offset from network address to reserve for default gateway",
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
			},
			"gateway_ref": {
				Type:        schema.TypeString,
				Description: "Reference id for gateway if created",
				Computed:    true,
			},
			"grid_ref": {
				Type:         schema.TypeString,
				Description:  "Ref for grid needed for restarting services.",
				Optional:     true,
				RequiredWith: []string{"restart_if_needed"},
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
			"network_view": {
				Type:        schema.TypeString,
				Description: "The name of the network view in which this network resides.",
				Optional:    true,
				Computed:    true,
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
							Required:    true,
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
			"parent_cidr": {
				Type:          schema.TypeString,
				Description:   "Parent container CIDR subnet",
				Optional:      true,
				Default:       "",
				AtLeastOneOf:  []string{"cidr", "parent_cidr", "ea_search"},
				ConflictsWith: []string{"cidr", "ea_search"},
				RequiredWith:  []string{"prefix_length"},
				ForceNew:      true,
			},
			"prefix_length": {
				Type:             schema.TypeInt,
				Description:      "Desired prefix size of requested network",
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
				ForceNew:         true,
			},
			"ref": {
				Type:        schema.TypeString,
				Description: "Reference id of network object.",
				Computed:    true,
			},
			"restart_if_needed": {
				Type:        schema.TypeBool,
				Description: "Restart dhcp services if needed.",
				Optional:    true,
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
				Code:        option.(map[string]interface{})["code"].(int),
				UseOption:   newBool(option.(map[string]interface{})["use_option"].(bool)),
				Value:       option.(map[string]interface{})["value"].(string),
				VendorClass: option.(map[string]interface{})["vendor_class"].(string),
			})
		}
	}

	eaMap := d.Get("extensible_attributes").(map[string]interface{})
	if len(eaMap) > 0 {
		eas, err := createExtensibleAttributesFromJSON(eaMap)
		if err != nil {
			return &network, err
		}
		network.ExtensibleAttributes = &eas
	}

	if client.OrchestratorEAs != nil && len(*client.OrchestratorEAs) > 0 {
		if network.ExtensibleAttributes == nil {
			network.ExtensibleAttributes = &infoblox.ExtensibleAttribute{}
		}
		for k, v := range *client.OrchestratorEAs {
			(*network.ExtensibleAttributes)[k] = v
		}
	}

	return &network, nil
}

func convertResourceDataToNetworkFromContainer(client *infoblox.Client, d *schema.ResourceData) (*infoblox.NetworkFromContainer, error) {
	var network infoblox.NetworkFromContainer

	prefix := d.Get("prefix_length").(int)
	parent_cidr := d.Get("parent_cidr").(string)

	if parent_cidr != "" {
		network.Network = infoblox.NetworkContainerFunction{
			Function:    "next_available_network",
			ResultField: "networks",
			Object:      "networkcontainer",
			ObjectParameters: map[string]string{
				"network": parent_cidr,
			},
			Parameters: map[string]int{
				"cidr": prefix,
			},
		}
	} else {
		ea_search_map := d.Get("ea_search").(map[string]interface{})
		ea_search := make(map[string]string)
		for k, v := range ea_search_map {
			ea_search[k] = v.(string)
		}
		network.Network = infoblox.NetworkContainerFunction{
			Function:         "next_available_network",
			ResultField:      "networks",
			Object:           "networkcontainer",
			ObjectParameters: ea_search,
			Parameters: map[string]int{
				"cidr": prefix,
			},
		}
	}
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
				Code:        option.(map[string]interface{})["code"].(int),
				UseOption:   newBool(option.(map[string]interface{})["use_option"].(bool)),
				Value:       option.(map[string]interface{})["value"].(string),
				VendorClass: option.(map[string]interface{})["vendor_class"].(string),
			})
		}
	}

	eaMap := d.Get("extensible_attributes").(map[string]interface{})
	if len(eaMap) > 0 {
		eas, err := createExtensibleAttributesFromJSON(eaMap)
		if err != nil {
			return &network, err
		}
		network.ExtensibleAttributes = &eas
	}

	if client.OrchestratorEAs != nil && len(*client.OrchestratorEAs) > 0 {
		if network.ExtensibleAttributes == nil {
			network.ExtensibleAttributes = &infoblox.ExtensibleAttribute{}
		}
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

	var network *infoblox.Network

	cidr := d.Get("cidr").(string)
	if cidr == "" {
		net, err := convertResourceDataToNetworkFromContainer(client, d)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
		cResult, err := client.CreateNetworkFromContainer(net)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
		network = &cResult
	} else {
		net, err := convertResourceDataToNetwork(client, d)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
		err = client.CreateNetwork(net)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
		network = net
	}

	if diags.HasError() {
		return diags
	}

	gw_offset := d.Get("gateway_offset").(int)
	if gw_offset > 0 {
		ip, err := ipmath.NewIP(network.CIDR)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
		err = ip.Add(gw_offset)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
		gw_reservation := infoblox.FixedAddress{
			IPAddress:   ip.ToIPString(),
			NetworkView: network.NetworkView,
			CIDR:        network.CIDR,
			Hostname:    d.Get("gateway_label").(string),
			MatchClient: "RESERVED",
		}
		err = client.CreateFixedAddress(&gw_reservation)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		d.Set("gateway_ref", gw_reservation.Ref)
		d.Set("gateway_ip", gw_reservation.IPAddress)
		gw_ea := d.Get("gateway_ea").(string)
		if gw_ea != "" {
			update := infoblox.Network{
				Ref: network.Ref,
			}
			eas := &infoblox.ExtensibleAttribute{}
			(*eas)[gw_ea] = infoblox.ExtensibleAttributeValue{
				Value: ip.ToIPString(),
			}
			update.ExtensibleAttributesAdd = eas
			updated_network, err := client.UpdateNetwork(network.Ref, update)
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
				return diags
			}
			network = &updated_network
		}
	}

	if d.Get("restart_if_needed").(bool) && len(network.Members) == 1 {
		err := client.RestartServices(d.Get("grid_ref").(string), infoblox.GridServiceRestartRequest{
			RestartOption: "RESTART_IF_NEEDED",
			Services:      []string{"DHCP"},
			Members:       []string{network.Members[0].Hostname},
		})
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
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
					Code:        option.(map[string]interface{})["code"].(int),
					UseOption:   newBool(option.(map[string]interface{})["use_option"].(bool)),
					Value:       option.(map[string]interface{})["value"].(string),
					VendorClass: option.(map[string]interface{})["vendor_class"].(string),
				})
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
			network.ExtensibleAttributesRemove = &infoblox.ExtensibleAttribute{}
			for _, v := range removeEAs {
				(*network.ExtensibleAttributesRemove)[v] = oldEAs[v]
			}
		}
		for k, v := range newEAs {
			if !Contains(oldKeys, k) || (Contains(oldKeys, k) && v.Value != oldEAs[k].Value) {
				if network.ExtensibleAttributesAdd == nil {
					network.ExtensibleAttributesAdd = &infoblox.ExtensibleAttribute{}
				}
				(*network.ExtensibleAttributesAdd)[k] = v
			}
		}
		if client.OrchestratorEAs != nil && len(*client.OrchestratorEAs) > 0 {
			if network.ExtensibleAttributesAdd == nil {
				network.ExtensibleAttributesAdd = &infoblox.ExtensibleAttribute{}
			}
			for k, v := range *client.OrchestratorEAs {
				(*network.ExtensibleAttributesAdd)[k] = v
			}
		}
	}

	if d.HasChange("gateway_ea") && !d.HasChange("gateway_offset") {
		old, _ := d.GetChange("extensible_attributes")
		oldEAs, err := createExtensibleAttributesFromJSON(old.(map[string]interface{}))
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
		if network.ExtensibleAttributesRemove == nil {
			network.ExtensibleAttributesRemove = &infoblox.ExtensibleAttribute{}
		}
		(*network.ExtensibleAttributesRemove)[d.Get("gateway_ea").(string)] = oldEAs[d.Get("gateway_ea").(string)]
		if network.ExtensibleAttributesAdd == nil {
			network.ExtensibleAttributesAdd = &infoblox.ExtensibleAttribute{}
		}
		gw_offset := d.Get("gateway_offset").(int)
		ip, err := ipmath.NewIP(d.Get("cidr").(string))
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
		err = ip.Add(gw_offset)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
		(*network.ExtensibleAttributesAdd)[d.Get("gateway_ea").(string)] = infoblox.ExtensibleAttributeValue{
			Value: ip.ToIPString(),
		}
	}

	changedNetwork, err := client.UpdateNetwork(d.Id(), network)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if d.HasChange("gateway_offset") {
		gw_ref := d.Get("gateway_ref").(string)
		if gw_ref != "" {
			err = client.DeleteFixedAddress(gw_ref)
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
				return diags
			}
		}
		gw_offset := d.Get("gateway_offset").(int)
		ip, err := ipmath.NewIP(changedNetwork.CIDR)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
		err = ip.Add(gw_offset)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
		gw_reservation := infoblox.FixedAddress{
			IPAddress:   ip.ToIPString(),
			NetworkView: changedNetwork.NetworkView,
			CIDR:        changedNetwork.CIDR,
			Hostname:    d.Get("gateway_label").(string),
			MatchClient: "RESERVED",
		}
		err = client.CreateFixedAddress(&gw_reservation)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		d.Set("gateway_ref", gw_reservation.Ref)
		d.Set("gateway_ip", gw_reservation.IPAddress)
		gw_ea := d.Get("gateway_ea").(string)
		if gw_ea != "" {
			update := infoblox.Network{
				Ref: changedNetwork.Ref,
			}
			eas := &infoblox.ExtensibleAttribute{}
			(*eas)[gw_ea] = infoblox.ExtensibleAttributeValue{
				Value: ip.ToIPString(),
			}
			update.ExtensibleAttributesAdd = eas
			updated_network, err := client.UpdateNetwork(changedNetwork.Ref, update)
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
				return diags
			}
			changedNetwork = updated_network
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
