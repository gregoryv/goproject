package main

import (
	"io"
	"os"
	"strings"

	"github.com/gregoryv/goproject"
	"github.com/gregoryv/nexus"
	"github.com/gregoryv/vt100"
)

func main() {
	root, _ := os.Getwd()
	project := goproject.New(root)
	showProject(os.Stdout, project)
}

func showProject(w io.Writer, project *goproject.Project) (int64, error) {
	fg := vt100.ForegroundColors()
	_ = vt100.BackgroundColors()
	vt := vt100.Attributes()

	p, err := nexus.NewPrinter(w)
	p.Print("\033[2J\033[f") // clear
	p.Println(fg.White, project.Root, vt.Reset)

	if v := project.Special(); len(v) > 0 {
		p.Println(fg.Yellow, strings.Join(v, "  "), vt.Reset)
	}
	p.Println()

	for _, f := range project.Files {
		p.Println(fg.White, f.Path, vt.Reset)
		if types := f.Types(); len(types) > 0 {
			p.Println(fg.Cyan, "  ", strings.Join(types, ", "), vt.Reset)
		}
	}

	return p.Written, *err
}
