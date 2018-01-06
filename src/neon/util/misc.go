package util

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"unicode/utf8"
)

// ToString returns a string from an interface
func ToString(object interface{}) (string, error) {
	str := reflect.ValueOf(object)
	if str.Kind() == reflect.Interface {
		str = str.Elem()
	}
	if str.Kind() == reflect.String {
		return str.Interface().(string), nil
	} else {
		return "", fmt.Errorf("%#v is not a string", str)
	}
}

// Return interface as a list of interfaces
func ToList(object interface{}) ([]interface{}, error) {
	slice := reflect.ValueOf(object)
	if slice.Kind() == reflect.Slice {
		result := make([]interface{}, slice.Len())
		for i := 0; i < slice.Len(); i++ {
			result[i] = slice.Index(i).Interface()
		}
		return result, nil
	} else {
		return nil, fmt.Errorf("must be a list")
	}
}

// ToSliceString return interface as a slice of strings.
func ToSliceString(object interface{}) ([]string, error) {
	slice := reflect.ValueOf(object)
	if slice.Kind() == reflect.Slice {
		result := make([]string, slice.Len())
		for i := 0; i < slice.Len(); i++ {
			value := slice.Index(i).Interface()
			var err error
			if result[i], err = ToString(value); err != nil {
				return nil, err
			}
		}
		return result, nil
	} else {
		return nil, fmt.Errorf("must be a slice of strings")
	}
}

// Return interface as a map with string keys and values
func ToMapStringString(object interface{}) (map[string]string, error) {
	value := reflect.ValueOf(object)
	if value.Kind() != reflect.Map {
		return nil, fmt.Errorf("object is not a map")
	}
	result := make(map[string]string)
	for _, key := range value.MapKeys() {
		keyString, err := ToString(key.Interface())
		if err != nil {
			return nil, err
		}
		valueString, err := ToString(value.MapIndex(key).Interface())
		if err != nil {
			return nil, err
		}
		result[keyString] = valueString
	}
	return result, nil
}

// ToMapStringInterface return interface as a map with string keys and interface
// values.
func ToMapStringInterface(object interface{}) (map[string]interface{}, error) {
	value := reflect.ValueOf(object)
	if value.Kind() != reflect.Map {
		return nil, fmt.Errorf("object is not a map")
	}
	result := make(map[string]interface{})
	for _, key := range value.MapKeys() {
		keyString, err := ToString(key.Interface())
		if err != nil {
			return nil, err
		}
		valueInterface := value.MapIndex(key).Interface()
		result[keyString] = valueInterface
	}
	return result, nil
}

// IsMap tells if given object is a map
func IsMap(object interface{}) bool {
	return reflect.ValueOf(object).Kind() == reflect.Map
}

// IsString tells if given object is a string
func IsString(object interface{}) bool {
	return reflect.ValueOf(object).Kind() == reflect.String
}

// IsSlice tells if given object is a slice.
func IsSlice(object interface{}) bool {
	return reflect.ValueOf(object).Kind() == reflect.Slice
}

// Return the maximum length of given lines
func MaxLineLength(lines []string) int {
	length := 0
	for _, line := range lines {
		if utf8.RuneCountInString(line) > length {
			length = utf8.RuneCountInString(line)
		}
	}
	return length
}

// Tells if we are running on windows
func Windows() bool {
	return runtime.GOOS == "windows"
}

// RemoveBlankLines removes blank lines of given string.
func RemoveBlankLines(text string) string {
	regex := regexp.MustCompile("(\n\\s*)+\n")
	return regex.ReplaceAllString(text, "\n")
}
