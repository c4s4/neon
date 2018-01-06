package build

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
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
	Desc TaskDesc
	Args TaskArgs
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
					Desc: desc,
					Args: args,
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

// Steps is a list of steps
type Steps []Step

func NewSteps(object interface{}) (Steps, error) {
	if reflect.ValueOf(object).IsNil() {
		return []Step{}, nil
	}
	if reflect.TypeOf(object).Kind() != reflect.Slice {
		return nil, fmt.Errorf("steps must be a slice")
	}
	len := reflect.ValueOf(object).Len()
	steps := make([]Step, len)
	for i := 0; i < len; i++ {
		step, err := NewStep(reflect.ValueOf(object).Index(i).Interface())
		if err != nil {
			return nil, err
		}
		steps[i] = step
	}
	return steps, nil
}

// Run steps in context
// - context: the context for running
// Return: an error if something went wrong
func (steps Steps) Run(context *Context) error {
	for index, step := range steps {
		err := step.Run(context)
		if err != nil {
			return fmt.Errorf("in step %d: %v", index+1, err)
		}
	}
	return nil
}
