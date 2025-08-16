package bootstrap

import (
	"github.com/h4rdc0m/aurora-api/api/commands"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aurora-api",
	Short: "Aurora API Server",
	Long: `
 ▄▄▄▄▄▄▄ ▄▄   ▄▄ ▄▄▄▄▄▄   ▄▄▄▄▄▄▄ ▄▄▄▄▄▄   ▄▄▄▄▄▄▄    ▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄ ▄▄▄ 
█       █  █ █  █   ▄  █ █       █   ▄  █ █       █  █       █       █   █
█   ▄   █  █ █  █  █ █ █ █   ▄   █  █ █ █ █   ▄   █  █   ▄   █    ▄  █   █
█  █▄█  █  █▄█  █   █▄▄█▄█  █ █  █   █▄▄█▄█  █▄█  █  █  █▄█  █   █▄█ █   █
█       █       █    ▄▄  █  █▄█  █    ▄▄  █       █  █       █    ▄▄▄█   █
█   ▄   █       █   █  █ █       █   █  █ █   ▄   █  █   ▄   █   █   █   █
█▄▄█ █▄▄█▄▄▄▄▄▄▄█▄▄▄█  █▄█▄▄▄▄▄▄▄█▄▄▄█  █▄█▄▄█ █▄▄█  █▄▄█ █▄▄█▄▄▄█   █▄▄▄█
version: 0.3.0

This is a command runner or cli for Aurora API. 
Using this we can use underlying dependency injection container for running scripts. 
Main advantage is that, we can use same services, repositories, infrastructure present in the application itself
`,
	TraverseChildren: true,
}

// App root of the application
type App struct {
	*cobra.Command
}

func NewApp() App {
	cmd := App{
		Command: rootCmd,
	}
	cmd.AddCommand(commands.GetSubCommands(CommonModules)...)
	return cmd
}

var RootApp = NewApp()
