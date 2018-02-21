package builtin

import (
	"regexp"
	"testing"
)

func TestNow(t *testing.T) {
	if !regexp.MustCompile(`\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\d`).MatchString(now()) {
		t.Errorf("Error builtin now")
	}
}
