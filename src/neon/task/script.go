package task

import (
	"fmt"
	"neon/build"
	"neon/util"
)

func init() {
	build.TaskMap["script"] = build.TaskDescriptor{
		Constructor: Script,
		Help: `Run an Anko script.

Arguments:

- script: the source of the script to run.

Examples:

    # build a classpath with all jar files in lib directory
    - script: |
        strings = import("strings")
        jars = find("lib", "*.jar")
        classpath = strings.Join(jars, ":")

Notes:

- The scripting language is Anko, which is a scriptable Go. For more information
  please refer to Anko site at http://github.com/mattn/anko. Thanks Mattn!
- Buitlin functions are functions you can access in scripts. To list them, you
  cas type 'neon -builtins', to get help on a given one, you may type for instance
  'neon -builtin find'.
- Properties can be accessed and set in scripts. Variables you define in scripts
  are readable as properties. In other words, scripts and properties share the
  same context.`,
	}
}

func Script(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"script"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	source, err := args.GetString("script")
	if err != nil {
		return nil, fmt.Errorf("parsing script task: %v", err)
	}
	return func() error {
		_, err := target.Build.Context.EvaluateExpression(source)
		if err != nil {
			return fmt.Errorf("evaluating script: %v", err)
		}
		return nil
	}, nil
}
