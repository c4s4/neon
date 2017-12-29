package build

import (
	"reflect"
	"fmt"
	"strings"
)

// character for expressions
const CHAR_EXPRESSION = `=`

// Map that gives constructor for given task name
var TaskMap map[string]TaskDesc = make(map[string]TaskDesc)

// A task descriptor is made of a task constructor and an help string
type TaskDesc struct {
	Args        reflect.Type
	Func        TaskFunc
	Help        string
}

// Type for task arguments as parsed in build file
type TaskArgs map[interface{}]interface{}

// Type for task function called to run task
type TaskFunc func(ctx *Context, args interface{}) error

// Validate task arguments against task arguments definition
// - args: task arguments parsed in build file
// - typ: the type of the arguments
// Return: an error (detailing the fault) if arguments are illegal
// NOTE: supported tags in argument types are:
// - optional: field might not be provided
func ValidateTaskArgs(args TaskArgs, typ reflect.Type) error {
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("params must be a pointer on a struct")
	}
	for i:=0; i<typ.NumField(); i++ {
		field := typ.Field(i)
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
func EvaluateTaskArgs(args TaskArgs, typ reflect.Type, context *Context) (interface{}, error) {
	var err error
	value := reflect.New(typ).Elem()
	for i:=0; i<value.NumField(); i++ {
		name := strings.ToLower(typ.Field(i).Name)
		if args[name] != nil {
			val := args[name]
			// evaluate strings and expressions
			if reflect.TypeOf(args[name]).Kind() == reflect.String {
				str := args[name].(string)
				if strings.HasPrefix(str, CHAR_EXPRESSION) {
					// if starts with ':' this is an expression
					val, err = context.EvaluateExpression(str[1:])
					if err != nil {
						return nil, err
					}
				} else {
					// if doesn't start with ':' this is a string
					if strings.HasPrefix(str, `\`+CHAR_EXPRESSION) {
						str = str[1:]
					}
					val, err = context.EvaluateString(str)
					if err != nil {
						return nil, err
					}
				}
			}
			// put value in params
			value.Field(i).Set(reflect.ValueOf(val))
		}
	}
	return value.Interface(), nil
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
