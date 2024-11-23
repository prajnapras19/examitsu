package participantsession

import "github.com/prajnapras19/project-form-exam-sman2/backend/lib"

type ParticipantSession struct {
	lib.BaseModel

	Serial        string
	ParticipantID uint
	IsAuthorized  bool
}
