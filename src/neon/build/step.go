package build

import (
	"fmt"
	"neon/util"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// A step has a Run() method
type Step interface {
	Run() error
}

// Make a step inside given target
func NewStep(target *Target, step interface{}) (Step, error) {
	switch step := step.(type) {
	case string:
		if strings.HasPrefix(step, "$") {
			return NewShellStep(target, step[1:])
		} else {
			return NewScriptStep(target, step)
		}
	case map[interface{}]interface{}:
		return NewTaskStep(target, step)
	default:
		return nil, fmt.Errorf("step must be string or map")
	}
}

// A shell step
type ShellStep struct {
	Target  *Target
	Command string
}

// Make a shell step
func NewShellStep(target *Target, shell string) (Step, error) {
	step := ShellStep{
		Target:  target,
		Command: shell,
	}
	return step, nil
}

// Run a shell step:
// - If running on windows, run shell with "cmd.exe"
// - Otherwise, run shell with "sh"
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
	command.Env, err = step.Target.Build.Context.EvaluateEnvironment()
	if err != nil {
		return fmt.Errorf("building environment: %v", err)
	}
	return command.Run()
}

// A script step
type ScriptStep struct {
	Target *Target
	Script string
}

// Make a script step
func NewScriptStep(target *Target, script string) (Step, error) {
	step := ScriptStep{
		Target: target,
		Script: script,
	}
	return step, nil
}

// Run a script step using Anko VM.
func (step ScriptStep) Run() error {
	_, err := step.Target.Build.Context.EvaluateExpression(step.Script)
	if err != nil {
		return fmt.Errorf("evaluating script: %v", err)
	}
	return nil
}

// Structure for a task step
type TaskStep struct {
	Target *Target
	Task   Task
}

// Make a task step
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

// Run a task step, calling the function for the step
func (step TaskStep) Run() error {
	return step.Task()
}
