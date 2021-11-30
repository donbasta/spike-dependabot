package dependabot

import (
	"dependabot/internal/client"
	"dependabot/internal/config"
	"dependabot/internal/db"
	"dependabot/internal/service"
	"log"

	entity "dependabot/internal/db/entity"
)

func Run() {
	mainConfig := config.ProvideConfig()

	db, err := db.ProvideDB(mainConfig)
	if err != nil {
		log.Fatal("initializing db failed with error: ", err)
	}
	log.Println("connected with db")

	err = db.AutoMigrate(
		&entity.MergeRequest{},
	)
	if err != nil {
		log.Fatal("error while migration with error: ", err)
	}

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

	service.CrawlGroups(client)
}
