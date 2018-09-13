package task

import (
	"fmt"
	"github.com/c4s4/neon/build"
	"github.com/c4s4/neon/util"
	"os"
	"os/exec"
	p "path"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "javac",
		Func: javac,
		Args: reflect.TypeOf(javacArgs{}),
		Help: `Compile Java source files.

Arguments:

- javac: glob of Java source files to compile (strings, file, wrap).
- source: directory of source files (string, file).
- exclude: glob of source files to exclude (strings, optional, file, wrap).
- dest: destination directory for generated classes (string, file).
- cp: classpath for compilation (string, optional).

Examples:

    # compile Java source files in src directory
    - javac:  '**/*.java'
      source: 'src'
      dest:   'build/classes'
    # compile Java source files in src directory with given classpath
    - javac:  '**/*.java'
      source: 'src'
      dest:   'build/classes'
      cp:     =classpath`,
	})
}

type javacArgs struct {
	Javac   []string `neon:"file,wrap"`
	Source  string   `neon:"file"`
	Exclude []string `neon:"optional,file,wrap"`
	Dest    string   `neon:"file"`
	Cp      string   `neon:"optional"`
}

func javac(context *build.Context, args interface{}) error {
	params := args.(javacArgs)
	if !util.DirExists(params.Dest) {
		err := os.MkdirAll(params.Dest, util.DirFileMode)
		if err != nil {
			return fmt.Errorf("making destination diectory: %v", err)
		}
	}
	sources, err := util.FindFiles(params.Source, params.Javac, params.Exclude, false)
	if err != nil {
		return fmt.Errorf("finding java source files: %v", err)
	}
	cp, err := context.EvaluateString(params.Cp)
	if err != nil {
		return fmt.Errorf("evaluating classpath: %v", err)
	}
	// run javac command
	context.Message("Compiling %d Java source file(s)", len(sources))
	opt := []string{"-d", params.Dest}
	if cp != "" {
		opt = append(opt, []string{"-cp", cp}...)
	}
	for _, s := range sources {
		opt = append(opt, p.Join(params.Source, s))
	}
	command := exec.Command("javac", opt...)
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current working directory: %v", err)
	}
	command.Dir = dir
	command.Env, err = context.EvaluateEnvironment()
	if err != nil {
		return fmt.Errorf("building environment: %v", err)
	}
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err = command.Run()
	if err != nil {
		return fmt.Errorf("compiling java source files: %v", err)
	}
	return nil
}
