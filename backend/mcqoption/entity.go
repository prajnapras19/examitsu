package mcqoption

import "github.com/prajnapras19/project-form-exam-sman2/backend/lib"

type McqOption struct {
	lib.BaseModel
	QuestionID  uint
	Description string
	Point       int
}
