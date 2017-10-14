package build

import (
	"strings"
	"testing"
)

func TestGetEnvironmentSimple(t *testing.T) {
	context := &VM{
		Environment: map[string]string{
			"FOO": "BAR",
		},
		Build: &Build{
			Dir: "dir",
		},
	}
	env, err := context.EvaluateEnvironment()
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
	context := &VM{
		Environment: map[string]string{
			"FOO": "BAR:${HOME}",
		},
		Build: &Build{
			Dir: "dir",
		},
	}
	env, err := context.EvaluateEnvironment()
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
