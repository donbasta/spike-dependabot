package dependabot

import (
	"dependabot/internal/config"
	"dependabot/internal/db"

	"github.com/gopaytech/go-commons/pkg/db/migration"
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
