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
	// missing fields count as zeros
	v, _ = NewVersion("0")
	o, _ = NewVersion("0.0")
	if v.Compare(o) != 0 {
		t.Errorf("Version comparison error: %d", v.Compare(o))
	}
	v, _ = NewVersion("1.0")
	o, _ = NewVersion("1.0.1")
	if v.Compare(o) >= 0 {
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

func TestCompareSuffixes(t *testing.T) {
	ordered := []string{"1.0-SNAPSHOT", "1.0-alpha", "1.0-BETA", "1.0-RC-1", "1.0-RC-2", "1.0", "1.0.1-SNAPSHOT"}
	for i := 0; i < len(ordered)-1; i++ {
		v, err := NewVersion(ordered[i])
		if err != nil {
			t.Fatal(err)
		}
		o, err := NewVersion(ordered[i+1])
		if err != nil {
			t.Fatal(err)
		}
		if v.Compare(o) >= 0 || o.Compare(v) <= 0 {
			t.Errorf("%s should be lower than %s", ordered[i], ordered[i+1])
		}
	}
}

func TestSortVersions(t *testing.T) {
	versions := []string{"v1.16.1", "1.10", "", "0.1.0", "1.2-RC1", "1.2", "1.2-BETA-2", "v1.9.1", "draft"}
	SortVersions(versions)
	expected := []string{"", "draft", "0.1.0", "1.2-BETA-2", "1.2-RC1", "1.2", "v1.9.1", "1.10", "v1.16.1"}
	Assert(versions, expected, t)
}

func TestVersionWithLeadingV(t *testing.T) {
	if !IsReleaseVersion("v1.2.3") || !IsReleaseVersion("V1.2.3-RC1") {
		t.Errorf("versions with leading v should be release versions")
	}
	if IsReleaseVersion("15fb03c") || IsReleaseVersion("UNKNOWN") || IsReleaseVersion("vv1.2") {
		t.Errorf("development versions should not be release versions")
	}
	v, err := NewVersion("v1.2.3")
	if err != nil || v.String != "1.2.3" {
		t.Fatalf("bad version: %v (%v)", v, err)
	}
	o, _ := NewVersion("1.2.3")
	if v.Compare(o) != 0 {
		t.Errorf("v1.2.3 should be equal to 1.2.3")
	}
}
