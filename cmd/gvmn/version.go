package main

import (
	"fmt"
	"os"
	"strings"

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
		fmt.Fprintf(os.Stderr, "usage: %s\n\n", cmd.UsageLine)
		fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(cmd.Long))
		return 2
	}

	fmt.Printf("gvmn version %s\n", gvmn.Version)
	return 0
}
