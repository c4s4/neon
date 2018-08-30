package build

import (
	"testing"
	"reflect"
	"github.com/fatih/color"
)

func TestParseAttributes(t *testing.T) {
	for key, value := range Attributes {
		attributes, err := ParseAttributes([]string{key})
		if err != nil {
			t.Errorf("error parsing attribute: %v", err)
		}
		if attributes[0] != value {
			t.Errorf("bad attribute parsing")
		}
	}
}

func TestParseAttributesError(t *testing.T) {
	_, err := ParseAttributes([]string{"bad_attribute"})
	if err == nil {
		t.Errorf("should have failed to parse attribute")
	}
}

func TestParseTheme(t *testing.T) {
	colors := &Colors{
		Title: []string{"FgYellow"},
		Ok:    []string{"FgGreen", "Bold"},
		Error: []string{"FgRed", "Bold"},
	}
	expected := &Theme{
		Title: []color.Attribute{color.FgYellow},
		Ok:    []color.Attribute{color.FgGreen, color.Bold},
		Error: []color.Attribute{color.FgRed, color.Bold},
	}
	theme, err := ParseTheme(colors)
	if err != nil {
		t.Errorf("error parsing theme: %v", err)
	}
	if !reflect.DeepEqual(theme, expected) {
		t.Errorf("error parsing theme")
	}
}
