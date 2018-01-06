package builtin

import (
	"testing"
)

func TestList(t *testing.T) {
	result := List("string")
	Assert(len(result), 1, t)
	Assert(result[0], "string", t)
	result = List([]interface{}{"string"})
	Assert(len(result), 1, t)
	Assert(result[0], "string", t)
}
