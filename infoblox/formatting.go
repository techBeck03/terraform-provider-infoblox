package infoblox

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/techBeck03/infoblox-go-sdk"
	"github.com/tidwall/gjson"
)

func prettyPrint(object interface{}) {
	output, _ := json.MarshalIndent(object, "", "    ")
	log.Printf("%s", string(output))
}

func newBool(b bool) *bool {
	return &b
}

func newExtensibleAttribute(ea infoblox.ExtensibleAttribute) *infoblox.ExtensibleAttribute {
	return &ea
}

// Contains tells whether a contains x.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// Keys returns the keys of a map[string]interface{} var as a slice
func Keys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func createExtensibleAttributesFromJSON(eaMap map[string]interface{}) (eas infoblox.ExtensibleAttribute, err error) {
	eas = infoblox.ExtensibleAttribute{}
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = fmt.Errorf(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()
	for k, v := range eaMap {
		parsed := gjson.Parse(v.(string))
		var ea infoblox.ExtensibleAttributeValue
		switch parsed.Get("type").Str {
		case "STRING":
			ea.Value = parsed.Get("value").Str
		case "ENUM":
			ea.Value = parsed.Get("value").Str
		case "EMAIL":
			ea.Value = parsed.Get("value").Str
		case "URL":
			ea.Value = parsed.Get("value").Str
		case "DATE":
			ea.Value = parsed.Get("value").Str
		case "INTEGER":
			ea.Value = parsed.Get("value").Int()
		}
		if parsed.Get("inheritance_source").Exists() {
			var inheritanceSource infoblox.InheritanceSource
			json.Unmarshal([]byte(parsed.Get("inheritance_source").String()), &inheritanceSource)
			ea.InheritanceSource = &inheritanceSource
		}
		ea.InheritanceOperation = parsed.Get("inheritance_operation").Str
		if parsed.Get("descendants_action").Exists() {
			var descendantsAction infoblox.DescendantsAction
			json.Unmarshal([]byte(parsed.Get("descendants_action").String()), &descendantsAction)
			ea.DescendantsAction = &descendantsAction
		}

		eas[k] = ea
	}

	return eas, err
}

func handleExtenisbleAttributesInheritanceValues(eas *infoblox.ExtensibleAttribute, d *schema.ResourceData) (infoblox.ExtensibleAttribute, error) {
	eaMap := d.Get("extensible_attributes").(map[string]interface{})
	configuredEAs, err := createExtensibleAttributesFromJSON(eaMap)
	newEas := make(infoblox.ExtensibleAttribute)
	if err != nil {
		return newEas, err
	}
	for k, v := range *eas {
		if ea, ok := configuredEAs[k]; ok {
			v.DescendantsAction = ea.DescendantsAction
			v.InheritanceOperation = ea.InheritanceOperation
		}
		newEas[k] = v
	}

	return newEas, nil
}

func isEmpty(object interface{}) bool {
	if object == nil {
		return true
	} else if object == "" {
		return true
	} else if object == false {
		return true
	}

	if reflect.ValueOf(object).Kind() == reflect.Map {
		if reflect.New(reflect.TypeOf(object)).Elem().Len() == 0 {
			return true
		}

	}

	if reflect.ValueOf(object).Kind() == reflect.Struct {
		empty := reflect.New(reflect.TypeOf(object)).Elem().Interface()
		if reflect.DeepEqual(object, empty) {
			return true
		}
	}
	return false
}

func remove(s []string, r string, createNew bool) []string {
	for i, v := range s {
		if v == r {
			if createNew {
				var m []string
				for j := 0; j < i; j++ {
					m = append(m, s[j])
				}
				for j := i + 1; j < len(s); j++ {
					m = append(m, s[j])
				}
				return m
			}
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
