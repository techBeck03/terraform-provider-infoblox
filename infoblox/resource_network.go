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
			optionCustomDiff,
		),
		Schema: map[string]*schema.Schema{
			"cidr": {
				Type:             schema.TypeString,
				Description:      "The network address in IPv4 Address/CIDR format.",
				Required:         true,
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
				Computed:    true,
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
	changedNetwork, err := client.UpdateNetwork(d.Id(), network)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
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
