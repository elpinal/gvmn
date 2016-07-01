package main

var cmdUse = &Command{
	Run:       runUse,
	UsageLine: "use ",
	Short:     "Use Go",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdUse.Flag.BoolVar(&flagA, "a", false, "")
}

// runUse executes use command and return exit code.
func runUse(args []string) int {

	return 0
}
