package gvmn

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func Use(version string) error {
	currentDir := filepath.Join(gvmnrootGo, "current")
	versionsDir := filepath.Join(gvmnrootGo, version)
	if !exist(versionsDir) {
		return fmt.Errorf("no installed go version specified")
	}
	if err := os.RemoveAll(currentDir); err != nil {
		return errors.Wrap(err, "failed to stop using former go version")
	}
	err := os.Symlink(versionsDir, currentDir)
	if err != nil {
		return errors.Wrap(err, "failed to create symbolic link")
	}
	return nil
}
