package build

import (
	"fmt"
	"strings"
	"sort"
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
	_, err := context.EvaluateExpression(step.Script)
	if err != nil {
		return fmt.Errorf("evaluating script: %v", err)
	}
	return nil
}

// Structure for a task step
type TaskStep struct {
	Target *Target
	Desc   TaskDesc
	Args   TaskArgs
}

// Make a task step
func NewTaskStep(target *Target, args TaskArgs) (Step, error) {
	for name, desc := range TaskMap {
		for field := range args {
			if name == field {
				err := ValidateTaskArgs(args, desc.Args)
				if err != nil {
					return nil, fmt.Errorf("parsing task '%s': %v", name, err)
				}
				step := TaskStep{
					Target: target,
					Desc:   desc,
					Args:   args,
				}
				return step, nil
			}
		}
	}
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
		return err
	}
	return step.Desc.Func(context, params)
}
