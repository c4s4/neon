package task

import (
	"fmt"
	"io"
	"neon/build"
	"neon/util"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "$",
		Func: shell,
		Args: reflect.TypeOf(shellArgs{}),
		Help: `Execute a command and return output and value.

Arguments:

- $: command to run (string or list of strings).
- +: options to pass on command line after command (strings, optional).
- n=: write command output into named property. Values for n are: 1 for stdout,
  2 for stderr and 3 for stdout and stderr.
- n>: write command output in named file. Values for n are: 1 for stdout,
  2 for stderr and 3 for stdout and stderr.
- n>>: append command output to named file. Values for n are: 1 for stdout,
  2 for stderr and 3 for stdout and stderr.
- nx: disable command output. Values for n are: 1 for stdout, 2 for stderr and
  3 for stdout and stderr.
- <: send given text to standard input of the process.
- :: print command on terminal before running it.

Examples:

    # execute ls command and get result in 'files' variable
    - $:  'ls -al'
      1=: 'files'
    # execute command as a list of strings and output on console
    - $: ['ls', '-al']
    # run pylint on all python files except those in venv
    - $: 'pylint'
      +: '=filter(find(".", "**/*.py"), "venv/**/*.py")'

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

type shellArgs struct {
	Shell []string `neon:"name=$,wrap"`
	Args  []string `neon:"name=+,expression,optional"`
	Del1  bool     `neon:"name=1x,bool,optional"`
	Del2  bool     `neon:"name=2x,bool,optional"`
	Del3  bool     `neon:"name=3x,bool,optional"`
	Red1  string   `neon:"name=1>,file,optional"`
	Red2  string   `neon:"name=2>,file,optional"`
	Red3  string   `neon:"name=3>,file,optional"`
	App1  string   `neon:"name=1>>,file,optional"`
	App2  string   `neon:"name=2>>,file,optional"`
	App3  string   `neon:"name=3>>,file,optional"`
	Var1  string   `neon:"name=1=,optional"`
	Var2  string   `neon:"name=2=,optional"`
	Var3  string   `neon:"name=3=,optional"`
	In    string   `neon:"name=<,optional"`
	Verb  bool     `neon:"name=:,bool,optional"`
}

func shell(context *build.Context, args interface{}) error {
	params := args.(shellArgs)
	// reader from stdin
	var stdin io.Reader = os.Stdin
	// writers to stdout and stderr
	stdout := getStdout(params)
	stderr := getStderr(params)
	// string builder to redirect in a property
	builder := &strings.Builder{}
	property := ""
	// redirect stdout in a file
	if params.Red1+params.Red3 != "" {
		filename := params.Red1 + params.Red3
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		stdout = append(stdout, file)
	}
	// redirect stderr in a file
	if params.Red2+params.Red3 != "" {
		filename := params.Red2 + params.Red3
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		stderr = append(stderr, file)
	}
	// append stdout to a file
	if params.App1+params.App3 != "" {
		filename := params.App1 + params.App3
		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, util.FileMode)
		if err != nil {
			return err
		}
		defer file.Close()
		stdout = append(stdout, file)
	}
	// append stderr to a file
	if params.App2+params.App3 != "" {
		filename := params.App2 + params.App3
		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, util.FileMode)
		if err != nil {
			return err
		}
		defer file.Close()
		stderr = append(stderr, file)
	}
	// write stdout in a property
	if params.Var1+params.Var3 != "" {
		stdout = append(stdout, builder)
		property = params.Var1 + params.Var3
	}
	// write stderr in a property
	if params.Var2+params.Var3 != "" {
		stderr = append(stderr, builder)
		property = params.Var2 + params.Var3
	}
	// write in standart input
	if params.In != "" {
		stdin = strings.NewReader(params.In)
	}
	// put writers in a multi writer
	multiStdout := io.MultiWriter(stdout...)
	multiStderr := io.MultiWriter(stderr...)
	err := run(params.Shell, params.Args, multiStdout, multiStderr, stdin, context, params.Verb)
	if property != "" {
		context.SetProperty(property, strings.TrimSpace(builder.String()))
	}
	if err != nil {
		return err
	}
	return nil
}

func getStdout(params shellArgs) []io.Writer {
	// output disabled on stdout
	if params.Del1 || params.Del3 {
		return []io.Writer{}
	}
	return []io.Writer{os.Stdout}
}

func getStderr(params shellArgs) []io.Writer {
	// output disabled on stderr
	if params.Del2 || params.Del3 {
		return []io.Writer{}
	}
	return []io.Writer{os.Stderr}
}

func run(command []string, args []string, stdout, stderr io.Writer, stdin io.Reader, context *build.Context, verbose bool) error {
	if args != nil {
		command = append(command, args...)
	}
	if len(command) == 0 {
		return fmt.Errorf("empty command")
	} else if len(command) < 2 {
		return runString(command[0], stdout, stderr, stdin, context, verbose)
	} else {
		return runList(command, stdout, stderr, stdin, context, verbose)
	}
}

func runList(cmd []string, stdout, stderr io.Writer, stdin io.Reader, context *build.Context, verbose bool) error {
	if verbose {
		context.Message("Running command: %s", strings.Join(cmd, " "))
	}
	command := exec.Command(cmd[0], cmd[1:]...)
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current working directory: %v", err)
	}
	command.Dir = dir
	command.Env, err = context.EvaluateEnvironment()
	if err != nil {
		return fmt.Errorf("building environment: %v", err)
	}
	command.Stdin = stdin
	command.Stdout = stdout
	command.Stderr = stderr
	err = command.Run()
	if err != nil {
		return fmt.Errorf("executing command: %v", err)
	}
	return nil
}

func runString(cmd string, stdout, stderr io.Writer, stdin io.Reader, context *build.Context, verbose bool) error {
	if verbose {
		context.Message("Running command: %s", cmd)
	}
	shell, err := context.Build.GetShell()
	if err != nil {
		return err
	}
	binary := shell[0]
	arguments := shell[1:]
	arguments = append(arguments, cmd)
	command := exec.Command(binary, arguments...)
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current working directory: %v", err)
	}
	command.Dir = dir
	command.Env, err = context.EvaluateEnvironment()
	if err != nil {
		return fmt.Errorf("building environment: %v", err)
	}
	command.Stdin = stdin
	command.Stdout = stdout
	command.Stderr = stderr
	err = command.Run()
	if err != nil {
		return fmt.Errorf("executing command: %v", err)
	}
	return nil
}
