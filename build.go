package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
)

type Build struct {
	File       string
	Dir        string
	Name       string
	Default    string
	Properties map[string]string
	Targets    map[string]*Target
}

func (build *Build) Init(file string) {
	build.File = file
	build.Dir = filepath.Dir(build.File)
	for name, target := range build.Targets {
		target.Init(build, name)
	}
}

func (build *Build) Run(targets []string) {
	for _, target := range targets {
		target := build.Target(target)
		target.Run()
	}
}

func (build *Build) Target(name string) *Target {
	if target, ok := build.Targets[name]; ok {
		return target
	} else {
		StopWithError(fmt.Sprintf("Target '%s' was not found", name), 6)
		return nil
	}
}

func (build *Build) ReplaceProperty(expression string) string {
	property := expression[2 : len(expression)-1]
	if value, ok := build.Properties[property]; ok {
		return value
	} else {
		StopWithError(fmt.Sprintf("Property '%s' was not found", property), 7)
		return ""
	}
}

func (build *Build) ReplaceProperties(command string) string {
	r := regexp.MustCompile("#{.*?}")
	replaced := r.ReplaceAllStringFunc(command, build.ReplaceProperty)
	return replaced
}

func (build *Build) Execute(cmd string) {
	cmd = build.ReplaceProperties(cmd)
	var command *exec.Cmd
	if runtime.GOOS == "windows" {
		command = exec.Command("cmd.exe", "/C", cmd)
	} else {
		command = exec.Command("sh", "-c", cmd)
	}
	command.Dir = build.Dir
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()
	StopOnError(err, "Error running command '"+cmd+"'", 5)
}
