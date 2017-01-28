package main

import (
	"flag"
	"fmt"
	"neon/build"
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

func ParseCommandLine() (string, bool, bool, bool, string, bool, bool, string, []string) {
	file := flag.String("file", DEFAULT_BUILD_FILE, "Build file to run")
	help := flag.Bool("build", false, "Print build help")
	verbose := flag.Bool("verbose", false, "Verbose build output")
	tasks := flag.Bool("tasks", false, "Print tasks list")
	task := flag.String("task", "", "Print help on given task")
	targs := flag.Bool("targets", false, "Print targets list")
	builtins := flag.Bool("builtins", false, "Print builtins list")
	builtin := flag.String("builtin", "", "Print help on given builtin")
	flag.Parse()
	targets := flag.Args()
	return *file, *help, *verbose, *tasks, *task, *targs, *builtins, *builtin, targets
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
	file, help, verbose, tasks, task, targs, builtins, builtin, targets := ParseCommandLine()
	path, err := FindBuildFile(file)
	if err != nil {
		util.PrintColor("%s %s", util.Red("ERROR"), err.Error())
		os.Exit(1)
	}
	build, err := build.NewBuild(path, verbose)
	if err != nil {
		util.PrintColor("%s %s", util.Red("ERROR"), err.Error())
		os.Exit(2)
	}
	if help {
		err = build.Help()
		if err != nil {
			util.PrintColor("%s %s", util.Red("ERROR"), err.Error())
			os.Exit(3)
		}
	} else if tasks {
		build.PrintTasks()
	} else if task != "" {
		build.PrintHelpTask(task)
	} else if builtins {
		build.PrintBuiltins()
	} else if builtin != "" {
		build.PrintHelpBuiltin(builtin)
	} else if targs {
		build.PrintTargets()
	} else {
		start := time.Now()
		err = build.Run(targets)
		duration := time.Now().Sub(start)
		if duration.Seconds() > 10 {
			build.Info("Build duration: %s", duration.String())
		}
		if err == nil {
			util.PrintColor("%s", util.Green("OK"))
		} else {
			util.PrintColor("%s %s", util.Red("ERROR"), err.Error())
			os.Exit(4)
		}
	}
}
