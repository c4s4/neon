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

func ParseCommandLine() (string, bool, bool, bool, string, bool, bool, string, bool, []string) {
	file := flag.String("file", DEFAULT_BUILD_FILE, "Build file to run")
	help := flag.Bool("build", false, "Print build help")
	timeit := flag.Bool("time", false, "Print build duration")
	tasks := flag.Bool("tasks", false, "Print tasks list")
	task := flag.String("task", "", "Print help on given task")
	targs := flag.Bool("targets", false, "Print targets list")
	builtins := flag.Bool("builtins", false, "Print builtins list")
	builtin := flag.String("builtin", "", "Print help on given builtin")
	refs := flag.Bool("refs", false, "Print tasks and builtins reference")
	flag.Parse()
	targets := flag.Args()
	return *file, *help, *timeit, *tasks, *task, *targs, *builtins, *builtin, *refs, targets
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
	file, help, timeit, tasks, task, targs, builtins, builtin, refs, targets := ParseCommandLine()
	if tasks {
		_build.PrintTasks()
	} else if task != "" {
		_build.PrintHelpTask(task)
	} else if builtins {
		_build.PrintBuiltins()
	} else if builtin != "" {
		_build.PrintHelpBuiltin(builtin)
	} else if refs {
		_build.PrintReference()
	}
	path, err := FindBuildFile(file)
	if err != nil {
		util.PrintColor("%s %s", util.Red("ERROR"), err.Error())
		os.Exit(1)
	}
	build, err := _build.NewBuild(path)
	if err != nil {
		util.PrintColor("%s %s", util.Red("ERROR"), err.Error())
		os.Exit(2)
	}
	if targs {
		build.PrintTargets()
	} else if help {
		err = build.Init()
		if err == nil {
			err = build.Help()
		}
		if err != nil {
			util.PrintColor("%s %s", util.Red("ERROR"), err.Error())
			os.Exit(3)
		}
	} else {
		err = build.Init()
		if err == nil {
			err = build.Run(targets)
		}
		duration := time.Now().Sub(start)
		if timeit || duration.Seconds() > 10 {
			_build.Info("Build duration: %s", duration.String())
		}
		if err == nil {
			util.PrintColor("%s", util.Green("OK"))
		} else {
			util.PrintColor("%s %s", util.Red("ERROR"), err.Error())
			os.Exit(4)
		}
	}
}
