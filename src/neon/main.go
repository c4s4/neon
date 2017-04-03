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

func ParseCommandLine() (string, bool, string, bool, bool, string, bool, bool, string, bool, []string) {
	file := flag.String("file", DEFAULT_BUILD_FILE, "Build file to run")
	help := flag.Bool("build", false, "Print build help")
	props := flag.String("props", "", "Build properties")
	timeit := flag.Bool("time", false, "Print build duration")
	tasks := flag.Bool("tasks", false, "Print tasks list")
	task := flag.String("task", "", "Print help on given task")
	targs := flag.Bool("targets", false, "Print targets list")
	builtins := flag.Bool("builtins", false, "Print builtins list")
	builtin := flag.String("builtin", "", "Print help on given builtin")
	refs := flag.Bool("refs", false, "Print tasks and builtins reference")
	flag.Parse()
	targets := flag.Args()
	return *file, *help, *props, *timeit, *tasks, *task, *targs, *builtins, *builtin,
		*refs, targets
}

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

func main() {
	start := time.Now()
	file, help, props, timeit, tasks, task, targs, builtins, builtin, refs, targets := ParseCommandLine()
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
	}
	// options that do require we load build file
	path, err := FindBuildFile(file)
	PrintError(err, 1)
	build, err := _build.NewBuild(path)
	PrintError(err, 2)
	if props != "" {
		err = build.SetProperties(props)
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

func PrintError(err error, code int) {
	if err != nil {
		util.PrintColor("%s %s", util.Red("ERROR"), err.Error())
		os.Exit(code)
	}
}
