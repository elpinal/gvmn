package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/elpinal/gvmn"
)

var cmdList = &Command{
	Run:       runList,
	UsageLine: "list",
	Short:     "list installed Go versions",
	Long: `
List lists installed Go versions.
Go versions is divided into states.
	`,
}

func init() {
	// Set your flag here like below.
	// cmdList.Flag.BoolVar(&flagA, "a", false, "")
}

func doOnce(f func()) func() {
	var done bool
	return func() {
		if done {
			return
		}
		f()
		done = true
	}
}

// runList executes list command and return exit code.
func runList(cmd *Command, args []string) int {
	if len(args) != 0 {
		fmt.Fprintf(os.Stderr, "usage: %s\n\n", cmd.UsageLine)
		fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(cmd.Long))
		return 2
	}
	l := lister{
		out: os.Stdout,
		err: os.Stderr,
	}
	return l.listMain()
}

type stringWriter interface {
	Write([]byte) (int, error)
	WriteString(string) (int, error)
}

type lister struct {
	out stringWriter
	err stringWriter
}

func newline(w stringWriter) {
	w.Write([]byte{'\n'})
}

func (l *lister) genHeader(header string) func() {
	return func() {
		newline(l.out)
		l.out.WriteString(header)
		newline(l.out)
	}
}

func (l *lister) listMain() int {
	list, err := gvmn.List()
	if err != nil {
		fmt.Fprintln(l.err, err)
		return 1
	}
	if list == nil {
		return 0
	}

	for _, info := range list {
		if info.Current {
			l.out.WriteString("Current:")
			newline(l.out)
			fmt.Fprint(l.out, "\t", info.Name)
			break
		}
	}

	ih := doOnce(l.genHeader("Installed:"))
	for _, info := range list {
		if info.Installed {
			ih()
			fmt.Fprint(l.out, "\t", info.Name)
		}
	}

	dh := doOnce(l.genHeader("Just downloaded; not installed:"))
	for _, info := range list {
		if !info.Current && !info.Installed {
			dh()
			fmt.Fprint(l.out, "\t", info.Name)
		}
	}

	newline(l.out)

	return 0
}
