package build

import (
	"fmt"
	anko_core "github.com/mattn/anko/builtins"
	"github.com/mattn/anko/parser"
	"github.com/mattn/anko/vm"
	zglob "github.com/mattn/go-zglob"
	"io/ioutil"
	"neon/util"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

// Build context
type VM struct {
	VM          *vm.Env
	Build       *Build
	Properties  []string
	Environment map[string]string
}

// NewVM Make a new virtual machine
func NewVM(build *Build) (*VM, error) {
	vm := vm.NewEnv()
	anko_core.LoadAllBuiltins(vm)
	LoadBuiltins(vm)
	properties := build.GetProperties()
	environment := build.GetEnvironment()
	context := &VM{
		VM:          vm,
		Build:       build,
		Properties:  properties.Fields(),
		Environment: environment,
	}
	for _, script := range build.Scripts {
		source, err := ioutil.ReadFile(script)
		if err != nil {
			return nil, fmt.Errorf("reading script '%s': %v", script, err)
		}
		_, err = vm.Execute(string(source))
		if err != nil {
			return nil, fmt.Errorf("evaluating script '%s': %v", script, FormatScriptError(err))
		}
	}
	err := context.SetInitialProperties(properties)
	if err != nil {
		return nil, fmt.Errorf("evaluating properties: %v", err)
	}
	return context, nil
}

// Set initial build properties
func (context *VM) SetInitialProperties(object util.Object) error {
	context.SetProperty("_OS", runtime.GOOS)
	context.SetProperty("_ARCH", runtime.GOARCH)
	context.SetProperty("_CPUS", runtime.NumCPU())
	context.SetProperty("_BASE", context.Build.Dir)
	context.SetProperty("_HERE", context.Build.Here)
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
func (context *VM) SetProperty(name string, value interface{}) {
	context.VM.Define(name, value)
}

// Get property value with given name
func (context *VM) GetProperty(name string) (interface{}, error) {
	value, err := context.VM.Get(name)
	if err != nil {
		return nil, err
	}
	return util.ValueToInterface(value), nil
}

// Evaluate given expression in context and return its value
func (context *VM) EvaluateExpression(source string) (interface{}, error) {
	value, err := context.VM.Execute(source)
	if err != nil {
		return nil, FormatScriptError(err)
	}
	return util.ValueToInterface(value), nil
}

// Evaluate a given object, that is replace '#{foo}' in strings with the value
// of property foo
func (context *VM) EvaluateObject(object interface{}) (interface{}, error) {
	switch value := object.(type) {
	case string:
		evaluated, err := context.EvaluateString(value)
		if err != nil {
			return nil, err
		}
		return evaluated, nil
	case bool:
		return value, nil
	case int:
		return value, nil
	case int32:
		return value, nil
	case int64:
		return value, nil
	case float64:
		return value, nil
	default:
		if value == nil {
			return nil, nil
		}
		switch reflect.TypeOf(object).Kind() {
		case reflect.Slice:
			slice := reflect.ValueOf(object)
			elements := make([]interface{}, slice.Len())
			for index := 0; index < slice.Len(); index++ {
				val, err := context.EvaluateObject(slice.Index(index).Interface())
				if err != nil {
					return nil, err
				}
				elements[index] = val
			}
			return elements, nil
		case reflect.Map:
			dict := reflect.ValueOf(object)
			elements := make(map[interface{}]interface{})
			for _, key := range dict.MapKeys() {
				keyEval, err := context.EvaluateObject(key.Interface())
				if err != nil {
					return nil, err
				}
				valueEval, err := context.EvaluateObject(dict.MapIndex(key).Interface())
				if err != nil {
					return nil, err
				}
				elements[keyEval] = valueEval
			}
			return elements, nil
		default:
			return nil, fmt.Errorf("no serializer for type '%T'", object)
		}
	}
}

// Evaluate a string by replacing '#{foo}' with value of property foo
func (context *VM) EvaluateString(text string) (string, error) {
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
func (context *VM) EvaluateEnvironment() ([]string, error) {
	environment := make(map[string]string)
	for _, line := range os.Environ() {
		index := strings.Index(line, "=")
		name := line[:index]
		value := line[index+1:]
		environment[name] = value
	}
	environment["_BASE"] = context.Build.Dir
	environment["_HERE"] = context.Build.Here
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

// Find files in the context:
// - dir: the search root directory
// - includes: the list of globs to include
// - excludes: the list of globs to exclude
// - folder: tells if we should include folders
// Return the list of files as a slice of strings
func (context *VM) FindFiles(dir string, includes, excludes []string, folder bool) ([]string, error) {
	eval, err := context.EvaluateString(dir)
	if err != nil {
		return nil, fmt.Errorf("evaluating source directory: %v", err)
	}
	dir = util.ExpandUserHome(eval)
	if dir != "" {
		oldDir, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("getting working directory: %v", err)
		}
		defer os.Chdir(oldDir)
		err = os.Chdir(dir)
		if err != nil {
			return nil, nil
		}
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
		pattern = util.ExpandUserHome(pattern)
		excluded = append(excluded, pattern)
	}
	var candidates []string
	for _, include := range included {
		list, _ := zglob.Glob(util.ExpandUserHome(include))
		for _, file := range list {
			stat, err := os.Stat(file)
			if err != nil {
				return nil, fmt.Errorf("stating file: %v", err)
			}
			if stat.Mode().IsRegular() || folder {
				candidates = append(candidates, file)
			}
		}
	}
	var files []string
	if excluded != nil {
		for index, file := range candidates {
			for _, exclude := range excluded {
				match, err := zglob.Match(exclude, file)
				if match || err != nil {
					candidates[index] = ""
				}
			}
		}
		for _, file := range candidates {
			if file != "" {
				files = append(files, file)
			}
		}
	} else {
		files = candidates
	}
	sort.Strings(files)
	return files, nil
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
