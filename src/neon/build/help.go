package build

import (
	"fmt"
	"neon/util"
	"sort"
	"strings"
	"unicode/utf8"
)

// Print help on build
func (build *Build) Help() error {
	// print build information
	if build.Name != "" {
		Info("name: %s", build.Name)
	}
	if build.Doc != "" {
		Info("doc: %s", build.Doc)
	}
	if len(build.Default) > 0 {
		defauls := "[" + strings.Join(build.Default, ", ") + "]"
		Info("default: %s", defauls)
	}
	Info("repository: %s", build.Repository)
	if build.Singleton != 0 {
		Info("singleton: %d", build.Singleton)
	}
	// print parent build files
	if len(build.Parents) > 0 {
		Info("")
		Info("extends:")
		for _, extend := range build.Extends {
			Info("- %s", extend)
		}
	}
	// print configuration files
	if len(build.Config) > 0 {
		Info("")
		Info("configuration:")
		for _, config := range build.Config {
			Info("- %s", config)
		}
	}
	// print context scripts
	if len(build.Config) > 0 {
		Info("")
		Info("context:")
		for _, script := range build.Scripts {
			Info("- %s", script)
		}
	}
	// print build properties
	length := util.MaxLength(build.Context.Properties)
	if len(build.Context.Properties) > 0 {
		Info("")
		Info("properties:")
		for _, name := range build.Context.Properties {
			value, err := build.Context.GetProperty(name)
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
	for name, _ := range build.Context.Environment {
		names = append(names, name)
	}
	length = util.MaxLength(names)
	sort.Strings(names)
	if len(build.Context.Environment) > 0 {
		Info("")
		Info("environment:")
		for _, name := range names {
			value := "\"" + build.Context.Environment[name] + "\""
			PrintProperty(name, value, []string{}, length)
		}
	}
	// print targets documentation
	targets := build.GetTargets()
	names = make([]string, 0)
	for name, _ := range targets {
		names = append(names, name)
	}
	length = util.MaxLength(names)
	sort.Strings(names)
	if len(names) > 0 {
		Info("")
		Info("targets:")
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
	Info("  %s: %s%s%s", name,
		strings.Repeat(" ", length-utf8.RuneCountInString(name)), doc, deps)
}

// Print build targets
func (build *Build) PrintTargets() {
	var targets []string
	for name, _ := range build.GetTargets() {
		targets = append(targets, name)
	}
	sort.Strings(targets)
	Info(strings.Join(targets, " "))
}

// Print tasks
func PrintTasks() {
	var tasks []string
	for name, _ := range TaskMap {
		tasks = append(tasks, name)
	}
	sort.Strings(tasks)
	Info(strings.Join(tasks, " "))
}

// Print help on tasks
func PrintHelpTask(task string) {
	descriptor, found := TaskMap[task]
	if found {
		Info(descriptor.Help)
	} else {
		Info("Task '%s' was not found", task)
	}
}

// Print builtins
func PrintBuiltins() {
	var builtins []string
	for name, _ := range BuiltinMap {
		builtins = append(builtins, name)
	}
	sort.Strings(builtins)
	Info(strings.Join(builtins, " "))
}

// Print help on builtins
func PrintHelpBuiltin(builtin string) {
	descriptor, found := BuiltinMap[builtin]
	if found {
		Info(descriptor.Help)
	} else {
		Info("Builtin '%s' was not found", builtin)
	}
}

// Print markdown reference for tasks and builtins
func PrintReference() {
	fmt.Println("Tasks Reference")
	fmt.Println("===============")
	fmt.Println()
	var tasks []string
	for name := range TaskMap {
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
	for name := range BuiltinMap {
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
