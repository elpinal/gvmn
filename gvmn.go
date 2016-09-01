package gvmn

import (
	"os"
	"path/filepath"
)

var RepoURL = "git://github.com/golang/go.git"

var (
	gvmnroot     string
	gvmnrootGo   string
	gvmnrootRepo string
)

// Setroot sets root as gvmn's root directory.
func SetRoot(root string) {
	gvmnroot = root
	gvmnrootGo = filepath.Join(root, "go")
	gvmnrootRepo = filepath.Join(root, "repo")
}

func exist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
