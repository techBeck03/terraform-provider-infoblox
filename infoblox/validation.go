package infoblox

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	validEAKeys = []string{
		"value",
		"type",
	}
)

// func validateEa(eaMap map[string]interface{}) (diags diag.Diagnostics) {
func validateEa(i interface{}, p cty.Path) (diags diag.Diagnostics) {
	log.Printf("%+v", p)
	prettyPrint(i)
	for _, v := range i.(map[string]interface{}) {
		var eaValue map[string]interface{}
		json.Unmarshal([]byte(v.(string)), &eaValue)
		for k := range eaValue {
			check := stringInSlice(validEAKeys, []string{k}, p[0].(cty.GetAttrStep).Name)
			if check.HasError() {
				diags = append(diags, check...)
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
	} else {
		return areEqual
	}
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
