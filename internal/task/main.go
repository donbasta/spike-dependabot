package task

import (
	"dependabot/internal/config"
	"log"
)

func Run(cfg *config.Main) {
	client, err := config.ProvideGitlabClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	projects, err := CrawlGroups(client)
	if err != nil {
		log.Fatal(err)
	}
	changes := CheckDependency(client, projects)
	UpdateProjects(changes)
}
