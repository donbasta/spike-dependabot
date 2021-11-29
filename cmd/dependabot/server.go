package dependabot

import (
	"dependabot/internal/config"

	"github.com/spf13/cobra"
)

type ServerCommand *cobra.Command

func NewServerCommand(
	config *config.Main,
) ServerCommand {
	command := &cobra.Command{
		Use:   "server",
		Short: "Run the server",
		Long:  "Run the server",
	}

	return command
}
