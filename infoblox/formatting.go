package infoblox

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/techBeck03/infoblox-go-sdk"
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

func createExtensibleAttributesFromJSON(client *infoblox.Client, eaMap map[string]interface{}) (eas infoblox.ExtensibleAttribute, err error) {
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
		var eaValue infoblox.ExtensibleAttributeJSONMapValue
		json.Unmarshal([]byte(v.(string)), &eaValue)
		var ea infoblox.ExtensibleAttributeValue
		switch eaValue.Type {
		case "STRING":
			ea.Value = eaValue.Value.(string)
		case "ENUM":
			ea.Value = eaValue.Value.(string)
		case "EMAIL":
			ea.Value = eaValue.Value.(string)
		case "URL":
			ea.Value = eaValue.Value.(string)
		case "DATE":
			ea.Value = eaValue.Value.(string)
		case "INTEGER":
			ea.Value = eaValue.Value.(int)
		}
		if strings.Contains(v.(string), "inheritance_source") {
			ea.InheritanceSource = eaValue.InheritanceSource
		}
		eas[k] = ea
	}

	return eas, err
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
