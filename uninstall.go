package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var cmdUninstall = &Command{
	Run:       runUninstall,
	UsageLine: "uninstall ",
	Short:     "Uninstall Go",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdUninstall.Flag.BoolVar(&flagA, "a", false, "")
}

// uninstall uninstalls specified versions of Go.
func uninstall(versions []string) error {
	for _, version := range versions {
		dir := filepath.Join(gvmnrootGo, version)
		if !exist(dir) {
			return fmt.Errorf("no Go version specified")
		}
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
	}
	return nil
}

// runUninstall executes uninstall command and return exit code.
func runUninstall(args []string) int {
	if len(args) == 0 {
		log.Print("gvmn uninstall: no Go versions specified")
		return 1
	}

	if err := uninstall(args); err != nil {
		log.Print(err)
		return 1
	}
	return 0
}
