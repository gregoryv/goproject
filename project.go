package main

import (
	"bytes"
	"io/fs"
	"os/exec"
	"path/filepath"
	"strings"
)

func LoadProject(root string) *Project {
	out, _ := exec.Command("go", "list", root).Output()
	project := Project{
		Root:    root,
		Package: string(bytes.TrimSpace(out)),
	}
	project.Update()
	return &project
}

type Project struct {
	Root    string
	Package string

	Readme    *File
	Changelog *File
	License   *File
	GoMod     *File
	GoFiles   []*File
	TestFiles []*File
}

func (me *Project) Update() {
	me.Reset()
	filepath.WalkDir(me.Root, me.load)
}

func (me *Project) Reset() {
	me.GoFiles = nil
	me.Readme = nil
	me.Changelog = nil
	me.License = nil
	me.GoMod = nil
}

func (me *Project) load(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if d.IsDir() && d.Name() == ".git" {
		return filepath.SkipDir
	}
	if d.IsDir() {
		return nil
	}

	f := &File{
		Path:     strings.TrimPrefix(path, me.Root+"/"),
		DirEntry: d,
	}
	me.AddFile(f)
	return nil
}

func (me *Project) AddFile(f *File) {
	switch f.Name() {
	case "README.md":
		me.Readme = f
	case "changelog.txt", "Changelog.md", "CHANGELOG.md", "CHANGELOG":
		me.Changelog = f
	case "go.mod":
		me.GoMod = f
	case "LICENSE", "license.txt":
		me.License = f
	case ".gitignore", ".onchange.sh", "go.sum":

	default:
		switch {
		case strings.HasSuffix(f.Name(), "_test.go"):
			me.TestFiles = append(me.TestFiles, f)

		case filepath.Ext(f.Name()) == ".go":
			me.GoFiles = append(me.GoFiles, f)
		}
	}
}

func (me *Project) Special() []string {
	special := []string{}
	if me.Readme != nil {
		special = append(special, me.Readme.Name())
	}
	if me.Changelog != nil {
		special = append(special, me.Changelog.Name())
	}
	if me.License != nil {
		special = append(special, me.License.Name())
	}
	if me.GoMod != nil {
		special = append(special, me.GoMod.Name())
	}
	return special
}
