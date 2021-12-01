package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"dependabot/internal/db/entity"
)

type (
	MergeRequestRepository interface {
		Create(tx *gorm.DB, mergeRequest *entity.MergeRequest) (uuid.UUID, error)
		Update(tx *gorm.DB, mergeRequest *entity.MergeRequest) error
		Delete(tx *gorm.DB, mergeRequest *entity.MergeRequest) (err error)

		GetAll(tx *gorm.DB) (mergeRequests []entity.MergeRequest, err error)
		GetAllMergeRequestByRepositoryID(tx *gorm.DB, repositoryID string) (mergeRequests []entity.MergeRequest, err error)
		GetLastMergeRequestByRepositoryID(tx *gorm.DB, repositoryID string) (mergeRequest entity.MergeRequest, err error)
	}

	mergeRequestRepository struct{}
)

func NewMergeRequestRepository() MergeRequestRepository {
	return &mergeRequestRepository{}
}

func (r *mergeRequestRepository) GetAllMergeRequestByRepositoryID(tx *gorm.DB, repositoryID string) (mergeRequests []entity.MergeRequest, err error) {
	err = tx.Where("repository_id = ?", repositoryID).Find(&mergeRequests).Error
	return
}

func (r *mergeRequestRepository) GetLastMergeRequestByRepositoryID(tx *gorm.DB, repositoryID string) (mergeRequest entity.MergeRequest, err error) {
	err = tx.Where("repository_id = ?", repositoryID).Order("created_at desc").Limit(1).Find(&mergeRequest).Error
	return
}

func (r *mergeRequestRepository) Create(tx *gorm.DB, mergeRequest *entity.MergeRequest) (uuid.UUID, error) {
	result := tx.Create(&mergeRequest)
	return mergeRequest.ID, result.Error
}

func (r *mergeRequestRepository) GetAll(tx *gorm.DB) (mergeRequests []entity.MergeRequest, err error) {
	err = tx.Find(&mergeRequests).Error
	return
}

func (r *mergeRequestRepository) Update(tx *gorm.DB, mergeRequest *entity.MergeRequest) error {
	err := tx.Save(&mergeRequest).Error
	return err
}

func (r *mergeRequestRepository) Delete(tx *gorm.DB, mergeRequest *entity.MergeRequest) (err error) {
	err = tx.Delete(&mergeRequest).Error
	return
}
