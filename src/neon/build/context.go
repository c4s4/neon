package build

import (
	"fmt"
	"github.com/mattn/anko/core"
	"github.com/mattn/anko/packages"
	"github.com/mattn/anko/parser"
	"github.com/mattn/anko/vm"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

const (
	propertyOs     = "_OS"
	propertyArch   = "_ARCH"
	propertyNcpu   = "_NCPU"
	propertyBase   = "_BASE"
	propertyHere   = "_HERE"
	propertyThread = "_thread"
	propertyInput  = "_input"
	environmentSep = "="
	environmentVar = "$"
)

var (
	regexpExp   = regexp.MustCompile(`\\*[#=]{.*?([^\\])}`)
	regexpEnv   = regexp.MustCompile(`\\*[#=$]{.*?}`)
	regexpParts = regexp.MustCompile(`(\\*)([#=$]){(.*)}`)
	regexpClose = regexp.MustCompile(`\\*}`)
)

// Context is the context of the build
// - VM: Anko VM that holds build properties
// - Build: the current build
// - Index: tracks steps index while running build
// - Stack: tracks targets calls
type Context struct {
	VM    *vm.Env
	Build *Build
	Stack *Stack
}

// NewContext make a new build context
// Return: a pointer to the context
func NewContext(build *Build) *Context {
	v := vm.NewEnv()
	core.Import(v)
	packages.DefineImport(v)
	LoadBuiltins(v)
	context := &Context{
		VM:    v,
		Build: build,
		Stack: NewStack(),
	}
	return context
}

// NewThreadContext builds a context for a thread by copying the build context
// - thread: the number of the thread, starting with 0
// - input: the thread input
// - output: the thread output
// Return: a pointer to the context
func (context *Context) NewThreadContext(thread int, input interface{}, output interface{}) *Context {
	another := &Context{
		VM:    context.VM.Copy(),
		Build: context.Build,
		Stack: context.Stack.Copy(),
	}
	another.SetProperty(propertyThread, thread)
	another.SetProperty(propertyInput, input)
	return another
}

// Init initializes context with build
// Return: an error if something went wrong
func (context *Context) Init() error {
	err := context.InitScripts()
	if err != nil {
		return fmt.Errorf("loading scripts: %v", err)
	}
	err = context.InitProperties()
	if err != nil {
		return fmt.Errorf("evaluating properties: %v", err)
	}
	return nil
}

// InitScripts loads build scripts in context
// Return: an error if something went wrong
func (context *Context) InitScripts() error {
	for _, script := range context.Build.Scripts {
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
// Return: an error if something went wrong
func (context *Context) InitProperties() error {
	context.SetProperty(propertyOs, runtime.GOOS)
	context.SetProperty(propertyArch, runtime.GOARCH)
	context.SetProperty(propertyNcpu, runtime.NumCPU())
	context.SetProperty(propertyBase, context.Build.Dir)
	context.SetProperty(propertyHere, context.Build.Here)
	todo := context.Build.Properties.Fields()
	var crash error
	for len(todo) > 0 {
		var done []string
		for _, name := range todo {
			value := context.Build.Properties[name]
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
	return value, nil
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
	// TODO: check that this is not necessary anymore
	// if !value.IsValid() {
	// 	return nil, nil
	// }
	return value, nil
}

// EvaluateString replaces '#{expression}' with the value of the expression
// - text: the string to evaluate
// Return:
// - evaluated string
// - an error if something went wrong
func (context *Context) EvaluateString(text string) (string, error) {
	var errors []error
	replaced := regexpExp.ReplaceAllStringFunc(text, func(expression string) string {
		parts := regexpParts.FindStringSubmatch(expression)
		prefix := parts[1]
		char := parts[2]
		source := parts[3]
		// expression was escaped
		if len(prefix)%2 == 1 {
			return strings.Repeat(`\`, len(prefix)/2) + char + `{` + source + `}`
		}
		// expression not escaped
		if len(prefix)%2 == 0 {
			prefix = strings.Repeat(`\`, len(prefix)/2)
		}
		// replace escaped closing curly braces "\}"" in expression
		source = regexpClose.ReplaceAllStringFunc(source, func(text string) string {
			text = text[:len(text)-1]
			slash := strings.Repeat(`\`, len(text)/2)
			return slash + `}`
		})
		value, err := context.EvaluateExpression(source)
		if err != nil {
			errors = append(errors, err)
			return ""
		}
		var str string
		str, err = PropertyToString(value, false)
		if err != nil {
			errors = append(errors, err)
			return ""
		}
		return prefix + str
	})
	if len(errors) > 0 {
		return replaced, errors[0]
	}
	return replaced, nil
}

// EvaluateObject recursively evaluates strings in a structure
// - object: the object to evaluate
// Return:
// - evaluated copy of object
// - an error if something went wrong
func (context *Context) EvaluateObject(object interface{}) (interface{}, error) {
	// if nil, return nil
	if object == nil {
		return nil, nil
	} else
	// we replace expressions in strings
	if reflect.TypeOf(object).Kind() == reflect.String {
		str := object.(string)
		if IsExpression(str) {
			value, err := context.EvaluateExpression(strings.TrimSpace(str[1:]))
			if err != nil {
				return nil, err
			}
			return value, nil
		}
		evaluated, err := context.EvaluateString(object.(string))
		if err != nil {
			return nil, err
		}
		return evaluated, nil
	}
	// we copy slices in a new slice with evaluated values
	if reflect.TypeOf(object).Kind() == reflect.Slice {
		source := reflect.ValueOf(object)
		dest := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(object).Elem()), source.Len(), source.Len())
		for i := 0; i < source.Len(); i++ {
			index := source.Index(i)
			evaluated, err := context.EvaluateObject(index.Interface())
			if err != nil {
				return nil, err
			}
			dest.Index(i).Set(reflect.ValueOf(evaluated))
		}
		return dest.Interface(), nil
	}
	// we copy maps in a new map with evaluated keys and values
	if reflect.TypeOf(object).Kind() == reflect.Map {
		source := reflect.ValueOf(object)
		dest := reflect.MakeMap(reflect.MapOf(reflect.TypeOf(object).Key(), reflect.TypeOf(object).Elem()))
		keys := source.MapKeys()
		for i := 0; i < len(keys); i++ {
			key := keys[i]
			keyEval, err := context.EvaluateObject(key.Interface())
			if err != nil {
				return nil, err
			}
			val := source.MapIndex(key)
			valEval, err := context.EvaluateObject(val.Interface())
			if err != nil {
				return nil, err
			}
			dest.SetMapIndex(reflect.ValueOf(keyEval), reflect.ValueOf(valEval))
		}
		return dest.Interface(), nil
	}
	// else we do nothing
	return object, nil
}

// EvaluateEnvironment evaluates environment variables in the context
// Return:
// - evaluated environment as a slice of strings
// - an error if something went wrong
func (context *Context) EvaluateEnvironment() ([]string, error) {
	environment := make(map[string]string)
	for _, line := range os.Environ() {
		index := strings.Index(line, environmentSep)
		name := line[:index]
		value := line[index+1:]
		environment[name] = value
	}
	environment[propertyBase] = context.Build.Dir
	environment[propertyHere] = context.Build.Here
	var variables []string
	for name := range context.Build.Environment {
		variables = append(variables, name)
	}
	sort.Strings(variables)
	for _, name := range variables {
		value := context.Build.Environment[name]
		var replaced string
		if IsExpression(value) {
			val, err := context.EvaluateExpression(value[1:])
			if err != nil {
				return nil, err
			}
			if val == nil || reflect.TypeOf(val).Kind() != reflect.String {
				return nil, fmt.Errorf("expression in environment must return a string")
			}
			replaced = val.(string)
		} else {
			replaced = regexpEnv.ReplaceAllStringFunc(value, func(expression string) string {
				parts := regexpParts.FindStringSubmatch(expression)
				prefix := parts[1]
				char := parts[2]
				source := parts[3]
				// expression was escaped
				if prefix == `\` {
					return char + `{` + source + `}`
				}
				// expression not escaped
				if prefix == `\\` {
					prefix = `\`
				}
				if char == environmentVar {
					val, ok := environment[source]
					if !ok {
						return prefix + `{` + source + `}`
					}
					return prefix + val
				}
				val, err := context.EvaluateExpression(source)
				if err != nil {
					return prefix + `{` + source + `}`
				}
				str, _ := PropertyToString(val, false)
				return prefix + str
			})
		}
		environment[name] = replaced
	}
	var lines []string
	for name, value := range environment {
		// unset variables with empty string value
		if value != "" {
			line := name + environmentSep + value
			lines = append(lines, line)
		}
	}
	return lines, nil
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
