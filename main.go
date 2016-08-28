package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mitchellh/go-homedir"
)

// A Command is an implementation of a gvmn command
type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(args []string) int

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'gvmn help' output.
	Short string

	// Long is the long message shown in the 'gvmn help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet
}

// Name returns the command's name: the first word in the usage line.
func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.Long))
	os.Exit(2)
}

// Commands lists the available commands and help topics.
// The order here is the order in which they are printed by 'gvmn help'.
var commands = []*Command{
	cmdInstall,
	cmdList,
	cmdUse,
	cmdUninstall,
}

var RepoURL = "git://github.com/golang/go.git"
var GvmnDir string

func exist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func main() {
	flag.Usage = usage
	flag.Parse()
	log.SetFlags(0)

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	if args[0] == "help" {
		help(args[1:])
		return
	}

	home, err := homedir.Dir()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	GvmnDir = filepath.Join(home, ".gvmn")

	if !exist(GvmnDir) {
		if err := os.MkdirAll(GvmnDir, 0777); err != nil {
			log.Print(err)
			os.Exit(1)
		}
	}
	etcDir := filepath.Join(GvmnDir, "etc")
	if !exist(etcDir) {
		if err := os.Mkdir(etcDir, 0777); err != nil {
			log.Print(err)
			os.Exit(1)
		}
	}
	loginFile := filepath.Join(GvmnDir, "etc", "login")
	if !exist(loginFile) {
		err := ioutil.WriteFile(loginFile, []byte(strings.TrimSpace(`
#!/bin/bash

__gvmn_configure_path()
{
  local gvmn_bin_path="$HOME/.gvmn/versions/current/bin"

  echo "$PATH" | grep -Fqv "$gvmn_bin_path" &&
    PATH="$gvmn_bin_path:$PATH"
}


__gvmn_configure_path

# __END__
		`)), 0666)
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] {
			cmd.Flag.Usage = func() { cmd.Usage() }

			cmd.Flag.Parse(args[1:])
			args = cmd.Flag.Args()

			os.Exit(cmd.Run(args))
		}
	}

	fmt.Fprintf(os.Stderr, "gvmn: unknown subcommand %q\nRun 'gvmn help' for usage.\n", args[0])
	os.Exit(2)
}

var usageTemplate = `gvmn is a tool for managing Go versions.

Usage:

	gvmn command [arguments]

The commands are:
{{range .}}
	{{.Name | printf "%-11s"}} {{.Short}}{{end}}

Use "gvmn help [command]" for more information about a command.

`

var helpTemplate = `usage: gvmn {{.UsageLine}}

{{.Long | trim}}
`

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace})
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func printUsage(w io.Writer) {
	bw := bufio.NewWriter(w)
	tmpl(bw, usageTemplate, commands)
	bw.Flush()
}

func usage() {
	printUsage(os.Stderr)
	os.Exit(2)
}

// help implements the 'help' command.
func help(args []string) {
	if len(args) == 0 {
		printUsage(os.Stdout)
		// not exit 2: succeeded at 'gvmn help'.
		return
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: gvmn help command\n\nToo many arguments given.\n")
		os.Exit(2) // failed at 'gvmn help'
	}

	arg := args[0]

	for _, cmd := range commands {
		if cmd.Name() == arg {
			tmpl(os.Stdout, helpTemplate, cmd)
			// not exit 2: succeeded at 'gvmn help cmd'.
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown help topic %#q.  Run 'gvmn help'.\n", arg)
	os.Exit(2) // failed at 'gvmn help cmd'
}
