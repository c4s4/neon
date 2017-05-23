package build

import (
	"fmt"
	"neon/util"
	"sort"
	"strings"
)

// Print help on build
func (build *Build) Help() error {
	newLine := false
	// print build documentation
	if build.Doc != "" {
		Info(build.Doc)
		newLine = true
	}
	// print build properties
	length := util.MaxLength(build.Context.Properties)
	if len(build.Context.Properties) > 0 {
		if newLine {
			Info("")
		}
		Info("Properties:")
		for _, name := range build.Context.Properties {
			value, err := build.Context.GetProperty(name)
			if err != nil {
				return fmt.Errorf("getting property '%s': %v", name, err)
			}
			valueStr, err := PropertyToString(value, true)
			if err != nil {
				return fmt.Errorf("formatting property '%s': %v", name, err)
			}
			PrintColorLine(name, valueStr, []string{}, length)
		}
		newLine = true
	}
	// print build environment
	var names []string
	for name, _ := range build.Context.Environment {
		names = append(names, name)
	}
	length = util.MaxLength(names)
	sort.Strings(names)
	if len(build.Context.Environment) > 0 {
		if newLine {
			Info("")
		}
		Info("Environment:")
		for _, name := range names {
			value := "\"" + build.Context.Environment[name] + "\""
			PrintColorLine(name, value, []string{}, length)
		}
		newLine = true
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
		if newLine {
			Info("")
		}
		Info("Targets:")
		for _, name := range names {
			target := targets[name]
			PrintColorLine(name, target.Doc, target.Depends, length)
		}
	}
	return nil
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
