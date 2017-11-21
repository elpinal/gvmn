package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mitchellh/go-homedir"

	"github.com/elpinal/gvmn"
)

// A Command is an implementation of a gvmn command
type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string) int

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
	cmdGet,
	cmdList,
	cmdUse,
	cmdRemove,
	cmdRun,
	cmdVersion,
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

	if root := os.Getenv("GVMNROOT"); root != "" {
		gvmn.SetRoot(root)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.Print(err)
			os.Exit(2)
		}
		gvmn.SetRoot(filepath.Join(home, ".gvmn"))
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] {
			cmd.Flag.Usage = func() { cmd.Usage() }

			cmd.Flag.Parse(args[1:])
			args = cmd.Flag.Args()

			os.Exit(cmd.Run(cmd, args))
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
