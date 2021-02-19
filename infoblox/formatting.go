package infoblox

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/techBeck03/infoblox-go-sdk"
)

func prettyPrint(object interface{}) {
	output, _ := json.MarshalIndent(object, "", "    ")
	log.Printf("%s", string(output))
}

func newBool(b bool) *bool {
	return &b
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
		switch eaValue.Type {
		case "STRING":
			eas[k] = infoblox.ExtensibleAttributeValue{
				Value: eaValue.Value.(string),
			}
		case "ENUM":
			eas[k] = infoblox.ExtensibleAttributeValue{
				Value: eaValue.Value.(string),
			}
		case "EMAIL":
			eas[k] = infoblox.ExtensibleAttributeValue{
				Value: eaValue.Value.(string),
			}
		case "URL":
			eas[k] = infoblox.ExtensibleAttributeValue{
				Value: eaValue.Value.(string),
			}
		case "DATE":
			eas[k] = infoblox.ExtensibleAttributeValue{
				Value: eaValue.Value.(string),
			}
		case "INTEGER":
			eas[k] = infoblox.ExtensibleAttributeValue{
				Value: eaValue.Value.(int),
			}
		}
	}

	return eas, err
}
