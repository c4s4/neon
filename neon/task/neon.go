package task

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/c4s4/neon/neon/build"
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
- properties: the properties to pass to the build (map, optional).

Examples:

    # run target 'hello' of build file 'bar/build.yml' with property name = 'world'
    - neon:    'bar/build.yml'
      targets: 'hello'
	  properties:
	    name: 'world'`,
	})
}

type neonArgs struct {
	Neon       string                      `neon:"file"`
	Targets    []string                    `neon:"optional,wrap"`
	Properties map[interface{}]interface{} `neon:"optional,properties"`
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
	if params.Properties != nil {
		for key, value := range params.Properties {
			name := fmt.Sprintf("%v", key)
			newContext.SetProperty(name, value)
		}
	}
	err = newBuild.Run(newContext, params.Targets)
	if err != nil {
		return fmt.Errorf("running build '%s': %v", params.Neon, err)
	}
	return nil
}
