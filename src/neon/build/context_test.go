package build

import (
	"strings"
	"testing"
)

func TestGetEnvironmentSimple(t *testing.T) {
	build := &Build{
		Dir: "dir",
	}
	context := &Context{
		Environment: map[string]string{
			"FOO": "BAR",
		},
	}
	env, err := context.EvaluateEnvironment(build)
	if err != nil {
		t.Errorf("Error getting environment: %v", err)
	}
	for _, line := range env {
		if line == "FOO=BAR" {
			return
		}
	}
	t.Error("Env FOO=BAR not found")
}

func TestGetEnvironmentComplex(t *testing.T) {
	build := &Build{
		Dir: "dir",
	}
	context := &Context{
		Environment: map[string]string{
			"FOO": "BAR:${HOME}",
		},
	}
	env, err := context.EvaluateEnvironment(build)
	if err != nil {
		t.Errorf("Error getting environment: %v", err)
	}
	var foo string
	for _, line := range env {
		if strings.HasPrefix(line, "HOME=") {
			foo = "FOO=BAR:" + line[5:]
		}
	}
	if foo == "" {
		return
	}
	for _, line := range env {
		if line == foo {
			return
		}
	}
	t.Error("Environment variable FOO not set correctly")
}
