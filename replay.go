package main

import (
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
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
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
			Path:     strings.TrimPrefix(path, root+"/"),
			DirEntry: d,
		}
		project.AddFile(f)
		return nil
	})
	return &project
}

type Project struct {
	Root  string
	Files []*File

	Readme    *File
	Changelog *File
	License   *File
	GoMod     *File
}

func (me *Project) WriteTo(w io.Writer) (int64, error) {
	fg := vt100.ForegroundColors()
	_ = vt100.BackgroundColors()
	vt := vt100.Attributes()

	p, err := nexus.NewPrinter(w)
	p.Print("\033[2J\033[f") // clear
	p.Println(fg.White, me.Root, vt.Reset)

	for _, f := range me.Files {
		p.Println(fg.Red, f.Path, vt.Reset)
	}

	return p.Written, *err
}

func (me *Project) AddFile(f *File) {
	switch f.Name() {
	case "README.md":
		me.Readme = f
	default:
		me.Files = append(me.Files, f)
	}
}

type File struct {
	Path string
	fs.DirEntry
	Types []string
}
