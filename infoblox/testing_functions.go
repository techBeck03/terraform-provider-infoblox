package infoblox

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func toHclString(value interface{}, isNested bool) string {
	// Ideally, we'd use a type switch here to identify slices and maps, but we can't do that, because Go doesn't
	// support generics, and the type switch only matches concrete types. So we could match []interface{}, but if
	// a user passes in []string{}, that would NOT match (the same logic applies to maps). Therefore, we have to
	// use reflection and manually convert into []interface{} and map[string]interface{}.

	if slice, isSlice := tryToConvertToGenericSlice(value); isSlice {
		return sliceToHclString(slice)
	} else if m, isMap := tryToConvertToGenericMap(value); isMap {
		return mapToHclString(m)
	} else {
		return primitiveToHclString(value, isNested)
	}
}

// Try to convert the given value to a generic slice. Return the slice and true if the underlying value itself was a
// slice and an empty slice and false if it wasn't. This is necessary because Go is a shitty language that doesn't
// have generics, nor useful utility methods built-in. For more info, see: http://stackoverflow.com/a/12754757/483528
func tryToConvertToGenericSlice(value interface{}) ([]interface{}, bool) {
	reflectValue := reflect.ValueOf(value)
	if reflectValue.Kind() != reflect.Slice {
		return []interface{}{}, false
	}

	genericSlice := make([]interface{}, reflectValue.Len())

	for i := 0; i < reflectValue.Len(); i++ {
		genericSlice[i] = reflectValue.Index(i).Interface()
	}

	return genericSlice, true
}

// Try to convert the given value to a generic map. Return the map and true if the underlying value itself was a
// map and an empty map and false if it wasn't. This is necessary because Go is a shitty language that doesn't
// have generics, nor useful utility methods built-in. For more info, see: http://stackoverflow.com/a/12754757/483528
func tryToConvertToGenericMap(value interface{}) (map[string]interface{}, bool) {
	reflectValue := reflect.ValueOf(value)
	if reflectValue.Kind() != reflect.Map {
		return map[string]interface{}{}, false
	}

	reflectType := reflect.TypeOf(value)
	if reflectType.Key().Kind() != reflect.String {
		return map[string]interface{}{}, false
	}

	genericMap := make(map[string]interface{}, reflectValue.Len())

	mapKeys := reflectValue.MapKeys()
	for _, key := range mapKeys {
		genericMap[key.String()] = reflectValue.MapIndex(key).Interface()
	}

	return genericMap, true
}

// Convert a slice to an HCL string. See ToHclString for details.
func sliceToHclString(slice []interface{}) string {
	hclValues := []string{}

	for _, value := range slice {
		hclValue := toHclString(value, true)
		hclValues = append(hclValues, hclValue)
	}

	return fmt.Sprintf("[%s]", strings.Join(hclValues, ", "))
}

// Convert a map to an HCL string. See ToHclString for details.
func mapToHclString(m map[string]interface{}) string {
	keyValuePairs := []string{}

	for key, value := range m {
		var keyValuePair string
		if _, isMap := tryToConvertToGenericMap(value); isMap {
			keyValuePair = fmt.Sprintf(`%s %s`, key, toHclString(value, true))
		} else {
			keyValuePair = fmt.Sprintf(`%s = %s`, key, toHclString(value, true))
		}
		keyValuePairs = append(keyValuePairs, keyValuePair)
	}

	return fmt.Sprintf("{\n%s\n}", strings.Join(keyValuePairs, "\n"))
}

// Convert a primitive, such as a bool, int, or string, to an HCL string. If this isn't a primitive, force its value
// using Sprintf. See ToHclString for details.
func primitiveToHclString(value interface{}, isNested bool) string {
	if value == nil {
		return "null"
	}

	switch v := value.(type) {

	case bool:
		return strconv.FormatBool(v)

	case string:
		// If string is nested in a larger data structure (e.g. list of string, map of string), ensure value is quoted
		if isNested {
			return fmt.Sprintf("\"%v\"", v)
		}

		return fmt.Sprintf("%v", v)

	default:
		return fmt.Sprintf("%v", v)
	}
}

func testAccCheckTestSliceVals(resourceName string, key string, expected []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		v, ok := rs.Primary.Attributes[fmt.Sprintf("%s.#", key)]
		if !ok {
			return fmt.Errorf("%s: Attribute '%s.#' not found", resourceName, key)
		}
		testCount, _ := strconv.Atoi(v)
		if testCount == 0 {
			return fmt.Errorf("No entries found in state for key: %s", key)
		}

		var sv []string
		for i := 0; i < testCount; i++ {
			sv = append(sv, rs.Primary.Attributes[fmt.Sprintf("%s.%d", key, i)])
		}

		diff := sliceDiff(expected, sv, true)
		if len(diff) > 0 {
			return fmt.Errorf("Set values: %s do not match expected value: %s", sv, expected)
		}
		return nil
	}
}
func sliceDiff(slice1 []string, slice2 []string, bidirectional bool) []string {
	var diff []string

	var loopCount int
	if bidirectional {
		loopCount = 2
	} else {
		loopCount = 1
	}

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < loopCount; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}

// composeConfig can be called to concatenate multiple strings to build test configurations
func composeConfig(config ...string) string {
	var str strings.Builder

	for _, conf := range config {
		str.WriteString(conf)
	}

	return str.String()
}
