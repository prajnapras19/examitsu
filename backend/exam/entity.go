package exam

import (
	"time"

	"github.com/prajnapras19/project-form-exam-sman2/backend/constants"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
	"gorm.io/gorm"
)

type Exam struct {
	lib.BaseModel
	Serial                 string
	Name                   string
	IsOpen                 bool
	AllowedDurationMinutes uint
}

type GetExamsFilter struct {
	SerialEqualsTo *lib.QueryFiltersEqualToString `json:"serial_equals_to"`
	IsOpenEqualsTo *lib.QueryFiltersEqualBool     `json:"is_open_equals_to"`
}

func (f *GetExamsFilter) Scope() []func(db *gorm.DB) *gorm.DB {
	scopes := []func(db *gorm.DB) *gorm.DB{}

	if f.SerialEqualsTo != nil {
		scopes = append(scopes, f.SerialEqualsTo.Scope(constants.Serial))
	}

	if f.IsOpenEqualsTo != nil {
		scopes = append(scopes, f.IsOpenEqualsTo.Scope(constants.IsOpen))
	}

	return scopes
}

// copy of participant.Participant
type Participant struct {
	lib.BaseModel
	ExamID                 uint
	Name                   string
	Password               string // not used anymore, but not dropped for backward compatibility with v.0.0
	AllowedDurationMinutes uint
	StartedAt              *time.Time
	EndedAt                *time.Time
}
