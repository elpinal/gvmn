package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/elpinal/gvmn"
)

var cmdGet = &Command{
	Run:       runGet,
	UsageLine: "get [-b] [-d] [-u] versions...",
	Short:     "download and install Go",
	Long: `
Get downloads the specified Go versions, and then installs them.

The versions are exepected as Git's references on the repository of Go.
A version named "stable" and "tip" are interpreted as the latest tag and the
latest developing version respectively.

The -b flag instructs get to download binaries of the Go versions.

The -d flag instructs get to stop after downloading the Go versions; that is,
it instructs get not to install the Go versions.

The -u flag instructs get to use the network to update the Go version.
By default, get uses the network to check out missing Go versions but does not use
it to look for updates to existing Go versions.
`,
}

var (
	getB bool
	getD bool
	getU bool
)

func init() {
	cmdGet.Flag.BoolVar(&getB, "b", false, "")
	cmdGet.Flag.BoolVar(&getD, "d", false, "")
	cmdGet.Flag.BoolVar(&getU, "u", false, "")
}

// runGet executes get command and returns exit code.
func runGet(_ *Command, args []string) int {
	err := getMain(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

type logger struct {
	w io.Writer
}

func (l *logger) Printf(format string, args ...interface{}) {
	fmt.Fprintf(l.w, format, args...)
}

func getMain(args []string) error {
	if len(args) == 0 {
		return errors.New("gvmn get: no go versions specified")
	}

	if getB {
		if err := getBinary(args); err != nil {
			return fmt.Errorf("getting binary: %v", err)
		}
		return nil
	}

	l := &logger{w: os.Stdout}

	for i, version := range args {
		switch version {
		case "stable":
			stable, err := gvmn.LatestTag()
			if err != nil {
				return fmt.Errorf("obtaining the latest tag: %v", err)
			}
			version = stable
			args[i] = stable
			l.Printf("stable: %s\n", stable)
		case "tip":
			tip, err := gvmn.Tip()
			if err != nil {
				return fmt.Errorf("obtaining the reference of tip: %v", err)
			}
			version = tip
			args[i] = tip
			l.Printf("tip: %s\n", tip)
		}

		if err := gvmn.Download(version, getU); err != nil {
			return fmt.Errorf("downloading (%s): %v", version, err)
		}
	}

	if getD {
		return nil
	}

	for _, version := range args {
		if err := gvmn.Install(version); err != nil {
			return fmt.Errorf("installing (%s): %v", version, err)
		}
	}

	return nil
}

func getBinary(versions []string) error {
	for _, version := range versions {
		if err := gvmn.GetBinary(version); err != nil {
			return err
		}
	}
	return nil
}
