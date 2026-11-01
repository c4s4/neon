package builtin

import (
	guuid "github.com/google/uuid"

	"github.com/c4s4/neon/neon/build"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "uuid",
		Func: uuid,
		Help: `Generate a random (version 4) UUID.

Returns:

- Generated UUID as a string.

Examples:

    # generate a UUID
    uuid()
    # returns: "9063eafd-40b9-48d9-bf90-ce44e9207821"`,
	})
}

func uuid() string {
	UUID, err := guuid.NewRandom()
	if err != nil {
		panic("ERROR generating UUID: " + err.Error())
	}
	return UUID.String()
}
