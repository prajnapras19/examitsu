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
	GetAllOpenedExams() ([]*Exam, error)
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
	if err == nil {
		r.cache.Del(context.Background(), r.GetAllOpenedExamsCacheKey())
	}
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

func (r *repository) GetAllOpenedExams() ([]*Exam, error) {
	var exams []*Exam

	cacheKey := r.GetAllOpenedExamsCacheKey()
	val, err := r.cache.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &exams)
		return exams, nil
	}

	err = r.db.Where("is_open").Find(&exams).Error
	if err != nil {
		return nil, err
	}

	res, _ := json.Marshal(exams)
	r.cache.Set(context.Background(), cacheKey, res, r.cfg.CacheTTL)
	return exams, nil
}

func (r *repository) UpdateExam(exam *Exam) error {
	currentData, err := r.GetExamBySerial(exam.Serial)
	if err != nil {
		return err
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
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

		if err := tx.Model(&Exam{}).
			Where("serial = ?", exam.Serial).
			Update("allowed_duration_minutes", exam.AllowedDurationMinutes).
			Error; err != nil {
			return err
		}

		var hangingParticipants []*Participant
		err := r.db.Where("started_at IS NOT NULL AND ended_at IS NULL").Find(&hangingParticipants).Error
		if err != nil {
			return err
		}

		var updatedExam Exam
		err = tx.Where("serial = ?", exam.Serial).First(&updatedExam).Error
		if err != nil {
			return err
		}

		if len(hangingParticipants) == 0 {
			return nil
		}

		err = r.db.Model(&Participant{}).Where("started_at IS NOT NULL AND ended_at IS NULL").Update("ended_at", updatedExam.UpdatedAt).Error

		if err == nil {
			for i := range hangingParticipants {
				r.cache.Del(context.Background(), r.GetParticipantByIDCacheKey(hangingParticipants[i].ID))
				r.cache.Del(context.Background(), r.GetParticipantByExamIDAndNameCacheKey(updatedExam.ID, hangingParticipants[i].Name))
			}
		}

		return err
	})
	if err == nil {
		r.cache.Del(context.Background(), r.GetExamByIDCacheKey(currentData.ID))
		r.cache.Del(context.Background(), r.GetExamBySerialCacheKey(currentData.Serial))
		r.cache.Del(context.Background(), r.GetAllOpenedExamsCacheKey())
	}
	return err
}

func (r *repository) DeleteExamBySerial(serial string) error {
	currentData, err := r.GetExamBySerial(serial)
	if err != nil {
		return err
	}

	res := r.db.Model(&Exam{}).Where("serial = ?", serial).Delete(&Exam{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		log.Printf("[exam][repository][DeleteExamBySerial] error: %s", res.Error)
		return lib.ErrExamNotFound
	}

	r.cache.Del(context.Background(), r.GetExamByIDCacheKey(currentData.ID))
	r.cache.Del(context.Background(), r.GetExamBySerialCacheKey(currentData.Serial))
	r.cache.Del(context.Background(), r.GetAllOpenedExamsCacheKey())

	return nil
}

func (r *repository) GetExamBySerialCacheKey(serial string) string {
	return fmt.Sprintf("exam:serial:%s", serial)
}

func (r *repository) GetExamByIDCacheKey(id uint) string {
	return fmt.Sprintf("exam:id:%d", id)
}

func (r *repository) GetParticipantByIDCacheKey(id uint) string {
	return fmt.Sprintf("participant:id:%d", id)
}

func (r *repository) GetParticipantByExamIDAndNameCacheKey(examID uint, name string) string {
	return fmt.Sprintf("participant:examID:%d:name:%s", examID, name)
}

func (r *repository) GetAllOpenedExamsCacheKey() string {
	return "exam:allOpened"
}
