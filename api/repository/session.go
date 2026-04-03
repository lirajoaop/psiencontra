package repository

import (
	"github.com/google/uuid"
	"github.com/joaop/psiencontra/api/schemas"
	"gorm.io/gorm"
)

type SessionRepo struct {
	DB *gorm.DB
}

func NewSessionRepo(db *gorm.DB) *SessionRepo {
	return &SessionRepo{DB: db}
}

func (r *SessionRepo) Create(s *schemas.Session) error {
	return r.DB.Create(s).Error
}

func (r *SessionRepo) FindByID(id uuid.UUID) (*schemas.Session, error) {
	var s schemas.Session
	err := r.DB.First(&s, "id = ?", id).Error
	return &s, err
}
