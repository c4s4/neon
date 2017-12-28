package build

import (
	"reflect"
	"fmt"
	"strings"
)

// Type for task arguments as parsed in build file
type TaskArgs map[string]interface{}

// Validate task arguments against task arguments definition
// - args: task arguments parsed in build file
// - params: instance of the task arguments type
// Return: an error (detailing the fault) if arguments are illegal
// NOTE: supported tags in argument types are:
// - optional: field might not be provided
func ValidateTaskArgs(args TaskArgs, params interface{}) error {
	if reflect.TypeOf(params).Kind() != reflect.Ptr {
		return fmt.Errorf("params must be a pointer on a struct")
	}
	st := reflect.TypeOf(params).Elem()
	if st.Kind() != reflect.Struct {
		return fmt.Errorf("params must be a pointer on a struct")
	}
	for i:=0; i<st.NumField(); i++ {
		field := st.Field(i)
		// check field is not missing
		argName := strings.ToLower(field.Name)
		if _, ok := args[argName]; !ok {
			if !FieldIs(field, "optional") {
				return fmt.Errorf("missing mandatory field '%s'", argName)
			}
		}
		// check field type
		value := args[argName]
		valueType := reflect.TypeOf(value)
		if field.Type != valueType {
			if !FieldIs(field, "optional") {
				return fmt.Errorf("field '%s' must be of type '%s' ('%s' provided)", argName, field.Type, valueType)
			}
		}
	}
	return nil
}

// Evaluate task arguments in given context to fill empty arguments
// - args: task arguments parsed in build file
// - params: pointer to the instance of the task arguments type to fill
// - context: the build context to evaluate arguments into
// Return: an error if something went wrong
func EvaluateTaskArgs(args TaskArgs, params interface{}, context *Context) error {
	st := reflect.TypeOf(params).Elem()
	value := reflect.ValueOf(params).Elem()
	for i:=0; i<value.NumField(); i++ {
		name := strings.ToLower(st.Field(i).Name)
		if args[name] != nil {
			value.Field(i).Set(reflect.ValueOf(args[name]))
		}
	}
	return nil
}

// FieldIs tells if given field tag contains quality
// - field: the struct field
// - quality: the tested quality (such as "optional")
func FieldIs(field reflect.StructField, quality string) bool {
	tag := string(field.Tag)
	qualities := strings.Split(tag, " ")
	for _, q := range qualities {
		if q == quality {
			return true
		}
	}
	return false
}
