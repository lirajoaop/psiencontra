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

// FindCompletedByUserID returns the user's sessions that have a Result
// attached, newest first. Sessions without a Result (abandoned mid-flow or
// where AI analysis failed) are filtered out at the SQL level via INNER JOIN
// so the history screen only lists consultable entries.
func (r *SessionRepo) FindCompletedByUserID(userID uuid.UUID) ([]schemas.Session, error) {
	var sessions []schemas.Session
	err := r.DB.
		Preload("Result").
		Joins("INNER JOIN results ON results.session_id = sessions.id AND results.deleted_at IS NULL").
		Where("sessions.user_id = ?", userID).
		Order("sessions.created_at DESC").
		Find(&sessions).Error
	return sessions, err
}
