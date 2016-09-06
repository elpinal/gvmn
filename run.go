package gvmn

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Run executes the specified Go version.
func Run(version string, args ...string) error {
	goCmd := filepath.Join(gvmnrootGo, version, "bin", "go")
	if !exist(goCmd) {
		return fmt.Errorf("no installed go version specified")
	}
	cmd := exec.Command(goCmd, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if _, isExitError := err.(*exec.ExitError); !isExitError {
			return err
		}
	}
	return nil
}
