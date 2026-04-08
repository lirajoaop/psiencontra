package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Response struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	SessionID    uuid.UUID      `gorm:"type:uuid;index" json:"session_id"`
	QuestionID   int            `json:"question_id"`
	QuestionText string         `json:"question_text"`
	AnswerType   string         `json:"answer_type"` // multiple_choice | open_ended
	AnswerValue  string         `json:"answer_value"`
}
