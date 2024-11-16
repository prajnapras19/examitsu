package mcqoption

import (
	"log"

	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
	"gorm.io/gorm"
)

type Repository interface {
	CreateMcqOption(mcqOption *McqOption) (*McqOption, error)
	GetMcqOptionsByQuestionID(questionID uint) ([]*McqOption, error)
	UpdateMcqOption(mcqOption *McqOption) error
	DeleteMcqOptionByID(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateMcqOption(mcqOption *McqOption) (*McqOption, error) {
	err := r.db.Create(mcqOption).Error
	return mcqOption, err
}

func (r *repository) GetMcqOptionsByQuestionID(questionID uint) ([]*McqOption, error) {
	var res []*McqOption
	err := r.db.Where("question_id = ?", questionID).Find(&res).Error
	return res, err
}

func (r *repository) UpdateMcqOption(mcqOption *McqOption) error {
	res := r.db.Updates(mcqOption)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		log.Printf("[mcqoption][repository][UpdateMcqOption] error: %s", res.Error)
		return lib.ErrMcqOptionNotFound
	}
	return nil
}

func (r *repository) DeleteMcqOptionByID(id uint) error {
	res := r.db.Model(&McqOption{}).Where("id = ?", id).Delete(&McqOption{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		log.Printf("[mcqoption][repository][DeleteMcqOptionByID] error: %s", res.Error)
		return lib.ErrMcqOptionNotFound
	}
	return nil
}
