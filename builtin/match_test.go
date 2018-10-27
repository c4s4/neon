package builtin

import (
	"testing"
)

func TestMatch(t *testing.T) {
	if !match(`n..n`, "neon") {
		t.Errorf("Match error")
	}
	if match(`n..n`, "zion") {
		t.Errorf("Match error")
	}
}
