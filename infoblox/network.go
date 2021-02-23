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
			eaCustomDiff,
			optionCustomDiff,
		),
		Schema: map[string]*schema.Schema{
			"ref": {
				Type:        schema.TypeString,
				Description: "Reference id of network object",
				Computed:    true,
			},
			"cidr": {
				Type:             schema.TypeString,
				Description:      "CIDR of network",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
			},
			"comment": {
				Type:        schema.TypeString,
				Description: "Comment string",
				Optional:    true,
				Computed:    true,
			},
			"disable_dhcp": {
				Type:        schema.TypeBool,
				Description: "Disable for DHCP",
				Optional:    true,
				Computed:    true,
			},
			"grid_ref": {
				Type:         schema.TypeString,
				Description:  "Ref for grid needed for restarting services",
				Optional:     true,
				RequiredWith: []string{"restart_if_needed"},
			},
			"restart_if_needed": {
				Type:        schema.TypeBool,
				Description: "Restart dhcp services if needed",
				Optional:    true,
			},
			"network_view": {
				Type:        schema.TypeString,
				Description: "Network view",
				Optional:    true,
				Computed:    true,
			},
			"member": {
				Type:        schema.TypeList,
				Description: "Grid members associated with network",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"struct": {
							Type:             schema.TypeString,
							Description:      "Struct type of member",
							Optional:         true,
							Default:          "dhcpmember",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"dhcpmember", "msdhcpserver"}, true)),
							StateFunc: func(val interface{}) string {
								return strings.ToLower(val.(string))
							},
						},
						"ip_v4_address": {
							Type:             schema.TypeString,
							Description:      "IPv4 address",
							Optional:         true,
							Computed:         true,
							ConflictsWith:    []string{"member.0.ip_v6_address"},
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv4Address),
						},
						"ip_v6_address": {
							Type:             schema.TypeString,
							Description:      "IPv6 address",
							Optional:         true,
							Computed:         true,
							ConflictsWith:    []string{"member.0.ip_v4_address"},
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
				Optional:    true,
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
							Optional:    true,
							Default:     true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "Value of option",
							Required:    true,
						},
						"vendor_class": {
							Type:        schema.TypeString,
							Description: "Value of option",
							Optional:    true,
							Default:     "DHCP",
						},
					},
				},
			},
			"extensible_attributes": {
				Type:             schema.TypeMap,
				Description:      "Extensible attributes of network",
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

	d.SetId(changedNetwork.Ref)
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
