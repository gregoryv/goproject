package main

import (
	"go/scanner"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gregoryv/nexus"
	"github.com/gregoryv/vt100"
)

func main() {
	root, _ := os.Getwd()
	project := NewProject(root)
	project.WriteTo(os.Stdout)
}

func NewProject(root string) *Project {
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
	Files     []*File
}

func (me *Project) Update() {
	// reset
	me.Files = nil
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

func (me *Project) WriteTo(w io.Writer) (int64, error) {
	fg := vt100.ForegroundColors()
	_ = vt100.BackgroundColors()
	vt := vt100.Attributes()

	p, err := nexus.NewPrinter(w)
	p.Print("\033[2J\033[f") // clear
	p.Println(fg.White, me.Root, vt.Reset)

	if v := me.Special(); len(v) > 0 {
		p.Println(fg.Yellow, strings.Join(v, "  "), vt.Reset)
	}
	p.Println()

	for _, f := range me.Files {
		p.Println(fg.White, f.Path, vt.Reset)
		if types := f.Types(); len(types) > 0 {
			p.Println(fg.Cyan, "  ", strings.Join(types, ", "), vt.Reset)
		}
	}

	return p.Written, *err
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
		me.Files = append(me.Files, f)
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
