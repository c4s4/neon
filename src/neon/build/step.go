package build

import (
	"fmt"
	"neon/util"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Step interface {
	Run() error
}

func NewStep(target *Target, step interface{}) (Step, error) {
	switch step := step.(type) {
	case string:
		return NewShellStep(target, step)
	case map[interface{}]interface{}:
		return NewTaskStep(target, step)
	default:
		return nil, fmt.Errorf("step must be string or map")
	}
}

type ShellStep struct {
	Target  *Target
	Command string
}

func NewShellStep(target *Target, shell string) (Step, error) {
	step := ShellStep{
		Target:  target,
		Command: shell,
	}
	return step, nil
}

func (step ShellStep) Run() error {
	cmd, err := step.Target.Build.Context.EvaluateString(step.Command)
	if err != nil {
		return fmt.Errorf("evaluating shell expression: %v", err)
	}
	var command *exec.Cmd
	if runtime.GOOS == "windows" {
		command = exec.Command("cmd.exe", "/C", cmd)
	} else {
		command = exec.Command("sh", "-c", cmd)
	}
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current working directory: %v", err)
	}
	command.Dir = dir
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env, err = step.Target.Build.Context.GetEnvironment()
	if err != nil {
		return fmt.Errorf("building environment: %v", err)
	}
	return command.Run()
}

type TaskStep struct {
	Target *Target
	Task   Task
}

func NewTaskStep(target *Target, m map[interface{}]interface{}) (Step, error) {
	object, err := util.NewObject(m)
	if err != nil {
		return nil, fmt.Errorf("a task must be a map with string keys")
	}
	fields := object.Fields()
	for name, descriptor := range TaskMap {
		for _, field := range fields {
			if name == field {
				task, err := descriptor.Constructor(target, object)
				if err != nil {
					return nil, fmt.Errorf("parsing task '%s': %v", name, err)
				}
				step := TaskStep{
					Target: target,
					Task:   task,
				}
				return step, nil
			}
		}
	}
	return nil, fmt.Errorf("unknown task '%s'", strings.Join(fields, "/"))
}

func (step TaskStep) Run() error {
	return step.Task()
}
