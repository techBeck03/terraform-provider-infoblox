package infoblox

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
	"github.com/tidwall/gjson"
)

var (
	validEAKeys = []string{
		"value",
		"type",
		"inheritance_operation",
		"descendants_action",
	}
	validEATypes = []string{
		"STRING",
		"ENUM",
		"EMAIL",
		"URL",
		"DATE",
		"INTEGER",
	}
	validInheritanceOperations = []string{
		"INHERIT",
		"DELETE",
		"OVERRIDE",
	}
	validOptionDeleteEAValues = []string{
		"REMOVE",
		"RETAIN",
	}
	validOptionWithEAValues = []string{
		"CONVERT",
		"INHERIT",
		"RETAIN",
	}
	validOptionWithoutEAValues = []string{
		"INHERIT",
		"NOT_INHERIT",
	}
)

// func validateEa(eaMap map[string]interface{}) (diags diag.Diagnostics) {
func validateEa(i interface{}, p cty.Path) (diags diag.Diagnostics) {
	for key, v := range i.(map[string]interface{}) {
		parsed := gjson.Parse(v.(string))
		for k := range parsed.Value().(map[string]interface{}) {
			check := stringInSlice(validEAKeys, []string{k}, p[0].(cty.GetAttrStep).Name)
			if check.HasError() {
				diags = append(diags, check...)
			}
		}
		if parsed.Get("type").Exists() && stringInSlice(validEATypes, []string{parsed.Get("type").String()}, key).HasError() {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Invalid value for extensible attribute: %s", key),
				Detail:   fmt.Sprintf("Expected `type` to be one of %s but found %s", strings.Join(validEATypes, ", "), parsed.Get("type").String()),
			})
		}
		if parsed.Get("inheritance_operation").Exists() && stringInSlice(validInheritanceOperations, []string{parsed.Get("inheritance_operation").String()}, key).HasError() {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Invalid value for extensible attribute: %s", key),
				Detail:   fmt.Sprintf("Expected inheritance_operation to be one of %s but found %s", strings.Join(validInheritanceOperations, ", "), parsed.Get("inheritance_operation").String()),
			})
		}
		if parsed.Get("descendants_action").Exists() {
			var descendantsAction infoblox.DescendantsAction
			err := json.Unmarshal([]byte(parsed.Get("descendants_action").Str), &descendantsAction)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("Invalid value for extensible attribute: %s", key),
					Detail:   fmt.Sprint(err),
				})
			}
			if descendantsAction.OptionDeleteEA != "" && stringInSlice(validOptionDeleteEAValues, []string{descendantsAction.OptionDeleteEA}, key).HasError() {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("Invalid value for extensible attribute: %s", key),
					Detail:   fmt.Sprintf("Expected option_delete_ea to be one of %s but found %s", strings.Join(validOptionDeleteEAValues, ", "), descendantsAction.OptionDeleteEA),
				})
			}
			if descendantsAction.OptionWithEA != "" && stringInSlice(validOptionWithEAValues, []string{descendantsAction.OptionWithEA}, key).HasError() {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("Invalid value for extensible attribute: %s", key),
					Detail:   fmt.Sprintf("Expected option_with_ea to be one of %s but found %s", strings.Join(validOptionWithEAValues, ", "), descendantsAction.OptionWithEA),
				})
			}
			if descendantsAction.OptionWithoutEA != "" && stringInSlice(validOptionWithoutEAValues, []string{descendantsAction.OptionWithoutEA}, key).HasError() {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("Invalid value for extensible attribute: %s", key),
					Detail:   fmt.Sprintf("Expected option_without_ea to be one of %s but found %s", strings.Join(validOptionWithoutEAValues, ", "), descendantsAction.OptionWithoutEA),
				})
			}
		}
	}
	return diags
}

func stringInSlice(valid []string, test []string, subject string) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, t := range test {
		matchFlag := false
		for _, v := range valid {
			if v == t {
				matchFlag = true
				break
			}
		}
		if !matchFlag {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Invalid value entered for: %s", subject),
				Detail:   fmt.Sprintf("%s is not one of supported values: %s", t, strings.Join(valid[:], ", ")),
			})
		}
	}
	return diags
}

func eaSuppressDiff(k, old, new string, d *schema.ResourceData) bool {
	areEqual, err := areEqualJSON(old, new)
	if err != nil {
		return false
	}
	return areEqual
}

func areEqualJSON(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 2 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}
