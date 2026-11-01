package task

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/c4s4/neon/neon/build"
	"github.com/c4s4/neon/neon/util"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "$",
		Func: shell,
		Args: reflect.TypeOf(shellArgs{}),
		Help: `Execute a command and return output and value.

Arguments:

- $: command to run (string or list of strings, other scalars are converted
  to strings).
- +: options to pass on command line after command (strings, other scalars
  are converted to strings, optional).
- n=: write command output into named property. Values for n are: 1 for stdout,
  2 for stderr and 3 for stdout and stderr (string, optional).
- n>: write command output in named file. Values for n are: 1 for stdout,
  2 for stderr and 3 for stdout and stderr (string, optional).
- n>>: append command output to named file. Values for n are: 1 for stdout,
  2 for stderr and 3 for stdout and stderr (string, optional).
- nx: disable command output. Values for n are: 1 for stdout, 2 for stderr and
  3 for stdout and stderr (boolean, optional).
- <: send given text to standard input of the process (string, other scalars
  are converted to strings, optional).
- :: print command on terminal before running it (boolean, optional).
- env: environment variables to set running command (map of strings, other
  scalars are converted to strings, optional).
- cafe: prevent the system from going to sleep while the command is running
  (boolean, optional). Only supported on Linux (systemd-inhibit) and MacOS
  (caffeinate).

Examples:

    # execute ls command, put result in 'files' variable and disable output
    - $:  'ls -al'
      1=: 'files'
      1x: true
    # execute command as a list of strings and output on console
    - $: ['ls', '-al']
    # run go test with options and packages defined as lists in properties
    - $: ['go', 'test', =TEST_OPTS, =PACKAGES]
    # run pylint on all python files except those in venv
    - $: 'pylint'
      +: '=filter(find(".", "**/*.py"), "venv/**/*.py")'
    # run echo command passing environment variable
    - $: 'echo "Hello $NAME!"'
      env:
        NAME: 'John'
    # run a long command preventing the system from going to sleep
    - $: 'sleep 3600'
      cafe: true

Notes:

- Scalars (numbers, booleans and dates) in command, options, standard input,
  environment values and file names of redirections are converted to strings.
  Thus you can write ['sleep', 1] instead of ['sleep', '1']. Floats are
  written without trailing zeros, so quote values such as versions ('1.0'
  instead of 1.0). Booleans give 'true' or 'false', unquoted dates (such as
  2026-11-01) are kept as written.
- Lists in command or options are flattened, each element being a separate
  argument. Thus you can write ['go', 'test', =OPTIONS, =PACKAGES] where
  OPTIONS and PACKAGES are lists.
- Commands defined as a string run in the shell defined by shell field at the
  root of the build file (or 'sh -c' on Unix and 'cmd /c' on Windows by
  default).
- Defining a command as a list of strings is useful on Windows. Default shell on
  Windows is 'cmd' which can't properly manage arguments with spaces.
- Argument of a command defined as a list won't be expanded by shell. Thus
  $USER won't be expanded for instance.
- The cafe option prevents the system from going to sleep while the command
  is running. It uses systemd-inhibit on Linux and caffeinate on MacOS. On
  other operating systems, or if the wrapper command is not available in
  PATH, it is silently ignored.`,
	})
}

type shellArgs struct {
	Shell []string          `neon:"name=$,wrap,strings"`
	Args  []string          `neon:"name=+,expression,wrap,optional,strings"`
	Del1  bool              `neon:"name=1x,bool,optional"`
	Del2  bool              `neon:"name=2x,bool,optional"`
	Del3  bool              `neon:"name=3x,bool,optional"`
	Red1  string            `neon:"name=1>,file,optional,strings"`
	Red2  string            `neon:"name=2>,file,optional,strings"`
	Red3  string            `neon:"name=3>,file,optional,strings"`
	App1  string            `neon:"name=1>>,file,optional,strings"`
	App2  string            `neon:"name=2>>,file,optional,strings"`
	App3  string            `neon:"name=3>>,file,optional,strings"`
	Var1  string            `neon:"name=1=,optional"`
	Var2  string            `neon:"name=2=,optional"`
	Var3  string            `neon:"name=3=,optional"`
	In    string            `neon:"name=<,optional,strings"`
	Verb  bool              `neon:"name=:,bool,optional"`
	Env   map[string]string `neon:"name=env,optional,strings"`
	Cafe  bool              `neon:"name=cafe,optional"`
}

func shell(context *build.Context, args interface{}) error {
	params := args.(shellArgs)
	// reader from stdin
	var stdin io.Reader = os.Stdin
	// writers to stdout and stderr
	stdout := getStdout(params)
	stderr := getStderr(params)
	// files are opened once, even if both outputs are redirected to the same
	// file, so that they share the same file offset
	files := make(map[string]*os.File)
	defer func() {
		for _, file := range files {
			_ = file.Close()
		}
	}()
	openFile := func(filename string, flag int) (io.Writer, error) {
		if file, ok := files[filename]; ok {
			return file, nil
		}
		file, err := os.OpenFile(filename, flag, util.FileMode)
		if err != nil {
			return nil, err
		}
		files[filename] = file
		return file, nil
	}
	// redirect output in a file
	redirections := []struct {
		filename string
		flag     int
		stdout   bool
		stderr   bool
	}{
		{params.Red1, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, true, false},
		{params.Red2, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, false, true},
		{params.Red3, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, true, true},
		{params.App1, os.O_WRONLY | os.O_CREATE | os.O_APPEND, true, false},
		{params.App2, os.O_WRONLY | os.O_CREATE | os.O_APPEND, false, true},
		{params.App3, os.O_WRONLY | os.O_CREATE | os.O_APPEND, true, true},
	}
	for _, redirection := range redirections {
		if redirection.filename == "" {
			continue
		}
		file, err := openFile(redirection.filename, redirection.flag)
		if err != nil {
			return err
		}
		if redirection.stdout {
			stdout = append(stdout, file)
		}
		if redirection.stderr {
			stderr = append(stderr, file)
		}
	}
	// write output in properties, with writers safe for concurrent use as
	// stdout and stderr are copied by distinct goroutines
	properties := make(map[string]*syncBuilder)
	captures := []struct {
		property string
		stdout   bool
		stderr   bool
	}{
		{params.Var1, true, false},
		{params.Var2, false, true},
		{params.Var3, true, true},
	}
	for _, capture := range captures {
		if capture.property == "" {
			continue
		}
		builder, ok := properties[capture.property]
		if !ok {
			builder = &syncBuilder{}
			properties[capture.property] = builder
		}
		if capture.stdout {
			stdout = append(stdout, builder)
		}
		if capture.stderr {
			stderr = append(stderr, builder)
		}
	}
	// write in standard input
	if params.In != "" {
		stdin = strings.NewReader(params.In)
	}
	err := run(params.Shell, params.Args, multiWriter(stdout), multiWriter(stderr), stdin, context, params.Env, params.Verb, params.Cafe)
	for property, builder := range properties {
		context.SetProperty(property, strings.TrimSpace(builder.String()))
	}
	return err
}

// syncBuilder is a strings.Builder safe for concurrent writes
type syncBuilder struct {
	mutex   sync.Mutex
	builder strings.Builder
}

// Write writes given bytes in the builder
func (b *syncBuilder) Write(p []byte) (int, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.builder.Write(p)
}

// String returns the content of the builder
func (b *syncBuilder) String() string {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.builder.String()
}

// multiWriter wraps writers in a single writer
func multiWriter(writers []io.Writer) io.Writer {
	switch len(writers) {
	case 0:
		return io.Discard
	case 1:
		return writers[0]
	default:
		return io.MultiWriter(writers...)
	}
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

func run(command []string, args []string, stdout, stderr io.Writer, stdin io.Reader, context *build.Context, env map[string]string, verbose bool, cafe bool) error {
	if args != nil {
		command = append(command, args...)
	}
	if len(command) == 0 {
		return fmt.Errorf("empty command")
	}
	if cafe && (util.GOOS == "darwin" || util.GOOS == "linux") {
		wrapper, err := cafeCommand(util.GOOS)
		if err != nil {
			return err
		}
		if _, err := exec.LookPath(wrapper[0]); err == nil {
			if len(command) == 1 {
				command[0] = strings.Join(wrapper, " ") + " " + command[0]
			} else {
				command = append(wrapper, command...)
			}
		}
	}
	if len(command) < 2 {
		return runString(command[0], stdout, stderr, stdin, context, env, verbose)
	} else {
		return runList(command, stdout, stderr, stdin, context, env, verbose)
	}
}

// cafeCommand returns the command to prevent the system from going to sleep
// while a command is running:
// - os: the operating system (such as "linux" or "darwin")
// Return:
// - the command as a list of strings
// - an error if the operating system is not supported
func cafeCommand(os string) ([]string, error) {
	switch os {
	case "darwin":
		return []string{"caffeinate", "-i"}, nil
	case "linux":
		return []string{"systemd-inhibit", "--what=idle", "--who=NeON",
			"--why=build", "--mode=block"}, nil
	default:
		return nil, fmt.Errorf("cafe option is not supported on %s", os)
	}
}

func runList(cmd []string, stdout, stderr io.Writer, stdin io.Reader, context *build.Context, env map[string]string, verbose bool) error {
	if verbose {
		context.MessageArgs("Running command: %s", strings.Join(cmd, " "))
	}
	environment, err := context.EvaluateEnvironment()
	if err != nil {
		return fmt.Errorf("building environment: %v", err)
	}
	for key, value := range env {
		environment[key] = value
	}
	environ := make([]string, 0, len(environment))
	for key, value := range environment {
		environ = append(environ, key+"="+value)
	}
	executable := cmd[0]
	executablePath, err := lookPath(executable, environment["PATH"])
	if err != nil {
		return fmt.Errorf("command '%s' was not found in PATH: %v", executable, err)
	}
	arguments := cmd[1:]
	command := exec.Command(executablePath, arguments...)
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current working directory: %v", err)
	}
	command.Dir = dir
	command.Stdin = stdin
	command.Stdout = stdout
	command.Stderr = stderr
	command.Env = environ
	err = command.Run()
	if err != nil {
		return fmt.Errorf("executing command: %v", err)
	}
	return nil
}

func runString(cmd string, stdout, stderr io.Writer, stdin io.Reader, context *build.Context, env map[string]string, verbose bool) error {
	if verbose {
		context.MessageArgs("Running command: %s", cmd)
	}
	shell, err := context.Build.GetShell()
	if err != nil {
		return err
	}
	environment, err := context.EvaluateEnvironment()
	if err != nil {
		return fmt.Errorf("building environment: %v", err)
	}
	for key, value := range env {
		environment[key] = value
	}
	environ := make([]string, 0, len(environment))
	for key, value := range environment {
		environ = append(environ, key+"="+value)
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
	command.Stdin = stdin
	command.Stdout = stdout
	command.Stderr = stderr
	command.Env = environ
	err = command.Run()
	if err != nil {
		return fmt.Errorf("executing command: %v", err)
	}
	return nil
}

// lookPath searches executable in given PATH, without changing the PATH of
// the process (which would be unsafe with threads):
// - executable: the executable to search
// - path: the PATH to search in, process PATH if empty
// Return: the path of the executable and an error if not found
func lookPath(executable, path string) (string, error) {
	if path == "" || strings.ContainsRune(executable, os.PathSeparator) || strings.ContainsRune(executable, '/') {
		return exec.LookPath(executable)
	}
	for _, dir := range filepath.SplitList(path) {
		if dir == "" {
			dir = "."
		}
		// LookPath checks file directly when it contains a separator (and
		// tries executable extensions on Windows)
		candidate, err := exec.LookPath(filepath.Join(dir, executable))
		if err == nil {
			return candidate, nil
		}
	}
	return "", exec.ErrNotFound
}
