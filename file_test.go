package goproject

import (
	"testing"
)

func TestParseTypes(t *testing.T) {
	f := &File{
		Path: "project.go",
	}
	if len(f.ParseTypes()) == 0 {
		t.Error("missing types in project.go")
	}
}

func TestParseVars(t *testing.T) {
	f := &File{
		Path: "file_test.go",
	}
	var X int // should be picked up
	_ = X
	if len(f.ParseVars()) == 0 {
		t.Error("no vars found in file_test.go")
	}
}

func TestTypes_Add(t *testing.T) {
	var v Types
	v.Add("")
	if len(v) > 0 {
		t.Error("empty value was added")
	}
}
