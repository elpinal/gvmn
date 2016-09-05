package gvmn

import (
	"fmt"
	"os"
	"path/filepath"
)

// Remove removes the specified Go versions.
func Remove(versions []string) error {
	for _, version := range versions {
		dir := filepath.Join(gvmnrootGo, version)
		if !exist(dir) {
			return fmt.Errorf("no go version specified")
		}
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
	}
	return nil
}
