package main

import (
	"github.com/mattn/anko/vm"
	"reflect"
	"regexp"
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

func (context *Context) ReplaceProperties(command string) (string, error) {
	r := regexp.MustCompile("#{.*?}")
	replaced := r.ReplaceAllStringFunc(command, context.replaceProperty)
	err := context.Error
	context.Error = nil
	return replaced, err
}
