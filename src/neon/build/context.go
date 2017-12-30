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

const (
	PROPERTY_OS     = "_OS"
	PROPERTY_ARCH   = "_ARCH"
	PROPERTY_NCPU   = "_NCPU"
	PROPERTY_BASE   = "_BASE"
	PROPERTY_HERE   = "_HERE"
	PROPERTY_THREAD = "_thread"
	PROPERTY_INPUT  = "_input"
	ENVIRONMENT_SEP = "="
)

var (
	REGEXP_EXP = regexp.MustCompile(`[#=]{.*?}`)
	REGEXP_ENV = regexp.MustCompile(`[$#=]{.*?}`)
)

// Context is the context of the build
// - VM: Anko VM that holds build properties
// - Index: tracks steps index while running build
// - Stack: tracks targets calls
type Context struct {
	VM    *vm.Env
	Index *Index
	Stack *Stack
}

// NewContext make a new build context
// Return: a pointer to the context
func NewContext() *Context {
	v := vm.NewEnv()
	anko_core.LoadAllBuiltins(v)
	LoadBuiltins(v)
	context := &Context{
		VM:          v,
		Index:       NewIndex(),
		Stack:       NewStack(),
	}
	return context
}

// NewThreadContext builds a context for a thread by copying the build context
// - thread: the number of the thread, starting with 0
// - input: the thread input
// - ouput: the thread output
// Return: a pointer to the context
func (context *Context) NewThreadContext(thread int, input interface{}, output interface{}) *Context {
	copy := &Context{
		VM:          context.VM.Copy(),
		Index:       context.Index.Copy(),
		Stack:       context.Stack.Copy(),
	}
	copy.SetProperty(PROPERTY_THREAD, thread)
	copy.SetProperty(PROPERTY_INPUT, input)
	return copy
}

// Init initializes context with build
// - build: the build
// Return: an error if something went wrong
func (context *Context) Init(build *Build) error {
	err := context.InitScripts(build)
	if err != nil {
		return fmt.Errorf("loading scripts: %v", err)
	}
	err = context.InitProperties(build)
	if err != nil {
		return fmt.Errorf("evaluating properties: %v", err)
	}
	return nil
}

// InitScript loads build scripts in context
// - build: the build
// Return: an error if something went wrong
func (context *Context) InitScripts(build *Build) error {
	for _, script := range build.Scripts {
		source, err := ioutil.ReadFile(script)
		if err != nil {
			return fmt.Errorf("reading script '%s': %v", script, err)
		}
		_, err = context.VM.Execute(string(source))
		if err != nil {
			return fmt.Errorf("evaluating script '%s': %v", script, FormatScriptError(err))
		}
	}
	return nil
}

// InitProperties sets build properties
// - build: the build
// Return: an error if something went wrong
func (context *Context) InitProperties(build *Build) error {
	context.SetProperty(PROPERTY_OS, runtime.GOOS)
	context.SetProperty(PROPERTY_ARCH, runtime.GOARCH)
	context.SetProperty(PROPERTY_NCPU, runtime.NumCPU())
	context.SetProperty(PROPERTY_BASE, build.Dir)
	context.SetProperty(PROPERTY_HERE, build.Here)
	todo := build.Properties.Fields()
	var crash error
	for len(todo) > 0 {
		var done []string
		for _, name := range todo {
			value := build.Properties[name]
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

// SetProperty sets given property in context
// - name: the name of the property
// - value: the value of the property
func (context *Context) SetProperty(name string, value interface{}) {
	context.VM.Define(name, value)
}

// GetProperty returns value of given property
// - name: the name of the property
// Return:
// - the value of the property
// - an error if something went wrong
func (context *Context) GetProperty(name string) (interface{}, error) {
	value, err := context.VM.Get(name)
	if err != nil {
		return nil, err
	}
	return util.ValueToInterface(value), nil
}

// EvaluateExpression evaluate given expression in the context
// - expression: the expression to evaluate
// Return:
// - the return value of the expression
// - an error if something went wrong
func (context *Context) EvaluateExpression(expression string) (interface{}, error) {
	value, err := context.VM.Execute(expression)
	if err != nil {
		return nil, FormatScriptError(err)
	}
	return util.ValueToInterface(value), nil
}

// EvaluateString replaces '#{expression}' with the value of the expression
// - text: the string to evaluate
// Return:
// - evaluated string
// - an error if something went wrong
func (context *Context) EvaluateString(text string) (string, error) {
	var errors []error
	replaced := REGEXP_EXP.ReplaceAllStringFunc(text, func(expression string) string {
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

// EvaluateRecursive recursively evaluates strings in a structure
// - object: the object to evaluate
// Return:
// - evaluated structure
// - an error if something went wrong
func (context *Context) EvaluateObject(object interface{}) (interface{}, error) {
	// we replace expressions in strings
	if reflect.TypeOf(object).Kind() == reflect.String {
		evaluated, err := context.EvaluateString(object.(string))
		if err != nil {
			return nil, err
		}
		return evaluated, nil
	}
	// we iterate through slices
	if reflect.TypeOf(object).Kind() == reflect.Slice {
		value := reflect.ValueOf(object)
		for i:=0; i<value.Len(); i++ {
			index := value.Index(i)
			evaluated, err := context.EvaluateObject(index.Interface())
			if err != nil {
				return nil, err
			}
			index.Set(reflect.ValueOf(evaluated))
		}
		return object, nil
	}
	// we iterate through maps
	if reflect.TypeOf(object).Kind() == reflect.Map {
		value := reflect.ValueOf(object)
		keys := value.MapKeys()
		for i:=0; i<len(keys); i++ {
			key := keys[i]
			keyEval, err := context.EvaluateObject(key.Interface())
			if err != nil {
				return nil, err
			}
			val := value.MapIndex(key)
			valEval, err := context.EvaluateObject(val.Interface())
			if err != nil {
				return nil, err
			}
			value.SetMapIndex(key, reflect.Value{})
			value.SetMapIndex(reflect.ValueOf(keyEval), reflect.ValueOf(valEval))
		}
		return object, nil
	}
	// else we do nothing
	return object, nil
}

// EvaluateEnvironment evaluates environment variables in the context
// - build: the build with environment to evaluate
// Return:
// - evaluated environment as a slice of strings
// - an error if something went wrong
func (context *Context) EvaluateEnvironment(build *Build) ([]string, error) {
	environment := make(map[string]string)
	for _, line := range os.Environ() {
		index := strings.Index(line, ENVIRONMENT_SEP)
		name := line[:index]
		value := line[index+1:]
		environment[name] = value
	}
	environment[PROPERTY_BASE] = build.Dir
	environment[PROPERTY_HERE] = build.Here
	var variables []string
	for name := range build.Environment {
		variables = append(variables, name)
	}
	sort.Strings(variables)
	for _, name := range variables {
		value := build.Environment[name]
		replaced := REGEXP_ENV.ReplaceAllStringFunc(value, func(expression string) string {
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

// RunSteps runs a list of steps in context
// Return: an error if something went wrong
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

// Message print a message on the console
// - text: the text to print on console
// - args: a slice of string arguments (as for fmt.Printf())
func (context *Context) Message(text string, args ...interface{}) {
	Message(text, args...)
}

// FormatScriptError adds line and column numbers on parser or vm errors.
// - err: the error to process
// Return: the processed error
func FormatScriptError(err error) error {
	if e, ok := err.(*parser.Error); ok {
		return fmt.Errorf("%s (at line %d, column %d)", err, e.Pos.Line, e.Pos.Column)
	} else if e, ok := err.(*vm.Error); ok {
		return fmt.Errorf("%s (at line %d, column %d)", err, e.Pos.Line, e.Pos.Column)
	} else {
		return err
	}
}
