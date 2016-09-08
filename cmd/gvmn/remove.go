package main

import (
	"log"

	"github.com/susp/gvmn"
)

var cmdRemove = &Command{
	Run:       runRemove,
	UsageLine: "remove versions...",
	Short:     "remove Go versions",
	Long:      `Remove removes the specified Go versions.`,
}

// runRemove executes remove command and return exit code.
func runRemove(cmd *Command, args []string) int {
	if len(args) == 0 {
		log.Print("gvmn remove: no go versions specified")
		return 1
	}
	if err := gvmn.Remove(args...); err != nil {
		log.Print(err)
		return 1
	}
	return 0
}
