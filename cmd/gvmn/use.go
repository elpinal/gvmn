package main

import (
	"log"

	"github.com/susp/gvmn"
)

var cmdUse = &Command{
	Run:       runUse,
	UsageLine: "use version",
	Short:     "select a Go version to use",
	Long:      `Use selects a Go version to use.`,
}

func init() {
	// Set your flag here like below.
	// cmdUse.Flag.BoolVar(&flagA, "a", false, "")
}

// runUse executes use command and return exit code.
func runUse(cmd *Command, args []string) int {
	if len(args) == 0 {
		log.Print("gvmn use: no go version specified")
		return 1
	}
	if err := gvmn.Use(args[0]); err != nil {
		log.Print(err)
		return 1
	}
	return 0
}
