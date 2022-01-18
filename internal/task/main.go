package task

import (
	"dependabot/internal/config"
	"dependabot/internal/task/checker"
	"dependabot/internal/task/crawler"
	"dependabot/internal/task/helper"
	"dependabot/internal/task/updater"
	"log"

	"github.com/gopaytech/go-commons/pkg/zlog"
)

func Run(cfg *config.Main) {
	client, err := config.ProvideGitlabClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	groupIds := helper.GetDefaultGroupList()

	zlog.Info("Start crawling repositories...")
	projects, err := crawler.CrawlMultipleGroups(client, groupIds)
	if err != nil {
		log.Fatal(err)
	}
	zlog.Info("Finished reading current dependencies from the repositories...")

	zlog.Info("Start fetching the dependencies' latest release from their source...")
	dependencyUpdates := checker.CheckMultipleProjectsDependencyUpdate(client, projects)
	zlog.Info("Finished getting latest release and making updates to made")

	zlog.Info("Start updating the repositories' dependencies...")
	updater.UpdateProjectDependencies(dependencyUpdates)
	zlog.Info("Finished updating the repositories...")

	zlog.Info("Job finished.")
}
