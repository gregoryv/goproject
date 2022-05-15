package main

import (
	"io"
	"os"
	"strings"

	"github.com/gregoryv/cmdline"
	"github.com/gregoryv/goproject"
	"github.com/gregoryv/nexus"
	"github.com/gregoryv/vt100"
)

func main() {
	var (
		wd, _ = os.Getwd()
		cli   = cmdline.NewBasicParser()
		root  = cli.NamedArg("DIR").String(wd)
	)
	cli.Parse()

	os.Chdir(root)
	project := goproject.LoadProject(".")
	showProject(os.Stdout, project)
}

func showProject(w io.Writer, project *goproject.Project) (int64, error) {
	fg := vt100.ForegroundColors()
	_ = vt100.BackgroundColors()
	vt := vt100.Attributes()

	p, err := nexus.NewPrinter(w)
	//p.Print("\033[2J\033[f") // clear
	p.Println(fg.White, project.Package, vt.Reset)

	if v := project.Special(); len(v) > 0 {
		p.Println(fg.Yellow, strings.Join(v, "  "), vt.Reset)
	}
	p.Println()

	var noTypes []string
	for _, f := range project.GoFiles {
		types := f.ParseTypes()
		if len(types) == 0 {
			noTypes = append(noTypes, f.Name())
			continue
		}

		p.Println(
			fg.White, f.Path,
			fg.Cyan, strings.Join(types, ", "), vt.Reset,
		)
	}
	p.Println()
	p.Println(vt.Dim, strings.Join(noTypes, ", "), "(without types)", vt.Reset)
	return p.Written, *err
}

func longest(files []*goproject.File) int {
	var l int
	for _, f := range files {
		if got := len(f.Name()); got > l {
			l = got
		}
	}
	return l
}
