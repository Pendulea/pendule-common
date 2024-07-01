package pcommon

import (
	"fmt"
	"reflect"
	"regexp"
	"unicode"
)

func ContainsDigit(s string) bool {
	matched, _ := regexp.MatchString(`[0-9]`, s)
	return matched
}

func ChunkString(s string, n int) []string {
	if n <= 0 {
		return []string{}
	}

	var result []string
	runes := []rune(s) // Convert the string to runes to handle Unicode characters properly

	for i := 0; i < len(runes); i += n {
		end := i + n
		if end > len(runes) {
			end = len(runes)
		}
		result = append(result, string(runes[i:end]))
	}

	return result
}

func Sort[T int64 | int](slice []T, desc bool) []T {
	ret := make([]T, len(slice))
	copy(ret, slice)
	if desc {
		for i := 0; i < len(ret); i++ {
			for j := i + 1; j < len(ret); j++ {
				if ret[i] < ret[j] {
					ret[i], ret[j] = ret[j], ret[i]
				}
			}
		}
	} else {
		for i := 0; i < len(ret); i++ {
			for j := i + 1; j < len(ret); j++ {
				if ret[i] > ret[j] {
					ret[i], ret[j] = ret[j], ret[i]
				}
			}
		}
	}
	return ret
}

func getFieldByJSONTag(item Data, jsonTag string) (interface{}, bool) {
	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Check fields in Quantity struct
	quantityVal := v.Field(0)
	quantityType := quantityVal.Type()
	for i := 0; i < quantityType.NumField(); i++ {
		field := quantityType.Field(i)
		tag := field.Tag.Get("json")
		if tag == jsonTag {
			return quantityVal.Field(i).Interface(), true
		}
	}

	// Check fields in QuantityTime struct itself
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		if tag == jsonTag {
			return v.Field(i).Interface(), true
		}
	}

	return nil, false
}

func filterToMap(data DataList, fields []ColumnName) ([]map[ColumnName]interface{}, error) {
	var filteredData []map[ColumnName]interface{}

	for _, item := range data.Map() {
		mappedItem := make(map[ColumnName]interface{})

		for _, field := range fields {
			if field == "time" {
				mappedItem[field] = reflect.ValueOf(item).FieldByName("Time").Interface()
				continue
			}
			v, ok := getFieldByJSONTag(item, string(field))
			if !ok {
				return nil, fmt.Errorf("field %s not found", field)
			} else {
				mappedItem[field] = v
			}
		}
		filteredData = append(filteredData, mappedItem)
	}

	return filteredData, nil
}

func isAlphanumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
