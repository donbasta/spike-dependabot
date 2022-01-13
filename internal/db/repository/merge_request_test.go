package repository_test

import (
	"dependabot/internal/config"
	"dependabot/internal/db"
	"dependabot/internal/db/entity"
	. "dependabot/internal/db/repository"
	"testing"

	"github.com/gopaytech/go-commons/pkg/db/migration"
	pgMigrate "github.com/gopaytech/go-commons/pkg/db/migration/postgresql"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type testMergeRequestContext struct {
	tx       *gorm.DB
	migrator migration.Migration
}

func (o *testMergeRequestContext) SetUp(t *testing.T) {
	main := config.ProvideConfig()

	idb, err := db.ProvideDB(main)
	assert.NoError(t, err)
	assert.NotNil(t, idb)

	sqlDb, err := idb.DB()
	assert.NoError(t, err)
	migrator, err := pgMigrate.WithInstanceMigrationTable(sqlDb, "file://./seed", "repository_test_migration")
	assert.NoError(t, err)

	o.tx = idb
	o.migrator = migrator

	err = migrator.Up()
	assert.NoError(t, err)
}

func (o *testMergeRequestContext) TearDown(t *testing.T) {
	err := o.migrator.Down()
	assert.NoError(t, err)
}

func TestMergeRequestRepository(t *testing.T) {
	ctx := &testMergeRequestContext{}
	ctx.SetUp(t)

	defer ctx.TearDown(t)

	repository := NewMergeRequestRepository()

	mergeRequest := &entity.MergeRequest{
		MergeRequestIID: "1",
		RepositoryURL:   faker.URL(),
		RepositoryID:    faker.Name(),
	}

	var mergeRequestID uuid.UUID
	t.Run("create", func(p *testing.T) {
		id, err := repository.Create(ctx.tx, mergeRequest)
		assert.Nil(p, err)
		assert.NotEqual(p, uuid.Nil, id)
		assert.NotNil(p, mergeRequest.CreatedAt)
		mergeRequestID = mergeRequest.ID
	})

	t.Run("get by id", func(t *testing.T) {
		existingMergeRequest, err := repository.GetByID(ctx.tx, mergeRequestID)
		assert.Nil(t, err)
		assert.NotNil(t, existingMergeRequest)
	})

	t.Run("get all by repository id", func(t *testing.T) {
		existingMergeRequest, err := repository.GetAllByRepositoryId(ctx.tx, "31999")
		assert.Nil(t, err)
		assert.NotNil(t, existingMergeRequest)
		assert.Equal(t, 2, len(existingMergeRequest))
	})

	t.Run("update", func(t *testing.T) {
		updatedURL := faker.URL()
		mergeRequest.RepositoryURL = updatedURL
		err := repository.Update(ctx.tx, mergeRequest)
		assert.Nil(t, err)

		existingMergeRequest, err := repository.GetByID(ctx.tx, mergeRequest.ID)
		assert.Nil(t, err)
		assert.Equal(t, updatedURL, existingMergeRequest.RepositoryURL)
	})

	t.Run("delete", func(t *testing.T) {
		existingMergeRequest, err := repository.GetByID(ctx.tx, mergeRequestID)
		assert.Nil(t, err)

		err = repository.Delete(ctx.tx, &existingMergeRequest)
		assert.Nil(t, err)

		_, err = repository.GetByID(ctx.tx, mergeRequestID)
		assert.NotNil(t, err)
	})

	t.Run("get all", func(t *testing.T) {
		mergeRequests, err := repository.GetAll(ctx.tx)
		assert.Nil(t, err)
		assert.NotEmpty(t, mergeRequests)
		assert.True(t, len(mergeRequests) >= 3)
	})
}
