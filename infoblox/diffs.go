package infoblox

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/techBeck03/go-ipmath"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

func makeEACustomDiff(arg string, ignored_eas ...string) func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
	return func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
		client := v.(*infoblox.Client)
		var eas infoblox.ExtensibleAttribute
		old, new := diff.GetChange(arg)
		eaMap := new.(map[string]interface{})
		for _, ignored_ea := range ignored_eas {
			eaMap[ignored_ea] = old.(map[string]interface{})[ignored_ea]
		}
		if diff.HasChange(arg) && len(eaMap) > 0 {
			localEAs, err := createExtensibleAttributesFromJSON(eaMap)
			if err != nil {
				return err
			}
			eas = localEAs
		}
		if len(old.(map[string]interface{})) > 0 {
			oldEAs, err := createExtensibleAttributesFromJSON(old.(map[string]interface{}))
			if err != nil {
				return err
			}
			newEAs, err := createExtensibleAttributesFromJSON(new.(map[string]interface{}))
			if err != nil {
				return err
			}
			for k, v := range oldEAs {
				if v.InheritanceSource != nil && (newEAs[k].Value == nil || newEAs[k].Value == v.Value || Contains(ignored_eas, k)) {
					if eas == nil {
						eas = infoblox.ExtensibleAttribute{
							k: v,
						}
					}
					(eas)[k] = v
				}
			}
		}
		if client.OrchestratorEAs != nil && len(*client.OrchestratorEAs) > 0 {

			for k, v := range *client.OrchestratorEAs {
				if len(eaMap) == 0 {
					eas = make(infoblox.ExtensibleAttribute)
				}
				(eas)[k] = v
			}
		}
		finalEas, err := client.ConvertEAsToJSONString(eas)
		if err != nil {
			return err
		}
		diff.SetNew(arg, finalEas)
		return nil
	}
}

func makeEACustomDiffNetwork(arg string) func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
	return func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
		client := v.(*infoblox.Client)
		var eas infoblox.ExtensibleAttribute
		old, new := diff.GetChange(arg)
		eaMap := new.(map[string]interface{})
		gateway_ea := diff.Get("gateway_ea").(string)
		if gateway_ea != "" && old != nil && old.(map[string]interface{})[gateway_ea] != nil {
			eaMap[gateway_ea] = old.(map[string]interface{})[gateway_ea]
		}
		if diff.HasChange(arg) && len(eaMap) > 0 {
			localEAs, err := createExtensibleAttributesFromJSON(eaMap)
			if err != nil {
				return err
			}
			eas = localEAs
		}
		if len(old.(map[string]interface{})) > 0 {
			oldEAs, err := createExtensibleAttributesFromJSON(old.(map[string]interface{}))
			if err != nil {
				return err
			}
			newEAs, err := createExtensibleAttributesFromJSON(new.(map[string]interface{}))
			if err != nil {
				return err
			}
			for k, v := range oldEAs {
				if v.InheritanceSource != nil && (newEAs[k].Value == nil || newEAs[k].Value == v.Value || k == gateway_ea) {
					if eas == nil {
						eas = infoblox.ExtensibleAttribute{
							k: v,
						}
					}
					(eas)[k] = v
				}
			}
		}
		if client.OrchestratorEAs != nil && len(*client.OrchestratorEAs) > 0 {

			for k, v := range *client.OrchestratorEAs {
				if len(eaMap) == 0 {
					eas = make(infoblox.ExtensibleAttribute)
				}
				(eas)[k] = v
			}
		}
		finalEas, err := client.ConvertEAsToJSONString(eas)
		if err != nil {
			return err
		}
		diff.SetNew(arg, finalEas)
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

func hostRecordAddressDiff(c context.Context, diff *schema.ResourceDiff, v interface{}) error {
	old, new := diff.GetChange("ip_v4_address")
	if diff.HasChange("ip_v4_address") {
		ipAddressList := new.([]interface{})
		atLeastOneOfFields := []string{
			"ip_address",
			"network",
			"range_function_string",
		}
		addressList := new.([]interface{})
		for k, address := range ipAddressList {
			addr := address.(map[string]interface{})
			if len(old.([]interface{})) > 0 {
				matchArgs := []string{}
				for _, f := range atLeastOneOfFields {
					if addr[f] != "" {
						matchArgs = append(matchArgs, f)
					}
				}
				if len(matchArgs) == 0 {
					return fmt.Errorf("At least one of %s required for ip_v4_address", strings.Join(atLeastOneOfFields, ", "))
				} else if len(matchArgs) > 1 {
					return fmt.Errorf("Only one of %s is allowed for ip_v4_address but found %s", strings.Join(atLeastOneOfFields, ", "), strings.Join(matchArgs, ", "))
				}
			}
			if addr["ip_address"].(string) == "" && len(old.([]interface{})) > 0 && old.([]interface{})[k].(map[string]interface{})["ip_address"].(string) != "" {
				addr["ip_address"] = old.([]interface{})[k].(map[string]interface{})["ip_address"].(string)
			}
			addressList[k] = addr
		}
		diff.SetNew("ip_v4_address", addressList)
	}
	return nil
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

func makeGTIPCheck(lowerIP string, higherIP string) func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
	return func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
		if lowerVal, ok := diff.GetOk(lowerIP); ok {
			lowerIPObj := ipmath.IP{
				Address: net.ParseIP(lowerVal.(string)),
			}
			if lowerIPObj.GT(net.ParseIP(diff.Get(higherIP).(string))) {
				return fmt.Errorf("`%s` must have an IP address lower than `%s`", lowerIP, higherIP)
			}
		}
		return nil
	}
}
