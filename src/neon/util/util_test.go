package util

import "testing"

// Make an assertion for testing purpose
func Assert(actual, expected string, t *testing.T) {
	if actual != expected {
		t.Errorf("actual \"%s\" != expected \"%s\"", actual, expected)
	}
}
