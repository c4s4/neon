package main

import (
	"flag"
	"fmt"
	_build "neon/build"
	_ "neon/builtin"
	_ "neon/task"
	"neon/util"
	"os"
	"path/filepath"
	"time"
)

const (
	DEFAULT_BUILD_FILE = "build.yml"
)

var VERSION string

// Parse command line and return parsed options
func ParseCommandLine() (string, bool, bool, string, bool, bool, string, bool, bool, string, bool, string, bool, []string) {
	file := flag.String("file", DEFAULT_BUILD_FILE, "Build file to run")
	info := flag.Bool("info", false, "Print build information")
	version := flag.Bool("version", false, "Print neon version")
	props := flag.String("props", "", "Build properties")
	timeit := flag.Bool("time", false, "Print build duration")
	tasks := flag.Bool("tasks", false, "Print tasks list")
	task := flag.String("task", "", "Print help on given task")
	targs := flag.Bool("targets", false, "Print targets list")
	builtins := flag.Bool("builtins", false, "Print builtins list")
	builtin := flag.String("builtin", "", "Print help on given builtin")
	refs := flag.Bool("refs", false, "Print tasks and builtins reference")
	install := flag.String("install", "", "Install given plugin")
	grey := flag.Bool("grey", false, "Print on terminal without colors")
	flag.Parse()
	targets := flag.Args()
	return *file, *info, *version, *props, *timeit, *tasks, *task, *targs, *builtins,
		*builtin, *refs, *install, *grey, targets
}

// Find build file and return its path
func FindBuildFile(name string) (string, error) {
	absolute, err := filepath.Abs(name)
	if err != nil {
		return "", fmt.Errorf("getting build file path: %v", err)
	}
	file := filepath.Base(absolute)
	dir := filepath.Dir(absolute)
	for {
		path := filepath.Join(dir, file)
		if util.FileExists(path) {
			return path, nil
		} else {
			parent := filepath.Dir(dir)
			if parent == dir {
				return "", fmt.Errorf("build file not found")
			}
			dir = parent
		}
	}
}

// Program entry point
func main() {
	start := time.Now()
	file, info, version, props, timeit, tasks, task, targs, builtins, builtin, refs, install, grey, targets := ParseCommandLine()
	// options that do not require we load build file
	_build.Grey = grey
	if tasks {
		_build.PrintTasks()
		return
	} else if task != "" {
		_build.PrintHelpTask(task)
		return
	} else if builtins {
		_build.PrintBuiltins()
		return
	} else if builtin != "" {
		_build.PrintHelpBuiltin(builtin)
		return
	} else if refs {
		_build.PrintReference()
		return
	} else if version {
		_build.Message(VERSION)
		return
	}
	// options that do require we load build file
	path, err := FindBuildFile(file)
	PrintError(err, 1)
	build, err := _build.NewBuild(path)
	PrintError(err, 2)
	if props != "" {
		err = build.SetCommandLineProperties(props)
		PrintError(err, 3)
	}
	if install != "" {
		err = build.Install(install)
		PrintError(err, 6)
		return
	} else if targs {
		build.PrintTargets()
		return
	} else if info {
		context, err := _build.NewContext(build)
		PrintError(err, 4)
		err = build.Info(context)
		PrintError(err, 4)
		return
	} else {
		context, err := _build.NewContext(build)
		PrintError(err, 5)
		err = build.Run(context, targets)
		PrintError(err, 5)
		duration := time.Now().Sub(start)
		if timeit || duration.Seconds() > 10 {
			_build.Message("Build duration: %s", duration.String())
		}
		PrintError(err, 5)
		_build.PrintOk()
		return
	}
}

// Print an error and exit if any
func PrintError(err error, code int) {
	if err != nil {
		_build.PrintError(err.Error())
		os.Exit(code)
	}
}
