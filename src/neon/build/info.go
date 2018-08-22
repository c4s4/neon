package build

import (
	"fmt"
	"neon/util"
	"path"
	"sort"
	"strings"
	"unicode/utf8"
)

// Info generates information about build on console.
// - context: context of the build
// Return: build info as a string and an error if something went wrong
func (build *Build) Info(context *Context) (string, error) {
	info := ""
	info += build.infoDoc()
	info += build.infoDefault()
	info += build.infoRepository()
	info += build.infoSingleton(context)
	info += build.infoParents()
	info += build.infoConfiguration()
	info += build.infoContext()
	info += "\n"
	envs := build.infoEnvironment()
	if envs != "" {
		info += envs + "\n"
	}
	props, err := build.infoProperties(context)
	if err != nil {
		return "", err
	}
	if props != "" {
		info += props + "\n"
	}
	targets := build.infoTargets()
	if targets != "" {
		info += targets
	}
	return strings.TrimSpace(info), nil
}

func (build *Build) infoDoc() string {
	info := ""
	if build.Doc != "" {
		info += "doc: " + build.Doc + "\n"
	}
	return info
}

func (build *Build) infoDefault() string {
	info := ""
	if len(build.Default) > 0 {
		defaults := "[" + strings.Join(build.Default, ", ") + "]"
		info += "default: " + defaults + "\n"
	}
	return info
}

func (build *Build) infoRepository() string {
	return "repository: " + build.Repository + "\n"
}

func (build *Build) infoSingleton(context *Context) string {
	info := ""
	if build.Singleton != "" {
		port, err := context.EvaluateExpression(build.Singleton)
		if err == nil {
			info += fmt.Sprintf("singleton: %v\n", port)
		}
	}
	return info
}

func (build *Build) infoParents() string {
	info := ""
	if len(build.Parents) > 0 {
		info += "extends:\n"
		for _, extend := range build.Extends {
			info += "- " + extend + "\n"
		}
	}
	return info
}

func (build *Build) infoConfiguration() string {
	info := ""
	if len(build.Config) > 0 {
		info += "configuration:\n"
		for _, config := range build.Config {
			info += "- " + config + "\n"
		}
	}
	return info
}

func (build *Build) infoContext() string {
	info := ""
	if len(build.Scripts) > 0 {
		info += "context:\n"
		for _, script := range build.Scripts {
			info += "- " + script + "\n"
		}
	}
	return info
}

func (build *Build) infoProperties(context *Context) (string, error) {
	info := ""
	length := util.MaxLineLength(build.Properties.Fields())
	if len(build.Properties) > 0 {
		info += "properties:\n"
		names := make([]string, 0, len(build.Properties))
		for name := range build.Properties {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			value, err := context.GetProperty(name)
			if err != nil {
				return "", fmt.Errorf("getting property '%s': %v", name, err)
			}
			valueStr, err := PropertyToString(value, true)
			if err != nil {
				return "", fmt.Errorf("formatting property '%s': %v", name, err)
			}
			info += FormatTarget(name, valueStr, []string{}, length) + "\n"
		}
	}
	return info, nil
}

func (build *Build) infoEnvironment() string {
	info := ""
	var names []string
	for name := range build.Environment {
		names = append(names, name)
	}
	length := util.MaxLineLength(names)
	sort.Strings(names)
	if len(build.Environment) > 0 {
		info += "environment:\n"
		for _, name := range names {
			value := "\"" + build.Environment[name] + "\""
			info += FormatTarget(name, value, []string{}, length) + "\n"
		}
	}
	return info
}

func (build *Build) infoTargets() string {
	info := ""
	targets := build.GetTargets()
	names := make([]string, 0)
	for name := range targets {
		if !strings.HasPrefix(name, "_") {
			names = append(names, name)
		}
	}
	length := util.MaxLineLength(names)
	sort.Strings(names)
	if len(names) > 0 {
		info += "targets:\n"
		for _, name := range names {
			target := targets[name]
			info += FormatTarget(name, target.Doc, target.Depends, length) + "\n"
		}
	}
	return info
}

// FormatTarget generates target documentation on console
// - name: the name of the target
// - doc: the target documentation
// - depends: targets on which this one depends
// - length: title length to align help on targets
// Return: target as a string
func FormatTarget(name, doc string, depends []string, length int) string {
	deps := ""
	if len(depends) > 0 {
		deps = "[" + strings.Join(depends, ", ") + "]"
	}
	if doc != "" {
		deps = " " + deps
	}
	return fmt.Sprintf("  %s: %s%s%s", name,
		strings.Repeat(" ", length-utf8.RuneCountInString(name)), doc, deps)
}

// FormatTargets generates help on targets
// Return: targets as a string
func (build *Build) FormatTargets() string {
	var targets []string
	for name := range build.GetTargets() {
		targets = append(targets, name)
	}
	sort.Strings(targets)
	return strings.Join(targets, " ")
}

// InfoTasks generates the list of tasks on the console.
// Return: list of tasks as a string
func InfoTasks() string {
	var tasks []string
	for name := range TaskMap {
		tasks = append(tasks, name)
	}
	sort.Strings(tasks)
	return strings.Join(tasks, " ")
}

// InfoTask generates help on given task.
// - task: name of the task to document.
// Return: task info as a string
func InfoTask(task string) string {
	descriptor, found := TaskMap[task]
	if found {
		return descriptor.Help
	}
	return fmt.Sprintf("Func '%s' was not found", task)
}

// InfoBuiltins generates the list of all builtins on console
// Return: builtins list as a string
func InfoBuiltins() string {
	var builtins []string
	for name := range BuiltinMap {
		builtins = append(builtins, name)
	}
	sort.Strings(builtins)
	return strings.Join(builtins, " ")
}

// InfoBuiltin generates help on given builtin:
// - builtin: the name of the builtin to document.
// Return: builtin info as a string
func InfoBuiltin(builtin string) string {
	descriptor, found := BuiltinMap[builtin]
	if found {
		return descriptor.Help
	}
	return fmt.Sprintf("Builtin '%s' was not found", builtin)
}

// InfoThemes generates the list of all available themes.
// Return: info about themes as a string
func InfoThemes() string {
	var themes []string
	for theme := range Themes {
		themes = append(themes, theme)
	}
	sort.Strings(themes)
	return strings.Join(themes, " ")
}

// InfoTemplates generates list of templates in repository:
// - repository: the NeON repository (defaults to '~/.neon')
// Return: template info as a string
func InfoTemplates(repository string) string {
	info := ""
	files, err := FindTemplates(repository)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		info += file + "\n"
	}
	return strings.TrimSpace(info)
}

// InfoParents generates list of parent build files in repository:
// - repository: the NeON repository (defaults to '~/.neon')
// Return: parents info as a string
func InfoParents(repository string) string {
	info := ""
	files, err := util.FindFiles(repository, []string{"*/*/*.yml"}, nil, false)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		file = util.PathToUnix(file)
		name := path.Base(file)
		if name != "CHANGELOG.yml" && name != "build.yml" {
			info += file + "\n"
		}
	}
	return strings.TrimSpace(info)
}

// InfoReference generates markdown reference for tasks and builtins on console.
// Return: reference as a string
func InfoReference() string {
	info := "Tasks Reference\n"
	info += "===============\n\n"
	var tasks []string
	for name := range TaskMap {
		tasks = append(tasks, name)
	}
	sort.Strings(tasks)
	for _, task := range tasks {
		info += task + "\n"
		info += strings.Repeat("-", len(task)) + "\n\n"
		info += TaskMap[task].Help + "\n\n"
	}
	info += "\n"
	info += "Builtins Reference\n"
	info += "==================\n\n"
	var builtins []string
	for name := range BuiltinMap {
		builtins = append(builtins, name)
	}
	sort.Strings(builtins)
	for _, builtin := range builtins {
		info += builtin + "\n"
		info += strings.Repeat("-", len(builtin)) + "\n\n"
		info += BuiltinMap[builtin].Help + "\n"
	}
	return info
}
