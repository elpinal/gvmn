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
func List() []Info {
	if !exist(gvmnrootGo) {
		return nil
	}
	current, _ := os.Readlink(filepath.Join(gvmnrootGo, "current"))
	currentVersion := filepath.Base(current)
	versions, err := ioutil.ReadDir(gvmnrootGo)
	if err != nil {
		return nil
	}
	var info []Info
	for _, version := range versions {
		ver := version.Name()
		if ver == "current" {
			continue
		}
		var installed bool
		if exist(filepath.Join(gvmnrootGo, version.Name(), "bin", "go")) {
			installed = true
		}
		info = append(info, Info{Name: ver, Current: ver == currentVersion, Installed: installed})
	}
	return info
}
