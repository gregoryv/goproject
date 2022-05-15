package goproject

import (
	"io/fs"
	"path/filepath"
	"testing"
)

func TestFile_Types(t *testing.T) {
	filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		f := &File{
			Path:     path,
			DirEntry: d,
		}
		if d.Name() == "project.go" && len(f.ParseTypes()) == 0 {
			t.Error("missing types in project.go")
		}
		return nil
	})
}

func TestTypes_Add(t *testing.T) {
	var v Types
	v.Add("")
	if len(v) > 0 {
		t.Error("empty value was added")
	}
}
