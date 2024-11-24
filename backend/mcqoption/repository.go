package mcqoption

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/prajnapras19/project-form-exam-sman2/backend/config"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
	redis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository interface {
	CreateMcqOption(mcqOption *McqOption) (*McqOption, error)
	GetMcqOptionsByQuestionID(questionID uint) ([]*McqOption, error)
	GetMcqOptionByID(id uint) (*McqOption, error)
	UpdateMcqOption(mcqOption *McqOption) error
	DeleteMcqOptionByID(id uint) error
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

func (r *repository) CreateMcqOption(mcqOption *McqOption) (*McqOption, error) {
	err := r.db.Create(mcqOption).Error
	if err != nil {
		r.cache.Del(context.Background(), r.GetMcqOptionByQuestionIDCacheKey(mcqOption.QuestionID))
	}
	return mcqOption, err
}

func (r *repository) GetMcqOptionsByQuestionID(questionID uint) ([]*McqOption, error) {
	var mcqOptions []*McqOption

	cacheKey := r.GetMcqOptionByQuestionIDCacheKey(questionID)
	val, err := r.cache.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &mcqOptions)
		return mcqOptions, nil
	}

	err = r.db.Where("question_id = ?", questionID).Find(&mcqOptions).Error
	if err != nil {
		return nil, err
	}

	res, _ := json.Marshal(mcqOptions)
	r.cache.Set(context.Background(), cacheKey, res, r.cfg.CacheTTL)
	return mcqOptions, err
}

func (r *repository) GetMcqOptionByID(id uint) (*McqOption, error) {
	var mcqOption McqOption

	cacheKey := r.GetMcqOptionByIDCacheKey(id)
	val, err := r.cache.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &mcqOption)
		return &mcqOption, nil
	}

	err = r.db.Where("id = ?", id).First(&mcqOption).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, lib.ErrMcqOptionNotFound
		}
		return nil, err
	}

	res, _ := json.Marshal(mcqOption)
	r.cache.Set(context.Background(), cacheKey, res, r.cfg.CacheTTL)
	return &mcqOption, nil
}

func (r *repository) UpdateMcqOption(mcqOption *McqOption) error {
	currentData, err := r.GetMcqOptionByID(mcqOption.ID)
	if err != nil {
		return err
	}

	res := r.db.Updates(mcqOption)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		log.Printf("[mcqoption][repository][UpdateMcqOption] error: %s", res.Error)
		return lib.ErrMcqOptionNotFound
	}

	r.cache.Del(context.Background(), r.GetMcqOptionByIDCacheKey(currentData.ID))
	r.cache.Del(context.Background(), r.GetMcqOptionByQuestionIDCacheKey(currentData.QuestionID))

	return nil
}

func (r *repository) DeleteMcqOptionByID(id uint) error {
	currentData, err := r.GetMcqOptionByID(id)
	if err != nil {
		return err
	}

	res := r.db.Model(&McqOption{}).Where("id = ?", id).Delete(&McqOption{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		log.Printf("[mcqoption][repository][DeleteMcqOptionByID] error: %s", res.Error)
		return lib.ErrMcqOptionNotFound
	}

	r.cache.Del(context.Background(), r.GetMcqOptionByIDCacheKey(currentData.ID))
	r.cache.Del(context.Background(), r.GetMcqOptionByQuestionIDCacheKey(currentData.QuestionID))

	return nil
}

func (r *repository) GetMcqOptionByIDCacheKey(id uint) string {
	return fmt.Sprintf("mcq_option_list:id:%d", id)
}

func (r *repository) GetMcqOptionByQuestionIDCacheKey(questionID uint) string {
	return fmt.Sprintf("mcq_option_list:questionID:%d", questionID)
}
