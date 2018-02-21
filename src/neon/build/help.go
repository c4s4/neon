package build

import (
	"fmt"
	"neon/util"
	"sort"
	"strings"
	"unicode/utf8"
)

// Info prints information about build on console.
// - context: context of the build
// Return: error if something went wrong
func (build *Build) Info(context *Context) error {
	build.infoDoc()
	build.infoDefault()
	build.infoRepository()
	build.infoSingleton(context)
	build.infoParents()
	build.infoConfiguration()
	build.infoContext()
	if err := build.infoProperties(context); err != nil {
		return err
	}
	build.infoEnvironment()
	build.infoTargets()
	return nil
}

func (build *Build) infoDoc() {
	if build.Doc != "" {
		Message("doc: %s", build.Doc)
	}
}

func (build *Build) infoDefault() {
	if len(build.Default) > 0 {
		defaults := "[" + strings.Join(build.Default, ", ") + "]"
		Message("default: %s", defaults)
	}
}

func (build *Build) infoRepository() {
	Message("repository: %s", build.Repository)
}

func (build *Build) infoSingleton(context *Context) {
	if build.Singleton != "" {
		port, err := context.EvaluateExpression(build.Singleton)
		if err == nil {
			Message("singleton: %d", port)
		}
	}
}

func (build *Build) infoParents() {
	if len(build.Parents) > 0 {
		Message("")
		Message("extends:")
		for _, extend := range build.Extends {
			Message("- %s", extend)
		}
	}
}

func (build *Build) infoConfiguration() {
	if len(build.Config) > 0 {
		Message("")
		Message("configuration:")
		for _, config := range build.Config {
			Message("- %s", config)
		}
	}
}

func (build *Build) infoContext() {
	if len(build.Scripts) > 0 {
		Message("")
		Message("context:")
		for _, script := range build.Scripts {
			Message("- %s", script)
		}
	}
}

func (build *Build) infoProperties(context *Context) error {
	length := util.MaxLineLength(build.Properties.Fields())
	if len(build.Properties) > 0 {
		Message("")
		Message("properties:")
		for name := range build.Properties {
			value, err := context.GetProperty(name)
			if err != nil {
				return fmt.Errorf("getting property '%s': %v", name, err)
			}
			valueStr, err := PropertyToString(value, true)
			if err != nil {
				return fmt.Errorf("formatting property '%s': %v", name, err)
			}
			PrintTarget(name, valueStr, []string{}, length)
		}
	}
	return nil
}

func (build *Build) infoEnvironment() {
	var names []string
	for name := range build.Environment {
		names = append(names, name)
	}
	length := util.MaxLineLength(names)
	sort.Strings(names)
	if len(build.Environment) > 0 {
		Message("")
		Message("environment:")
		for _, name := range names {
			value := "\"" + build.Environment[name] + "\""
			PrintTarget(name, value, []string{}, length)
		}
	}
}

func (build *Build) infoTargets() {
	targets := build.GetTargets()
	names := make([]string, 0)
	for name := range targets {
		names = append(names, name)
	}
	length := util.MaxLineLength(names)
	sort.Strings(names)
	if len(names) > 0 {
		Message("")
		Message("targets:")
		for _, name := range names {
			target := targets[name]
			PrintTarget(name, target.Doc, target.Depends, length)
		}
	}
}

// PrintTarget prints target documentation on console
// - name: the name of the target
// - doc: the target documentation
// - depends: targets on which this one depends
// - length: title length to align help on targets
func PrintTarget(name, doc string, depends []string, length int) {
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

// PrintTargets prints help on targets
func (build *Build) PrintTargets() {
	var targets []string
	for name := range build.GetTargets() {
		targets = append(targets, name)
	}
	sort.Strings(targets)
	Message(strings.Join(targets, " "))
}

// PrintTasks prints the list of tasks on the console.
func PrintTasks() {
	var tasks []string
	for name := range TaskMap {
		tasks = append(tasks, name)
	}
	sort.Strings(tasks)
	Message(strings.Join(tasks, " "))
}

// PrintHelpTask prints help on given task.
// - task: name of the task to document.
func PrintHelpTask(task string) {
	descriptor, found := TaskMap[task]
	if found {
		Message(descriptor.Help)
	} else {
		Message("Func '%s' was not found", task)
	}
}

// PrintBuiltins prints the list of all builtins on console
func PrintBuiltins() {
	var builtins []string
	for name := range BuiltinMap {
		builtins = append(builtins, name)
	}
	sort.Strings(builtins)
	Message(strings.Join(builtins, " "))
}

// PrintHelpBuiltin prints help on given builtin:
// - builtin: the name of the builtin to document.
func PrintHelpBuiltin(builtin string) {
	descriptor, found := BuiltinMap[builtin]
	if found {
		Message(descriptor.Help)
	} else {
		Message("Builtin '%s' was not found", builtin)
	}
}

// PrintReference prints markdown reference for tasks and builtins on console.
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
