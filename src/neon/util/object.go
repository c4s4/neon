package util

import (
	"fmt"
	"reflect"
	"sort"
)

type Object map[string]interface{}

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
	} else {
		return nil, err
	}
}

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
	} else {
		return nil, fmt.Errorf("field must be a list")
	}
}

func (object Object) GetListStrings(field string) ([]string, error) {
	thing, ok := object[field]
	if !ok {
		return make([]string, 0), nil
	}
	err := fmt.Errorf("field must be a map with string keys")
	slice := reflect.ValueOf(thing)
	if slice.Kind() == reflect.Slice {
		result := make([]string, slice.Len())
		for i := 0; i < slice.Len(); i++ {
			value := slice.Index(i)
			str, ok := value.Interface().(string)
			if !ok {
				return nil, err
			}
			result[i] = str
		}
		return result, nil
	} else {
		return nil, err
	}
}

func (object Object) GetListStringsOrString(field string) ([]string, error) {
	thing, ok := object[field]
	if !ok {
		return make([]string, 0), nil
	}
	err := fmt.Errorf("field must be a map with string keys")
	slice := reflect.ValueOf(thing)
	if slice.Kind() == reflect.Slice {
		result := make([]string, slice.Len())
		for i := 0; i < slice.Len(); i++ {
			value := slice.Index(i)
			str, ok := value.Interface().(string)
			if !ok {
				return nil, err
			}
			result[i] = str
		}
		return result, nil
	} else if slice.Kind() == reflect.String {
		str, _ := slice.Interface().(string)
		return []string{str}, nil
	} else {
		return nil, err
	}
}

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

func (object Object) CheckFields(fields []string) error {
	for entry, _ := range object {
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

func (object Object) Copy() Object {
	copy := make(map[string]interface{})
	for name, value := range object {
		copy[name] = value
	}
	return copy
}

func (object Object) Fields() []string {
	fields := make([]string, len(object))
	index := 0
	for name, _ := range object {
		fields[index] = name
		index++
	}
	sort.Strings(fields)
	return fields
}

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
