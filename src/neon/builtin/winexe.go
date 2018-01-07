package builtin

import (
	"neon/build"
	"neon/util"
	"strings"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "winexe",
		Func: Winexe,
		Help: `Add '.exe' or '.bat' extensions depending on platform:
- command will stay command on Unix and will become command.exe on Windows.
- script.sh will stay script.sh on Unix and will become script.bat on Windows.
It will also replace / with \ in the executable path.

Arguments:

- The command to process.

Returns:

- Command adapted to host system.

Examples:

    // run command foo on unix and windows
    run(winexe("bin/foo"))
    // will run bin/foo on unix and bin\foo.exe on windows
    // run script script.sh unix and windows
    run(winexe("script.sh"))
    // will run script.sh on unix and script.bat on windows`,
	})
}

func Winexe(command string) string {
	if util.Windows() {
		if strings.HasSuffix(command, ".sh") {
			command = command[:len(command)-3] + ".bat"
		} else {
			command = command + ".exe"
		}
		return strings.Replace(command, "/", "\\", -1)
	} else {
		return command
	}
}
