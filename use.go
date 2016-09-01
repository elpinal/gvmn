package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

var cmdUse = &Command{
	Run:       runUse,
	UsageLine: "use version",
	Short:     "select a Go version to use",
	Long: `
Use selects a Go version to use.
	`,
}

func init() {
	// Set your flag here like below.
	// cmdUse.Flag.BoolVar(&flagA, "a", false, "")
}

// runUse executes use command and return exit code.
func runUse(args []string) int {
	if len(args) == 0 {
		log.Print("gvmn use: no go version specified")
		return 1
	}
	currentDir := filepath.Join(gvmnrootGo, "current")
	version := args[0]
	versionsDir := filepath.Join(gvmnrootGo, version)
	if !exist(versionsDir) {
		log.Print("no installed go version specified")
		return 1
	}
	if err := os.RemoveAll(currentDir); err != nil {
		log.Print(errors.Wrap(err, "failed to stop using former go version"))
		return 1
	}
	err := os.Symlink(versionsDir, currentDir)
	if err != nil {
		log.Print(errors.Wrap(err, "failed to create symbolic link"))
		return 1
	}

	return 0
}
