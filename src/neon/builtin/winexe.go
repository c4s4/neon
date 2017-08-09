package builtin

import (
	"neon/build"
	"neon/util"
	"strings"
)

func init() {
	build.BuiltinMap["winexe"] = build.BuiltinDescriptor{
		Function: Winexe,
		Help: `Add '.exe' or '.bat' extensions depending on platform:
- command will stay command on Unix and will become command.exe on Windows.
- script.sh will stay script.sh on Unix and will become script.bat on Windows.

Arguments:

- The command to process.

Returns:

- Command adapted to host system.

Examples:

    // run command foo on unix and windows
    run(winexe("foo"))
    // will run foo on unix and foo.exe on windows
    // run script script.sh unix and windows
    run(winexe("script.sh"))
    // will run script.sh on unix and script.bat on windows`,
	}
}

func Winexe(command string) string {
	if strings.HasSuffix(command, ".sh") {
		if util.Windows() {
			return command[:len(command)-3] + ".bat"
		} else {
			return command
		}
	} else {
		if util.Windows() {
			return command + ".exe"
		} else {
			return command
		}
	}
}
