package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var cmdRemove = &Command{
	Run:       runRemove,
	UsageLine: "remove versions...",
	Short:     "remove Go versions",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdRemove.Flag.BoolVar(&flagA, "a", false, "")
}

// remove removes the specified Go versions.
func remove(versions []string) error {
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

// runRemove executes remove command and return exit code.
func runRemove(args []string) int {
	if len(args) == 0 {
		log.Print("gvmn remove: no go versions specified")
		return 1
	}

	if err := remove(args); err != nil {
		log.Print(err)
		return 1
	}
	return 0
}
