package submission

import (
	"fmt"
	"time"

	"github.com/prajnapras19/project-form-exam-sman2/backend/constants"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
)

type Submission struct {
	lib.BaseModel

	ParticipantID uint
	QuestionID    uint
	McqOptionID   uint
}

type ExamSessionSubmissionCacheObject struct {
	ParticipantID uint
	QuestionID    uint
	McqOptionID   uint
	Timestamp     time.Time
}

func (e *ExamSessionSubmissionCacheObject) GetKey() string {
	return fmt.Sprintf("%s:%d:%d", constants.ExamSessionSubmissionCacheObjectKeyPrefix, e.ParticipantID, e.QuestionID)
}
