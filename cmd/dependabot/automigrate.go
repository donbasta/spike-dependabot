package dependabot

import (
	"dependabot/internal/config"
	"dependabot/internal/db"
	"dependabot/internal/db/entity"
	"log"
)

func AutoMigrate() {
	main := config.ProvideConfig()
	db, err := db.ProvideDB(main)
	if err != nil {
		log.Fatal("initialize db failed: ", err)
	}
	err = db.AutoMigrate(
		&entity.MergeRequest{},
	)
	if err != nil {
		log.Fatal("error while migration with error: ", err)
	}
}
