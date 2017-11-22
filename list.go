package gvmn

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// A Info represents Go name and states.
type Info struct {
	Name      string
	Current   bool
	Installed bool
}

// List returns information of installed Go versions.
func List() (current *Info, all []Info, err error) {
	if !exist(gvmnrootGo) {
		return
	}

	dir := filepath.Join(gvmnrootGo, "current")
	currentPath, e := os.Readlink(dir)
	// Ignore, if any, an error if no version is activated.
	if e != nil && exist(dir) {
		err = e
		return
	}
	currentVersion := filepath.Base(currentPath)
	versions, e := ioutil.ReadDir(gvmnrootGo)
	if e != nil {
		err = e
		return
	}

	for _, version := range versions {
		ver := version.Name()
		if ver == "current" {
			continue
		}
		var installed bool
		if exist(filepath.Join(gvmnrootGo, version.Name(), "bin", "go")) {
			installed = true
		}
		i := Info{Name: ver, Current: ver == currentVersion, Installed: installed}
		all = append(all, i)
		if ver == currentVersion {
			current = &i
		}
	}
	return
}
