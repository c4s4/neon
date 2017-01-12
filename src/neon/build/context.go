package build

import (
	"fmt"
	anko_core "github.com/mattn/anko/builtins"
	"github.com/mattn/anko/vm"
	"neon/builtin"
	"neon/util"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Context struct {
	Env        *vm.Env
	Build      *Build
	Properties []string
	Error      error
}

func NewContext(build *Build, object util.Object) (*Context, error) {
	env := vm.NewEnv()
	anko_core.LoadAllBuiltins(env)
	context := &Context{
		Env:        env,
		Build:      build,
		Properties: object.Fields(),
	}
	err := context.SetProperties(object)
	if err != nil {
		return nil, err
	}
	context.AddBuiltins()
	return context, nil
}

func (context *Context) AddBuiltins() {
	context.Env.Define("find", builtin.Find)
}

func (context *Context) Evaluate(source string) (interface{}, error) {
	value, err := context.Env.Execute(source)
	return value.Interface(), err
}

func (context *Context) SetProperty(name string, value interface{}) {
	context.Env.Define(name, value)
}

func (context *Context) SetProperties(object util.Object) error {
	todo := object.Fields()
	var crash error
	for len(todo) > 0 {
		var done []string
		for _, name := range todo {
			value := object[name]
			str, ok := value.(string)
			if ok {
				eval, err := context.ReplaceProperties(str)
				if err == nil {
					context.SetProperty(name, eval)
					done = append(done, name)
				} else {
					crash = err
				}
			} else {
				context.SetProperty(name, value)
				done = append(done, name)
			}
		}
		if len(done) == 0 {
			return fmt.Errorf("evaluating properties: %v", crash)
		}
		var next []string
		for _, name := range todo {
			found := false
			for _, n := range done {
				if name == n {
					found = true
				}
			}
			if !found {
				next = append(next, name)
			}
		}
		todo = next
	}
	return nil
}

func (context *Context) GetProperty(name string) (interface{}, error) {
	value, err := context.Env.Get(name)
	if err != nil {
		return nil, err
	}
	return value.Interface(), nil
}

func (context *Context) replaceProperty(expression string) string {
	name := expression[2 : len(expression)-1]
	value, err := context.GetProperty(name)
	if err != nil {
		context.Error = err
	}
	var str string
	if err == nil {
		str, err = PropertyToString(value, false)
		if err != nil {
			context.Error = err
		}
	}
	return str
}

func (context *Context) ReplaceProperties(text string) (string, error) {
	r := regexp.MustCompile("#{.*?}")
	replaced := r.ReplaceAllStringFunc(text, context.replaceProperty)
	err := context.Error
	context.Error = nil
	return replaced, err
}

func PropertyToString(object interface{}, quotes bool) (string, error) {
	switch value := object.(type) {
	case bool:
		return strconv.FormatBool(value), nil
	case string:
		if quotes {
			return "\"" + value + "\"", nil
		} else {
			return value, nil
		}
	case int:
		return strconv.Itoa(value), nil
	case int32:
		return strconv.Itoa(int(value)), nil
	case int64:
		return strconv.Itoa(int(value)), nil
	case float64:
		return strconv.FormatFloat(value, 'g', -1, 64), nil
	default:
		switch reflect.TypeOf(object).Kind() {
		case reflect.Slice:
			slice := reflect.ValueOf(object)
			elements := make([]string, slice.Len())
			for index := 0; index < slice.Len(); index++ {
				str, err := PropertyToString(slice.Index(index).Interface(), quotes)
				if err != nil {
					return "", err
				}
				elements[index] = str
			}
			return "[" + strings.Join(elements, ", ") + "]", nil
		case reflect.Map:
			dict := reflect.ValueOf(object)
			elements := make(map[string]string)
			var keys []string
			for _, key := range dict.MapKeys() {
				value := dict.MapIndex(key)
				keyStr, err := PropertyToString(key.Interface(), quotes)
				if err != nil {
					return "", err
				}
				keys = append(keys, keyStr)
				valueStr, err := PropertyToString(value.Interface(), quotes)
				if err != nil {
					return "", err
				}
				elements[keyStr] = valueStr
			}
			sort.Strings(keys)
			pairs := make([]string, len(keys))
			for index, key := range keys {
				pairs[index] = key + ": " + elements[key]
			}
			return "[" + strings.Join(pairs, ", ") + "]", nil
		default:
			return "", fmt.Errorf("no serializer for type '%T'", object)
		}
	}
}
