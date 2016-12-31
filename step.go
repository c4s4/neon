package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type Step struct {
	Target  *Target
	Command string
}

func NewStep(target *Target, cmd interface{}) (*Step, error) {
	command, ok := cmd.(string)
	if !ok {
		return nil, fmt.Errorf("step must be string")
	}
	step := &Step{
		Target:  target,
		Command: command,
	}
	return step, nil
}

func (step *Step) Run() error {
	cmd := step.Target.Build.Context.ReplaceProperties(step.Command)
	var command *exec.Cmd
	if runtime.GOOS == "windows" {
		command = exec.Command("cmd.exe", "/C", cmd)
	} else {
		command = exec.Command("sh", "-c", cmd)
	}
	command.Dir = step.Target.Build.Dir
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}
