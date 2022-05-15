package goproject

import (
	"os"
	"testing"
)

func TestProject(t *testing.T) {
	wd, _ := os.Getwd()
	_ = New(wd)
}
