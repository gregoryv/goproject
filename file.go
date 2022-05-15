package goproject

import (
	"go/scanner"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
)

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
