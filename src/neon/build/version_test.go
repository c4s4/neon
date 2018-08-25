package build

import (
	"testing"
)

func TestNewVersion(t *testing.T) {
	expected := Version{
		String: "1.2.3",
		Fields: []int{1, 2, 3},
	}
	actual, err := NewVersion("1.2.3")
	if err != nil {
		t.Error("Error parsing version", err)
	}
	if expected.String != actual.String {
		t.Error("Bad String field")
	}
	if len(expected.Fields) != len(actual.Fields) {
		t.Error("Bad Fields length")
	}
	for i := 0; i < len(actual.Fields); i++ {
		if expected.Fields[i] != actual.Fields[i] {
			t.Error("Bad Field", i)
		}
	}
	_, err = NewVersion("")
	if err == nil {
		t.Error("Error parsing version")
	}
	_, err = NewVersion(" 1.2.3")
	if err == nil {
		t.Error("Error parsing version")
	}
	_, err = NewVersion("1.2.")
	if err == nil {
		t.Error("Error parsing version")
	}
	_, err = NewVersion(".1.2")
	if err == nil {
		t.Error("Error parsing version")
	}
	_, err = NewVersion("a")
	if err == nil {
		t.Error("Error parsing version")
	}
	actual, err = NewVersion("1.2.3-SNAPSHOT")
	if err != nil {
		t.Error("Error parsing version", err)
	}
	if expected.String != actual.String {
		t.Error("Bad String field")
	}
	if len(expected.Fields) != len(actual.Fields) {
		t.Error("Bad Fields length")
	}
	for i := 0; i < len(actual.Fields); i++ {
		if expected.Fields[i] != actual.Fields[i] {
			t.Error("Bad Field", i)
		}
	}
}

func TestCompare(t *testing.T) {
	v, _ := NewVersion("0")
	o, _ := NewVersion("0")
	if v.Compare(o) != 0 {
		t.Errorf("Version comparison error: %d", v.Compare(o))
	}
	v, _ = NewVersion("0.1.2")
	o, _ = NewVersion("0.1.2")
	if v.Compare(o) != 0 {
		t.Errorf("Version comparison error: %d", v.Compare(o))
	}
	v, _ = NewVersion("0")
	o, _ = NewVersion("0.0")
	if v.Compare(o) >= 0 {
		t.Errorf("Version comparison error: %d", v.Compare(o))
	}
	v, _ = NewVersion("0.0")
	o, _ = NewVersion("0")
	if v.Compare(o) <= 0 {
		t.Errorf("Version comparison error: %d", v.Compare(o))
	}
	v, _ = NewVersion("0")
	o, _ = NewVersion("1")
	if v.Compare(o) >= 0 {
		t.Errorf("Version comparison error: %d", v.Compare(o))
	}
	v, _ = NewVersion("0.1.2")
	o, _ = NewVersion("0.1.3")
	if v.Compare(o) >= 0 {
		t.Errorf("Version comparison error: %d", v.Compare(o))
	}
}
