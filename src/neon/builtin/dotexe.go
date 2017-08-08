package builtin

import (
	"neon/build"
	"neon/util"
)

func init() {
	build.BuiltinMap["dotexe"] = build.BuiltinDescriptor{
		Function: Dotexe,
		Help: `Add '.exe' extension or not depending on platform.

Arguments:

- The command to process.

Returns:

- Command with '.exe' added (on Windows) or not (on Unix).

Examples:

    // run command foo on windows and linux
    run(dotexe("foo"))
    // runs 'foo.exe' on windows and 'foo' on unix`,
	}
}

func Dotexe(command string) string {
	if util.Windows() {
		return command + ".exe"
	} else {
		return command
	}
}
