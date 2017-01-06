package main

import (
	"fmt"
	"github.com/mattn/anko/vm"
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

func NewContext(build *Build, object Object) (*Context, error) {
	context := &Context{
		Env:        vm.NewEnv(),
		Build:      build,
		Properties: object.Fields(),
	}
	err := context.SetProperties(object)
	if err != nil {
		return nil, err
	}
	return context, nil
}

func (context *Context) SetProperty(name string, value interface{}) {
	context.Env.Define(name, value)
}

func (context *Context) SetProperties(object Object) error {
	todo, _ := NewObject(object)
	length := len(todo)
	list := make([]string, len(todo)+1)
	var err error
	for length < len(list) && len(todo) > 0 {
		list = todo.Fields()
		for _, field := range list {
			value := todo[field]
			str, ok := value.(string)
			if ok {
				replaced, err := context.ReplaceProperties(str)
				if err != nil {
					continue
				} else {
					context.SetProperty(field, replaced)
					delete(todo, field)
				}
			} else {
				context.SetProperty(field, value)
				delete(todo, field)
			}
		}
		length = len(todo)
	}
	if len(todo) > 0 {
		return err
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
	context.Error = err
	return reflect.ValueOf(value).String()
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
