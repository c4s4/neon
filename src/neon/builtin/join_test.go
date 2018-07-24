package builtin

import (
	"testing"
)

func TestJoin(t *testing.T) {
	if join([]string{"foo", "bar"}, " ") != "foo bar" {
		t.Errorf("Error builtin join")
	}
	if join([]interface{}{"foo", "bar"}, " ") != "foo bar" {
		t.Errorf("Error builtin join")
	}
}

func TestJoinPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	join(1, " ")
}

func TestJoinPanic2(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	join([]int{1}, " ")
}
