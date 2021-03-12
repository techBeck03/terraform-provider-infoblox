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
	fixedAddressRequiredIPFields = []string{
		"network",
		"ip_v4_address",
		"range_function_string",
	}
)

func resourceFixedAddress() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFixedAddressCreate,
		ReadContext:   resourceFixedAddressRead,
		UpdateContext: resourceFixedAddressUpdate,
		DeleteContext: resourceFixedAddressDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
		CustomizeDiff: customdiff.Sequence(
			makeEACustomDiff("extensible_attributes"),
		),
		Schema: map[string]*schema.Schema{
			"ref": {
				Type:        schema.TypeString,
				Description: "Reference id of host fixed address object",
				Computed:    true,
			},
			"hostname": {
				Type:        schema.TypeString,
				Description: "Hostname of host fixed address",
				Optional:    true,
				Computed:    true,
			},
			"cidr": {
				Type:             schema.TypeString,
				Description:      "Hostname of host fixed address",
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
			},
			"comment": {
				Type:        schema.TypeString,
				Description: "Comment string",
				Optional:    true,
				Computed:    true,
			},
			"ip_address": {
				Type:             schema.TypeString,
				Description:      "IPv4 address",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv4Address),
			},
			"network_view": {
				Type:        schema.TypeString,
				Description: "Network view",
				Optional:    true,
				Computed:    true,
			},
			"mac": {
				Type:        schema.TypeString,
				Description: "MAC address",
				Optional:    true,
				Computed:    true,
			},
			"match_client": {
				Type:             schema.TypeString,
				Description:      "Match client",
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"MAC_ADDRESS", "CLIENT_ID", "RESERVED", "CIRCUIT_ID", "REMOTE_ID"}, false)),
			},
			"disable": {
				Type:        schema.TypeBool,
				Description: "Disabled",
				Optional:    true,
				Default:     false,
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
			"member": {
				Type:        schema.TypeList,
				Description: "Grid member associated with range",
				Required:    true,
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
				Description:      "Extensible attributes of range",
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

func convertFixedAddressToResourceData(client *infoblox.Client, d *schema.ResourceData, fixedAddress *infoblox.FixedAddress) diag.Diagnostics {
	var diags diag.Diagnostics

	d.Set("ref", fixedAddress.Ref)
	d.Set("cidr", fixedAddress.CIDR)
	d.Set("ip_address", fixedAddress.IPAddress)
	d.Set("hostname", fixedAddress.Hostname)
	d.Set("comment", fixedAddress.Comment)
	d.Set("disable", fixedAddress.Disable)
	d.Set("network_view", fixedAddress.NetworkView)
	d.Set("mac", fixedAddress.Mac)
	d.Set("match_client", fixedAddress.MatchClient)

	var optionList []map[string]interface{}
	for _, option := range fixedAddress.Options {
		optionList = append(optionList, map[string]interface{}{
			"name":         option.Name,
			"code":         option.Code,
			"use_option":   option.UseOption,
			"value":        option.Value,
			"vendor_class": option.VendorClass,
		})
	}

	d.Set("option", optionList)

	eas, err := client.ConvertEAsToJSONString(*fixedAddress.ExtensibleAttributes)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	} else {
		d.Set("extensible_attributes", eas)
	}

	return diags
}

func convertResourceDataToFixedAddress(client *infoblox.Client, d *schema.ResourceData) (*infoblox.FixedAddress, error) {
	var fixedAddress infoblox.FixedAddress

	fixedAddress.CIDR = d.Get("cidr").(string)
	fixedAddress.IPAddress = d.Get("ip_address").(string)
	fixedAddress.Hostname = d.Get("hostname").(string)
	fixedAddress.Comment = d.Get("comment").(string)
	fixedAddress.Disable = newBool(d.Get("disable").(bool))
	fixedAddress.NetworkView = d.Get("network_view").(string)
	fixedAddress.Mac = d.Get("mac").(string)
	fixedAddress.MatchClient = d.Get("match_client").(string)

	optionList := d.Get("option").(*schema.Set).List()
	fixedAddress.Options = []infoblox.Option{}
	if len(optionList) > 0 {
		for _, option := range optionList {
			fixedAddress.Options = append(fixedAddress.Options, infoblox.Option{
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
			return &fixedAddress, err
		}
		fixedAddress.ExtensibleAttributes = &eas
	}

	if client.OrchestratorEAs != nil && len(*client.OrchestratorEAs) > 0 {
		for k, v := range *client.OrchestratorEAs {
			(*fixedAddress.ExtensibleAttributes)[k] = v
		}
	}

	return &fixedAddress, nil
}

func resourceFixedAddressRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics
	ref := d.Id()

	fixedAddress, err := client.GetFixedAddressByRef(ref, nil)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	check := convertFixedAddressToResourceData(client, d, &fixedAddress)
	if check.HasError() {
		return check
	}

	d.SetId(fixedAddress.Ref)

	return diags
}

func resourceFixedAddressCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)
	var member *infoblox.Member

	var diags diag.Diagnostics

	fixedAddress, err := convertResourceDataToFixedAddress(client, d)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	err = client.CreateFixedAddress(fixedAddress)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	memberList := d.Get("member").([]interface{})
	if len(memberList) > 0 {
		var members []infoblox.Member
		for _, member := range memberList {
			members = append(members, infoblox.Member{
				StructType:  member.(map[string]interface{})["struct"].(string),
				Hostname:    member.(map[string]interface{})["hostname"].(string),
				IPV4Address: member.(map[string]interface{})["ip_v4_address"].(string),
				IPV6Address: member.(map[string]interface{})["ip_v6_address"].(string),
			})
		}
		member = &members[0]
	}

	if d.Get("restart_if_needed").(bool) && member != nil {
		err = client.RestartServices(d.Get("grid_ref").(string), infoblox.GridServiceRestartRequest{
			RestartOption: "RESTART_IF_NEEDED",
			Services:      []string{"DHCP"},
			Members:       []string{member.Hostname},
		})
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
	}

	if diags.HasError() {
		return diags
	}

	d.SetId(fixedAddress.Ref)
	return resourceFixedAddressRead(ctx, d, m)
}

func resourceFixedAddressUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	client := m.(*infoblox.Client)
	var member *infoblox.Member

	var fixedAddress infoblox.FixedAddress

	if d.HasChange("cidr") {
		fixedAddress.CIDR = d.Get("cidr").(string)
	}
	if d.HasChange("ip_address") {
		fixedAddress.IPAddress = d.Get("ip_address").(string)
	}
	if d.HasChange("hostname") {
		fixedAddress.Hostname = d.Get("hostname").(string)
	}
	if d.HasChange("comment") {
		fixedAddress.Comment = d.Get("comment").(string)
	}
	if d.HasChange("disable") {
		fixedAddress.Disable = newBool(d.Get("disable").(bool))
	}
	if d.HasChange("network_view") {
		fixedAddress.NetworkView = d.Get("network_view").(string)
	}
	if d.HasChange("mac") {
		fixedAddress.Mac = d.Get("mac").(string)
	}
	if d.HasChange("match_client") {
		fixedAddress.MatchClient = d.Get("match_client").(string)
	}

	memberList := d.Get("member").([]interface{})
	if len(memberList) > 0 {
		var members []infoblox.Member
		for _, member := range memberList {
			members = append(members, infoblox.Member{
				StructType:  member.(map[string]interface{})["struct"].(string),
				Hostname:    member.(map[string]interface{})["hostname"].(string),
				IPV4Address: member.(map[string]interface{})["ip_v4_address"].(string),
				IPV6Address: member.(map[string]interface{})["ip_v6_address"].(string),
			})
		}
		member = &members[0]
	}

	if d.HasChange("option") {
		optionList := d.Get("option").(*schema.Set).List()
		fixedAddress.Options = []infoblox.Option{}
		if len(optionList) > 0 {
			for _, option := range optionList {
				fixedAddress.Options = append(fixedAddress.Options, infoblox.Option{
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
			fixedAddress.ExtensibleAttributes = &eas
		}
		if client.OrchestratorEAs != nil && len(*client.OrchestratorEAs) > 0 {
			for k, v := range *client.OrchestratorEAs {
				(*fixedAddress.ExtensibleAttributes)[k] = v
			}
		}
	}

	changedFixedAddress, err := client.UpdateFixedAddress(d.Id(), fixedAddress)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if d.Get("restart_if_needed").(bool) && member != nil {
		err = client.RestartServices(d.Get("grid_ref").(string), infoblox.GridServiceRestartRequest{
			RestartOption: "RESTART_IF_NEEDED",
			Services:      []string{"DHCP"},
			Members:       []string{member.Hostname},
		})
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
	}

	d.SetId(changedFixedAddress.Ref)
	return resourceFixedAddressRead(ctx, d, m)
}

func resourceFixedAddressDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)
	var member *infoblox.Member

	var diags diag.Diagnostics
	ref := d.Id()

	err := client.DeleteFixedAddress(ref)
	if err != nil {
		return diag.FromErr(err)
	}

	memberList := d.Get("member").([]interface{})
	if len(memberList) > 0 {
		var members []infoblox.Member
		for _, member := range memberList {
			members = append(members, infoblox.Member{
				StructType:  member.(map[string]interface{})["struct"].(string),
				Hostname:    member.(map[string]interface{})["hostname"].(string),
				IPV4Address: member.(map[string]interface{})["ip_v4_address"].(string),
				IPV6Address: member.(map[string]interface{})["ip_v6_address"].(string),
			})
		}
		member = &members[0]
	}

	if d.Get("restart_if_needed").(bool) && member != nil {
		err = client.RestartServices(d.Get("grid_ref").(string), infoblox.GridServiceRestartRequest{
			RestartOption: "RESTART_IF_NEEDED",
			Services:      []string{"DHCP"},
			Members:       []string{member.Hostname},
		})
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
	}

	return diags
}
