package submission

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/prajnapras19/project-form-exam-sman2/backend/config"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
	redis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository interface {
	GetSubmissionByParticipantIDAndQuestionID(participantID uint, questionID uint) (*Submission, error)
	SaveCacheObject(cacheObject *ExamSessionSubmissionCacheObject) error
	UpsertSubmissionInDB(cacheObject *ExamSessionSubmissionCacheObject) error
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

func (r *repository) GetSubmissionByParticipantIDAndQuestionID(participantID uint, questionID uint) (*Submission, error) {
	var submission Submission

	cacheKey := (&ExamSessionSubmissionCacheObject{
		ParticipantID: participantID,
		QuestionID:    questionID,
	}).GetKey()

	val, err := r.cache.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &submission)
		return &submission, nil
	}

	err = r.db.Where("participant_id = ? AND question_id AND not_archived", participantID, questionID).First(&submission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, lib.ErrSubmissionNotFound
		}
		return nil, err
	}

	res, _ := json.Marshal(&ExamSessionSubmissionCacheObject{
		ParticipantID: submission.ParticipantID,
		QuestionID:    submission.QuestionID,
		McqOptionID:   submission.McqOptionID,
		Timestamp:     submission.UpdatedAt,
	})
	r.cache.Set(context.Background(), cacheKey, res, r.cfg.CacheTTL)
	return &submission, nil
}

func (r *repository) SaveCacheObject(cacheObject *ExamSessionSubmissionCacheObject) error {
	res, _ := json.Marshal(cacheObject)
	return r.cache.Set(context.Background(), cacheObject.GetKey(), res, r.cfg.CacheTTL).Err()
}

func (r *repository) UpsertSubmissionInDB(cacheObject *ExamSessionSubmissionCacheObject) error {
	var submission Submission
	err := r.db.Where("participant_id = ? AND question_id AND not_archived", cacheObject.ParticipantID, cacheObject.QuestionID).First(&submission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r.db.Create(&Submission{
				BaseModel: lib.BaseModel{
					Model: gorm.Model{
						CreatedAt: cacheObject.Timestamp,
						UpdatedAt: cacheObject.Timestamp,
					},
				},
				ParticipantID: cacheObject.ParticipantID,
				QuestionID:    cacheObject.QuestionID,
				McqOptionID:   cacheObject.McqOptionID,
			}).Error
		}
		return err
	}

	if submission.UpdatedAt != cacheObject.Timestamp {
		return r.db.Where("participant_id = ? AND question_id AND not_archived", cacheObject.ParticipantID, cacheObject.QuestionID).Update("mcq_option_id", cacheObject.McqOptionID).Error
	}
	return nil
}
