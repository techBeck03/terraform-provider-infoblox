package infoblox

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

func eaCustomDiff(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
	client := v.(*infoblox.Client)
	if len(*client.OrchestratorEAs) > 0 {
		eaMap := diff.Get("extensible_attributes").(map[string]interface{})
		var eas infoblox.ExtensibleAttribute
		if len(eaMap) > 0 {
			localEAs, err := createExtensibleAttributesFromJSON(client, eaMap)
			if err != nil {
				return err
			}
			eas = localEAs
		}
		for k, v := range *client.OrchestratorEAs {
			(eas)[k] = v
		}
		finalEas, err := client.ConvertEAsToJSONString(eas)
		if err != nil {
			return err
		}
		diff.SetNew("extensible_attributes", finalEas)
	}
	return nil
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
