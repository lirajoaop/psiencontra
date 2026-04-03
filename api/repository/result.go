package repository

import (
	"github.com/google/uuid"
	"github.com/joaop/psiencontra/api/schemas"
	"gorm.io/gorm"
)

type ResultRepo struct {
	DB *gorm.DB
}

func NewResultRepo(db *gorm.DB) *ResultRepo {
	return &ResultRepo{DB: db}
}

func (r *ResultRepo) Create(result *schemas.Result) error {
	return r.DB.Create(result).Error
}

func (r *ResultRepo) FindBySessionID(sessionID uuid.UUID) (*schemas.Result, error) {
	var result schemas.Result
	err := r.DB.Where("session_id = ?", sessionID).First(&result).Error
	return &result, err
}
