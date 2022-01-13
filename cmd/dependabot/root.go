package dependabot

import (
	"github.com/spf13/cobra"
)

var (
	cliName = "dependabot"
)

func NewRootCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   cliName,
		Short: "Root command for dependabot",
		Long:  "Root command for dependabot",
		Run: func(c *cobra.Command, args []string) {
			c.HelpFunc()(c, args)
		},
	}

	serverCommand, err := InitializeServerCommand()
	if err != nil {
		panic(err)
	}

	command.AddCommand(
		serverCommand,
		newAutoMigrate(),
		newMigrate(),
	)
	return command
}
