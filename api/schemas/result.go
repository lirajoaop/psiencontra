package schemas

import (
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Result struct {
	gorm.Model
	SessionID       uuid.UUID       `gorm:"type:uuid;uniqueIndex" json:"session_id"`
	ApproachScores  json.RawMessage `gorm:"type:jsonb" json:"approach_scores"`
	FieldScores     json.RawMessage `gorm:"type:jsonb" json:"field_scores"`
	Explanation     string          `gorm:"type:text" json:"explanation"`
	ApproachDetails json.RawMessage `gorm:"type:jsonb" json:"approach_details"`
	FieldDetails    json.RawMessage `gorm:"type:jsonb" json:"field_details"`
	AIProvider      string          `json:"ai_provider"` // gemini | groq
}
