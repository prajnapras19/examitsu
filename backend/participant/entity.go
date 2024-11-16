package participant

import (
	"time"

	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
)

type Participant struct {
	lib.BaseModel
	ExamID    uint
	Name      string
	Password  string
	StartedAt *time.Time
	EndedAt   *time.Time
}
