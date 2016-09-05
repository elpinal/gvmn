// Package gvmn provides support for managing Go versions.
//
// First of all, SetRoot to determine root directory for gvmn:
//
//     gvmn.SetRoot("/path/to/root")
//
// Then, Get a specific Go version:
//
//     err := gvmn.Get("go1.7")
//     ...
//
package gvmn

import (
	"os"
	"path/filepath"
)

// RepoURL indicates the Go original repository url.
var RepoURL = "git://github.com/golang/go.git"

var (
	gvmnroot     string
	gvmnrootGo   string
	gvmnrootRepo string
)

// SetRoot sets root as gvmn's root directory.
func SetRoot(root string) {
	gvmnroot = root
	gvmnrootGo = filepath.Join(root, "go")
	gvmnrootRepo = filepath.Join(root, "repo")
}

func exist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
