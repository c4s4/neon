package main

import (
	"fmt"
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
	return &context
}

func (context *Context) ReplaceProperty(expression string) string {
	property := expression[2 : len(expression)-1]
	if value, ok := context.Build.Properties[property]; ok {
		return value
	} else {
		StopWithError(fmt.Sprintf("Property '%s' was not found", property), 7)
		return ""
	}
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
