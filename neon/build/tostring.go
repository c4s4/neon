package build

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// PropertyToString returns a string representation of given property:
// - object: the body of the property as an interface
// - quotes: tells if we want quotes around strings
// Return:
// - string representation of the property
// - an error if something went wrong
func PropertyToString(object interface{}, quotes bool) (string, error) {
	switch value := object.(type) {
	case bool:
		return strconv.FormatBool(value), nil
	case string:
		if quotes {
			return "\"" + value + "\"", nil
		}
		return value, nil
	case int:
		return strconv.Itoa(value), nil
	case int32:
		return strconv.Itoa(int(value)), nil
	case int64:
		return strconv.Itoa(int(value)), nil
	case float64:
		return strconv.FormatFloat(value, 'g', -1, 64), nil
	default:
		if value == nil {
			return "null", nil
		}
		switch reflect.TypeOf(object).Kind() {
		case reflect.Slice:
			str, err := sliceToString(object, quotes)
			if err != nil {
				return "", err
			}
			return str, nil
		case reflect.Map:
			str, err := mapToString(object, quotes)
			if err != nil {
				return "", err
			}
			return str, nil
		default:
			return "", fmt.Errorf("no serializer for type '%T'", object)
		}
	}
}

func sliceToString(object interface{}, quotes bool) (string, error) {
	slice := reflect.ValueOf(object)
	elements := make([]string, slice.Len())
	for index := 0; index < slice.Len(); index++ {
		str, err := PropertyToString(slice.Index(index).Interface(), quotes)
		if err != nil {
			return "", err
		}
		elements[index] = str
	}
	return "[" + strings.Join(elements, ", ") + "]", nil
}

func mapToString(object interface{}, quotes bool) (string, error) {
	dict := reflect.ValueOf(object)
	elements := make(map[string]string)
	var keys []string
	for _, key := range dict.MapKeys() {
		value := dict.MapIndex(key)
		keyStr, err := PropertyToString(key.Interface(), quotes)
		if err != nil {
			return "", err
		}
		keys = append(keys, keyStr)
		valueStr, err := PropertyToString(value.Interface(), quotes)
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
	return "{" + strings.Join(pairs, ", ") + "}", nil
}
