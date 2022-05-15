package goproject

import (
	"go/scanner"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func New(root string) *Project {
	project := Project{
		Root: root,
	}
	project.Update()
	return &project
}

type Project struct {
	Root string

	Readme    *File
	Changelog *File
	License   *File
	GoMod     *File
	GoFiles   []*File
}

func (me *Project) Update() {
	// reset
	me.GoFiles = nil
	me.Readme = nil
	me.Changelog = nil
	me.License = nil
	me.GoMod = nil

	filepath.WalkDir(me.Root, func(path string, d fs.DirEntry, err error) error {
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
	})
}

func (me *Project) AddFile(f *File) {
	switch f.Name() {
	case "README.md":
		me.Readme = f
	case "changelog.txt", "Changelog.md", "CHANGELOG.md":
		me.Changelog = f
	case "go.mod":
		me.GoMod = f
	case "LICENSE", "license.txt":
		me.License = f
	case ".gitignore", ".onchange.sh", "go.sum":

	default:
		if filepath.Ext(f.Name()) == ".go" {
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

type File struct {
	Path string
	fs.DirEntry
}

func (me *File) Types() []string {
	if filepath.Ext(me.Path) != ".go" {
		return nil
	}
	src, err := os.ReadFile(me.Path)
	if err != nil {
		panic(err.Error())
	}

	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile(me.Path, fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)

	res := []string{}
	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		if tok != token.TYPE {
			continue
		}
		_, _, lit = s.Scan()
		if lit == "" {
			continue
		}
		res = append(res, lit)
	}
	return res
}
