package build

import (
	"fmt"
	"strings"
	"sort"
	"reflect"
)

// A step has a Run() method
type Step interface {
	Run(context *Context) error
}

// Make a step inside given target
func NewStep(step interface{}) (Step, error) {
	switch step := step.(type) {
	case string:
		return NewScriptStep(step)
	case map[interface{}]interface{}:
		return NewTaskStep(step)
	default:
		return nil, fmt.Errorf("a step must be a string or a map (%v provided)", reflect.TypeOf(step))
	}
}

// A script step
type ScriptStep struct {
	Script string
}

// Make a script step
func NewScriptStep(script string) (Step, error) {
	step := ScriptStep{
		Script: script,
	}
	return step, nil
}

// Run a script step using Anko VM.
func (step ScriptStep) Run(context *Context) error {
	_, err := context.EvaluateExpression(step.Script)
	if err != nil {
		return fmt.Errorf("evaluating script: %v", err)
	}
	return nil
}

// Structure for a task step
type TaskStep struct {
	Desc   TaskDesc
	Args   TaskArgs
}

// Make a task step
func NewTaskStep(args TaskArgs) (Step, error) {
	// find the task in the map
	for name, desc := range TaskMap {
		for field := range args {
			if name == field {
				err := ValidateTaskArgs(args, desc.Args)
				if err != nil {
					return nil, fmt.Errorf("parsing task '%s': %v", name, err)
				}
				step := TaskStep{
					Desc:   desc,
					Args:   args,
				}
				return step, nil
			}
		}
	}
	// task was not found, build error message
	var fields []string
	for field := range args {
		fields = append(fields, field.(string))
	}
	fields = sort.StringSlice(fields)
	return nil, fmt.Errorf("unknown task '%s'", strings.Join(fields, "/"))
}

// Run a task step, calling the function for the step
func (step TaskStep) Run(context *Context) error {
	params, err := EvaluateTaskArgs(step.Args, step.Desc.Args, context)
	if err != nil {
		return fmt.Errorf("in task '%s': %v", step.Desc.Name, err)
	}
	return step.Desc.Func(context, params)
}
