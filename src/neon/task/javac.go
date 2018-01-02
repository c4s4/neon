package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"os"
	"os/exec"
	"path"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "javac",
		Func: Javac,
		Args: reflect.TypeOf(JavacArgs{}),
		Help: `Compile Java source files.

Arguments:

- javac: the glob for Java source files.
- source: directory for source files.
- exclude: glob for source files to exclude (optional).
- dest: destination directory for generated classes.
- cp: classpath for compilation.

Examples:

    # compile Java source files in src directory
    - javac:  '**/*.java'
      source: 'src'
      dest:   'build/classes'
    # compile Java source files in src directory with given classpath
    - javac:  '**/*.java'
      source: 'src'
      dest:   'build/classes'
      cp:     '#{classpath}'`,
	})
}

type JavacArgs struct {
	Javac   []string `file wrap`
	Source  string   `file`
	Exclude []string `optional file wrap`
	Dest    string   `file`
	Cp      string   `optional`
}

func Javac(context *build.Context, args interface{}) error {
	params := args.(JavacArgs)
	if !util.DirExists(params.Dest) {
		err := os.MkdirAll(params.Dest, util.DIR_FILE_MODE)
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
		opt = append(opt, path.Join(params.Source, s))
	}
	command := exec.Command("javac", opt...)
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current working directory: %v", err)
	}
	command.Dir = dir
	// FIXME
	//command.Env, err = context.EvaluateEnvironment(target.Build)
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
