package task

import (
	"dependabot/internal/config"
	"dependabot/internal/task/checker"
	"dependabot/internal/task/crawler"
	"dependabot/internal/task/helper"
	"dependabot/internal/task/updater"
	"log"
)

func Run(cfg *config.Main) {
	client, err := config.ProvideGitlabClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	groupIds := helper.GetDefaultGroupList()

	projects, err := crawler.CrawlMultipleGroups(client, groupIds)
	if err != nil {
		log.Fatal(err)
	}

	dependencyUpdates := checker.CheckMultipleProjectsDependencyUpdate(client, projects)

	updater.UpdateProjectDependencies(dependencyUpdates)
}
