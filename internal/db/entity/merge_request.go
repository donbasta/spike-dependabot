package entity

import (
	"dependabot/internal/db"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MergeRequest struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;not null;primaryKey"`
	MergeRequestIID string    `json:"merge_requests_id" gorm:"type:varchar(255);not null"`
	RepositoryURL   string    `json:"repository_url" gorm:"type:varchar(255);not null"`
	RepositoryID    string    `json:"repository_id" gorm:"type:varchar(255);not null"`
	db.Timestamp
}

func (m *MergeRequest) BeforeCreate(_ *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return
}
