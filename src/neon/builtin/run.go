package builtin

import (
	"neon/build"
	"os/exec"
	"strings"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "run",
		Func: Run,
		Help: `Run given command and return output.

Arguments:

- The command to run.
- The arguments of the command as strings.

Returns:

- The standard and error output of the command as a string.
- If the command fails, this will cause the script failure.

Examples:

    # zip files of foo directory in bar.zip file
    run("zip", "-r", "bar.zip", "foo")
    # returns: the trimed output of the command`,
	})
}

func Run(command string, params ...string) string {
	cmd := exec.Command(command, params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(err.Error())
	}
	return strings.TrimSpace(string(output))
}
