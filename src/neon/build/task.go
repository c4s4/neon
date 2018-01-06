package build

import (
	"fmt"
	"neon/util"
	"reflect"
	"strings"
)

// character for expressions
const (
	// character that start expressions
	CHAR_EXPRESSION = "="
	// expression curly brace
	CURLY_EXPRESSION = "{"
	// tag separator
	TAG_SEPARATOR = " "
	// field might not be provided
	FIELD_OPTIONAL = "optional"
	// field is a file name that is expanded for user home
	FIELD_FILE = "file"
	// field is an expression
	FIELD_EXPRESSION = "expression"
	// field should be rapped in a slice of its type
	FIELD_WRAP = "wrap"
	// field is a list of steps
	FIELD_STEPS = "steps"
	// field has a different name
	FIELD_NAME = "name"
)

// Map that gives constructor for given task name
var TaskMap = make(map[string]TaskDesc)

// AddTask adds a task in the map:
// - task: description of the task
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
// - typ: type of the arguments
// Return: an error (detailing the fault) if arguments are illegal
func ValidateTaskArgs(args TaskArgs, typ reflect.Type) error {
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("params must be a pointer on a struct")
	}
	var fields []string
	// iterate on fields of the parameters types and check argument types
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		name := strings.ToLower(field.Name)
		if field.Tag.Get(FIELD_NAME) != "" {
			name = field.Tag.Get(FIELD_NAME)
		}
		fields = append(fields, name)
		// check field is not missing
		if _, ok := args[name]; !ok {
			if !FieldIs(field, FIELD_OPTIONAL) {
				return fmt.Errorf("missing mandatory field '%s'", name)
			}
		}
		value := args[name]
		// parse steps fields
		if FieldIs(field, "steps") && value != nil {
			steps, err := NewSteps(value)
			if err != nil {
				return fmt.Errorf("parsing field '%s': %v", name, err)
			}
			args[name] = steps
		}
		// check field type
		if !CheckType(field, value) {
			return fmt.Errorf("field '%s' must be of type '%s' ('%s' provided)",
				name, field.Type, reflect.TypeOf(value))
		}
	}
	// check that we don't have unknown args
	if !(typ.NumField() == 0 && len(args) == 1) {
		for name := range args {
			found := false
			for _, n := range fields {
				if n == name {
					found = true
					continue
				}
			}
			if !found {
				return fmt.Errorf("unknown parameter '%s'", name)
			}
		}
	}
	return nil
}

// CheckType checks that given value is compatible with field type
// - field: the field of the parameters as reflect.StructField
// - value: the value of the argument
// Return: a bool that tells if type is OK
func CheckType(field reflect.StructField, value interface{}) bool {
	valueType := reflect.TypeOf(value)
	// if field is optional and argument nil, it's OK
	if FieldIs(field, FIELD_OPTIONAL) && value == nil {
		return true
	}
	// if argument is an expression it's OK whatever the type
	if valueType.Kind() == reflect.String &&
		(IsExpression(value.(string)) ||
			FieldIs(field, FIELD_EXPRESSION)) {
		return true
	}
	// if type of field is slice of the type of the argument and wrap, it's OK
	if field.Type.Kind() == reflect.Slice &&
		reflect.SliceOf(valueType) == field.Type &&
		FieldIs(field, FIELD_WRAP) {
		return true
	}
	// if type is slice and argument is steps, it's OK
	if field.Type.Kind() == reflect.Slice &&
		FieldIs(field, FIELD_STEPS) {
		return true
	}
	// check that value is of given type
	return IsValueOfType(value, field.Type)
}

// IsValueOfType tells if a value is of given type
// - value: the value to test as an interface{}
// - type: the type to check as a reflect.Type
// Return: a bool telling if value is of given type
func IsValueOfType(value interface{}, typ reflect.Type) bool {
	// if value is of given type it's true
	if reflect.TypeOf(value) == typ {
		return true
	}
	// if value and type are slices, test elements
	if reflect.TypeOf(value).Kind() == reflect.Slice && typ.Kind() == reflect.Slice {
		return IsValueOfType(reflect.ValueOf(value).Index(0).Interface(), typ.Elem())
	}
	// if value and type are maps, test key and value
	if reflect.TypeOf(value).Kind() == reflect.Map && typ.Kind() == reflect.Map {
		key := reflect.ValueOf(value).MapKeys()[0].Interface()
		val := reflect.ValueOf(value).MapIndex(reflect.ValueOf(key)).Interface()
		return IsValueOfType(key, typ.Key()) && IsValueOfType(val, typ.Elem())
	}
	// else return false
	return false
}

// Build task arguments from task params and return it
// - args: task arguments parsed in build file
// - typ: the type of the arguments
// - context: the build context to evaluate arguments into
// Return:
// - result: as an interface{}
// - error: if something went wrong
func EvaluateTaskArgs(args TaskArgs, typ reflect.Type, context *Context) (interface{}, error) {
	var err error
	value := reflect.New(typ).Elem()
	for i := 0; i < value.NumField(); i++ {
		name := strings.ToLower(typ.Field(i).Name)
		field := typ.Field(i)
		if field.Tag.Get(FIELD_NAME) != "" {
			name = field.Tag.Get(FIELD_NAME)
		}
		if args[name] != nil {
			val := args[name]
			// evaluate expressions in context
			if reflect.TypeOf(val).Kind() == reflect.String &&
				(IsExpression(val.(string)) || FieldIs(field, FIELD_EXPRESSION)) {
				str := val.(string)
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
			// evaluate arguments
			val, err = context.EvaluateObject(val)
			if err != nil {
				return nil, err
			}
			// evaluate strings to replace "={expression}" with its value
			if reflect.TypeOf(val).Kind() == reflect.String {
				str := val.(string)
				// replace '\=' with '='
				if strings.HasPrefix(str, `\`+CHAR_EXPRESSION) {
					str = str[1:]
				}
				// expand home if field tagged 'file'
				if FieldIs(field, FIELD_FILE) {
					str = util.ExpandUserHome(str)
				}
				val = str
			}
			// wrap values if necessary
			if FieldIs(field, FIELD_WRAP) && !(reflect.TypeOf(val).Kind() == reflect.Slice) {
				slice := reflect.New(field.Type).Elem()
				slice = reflect.Append(slice, reflect.ValueOf(val))
				val = slice.Interface()
			}
			// put value in params
			CopyValue(reflect.ValueOf(val), value.Field(i))
		}
	}
	return value.Interface(), nil
}

// CopyValue copy given avlue in another
// - orig: origin value
// - dest: destination value
func CopyValue(orig, dest reflect.Value) {
	// loop on slices
	if orig.Type().Kind() == reflect.Slice && dest.Type().Kind() == reflect.Slice {
		new := reflect.MakeSlice(dest.Type(), orig.Len(), orig.Len())
		for i := 0; i < orig.Len(); i++ {
			CopyValue(orig.Index(i), new.Index(i))
		}
		dest.Set(new)
	} else
	// loop on maps
	if orig.Type().Kind() == reflect.Map && dest.Type().Kind() == reflect.Map {
		keyType := dest.Type().Key()
		valType := dest.Type().Elem()
		new := reflect.MakeMap(reflect.MapOf(keyType, valType))
		for _, key := range orig.MapKeys() {
			if key.Type().Kind() == reflect.Interface {
				key = key.Elem()
			}
			val := orig.MapIndex(key)
			if val.Type().Kind() == reflect.Interface {
				val = val.Elem()
			}
			new.SetMapIndex(key, val)
		}
		dest.Set(new)
	} else
	// get value of interfaces
	if orig.Kind() == reflect.Interface {
		CopyValue(orig.Elem(), dest)
	} else
	// other types
	{
		dest.Set(orig)
	}
}

// FieldIs tells if given field tag contains quality
// - field: the struct field
// - quality: the tested quality (such as "optional")
func FieldIs(field reflect.StructField, quality string) bool {
	tag := string(field.Tag)
	qualities := strings.Split(tag, TAG_SEPARATOR)
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
	return s[0:1] == CHAR_EXPRESSION && s[1:2] != CURLY_EXPRESSION
}
