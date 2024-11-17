package exam

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
	CreateExam(exam *Exam) (*Exam, error)
	GetExamBySerial(serial string) (*Exam, error)
	GetExams(pagination *lib.QueryPagination, filter *GetExamsFilter) ([]*Exam, error)
	UpdateExam(exam *Exam) error
	DeleteExamBySerial(serial string) error
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

func (r *repository) CreateExam(exam *Exam) (*Exam, error) {
	err := r.db.Create(exam).Error
	return exam, err
}

func (r *repository) GetExamBySerial(serial string) (*Exam, error) {
	var exam Exam

	cacheKey := r.GetExamBySerialCacheKey(serial)
	val, err := r.cache.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &exam)
		return &exam, nil
	}

	err = r.db.Where("serial = ?", serial).First(&exam).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, lib.ErrExamNotFound
		}
		return nil, err
	}
	res, _ := json.Marshal(exam)
	r.cache.Set(context.Background(), cacheKey, res, r.cfg.CacheTTL)
	return &exam, nil
}

func (r *repository) GetExams(pagination *lib.QueryPagination, filter *GetExamsFilter) ([]*Exam, error) {
	var res []*Exam
	err := r.db.Scopes(append(filter.Scope(), pagination.Scope())...).Find(&res).Error
	return res, err
}

func (r *repository) UpdateExam(exam *Exam) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Exam{}).
			Where("serial = ?", exam.Serial).
			Update("name", exam.Name).
			Error; err != nil {
			return err
		}

		if err := tx.Model(&Exam{}).
			Where("serial = ?", exam.Serial).
			Update("is_open", exam.IsOpen).
			Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	cacheKey := r.GetExamBySerialCacheKey(exam.Serial)
	r.cache.Del(context.Background(), cacheKey)
	return nil
}

func (r *repository) DeleteExamBySerial(serial string) error {
	res := r.db.Model(&Exam{}).Where("serial = ?", serial).Delete(&Exam{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		log.Printf("[exam][repository][DeleteExamBySerial] error: %s", res.Error)
		return lib.ErrExamNotFound
	}
	cacheKey := r.GetExamBySerialCacheKey(serial)
	r.cache.Del(context.Background(), cacheKey)
	return nil
}

func (r *repository) GetExamBySerialCacheKey(serial string) string {
	return fmt.Sprintf("exam:%s", serial)
}
