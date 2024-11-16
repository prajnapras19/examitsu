package exam

import (
	"github.com/prajnapras19/project-form-exam-sman2/backend/constants"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
	"gorm.io/gorm"
)

type Exam struct {
	lib.BaseModel
	Serial string
	Name   string
	IsOpen bool
}

type GetExamsFilter struct {
	SerialEqualsTo *lib.QueryFiltersEqualToString `json:"serial_equals_to"`
}

func (f *GetExamsFilter) Scope() []func(db *gorm.DB) *gorm.DB {
	scopes := []func(db *gorm.DB) *gorm.DB{}

	if f.SerialEqualsTo != nil {
		scopes = append(scopes, f.SerialEqualsTo.Scope(constants.Serial))
	}

	return scopes
}
