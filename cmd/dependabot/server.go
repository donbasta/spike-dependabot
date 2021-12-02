package dependabot

import (
	"dependabot/internal/config"
	"dependabot/internal/task"

	"github.com/spf13/cobra"
)

type ServerCommand *cobra.Command

func NewServerCommand(
	cfg *config.Main,
) ServerCommand {
	command := &cobra.Command{
		Use:   "task",
		Short: "Run the server",
		Long:  "Run the server",
		Run: func(c *cobra.Command, args []string) {
			task.Run(cfg)
		},
	}

	return command
}

func InitializeServerCommand() (ServerCommand, error) {
	main := config.ProvideConfig()
	serverCommand := NewServerCommand(main)
	return serverCommand, nil
}
