package build

import (
	"fmt"
	"neon/util"
	"sort"
	"strings"
	"unicode/utf8"
)

// Print help on build
func (build *Build) Info(context *Context) error {
	// print build information
	if build.Doc != "" {
		Message("doc: %s", build.Doc)
	}
	if len(build.Default) > 0 {
		defaults := "[" + strings.Join(build.Default, ", ") + "]"
		Message("default: %s", defaults)
	}
	Message("repository: %s", build.Repository)
	if build.Singleton != 0 {
		Message("singleton: %d", build.Singleton)
	}
	// print parent build files
	if len(build.Parents) > 0 {
		Message("")
		Message("extends:")
		for _, extend := range build.Extends {
			Message("- %s", extend)
		}
	}
	// print configuration files
	if len(build.Config) > 0 {
		Message("")
		Message("configuration:")
		for _, config := range build.Config {
			Message("- %s", config)
		}
	}
	// print context scripts
	if len(build.Scripts) > 0 {
		Message("")
		Message("context:")
		for _, script := range build.Scripts {
			Message("- %s", script)
		}
	}
	// print build properties
	length := util.MaxLength(context.Properties)
	if len(context.Properties) > 0 {
		Message("")
		Message("properties:")
		for _, name := range context.Properties {
			value, err := context.GetProperty(name)
			if err != nil {
				return fmt.Errorf("getting property '%s': %v", name, err)
			}
			valueStr, err := PropertyToString(value, true)
			if err != nil {
				return fmt.Errorf("formatting property '%s': %v", name, err)
			}
			PrintProperty(name, valueStr, []string{}, length)
		}
	}
	// print build environment
	var names []string
	for name := range context.Environment {
		names = append(names, name)
	}
	length = util.MaxLength(names)
	sort.Strings(names)
	if len(context.Environment) > 0 {
		Message("")
		Message("environment:")
		for _, name := range names {
			value := "\"" + context.Environment[name] + "\""
			PrintProperty(name, value, []string{}, length)
		}
	}
	// print targets documentation
	targets := build.GetTargets()
	names = make([]string, 0)
	for name := range targets {
		names = append(names, name)
	}
	length = util.MaxLength(names)
	sort.Strings(names)
	if len(names) > 0 {
		Message("")
		Message("targets:")
		for _, name := range names {
			target := targets[name]
			PrintProperty(name, target.Doc, target.Depends, length)
		}
	}
	return nil
}

func PrintProperty(name, doc string, depends []string, length int) {
	deps := ""
	if len(depends) > 0 {
		deps = "[" + strings.Join(depends, ", ") + "]"
	}
	if doc != "" {
		deps = " " + deps
	}
	Message("  %s: %s%s%s", name,
		strings.Repeat(" ", length-utf8.RuneCountInString(name)), doc, deps)
}

// Print build targets
func (build *Build) PrintTargets() {
	var targets []string
	for name := range build.GetTargets() {
		targets = append(targets, name)
	}
	sort.Strings(targets)
	Message(strings.Join(targets, " "))
}

// Print tasks
func PrintTasks() {
	var tasks []string
	for name := range TaskMap {
		tasks = append(tasks, name)
	}
	sort.Strings(tasks)
	Message(strings.Join(tasks, " "))
}

// Print help on tasks
func PrintHelpTask(task string) {
	descriptor, found := TaskMap[task]
	if found {
		Message(descriptor.Help)
	} else {
		Message("Task '%s' was not found", task)
	}
}

// Print builtins
func PrintBuiltins() {
	var builtins []string
	for name := range BuiltinMap {
		builtins = append(builtins, name)
	}
	sort.Strings(builtins)
	Message(strings.Join(builtins, " "))
}

// Print help on builtins
func PrintHelpBuiltin(builtin string) {
	descriptor, found := BuiltinMap[builtin]
	if found {
		Message(descriptor.Help)
	} else {
		Message("Builtin '%s' was not found", builtin)
	}
}

// Print markdown reference for tasks and builtins
func PrintReference() {
	Message("Tasks Reference")
	Message("===============")
	Message("")
	var tasks []string
	for name := range TaskMap {
		tasks = append(tasks, name)
	}
	sort.Strings(tasks)
	for _, task := range tasks {
		Message(task)
		Message(strings.Repeat("-", len(task)))
		Message("")
		Message(TaskMap[task].Help)
		Message("")
	}
	Message("")
	Message("Builtins Reference")
	Message("==================")
	Message("")
	var builtins []string
	for name := range BuiltinMap {
		builtins = append(builtins, name)
	}
	sort.Strings(builtins)
	for _, builtin := range builtins {
		Message(builtin)
		Message(strings.Repeat("-", len(builtin)))
		Message("")
		Message(BuiltinMap[builtin].Help)
		Message("")
	}
}
