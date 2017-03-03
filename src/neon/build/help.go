package build

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"
)

func (build *Build) Help() error {
	newLine := false
	// print build documentation
	if build.Doc != "" {
		Info(build.Doc)
		newLine = true
	}
	// print build properties
	length := maxLength(build.Context.Properties)
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
	length = maxLength(names)
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
	length = maxLength(names)
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

func (build *Build) PrintTargets() {
	var targets []string
	for name, _ := range build.GetTargets() {
		targets = append(targets, name)
	}
	sort.Strings(targets)
	Info(strings.Join(targets, " "))
}

func (build *Build) PrintTasks() {
	var tasks []string
	for name, _ := range TaskMap {
		tasks = append(tasks, name)
	}
	sort.Strings(tasks)
	Info(strings.Join(tasks, " "))
}

func (build *Build) PrintHelpTask(task string) {
	descriptor, found := TaskMap[task]
	if found {
		Info(descriptor.Help)
	} else {
		Info("Task '%s' was not found", task)
	}
}

func PrintBuiltins() {
	var builtins []string
	for name, _ := range BuiltinMap {
		builtins = append(builtins, name)
	}
	sort.Strings(builtins)
	Info(strings.Join(builtins, " "))
}

func PrintHelpBuiltin(builtin string) {
	descriptor, found := BuiltinMap[builtin]
	if found {
		Info(descriptor.Help)
	} else {
		Info("Builtin '%s' was not found", builtin)
	}
}

func PrintReference() {
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
