package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "$",
		Func: Shell,
		Args: reflect.TypeOf(ShellArgs{}),
		Help: `Execute a command and return output and value.

Arguments:

- $: command to run as a string or a list of strings.
- =: name of the variable to store trimmed output into (optional, output to
  console if not set).

Examples:

    # execute ls command and get result in 'files' variable
    - $: 'ls -al'
      =: 'files'
    # execute command as a list of strings
    - $: ['ls', '-al']

Notes:

- Commands defined as a string run in the shell defined by shell field at the
  root of the build file (or 'sh -c' on Unix and 'cmd /c' on Windows by
  default).
- Defining a command as a list of strings is useful on Windows. Default shell on
  Windows is 'cmd' which can't properly manage arguments with spaces.
- Argument of a command defined as a list won't be expanded by shell. Thus
  $USER won't be expanded for instance.`,
	})
}

type ShellArgs struct {
	Shell []string `name:"$" wrap`
	To    string   `name:"=" optional`
}

func Shell(context *build.Context, args interface{}) error {
	params := args.(ShellArgs)
	// FIXME
	//commands, err := NewCommands(context.Build, params.Shell)
	commands, err := NewCommands(nil, params.Shell)
	if err != nil {
		return err
	}
	output, err := commands.Run(params.To == "", context)
	if err != nil {
		if output != "" {
			context.Message(output)
		}
		return err
	}
	if params.To != "" {
		context.SetProperty(params.To, strings.TrimSpace(string(output)))
	}
	return nil
}

// Commands lists commands by operating system
type Commands struct {
	Build    *build.Build
	Commands map[string]CommandFunc
}

// GetCommand return a command depending on current operating system and
// default command
func (c Commands) GetCommand() (CommandFunc, error) {
	for system, command := range c.Commands {
		if system != "default" && system == runtime.GOOS {
			return command, nil
		}
	}
	command, ok := c.Commands["default"]
	if !ok {
		return nil, fmt.Errorf("no command found for '%s'", runtime.GOOS)
	}
	return command, nil
}

// Run execute a command and return its output and an error (if command
// returned a value different from 0). Arguments:
// - pipe tells if we should print the output of the command on the console.
func (c Commands) Run(pipe bool, context *build.Context) (string, error) {
	command, err := c.GetCommand()
	if err != nil {
		return "", err
	}
	output, err := command.Run(c.Build, pipe, context)
	output = util.RemoveBlankLines(output)
	output = strings.TrimSuffix(output, "\n")
	return output, err
}

// Command is the interface for a command.
type CommandFunc interface {
	Run(build *build.Build, pipe bool, context *build.Context) (string, error)
}

// CommandList is a command as a list with executable and arguments.
type CommandList struct {
	Parts []string
}

// Run execute the command and returns its output and an error. Arguments:
// - pipe tells if we should print output of the command on the console.
func (c CommandList) Run(build *build.Build, pipe bool, context *build.Context) (string, error) {
	parts := make([]string, len(c.Parts))
	var err error
	for i := 0; i < len(c.Parts); i++ {
		parts[i], err = context.EvaluateString(c.Parts[i])
		if err != nil {
			return "", err
		}
	}
	command := exec.Command(parts[0], parts[1:]...)
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting current working directory: %v", err)
	}
	command.Dir = dir
	command.Env, err = context.EvaluateEnvironment(build)
	if err != nil {
		return "", fmt.Errorf("building environment: %v", err)
	}
	if pipe {
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		err = command.Run()
		if err != nil {
			return "", fmt.Errorf("executing command: %v", err)
		}
		return "", nil
	} else {
		bytes, err := command.CombinedOutput()
		if err != nil {
			return string(bytes), fmt.Errorf("executing command: %v", err)
		}
		return string(bytes), nil
	}
}

// CommandShell is a command to run in a shell and is made of a single string.
type CommandShell struct {
	Script string
}

// Run execute the command and returns its output and an error. Arguments:
// - pipe tells if we should print output of the command on the console.
func (c CommandShell) Run(build *build.Build, pipe bool, context *build.Context) (string, error) {
	shell, err := build.GetShell()
	if err != nil {
		return "", err
	}
	binary := shell[0]
	arguments := shell[1:]
	cmd, err := context.EvaluateString(c.Script)
	if err != nil {
		return "", err
	}
	arguments = append(arguments, cmd)
	command := exec.Command(binary, arguments...)
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting current working directory: %v", err)
	}
	command.Dir = dir
	command.Env, err = context.EvaluateEnvironment(build)
	if err != nil {
		return "", fmt.Errorf("building environment: %v", err)
	}
	if pipe {
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		err = command.Run()
		if err != nil {
			return "", fmt.Errorf("executing command: %v", err)
		}
		return "", nil
	} else {
		bytes, err := command.CombinedOutput()
		if err != nil {
			return string(bytes), fmt.Errorf("executing command: %v", err)
		}
		return string(bytes), nil
	}
}

// NewCommands parses a step of the build file to build a command. Arguments:
// - build is a reference to the build.
// - object is the parsed step.
func NewCommands(build *build.Build, object interface{}) (*Commands, error) {
	if !util.IsMap(object) {
		m := map[string]interface{}{
			"default": object,
		}
		return NewCommands(build, m)
	}
	commands := make(map[string]CommandFunc)
	m, _ := util.ToMapStringInterface(object)
	for os, cmd := range m {
		if util.IsSlice(cmd) {
			c, err := util.ToSliceString(cmd)
			if err != nil {
				return nil, err
			}
			commands[os] = CommandList{Parts: c}
		} else if util.IsString(cmd) {
			s := cmd.(string)
			commands[os] = CommandShell{Script: s}
		} else {
			return nil, fmt.Errorf("command must a string or liste of strings")
		}
	}
	return &Commands{
		Build:    build,
		Commands: commands,
	}, nil
}
