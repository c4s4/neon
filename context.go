package main

import (
	"github.com/mattn/anko/vm"
	"os"
	"os/exec"
	"regexp"
	"runtime"
)

type Context struct {
	Env   *vm.Env
	Build *Build
}

func NewContext(build *Build) *Context {
	context := Context{
		Env:   vm.NewEnv(),
		Build: build,
	}
	for name, value := range build.Properties {
		context.SetProperty(name, value)
	}
	return &context
}

func (context *Context) SetProperty(name string, value interface{}) {
	context.Env.Define(name, value)
}

func (context *Context) GetProperty(name string) interface{} {
	value, err := context.Env.Get(name)
	StopOnError(err, "Error getting value of '"+name+"'", 10)
	return value
}

func (context *Context) ReplaceProperty(expression string) string {
	name := expression[2 : len(expression)-1]
	value, err := context.Env.Get(name)
	StopOnError(err, "Error getting value of '"+name+"'", 10)
	return value.String()
}

func (context *Context) ReplaceProperties(command string) string {
	r := regexp.MustCompile("#{.*?}")
	replaced := r.ReplaceAllStringFunc(command, context.ReplaceProperty)
	return replaced
}

func (context *Context) Execute(cmd string) {
	cmd = context.ReplaceProperties(cmd)
	var command *exec.Cmd
	if runtime.GOOS == "windows" {
		command = exec.Command("cmd.exe", "/C", cmd)
	} else {
		command = exec.Command("sh", "-c", cmd)
	}
	command.Dir = context.Build.Dir
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()
	StopOnError(err, "Error running command '"+cmd+"'", 5)
}
