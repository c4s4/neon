package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func init() {
	build.TaskMap["$"] = build.TaskDescriptor{
		Constructor: Shell,
		Help: `Execute a command and return output and value.

Arguments:

- $: command to run as a string or a list of strings. You can also provide a
  map of commands per operating system ("default" defines command to run on
  operating systems that are not in the map).
- =: name of the variable to store trimed output into (optional, output to
  console if not set).

Examples:

    # execute ls command and get result in 'files' variable
    - $: 'ls'
      =: 'files'
    # execute dir command on windows and ls on other OS
    - $:
    	windows: 'dir'
    	default: 'ls'
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
	}
}

// Commands lists commands by operating system
type Commands struct {
	Build    *build.Build
	Commands map[string]Command
}

// GetCommand return a command depending on current operating system and
// default command
func (c Commands) GetCommand() (Command, error) {
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
func (c Commands) Run(pipe bool) (string, error) {
	command, err := c.GetCommand()
	if err != nil {
		return "", err
	}
	output, err := command.Run(c.Build, pipe)
	output = util.RemoveBlankLines(output)
	output = strings.TrimSuffix(output, "\n")
	return output, err
}

// Command is the interface for a command.
type Command interface {
	Run(build *build.Build, pipe bool) (string, error)
}

// CommandList is a command as a list with executable and arguments.
type CommandList struct {
	Parts []string
}

// Run execute the command and returns its output and an error. Arguments:
// - pipe tells if we should print output of the command on the console.
func (c CommandList) Run(build *build.Build, pipe bool) (string, error) {
	parts := make([]string, len(c.Parts))
	var err error
	for i := 0; i < len(c.Parts); i++ {
		parts[i], err = build.Context.EvaluateString(c.Parts[i])
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
	command.Env, err = build.Context.EvaluateEnvironment()
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
func (c CommandShell) Run(build *build.Build, pipe bool) (string, error) {
	shell, err := build.GetShell()
	if err != nil {
		return "", err
	}
	binary := shell[0]
	arguments := shell[1:]
	cmd, err := build.Context.EvaluateString(c.Script)
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
	command.Env, err = build.Context.EvaluateEnvironment()
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
	commands := make(map[string]Command)
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

// Shell is the function to build a shell task.
// Arguments:
// - target in which will run the task.
// - args of the task.
// Returns the task and an error if any.
func Shell(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"$", "="}
	if err := CheckFields(args, fields, fields[:1]); err != nil {
		return nil, err
	}
	commands, err := NewCommands(target.Build, args["$"])
	if err != nil {
		return nil, err
	}
	var variable string
	var ok bool
	if args.HasField("=") {
		variable, ok = args["="].(string)
		if !ok {
			return nil, fmt.Errorf("argument = of task $ must be a string")
		}
	}
	return func() error {
		_variable, _err := target.Build.Context.EvaluateString(variable)
		if _err != nil {
			return fmt.Errorf("processing output argument: %v", _err)
		}
		_output, _err := commands.Run(_variable == "")
		if _err != nil {
			if _output != "" {
				build.Message(_output)
			}
			return _err
		}
		if _variable != "" {
			target.Build.Context.SetProperty(_variable, strings.TrimSpace(string(_output)))
		}
		return nil
	}, nil
}
