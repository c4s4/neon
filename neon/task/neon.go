package task

import (
	"fmt"
	"github.com/c4s4/neon/neon/build"
	"os"
	"path/filepath"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "neon",
		Func: neon,
		Args: reflect.TypeOf(neonArgs{}),
		Help: `Run a NeON build.

Arguments:

- neon: the build file to run (string).
- targets: the target(s) to run (strings, wrap, optional).

Examples:

    # run target 'foo' of build file 'bar/build.yml'
    - neon:    'bar/build.yml'
      targets: 'foo'`,
	})
}

type neonArgs struct {
	Neon    string   `neon:"file"`
	Targets []string `neon:"optional,wrap"`
}

func neon(context *build.Context, args interface{}) error {
	params := args.(neonArgs)
	// FIXME: path relative to build directory
	path, err := filepath.Abs(params.Neon)
	if err != nil {
		return fmt.Errorf("getting build file path: %v", err)
	}
	base := filepath.Dir(path)
	newBuild, err := build.NewBuild(path, base, context.Build.Repository, context.Build.Template)
	if err != nil {
		return fmt.Errorf("instantiating build: %v", err)
	}
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(dir)
	}()
	if err := os.Chdir(newBuild.Dir); err != nil {
		return err
	}
	newContext := build.NewContext(newBuild)
	err = newContext.Init()
	if err != nil {
		return fmt.Errorf("initializing build context: %v", err)
	}
	err = newBuild.Run(newContext, params.Targets)
	if err != nil {
		return fmt.Errorf("running build '%s': %v", params.Neon, err)
	}
	return nil
}
