package build

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"neon/util"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"
)

type Build struct {
	File        string
	Dir         string
	Here        string
	Name        string
	Default     []string
	Doc         string
	Scripts     []string
	Verbose     bool
	Properties  util.Object
	Environment map[string]string
	Targets     map[string]*Target
	Context     *Context
	Parents     []*Build
	Index       *Index
}

func NewBuild(file string, verbose bool) (*Build, error) {
	build := &Build{}
	build.Verbose = verbose
	build.Debug("Loading build file '%s'", file)
	if strings.HasPrefix(file, "~/") {
		user, _ := user.Current()
		home := user.HomeDir
		file = filepath.Join(home, file[2:])
	}
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
	err = object.CheckFields([]string{"name", "doc", "default", "context",
		"extends", "properties", "environment", "targets"})
	if err != nil {
		return nil, fmt.Errorf("parsing build file: %v", err)
	}
	if object.HasField("name") {
		name, err := object.GetString("name")
		if err != nil {
			return nil, fmt.Errorf("getting build name: %v", err)
		}
		build.Name = name
	}
	if object.HasField("default") {
		list, err := object.GetListStringsOrString("default")
		if err != nil {
			return nil, fmt.Errorf("getting default targets: %v", err)
		}
		build.Default = list
	}
	if object.HasField("doc") {
		doc, err := object.GetString("doc")
		if err != nil {
			return nil, fmt.Errorf("getting build doc: %v", err)
		}
		build.Doc = doc
	}
	if object.HasField("context") {
		scripts, err := object.GetListStringsOrString("context")
		if err != nil {
			return nil, fmt.Errorf("getting context: %v", err)
		}
		build.Scripts = scripts
	}
	if object.HasField("extends") {
		parents, err := object.GetListStringsOrString("extends")
		if err != nil {
			return nil, fmt.Errorf("parsing parents: %v", err)
		}
		var extends []*Build
		for _, parent := range parents {
			parent, err = ExpandNeonPath(parent)
			if err != nil {
				return nil, fmt.Errorf("expanding neon path: %v", err)
			}
			extend, err := NewBuild(parent, verbose)
			if err != nil {
				return nil, fmt.Errorf("parsing parent '%s': %v", parent, err)
			}
			extends = append(extends, extend)
		}
		build.Parents = extends
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
	build.Properties = properties
	environment := make(map[string]string)
	if object.HasField("environment") {
		env, err := object.GetObject("environment")
		if err != nil {
			return nil, fmt.Errorf("parsing environmen: %v", err)
		}
		environment, err = env.ToMapStringString()
		if err != nil {
			return nil, fmt.Errorf("getting environment: %v", err)
		}
	}
	build.Environment = environment
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

func (build *Build) GetProperties() util.Object {
	var properties = make(map[string]interface{})
	for _, parent := range build.Parents {
		for name, value := range parent.GetProperties() {
			properties[name] = value
		}
	}
	for name, value := range build.Properties {
		properties[name] = value
	}
	return properties
}

func (build *Build) GetEnvironment() map[string]string {
	var environment = make(map[string]string)
	for _, parent := range build.Parents {
		for name, value := range parent.GetEnvironment() {
			environment[name] = value
		}
	}
	for name, value := range build.Environment {
		environment[name] = value
	}
	return environment
}

func (build *Build) GetTargets() map[string]*Target {
	var targets = make(map[string]*Target)
	for _, parent := range build.Parents {
		for name, target := range parent.GetTargets() {
			targets[name] = target
		}
	}
	for name, target := range build.Targets {
		targets[name] = target
	}
	return targets
}

func (build *Build) SetContext(context *Context) {
	build.Context = context
	for _, parent := range build.Parents {
		parent.SetContext(context)
	}
}

func (build *Build) SetDir(dir string) {
	build.Dir = dir
	for _, parent := range build.Parents {
		parent.SetDir(dir)
	}
}

func (build *Build) Init() error {
	context, err := NewContext(build)
	if err != nil {
		return fmt.Errorf("evaluating context: %v", err)
	}
	build.SetDir(build.Dir)
	build.SetContext(context)
	return nil
}

func (build *Build) GetDefault() []string {
	if len(build.Default) > 0 {
		return build.Default
	} else {
		for _, parent := range build.Parents {
			if len(parent.Default) > 0 {
				return parent.Default
			}
		}
		for _, parent := range build.Parents {
			parentDefault := parent.GetDefault()
			if len(parentDefault) > 0 {
				return parentDefault
			}
		}
	}
	return build.Default
}

func (build *Build) Run(targets []string) error {
	if len(targets) == 0 {
		targets = build.GetDefault()
		if len(targets) == 0 {
			return fmt.Errorf("no default target")
		}
	}
	for _, target := range targets {
		err := build.RunTarget(target, NewStack())
		if err != nil {
			return err
		}
	}
	return nil
}

func (build *Build) GetTarget(name string) *Target {
	target, found := build.Targets[name]
	if found {
		return target
	} else {
		for _, parent := range build.Parents {
			target = parent.GetTarget(name)
			if target != nil {
				return target
			}
		}
	}
	return nil
}

func (build *Build) RunTarget(name string, stack *Stack) error {
	target := build.GetTarget(name)
	if target == nil {
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
	targets := build.GetTargets()
	names = make([]string, 0)
	for name, _ := range targets {
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
			target := targets[name]
			build.PrintColorLine(name, target.Doc, target.Depends, length)
		}
	}
	return nil
}

func (build *Build) PrintTargets() {
	var targets []string
	for name, _ := range build.GetTargets() {
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

func (build *Build) PrintReference() {
	fmt.Println("Tasks Reference")
	fmt.Println("===============")
	fmt.Println()
	var tasks []string
	for name, _ := range TaskMap {
		tasks = append(tasks, name)
	}
	sort.Strings(tasks)
	for _, task := range tasks {
		fmt.Println(task)
		fmt.Println(strings.Repeat("-", len(task)))
		fmt.Println()
		fmt.Println(TaskMap[task].Help)
		fmt.Println()
	}
	fmt.Println()
	fmt.Println("Builtins Reference")
	fmt.Println("==================")
	fmt.Println()
	var builtins []string
	for name, _ := range BuiltinMap {
		builtins = append(builtins, name)
	}
	sort.Strings(builtins)
	for _, builtin := range builtins {
		fmt.Println(builtin)
		fmt.Println(strings.Repeat("-", len(builtin)))
		fmt.Println()
		fmt.Println(BuiltinMap[builtin].Help)
		fmt.Println()
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

func ExpandNeonPath(path string) (string, error) {
	if strings.HasPrefix(path, ":") {
		parts := strings.Split(path[1:], "/")
		if len(parts) < 2 || len(parts) > 3 {
			return "", fmt.Errorf("Bad Neon path '%s'", path)
		}
		if len(parts) == 2 {
			parts = []string{parts[0], "latest", parts[1]}
		}
		return fmt.Sprintf("~/.neon/%s/%s/%s", parts[0], parts[1], parts[2]), nil
	} else {
		return path, nil
	}
}
