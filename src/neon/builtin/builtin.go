package builtin

import (
	"fmt"
	"github.com/mattn/anko/vm"
	zglob "github.com/mattn/go-zglob"
	"os"
	"sort"
	"time"
)

var Builtins = map[string]interface{}{
	"find":  Find,
	"throw": Throw,
	"now":   Now,
}

func AddBuiltins(vm *vm.Env) {
	for name, function := range Builtins {
		vm.Define(name, function)
	}
}

func Find(dir string, patterns ...string) []string {
	oldDir, err := os.Getwd()
	if err != nil {
		return nil
	}
	err = os.Chdir(dir)
	if err != nil {
		return nil
	}
	var files []string
	for _, pattern := range patterns {
		f, _ := zglob.Glob(pattern)
		for _, e := range f {
			files = append(files, e)
		}
	}
	sort.Strings(files)
	os.Chdir(oldDir)
	return files
}

func Throw(message string) error {
	return fmt.Errorf(message)
}

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
