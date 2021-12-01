package dependabot

import (
	"dependabot/internal/client"
	"dependabot/internal/config"
	"dependabot/internal/service"
	"log"
)

func Run() {
	mainConfig := config.ProvideConfig()

	// AutoMigrate()

	migration, err := InitializeMigration()
	if err != nil {
		log.Fatal("initializing db failed with error: ", err)
	}
	err = migration.Up()
	if err != nil {
		log.Fatal("Migrate up failing with error: ", err)
	}
	log.Println("migration finished")

	client := client.CreateClient(mainConfig.Git.Token)

	projects, err := service.CrawlGroups(client)
	if err != nil {
		log.Fatal(err)
	}
	service.CheckDependency(client, projects)
}
