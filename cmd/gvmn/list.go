package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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
func runList(cmd *Command, args []string) int {
	if len(args) != 0 {
		fmt.Fprintf(os.Stderr, "usage: %s\n\n", cmd.UsageLine)
		fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(cmd.Long))
		return 2
	}
	if err := gvmn.List(); err != nil {
		log.Print(err)
		return 1
	}
	return 0
}
