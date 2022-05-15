package goproject

import (
	"go/scanner"
	"go/token"
	"io/fs"
	"os"
)

type File struct {
	Path string
	fs.DirEntry
}

func (me *File) ParseTypes() Types {
	src, _ := os.ReadFile(me.Path)

	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile(me.Path, fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)

	var res Types
	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		if tok != token.TYPE {
			continue
		}
		_, _, lit = s.Scan()
		res.Add(lit)
	}
	return res
}

type Types []string

func (me *Types) Add(v string) {
	if v == "" {
		return
	}
	*me = append(*me, v)
}
