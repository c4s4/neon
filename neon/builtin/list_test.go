package builtin

import (
	"testing"
)

func TestList(t *testing.T) {
	result := list("string")
	Assert(len(result), 1, t)
	Assert(result[0], "string", t)
	result = list([]interface{}{"string"})
	Assert(len(result), 1, t)
	Assert(result[0], "string", t)
}
