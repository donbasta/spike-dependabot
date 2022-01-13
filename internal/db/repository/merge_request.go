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

		GetByID(tx *gorm.DB, id uuid.UUID) (mergeRequest entity.MergeRequest, err error)
		GetAll(tx *gorm.DB) (mergeRequests []entity.MergeRequest, err error)
		GetAllByRepositoryId(tx *gorm.DB, repositoryID string) (mergeRequests []entity.MergeRequest, err error)
	}

	mergeRequestRepository struct{}
)

func NewMergeRequestRepository() MergeRequestRepository {
	return &mergeRequestRepository{}
}

func (r *mergeRequestRepository) GetByID(tx *gorm.DB, id uuid.UUID) (mergeRequest entity.MergeRequest, err error) {
	err = tx.First(&mergeRequest, id).Error
	return
}

func (r *mergeRequestRepository) GetAllByRepositoryId(tx *gorm.DB, repositoryID string) (mergeRequests []entity.MergeRequest, err error) {
	err = tx.Where("repository_id = ?", repositoryID).Find(&mergeRequests).Error
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
