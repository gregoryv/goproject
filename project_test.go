package goproject

import (
	"errors"
	"os"
	"testing"
)

func TestLoadProject(t *testing.T) {
	wd, _ := os.Getwd()
	project := LoadProject(wd)

	t.Run("default content", func(t *testing.T) {
		if len(project.GoFiles) <= 2 {
			t.Error(project.GoFiles)
		}
	})

	t.Run("load", func(t *testing.T) {
		if err := project.load("", nil, errors.New("x")); err == nil {
			t.Error("expected an error")
		}
	})

	t.Run("Special", func(t *testing.T) {
		if got := project.Special(); len(got) < 4 {
			t.Error(got)
		}
	})

}
