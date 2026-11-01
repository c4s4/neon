package task

import (
	"fmt"
	"os"
	"reflect"

	"github.com/c4s4/neon/neon/build"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "dotenv",
		Func: dotenv,
		Args: reflect.TypeOf(dotenvArgs{}),
		Help: `Load given dotenv file in environment.

Arguments:

- dotenv: name of dotenv file to load (string, file).

Examples:

    # load ".env" file in environment
    - dotenv: '.env'`,
	})
}

type dotenvArgs struct {
	Dotenv string `neon:"file"`
}

func dotenv(context *build.Context, args interface{}) error {
	params := args.(dotenvArgs)
	context.MessageArgs("Loading environment in dotenv file %s", params.Dotenv)
	err := LoadEnv(params.Dotenv)
	if err != nil {
		return fmt.Errorf("loading dotenv file: %v", err)
	}
	return nil
}

// LoadEnv loads environment in given file
func LoadEnv(filename string) error {
	environment, err := build.LoadDotEnv(filename, "")
	if err != nil {
		return err
	}
	for name, value := range environment {
		if err := os.Setenv(name, value); err != nil {
			return fmt.Errorf("setting environment variable '%s': %v", name, err)
		}
	}
	return nil
}
