package exam

import "github.com/prajnapras19/project-form-exam-sman2/backend/lib"

type Exam struct {
	lib.BaseModel
	Serial string
	Name   string
	IsOpen bool
}
