// +build !windows

package builtin

import (
	"reflect"
	"testing"
)

func TestFindInPath(t *testing.T) {
	actual := findInPath("ls")
	if !reflect.DeepEqual(actual, []string{"/bin/ls"}) &&
		!reflect.DeepEqual(actual, []string{"/usr/bin/ls", "/bin/ls"}) {
		t.Errorf("findinpath test failed: bad ls locations: %v", actual)
	}
}
