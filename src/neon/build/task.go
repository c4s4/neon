package build

import (
	"reflect"
	"fmt"
	"strings"
	"neon/util"
)

// character for expressions
const (
	// character that start expressions
	CHAR_EXPRESSION = "="
	// field might not be provided
	FIELD_OPTIONAL = "optional"
	// field is a file name that is expanded for user home
	FIELD_FILE = "file"
	// field is an expression
	FIELD_EXPRESSION = "expression"
)

// Map that gives constructor for given task name
var TaskMap map[string]TaskDesc = make(map[string]TaskDesc)

func AddTask(task TaskDesc) {
	if _, ok := TaskMap[task.Name]; ok {
		panic(fmt.Errorf("task '%s' already defined", task.Name))
	}
	TaskMap[task.Name] = task
}

// A task descriptor is made of a task constructor and an help string
type TaskDesc struct {
	Name string
	Args reflect.Type
	Func TaskFunc
	Help string
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
func ValidateTaskArgs(args TaskArgs, typ reflect.Type) error {
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("params must be a pointer on a struct")
	}
	for i:=0; i<typ.NumField(); i++ {
		field := typ.Field(i)
		// check field is not missing
		argName := strings.ToLower(field.Name)
		if _, ok := args[argName]; !ok {
			if !FieldIs(field, FIELD_OPTIONAL) {
				return fmt.Errorf("missing mandatory field '%s'", argName)
			}
		}
		// check field type
		value := args[argName]
		valueType := reflect.TypeOf(value)
		if !(field.Type == valueType || (value == nil && FieldIs(field, FIELD_OPTIONAL))) {
			// if expression deffer type check after evaluation
			if !(valueType.Kind() == reflect.String &&
				 (IsExpression(value.(string)) || FieldIs(field, FIELD_EXPRESSION))) {
				return fmt.Errorf("field '%s' must be of type '%s' ('%s' provided)", argName, field.Type, valueType)
			}
		}
	}
	return nil
}

// Evaluate task arguments in given context to fill empty arguments
// - args: task arguments parsed in build file
// - typ: the type of the arguments
// - context: the build context to evaluate arguments into
// Return:
// - result: as an interface{}
// - error: if something went wrong
func EvaluateTaskArgs(args TaskArgs, typ reflect.Type, context *Context) (interface{}, error) {
	var err error
	value := reflect.New(typ).Elem()
	for i:=0; i<value.NumField(); i++ {
		name := strings.ToLower(typ.Field(i).Name)
		if args[name] != nil {
			val := args[name]
			field := typ.Field(i)
			// evaluate expressions in context
			if reflect.TypeOf(val).Kind() == reflect.String &&
				(IsExpression(args[name].(string)) || FieldIs(field, FIELD_EXPRESSION)) {
				str := args[name].(string)
				if IsExpression(str) {
					str = str[1:]
				}
				val, err = context.EvaluateExpression(str)
				if err != nil {
					return nil, err
				}
				expected := field.Type
				actual := reflect.TypeOf(val)
				if actual != expected {
					return nil, fmt.Errorf("bad expression return type, expected '%s' but '%s' was returned",
						expected.Name(), actual.Name())
				}
			}
			// evaluate strings to replace "={expression}" with its value
			if reflect.TypeOf(val).Kind() == reflect.String {
				str := args[name].(string)
				// replace '\=' with '='
				if strings.HasPrefix(str, `\`+CHAR_EXPRESSION) {
					str = str[1:]
				}
				// evaluate string
				str, err = context.EvaluateString(str)
				if err != nil {
					return nil, err
				}
				// expand home if field tagged 'file'
				if FieldIs(field, FIELD_FILE) {
					str = util.ExpandUserHome(str)
				}
				val = str
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

// IsExpression tells if given string is an expression
// - s: the string to test
// Return: a bool that tells if the string is an expression
func IsExpression(s string) bool {
	return s[0:1] == CHAR_EXPRESSION && s[1:2] != "{"
}
