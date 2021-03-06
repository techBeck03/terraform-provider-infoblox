package infoblox

import (
	"bytes"
	"context"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/techBeck03/go-ipmath"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

func makeEACustomDiff(arg string) func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
	return func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
		client := v.(*infoblox.Client)
		if len(*client.OrchestratorEAs) > 0 {
			old, new := diff.GetChange(arg)
			eaMap := new.(map[string]interface{})
			var eas infoblox.ExtensibleAttribute
			if len(eaMap) > 0 {
				localEAs, err := createExtensibleAttributesFromJSON(client, eaMap)
				if err != nil {
					return err
				}
				eas = localEAs
				if len(old.(map[string]interface{})) > 0 {
					oldEAs, err := createExtensibleAttributesFromJSON(client, old.(map[string]interface{}))
					if err != nil {
						return err
					}
					newEAs, err := createExtensibleAttributesFromJSON(client, new.(map[string]interface{}))
					if err != nil {
						return err
					}
					for k, v := range oldEAs {
						if v.InheritanceSource != nil && (v.Value == newEAs[k].Value || newEAs[k].Value == nil) {
							(eas)[k] = v
						}
					}
				}
			}
			for k, v := range *client.OrchestratorEAs {
				if len(eaMap) == 0 {
					eas = make(infoblox.ExtensibleAttribute)
				}
				(eas)[k] = v
			}

			finalEas, err := client.ConvertEAsToJSONString(eas)
			if err != nil {
				return err
			}
			diff.SetNew(arg, finalEas)
		}
		return nil
	}
}

func optionCustomDiff(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
	old, new := diff.GetChange("option")
	defaultFlag := false
	var leaseOption map[string]interface{}
	optionList := old.(*schema.Set).List()

	if len(optionList) > 0 {
		for _, option := range optionList {
			if option.(map[string]interface{})["code"].(int) == 51 {
				defaultFlag = true
				leaseOption = option.(map[string]interface{})
				break
			}
		}
		if defaultFlag {
			newOptions := new.(*schema.Set).List()
			newOptions = append(newOptions, leaseOption)
			diff.SetNew("option", newOptions)
		}

	}
	return nil
}

func makeAddressCompareCustomDiff(low string, high string) func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
	return func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
		if _, ok := diff.GetOk(low); ok {
			lowAddress := net.ParseIP(diff.Get(low).(string))
			highAddress := net.ParseIP(diff.Get(high).(string))

			if bytes.Compare(lowAddress, highAddress) >= 0 {
				return fmt.Errorf("IP Address `%s` must be lower than `%s`", low, high)
			}
		}

		return nil
	}
}

func rangeForceNew(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
	old, new := diff.GetChange("sequential_count")
	if _, ok := diff.GetOk("sequential_count"); ok && old == 0 {
		diff.ForceNew("sequential_count")
	} else if _, ok := diff.GetOk("start_address"); ok && old != 0 && new == 0 {
		diff.ForceNew("start_address")
	}

	return nil
}

func makeCidrContainsIPCheck(cidrArg string, ipArgs []string) func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
	return func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
		_, cidrNet, err := net.ParseCIDR(diff.Get(cidrArg).(string))
		if err != nil {
			return err
		}
		for _, ipArg := range ipArgs {
			if ip, ok := diff.GetOk(ipArg); ok {
				if cidrNet.Contains(net.ParseIP(ip.(string))) != true {
					return fmt.Errorf("Argument: %s contains ip: %s which is not in CIDR: %s", ipArg, ip.(string), cidrNet)
				}
			}
		}

		return nil
	}
}

func makeLowerThanIPCheck(lowerIP string, higherIP string) func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
	return func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
		if lowerVal, ok := diff.GetOk(lowerIP); ok {
			lowerIPObj := ipmath.IP{
				Address: net.ParseIP(lowerVal.(string)),
			}
			if lowerIPObj.GTE(net.ParseIP(diff.Get(higherIP).(string))) {
				return fmt.Errorf("`%s` must have an IP address lower than `%s`", lowerIP, higherIP)
			}
		}
		return nil
	}
}
