package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

var cmdList = &Command{
	Run:       runList,
	UsageLine: "list ",
	Short:     "List versions of Go",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdList.Flag.BoolVar(&flagA, "a", false, "")
}

// runList executes list command and return exit code.
func runList(args []string) int {
	current, _ := os.Readlink(filepath.Join(GvmnDir, "versions", "current"))
	currentVersion := filepath.Base(current)
	versions, err := ioutil.ReadDir(filepath.Join(GvmnDir, "versions"))
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "failed to list versions of Go"))
	}
	for _, version := range versions {
		ver := version.Name()
		if ver == "current" {
			continue
		}
		var mark string
		if ver == currentVersion {
			mark = "*"
		} else {
			mark = " "
		}
		fmt.Println(mark, ver)
	}
	return 0
}
