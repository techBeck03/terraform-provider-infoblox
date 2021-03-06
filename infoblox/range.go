package infoblox

import (
	"context"
	"net"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/techBeck03/go-ipmath"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

var (
	rangeRequiredFields = []string{
		"start_address",
		"sequential_count",
	}
	dhcpRequiredFields = []string{
		"disable_dhcp",
		"member",
	}
)

func resourceRange() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRangeCreate,
		ReadContext:   resourceRangeRead,
		UpdateContext: resourceRangeUpdate,
		DeleteContext: resourceRangeDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
		CustomizeDiff: customdiff.Sequence(
			makeEACustomDiff("extensible_attributes"),
			makeAddressCompareCustomDiff("start_address", "end_address"),
			rangeForceNew,
			makeCidrContainsIPCheck("cidr", []string{"start_address", "end_address"}),
			makeLowerThanIPCheck("start_address", "end_address"),
		),
		Schema: map[string]*schema.Schema{
			"ref": {
				Type:        schema.TypeString,
				Description: "Reference id of range object",
				Optional:    true,
				Computed:    true,
			},
			"comment": {
				Type:        schema.TypeString,
				Description: "Comment string",
				Optional:    true,
				Computed:    true,
			},
			"cidr": {
				Type:             schema.TypeString,
				Description:      "Network for range in CIDR notation",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
			},
			"disable_dhcp": {
				Type:          schema.TypeBool,
				Description:   "Disable for DHCP",
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"member"},
				AtLeastOneOf:  dhcpRequiredFields,
			},
			"network_view": {
				Type:        schema.TypeString,
				Description: "Network view",
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
			},
			"start_address": {
				Type:             schema.TypeString,
				Description:      "Starting IP address",
				Optional:         true,
				Computed:         true,
				ConflictsWith:    []string{"sequential_count"},
				RequiredWith:     []string{"end_address"},
				AtLeastOneOf:     rangeRequiredFields,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv4Address),
			},
			"end_address": {
				Type:             schema.TypeString,
				Description:      "Starting IP address",
				Optional:         true,
				Computed:         true,
				RequiredWith:     []string{"start_address"},
				ConflictsWith:    []string{"sequential_count"},
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv4Address),
			},
			"sequential_count": {
				Type:          schema.TypeInt,
				Description:   "Sequential count of addresses",
				Optional:      true,
				ConflictsWith: []string{"start_address", "end_address"},
				AtLeastOneOf:  rangeRequiredFields,
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
				Type:          schema.TypeList,
				Description:   "Grid member associated with range",
				Optional:      true,
				ConflictsWith: []string{"disable_dhcp"},
				AtLeastOneOf:  dhcpRequiredFields,
				MaxItems:      1,
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

func convertRangeToResourceData(client *infoblox.Client, d *schema.ResourceData, addressRange *infoblox.Range) diag.Diagnostics {
	var diags diag.Diagnostics

	d.Set("ref", addressRange.Ref)
	d.Set("cidr", addressRange.CIDR)
	d.Set("comment", addressRange.Comment)
	d.Set("start_address", addressRange.StartAddress)
	d.Set("end_address", addressRange.EndAddress)
	d.Set("network_view", addressRange.NetworkView)
	d.Set("disable_dhcp", addressRange.DisableDHCP)

	var memberList []map[string]interface{}
	if addressRange.Member != nil {
		memberList = append(memberList, map[string]interface{}{
			"struct":        addressRange.Member.StructType,
			"hostname":      addressRange.Member.Hostname,
			"ip_v4_address": addressRange.Member.IPV4Address,
			"ip_v6_address": addressRange.Member.IPV6Address,
		})

		d.Set("member", memberList)
	}

	eas, err := client.ConvertEAsToJSONString(*addressRange.ExtensibleAttributes)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	} else {
		d.Set("extensible_attributes", eas)
	}

	return diags
}

func convertResourceDataToRange(client *infoblox.Client, d *schema.ResourceData) (*infoblox.Range, error) {
	var addressRange infoblox.Range

	addressRange.CIDR = d.Get("cidr").(string)
	addressRange.Comment = d.Get("comment").(string)
	addressRange.StartAddress = d.Get("start_address").(string)
	addressRange.EndAddress = d.Get("end_address").(string)
	addressRange.NetworkView = d.Get("network_view").(string)
	addressRange.DisableDHCP = newBool(d.Get("disable_dhcp").(bool))

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
		addressRange.Member = &members[0]
	}

	eaMap := d.Get("extensible_attributes").(map[string]interface{})
	if len(eaMap) > 0 {
		eas, err := createExtensibleAttributesFromJSON(client, eaMap)
		if err != nil {
			return &addressRange, err
		}
		addressRange.ExtensibleAttributes = &eas
	}

	if client.OrchestratorEAs != nil && len(*client.OrchestratorEAs) > 0 {
		for k, v := range *client.OrchestratorEAs {
			(*addressRange.ExtensibleAttributes)[k] = v
		}
	}

	return &addressRange, nil
}

func resourceRangeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics
	ref := d.Id()

	addressRange, err := client.GetRangeByRef(ref, nil)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	check := convertRangeToResourceData(client, d, &addressRange)
	if check.HasError() {
		return check
	}

	d.SetId(addressRange.Ref)

	return diags
}

func resourceRangeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics

	addressRange, err := convertResourceDataToRange(client, d)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	count, countOk := d.GetOk("sequential_count")
	if countOk {
		err = client.CreateSequentialRange(addressRange, infoblox.AddressQuery{
			CIDR:  addressRange.CIDR,
			Count: count.(int),
		})
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
	} else {
		err = client.CreateRange(addressRange)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
	}

	if d.Get("restart_if_needed").(bool) && addressRange.Member != nil {
		err = client.RestartServices(d.Get("grid_ref").(string), infoblox.GridServiceRestartRequest{
			RestartOption: "RESTART_IF_NEEDED",
			Services:      []string{"DHCP"},
			Members:       []string{addressRange.Member.Hostname},
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

	d.SetId(addressRange.Ref)
	return resourceRangeRead(ctx, d, m)
}

func resourceRangeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	client := m.(*infoblox.Client)

	var addressRange infoblox.Range

	if d.HasChange("comment") {
		addressRange.Comment = d.Get("comment").(string)
	}
	if d.HasChange("disable_dhcp") {
		addressRange.DisableDHCP = newBool(d.Get("disable_dhcp").(bool))
	}
	if d.HasChange("sequential_count") {
		addressRange.StartAddress = d.Get("start_address").(string)
		old, new := d.GetChange("sequential_count")
		if new.(int) < old.(int) {
			endAddress := ipmath.IP{
				Address: net.ParseIP(d.Get("end_address").(string)),
			}
			endAddressFinal, err := endAddress.Subtract(old.(int) - new.(int))
			if err != nil {
				d.Set("sequential_count", old.(int))
				return diag.FromErr(err)
			}
			addressRange.EndAddress = endAddressFinal.String()
		} else {
			endAddress := ipmath.IP{
				Address: net.ParseIP(d.Get("end_address").(string)),
			}
			endAddressFinal, err := endAddress.Add(new.(int) - old.(int))
			if err != nil {
				d.Set("sequential_count", old.(int))
				return diag.FromErr(err)
			}
			startAddress := ipmath.IP{
				Address: net.ParseIP(d.Get("end_address").(string)),
			}
			startAddressFinal, err := startAddress.Inc()
			if err != nil {
				d.Set("sequential_count", old.(int))
				return diag.FromErr(err)
			}
			check, err := client.GetSequentialAddressRange(infoblox.AddressQuery{
				CIDR:         d.Get("cidr").(string),
				StartAddress: startAddressFinal.String(),
				EndAddress:   endAddressFinal.String(),
				Count:        new.(int) - old.(int),
			})
			if err != nil {
				d.Set("sequential_count", old.(int))
				return diag.FromErr(err)
			}
			if check == nil {
				d.Set("sequential_count", old.(int))
				return diag.Errorf("Sequential address count increase overlaps with another range or USED IP")
			}
			addressRange.EndAddress = endAddressFinal.String()
		}
	}

	if d.HasChange("start_address") {
		old, new := d.GetChange("start_address")
		oldIP := ipmath.IP{
			Address: net.ParseIP(old.(string)),
		}
		newIP := ipmath.IP{
			Address: net.ParseIP(new.(string)),
		}
		if newIP.LT(oldIP.Address) {
			endAddress, err := oldIP.Dec()
			if err != nil {
				return diag.FromErr(err)
			}
			_, err = client.GetSequentialAddressRange(infoblox.AddressQuery{
				CIDR:         d.Get("cidr").(string),
				StartAddress: newIP.Address.String(),
				EndAddress:   endAddress.String(),
				Count:        newIP.Difference(oldIP.Address),
			})
			if err != nil {
				d.Set("start_address", old.(string))
				return diag.Errorf("Decreasing the `start_address` causes an overlap with another range or USED IP")
			}
		}
		addressRange.StartAddress = new.(string)
	}

	if d.HasChange("end_address") {
		old, new := d.GetChange("end_address")
		oldIP := ipmath.IP{
			Address: net.ParseIP(old.(string)),
		}
		newIP := ipmath.IP{
			Address: net.ParseIP(new.(string)),
		}
		if newIP.GT(oldIP.Address) {
			startAddress, err := oldIP.Inc()
			if err != nil {
				return diag.FromErr(err)
			}
			_, err = client.GetSequentialAddressRange(infoblox.AddressQuery{
				CIDR:         d.Get("cidr").(string),
				StartAddress: startAddress.String(),
				EndAddress:   oldIP.Address.String(),
				Count:        oldIP.Difference(newIP.Address),
			})
			if err != nil {
				d.Set("end_address", old.(string))
				return diag.Errorf("Increasing the `end_address` causes an overlap with another range or USED IP")
			}
		}
		addressRange.EndAddress = new.(string)
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
		addressRange.Member = &members[0]
	}

	if d.HasChange("extensible_attributes") {
		eaMap := d.Get("extensible_attributes").(map[string]interface{})
		if len(eaMap) > 0 {
			eas, err := createExtensibleAttributesFromJSON(client, eaMap)
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
				return diags
			}
			addressRange.ExtensibleAttributes = &eas
		}
		if client.OrchestratorEAs != nil && len(*client.OrchestratorEAs) > 0 {
			for k, v := range *client.OrchestratorEAs {
				(*addressRange.ExtensibleAttributes)[k] = v
			}
		}
	}
	changedRange, err := client.UpdateRange(d.Id(), addressRange)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if d.Get("restart_if_needed").(bool) && addressRange.Member != nil {
		err = client.RestartServices(d.Get("grid_ref").(string), infoblox.GridServiceRestartRequest{
			RestartOption: "RESTART_IF_NEEDED",
			Services:      []string{"DHCP"},
			Members:       []string{addressRange.Member.Hostname},
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

	d.SetId(changedRange.Ref)
	return resourceRangeRead(ctx, d, m)
}

func resourceRangeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics
	ref := d.Id()

	addressRange, err := convertResourceDataToRange(client, d)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	err = client.DeleteRange(ref)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.Get("restart_if_needed").(bool) && addressRange.Member != nil {
		err = client.RestartServices(d.Get("grid_ref").(string), infoblox.GridServiceRestartRequest{
			RestartOption: "RESTART_IF_NEEDED",
			Services:      []string{"DHCP"},
			Members:       []string{addressRange.Member.Hostname},
		})
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
	}

	if diags.HasError() {
		return diags
	}

	return diags
}
