package participant

import (
	"time"

	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
)

type Participant struct {
	lib.BaseModel
	ExamID                 uint
	Name                   string
	Password               string // not used anymore, but not dropped for backward compatibility with v.0.0
	AllowedDurationMinutes uint
	StartedAt              *time.Time
	EndedAt                *time.Time
}

type ParticipantTotalPoint struct {
	ParticipantID uint
	TotalPoint    int
}
