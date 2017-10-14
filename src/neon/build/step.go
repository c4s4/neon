package build

import (
	"fmt"
	"neon/util"
	"strings"
)

// A step has a Run() method
type Step interface {
	Run(context *Context) error
}

// Make a step inside given target
func NewStep(target *Target, step interface{}) (Step, error) {
	switch step := step.(type) {
	case string:
		return NewScriptStep(target, step)
	case map[interface{}]interface{}:
		return NewTaskStep(target, step)
	default:
		return nil, fmt.Errorf("step must be string or map")
	}
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
func (step ScriptStep) Run(context *Context) error {
	_, err := context.VM.EvaluateExpression(step.Script)
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
func (step TaskStep) Run(context *Context) error {
	return step.Task(context)
}
