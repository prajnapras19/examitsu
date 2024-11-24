package question

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/prajnapras19/project-form-exam-sman2/backend/config"
	"github.com/prajnapras19/project-form-exam-sman2/backend/constants"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
	redis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository interface {
	CreateQuestion(question *Question) (*Question, error)
	GetQuestionsIDOnly(pagination *lib.QueryPagination, filter *GetQuestionsFilter) ([]*Question, error)
	GetQuestionByID(id uint) (*Question, error)
	GetQuestions(pagination *lib.QueryPagination, filter *GetQuestionsFilter) ([]*Question, error)
	GetQuestionsIDByExamID(examID uint) ([]*Question, error)
	UpdateQuestionDataByID(question *Question) error
	DeleteQuestionByID(id uint) error
}

type repository struct {
	cfg   *config.Config
	db    *gorm.DB
	cache *redis.Client
}

func NewRepository(
	cfg *config.Config,
	db *gorm.DB,
	cache *redis.Client,
) Repository {
	return &repository{
		cfg:   cfg,
		db:    db,
		cache: cache,
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
		r.cache.Del(context.Background(), r.GetQuestionsIDByExamIDCacheKey(question.ExamID))
		return nil
	})
	return question, err
}

func (r *repository) GetQuestionByID(id uint) (*Question, error) {
	var question Question

	cacheKey := r.GetQuestionByIDCacheKey(id)
	val, err := r.cache.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &question)
		return &question, nil
	}

	err = r.db.Where("id = ?", id).First(&question).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, lib.ErrQuestionNotFound
		}
		return nil, err
	}

	res, _ := json.Marshal(question)
	r.cache.Set(context.Background(), cacheKey, res, r.cfg.CacheTTL)
	return &question, nil
}

func (r *repository) GetQuestionsIDOnly(pagination *lib.QueryPagination, filter *GetQuestionsFilter) ([]*Question, error) {
	var res []*Question
	err := r.db.Select("id").Scopes(append(filter.Scope(), pagination.Scope())...).Find(&res).Error
	return res, err
}

func (r *repository) GetQuestionsIDByExamID(examID uint) ([]*Question, error) {
	var questions []*Question

	cacheKey := r.GetQuestionsIDByExamIDCacheKey(examID)
	val, err := r.cache.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &questions)
		return questions, nil
	}

	err = r.db.Select("id").Order("order_number ASC").Where("exam_id = ?", examID).Find(&questions).Error
	if err != nil {
		return nil, err
	}

	res, _ := json.Marshal(questions)
	r.cache.Set(context.Background(), cacheKey, res, r.cfg.CacheTTL)
	return questions, err
}

func (r *repository) GetQuestions(pagination *lib.QueryPagination, filter *GetQuestionsFilter) ([]*Question, error) {
	var res []*Question
	err := r.db.Scopes(append(filter.Scope(), pagination.Scope())...).Find(&res).Error
	return res, err
}

func (r *repository) UpdateQuestionDataByID(question *Question) error {
	currentData, err := r.GetQuestionByID(question.ID)
	if err != nil {
		return err
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&Question{}).
			Where("id = ?", question.ID).
			Update("data", question.Data).
			Error
	})

	if err == nil {
		r.cache.Del(context.Background(), r.GetQuestionByIDCacheKey(currentData.ID))
		r.cache.Del(context.Background(), r.GetQuestionsIDByExamIDCacheKey(currentData.ExamID))
	}
	return err
}

func (r *repository) DeleteQuestionByID(id uint) error {
	currentData, err := r.GetQuestionByID(id)
	if err != nil {
		return err
	}

	res := r.db.Model(&Question{}).Where("id = ?", id).Delete(&Question{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		log.Printf("[question][repository][DeleteQuestionByID] error: %s", res.Error)
		return lib.ErrQuestionNotFound
	}

	r.cache.Del(context.Background(), r.GetQuestionByIDCacheKey(currentData.ID))
	r.cache.Del(context.Background(), r.GetQuestionsIDByExamIDCacheKey(currentData.ExamID))

	return nil
}

func (r *repository) GetQuestionByIDCacheKey(id uint) string {
	return fmt.Sprintf("question:id:%d", id)
}

func (r *repository) GetQuestionsIDByExamIDCacheKey(examID uint) string {
	return fmt.Sprintf("question_id_list:examID:%d", examID)
}
