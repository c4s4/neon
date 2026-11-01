package build

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/mattn/anko/core"
	"github.com/mattn/anko/env"
	_ "github.com/mattn/anko/packages"
	"github.com/mattn/anko/vm"

	"github.com/c4s4/neon/neon/util"
)

const (
	propertyOs     = "_OS"
	propertyArch   = "_ARCH"
	propertyNcpu   = "_NCPU"
	propertyBase   = "_BASE"
	propertyHere   = "_HERE"
	propertyRepo   = "_REPO"
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
// - Stack: tracks targets calls
// - History: tracks targets that already ran
type Context struct {
	VM      *env.Env
	Build   *Build
	Stack   *Stack
	History *History
	// Calls lists build files called with neon task leading to this context
	Calls []string
}

// NewContext makes a new build context
// Return: a pointer to the context
func NewContext(build *Build) *Context {
	e := env.NewEnv()
	core.Import(e)
	core.ImportToX(e)
	LoadBuiltins(e)
	context := &Context{
		VM:      e,
		Build:   build,
		Stack:   NewStack(),
		History: NewHistory(),
	}
	return context
}

// Copy performs a deep copy of the Context
// Return: a pointer to the context copy
func (context *Context) Copy() *Context {
	another := &Context{
		VM:      context.VM.DeepCopy(),
		Build:   context.Build,
		Stack:   context.Stack.Copy(),
		History: context.History.Copy(),
		Calls:   append([]string{}, context.Calls...),
	}
	return another
}

// Init initializes context with build
// Return: an error if something went wrong
func (context *Context) Init() error {
	return context.InitWithProperties(nil)
}

// InitWithProperties initializes context with build, overriding given
// properties. Build properties are evaluated with overridden values.
// - overrides: properties that override those of the build
// Return: an error if something went wrong
func (context *Context) InitWithProperties(overrides map[string]interface{}) error {
	if err := context.InitScripts(); err != nil {
		return fmt.Errorf("loading scripts: %v", err)
	}
	if err := context.InitProperties(overrides); err != nil {
		return fmt.Errorf("evaluating properties: %v", err)
	}
	if err := context.InitEnvironment(); err != nil {
		return fmt.Errorf("evaluating environment: %v", err)
	}
	return nil
}

// InitScripts loads build scripts in context
// Return: an error if something went wrong
func (context *Context) InitScripts() error {
	scripts, err := context.Build.GetScripts()
	if err != nil {
		return err
	}
	for _, path := range scripts {
		source, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading script '%s': %v", path, err)
		}
		_, err = vm.Execute(context.VM, nil, string(source))
		if err != nil {
			return fmt.Errorf("evaluating script '%s': %v", path, FormatScriptError(err))
		}
	}
	return nil
}

// InitProperties sets build properties
// - overrides: properties that override those of the build (may be nil)
// Return: an error if something went wrong
func (context *Context) InitProperties(overrides map[string]interface{}) error {
	context.SetProperty(propertyOs, runtime.GOOS)
	context.SetProperty(propertyArch, runtime.GOARCH)
	context.SetProperty(propertyNcpu, runtime.NumCPU())
	context.SetProperty(propertyBase, context.Build.Dir)
	context.SetProperty(propertyHere, context.Build.Here)
	context.SetProperty(propertyRepo, context.Build.Repository)
	for name, value := range overrides {
		context.SetProperty(name, value)
	}
	var todo []string
	for _, name := range context.Build.Properties.Fields() {
		if _, ok := overrides[name]; !ok {
			todo = append(todo, name)
		}
	}
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

// InitEnvironment sets environment variables
// Return: an error if something went wrong
func (context *Context) InitEnvironment() error {
	environment, err := context.EvaluateEnvironment()
	if err != nil {
		return err
	}
	for name, value := range environment {
		if err := os.Setenv(name, value); err != nil {
			return fmt.Errorf("setting environment variable '%s': %w", name, err)
		}
	}
	return nil
}

// SetProperty sets given property in context
// - name: the name of the property
// - value: the value of the property
func (context *Context) SetProperty(name string, value interface{}) {
	_ = context.VM.Define(name, value)
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

// DelProperty deletes given property
// - name: the name of the property
func (context *Context) DelProperty(name string) {
	context.VM.Delete(name)
}

// EvaluateExpression evaluates given expression in the context
// - expression: the expression to evaluate
// Return:
// - the return value of the expression
// - an error if something went wrong
func (context *Context) EvaluateExpression(expression string) (interface{}, error) {
	value, err := vm.Execute(context.VM, nil, expression)
	if err != nil {
		return nil, FormatScriptError(err)
	}
	return value, nil
}

// EvaluateString replaces '={expression}' (or '#{expression}') with the value
// of the expression
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
		// replace escaped closing curly braces "\}" in expression
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
	switch value := object.(type) {
	case nil:
		return nil, nil
	case string:
		return context.evaluateStringObject(value)
	case []interface{}:
		// lists decoded from YAML
		list := make([]interface{}, len(value))
		for i, element := range value {
			evaluated, err := context.EvaluateObject(element)
			if err != nil {
				return nil, err
			}
			list[i] = evaluated
		}
		return list, nil
	case map[interface{}]interface{}:
		// maps decoded from YAML
		result := make(map[interface{}]interface{}, len(value))
		for key, element := range value {
			keyEval, err := context.EvaluateObject(key)
			if err != nil {
				return nil, err
			}
			if keyEval != nil && !reflect.TypeOf(keyEval).Comparable() {
				return nil, fmt.Errorf("map key '%v' must be a scalar", keyEval)
			}
			valEval, err := context.EvaluateObject(element)
			if err != nil {
				return nil, err
			}
			result[keyEval] = valEval
		}
		return result, nil
	}
	// other slices and maps (such as values returned by expressions) are
	// copied with reflection to keep their type
	switch reflect.TypeOf(object).Kind() {
	case reflect.Slice:
		return context.evaluateSlice(object)
	case reflect.Map:
		return context.evaluateMap(object)
	}
	// other values (numbers, booleans...) are returned as is
	return object, nil
}

// evaluateStringObject evaluates a string: an expression if it starts with
// '=', or a string in which "={expression}" are replaced
func (context *Context) evaluateStringObject(str string) (interface{}, error) {
	if IsExpression(str) {
		return context.EvaluateExpression(strings.TrimSpace(str[1:]))
	}
	return context.EvaluateString(str)
}

// evaluateSlice evaluates elements of a typed slice
func (context *Context) evaluateSlice(object interface{}) (interface{}, error) {
	source := reflect.ValueOf(object)
	dest := reflect.MakeSlice(source.Type(), source.Len(), source.Len())
	for i := 0; i < source.Len(); i++ {
		evaluated, err := context.EvaluateObject(source.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		if err := setValue(dest.Index(i), evaluated); err != nil {
			return nil, err
		}
	}
	return dest.Interface(), nil
}

// evaluateMap evaluates keys and values of a typed map
func (context *Context) evaluateMap(object interface{}) (interface{}, error) {
	source := reflect.ValueOf(object)
	dest := reflect.MakeMapWithSize(source.Type(), source.Len())
	iter := source.MapRange()
	for iter.Next() {
		keyEval, err := context.EvaluateObject(iter.Key().Interface())
		if err != nil {
			return nil, err
		}
		valEval, err := context.EvaluateObject(iter.Value().Interface())
		if err != nil {
			return nil, err
		}
		keyValue, err := toValue(keyEval, dest.Type().Key())
		if err != nil {
			return nil, err
		}
		valValue, err := toValue(valEval, dest.Type().Elem())
		if err != nil {
			return nil, err
		}
		dest.SetMapIndex(keyValue, valValue)
	}
	return dest.Interface(), nil
}

// toValue converts an evaluated object to a reflect value of given type. A nil
// object gives the zero value of the type (nil for interfaces).
func toValue(object interface{}, typ reflect.Type) (reflect.Value, error) {
	if object == nil {
		return reflect.Zero(typ), nil
	}
	value := reflect.ValueOf(object)
	if !value.Type().AssignableTo(typ) {
		return reflect.Value{}, fmt.Errorf("value '%v' of type '%v' can't be assigned to type '%v'", object, value.Type(), typ)
	}
	return value, nil
}

// setValue sets an evaluated object in given destination value
func setValue(dest reflect.Value, object interface{}) error {
	value, err := toValue(object, dest.Type())
	if err != nil {
		return err
	}
	dest.Set(value)
	return nil
}

// EvaluateEnvironment evaluates environment variables in the context
// Return:
// - evaluated environment as a map of strings
// - an error if something went wrong
func (context *Context) EvaluateEnvironment() (map[string]string, error) {
	environment := make(map[string]string)
	for _, line := range os.Environ() {
		index := strings.Index(line, environmentSep)
		if index < 0 {
			continue
		}
		environment[line[:index]] = line[index+1:]
	}
	for _, filename := range context.Build.DotEnv {
		env, err := LoadDotEnv(filename, context.Build.Dir)
		if err != nil {
			return nil, err
		}
		for name, value := range env {
			environment[name] = value
		}
	}
	var variables []string
	for name := range context.Build.Environment {
		variables = append(variables, name)
	}
	sort.Strings(variables)
	// variables are evaluated after variables of the build they reference
	resolved := make(map[string]bool)
	var resolve func(name string, chain []string) error
	resolve = func(name string, chain []string) error {
		if resolved[name] {
			return nil
		}
		for _, n := range chain {
			if n == name {
				return fmt.Errorf("cycle in environment variables: %s -> %s", strings.Join(chain, " -> "), name)
			}
		}
		value := context.Build.Environment[name]
		for _, reference := range environmentReferences(value) {
			if _, ok := context.Build.Environment[reference]; ok && reference != name {
				if err := resolve(reference, append(chain, name)); err != nil {
					return err
				}
			}
		}
		replaced, err := context.evaluateEnvironmentValue(value, environment)
		if err != nil {
			return fmt.Errorf("evaluating environment variable '%s': %v", name, err)
		}
		// variables set to empty strings are unset (as documented)
		if replaced == "" {
			delete(environment, name)
		} else {
			environment[name] = replaced
		}
		resolved[name] = true
		return nil
	}
	for _, name := range variables {
		if err := resolve(name, nil); err != nil {
			return nil, err
		}
	}
	environment[propertyBase] = context.Build.Dir
	environment[propertyHere] = context.Build.Here
	return environment, nil
}

// environmentReferences returns names of variables referenced with ${NAME}
// (and not escaped) in given value
func environmentReferences(value string) []string {
	var references []string
	for _, expression := range regexpEnv.FindAllString(value, -1) {
		parts := regexpParts.FindStringSubmatch(expression)
		if len(parts[1])%2 == 0 && parts[2] == environmentVar {
			references = append(references, parts[3])
		}
	}
	return references
}

// evaluateEnvironmentValue evaluates the value of an environment variable,
// replacing ${NAME} with the value of the environment variable and ={expr} or
// #{expr} with the value of the expression. A value starting with '=' is an
// expression that must return a string.
func (context *Context) evaluateEnvironmentValue(value string, environment map[string]string) (string, error) {
	if IsExpression(value) {
		val, err := context.EvaluateExpression(value[1:])
		if err != nil {
			return "", err
		}
		str, ok := val.(string)
		if !ok {
			return "", fmt.Errorf("expression in environment must return a string")
		}
		return str, nil
	}
	var errs []error
	replaced := regexpEnv.ReplaceAllStringFunc(value, func(expression string) string {
		parts := regexpParts.FindStringSubmatch(expression)
		backslashes := parts[1]
		char := parts[2]
		source := parts[3]
		prefix := strings.Repeat(`\`, len(backslashes)/2)
		// expression was escaped
		if len(backslashes)%2 == 1 {
			return prefix + char + `{` + source + `}`
		}
		if char == environmentVar {
			val, ok := environment[source]
			if !ok {
				// undefined variables are left unchanged
				return prefix + char + `{` + source + `}`
			}
			return prefix + val
		}
		val, err := context.EvaluateExpression(source)
		if err != nil {
			errs = append(errs, err)
			return ""
		}
		str, err := PropertyToString(val, false)
		if err != nil {
			errs = append(errs, err)
			return ""
		}
		return prefix + str
	})
	if len(errs) > 0 {
		return "", errs[0]
	}
	return replaced, nil
}

// LoadDotEnv loads .env file
// - filename: the name of the file to load
// - buildDir: the directory relative paths are resolved from
// Return:
// - a map of environment variables
// - an error if something went wrong
func LoadDotEnv(filename, buildDir string) (map[string]string, error) {
	filename = util.ExpandUserHome(filename)
	if !filepath.IsAbs(filename) && buildDir != "" {
		filename = filepath.Join(buildDir, filename)
	}
	filename = filepath.Clean(filename)
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening dotenv file '%s': %w", filename, err)
	}
	defer func() {
		_ = file.Close()
	}()
	environment, err := ReadDotEnv(file)
	if err != nil {
		return nil, fmt.Errorf("reading dotenv file '%s': %w", filename, err)
	}
	return environment, nil
}

// ReadDotEnv parses dotenv content: lines such as NAME=value, with optional
// "export " prefix and optional single or double quotes around value. Empty
// lines and lines starting with '#' are ignored.
// - reader: the reader for dotenv content
// Return:
// - a map of environment variables
// - an error if something went wrong
func ReadDotEnv(reader io.Reader) (map[string]string, error) {
	environment := make(map[string]string)
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		index := strings.Index(line, "=")
		if index < 0 {
			return nil, fmt.Errorf("bad environment line: '%s'", line)
		}
		name := strings.TrimSpace(line[:index])
		value := strings.TrimSpace(line[index+1:])
		if len(value) >= 2 && (value[0] == '"' || value[0] == '\'') && value[len(value)-1] == value[0] {
			value = value[1 : len(value)-1]
		}
		environment[name] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return environment, nil
}

// Message prints a message on the console
// - text: the text to print on console
func (context *Context) Message(text string) {
	Message(text)
}

// MessageArgs prints a message on the console
// - text: the text to print on console, with formatting verbs
// - args: the arguments (as for fmt.Printf())
func (context *Context) MessageArgs(text string, args ...interface{}) {
	MessageArgs(text, args...)
}

// FormatScriptError adds line and column numbers on parser or vm errors.
// - err: the error to process
// Return: the processed error
func FormatScriptError(err error) error {
	if e, ok := err.(*vm.Error); ok {
		return fmt.Errorf("%s (at line %d, column %d)", err, e.Pos.Line, e.Pos.Column)
	} else {
		return err
	}
}
