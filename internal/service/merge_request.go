package service

import (
	"dependabot/internal/db/entity"
	"dependabot/internal/db/repository"
	"dependabot/internal/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MergeRequestService interface {
	GetAll() (mergeRequests []entity.MergeRequest, err error)
	GetAllByRepositoryId(repositoryId string) (mergeRequests *[]entity.MergeRequest, err error)
	Create(mergeRequest *entity.MergeRequest) (uuid.UUID, error)
	Delete(mergeRequest *entity.MergeRequest) error
}

type mergeRequestService struct {
	db         *gorm.DB
	repository repository.MergeRequestRepository
}

func (m *mergeRequestService) GetAll() (mergeRequests []entity.MergeRequest, err error) {
	mergeRequests, err = m.repository.GetAll(m.db)
	if err != nil {
		return nil, errors.NewDatabaseError(err, "Get all merge requests")
	}
	return
}

func (m *mergeRequestService) GetAllByRepositoryId(repositoryId string) (mergeRequests *[]entity.MergeRequest, err error) {
	existingMergeRequests, err := m.repository.GetAllByRepositoryId(m.db, repositoryId)
	if err != nil {
		return nil, errors.NewDatabaseError(err, "Get merge requests by repository ID %s", repositoryId)
	}
	return &existingMergeRequests, nil
}

func (m *mergeRequestService) Create(mergeRequest *entity.MergeRequest) (uuid.UUID, error) {
	id, err := m.repository.Create(m.db, mergeRequest)
	if err != nil {
		return uuid.Nil, errors.NewDatabaseError(err, "create job error")
	}
	return id, nil
}

func (m *mergeRequestService) Delete(mergeRequest *entity.MergeRequest) error {
	err := m.repository.Delete(m.db, mergeRequest)
	if err != nil {
		return errors.NewDatabaseError(err, "delete merge request %s", mergeRequest.ID)
	}
	return nil
}

func NewMergeRequestService(
	db *gorm.DB,
	repository repository.MergeRequestRepository,
) MergeRequestService {
	return &mergeRequestService{
		db:         db,
		repository: repository,
	}
}
