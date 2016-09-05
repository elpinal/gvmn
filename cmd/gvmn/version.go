package main

import (
	"fmt"
	"os"

	"github.com/susp/gvmn"
)

var cmdVersion = &Command{
	Run:       runVersion,
	UsageLine: "version",
	Short:     "print gvmn version",
	Long:      `Version prints the gvmn version.`,
}

func runVersion(cmd *Command, args []string) int {
	if len(args) != 0 {
		fmt.Fprint(os.Stderr, "usage: version\n\n")
		fmt.Fprint(os.Stderr, "Version prints the gvmn version.\n")
		return 2
	}

	fmt.Printf("gvmn version %s\n", gvmn.Version)
	return 0
}
