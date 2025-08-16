package commands

import (
	"context"
	"github.com/h4rdc0m/aurora-api/domain/common"
	"github.com/h4rdc0m/aurora-api/infrastructure/logging"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var cmds = map[string]common.Command{
	"app:serve": NewServeCommand(),
}

// GetSubCommands gives a list of sub commands
func GetSubCommands(opt fx.Option) []*cobra.Command {
	var subCommands []*cobra.Command
	for name, cmd := range cmds {
		subCommands = append(subCommands, WrapSubCommand(name, cmd, opt))
	}
	return subCommands
}

func WrapSubCommand(name string, cmd common.Command, opt fx.Option) *cobra.Command {
	wrappedCmd := &cobra.Command{
		Use:   name,
		Short: cmd.Short(),
		Run: func(c *cobra.Command, args []string) {
			opts := fx.Options(
				fx.WithLogger(logging.NewFxLogger),
				fx.Invoke(cmd.Run()),
			)
			ctx := context.Background()
			app := fx.New(opt, opts)
			err := app.Start(ctx)
			defer func(app *fx.App, ctx context.Context) {
				err := app.Stop(ctx)
				if err != nil {
					panic(err)
				}
			}(app, ctx)
			if err != nil {
				panic(err)
			}
		},
	}
	cmd.Setup(wrappedCmd)
	return wrappedCmd
}
