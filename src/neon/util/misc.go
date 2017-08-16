package util

import (
	"fmt"
	"net"
	"reflect"
	"time"
	"unicode/utf8"
	"runtime"
	"testing"
	"regexp"
)

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
			value := slice.Index(i)
			if value.Kind() != reflect.String {
				fmt.Errorf("must be a slice of strings")
			}
			result[i] = value.Interface().(string)
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
		keyString := key.Interface().(string)
		valueString := value.MapIndex(key).Interface().(string)
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
		keyString := key.Interface().(string)
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

// Return a reflect.Value as an interface
func ValueToInterface(value reflect.Value) interface{} {
	if value.IsValid() {
		kind := value.Kind()
		if kind == reflect.Slice {
			result := make([]interface{}, value.Len())
			for i := 0; i < value.Len(); i++ {
				result[i] = ValueToInterface(value.Index(i))
			}
			return result
		} else if kind == reflect.Map {
			result := make(map[interface{}]interface{})
			for _, key := range value.MapKeys() {
				keyInterface := ValueToInterface(key)
				valueInterface := ValueToInterface(value.MapIndex(key))
				result[keyInterface] = valueInterface
			}
			return result
		} else {
			return value.Interface()
		}
	} else {
		return nil
	}
}

// Return the maximum length of given lines
func MaxLength(lines []string) int {
	length := 0
	for _, line := range lines {
		if utf8.RuneCountInString(line) > length {
			length = utf8.RuneCountInString(line)
		}
	}
	return length
}

// Run a TCP server on given port to ensure that a single instance is running
// on a machine. Fails if another instance is already running.
func Singleton(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	go func() {
		for {
			listener.Accept()
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return nil
}

// Tells if we are running on windows
func Windows() bool {
	return 	runtime.GOOS == "windows"
}

// Make an assertion for testing purpose
func Assert(actual, expected string, t *testing.T) {
	if actual != expected {
		t.Errorf("actual \"%s\" != expected \"%s\"", actual, expected)
	}
}

// RemoveBlankLines removes blank lines of given string.
func RemoveBlankLines(text string) string {
	regex := regexp.MustCompile("(\n\\s*)+\n")
	return regex.ReplaceAllString(text, "\n")
}