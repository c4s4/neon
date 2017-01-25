package build

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"neon/util"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"
)

type Build struct {
	File    string
	Dir     string
	Here    string
	Name    string
	Default []string
	Doc     string
	Context *Context
	Targets map[string]*Target
	Verbose bool
}

func NewBuild(file string, verbose bool) (*Build, error) {
	build := &Build{}
	build.Verbose = verbose
	build.Debug("Loading build file '%s'", file)
	source, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("loading build file '%s': %v", file, err)
	}
	build.Debug("Parsing build file")
	var object util.Object
	err = yaml.Unmarshal(source, &object)
	if err != nil {
		return nil, fmt.Errorf("build must be a map with string keys")
	}
	build.Debug("Build structure: %#v", object)
	build.Debug("Reading build first level fields")
	err = object.CheckFields([]string{"name", "default", "doc", "properties",
		"environment", "targets"})
	if err != nil {
		return nil, fmt.Errorf("parsing build file: %v", err)
	}
	if object.HasField("name") {
		str, err := object.GetString("name")
		if err == nil {
			build.Name = str
		}
	}
	if object.HasField("default") {
		list, err := object.GetListStringsOrString("default")
		if err == nil {
			build.Default = list
		}
	}
	if object.HasField("doc") {
		str, err := object.GetString("doc")
		if err == nil {
			build.Doc = str
		}
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return nil, fmt.Errorf("getting build file path: %v", err)
	}
	build.File = filepath.Base(path)
	build.Dir = filepath.Dir(path)
	here, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getting current directory: %v", err)
	}
	build.Here = here
	properties := make(map[string]interface{})
	if object.HasField("properties") {
		properties, err = object.GetObject("properties")
		if err != nil {
			return nil, fmt.Errorf("parsing properties: %v", err)
		}
	}
	environment := make(map[string]interface{})
	if object.HasField("environment") {
		environment, err = object.GetObject("environment")
		if err != nil {
			return nil, fmt.Errorf("parsing environmen: %v", err)
		}
	}
	context, err := NewContext(build, properties, environment)
	if err != nil {
		return nil, fmt.Errorf("building context: %v", err)
	}
	build.Context = context
	targets := util.Object(make(map[string]interface{}))
	if object.HasField("targets") {
		targets, err = object.GetObject("targets")
		if err != nil {
			return nil, fmt.Errorf("parsing targets: %v", err)
		}
	}
	build.Targets = make(map[string]*Target)
	for name, _ := range targets {
		object, err := targets.GetObject(name)
		if err != nil {
			return nil, fmt.Errorf("parsing target '%s': %v", name, err)
		}
		target, err := NewTarget(build, name, object)
		if err != nil {
			return nil, fmt.Errorf("parsing target '%s': %v", name, err)
		}
		build.Targets[name] = target
	}
	return build, nil
}

func (build *Build) Run(targets []string) error {
	if len(targets) == 0 {
		if len(build.Default) == 0 {
			return fmt.Errorf("no default target")
		}
		targets = build.Default
	}
	for _, target := range targets {
		err := build.RunTarget(target, NewStack())
		if err != nil {
			return err
		}
	}
	return nil
}

func (build *Build) RunTarget(name string, stack *Stack) error {
	target, ok := build.Targets[name]
	if !ok {
		return fmt.Errorf("target '%s' not found", name)
	}
	err := target.Run(stack)
	if err != nil {
		return fmt.Errorf("running target '%s': %v", name, err)
	}
	return nil
}

func (build *Build) Help() error {
	newLine := false
	// print build documentation
	if build.Doc != "" {
		build.Info(build.Doc)
		newLine = true
	}
	// print build properties
	length := maxLength(build.Context.Properties)
	if len(build.Context.Properties) > 0 {
		if newLine {
			build.Info("")
		}
		build.Info("Properties:")
		for _, name := range build.Context.Properties {
			value, err := build.Context.GetProperty(name)
			if err != nil {
				return fmt.Errorf("getting property '%s': %v", name, err)
			}
			valueStr, err := PropertyToString(value, true)
			if err != nil {
				return fmt.Errorf("formatting property '%s': %v", name, err)
			}
			build.PrintColorLine(name, valueStr, []string{}, length)
		}
		newLine = true
	}
	// print build environment
	var names []string
	for name, _ := range build.Context.Environment {
		names = append(names, name)
	}
	length = maxLength(names)
	sort.Strings(names)
	if len(build.Context.Environment) > 0 {
		if newLine {
			build.Info("")
		}
		build.Info("Environment:")
		for _, name := range names {
			value := "\"" + build.Context.Environment[name] + "\""
			build.PrintColorLine(name, value, []string{}, length)
		}
		newLine = true
	}
	// print targets documentation
	names = make([]string, 0)
	for name, _ := range build.Targets {
		names = append(names, name)
	}
	length = maxLength(names)
	sort.Strings(names)
	if len(names) > 0 {
		if newLine {
			build.Info("")
		}
		build.Info("Targets:")
		for _, name := range names {
			target := build.Targets[name]
			build.PrintColorLine(name, target.Doc, target.Depends, length)
		}
	}
	return nil
}

func (build *Build) PrintTargets() {
	var targets []string
	for name, _ := range build.Targets {
		targets = append(targets, name)
	}
	sort.Strings(targets)
	build.Info(strings.Join(targets, " "))
}

func (build *Build) PrintTasks() {
	var tasks []string
	for name, _ := range TaskMap {
		tasks = append(tasks, name)
	}
	sort.Strings(tasks)
	build.Info(strings.Join(tasks, " "))
}

func (build *Build) PrintHelpTask(task string) {
	descriptor, found := TaskMap[task]
	if found {
		build.Info(descriptor.Help)
	} else {
		build.Info("Task '%s' was not found", task)
	}
}

func (build *Build) PrintBuiltins() {
	var builtins []string
	for name, _ := range BuiltinMap {
		builtins = append(builtins, name)
	}
	sort.Strings(builtins)
	build.Info(strings.Join(builtins, " "))
}

func (build *Build) PrintHelpBuiltin(builtin string) {
	descriptor, found := BuiltinMap[builtin]
	if found {
		build.Info(descriptor.Help)
	} else {
		build.Info("Builtin '%s' was not found", builtin)
	}
}

func (build *Build) Info(message string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(message, args...))
}

func (build *Build) Debug(message string, args ...interface{}) {
	if build.Verbose {
		fmt.Println(fmt.Sprintf(message, args...))
	}
}

func (build *Build) PrintColorLine(name, doc string, depends []string, length int) {
	deps := ""
	if len(depends) > 0 {
		deps = "[" + strings.Join(depends, ", ") + "]"
	}
	if doc != "" {
		deps = " " + deps
	}
	util.PrintColor("%s%s %s%s", util.Yellow(name),
		strings.Repeat(" ", length-utf8.RuneCountInString(name)), doc, deps)
}

func maxLength(lines []string) int {
	length := 0
	for _, line := range lines {
		if utf8.RuneCountInString(line) > length {
			length = utf8.RuneCountInString(line)
		}
	}
	return length
}
