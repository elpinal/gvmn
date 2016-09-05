package gvmn

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// List prints installed Go versions.
// A currently used Go version is marked by *.
func List() error {
	if !exist(gvmnrootGo) {
		return nil
	}
	current, _ := os.Readlink(filepath.Join(gvmnrootGo, "current"))
	currentVersion := filepath.Base(current)
	versions, err := ioutil.ReadDir(gvmnrootGo)
	if err != nil {
		return errors.Wrap(err, "ReadDir")
	}
	for _, version := range versions {
		ver := version.Name()
		if ver == "current" {
			continue
		}
		var mark string
		if ver == currentVersion {
			mark = "*"
		} else {
			mark = " "
		}
		fmt.Println(mark, ver)
	}
	return nil
}
