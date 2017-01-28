package build

import (
	"testing"
)

func TestExpand(t *testing.T) {
	index := NewIndex()
	if index.Len() != 1 {
		t.Errorf("Bad index length: %d", index.Len())
	}
	index.Expand()
	if index.Len() != 2 {
		t.Errorf("Bad index length: %d", index.Len())
	}
}

func TestShrink(t *testing.T) {
	index := NewIndex()
	index.Expand()
	if index.Len() != 2 {
		t.Errorf("Bad index length: %d", index.Len())
	}
	index.Shrink()
	if index.Len() != 1 {
		t.Errorf("Bad index length: %d", index.Len())
	}
	index.Shrink()
	if index.Len() != 0 {
		t.Errorf("Bad index length: %d", index.Len())
	}
}

func TestSet(t *testing.T) {
	index := NewIndex()
	if index.Index[0] != 0 {
		t.Errorf("Bad index value: %d", index.Index[0])
	}
	index.Set(1)
	if index.Index[0] != 1 {
		t.Errorf("Bad index value: %d", index.Index[0])
	}
}

func TestString(t *testing.T) {
	index := NewIndex()
	if index.String() != "1" {
		t.Errorf("Bad string value: %s", index.String())
	}
	index.Set(1)
	if index.String() != "2" {
		t.Errorf("Bad string value: %s", index.String())
	}
	index.Expand()
	if index.String() != "2.1" {
		t.Errorf("Bad string value: %s", index.String())
	}
}
