package build

import (
	"fmt"
	anko_core "github.com/mattn/anko/builtins"
	"github.com/mattn/anko/vm"
	zglob "github.com/mattn/go-zglob"
	"io/ioutil"
	"neon/util"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

type Context struct {
	VM          *vm.Env
	Build       *Build
	Properties  []string
	Environment map[string]string
}

func NewContext(build *Build) (*Context, error) {
	vm := vm.NewEnv()
	anko_core.LoadAllBuiltins(vm)
	LoadBuiltins(vm)
	properties := build.GetProperties()
	environment := build.GetEnvironment()
	context := &Context{
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
			return nil, fmt.Errorf("evaluating script '%s': %v", script, err)
		}
	}
	err := context.SetProperties(properties)
	if err != nil {
		return nil, fmt.Errorf("evaluating properties: %v", err)
	}
	return context, nil
}

func (context *Context) Evaluate(source string) (interface{}, error) {
	value, err := context.VM.Execute(source)
	if err != nil {
		return nil, err
	}
	if value.IsValid() {
		return value.Interface(), nil
	} else {
		return nil, nil
	}
}

func (context *Context) SetProperty(name string, value interface{}) {
	context.VM.Define(name, value)
}

func (context *Context) GetProperty(name string) (interface{}, error) {
	value, err := context.VM.Get(name)
	if err != nil {
		return nil, err
	}
	if value.IsValid() {
		return value.Interface(), nil
	} else {
		return nil, nil
	}
}

func (context *Context) SetProperties(object util.Object) error {
	context.SetProperty("OS", runtime.GOOS)
	context.SetProperty("ARCH", runtime.GOARCH)
	context.SetProperty("BASE", context.Build.Dir)
	context.SetProperty("HERE", context.Build.Here)
	todo := object.Fields()
	var crash error
	for len(todo) > 0 {
		var done []string
		for _, name := range todo {
			value := object[name]
			str, ok := value.(string)
			if ok {
				eval, err := context.ReplaceProperties(str)
				if err == nil {
					context.SetProperty(name, eval)
					done = append(done, name)
				} else {
					crash = err
				}
			} else {
				context.SetProperty(name, value)
				done = append(done, name)
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

func (context *Context) ReplaceProperties(text string) (string, error) {
	r := regexp.MustCompile("#{.*?}")
	var err error
	replaced := r.ReplaceAllStringFunc(text, func(expression string) string {
		name := expression[2 : len(expression)-1]
		var value interface{}
		value, err = context.Evaluate(name)
		if err != nil {
			return ""
		} else {
			var str string
			str, err = PropertyToString(value, false)
			if err != nil {
				return ""
			} else {
				return str
			}
		}
	})
	return replaced, err
}

func (context *Context) GetEnvironment() ([]string, error) {
	environment := make(map[string]string)
	for _, line := range os.Environ() {
		index := strings.Index(line, "=")
		name := line[:index]
		value := line[index+1:]
		environment[name] = value
	}
	environment["BASE"] = context.Build.Dir
	environment["HERE"] = context.Build.Here
	var variables []string
	for name, _ := range context.Environment {
		variables = append(variables, name)
	}
	sort.Strings(variables)
	for _, name := range variables {
		value := context.Environment[name]
		r := regexp.MustCompile("[\\$#]{.*?}")
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
				value, err := context.GetProperty(name)
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

func (context *Context) FindFiles(dir string, includes, excludes []string) ([]string, error) {
	eval, err := context.ReplaceProperties(dir)
	if err != nil {
		return nil, fmt.Errorf("evaluating source directory: %v", err)
	}
	dir = eval
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
		pattern, err := context.ReplaceProperties(include)
		if err != nil {
			return nil, fmt.Errorf("evaluating pattern: %v", err)
		}
		included = append(included, pattern)
	}
	var excluded []string
	for _, exclude := range excludes {
		pattern, err := context.ReplaceProperties(exclude)
		if err != nil {
			return nil, fmt.Errorf("evaluating pattern: %v", err)
		}
		excluded = append(excluded, pattern)
	}
	var candidates []string
	for _, include := range included {
		list, _ := zglob.Glob(include)
		for _, file := range list {
			stat, err := os.Stat(file)
			if err != nil {
				return nil, fmt.Errorf("stating file: %v", err)
			}
			if stat.Mode().IsRegular() {
				candidates = append(candidates, file)
			}
		}
	}
	var files []string
	if excluded != nil {
		for _, file := range candidates {
			for _, exclude := range excluded {
				match, err := zglob.Match(exclude, file)
				if !match && err == nil {
					files = append(files, file)
				}
			}
		}
	} else {
		files = candidates
	}
	sort.Strings(files)
	return files, nil
}

func PropertyToString(object interface{}, quotes bool) (string, error) {
	switch value := object.(type) {
	case bool:
		return strconv.FormatBool(value), nil
	case string:
		if quotes {
			return "\"" + value + "\"", nil
		} else {
			return value, nil
		}
	case int:
		return strconv.Itoa(value), nil
	case int32:
		return strconv.Itoa(int(value)), nil
	case int64:
		return strconv.Itoa(int(value)), nil
	case float64:
		return strconv.FormatFloat(value, 'g', -1, 64), nil
	default:
		if value == nil {
			return "null", nil
		}
		switch reflect.TypeOf(object).Kind() {
		case reflect.Slice:
			slice := reflect.ValueOf(object)
			elements := make([]string, slice.Len())
			for index := 0; index < slice.Len(); index++ {
				str, err := PropertyToString(slice.Index(index).Interface(), quotes)
				if err != nil {
					return "", err
				}
				elements[index] = str
			}
			return "[" + strings.Join(elements, ", ") + "]", nil
		case reflect.Map:
			dict := reflect.ValueOf(object)
			elements := make(map[string]string)
			var keys []string
			for _, key := range dict.MapKeys() {
				value := dict.MapIndex(key)
				keyStr, err := PropertyToString(key.Interface(), quotes)
				if err != nil {
					return "", err
				}
				keys = append(keys, keyStr)
				valueStr, err := PropertyToString(value.Interface(), quotes)
				if err != nil {
					return "", err
				}
				elements[keyStr] = valueStr
			}
			sort.Strings(keys)
			pairs := make([]string, len(keys))
			for index, key := range keys {
				pairs[index] = key + ": " + elements[key]
			}
			return "[" + strings.Join(pairs, ", ") + "]", nil
		default:
			return "", fmt.Errorf("no serializer for type '%T'", object)
		}
	}
}
