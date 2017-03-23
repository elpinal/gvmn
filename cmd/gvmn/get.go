package main

import (
	"log"

	"github.com/susp/gvmn"
)

var cmdGet = &Command{
	Run:       runGet,
	UsageLine: "get [-b] [-d] versions...",
	Short:     "download and install Go",
	Long: `
Get downloads the specified Go versions, and then installs them.

The -b flag instructs get to download binaries of the Go versions.

The -d flag instructs get to stop after downloading the Go versions; that is,
it instructs get not to install the Go versions.
`,
}

var (
	getB bool
	getD bool
)

func init() {
	cmdGet.Flag.BoolVar(&getB, "b", false, "")
	cmdGet.Flag.BoolVar(&getD, "d", false, "")
}

// runGet executes get command and return exit code.
func runGet(cmd *Command, args []string) int {
	if len(args) == 0 {
		log.Print("gvmn get: no go versions specified")
		return 1
	}

	if getB {
		if err := getBinary(args); err != nil {
			log.Print(err)
			return 1
		}
		return 0
	}

	for i, version := range args {
		if version == "latest" {
			latest, err := gvmn.LatestTag()
			if err != nil {
				log.Print(err)
				return 1
			}
			version = latest
			args[i] = latest
		}

		if err := gvmn.Download(version, false); err != nil {
			log.Print(err)
			return 1
		}
	}

	if getD {
		return 0
	}

	for _, version := range args {
		if err := gvmn.Install(version); err != nil {
			log.Print(err)
			return 1
		}
	}

	if err := gvmn.Use(args[len(args)-1:][0]); err != nil {
		log.Print(err)
		return 1
	}

	return 0
}

func getBinary(versions []string) error {
	for _, version := range versions {
		if err := gvmn.GetBinary(version); err != nil {
			return err
		}
	}
	return gvmn.Use(versions[len(versions)-1:][0])
}
