package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/susp/gvmn"
)

var cmdUse = &Command{
	Run:       runUse,
	UsageLine: "use version",
	Short:     "select a Go version to use",
	Long:      `Use selects a Go version to use.`,
}

// runUse executes use command and return exit code.
func runUse(cmd *Command, args []string) int {
	if len(args) == 0 {
		log.Print("gvmn use: no go version specified")
		return 1
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: %s\n\n", cmd.UsageLine)
		fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(cmd.Long))
		return 2
	}
	if err := gvmn.Use(args[0]); err != nil {
		log.Print(err)
		return 1
	}
	return 0
}
