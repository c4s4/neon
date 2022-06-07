package build

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Step has a Run() method
type Step interface {
	Run(context *Context) error
}

// NewStep makes a new step
// - step: body of the step as an interface
// Return:
// - built step
// - error if something went wrong
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

// ScriptStep is made of a string
type ScriptStep struct {
	Script string
}

// NewScriptStep makes a new script step
// - script: the script as a string
// - built step
// - error if something went wrong
func NewScriptStep(script string) (Step, error) {
	step := ScriptStep{
		Script: script,
	}
	return step, nil
}

// Run a script step using Anko VM.
// - context: the build context to run tha script
// Return: an error if something went wrong
func (step ScriptStep) Run(context *Context) error {
	_, err := context.EvaluateExpression(step.Script)
	if err != nil {
		return fmt.Errorf("evaluating script: %v", err)
	}
	return nil
}

// TaskStep for a task step
type TaskStep struct {
	Desc TaskDesc
	Args TaskArgs
}

// NewTaskStep makes a task step
// - args: task args
// Return:
// - built step
// - error if something went wrong
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
	sort.Strings(fields)
	return nil, fmt.Errorf("unknown task '%s'", strings.Join(fields, "/"))
}

// Run a task step, calling the function for the step
// - context: build context
// Return: an error if something went wrong
func (step TaskStep) Run(context *Context) error {
	params, err := EvaluateTaskArgs(step.Args, step.Desc.Args, context)
	if err != nil {
		return fmt.Errorf("in task '%s': %v", step.Desc.Name, err)
	}
	return step.Desc.Func(context, params)
}

// Steps is a list of steps
type Steps []Step

// NewSteps makes a new steps
// - object: body of the steps as an interface
// Return:
// - steps
// - an error if something went wrong
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
