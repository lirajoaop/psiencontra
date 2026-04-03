package schemas

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Response struct {
	gorm.Model
	SessionID    uuid.UUID `gorm:"type:uuid;index" json:"session_id"`
	QuestionID   int       `json:"question_id"`
	QuestionText string    `json:"question_text"`
	AnswerType   string    `json:"answer_type"` // multiple_choice | open_ended
	AnswerValue  string    `json:"answer_value"`
}
