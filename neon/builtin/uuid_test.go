package builtin

import (
	"testing"
)

func TestUuid(t *testing.T) {
	AssertRegexp(uuid(), `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, t)
}
