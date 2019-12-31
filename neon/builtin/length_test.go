package builtin

import (
	"testing"
)

func TestLength(t *testing.T) {
	Assert(length("Hello World!"), 12, t)
	Assert(length("éà&`"), 4, t)
}
