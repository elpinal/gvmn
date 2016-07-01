package main

var cmdInstall = &Command{
	Run:       runInstall,
	UsageLine: "install ",
	Short:     "Install Go",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdInstall.Flag.BoolVar(&flagA, "a", false, "")
}

// runInstall executes install command and return exit code.
func runInstall(args []string) int {

	return 0
}
