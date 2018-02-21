package util

import (
	"fmt"
	"reflect"
	"sort"
)

// Object is a dictionary of interfaces
type Object map[string]interface{}

// NewObject makes an object from an interface:
// - thing: the thing to convert to an object
// Return:
// - converted object
// - an error if something went wrong
func NewObject(thing interface{}) (Object, error) {
	err := fmt.Errorf("field must be a map with string keys")
	value := reflect.ValueOf(thing)
	if value.Kind() == reflect.Map {
		result := make(map[string]interface{})
		for _, key := range value.MapKeys() {
			str, ok := key.Interface().(string)
			if !ok {
				return nil, err
			}
			result[str] = value.MapIndex(key).Interface()
		}
		return result, nil
	}
	return nil, err
}

// GetString returns an object field as a string:
// - field: name of the field to get
// Return:
// - string content of the field
// - an error if something went wrong
func (object Object) GetString(field string) (string, error) {
	value, ok := object[field]
	if !ok {
		return "", fmt.Errorf("field '%s' not found", field)
	}
	str, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("field '%s' must be a string", field)
	}
	return str, nil
}

// GetBoolean returns an object field as a boolean:
// - field: name of the field to get
// Return:
// - boolean content of the field
// - an error if something went wrong
func (object Object) GetBoolean(field string) (bool, error) {
	value, ok := object[field]
	if !ok {
		return false, fmt.Errorf("field '%s' not found", field)
	}
	boolean, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("field '%s' must be a boolean", field)
	}
	return boolean, nil
}

// GetInteger returns an object field as a integer:
// - field: name of the field to get
// Return:
// - integer content of the field
// - an error if something went wrong
func (object Object) GetInteger(field string) (int, error) {
	value, ok := object[field]
	if !ok {
		return 0, fmt.Errorf("field '%s' not found", field)
	}
	integer, ok := value.(int)
	if !ok {
		return 0, fmt.Errorf("field '%s' must be an integer", field)
	}
	return integer, nil
}

// GetList returns an object field as a slice of interfaces:
// - field: name of the field to get
// Return:
// - content of the field as a slice of interfaces
// - an error if something went wrong
func (object Object) GetList(field string) ([]interface{}, error) {
	thing, ok := object[field]
	if !ok {
		return make([]interface{}, 0), nil
	}
	slice := reflect.ValueOf(thing)
	if slice.Kind() == reflect.Slice {
		result := make([]interface{}, slice.Len())
		for i := 0; i < slice.Len(); i++ {
			result[i] = slice.Index(i).Interface()
		}
		return result, nil
	}
	return nil, fmt.Errorf("field must be a list")
}

// GetListStrings returns an object field as a slice of strings:
// - field: name of the field to get
// Return:
// - content of the field as a slice of strings
// - an error if something went wrong
func (object Object) GetListStrings(field string) ([]string, error) {
	thing, ok := object[field]
	if !ok {
		return make([]string, 0), nil
	}
	slice := reflect.ValueOf(thing)
	if slice.Kind() == reflect.Slice {
		result := make([]string, slice.Len())
		for i := 0; i < slice.Len(); i++ {
			value := slice.Index(i)
			str, err := ToString(value.Interface())
			if err != nil {
				return nil, err
			}
			result[i] = str
		}
		return result, nil
	}
	return nil, fmt.Errorf("field must be a map with string keys")
}

// GetListStringsOrString returns an object field as a slice of strings:
// - field: name of the field to get
// Return:
// - content of the field as a slice of strings
// - an error if something went wrong
func (object Object) GetListStringsOrString(field string) ([]string, error) {
	thing, ok := object[field]
	if !ok {
		return make([]string, 0), nil
	}
	slice := reflect.ValueOf(thing)
	if slice.Kind() == reflect.Slice {
		result := make([]string, slice.Len())
		for i := 0; i < slice.Len(); i++ {
			value := slice.Index(i)
			str, err := ToString(value.Interface())
			if err != nil {
				return nil, err
			}
			result[i] = str
		}
		return result, nil
	}
	if slice.Kind() == reflect.String {
		str, err := ToString(slice.Interface())
		if err != nil {
			return nil, err
		}
		return []string{str}, nil
	}
	return nil, fmt.Errorf("field must be a map with string keys")
}

// GetObject returns an object field as an object:
// - field: name of the field to get
// Return:
// - content of the field as an object
// - an error if something went wrong
func (object Object) GetObject(field string) (Object, error) {
	value, ok := object[field]
	if !ok {
		return nil, fmt.Errorf("field '%s' not found", field)
	}
	object, err := NewObject(value)
	if err != nil {
		return nil, fmt.Errorf("getting field '%s': %v", field, err)
	}
	return object, nil
}

// GetMapStringString returns an object field as a map with string keys and values:
// - field: name of the field to get
// Return:
// - content of the field map with string keys and values
// - an error if something went wrong
func (object Object) GetMapStringString(field string) (map[string]string, error) {
	value, ok := object[field]
	if !ok {
		return nil, fmt.Errorf("field '%s' not found", field)
	}
	result, err := ToMapStringString(value)
	if err != nil {
		return nil, fmt.Errorf("field must be a map string string")
	}
	return result, nil
}

// CheckFields checks that object has no field whose name is not in given list:
// - fields: list of fields to check
// Return: an error if something went wrong
func (object Object) CheckFields(fields []string) error {
	for entry := range object {
		found := false
		for _, field := range fields {
			if field == entry {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("unknown field '%s'", entry)
		}
	}
	return nil
}

// Copy returns a copy of the object
func (object Object) Copy() Object {
	copy := make(map[string]interface{})
	for name, value := range object {
		copy[name] = value
	}
	return copy
}

// Fields return fields of the object as a list of strings
func (object Object) Fields() []string {
	fields := make([]string, len(object))
	index := 0
	for name := range object {
		fields[index] = name
		index++
	}
	sort.Strings(fields)
	return fields
}

// HasField tells if object has given field:
// - field: name of the field to check
// Return: a boolean that tells if object has field
func (object Object) HasField(field string) bool {
	for name := range object {
		if name == field {
			return true
		}
	}
	return false
}

// ToMapStringString returns an object as a map with string keys and values.
// Return:
// - object as a map with string keys and values
// - an error if something went wrong
func (object Object) ToMapStringString() (map[string]string, error) {
	mapStringString := make(map[string]string)
	for name, value := range object {
		str, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("field '%s' is not a string", name)
		}
		mapStringString[name] = str
	}
	return mapStringString, nil
}
