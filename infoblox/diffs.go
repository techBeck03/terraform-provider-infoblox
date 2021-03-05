package infoblox

import (
	"bytes"
	"context"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
						if v.InheritanceSource != nil && v.Value == newEAs[k].Value {
							(eas)[k] = v
						}
						// e := reflect.ValueOf(&v).Elem()
						// for i := 0; i < e.NumField(); i++ {
						// 	// if e.Type().Field(i).Name == "InheritanceSource" && Contains(Keys(new.(map[string]interface{})), k) != true {

						// 	if e.Type().Field(i).Name == "InheritanceSource" && new.(map[string]interface{})), k) != true {
						// 		(eas)[k] = v
						// 	}
						// }
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
	old, _ := diff.GetChange("sequential_count")
	if _, ok := diff.GetOk("sequential_count"); ok && old == nil {
		diff.ForceNew("sequential_count")
	} else {
		diff.ForceNew("start_address")
	}

	return nil
}
