package builtin

import (
	"testing"
)

func TestSortVersions(t *testing.T) {
	versions := []string{"1.10", "1.1", "1.2"}
	sortVersions(versions)
	expected := []string{"1.1", "1.2", "1.10"}
	Assert(versions, expected, t)
}
