package question

import (
	"log"

	"github.com/prajnapras19/project-form-exam-sman2/backend/constants"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
	"gorm.io/gorm"
)

type Repository interface {
	CreateQuestion(question *Question) (*Question, error)
	GetQuestions(pagination *lib.QueryPagination, filter *GetQuestionsFilter) ([]*Question, error)
	UpdateQuestionDataByID(question *Question) error
	DeleteQuestionByID(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateQuestion(question *Question) (*Question, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(question).Error; err != nil {
			return err
		}
		if err := tx.Model(question).Where("id = ?", question.ID).Update(constants.OrderNumber, question.ID).Error; err != nil {
			return err
		}
		return nil
	})
	return question, err
}

func (r *repository) GetQuestions(pagination *lib.QueryPagination, filter *GetQuestionsFilter) ([]*Question, error) {
	var res []*Question
	err := r.db.Scopes(append(filter.Scope(), pagination.Scope())...).Find(&res).Error
	return res, err
}

func (r *repository) UpdateQuestionDataByID(question *Question) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&Question{}).
			Where("id = ?", question.ID).
			Update("data", question.Data).
			Error
	})
}

func (r *repository) DeleteQuestionByID(id uint) error {
	res := r.db.Model(&Question{}).Where("id = ?", id).Delete(&Question{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		log.Printf("[question][repository][DeleteQuestionByID] error: %s", res.Error)
		return lib.ErrQuestionNotFound
	}
	return nil
}
