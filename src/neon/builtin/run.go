package builtin

import (
	"neon/build"
	"os/exec"
)

func init() {
	build.BuiltinMap["run"] = build.BuiltinDescriptor{
		Function: Run,
		Help:     "Run given command and return output",
	}
}

func Run(command string, params ...string) string {
	cmd := exec.Command(command, params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(err.Error())
	}
	return string(output)
}
