package dependabot

import (
	"dependabot/internal/client"
	"dependabot/internal/config"
	"dependabot/internal/service"
)

func Run() {
	mainConfig := config.ProvideConfig()

	client := client.CreateClient(mainConfig.Git.Token)

	service.CrawlGroups(client)
}
