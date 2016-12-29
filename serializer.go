package main

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func Serialize(object interface{}) (string, error) {
	switch value := object.(type) {
	case bool:
		return strconv.FormatBool(value), nil
	case string:
		return "\"" + value + "\"", nil
	case int:
		return strconv.Itoa(value), nil
	case float64:
		return strconv.FormatFloat(value, 'g', -1, 64), nil
	default:
		switch reflect.TypeOf(object).Kind() {
		case reflect.Slice:
			slice := reflect.ValueOf(object)
			elements := make([]string, slice.Len())
			for index := 0; index < slice.Len(); index++ {
				str, err := Serialize(slice.Index(index).Interface())
				if err != nil {
					return "", err
				}
				elements[index] = str
			}
			return "[" + strings.Join(elements, ", ") + "]", nil
		case reflect.Map:
			dict := reflect.ValueOf(object)
			elements := make(map[string]string)
			var keys []string
			for _, key := range dict.MapKeys() {
				value := dict.MapIndex(key)
				keyStr, err := Serialize(key.Interface())
				if err != nil {
					return "", err
				}
				keys = append(keys, keyStr)
				valueStr, err := Serialize(value.Interface())
				if err != nil {
					return "", err
				}
				elements[keyStr] = valueStr
			}
			sort.Strings(keys)
			pairs := make([]string, len(keys))
			for index, key := range keys {
				pairs[index] = key + ": " + elements[key]
			}
			return "[" + strings.Join(pairs, ", ") + "]", nil
		default:
			return "", fmt.Errorf("no serializer for type '%T'", object)
		}
	}
}
