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

// Parse command line and return parsed options
func ParseCommandLine() (string, bool, bool, string, bool, bool, string, bool, bool, string, bool, string, []string) {
	file := flag.String("file", DEFAULT_BUILD_FILE, "Build file to run")
	help := flag.Bool("build", false, "Print build help")
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
	flag.Parse()
	targets := flag.Args()
	return *file, *help, *version, *props, *timeit, *tasks, *task, *targs, *builtins,
		*builtin, *refs, *install, targets
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
	file, help, version, props, timeit, tasks, task, targs, builtins, builtin, refs, install, targets := ParseCommandLine()
	// options that do not require we load build file
	if tasks {
		_build.PrintTasks()
		os.Exit(0)
	} else if task != "" {
		_build.PrintHelpTask(task)
		os.Exit(0)
	} else if builtins {
		_build.PrintBuiltins()
		os.Exit(0)
	} else if builtin != "" {
		_build.PrintHelpBuiltin(builtin)
		os.Exit(0)
	} else if refs {
		_build.PrintReference()
		os.Exit(0)
	} else if version {
		fmt.Println(VERSION)
		os.Exit(0)
	}
 	// options that do require we load build file
	path, err := FindBuildFile(file)
	PrintError(err, 1)
	build, err := _build.NewBuild(path)
	if build != nil && install != "" {
		err = build.Install(install)
		if err == nil {
			util.PrintColor("%s", util.Green("OK"))
			os.Exit(0)
		} else {
			PrintError(err, 6)
		}
	}
	PrintError(err, 2)
	if props != "" {
		err = build.SetCommandLineProperties(props)
		PrintError(err, 3)
	}
	err = build.Init()
	PrintError(err, 4)
	if targs {
		build.PrintTargets()
	} else if help {
		err = build.Help()
	} else {
		err = build.Run(targets)
		duration := time.Now().Sub(start)
		if timeit || duration.Seconds() > 10 {
			_build.Info("Build duration: %s", duration.String())
		}
		if err == nil {
			util.PrintColor("%s", util.Green("OK"))
		} else {
			PrintError(err, 5)
		}
	}
}

// Print an error and exit if any
func PrintError(err error, code int) {
	if err != nil {
		util.PrintColor("%s %s", util.Red("ERROR"), err.Error())
		os.Exit(code)
	}
}
