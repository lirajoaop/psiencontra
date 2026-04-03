package repository

import (
	"github.com/google/uuid"
	"github.com/joaop/psiencontra/api/schemas"
	"gorm.io/gorm"
)

type ResponseRepo struct {
	DB *gorm.DB
}

func NewResponseRepo(db *gorm.DB) *ResponseRepo {
	return &ResponseRepo{DB: db}
}

func (r *ResponseRepo) CreateBatch(responses []schemas.Response) error {
	return r.DB.Create(&responses).Error
}

func (r *ResponseRepo) FindBySessionID(sessionID uuid.UUID) ([]schemas.Response, error) {
	var responses []schemas.Response
	err := r.DB.Where("session_id = ?", sessionID).Order("question_id").Find(&responses).Error
	return responses, err
}
