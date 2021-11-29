package db

import (
	"dependabot/internal/config"
	"sync"

	"github.com/gopaytech/go-commons/pkg/db"
	"github.com/gopaytech/go-commons/pkg/db/gorm/postgresql"
	"github.com/gopaytech/go-commons/pkg/db/migration"
	pgMigrate "github.com/gopaytech/go-commons/pkg/db/migration/postgresql"
	"github.com/spf13/cast"

	"gorm.io/gorm"
)

var gormDb *gorm.DB
var once sync.Once

func dbConfig(config *config.Main) db.Config {
	return db.Config{
		Username:     config.Database.User,
		Password:     config.Database.Password,
		Host:         config.Database.Host,
		Port:         cast.ToInt(config.Database.Port),
		DatabaseName: config.Database.Name,
	}
}

func ProvideDB(config *config.Main) (*gorm.DB, error) {
	dbConf := dbConfig(config)
	var errObj error

	once.Do(func() {
		g, err := postgresql.ConnectDefault(dbConf)
		if err != nil {
			errObj = err
		}
		gormDb = g
	})

	return gormDb, errObj
}

func ProvideMigration(gdb *gorm.DB) (migration.Migration, error) {
	sqlDb, err := gdb.DB()
	if err != nil {
		return nil, err
	}
	return pgMigrate.WithInstance(sqlDb, "file://./migrations")
}
