package dependabot

import (
	"dependabot/internal/config"
	"dependabot/internal/db"
	"dependabot/internal/db/entity"
	"log"

	"github.com/spf13/cobra"
)

func newAutoMigrate() *cobra.Command {
	var command = &cobra.Command{
		Use:   "automigrate",
		Short: "Run gorm auto migration",
		Long:  "Run gorm auto migration",
		Run: func(c *cobra.Command, args []string) {
			cfg := config.Config()
			if cfg == nil {
				log.Fatal("parsing config failed")
				return
			}

			db, err := db.ProvideDB(cfg)
			if err != nil {
				log.Fatal("initialize db failed: ", err)
				return
			}

			err = db.AutoMigrate(
				&entity.MergeRequest{},
			)
			if err != nil {
				log.Fatal("error while migration with error: ", err)
			}
		},
	}

	return command
}
