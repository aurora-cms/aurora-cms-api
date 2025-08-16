package common

import "github.com/spf13/cobra"

type CommandRunner interface {
}

type Command interface {
	// Short returns a short description of the command.
	// The string is shown in the help output.
	Short() string

	// Setup is used to set up flags or pre-run logic for the command.
	//
	// Example:
	// cms.Flags().IntVar(&r.num, "num", 0, "Number of times to run the command")
	Setup(cmd *cobra.Command)

	// Run the command runner
	// Returns a CommandRunner which is a function with dependency-injected arguments.
	//
	// Example:
	// Command{
	// 	 Run: func(l lib.Logger) {
	// 		   l.Info("Running command with num:", r.num)
	// 	 },
	// }
	Run() CommandRunner
}
