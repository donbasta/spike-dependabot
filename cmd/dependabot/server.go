package dependabot

import (
	"dependabot/internal/client"
	"dependabot/internal/config"
	"dependabot/internal/service"
	"log"

	"github.com/spf13/cobra"
)

type ServerCommand *cobra.Command

func NewServerCommand(
	cfg *config.Main,
) ServerCommand {
	command := &cobra.Command{
		Use:   "server",
		Short: "Run the server",
		Long:  "Run the server",
		Run: func(c *cobra.Command, args []string) {
			client := client.CreateClient(cfg.Git.Token)

			projects, err := service.CrawlGroups(client)
			if err != nil {
				log.Fatal(err)
			}
			changes := service.CheckDependency(client, projects)
			service.UpdateProjects(changes)
		},
	}

	return command
}

func InitializeServerCommand() (ServerCommand, error) {
	main := config.ProvideConfig()
	serverCommand := NewServerCommand(main)
	return serverCommand, nil
}
