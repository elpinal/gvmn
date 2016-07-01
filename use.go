package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

var cmdUse = &Command{
	Run:       runUse,
	UsageLine: "use ",
	Short:     "Use Go",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdUse.Flag.BoolVar(&flagA, "a", false, "")
}

// runUse executes use command and return exit code.
func runUse(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "gvmn use: no Go version specified")
		return 1
	}
	currentDir := filepath.Join(GvmnDir, "versions", "current")
	version := args[0]
	versionsDir := filepath.Join(GvmnDir, "versions", version)
	if _, err := os.Stat(versionsDir); err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "no installed version of Go specified"))
		return 1
	}
	if _, err := os.Stat(currentDir); err == nil {
		if err := os.RemoveAll(currentDir); err != nil {
			fmt.Fprintln(os.Stderr, errors.Wrap(err, "failed unuse former version of Go"))
			return 1
		}
	}
	err := os.Symlink(versionsDir, currentDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "failed create symbolic link"))
		return 1
	}

	return 0
}
