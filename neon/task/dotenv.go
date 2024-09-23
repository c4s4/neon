package task

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/c4s4/neon/neon/build"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "dotenv",
		Func: dotenv,
		Args: reflect.TypeOf(dotenvArgs{}),
		Help: `Load given dotenv file in environment.

Arguments:

- dotenv: name of dotenv file to load.

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
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		bytes, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		line := strings.TrimSpace(string(bytes))
		if line[0] == '#' {
			continue
		}
		index := strings.Index(line, "=")
		if index < 0 {
			return fmt.Errorf("bad environment line: '%s'", line)
		}
		name := strings.TrimSpace(line[:index])
		value := strings.TrimSpace(line[index+1:])
		os.Setenv(name, value)
	}
	return nil
}
