package main

import (
	"github.com/mattn/anko/vm"
	"regexp"
	"sort"
)

type Context struct {
	Env        *vm.Env
	Build      *Build
	Properties []string
}

func NewContext(build *Build, object Object) *Context {
	context := &Context{
		Env:   vm.NewEnv(),
		Build: build,
	}
	var properties []string
	for name, value := range object {
		context.SetProperty(name, value)
		properties = append(properties, name)
	}
	sort.Strings(properties)
	context.Properties = properties
	return context
}

func (context *Context) SetProperty(name string, value interface{}) {
	context.Env.Define(name, value)
}

func (context *Context) GetProperty(name string) (interface{}, error) {
	value, err := context.Env.Get(name)
	if err != nil {
		return nil, err
	}
	return value.Interface(), nil
}

func (context *Context) ReplaceProperty(expression string) string {
	name := expression[2 : len(expression)-1]
	value, _ := context.Env.Get(name)
	return value.String()
}

func (context *Context) ReplaceProperties(command string) string {
	r := regexp.MustCompile("#{.*?}")
	replaced := r.ReplaceAllStringFunc(command, context.ReplaceProperty)
	return replaced
}
