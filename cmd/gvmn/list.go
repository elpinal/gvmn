package main

import (
	"log"

	"github.com/susp/gvmn"
)

var cmdList = &Command{
	Run:       runList,
	UsageLine: "list",
	Short:     "list installed Go versions",
	Long: `
List lists installed Go versions.
A Go version selected by gvmn use is marked.
	`,
}

func init() {
	// Set your flag here like below.
	// cmdList.Flag.BoolVar(&flagA, "a", false, "")
}

// runList executes list command and return exit code.
func runList(args []string) int {
	if err := gvmn.List(); err != nil {
		log.Print(err)
		return 1
	}
	return 0
}
