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
        name: 'world'

Notes:

- Relative build file paths are relative to the directory of the calling
  build (or current directory for templates).
- Properties passed to the build override those defined in the build file
  and are used to evaluate other properties of the build.
- Recursive calls of a build file are detected and raise an error.
- The called build uses its own repository, unless one was set with -repo
  option on command line.
- The called build runs in its own directory, thus it can't run in threads.`,
	})
}

type neonArgs struct {
	Neon       string                      `neon:"file"`
	Targets    []string                    `neon:"optional,wrap"`
	Properties map[interface{}]interface{} `neon:"optional,properties"`
}

func neon(context *build.Context, args interface{}) error {
	params := args.(neonArgs)
	// relative paths are relative to the build directory
	file := params.Neon
	if !filepath.IsAbs(file) {
		dir := context.Build.Dir
		if context.Build.Template {
			dir = context.Build.Here
		}
		file = filepath.Join(dir, file)
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return fmt.Errorf("getting build file path: %v", err)
	}
	// detect recursive calls
	for _, call := range callChain(context) {
		if call == path {
			return fmt.Errorf("recursive call of build file '%s'", params.Neon)
		}
	}
	base := filepath.Dir(path)
	// called build uses its own repository, unless set on command line
	newBuild, err := build.NewBuild(path, base, build.CommandLineRepository, context.Build.Template)
	if err != nil {
		return fmt.Errorf("instantiating build: %v", err)
	}
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current directory: %v", err)
	}
	defer func() {
		_ = build.Chdir(dir)
	}()
	if err := build.Chdir(newBuild.Dir); err != nil {
		return err
	}
	newContext := build.NewContext(newBuild)
	newContext.Calls = append(callChain(context), path)
	properties := make(map[string]interface{})
	for key, value := range params.Properties {
		properties[fmt.Sprintf("%v", key)] = value
	}
	err = newContext.InitWithProperties(properties)
	if err != nil {
		return fmt.Errorf("initializing build context: %v", err)
	}
	err = newBuild.Run(newContext, params.Targets)
	if err != nil {
		return fmt.Errorf("running build '%s': %v", params.Neon, err)
	}
	return nil
}

// callChain returns the list of build files being run, including the
// current one
func callChain(context *build.Context) []string {
	calls := append([]string{}, context.Calls...)
	if len(calls) == 0 && context.Build.FileDir != "" {
		calls = append(calls, filepath.Join(context.Build.FileDir, context.Build.File))
	}
	return calls
}
