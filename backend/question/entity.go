package question

import (
	"github.com/prajnapras19/project-form-exam-sman2/backend/constants"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
	"gorm.io/gorm"
)

type Question struct {
	lib.BaseModel
	ExamID      uint
	OrderNumber uint
	Data        string
}

type GetQuestionsFilter struct {
	IDEqualsTo         *lib.QueryFiltersEqualToUint   `json:"id_equals_to"`
	ExamIDEqualsTo     *lib.QueryFiltersEqualToUint   `json:"-"`
	ExamSerialEqualsTo *lib.QueryFiltersEqualToString `json:"exam_serial_equals_to"`
}

func (f *GetQuestionsFilter) Scope() []func(db *gorm.DB) *gorm.DB {
	scopes := []func(db *gorm.DB) *gorm.DB{}

	if f.IDEqualsTo != nil {
		scopes = append(scopes, f.IDEqualsTo.Scope(constants.ID))
	}
	if f.ExamIDEqualsTo != nil {
		scopes = append(scopes, f.ExamIDEqualsTo.Scope(constants.ExamID))
	}

	return scopes
}
