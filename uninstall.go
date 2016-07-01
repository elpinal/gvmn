package main

var cmdUninstall = &Command{
	Run:       runUninstall,
	UsageLine: "uninstall ",
	Short:     "Uninstall Go",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdUninstall.Flag.BoolVar(&flagA, "a", false, "")
}

// runUninstall executes uninstall command and return exit code.
func runUninstall(args []string) int {

	return 0
}
