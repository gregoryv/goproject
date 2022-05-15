package goproject

import (
	"os"
	"testing"
)

func TestProject(t *testing.T) {
	wd, _ := os.Getwd()
	project := New(wd)
	if len(project.GoFiles) <= 2 {
		t.Error(project.GoFiles)
	}
}
