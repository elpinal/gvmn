package gvmn

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func List() error {
	current, _ := os.Readlink(filepath.Join(gvmnrootGo, "current"))
	currentVersion := filepath.Base(current)
	versions, err := ioutil.ReadDir(gvmnrootGo)
	if err != nil {
		return errors.Wrap(err, "failed to list installed go versions")
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
