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
	GetExamByID(id uint) (*Exam, error)
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

func (r *repository) GetExamByID(id uint) (*Exam, error) {
	var exam Exam

	cacheKey := r.GetExamByIDCacheKey(id)
	val, err := r.cache.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &exam)
		return &exam, nil
	}

	err = r.db.Where("id = ?", id).First(&exam).Error
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

func (r *repository) GetExamBySerial(serial string) (*Exam, error) {
	var exam Exam

	cacheKey := r.GetExamBySerialCacheKey(serial)
	val, err := r.cache.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &exam)
		return &exam, nil
	}

	err = r.db.Where("serial = ? AND not_archived", serial).First(&exam).Error
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
	currentData, err := r.GetExamBySerial(exam.Serial)
	if err != nil {
		return err
	}

	r.cache.Del(context.Background(), r.GetExamByIDCacheKey(currentData.ID))
	r.cache.Del(context.Background(), r.GetExamBySerialCacheKey(currentData.Serial))

	return r.db.Transaction(func(tx *gorm.DB) error {
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
}

func (r *repository) DeleteExamBySerial(serial string) error {
	currentData, err := r.GetExamBySerial(serial)
	if err != nil {
		return err
	}

	r.cache.Del(context.Background(), r.GetExamByIDCacheKey(currentData.ID))
	r.cache.Del(context.Background(), r.GetExamBySerialCacheKey(currentData.Serial))

	res := r.db.Model(&Exam{}).Where("serial = ?", serial).Delete(&Exam{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		log.Printf("[exam][repository][DeleteExamBySerial] error: %s", res.Error)
		return lib.ErrExamNotFound
	}
	return nil
}

func (r *repository) GetExamBySerialCacheKey(serial string) string {
	return fmt.Sprintf("exam:serial:%s", serial)
}

func (r *repository) GetExamByIDCacheKey(id uint) string {
	return fmt.Sprintf("exam:id:%d", id)
}
