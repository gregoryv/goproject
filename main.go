package main

import (
	"io"
	"os"
	"strings"
	"unicode"

	"github.com/gregoryv/cmdline"
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
	project := LoadProject(".")
	showProject(os.Stdout, project)
}

func showProject(w io.Writer, project *Project) (int64, error) {
	fg := vt100.ForegroundColors()
	vt := vt100.Attributes()

	p, err := nexus.NewPrinter(w)
	p.Print(vt.Bright, project.Package, vt.Reset, "\n")

	if v := project.Special(); len(v) > 0 {
		p.Print(fg.Cyan, "", vt.Dim, strings.Join(v, "  "), vt.Reset, "\n")
	}

	var noTypes []string
	for _, f := range project.GoFiles {
		types := f.ParseTypes()
		if len(types) == 0 {
			noTypes = append(noTypes, f.Name())
			continue
		}
		vars := f.ParseVars()
		p.Print(
			vt.Dim, "", fg.Cyan, f.Path, " ",
			vt.Reset, strings.Join(public(types), ", "), " ",
			vt.Dim, fg.Yellow, strings.Join(private(types), ", "), " ",
			fg.Magenta, strings.Join(public(vars), ", "), vt.Reset,
			"\n",
		)
	}
	p.Print(vt.Dim, "", strings.Join(noTypes, ", "), " (without types)", vt.Reset, "\n")
	return p.Written, *err
}

// public returns a slice of all words starting with uppercase letter
func public(v []string) []string {
	res := make([]string, 0)
	for _, name := range v {
		for i, r := range name {
			if i > 0 {
				break
			}
			if unicode.IsUpper(r) {
				res = append(res, name)
			}
		}
	}
	return res
}

// public returns a slice of all words starting with lowercase letter
func private(v []string) []string {
	res := make([]string, 0)
	for _, name := range v {
		for i, r := range name {
			if i > 0 {
				break
			}
			if unicode.IsLower(r) {
				res = append(res, name)
			}
		}
	}
	return res
}

func longest(files []*File) int {
	var l int
	for _, f := range files {
		if got := len(f.Name()); got > l {
			l = got
		}
	}
	return l
}
