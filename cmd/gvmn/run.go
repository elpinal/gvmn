package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/elpinal/gvmn"
)

var cmdRun = &Command{
	Run:       runRun,
	UsageLine: "run version command [arguments]",
	Short:     "execute the specified Go version",
	Long:      `Run executes the specified Go version.`,
}

func runRun(cmd *Command, args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: %s\n\n", cmd.UsageLine)
		fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(cmd.Long))
		return 2
	}

	var goArgs []string
	if len(args) != 1 {
		goArgs = args[1:]
	}
	if err := gvmn.Run(args[0], goArgs...); err != nil {
		log.Print(err)
		return 1
	}
	return 0
}
