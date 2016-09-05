package main

import (
	"log"

	"github.com/susp/gvmn"
)

var cmdGet = &Command{
	Run:       runGet,
	UsageLine: "get [-d] version",
	Short:     "download and install Go",
	Long:      `
Get downloads the specified Go version, and then installs it.

The -d flag instructs get to stop after downloading the Go version; that is,
it instructs get not to install the Go version.
`,
}

var getD bool

func init() {
	cmdGet.Flag.BoolVar(&getD, "d", false, "")
}

// runGet executes get command and return exit code.
func runGet(cmd *Command, args []string) int {
	if len(args) == 0 {
		log.Print("gvmn get: no go version specified")
		return 1
	}

	version := args[0]
	if version == "latest" {
		latest, err := gvmn.LatestTag()
		if err != nil {
			log.Print(err)
			return 1
		}
		version = latest
	}

	if err := gvmn.Download(version); err != nil {
		log.Print(err)
		return 1
	}

	if getD {
		return 0
	}

	if err := gvmn.Install(version); err != nil {
		log.Print(err)
		return 1
	}

	return 0
}
