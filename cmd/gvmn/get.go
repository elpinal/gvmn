package main

import (
	"log"

	"github.com/susp/gvmn"
)

var cmdGet = &Command{
	Run:       runGet,
	UsageLine: "get version",
	Short:     "download and install Go",
	Long: `
Get downloads the specified Go version, and then installs it.
	`,
}

func init() {
	// Set your flag here like below.
	// cmdGet.Flag.BoolVar(&flagA, "a", false, "")
}

// runGet executes get command and return exit code.
func runGet(args []string) int {
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

	if err := gvmn.Get(version); err != nil {
		log.Print(err)
		return 1
	}

	return 0
}
