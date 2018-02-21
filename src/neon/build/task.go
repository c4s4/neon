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
	CharExpression = "="
	// expression curly brace
	CurlyExpression = "{"
	// tag name
	NeonTag = "neon"
	// tag separator
	TagSeparator = ","
	// field might not be provided
	FieldOptional = "optional"
	// field is a file name that is expanded for user home
	FieldFile = "file"
	// field is an expression
	FieldExpression = "expression"
	// field should be rapped in a slice of its type
	FieldWrap = "wrap"
	// field is a list of steps
	FieldSteps = "steps"
	// field has a different name
	FieldName = "name"
)

// TaskMap is a map that gives constructor for given task name
var TaskMap = make(map[string]TaskDesc)

// AddTask adds a task in the map:
// - task: description of the task
func AddTask(task TaskDesc) {
	if _, ok := TaskMap[task.Name]; ok {
		panic(fmt.Errorf("task '%s' already defined", task.Name))
	}
	err := CheckTaskArgs(TaskMap, task)
	if err != nil {
		panic(err)
	}
	TaskMap[task.Name] = task
}

// CheckTaskArgs checks that task argument names don't collide with with the
// name of an existing task. For instance, a task can't have an argument named
// copy as a task named copy already exists.
// - map: the tasks map.
// - task: the task description.
// Return: an error if argument collides.
func CheckTaskArgs(m map[string]TaskDesc, t TaskDesc) error {
	for n := range m {
		if HasField(t.Args, n) {
			return fmt.Errorf("task '%s' cannot have a field named '%s' as a task with this name exists", t.Name, n)
		}
	}
	return nil
}

// HasField tells if given structure has named field.
// - t: the type to check (must be a struct).
// - n: the name of the field to check as a string.
// Return: a boolean that tells if field exists.
func HasField(t reflect.Type, n string) bool {
	if t.Kind() != reflect.Struct {
		panic("task args must be a struct")
	}
	for i := 0; i < t.NumField(); i++ {
		if strings.ToLower(t.Field(i).Name) == n {
			return true
		}
	}
	return false
}

// TaskDesc is a task descriptor is made of a task constructor and an help string
type TaskDesc struct {
	Name string
	Args reflect.Type
	Func TaskFunc
	Help string
}

// TaskArgs is a type for task arguments as parsed in build file
type TaskArgs map[interface{}]interface{}

// TaskFunc is a type for task function called to run task
type TaskFunc func(ctx *Context, args interface{}) error

// ValidateTaskArgs validates arguments against task arguments definition
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
		name, err := checkArgumentType(field, args)
		if err != nil {
			return err
		}
		fields = append(fields, name)
	}
	// check that we don't have unknown args
	if err := checkUnknownArgs(fields, args); err != nil {
		return err
	}
	return nil
}

func checkArgumentType(field reflect.StructField, args TaskArgs) (string, error) {
	name := GetQuality(field, FieldName)
	if name == "" {
		name = strings.ToLower(field.Name)
	}
	// check field is not missing
	if _, ok := args[name]; !ok {
		if !FieldIs(field, FieldOptional) {
			return "", fmt.Errorf("missing mandatory field '%s'", name)
		}
	}
	value := args[name]
	// parse steps fields
	if FieldIs(field, "steps") && value != nil {
		steps, err := NewSteps(value)
		if err != nil {
			return "", fmt.Errorf("parsing field '%s': %v", name, err)
		}
		args[name] = steps
	}
	// check field type
	if !CheckType(field, value) {
		return "", fmt.Errorf("field '%s' must be of type '%s' ('%s' provided)",
			name, field.Type, reflect.TypeOf(value))
	}
	return name, nil
}

func checkUnknownArgs(fields []string, args TaskArgs) error {
	if !(len(fields) == 0 && len(args) == 1) {
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
	if FieldIs(field, FieldOptional) && value == nil {
		return true
	}
	// if argument is an expression it's OK whatever the type
	if valueType.Kind() == reflect.String &&
		(IsExpression(value.(string)) ||
			FieldIs(field, FieldExpression)) {
		return true
	}
	// if type of field is slice of the type of the argument and wrap, it's OK
	if field.Type.Kind() == reflect.Slice &&
		reflect.SliceOf(valueType) == field.Type &&
		FieldIs(field, FieldWrap) {
		return true
	}
	// if type is slice and argument is steps, it's OK
	if field.Type.Kind() == reflect.Slice &&
		FieldIs(field, FieldSteps) {
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

// EvaluateTaskArgs builds task arguments from task params and return it
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
		field := typ.Field(i)
		name := GetQuality(field, FieldName)
		if name == "" {
			name = strings.ToLower(field.Name)
		}
		if args[name] != nil {
			val := args[name]
			// evaluate expressions in context
			if reflect.TypeOf(val).Kind() == reflect.String &&
				(IsExpression(val.(string)) || FieldIs(field, FieldExpression)) {
				val, err = evaluateExpression(field, val, context)
				if err != nil {
					return nil, err
				}
			}
			// evaluate arguments
			val, err = context.EvaluateObject(val)
			if err != nil {
				return nil, err
			}
			// evaluate strings to replace "={expression}" with its value
			if reflect.TypeOf(val).Kind() == reflect.String {
				val = evaluateStrings(val, field)
			}
			// wrap values if necessary
			if FieldIs(field, FieldWrap) && !(reflect.TypeOf(val).Kind() == reflect.Slice) {
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

func evaluateStrings(val interface{}, field reflect.StructField) interface{} {
	str := val.(string)
	// replace '\=' with '='
	if strings.HasPrefix(str, `\`+CharExpression) {
		str = str[1:]
	}
	// expand home if field tagged 'file'
	if FieldIs(field, FieldFile) {
		str = util.ExpandUserHome(str)
	}
	val = str
	return val
}

func evaluateExpression(field reflect.StructField, val interface{}, context *Context) (interface{}, error) {
	var err error
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
		// we accept if expected is slice of interfaces and actual is slice
		if !(expected.Kind() == reflect.Slice && actual.Kind() == reflect.Slice) &&
			// or if slice and wrap
			!(expected == reflect.SliceOf(actual) && FieldIs(field, FieldWrap)) {
			return nil, fmt.Errorf("bad expression return type, expected '%v' but '%v' was returned",
				expected, actual)
		}
	}
	return val, nil
}

// CopyValue copy given value in another
// - orig: origin value
// - dest: destination value
func CopyValue(orig, dest reflect.Value) {
	// loop on slices
	if orig.Type().Kind() == reflect.Slice && dest.Type().Kind() == reflect.Slice {
		newSlice := reflect.MakeSlice(dest.Type(), orig.Len(), orig.Len())
		for i := 0; i < orig.Len(); i++ {
			CopyValue(orig.Index(i), newSlice.Index(i))
		}
		dest.Set(newSlice)
	} else
	// loop on maps
	if orig.Type().Kind() == reflect.Map && dest.Type().Kind() == reflect.Map {
		keyType := dest.Type().Key()
		valType := dest.Type().Elem()
		newMap := reflect.MakeMap(reflect.MapOf(keyType, valType))
		for _, key := range orig.MapKeys() {
			if key.Type().Kind() == reflect.Interface {
				key = key.Elem()
			}
			val := orig.MapIndex(key)
			if val.Type().Kind() == reflect.Interface {
				val = val.Elem()
			}
			newMap.SetMapIndex(key, val)
		}
		dest.Set(newMap)
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
	tags := field.Tag.Get(NeonTag)
	qualities := strings.Split(tags, TagSeparator)
	for _, q := range qualities {
		if q == quality {
			return true
		}
	}
	return false
}

// GetQuality returns value of given quality
// - field: the field to examine
// - quality: quality to get
func GetQuality(field reflect.StructField, quality string) string {
	tags := field.Tag.Get(NeonTag)
	qualities := strings.Split(tags, TagSeparator)
	for _, q := range qualities {
		prefix := quality + "="
		if strings.HasPrefix(q, prefix) {
			return q[len(prefix):]
		}
	}
	return ""
}

// IsExpression tells if given string is an expression
// - s: the string to test
// Return: a bool that tells if the string is an expression
func IsExpression(s string) bool {
	if len(s) < 2 {
		return false
	}
	return s[0:1] == CharExpression && s[1:2] != CurlyExpression
}
