package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UserID    *uuid.UUID     `gorm:"type:uuid;index" json:"user_id,omitempty"`
	Responses []Response     `gorm:"foreignKey:SessionID" json:"responses,omitempty"`
	Result    *Result        `gorm:"foreignKey:SessionID" json:"result,omitempty"`
}

func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
