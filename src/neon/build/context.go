package build

import (
	"fmt"
	"io/ioutil"
	"neon/util"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"
	anko_core "github.com/c4s4/anko/builtins"
	"github.com/c4s4/anko/parser"
	"github.com/c4s4/anko/vm"
)

// Build context
type Context struct {
	VM          *vm.Env
	Properties  []string
	Environment map[string]string
	Index       *Index
	Stack       *Stack
}

// NewContext make a new build context
func NewContext(build *Build) (*Context, error) {
	v := vm.NewEnv()
	anko_core.LoadAllBuiltins(v)
	LoadBuiltins(v)
	properties := build.GetProperties()
	environment := build.GetEnvironment()
	context := &Context{
		VM:          v,
		Properties:  properties.Fields(),
		Environment: environment,
		Index:       NewIndex(),
		Stack:       NewStack(),
	}
	for _, script := range build.Scripts {
		source, err := ioutil.ReadFile(script)
		if err != nil {
			return nil, fmt.Errorf("reading script '%s': %v", script, err)
		}
		_, err = v.Execute(string(source))
		if err != nil {
			return nil, fmt.Errorf("evaluating script '%s': %v", script, FormatScriptError(err))
		}
	}
	err := context.setInitialProperties(build, properties)
	if err != nil {
		return nil, fmt.Errorf("evaluating properties: %v", err)
	}
	return context, nil
}

// NewThreadContext builds a context in a thread
func (context *Context) NewThreadContext(thread int, input interface{}, output interface{}) *Context {
	copy := context.Copy()
	copy.SetProperty("_thread", thread)
	copy.SetProperty("_input", input)
	return copy
}

func (context *Context) Copy() *Context {
	properties := make([]string, len(context.Properties))
	for i := 0; i < len(context.Properties); i++ {
		properties[i] = context.Properties[i]
	}
	environment := make(map[string]string)
	for name, value := range context.Environment {
		environment[name] = value
	}
	copy := Context{
		VM:          context.VM.Copy(),
		Properties:  properties,
		Environment: environment,
		Index:       context.Index.Copy(),
		Stack:       context.Stack.Copy(),
	}
	return &copy
}

// Set initial build properties
func (context *Context) setInitialProperties(build *Build, object util.Object) error {
	context.SetProperty("_OS", runtime.GOOS)
	context.SetProperty("_ARCH", runtime.GOARCH)
	context.SetProperty("_NCPU", runtime.NumCPU())
	context.SetProperty("_BASE", build.Dir)
	context.SetProperty("_HERE", build.Here)
	todo := object.Fields()
	var crash error
	for len(todo) > 0 {
		var done []string
		for _, name := range todo {
			value := object[name]
			eval, err := context.EvaluateObject(value)
			if err == nil {
				context.SetProperty(name, eval)
				done = append(done, name)
			} else {
				crash = err
			}
		}
		if len(done) == 0 {
			return crash
		}
		var next []string
		for _, name := range todo {
			found := false
			for _, n := range done {
				if name == n {
					found = true
				}
			}
			if !found {
				next = append(next, name)
			}
		}
		todo = next
	}
	return nil
}

// Set property with given to given value
func (context *Context) SetProperty(name string, value interface{}) {
	context.VM.Define(name, value)
}

// Get property value with given name
func (context *Context) GetProperty(name string) (interface{}, error) {
	value, err := context.VM.Get(name)
	if err != nil {
		return nil, err
	}
	return util.ValueToInterface(value), nil
}

// Evaluate given expression in context and return its value
func (context *Context) EvaluateExpression(source string) (interface{}, error) {
	value, err := context.VM.Execute(source)
	if err != nil {
		return nil, FormatScriptError(err)
	}
	return util.ValueToInterface(value), nil
}

// Evaluate a given object, that is replace '#{foo}' in strings with the value
// of property foo
func (context *Context) EvaluateObject(object interface{}) (interface{}, error) {
	// we replace #{expression} in strings with the result of the expression
	if reflect.TypeOf(object).Kind() == reflect.String {
		evaluated, err := context.EvaluateString(object.(string))
		if err != nil {
			return nil, err
		}
		return evaluated, nil
	}
	// we go inside slices and maps to process strings
	if reflect.TypeOf(object).Kind() == reflect.Slice ||
		reflect.TypeOf(object).Kind() == reflect.Map {
		value := reflect.ValueOf(object)
		for i:=0; i<value.Len(); i++ {
			index := value.Index(i)
			evaluated, err := context.EvaluateObject(index)
			if err != nil {
				return nil, err
			}
			index.Set(reflect.ValueOf(evaluated))
		}
		return object, nil
	}
	// else we do nothing
	return object, nil
}

// Evaluate a string by replacing '#{foo}' with value of property foo
func (context *Context) EvaluateString(text string) (string, error) {
	r := regexp.MustCompile(`#{.*?}`)
	var errors []error
	replaced := r.ReplaceAllStringFunc(text, func(expression string) string {
		name := expression[2 : len(expression)-1]
		var value interface{}
		value, err := context.EvaluateExpression(name)
		if err != nil {
			errors = append(errors, err)
			return ""
		} else {
			var str string
			str, err = PropertyToString(value, false)
			if err != nil {
				errors = append(errors, err)
				return ""
			} else {
				return str
			}
		}
	})
	if len(errors) > 0 {
		return replaced, errors[0]
	} else {
		return replaced, nil
	}
}

// Evaluate environment in context and return it as a slice of strings
func (context *Context) EvaluateEnvironment(build *Build) ([]string, error) {
	environment := make(map[string]string)
	for _, line := range os.Environ() {
		index := strings.Index(line, "=")
		name := line[:index]
		value := line[index+1:]
		environment[name] = value
	}
	environment["_BASE"] = build.Dir
	environment["_HERE"] = build.Here
	var variables []string
	for name := range context.Environment {
		variables = append(variables, name)
	}
	sort.Strings(variables)
	for _, name := range variables {
		value := context.Environment[name]
		r := regexp.MustCompile(`[$#]{.*?}`)
		replaced := r.ReplaceAllStringFunc(value, func(expression string) string {
			name := expression[2 : len(expression)-1]
			if expression[0:1] == "$" {
				value, ok := environment[name]
				if !ok {
					return expression
				} else {
					return value
				}
			} else {
				value, err := context.EvaluateExpression(name)
				if err != nil {
					return expression
				} else {
					str, _ := PropertyToString(value, false)
					return str
				}
			}
		})
		environment[name] = replaced
	}
	var lines []string
	for name, value := range environment {
		line := name + "=" + value
		lines = append(lines, line)
	}
	return lines, nil
}

// Run steps in context
func (context *Context) Run(steps []Step) error {
	context.Index.Expand()
	for index, step := range steps {
		context.Index.Set(index)
		err := step.Run(context)
		if err != nil {
			return err
		}
	}
	context.Index.Shrink()
	return nil
}

// Find files in the context:
// - dir: the search root directory
// - includes: the list of globs to include
// - excludes: the list of globs to exclude
// - folder: tells if we should include folders
// Return the list of files as a slice of strings
func (context *Context) FindFiles(dir string, includes, excludes []string, folder bool) ([]string, error) {
	dir, err := context.EvaluateString(dir)
	if err != nil {
		return nil, fmt.Errorf("evaluating source directory: %v", err)
	}
	var included []string
	for _, include := range includes {
		pattern, err := context.EvaluateString(include)
		if err != nil {
			return nil, fmt.Errorf("evaluating pattern: %v", err)
		}
		included = append(included, pattern)
	}
	var excluded []string
	for _, exclude := range excludes {
		pattern, err := context.EvaluateString(exclude)
		if err != nil {
			return nil, fmt.Errorf("evaluating pattern: %v", err)
		}
		excluded = append(excluded, pattern)
	}
	return util.FindFiles(dir, included, excluded, folder)
}

// Message print a message on the console
func (context *Context) Message(text string, args ...interface{}) {
	Message(text, args...)
}

// FormatScriptError adds line and column numbers on parser or vm errors.
func FormatScriptError(err error) error {
	if e, ok := err.(*parser.Error); ok {
		return fmt.Errorf("%s (at line %d, column %d)", err, e.Pos.Line, e.Pos.Column)
	} else if e, ok := err.(*vm.Error); ok {
		return fmt.Errorf("%s (at line %d, column %d)", err, e.Pos.Line, e.Pos.Column)
	} else {
		return err
	}
}
