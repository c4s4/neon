package util

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"unicode/utf8"
)

// GOOS is the OS name
var GOOS = runtime.GOOS

// ToString returns a string from an interface:
// - object: the string as an interface
// Return:
// - converted string
// - an error if something went wrong
func ToString(object interface{}) (string, error) {
	str := reflect.ValueOf(object)
	if str.Kind() == reflect.Interface {
		str = str.Elem()
	}
	if str.Kind() == reflect.String {
		return str.Interface().(string), nil
	}
	return "", fmt.Errorf("%#v is not a string", str)
}

// ToSliceString return interface as a slice of strings:
// - object: the slice of strings as an interface
// Return:
// - converted slice of strings
// - an error if something went wrong
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
	}
	return nil, fmt.Errorf("must be a slice of strings")
}

// ToMapStringString return interface as a map with string keys and values:
// - object: the maps of strings as an interface
// Return:
// - converted map of strings
// - an error if something went wrong
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

// ToMapStringInterface return interface as a map with string keys and
// interface values:
// - object: the maps of interfaces as an interface
// Return:
// - converted map of interfaces
// - an error if something went wrong
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

// IsMap tells if given object is a map:
// - object: object to examine
// Return: a boolean that tells if object is a map
func IsMap(object interface{}) bool {
	return reflect.ValueOf(object).Kind() == reflect.Map
}

// IsSlice tells if given object is a slice:
// - object: object to examine
// Return: a boolean that tells if object is a slice
func IsSlice(object interface{}) bool {
	return reflect.ValueOf(object).Kind() == reflect.Slice
}

// MaxLineLength returns the maximum length of given lines:
// - lines: lines to examine
// Return: maximum length of lines as an integer
func MaxLineLength(lines []string) int {
	length := 0
	for _, line := range lines {
		if utf8.RuneCountInString(line) > length {
			length = utf8.RuneCountInString(line)
		}
	}
	return length
}

// Windows tells if we are running on windows
// Return: a boolean that tells if we are running on windows
func Windows() bool {
	return GOOS == "windows"
}

// RemoveBlankLines removes blank lines of given string:
// - text: to text to process
// Return: a string without blank lines
func RemoveBlankLines(text string) string {
	regex := regexp.MustCompile("(\n\\s*)+\n")
	return regex.ReplaceAllString(text, "\n")
}
