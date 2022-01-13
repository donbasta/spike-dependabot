package dependabot

import (
	"dependabot/internal/config"
	"dependabot/internal/db"
	"fmt"

	"github.com/gopaytech/go-commons/pkg/db/migration"
	"github.com/gopaytech/go-commons/pkg/zlog"
	"github.com/spf13/cobra"
)

func InitializeMigration() (migration.Migration, error) {
	main := config.ProvideConfig()

	gormDB, err := db.ProvideDB(main)
	if err != nil {
		return nil, err
	}

	migration, err := db.ProvideMigration(gormDB)
	if err != nil {
		return nil, err
	}

	return migration, nil
}

func newMigrate() *cobra.Command {
	command := &cobra.Command{
		Use:     "migrate",
		Short:   "Run db migration",
		Long:    "Run db migration",
		Aliases: []string{"m"},
		RunE: func(c *cobra.Command, args []string) error {
			cfg := config.Config()
			if cfg == nil {
				return fmt.Errorf("parsing config failed %v", cfg)
			}

			migration, err := InitializeMigration()
			if err != nil {
				return err
			}
			err = migration.Up()
			if err != nil {
				zlog.Error(err, "migrate up fail")
			}
			zlog.Info("migration finished")
			return nil
		},
	}

	return command
}
